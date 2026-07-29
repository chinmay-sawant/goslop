package engine

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sort"

	"golang.org/x/sync/errgroup"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

// Analyzer is the language-agnostic static analysis orchestrator.
type Analyzer struct {
	registry *Registry
	ctx      *core.ScanContext
	workers  int
	walkOpts *WalkOptions
}

// AnalyzerBuilder configures an Analyzer.
type AnalyzerBuilder struct {
	registry *Registry
	ctx      *core.ScanContext
	workers  int
	walkOpts *WalkOptions
}

// NewAnalyzerBuilder starts a fluent builder.
func NewAnalyzerBuilder() *AnalyzerBuilder {
	return &AnalyzerBuilder{}
}

// Registry sets the detector/language registry.
func (b *AnalyzerBuilder) Registry(r *Registry) *AnalyzerBuilder {
	b.registry = r
	return b
}

// ScanContext sets rule filters and scan options.
func (b *AnalyzerBuilder) ScanContext(ctx *core.ScanContext) *AnalyzerBuilder {
	b.ctx = ctx
	return b
}

// Workers sets the parallel file worker count (0 = GOMAXPROCS).
func (b *AnalyzerBuilder) Workers(n int) *AnalyzerBuilder {
	b.workers = n
	return b
}

// WalkOptions sets filesystem discovery options.
func (b *AnalyzerBuilder) WalkOptions(opts WalkOptions) *AnalyzerBuilder {
	b.walkOpts = &opts
	return b
}

// Build constructs the Analyzer.
func (b *AnalyzerBuilder) Build() *Analyzer {
	ctx := b.ctx
	if ctx == nil {
		ctx = core.DefaultScanContext()
	}
	workers := b.workers
	if workers == 0 {
		workers = ctx.Workers
	}
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	return &Analyzer{
		registry: b.registry,
		ctx:      ctx,
		workers:  workers,
		walkOpts: b.walkOpts,
	}
}

// NewAnalyzer is a convenience constructor.
func NewAnalyzer(ctx *core.ScanContext, registry *Registry) *Analyzer {
	return NewAnalyzerBuilder().Registry(registry).ScanContext(ctx).Build()
}

// ScanContext returns the scan policy used by this analyzer.
func (a *Analyzer) ScanContext() *core.ScanContext {
	return a.ctx
}

// AnalyzePaths walks paths, runs matching detectors in parallel, finalizes, and returns results.
func (a *Analyzer) AnalyzePaths(paths []string) (*AnalysisResult, error) {
	if a == nil {
		return nil, fmt.Errorf("nil analyzer")
	}
	if a.registry == nil {
		return nil, fmt.Errorf("analyzer has no registry")
	}
	if len(paths) == 0 {
		paths = []string{"."}
	}
	ctx := a.ctx
	if ctx == nil {
		ctx = core.DefaultScanContext()
	}

	for _, det := range a.registry.Detectors() {
		det.BeginScan(ctx)
	}
	defer func() {
		for _, det := range a.registry.Detectors() {
			det.EndScan()
		}
	}()

	roots := distinctRoots(paths)
	for _, plugin := range a.registry.Plugins() {
		plugin.PrepareProject(ctx, roots)
	}

	walkOpts := DefaultWalkOptions()
	if a.walkOpts != nil {
		walkOpts = *a.walkOpts
	}
	if ctx.IncludeTests {
		walkOpts.IncludeTests = true
	}

	entries, err := CollectFiles(paths, walkOpts, a.registry.ExtensionMap())
	if err != nil {
		return nil, err
	}

	stats := &ScanStats{
		DetectorsLoaded: a.registry.DetectorCount(),
	}

	results := make([]fileResult, len(entries))
	g := new(errgroup.Group)
	g.SetLimit(a.workers)

	for i, entry := range entries {
		i, entry := i, entry
		g.Go(func() error {
			results[i] = a.scanOne(ctx, entry)
			return nil
		})
	}
	_ = g.Wait()

	var (
		findings    []rules.Finding
		scanErrors  []ScanError
		sourceCache map[string]string
	)
	if ctx.RetainSources {
		sourceCache = make(map[string]string)
	}

	for _, fr := range results {
		if fr.err != nil {
			if se, ok := fr.err.(*ScanError); ok {
				scanErrors = append(scanErrors, *se)
			} else {
				scanErrors = append(scanErrors, ScanError{
					Path:    fr.path,
					Kind:    ScanErrorEngine,
					Message: fr.err.Error(),
				})
			}
			stats.FilesErrored++
			continue
		}
		stats.FilesScanned++
		stats.BytesScanned += fr.bytes
		findings = append(findings, fr.findings...)
		if sourceCache != nil && fr.source != "" {
			sourceCache[fr.path] = fr.source
		}
	}

	for _, det := range a.registry.Detectors() {
		if !anyRuleAllowed(ctx, det.RuleIDs()) {
			continue
		}
		det.Finalize(ctx, &findings)
	}

	findings = filterFindings(ctx, findings)

	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		if findings[i].Column != findings[j].Column {
			return findings[i].Column < findings[j].Column
		}
		return findings[i].RuleID < findings[j].RuleID
	})

	stats.FindingsTotal = len(findings)

	return &AnalysisResult{
		Findings:    findings,
		Errors:      scanErrors,
		Stats:       stats,
		SourceCache: sourceCache,
	}, nil
}

type fileResult struct {
	findings []rules.Finding
	err      error
	bytes    int64
	source   string
	path     string
}

func (a *Analyzer) scanOne(ctx *core.ScanContext, entry ScanEntry) fileResult {
	source, err := ReadUTF8(entry.Path)
	if err != nil {
		return fileResult{err: err, path: entry.Path}
	}

	display := entry.Path
	if rel, rerr := filepath.Rel(".", entry.Path); rerr == nil && !filepath.IsAbs(rel) {
		display = rel
	}

	unit, perr := a.buildUnit(entry, display, source)
	if perr != nil {
		return fileResult{err: perr, path: display, bytes: int64(len(source))}
	}
	defer closeUnitTree(unit)

	var out []rules.Finding
	for _, idx := range a.registry.DetectorIndices(entry.Language) {
		det := a.registry.Detector(idx)
		if det == nil {
			continue
		}
		if !anyRuleAllowed(ctx, det.RuleIDs()) {
			continue
		}
		det.Run(ctx, unit, &out)
		det.AccumulateState(ctx, unit)
	}

	return fileResult{
		findings: out,
		bytes:    int64(len(source)),
		source:   source,
		path:     display,
	}
}

// closeUnitTree frees an attached tree-sitter tree if the unit owns one.
func closeUnitTree(unit *core.ParsedUnit) {
	if unit == nil || unit.Tree == nil {
		return
	}
	if c, ok := unit.Tree.(interface{ Close() }); ok {
		c.Close()
	}
	unit.Tree = nil
}

func (a *Analyzer) buildUnit(entry ScanEntry, display, source string) (*core.ParsedUnit, error) {
	plugin := a.registry.PluginForID(entry.Language)
	if plugin != nil {
		unit, err := plugin.ParseSource(entry.Path, source)
		if err != nil {
			return nil, &ScanError{
				Path:    entry.Path,
				Kind:    ScanErrorParse,
				Message: err.Error(),
			}
		}
		if unit != nil {
			if unit.DisplayPath == "" {
				unit.DisplayPath = display
			}
			if unit.Path == "" {
				unit.Path = entry.Path
			}
			return unit, nil
		}
	}
	unit := core.NewParsedUnit(entry.Language, entry.Path, source)
	unit.DisplayPath = display
	return unit, nil
}

func anyRuleAllowed(ctx *core.ScanContext, ruleIDs []string) bool {
	for _, id := range ruleIDs {
		if ctx.Allows(id) {
			return true
		}
	}
	return false
}

func filterFindings(ctx *core.ScanContext, findings []rules.Finding) []rules.Finding {
	if len(findings) == 0 {
		return findings
	}
	out := findings[:0]
	for i := range findings {
		f := findings[i]
		if !ctx.Allows(f.RuleID) {
			continue
		}
		out = append(out, f)
	}
	return out
}

func distinctRoots(paths []string) []string {
	seen := make(map[string]struct{})
	var roots []string
	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			abs = p
		}
		root := abs
		if filepath.Ext(abs) != "" {
			root = filepath.Dir(abs)
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		roots = append(roots, root)
	}
	if len(roots) == 0 {
		return []string{"."}
	}
	return roots
}

package engine

import (
	"errors"
	"path/filepath"
	"runtime"
	"sort"
	"sync/atomic"

	"golang.org/x/sync/errgroup"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/engine/baseline"
	"github.com/chinmay/codehound/internal/engine/cache"
	"github.com/chinmay/codehound/internal/engine/ignore"
	"github.com/chinmay/codehound/internal/rules"
)

// Analyzer is the language-agnostic static analysis orchestrator.
type Analyzer struct {
	registry *Registry
	ctx      *core.ScanContext
	workers  int
	walkOpts *WalkOptions
	cache    *cache.Store
	baseline *baseline.Baseline
	// projectRoot is used for cache-relative paths (first scan root).
	projectRoot string
}

// AnalyzerBuilder configures an Analyzer.
type AnalyzerBuilder struct {
	registry    *Registry
	ctx         *core.ScanContext
	workers     int
	walkOpts    *WalkOptions
	cache       *cache.Store
	baseline    *baseline.Baseline
	projectRoot string
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

// Cache sets the incremental analysis cache (nil disables).
func (b *AnalyzerBuilder) Cache(store *cache.Store) *AnalyzerBuilder {
	b.cache = store
	return b
}

// Baseline sets the baseline filter (nil disables).
func (b *AnalyzerBuilder) Baseline(bl *baseline.Baseline) *AnalyzerBuilder {
	b.baseline = bl
	return b
}

// ProjectRoot sets the root used for cache-relative paths.
func (b *AnalyzerBuilder) ProjectRoot(root string) *AnalyzerBuilder {
	b.projectRoot = root
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
		registry:    b.registry,
		ctx:         ctx,
		workers:     workers,
		walkOpts:    b.walkOpts,
		cache:       b.cache,
		baseline:    b.baseline,
		projectRoot: b.projectRoot,
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

// SetCache attaches or replaces the incremental cache.
func (a *Analyzer) SetCache(store *cache.Store) {
	if a != nil {
		a.cache = store
	}
}

// SetBaseline attaches or replaces the baseline filter.
func (a *Analyzer) SetBaseline(bl *baseline.Baseline) {
	if a != nil {
		a.baseline = bl
	}
}

// AnalyzePaths walks paths, runs matching detectors in parallel, finalizes, and returns results.
func (a *Analyzer) AnalyzePaths(paths []string) (*AnalysisResult, error) {
	if a == nil {
		return nil, errors.New("nil analyzer")
	}
	if a.registry == nil {
		return nil, errors.New("analyzer has no registry")
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

	// Project root for relative cache keys: first directory root or cwd.
	projectRoot := a.projectRoot
	if projectRoot == "" {
		if len(roots) > 0 {
			projectRoot = roots[0]
		} else {
			projectRoot = "."
		}
	}

	store := a.cache
	if ctx.NoCache {
		store = nil
	}
	if store != nil {
		store.EnsureRuleConfigHash(ctx.RuleConfigFingerprint())
	}

	stats := &ScanStats{
		DetectorsLoaded: a.registry.DetectorCount(),
	}

	results := make([]fileResult, len(entries))
	g := new(errgroup.Group)
	g.SetLimit(a.workers)

	var cacheHits, cacheMisses atomic.Int64
	var suppressedTotal atomic.Int64

	scannedFiles := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		rel := relPath(projectRoot, e.Path)
		scannedFiles[cache.NormalizePath(rel)] = struct{}{}
	}

	for i, entry := range entries {
		i, entry := i, entry
		g.Go(func() error {
			fr := a.scanOne(ctx, entry, projectRoot, store, &cacheHits, &cacheMisses, &suppressedTotal)
			results[i] = fr
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
			var se *ScanError
			if errors.As(fr.err, &se) {
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

	// Baseline filter (after finalize + only/skip).
	baselinedCount := 0
	if a.baseline != nil && !ctx.NoBaseline {
		findings, baselinedCount = a.baseline.Filter(findings, ctx.ShowBaselined)
	}

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

	// Cache prune + flush.
	if store != nil {
		_, _ = store.Prune(scannedFiles)
		_ = store.Flush()
	}

	suppressed := int(suppressedTotal.Load())
	stats.FindingsTotal = len(findings)
	stats.FindingsSuppressed = suppressed
	stats.FindingsBaselined = baselinedCount
	stats.CacheHits = int(cacheHits.Load())
	stats.CacheMisses = int(cacheMisses.Load())

	return &AnalysisResult{
		Findings:        findings,
		Errors:          scanErrors,
		Stats:           stats,
		SourceCache:     sourceCache,
		SuppressedCount: suppressed,
		BaselinedCount:  baselinedCount,
	}, nil
}

type fileResult struct {
	findings []rules.Finding
	err      error
	bytes    int64
	source   string
	path     string
}

func (a *Analyzer) scanOne(
	ctx *core.ScanContext,
	entry ScanEntry,
	projectRoot string,
	store *cache.Store,
	hits, misses *atomic.Int64,
	suppressedTotal *atomic.Int64,
) fileResult {
	source, err := ReadUTF8(entry.Path)
	if err != nil {
		return fileResult{err: err, path: entry.Path}
	}

	display := entry.Path
	if rel, rerr := filepath.Rel(".", entry.Path); rerr == nil && !filepath.IsAbs(rel) {
		display = rel
	}
	cacheRel := cache.NormalizePath(relPath(projectRoot, entry.Path))
	contentHash := cache.ContentHash(source)

	// Warm cache hit: skip parse + detect.
	if store != nil && store.ShouldCacheBytes(int64(len(source))) {
		if kind, entryData := store.Lookup(cacheRel, contentHash); kind == cache.LookupHit && entryData != nil {
			hits.Add(1)
			findings := append([]rules.Finding(nil), entryData.Findings...)
			// Re-apply display path? findings already have paths from prior run.
			return fileResult{
				findings: findings,
				bytes:    int64(len(source)),
				source:   source,
				path:     display,
			}
		}
		misses.Add(1)
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

	// Inline / file ignore suppressions.
	var suppressed int
	out, suppressed = ignore.Apply(source, out, ignore.ApplyOptions{ShowIgnored: ctx.ShowIgnored})
	if suppressed > 0 {
		suppressedTotal.Add(int64(suppressed))
	}

	if store != nil && store.ShouldCacheBytes(int64(len(source))) {
		_ = store.Put(cacheRel, contentHash, out, suppressed)
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
			// May be a file; only use Dir when it is not a directory.
			// Heuristic: if extension present, prefer parent (works for .go).
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

func relPath(root, path string) string {
	if root == "" {
		root = "."
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return rel
}

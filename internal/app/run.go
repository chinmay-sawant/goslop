// Package app orchestrates the goslop CLI.
package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chinmay-sawant/goslop/internal/cli"
	"github.com/chinmay-sawant/goslop/internal/config"
	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/engine/baseline"
	"github.com/chinmay-sawant/goslop/internal/engine/cache"
	"github.com/chinmay-sawant/goslop/internal/export"
	"github.com/chinmay-sawant/goslop/internal/reporting"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Run is the CLI entry used by cmd/goslop.
func Run(args []string) error {
	return run(args, os.Stdout, os.Stderr)
}

func run(args []string, stdout, stderr io.Writer) error {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	opts, err := cli.ParseWithOutput(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return &ExitCodeError{Code: ExitConfig, Err: err}
	}

	if opts.Command == "init" {
		return runInit(stdout)
	}
	if opts.Version {
		_, _ = fmt.Fprintln(stdout, Version)
		return nil
	}
	// Load goslop.toml early so --list-rules / --explain can honor languages.
	merged, merr := config.LoadAndMerge(config.MergeInput{
		Only:           opts.Only,
		Skip:           opts.Skip,
		IncludeTests:   opts.IncludeTests,
		NoCache:        opts.NoCache,
		CacheDir:       opts.CacheDir,
		NoBaseline:     opts.NoBaseline,
		BaselineFile:   opts.BaselineFile,
		Taint:          opts.Taint,
		NoTaint:        opts.NoTaint,
		TaintShowPaths: opts.TaintShowPaths,
		NoFail:         opts.NoFail,
		ConfigPath:     opts.ConfigPath,
		Paths:          opts.Paths,
	})
	if merr != nil {
		return &ExitCodeError{Code: ExitConfig, Err: merr}
	}

	// Compose a registry from built-in plugins for the enabled languages.
	// Default (unset languages) is Go-only via config.DefaultEnabledLanguages.
	// Unknown language tokens are rejected at config load; missing built-ins error here.
	reg, rerr := engine.NewRegistryWithLanguages(merged.Languages...)
	if rerr != nil {
		return &ExitCodeError{Code: ExitConfig, Err: rerr}
	}

	if opts.ListRules {
		return listRules(stdout, reg)
	}
	if opts.ExplainRule != "" {
		return explainRule(stdout, reg, opts.ExplainRule)
	}

	profile, ok := core.ParseProfile(opts.Profile)
	if !ok {
		return &ExitCodeError{
			Code: ExitConfig,
			Err:  fmt.Errorf("unknown profile %q", opts.Profile),
		}
	}

	plan := resolveScanPlan(profile, opts, merged)
	scope := engine.ResolveScanScope(opts.Paths)
	ctx := plan.context

	// Cache open / rebuild / prune.
	var store *cache.Store
	if !plan.cacheDisabled {
		cacheDir := plan.cacheDir
		openOpts := plan.cacheOptions
		if opts.RebuildCache {
			store, err = cache.Rebuild(cacheDir, openOpts)
			if err != nil {
				if opts.PruneCache {
					return &ExitCodeError{Code: ExitInternal, Err: fmt.Errorf("rebuild cache before pruning: %w", err)}
				}
				_, _ = fmt.Fprintf(stderr, "warning: could not rebuild cache: %v\n", err)
				store = nil
			}
		} else {
			store, err = cache.Open(cacheDir, openOpts)
			if err != nil {
				if opts.PruneCache {
					return &ExitCodeError{Code: ExitInternal, Err: fmt.Errorf("open cache for pruning: %w", err)}
				}
				_, _ = fmt.Fprintf(stderr, "warning: could not open cache: %v\n", err)
				store = nil
			}
		}
	}

	if opts.PruneCache {
		return runPruneCache(opts, reg, store, scope, stdout, stderr)
	}

	// Baseline discovery / load.
	var bl *baseline.Baseline
	if !plan.baselineDisabled {
		path := plan.baselineFile
		if path == "" {
			path = baseline.Discover(".")
		}
		if path != "" {
			if loaded, lerr := baseline.Load(path); lerr != nil {
				_, _ = fmt.Fprintf(stderr, "warning: could not load baseline %s: %v\n", path, lerr)
			} else {
				bl = loaded
			}
		}
	}

	analyzer := engine.NewAnalyzerBuilder().
		Registry(reg).
		ScanContext(ctx).
		WalkOptions(plan.walkOptions).
		Cache(store).
		Baseline(bl).
		ProjectRoot(scope.ProjectRoot()).
		Build()

	t0 := time.Now()
	res, err := analyzer.AnalyzePaths(opts.Paths)
	wall := time.Since(t0)
	if err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}
	if err := flushCache(store); err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}

	findings := res.Findings
	if findings == nil {
		findings = []rules.Finding{}
	}

	// When --no-cache, treat every scanned file as a cache miss for summary
	// parity with Rust ("0 hits, N misses (full re-analysis)").
	if opts.NoCache && res.Stats != nil && res.Stats.CacheHits+res.Stats.CacheMisses == 0 {
		res.Stats.CacheMisses = res.Stats.FilesScanned
	}

	// Product-style scan summary (Rust make run / --no-terminal).
	// Always on stderr so JSON/SARIF stdout stays pure for tooling.
	printScanSummary(stderr, res, findings, wall)
	reportScanDiagnostics(stderr, res.Diagnostics)
	if len(res.Errors) > 0 {
		reportScanErrors(stderr, res.Errors)
		return &ExitCodeError{
			Code: ExitInternal,
			Err:  fmt.Errorf("analysis incomplete: %d file(s) could not be analyzed", len(res.Errors)),
		}
	}

	// --no-terminal: summary only (no per-finding text dump). Still emit
	// machine formats when explicitly requested (json/sarif).
	if !(opts.NoTerminal && opts.Format == cli.FormatText) {
		rep, rerr := reporting.New(string(opts.Format))
		if rerr != nil {
			return &ExitCodeError{Code: ExitConfig, Err: rerr}
		}
		switch r := rep.(type) {
		case reporting.JSONReporter:
			r.Version = Version
			rep = r
		case reporting.SARIFReporter:
			r.Version = Version
			rep = r
		}
		if err := rep.Write(findings, stdout); err != nil {
			return &ExitCodeError{Code: ExitInternal, Err: err}
		}
	}

	if opts.ExportContext || opts.ExportChunks {
		sum, xerr := export.ExportFindings(findings, export.Options{
			ExportContext:    opts.ExportContext,
			ExportChunks:     opts.ExportChunks,
			ChunkSize:        opts.ChunkSize,
			ContextOutputDir: opts.ContextDir,
			ChunksOutputDir:  opts.ChunksDir,
			WholeFunction:    merged.ExportWholeFunction, // nil → default true
		}, res.SourceCache)
		if xerr != nil {
			return &ExitCodeError{Code: ExitInternal, Err: xerr}
		}
		if opts.ExportContext {
			_, _ = fmt.Fprintf(stderr, "exported %d context file(s) to %s",
				sum.ContextFilesWritten, export.Options{
					ContextOutputDir: opts.ContextDir,
				}.Normalize().ContextOutputDir)
			if opts.ExportChunks {
				_, _ = fmt.Fprintf(stderr, "; exported %d chunk file(s) to %s",
					sum.ChunkFilesWritten, export.Options{
						ChunksOutputDir: opts.ChunksDir,
					}.Normalize().ChunksOutputDir)
			}
			_, _ = fmt.Fprintln(stderr)
		} else if opts.ExportChunks {
			_, _ = fmt.Fprintf(stderr, "exported %d chunk file(s) to %s\n",
				sum.ChunkFilesWritten, export.Options{
					ChunksOutputDir: opts.ChunksDir,
				}.Normalize().ChunksOutputDir)
		}
	}

	if !opts.NoFail && res.ShouldFail(ctx.FailPolicy) {
		return &ExitCodeError{Code: ExitFailing}
	}
	return nil
}

func runPruneCache(opts *cli.Options, reg *engine.Registry, store *cache.Store, scope engine.ScanScope, stdout, stderr io.Writer) error {
	if store == nil {
		_, _ = fmt.Fprintln(stdout, "cache disabled; nothing to prune")
		return nil
	}
	walkOpts := engine.DefaultWalkOptions()
	if opts.IncludeTests {
		walkOpts.IncludeTests = true
	}
	entries, _, err := engine.CollectFiles(opts.Paths, walkOpts, reg.ExtensionMap())
	if err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}
	scanned := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		scanned[scope.CacheKey(e.Path)] = struct{}{}
	}
	pruned, err := store.Prune(scanned)
	if err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}
	orphaned, err := store.CleanOrphans()
	if err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}
	if err := store.Flush(); err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}
	if pruned > 0 || orphaned > 0 {
		_, _ = fmt.Fprintf(stdout, "Pruned %d stale manifest entries and removed %d orphaned cache files from %s\n",
			pruned, orphaned, store.Dir())
	} else {
		_, _ = fmt.Fprintf(stdout, "Cache at %s is clean (0 stale entries, 0 orphans)\n", store.Dir())
	}
	_ = stderr
	return nil
}

// flushCache reconciles orphaned cache files after analyzer-owned prune and
// flush have completed successfully.
func flushCache(store *cache.Store) error {
	if store == nil {
		return nil
	}
	if _, err := store.CleanOrphans(); err != nil {
		return fmt.Errorf("prune cache: %w", err)
	}
	return nil
}

func printScanSummary(w io.Writer, res *engine.AnalysisResult, findings []rules.Finding, wall time.Duration) {
	if w == nil || res == nil || res.Stats == nil {
		return
	}
	st := res.Stats
	// scanned 78 files (28120 lines) in 479.5ms
	_, _ = fmt.Fprintf(w, "scanned %d files (%d lines) in %s\n",
		st.FilesScanned, st.LinesScanned, formatWall(wall))
	// cache line (always when we have hits/misses, including full re-analysis)
	if st.CacheHits+st.CacheMisses > 0 {
		suffix := ""
		switch {
		case st.CacheHits > 0 && st.CacheMisses == 0:
			suffix = " (results from cache; not re-analyzed)"
		case st.CacheHits == 0:
			suffix = " (full re-analysis)"
		}
		_, _ = fmt.Fprintf(w, "  cache: %d hits, %d misses%s\n", st.CacheHits, st.CacheMisses, suffix)
	}
	if st.FilesSkipped > 0 {
		_, _ = fmt.Fprintf(w, "  skipped %d files\n", st.FilesSkipped)
	}
	if len(res.Errors) > 0 {
		_, _ = fmt.Fprintf(w, "  incomplete: %d file(s) could not be analyzed\n", len(res.Errors))
	}
	if len(findings) == 0 {
		if len(res.Errors) > 0 {
			return
		}
		_, _ = fmt.Fprintln(w, "no slop detected")
		return
	}
	_, _ = fmt.Fprintf(w, "%d findings\n", len(findings))
	// severity histogram (Rust order: high, info, low, medium)
	var high, med, low, info int
	counts := map[string]int{}
	exampleCount := 0
	for _, f := range findings {
		counts[f.RuleID]++
		if isExampleDemoFinding(f) {
			exampleCount++
		}
		switch f.Severity {
		case rules.SeverityHigh, rules.SeverityCritical:
			high++
		case rules.SeverityMedium:
			med++
		case rules.SeverityLow:
			low++
		case rules.SeverityInfo:
			info++
		default:
			info++
		}
	}
	_, _ = fmt.Fprintf(w, "  severity: %d high, %d info, %d low, %d medium\n", high, info, low, med)
	// top 5 rules
	type pair struct {
		id string
		n  int
	}
	top := make([]pair, 0, len(counts))
	for id, n := range counts {
		top = append(top, pair{id, n})
	}
	sort.Slice(top, func(i, j int) bool {
		if top[i].n != top[j].n {
			return top[i].n > top[j].n
		}
		return top[i].id < top[j].id
	})
	if len(top) > 5 {
		top = top[:5]
	}
	parts := make([]string, 0, len(top))
	for _, p := range top {
		parts = append(parts, fmt.Sprintf("%s ×%d", p.id, p.n))
	}
	_, _ = fmt.Fprintf(w, "  top rules: %s\n", strings.Join(parts, ", "))
	if exampleCount > 0 {
		_, _ = fmt.Fprintf(w, "  example findings: %d (of %d total)\n", exampleCount, len(findings))
	}
}

func reportScanErrors(w io.Writer, scanErrors []engine.ScanError) {
	for _, scanErr := range scanErrors {
		_, _ = fmt.Fprintf(w, "analysis error [%s]: %s\n", scanErr.Kind, scanErr.Error())
	}
}

func reportScanDiagnostics(w io.Writer, diagnostics []engine.ScanDiagnostic) {
	for _, diagnostic := range diagnostics {
		_, _ = fmt.Fprintf(w, "analysis warning [%s]: %s\n", diagnostic.Kind, diagnostic.String())
	}
}

// formatWall matches Rust product summary style (e.g. 479.5ms, 1.2s).
func formatWall(d time.Duration) string {
	if d < time.Second {
		ms := float64(d) / float64(time.Millisecond)
		return fmt.Sprintf("%.1fms", ms)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// isExampleDemoFinding reports findings under sample/demo paths (Rust
// EXAMPLE_PATH_COMPONENTS: examples, example, sampledata, samples).
func isExampleDemoFinding(f rules.Finding) bool {
	for _, t := range f.Tags {
		if t == "example" {
			return true
		}
	}
	p := strings.ToLower(filepath.ToSlash(f.File))
	for _, part := range strings.Split(p, "/") {
		switch part {
		case "examples", "example", "sampledata", "samples":
			return true
		}
	}
	return false
}

func listRules(w io.Writer, reg *engine.Registry) error {
	if reg == nil {
		reg = engine.DefaultRegistry()
	}
	ids := reg.AllRuleIDs()
	if len(ids) == 0 {
		_, _ = fmt.Fprintln(w, "no rules registered")
		return nil
	}
	// Stable sort already in AllRuleIDs; print once per id via detectors.
	seen := map[string]struct{}{}
	for _, d := range reg.Detectors() {
		for _, id := range d.RuleIDs() {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			mat := rules.InferMaturity(id)
			title := ""
			if meta := d.MetadataFor(id); meta != nil {
				title = meta.Title
				mat = meta.EffectiveMaturity()
			}
			if title != "" {
				_, _ = fmt.Fprintf(w, "[%s]\t%s\t%s\n", mat, id, title)
			} else {
				_, _ = fmt.Fprintf(w, "[%s]\t%s\n", mat, id)
			}
		}
	}
	return nil
}

func explainRule(w io.Writer, reg *engine.Registry, ruleID string) error {
	if reg == nil {
		reg = engine.DefaultRegistry()
	}
	for _, d := range reg.Detectors() {
		meta := d.MetadataFor(ruleID)
		if meta == nil {
			continue
		}
		mat := meta.EffectiveMaturity()
		_, _ = fmt.Fprintf(w, "rule: %s\n", meta.ID)
		_, _ = fmt.Fprintf(w, "title: %s\n", meta.Title)
		_, _ = fmt.Fprintf(w, "severity: %s\n", meta.Severity.String())
		_, _ = fmt.Fprintf(w, "pack: %s\n", meta.Pack.String())
		_, _ = fmt.Fprintf(w, "maturity: %s\n", mat)
		if meta.QuarantineReason != "" {
			_, _ = fmt.Fprintf(w, "quarantine: %s\n", meta.QuarantineReason)
		}
		if meta.Description != "" {
			_, _ = fmt.Fprintf(w, "description: %s\n", meta.Description)
		}
		if meta.Fix != "" {
			_, _ = fmt.Fprintf(w, "fix: %s\n", meta.Fix)
		}
		// Pack eligibility hints
		rec := core.BuildScanContext(core.ProfileRecommended, nil, nil)
		sec := core.BuildScanContext(core.ProfileSecurity, nil, nil)
		_, _ = fmt.Fprintf(w, "recommended_pack: %v\n", rec.Allows(ruleID))
		_, _ = fmt.Fprintf(w, "security_pack: %v\n", sec.Allows(ruleID))
		_, _ = fmt.Fprintf(w, "all_pack: true\n")
		return nil
	}
	return &ExitCodeError{
		Code: ExitConfig,
		Err:  fmt.Errorf("unknown rule %q", ruleID),
	}
}

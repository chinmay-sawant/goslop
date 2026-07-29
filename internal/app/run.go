// Package app orchestrates the CodeHound CLI.
package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chinmay/codehound/internal/cli"
	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/engine"
	"github.com/chinmay/codehound/internal/engine/baseline"
	"github.com/chinmay/codehound/internal/engine/cache"
	"github.com/chinmay/codehound/internal/reporting"
	"github.com/chinmay/codehound/internal/rules"
)

// Run is the CLI entry used by cmd/codehound.
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
		return runInit()
	}
	if opts.Version {
		_, _ = fmt.Fprintln(stdout, Version)
		return nil
	}
	if opts.ListRules {
		return listRules(stdout)
	}

	profile, ok := core.ParseProfile(opts.Profile)
	if !ok {
		return &ExitCodeError{
			Code: ExitConfig,
			Err:  fmt.Errorf("unknown profile %q", opts.Profile),
		}
	}

	ctx := core.NewScanContext(profile, opts.Only, opts.Skip)
	ctx.IncludeTests = opts.IncludeTests
	ctx.NoCache = opts.NoCache
	ctx.ShowIgnored = opts.ShowIgnored
	ctx.ShowBaselined = opts.ShowBaselined
	ctx.NoBaseline = opts.NoBaseline
	if opts.NoTaint {
		ctx.TaintEnabled = false
	} else if opts.Taint {
		ctx.TaintEnabled = true
	}
	if opts.TaintDepth > 0 {
		ctx.TaintMaxDepth = opts.TaintDepth
	}
	ctx.TaintShowPaths = opts.TaintShowPaths

	reg := engine.DefaultRegistry()

	// Cache open / rebuild / prune.
	var store *cache.Store
	if !opts.NoCache {
		cacheDir := opts.CacheDir
		if cacheDir == "" {
			cacheDir = cache.DEFAULT_CACHE_DIR
		}
		openOpts := cache.OpenOptions{
			MaxSizeMB:     500,
			MaxFileSizeMB: 4,
			ToolVersion:   Version,
		}
		if opts.RebuildCache {
			store, err = cache.Rebuild(cacheDir, openOpts)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "warning: could not rebuild cache: %v\n", err)
				store = nil
			}
		} else {
			store, err = cache.Open(cacheDir, openOpts)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "warning: could not open cache: %v\n", err)
				store = nil
			}
		}
	}

	if opts.PruneCache {
		return runPruneCache(opts, reg, store, stdout, stderr)
	}

	// Baseline discovery / load.
	var bl *baseline.Baseline
	if !opts.NoBaseline {
		path := opts.BaselineFile
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

	// Project root for cache-relative keys: first path if dir, else parent.
	projectRoot := "."
	if len(opts.Paths) > 0 {
		p := opts.Paths[0]
		if fi, serr := os.Stat(p); serr == nil && fi.IsDir() {
			if abs, aerr := filepath.Abs(p); aerr == nil {
				projectRoot = abs
			} else {
				projectRoot = p
			}
		} else if abs, aerr := filepath.Abs(filepath.Dir(p)); aerr == nil {
			projectRoot = abs
		}
	}

	analyzer := engine.NewAnalyzerBuilder().
		Registry(reg).
		ScanContext(ctx).
		Cache(store).
		Baseline(bl).
		ProjectRoot(projectRoot).
		Build()

	res, err := analyzer.AnalyzePaths(opts.Paths)
	if err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}

	findings := res.Findings
	if findings == nil {
		findings = []rules.Finding{}
	}

	rep, err := reporting.New(string(opts.Format))
	if err != nil {
		return &ExitCodeError{Code: ExitConfig, Err: err}
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

	if res.ShouldFail(ctx.FailPolicy) {
		return &ExitCodeError{Code: ExitFailing}
	}
	return nil
}

func runPruneCache(opts *cli.Options, reg *engine.Registry, store *cache.Store, stdout, stderr io.Writer) error {
	if store == nil {
		_, _ = fmt.Fprintln(stdout, "cache disabled; nothing to prune")
		return nil
	}
	walkOpts := engine.DefaultWalkOptions()
	if opts.IncludeTests {
		walkOpts.IncludeTests = true
	}
	entries, err := engine.CollectFiles(opts.Paths, walkOpts, reg.ExtensionMap())
	if err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}
	scanned := make(map[string]struct{}, len(entries))
	// Prune uses project-relative keys; normalize to basenames-relative as stored.
	// Use path as absolute normalized for best-effort match with scan-time keys.
	projectRoot := "."
	if len(opts.Paths) > 0 {
		if abs, aerr := filepath.Abs(opts.Paths[0]); aerr == nil {
			if fi, serr := os.Stat(abs); serr == nil && !fi.IsDir() {
				projectRoot = filepath.Dir(abs)
			} else {
				projectRoot = abs
			}
		}
	}
	for _, e := range entries {
		rel, rerr := filepath.Rel(projectRoot, e.Path)
		if rerr != nil {
			rel = e.Path
		}
		scanned[cache.NormalizePath(rel)] = struct{}{}
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

func listRules(w io.Writer) error {
	reg := engine.DefaultRegistry()
	ids := reg.AllRuleIDs()
	if len(ids) == 0 {
		_, _ = fmt.Fprintln(w, "no rules registered")
		return nil
	}
	for _, d := range reg.Detectors() {
		for _, id := range d.RuleIDs() {
			title := ""
			if meta := d.MetadataFor(id); meta != nil {
				title = meta.Title
			}
			if title != "" {
				_, _ = fmt.Fprintf(w, "%s\t%s\n", id, title)
			} else {
				_, _ = fmt.Fprintln(w, id)
			}
		}
	}
	return nil
}

package app

import (
	"github.com/chinmay-sawant/goslop/internal/cli"
	"github.com/chinmay-sawant/goslop/internal/config"
	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/engine/cache"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// scanPlan is the resolved, application-facing scan configuration. It keeps
// the CLI/config precedence boundary in one place before engine construction.
type scanPlan struct {
	context          *core.ScanContext
	walkOptions      engine.WalkOptions
	cacheDir         string
	cacheOptions     cache.OpenOptions
	cacheDisabled    bool
	baselineDisabled bool
	baselineFile     string
}

func resolveScanPlan(profile core.ScanProfile, opts *cli.Options, merged *config.Merged) scanPlan {
	ctx := core.NewScanContext(profile, merged.Only, merged.Skip)
	ctx.IncludeTests = merged.IncludeTests
	ctx.NoCache = merged.NoCache
	ctx.ShowIgnored = opts.ShowIgnored
	ctx.ShowBaselined = opts.ShowBaselined
	ctx.NoBaseline = merged.NoBaseline
	if merged.FailPolicy != nil {
		ctx.FailPolicy = *merged.FailPolicy
	}
	if merged.BadPracticesEnabled != nil {
		ctx.BadPracticesEnabled = *merged.BadPracticesEnabled
	}
	if merged.BPSeverity != nil {
		ctx.BadPracticeSeverity = merged.BPSeverity
	}
	if len(merged.SeverityOverrides) > 0 {
		ctx.SeverityOverrides = make(map[string]rules.Severity, len(merged.SeverityOverrides))
		for id, severity := range merged.SeverityOverrides {
			ctx.SeverityOverrides[id] = severity
		}
	}
	if merged.NoTaint {
		ctx.TaintEnabled = false
	} else if merged.Taint {
		ctx.TaintEnabled = true
	}
	if opts.TaintDepth > 0 {
		ctx.TaintMaxDepth = opts.TaintDepth
	}
	ctx.TaintShowPaths = merged.TaintShowPaths
	ctx.RetainSources = opts.ExportContext || opts.ExportChunks

	cacheDir := merged.CacheDir
	if cacheDir == "" {
		cacheDir = cache.DEFAULT_CACHE_DIR
	}
	cacheOptions := cache.OpenOptions{
		MaxSizeMB:        500,
		EvictTargetRatio: 0.9,
		MaxFileSizeMB:    4,
		ToolVersion:      Version,
	}
	if merged.CacheMaxSizeMB != nil {
		cacheOptions.MaxSizeMB = *merged.CacheMaxSizeMB
	}
	if merged.CacheEvictTargetRatio != nil {
		cacheOptions.EvictTargetRatio = *merged.CacheEvictTargetRatio
	}
	if merged.CacheMaxFileSizeMB != nil {
		cacheOptions.MaxFileSizeMB = *merged.CacheMaxFileSizeMB
	}

	walkOptions := engine.DefaultWalkOptions()
	walkOptions.IncludeTests = merged.IncludeTests
	walkOptions.Include = merged.Include
	walkOptions.Exclude = merged.Exclude

	return scanPlan{
		context:          ctx,
		walkOptions:      walkOptions,
		cacheDir:         cacheDir,
		cacheOptions:     cacheOptions,
		cacheDisabled:    merged.NoCache,
		baselineDisabled: merged.NoBaseline,
		baselineFile:     merged.BaselineFile,
	}
}

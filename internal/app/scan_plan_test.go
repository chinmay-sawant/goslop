package app

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/cli"
	"github.com/chinmay-sawant/goslop/internal/config"
	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestResolveScanPlanAppliesConfiguredEvictionThreshold(t *testing.T) {
	ratio := 0.5
	maxSize := uint64(128)
	plan := resolveScanPlan(core.ProfileRecommended, &cli.Options{}, &config.Merged{
		CacheEvictTargetRatio: &ratio,
		CacheMaxSizeMB:        &maxSize,
	})
	if plan.cacheOptions.EvictTargetRatio != ratio {
		t.Fatalf("eviction target ratio=%v want %v", plan.cacheOptions.EvictTargetRatio, ratio)
	}
	if plan.cacheOptions.MaxSizeMB != maxSize {
		t.Fatalf("cache max size=%d want %d", plan.cacheOptions.MaxSizeMB, maxSize)
	}
}

func TestResolveScanPlanAppliesRemainingMergedRuntimeFields(t *testing.T) {
	severity := rules.SeverityCritical
	maxFileSize := uint64(9)
	plan := resolveScanPlan(core.ProfileAll, &cli.Options{
		ShowIgnored:   true,
		ShowBaselined: true,
		TaintDepth:    3,
	}, &config.Merged{
		Only:                []string{"PERF-6"},
		Skip:                []string{"CWE-78"},
		Include:             []string{"cmd/**"},
		Exclude:             []string{"vendor/**"},
		IncludeTests:        true,
		NoCache:             true,
		CacheDir:            "custom-cache",
		NoBaseline:          true,
		BaselineFile:        "baseline.json",
		Taint:               true,
		TaintShowPaths:      true,
		BadPracticesEnabled: boolPointer(false),
		BPSeverity:          &severity,
		SeverityOverrides: map[string]rules.Severity{
			"PERF-6": rules.SeverityHigh,
		},
		CacheMaxFileSizeMB: &maxFileSize,
	})

	ctx := plan.context
	if len(ctx.Only) != 1 || ctx.Only[0] != "PERF-6" || len(ctx.Skip) != 1 || ctx.Skip[0] != "CWE-78" {
		t.Fatalf("rule filters only=%v skip=%v", ctx.Only, ctx.Skip)
	}
	if !ctx.IncludeTests || !ctx.NoCache || !ctx.NoBaseline || !ctx.TaintEnabled || !ctx.TaintShowPaths || ctx.TaintMaxDepth != 3 {
		t.Fatalf("context flags=%+v", ctx)
	}
	if !ctx.ShowIgnored || !ctx.ShowBaselined || ctx.BadPracticesEnabled || ctx.BadPracticeSeverity == nil || *ctx.BadPracticeSeverity != severity ||
		ctx.SeverityOverrides["PERF-6"] != rules.SeverityHigh {
		t.Fatalf("context overrides=%+v", ctx)
	}
	if plan.cacheDir != "custom-cache" || !plan.cacheDisabled || !plan.baselineDisabled || plan.baselineFile != "baseline.json" ||
		plan.cacheOptions.MaxFileSizeMB != maxFileSize {
		t.Fatalf("plan=%+v", plan)
	}
	if !plan.walkOptions.IncludeTests || len(plan.walkOptions.Include) != 1 || plan.walkOptions.Include[0] != "cmd/**" ||
		len(plan.walkOptions.Exclude) != 1 || plan.walkOptions.Exclude[0] != "vendor/**" {
		t.Fatalf("walk options=%+v", plan.walkOptions)
	}
}

func boolPointer(v bool) *bool { return &v }

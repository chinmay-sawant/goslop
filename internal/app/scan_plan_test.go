package app

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/cli"
	"github.com/chinmay-sawant/goslop/internal/config"
	"github.com/chinmay-sawant/goslop/internal/core"
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

package perf_test

import (
	"testing"

	"github.com/chinmay/goslop/internal/lang/go/detectors/perf"
)

func TestBatch2RuleIDsRegistered(t *testing.T) {
	d := perf.NewGoPerfScan()
	ids := d.RuleIDs()
	want := []string{
		"PERF-61", "PERF-71", "PERF-76", "PERF-91", "PERF-98",
		"PERF-101", "PERF-103", "PERF-105", "PERF-106", "PERF-111",
	}
	have := map[string]bool{}
	for _, id := range ids {
		have[id] = true
	}
	for _, id := range want {
		if !have[id] {
			t.Errorf("missing registered rule %s (have %d rules)", id, len(ids))
		}
	}
}

func TestFixturePERF061(t *testing.T) {
	runFixtureRule(t, "PERF-061-vulnerable.txt", "PERF-61", true)
	runFixtureRule(t, "PERF-061-safe.txt", "PERF-61", false)
}

func TestFixturePERF067(t *testing.T) {
	runFixtureRule(t, "PERF-067-vulnerable.txt", "PERF-67", true)
	runFixtureRule(t, "PERF-067-safe.txt", "PERF-67", false)
}

func TestFixturePERF071(t *testing.T) {
	runFixtureRule(t, "PERF-071-vulnerable.txt", "PERF-71", true)
	runFixtureRule(t, "PERF-071-safe.txt", "PERF-71", false)
}

func TestFixturePERF076(t *testing.T) {
	runFixtureRule(t, "PERF-076-vulnerable.txt", "PERF-76", true)
	runFixtureRule(t, "PERF-076-safe.txt", "PERF-76", false)
}

func TestFixturePERF091(t *testing.T) {
	runFixtureRule(t, "PERF-091-vulnerable.txt", "PERF-91", true)
	runFixtureRule(t, "PERF-091-safe.txt", "PERF-91", false)
}

func TestFixturePERF098(t *testing.T) {
	runFixtureRule(t, "PERF-098-vulnerable.txt", "PERF-98", true)
	runFixtureRule(t, "PERF-098-safe.txt", "PERF-98", false)
}

func TestFixturePERF101(t *testing.T) {
	runFixtureRule(t, "PERF-101-vulnerable.txt", "PERF-101", true)
	runFixtureRule(t, "PERF-101-safe.txt", "PERF-101", false)
}

func TestFixturePERF103(t *testing.T) {
	runFixtureRule(t, "PERF-103-vulnerable.txt", "PERF-103", true)
	runFixtureRule(t, "PERF-103-safe.txt", "PERF-103", false)
}

func TestFixturePERF105(t *testing.T) {
	runFixtureRule(t, "PERF-105-vulnerable.txt", "PERF-105", true)
	runFixtureRule(t, "PERF-105-safe.txt", "PERF-105", false)
}

func TestFixturePERF106(t *testing.T) {
	runFixtureRule(t, "PERF-106-vulnerable.txt", "PERF-106", true)
	runFixtureRule(t, "PERF-106-safe.txt", "PERF-106", false)
}

func TestFixturePERF110(t *testing.T) {
	runFixtureRule(t, "PERF-110-vulnerable.txt", "PERF-110", true)
	runFixtureRule(t, "PERF-110-safe.txt", "PERF-110", false)
}

func TestFixturePERF111(t *testing.T) {
	runFixtureRule(t, "PERF-111-vulnerable.txt", "PERF-111", true)
	runFixtureRule(t, "PERF-111-safe.txt", "PERF-111", false)
}

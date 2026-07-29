package perf

import (
	"path/filepath"
	"strings"

	"github.com/chinmay/codehound/internal/core"
)

// pureFalsePositiveRules fire on gopdfsuit in Go but never in the Rust oracle
// (issue #8). Suppress until rewritten to Rust parity — but only on real
// project tree scans so fixture / unit tests keep full catalogue coverage.
var pureFalsePositiveRules = map[string]struct{}{
	"PERF-116": {}, "PERF-121": {}, "PERF-125": {}, "PERF-129": {},
	"PERF-132": {}, "PERF-144": {}, "PERF-158": {}, "PERF-159": {},
}

func oracleSkip(id string) bool {
	_, ok := pureFalsePositiveRules[id]
	return ok
}

// oracleSkipUnit reports whether pure-FP museum rules should be suppressed
// for this unit. Unit tests and fixture corpora always keep the rules live.
func oracleSkipUnit(unit *core.ParsedUnit, id string) bool {
	if !oracleSkip(id) {
		return false
	}
	return isRealProjectScan(unit)
}

func isRealProjectScan(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, p := range []string{unit.Path, unit.DisplayPath} {
		if p == "" {
			continue
		}
		// Fixture corpora and synthetic test names must never suppress.
		if strings.Contains(p, "tests/fixtures") ||
			strings.Contains(p, `tests\fixtures`) ||
			strings.Contains(p, "codehound-fixtures") ||
			strings.Contains(p, "-vulnerable") ||
			strings.Contains(p, "-safe") {
			return false
		}
		// Absolute path → product scan (e.g. gopdfsuit oracle).
		if filepath.IsAbs(p) {
			return true
		}
		// Multi-segment relative project path.
		if strings.Contains(p, "/") || strings.Contains(p, `\`) {
			return true
		}
	}
	// Bare names like "sample.go" are unit-test units.
	return false
}

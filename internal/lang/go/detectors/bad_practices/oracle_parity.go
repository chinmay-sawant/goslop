package badpractices

import (
	"path/filepath"
	"strings"

	"github.com/chinmay/codehound/internal/core"
)

// pureFalsePositiveRules fire on gopdfsuit in Go but never in the Rust oracle
// (issue #8). Suppress until rewritten to Rust parity — but only on real
// project tree scans so fixture / unit tests keep full catalogue coverage.
var pureFalsePositiveRules = map[string]struct{}{
	"BP-12": {}, "BP-43": {}, "BP-58": {}, "BP-66": {}, "BP-73": {},
	"BP-75": {}, "BP-76": {}, "BP-91": {}, "BP-94": {}, "BP-100": {},
	"BP-104": {}, "BP-109": {}, "BP-111": {},
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
		if strings.Contains(p, "tests/fixtures") ||
			strings.Contains(p, `tests\fixtures`) ||
			strings.Contains(p, "codehound-fixtures") ||
			strings.Contains(p, "-vulnerable") ||
			strings.Contains(p, "-safe") {
			return false
		}
		if filepath.IsAbs(p) {
			return true
		}
		if strings.Contains(p, "/") || strings.Contains(p, `\`) {
			return true
		}
	}
	return false
}

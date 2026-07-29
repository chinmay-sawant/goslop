// Package detectors aggregates Go language detectors.
package detectors

import (
	"github.com/chinmay/codehound/internal/core"
	badpractices "github.com/chinmay/codehound/internal/lang/go/detectors/bad_practices"
	"github.com/chinmay/codehound/internal/lang/go/detectors/cwe"
	"github.com/chinmay/codehound/internal/lang/go/detectors/cwe/taint"
	"github.com/chinmay/codehound/internal/lang/go/detectors/perf"
)

// All returns the detector set registered with the Go language plugin.
func All() []core.Detector {
	return []core.Detector{
		// Structural CWE catalogue (includes taint-lite seeds when taint is off).
		cwe.NewGoCweScan(),
		// Experimental taint graph (CWE-22/78/79/89 when --taint / security).
		taint.NewDetector(taint.MetaSet{
			CWE22: &cwe.MetaCWE22,
			CWE78: &cwe.MetaCWE78,
			CWE79: &cwe.MetaCWE79,
			CWE89: &cwe.MetaCWE89,
		}),
		perf.NewGoPerfScan(),
		badpractices.NewGoBadPracticeScan(),
	}
}

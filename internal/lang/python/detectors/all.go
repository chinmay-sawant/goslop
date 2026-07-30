// Package detectors aggregates Python language detectors.
package detectors

import (
	"github.com/chinmay-sawant/goslop/internal/core"
	badpractices "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices"
)

// All returns the detector set registered with the Python language plugin.
// Keep additive: CWE/PERF streams may append later without restructuring.
func All() []core.Detector {
	return []core.Detector{
		badpractices.NewPythonBadPracticeScan(),
	}
}

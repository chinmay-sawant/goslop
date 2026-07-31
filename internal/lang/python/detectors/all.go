// Package detectors aggregates Python language detectors.
//
// Active families: CWE priority batch (#52), BP-PY priority subset (#53), and PERF-PY (#54).
package detectors

import (
	"github.com/chinmay-sawant/goslop/internal/core"
	badpractices "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors/perf"
)

// All returns the detector set registered with the Python language plugin.
func All() []core.Detector {
	return []core.Detector{
		cwe.NewPyCweScan(),
		badpractices.NewPythonBadPracticeScan(),
		perf.NewPythonPerfScan(),
	}
}

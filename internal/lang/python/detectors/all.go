// Package detectors aggregates Python language detectors.
//
// Active families: CWE priority batch (#52) and BP-PY priority subset (#53).
// PERF (#54) is deferred. Keep All() additive when new families land.
package detectors

import (
	"github.com/chinmay-sawant/goslop/internal/core"
	badpractices "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe"
)

// All returns the detector set registered with the Python language plugin.
func All() []core.Detector {
	return []core.Detector{
		cwe.NewPyCweScan(),
		badpractices.NewPythonBadPracticeScan(),
	}
}

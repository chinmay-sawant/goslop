// Package detectors aggregates Python language detectors.
//
// Currently CWE priority batch only (#52). BP (#53) will register additively
// via All() when that stream lands — keep this surface additive.
package detectors

import (
	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe"
)

// All returns the detector set registered with the Python language plugin.
func All() []core.Detector {
	return []core.Detector{
		cwe.NewPyCweScan(),
	}
}

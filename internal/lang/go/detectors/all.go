// Package detectors aggregates Go language detectors.
package detectors

import (
	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/detectors/cwe"
	"github.com/chinmay/codehound/internal/lang/go/detectors/perf"
)

// All returns the seed detector set registered with the Go language plugin.
func All() []core.Detector {
	return []core.Detector{
		cwe.NewCWE78(),
		cwe.NewCWE89(),
		perf.NewPERF116(),
	}
}

package perf

// pureFalsePositiveRules fire on gopdfsuit in Go but never in the Rust oracle
// (issue #8). Suppress until rewritten to Rust parity.
var pureFalsePositiveRules = map[string]struct{}{
	"PERF-116": {}, "PERF-121": {}, "PERF-125": {}, "PERF-129": {},
	"PERF-132": {}, "PERF-144": {}, "PERF-158": {}, "PERF-159": {},
}

func oracleSkip(id string) bool {
	_, ok := pureFalsePositiveRules[id]
	return ok
}

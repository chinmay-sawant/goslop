package badpractices

// pureFalsePositiveRules fire on gopdfsuit in Go but never in the Rust oracle
// (issue #8). Their heuristics are too broad; suppress emission until rewritten.
var pureFalsePositiveRules = map[string]struct{}{
	"BP-12": {}, "BP-43": {}, "BP-58": {}, "BP-66": {}, "BP-73": {},
	"BP-75": {}, "BP-76": {}, "BP-91": {}, "BP-94": {}, "BP-100": {},
	"BP-104": {}, "BP-109": {}, "BP-111": {},
}

func oracleSkip(id string) bool {
	_, ok := pureFalsePositiveRules[id]
	return ok
}

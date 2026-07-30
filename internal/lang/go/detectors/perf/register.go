package perf

import (
	"sync"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

var registerMu sync.Mutex

// RegisterRule appends a PERF rule to the global catalogue.
// Call from init() in batch files so multiple domains can register in parallel.
// Optional gates are any-of SourceIndex needles (must be in perfNeedles); when
// set, Run skips the rule if none are present. Nil/empty gates = always run.
func RegisterRule(id string, fn ruleFn, meta *rules.RuleMetadata, gates ...string) {
	if id == "" || fn == nil || meta == nil {
		return
	}
	var gateCopy []string
	if len(gates) > 0 {
		gateCopy = append([]string(nil), gates...)
	}
	registerMu.Lock()
	defer registerMu.Unlock()
	entry := ruleEntry{id: id, fn: fn, meta: meta, gates: gateCopy}
	// replace if already present (idempotent re-register)
	for i, e := range goPerfRules {
		if e.id == id {
			goPerfRules[i] = entry
			if metaByID == nil {
				metaByID = map[string]*rules.RuleMetadata{}
			}
			metaByID[id] = meta
			return
		}
	}
	goPerfRules = append(goPerfRules, entry)
	if metaByID == nil {
		metaByID = map[string]*rules.RuleMetadata{}
	}
	metaByID[id] = meta
}

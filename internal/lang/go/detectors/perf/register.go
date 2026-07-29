package perf

import (
	"sync"

	"github.com/chinmay/codehound/internal/rules"
)

var registerMu sync.Mutex

// RegisterRule appends a PERF rule to the global catalogue.
// Call from init() in batch files so multiple domains can register in parallel.
func RegisterRule(id string, fn ruleFn, meta *rules.RuleMetadata) {
	if id == "" || fn == nil || meta == nil {
		return
	}
	registerMu.Lock()
	defer registerMu.Unlock()
	// replace if already present (idempotent re-register)
	for i, e := range goPerfRules {
		if e.id == id {
			goPerfRules[i] = ruleEntry{id: id, fn: fn, meta: meta}
			metaByID[id] = meta
			return
		}
	}
	goPerfRules = append(goPerfRules, ruleEntry{id: id, fn: fn, meta: meta})
	if metaByID == nil {
		metaByID = map[string]*rules.RuleMetadata{}
	}
	metaByID[id] = meta
}

package badpractices

import (
	"sync"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

// ruleFn is one BP detector over a prebuilt fact bag.
type ruleFn func(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding)

type ruleEntry struct {
	id   string
	fn   ruleFn
	meta *rules.RuleMetadata
}

var (
	registerMu sync.Mutex
	goBPRules  []ruleEntry
)

// RegisterRule appends a BP rule to the global catalogue.
// Prefer calling from init() so multiple domains can register.
func RegisterRule(id string, fn ruleFn) {
	if id == "" || fn == nil {
		return
	}
	meta := MetadataForID(id)
	if meta == nil {
		// Fall back to a stub so detectors still register even if catalogue is incomplete.
		meta = &rules.RuleMetadata{
			ID:       id,
			Title:    id,
			Severity: rules.SeverityLow,
			Pack:     rules.PackBadPractice,
		}
	}
	registerMu.Lock()
	defer registerMu.Unlock()
	for i, e := range goBPRules {
		if e.id == id {
			goBPRules[i] = ruleEntry{id: id, fn: fn, meta: meta}
			return
		}
	}
	goBPRules = append(goBPRules, ruleEntry{id: id, fn: fn, meta: meta})
}

// registeredRuleIDs returns a snapshot of implemented rule ids (sorted later by caller).
func registeredRuleIDs() []string {
	registerMu.Lock()
	defer registerMu.Unlock()
	ids := make([]string, len(goBPRules))
	for i, e := range goBPRules {
		ids[i] = e.id
	}
	return ids
}

func snapshotRules() []ruleEntry {
	registerMu.Lock()
	defer registerMu.Unlock()
	out := make([]ruleEntry, len(goBPRules))
	copy(out, goBPRules)
	return out
}

package badpractices

import (
	"sync"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// ruleFn is one BP-PY detector over a prebuilt fact bag.
type ruleFn func(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding)

type ruleEntry struct {
	id   string
	fn   ruleFn
	meta *rules.RuleMetadata
}

var (
	registerMu  sync.Mutex
	pythonRules []ruleEntry
)

// RegisterRule appends a BP-PY rule to the global catalogue.
// Prefer calling from init() so multiple domains can register.
func RegisterRule(id string, fn ruleFn) {
	if id == "" || fn == nil {
		return
	}
	meta := MetadataForID(id)
	if meta == nil {
		meta = &rules.RuleMetadata{
			ID:       id,
			Title:    id,
			Severity: rules.SeverityLow,
			Pack:     rules.PackBadPractice,
		}
	}
	registerMu.Lock()
	defer registerMu.Unlock()
	for i, e := range pythonRules {
		if e.id == id {
			pythonRules[i] = ruleEntry{id: id, fn: fn, meta: meta}
			return
		}
	}
	pythonRules = append(pythonRules, ruleEntry{id: id, fn: fn, meta: meta})
}

func snapshotRules() []ruleEntry {
	registerMu.Lock()
	defer registerMu.Unlock()
	out := make([]ruleEntry, len(pythonRules))
	copy(out, pythonRules)
	return out
}

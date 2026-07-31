package perf

import (
	"sync"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

type ruleFn func(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding)

type ruleEntry struct {
	id string
	fn ruleFn
}

var (
	registerMu  sync.Mutex
	pythonRules []ruleEntry
)

// RegisterRule adds or replaces a PERF-PY rule during package initialization.
func RegisterRule(id string, fn ruleFn) {
	if id == "" || fn == nil || MetadataForID(id) == nil {
		return
	}
	registerMu.Lock()
	defer registerMu.Unlock()
	for i, entry := range pythonRules {
		if entry.id == id {
			pythonRules[i] = ruleEntry{id: id, fn: fn}
			return
		}
	}
	pythonRules = append(pythonRules, ruleEntry{id: id, fn: fn})
}

func snapshotRules() []ruleEntry {
	registerMu.Lock()
	defer registerMu.Unlock()
	out := make([]ruleEntry, len(pythonRules))
	copy(out, pythonRules)
	return out
}

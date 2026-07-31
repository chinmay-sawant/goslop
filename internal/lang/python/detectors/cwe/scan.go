// Package cwe implements pure-Go source-pattern CWE heuristics for Python.
//
// No CGO / tree-sitter: detectors scan unit.Source with needle prefilters
// (ast.SourceIndex) and light call-site heuristics. Priority batch: CWE-22,
// CWE-78, CWE-79, CWE-89, CWE-502 (issue #52). Full 344-rule catalogue and
// inter-procedural taint are out of scope.
package cwe

import (
	"sort"
	"sync"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// pyCweRuleFn is one CWE detector over a prebuilt fact bag.
type pyCweRuleFn func(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding)

type pyCweRuleEntry struct {
	id   string
	fn   pyCweRuleFn
	meta *rules.RuleMetadata
	// gates are optional SourceIndex needles (any-of). When non-empty, Run
	// skips the rule when facts.Index has none of them. Nil/empty = always run.
	// Must be FN-safe: only needles the rule needs before it can emit.
	gates []string
}

var (
	pyCweRegisterMu sync.Mutex
	pyCweRules      []pyCweRuleEntry
	metaByID        map[string]*rules.RuleMetadata
)

// RegisterRule appends a Python CWE rule to the global catalogue.
// Call from init() so rule files can register independently.
// Optional gates are any-of SourceIndex needles; when set, the rule is skipped
// if none are present in the file (see facts.Index / pyCweNeedles).
func RegisterRule(id string, fn pyCweRuleFn, meta *rules.RuleMetadata, gates ...string) {
	if id == "" || fn == nil || meta == nil {
		return
	}
	var gateCopy []string
	if len(gates) > 0 {
		gateCopy = append([]string(nil), gates...)
	}
	pyCweRegisterMu.Lock()
	defer pyCweRegisterMu.Unlock()
	entry := pyCweRuleEntry{id: id, fn: fn, meta: meta, gates: gateCopy}
	for i, e := range pyCweRules {
		if e.id == id {
			pyCweRules[i] = entry
			if metaByID == nil {
				metaByID = map[string]*rules.RuleMetadata{}
			}
			metaByID[id] = meta
			return
		}
	}
	pyCweRules = append(pyCweRules, entry)
	if metaByID == nil {
		metaByID = map[string]*rules.RuleMetadata{}
	}
	metaByID[id] = meta
}

// PyCweScan is the unified Python CWE detector (one Detector, many rules).
// Mirrors Go GoCweScan: build facts once, then dispatch allowed rules.
type PyCweScan struct {
	core.BaseDetector
}

// NewPyCweScan returns the unified Python CWE detector bundle.
func NewPyCweScan() *PyCweScan {
	return &PyCweScan{}
}

// Language implements core.Detector.
func (d *PyCweScan) Language() core.LanguageID { return core.LanguagePython }

// Immutable catalogue snapshot after init() registration.
var (
	pyCweCatalogueOnce sync.Once
	pyCweCatalogue     []pyCweRuleEntry
	pyCweRuleIDs       []string
)

func pyCweCatalogueSnapshot() []pyCweRuleEntry {
	pyCweCatalogueOnce.Do(func() {
		pyCweRegisterMu.Lock()
		pyCweCatalogue = append([]pyCweRuleEntry(nil), pyCweRules...)
		ids := make([]string, len(pyCweRules))
		for i, e := range pyCweRules {
			ids[i] = e.id
		}
		pyCweRegisterMu.Unlock()
		sort.Strings(ids)
		pyCweRuleIDs = ids
	})
	return pyCweCatalogue
}

// RuleIDs implements core.Detector.
func (d *PyCweScan) RuleIDs() []string {
	_ = pyCweCatalogueSnapshot()
	return pyCweRuleIDs
}

// MetadataFor implements core.Detector.
func (d *PyCweScan) MetadataFor(ruleID string) *rules.RuleMetadata {
	pyCweRegisterMu.Lock()
	defer pyCweRegisterMu.Unlock()
	if m, ok := metaByID[ruleID]; ok {
		return m
	}
	return nil
}

// Run implements core.Detector.
func (d *PyCweScan) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if unit.Source == "" {
		return
	}
	all := pyCweCatalogueSnapshot()
	hasAllowedRule := false
	for _, e := range all {
		if ctx != nil && !ctx.Allows(e.id) {
			continue
		}
		hasAllowedRule = true
		break
	}
	if !hasAllowedRule {
		return
	}

	facts := BuildFacts(unit)
	for _, e := range all {
		if ctx != nil && !ctx.Allows(e.id) {
			continue
		}
		if len(e.gates) > 0 && !facts.Index.HasAny(e.gates) {
			continue
		}
		e.fn(unit, facts, out)
	}
}

// RegisteredRuleCount returns how many CWE rules are registered (tests).
func RegisteredRuleCount() int {
	pyCweRegisterMu.Lock()
	defer pyCweRegisterMu.Unlock()
	return len(pyCweRules)
}

// RegisteredRuleIDs returns a sorted copy of registered rule ids (tests).
func RegisteredRuleIDs() []string {
	_ = pyCweCatalogueSnapshot()
	out := make([]string, len(pyCweRuleIDs))
	copy(out, pyCweRuleIDs)
	return out
}

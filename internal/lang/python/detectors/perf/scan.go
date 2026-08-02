// Package perf implements source-only Python performance heuristics.
package perf

import (
	"sort"
	"sync"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// PythonPerfScan is the unified Python PERF-PY detector.
type PythonPerfScan struct{ core.BaseDetector }

// NewPythonPerfScan returns the PERF-PY detector bundle.
func NewPythonPerfScan() *PythonPerfScan { return &PythonPerfScan{} }

// Language implements core.Detector.
func (d *PythonPerfScan) Language() core.LanguageID { return core.LanguagePython }

var (
	perfCatalogueOnce sync.Once
	perfCatalogue     []ruleEntry
	perfRuleIDs       []string
)

func perfCatalogueSnapshot() []ruleEntry {
	perfCatalogueOnce.Do(func() {
		perfCatalogue = snapshotRules()
		perfRuleIDs = make([]string, len(perfCatalogue))
		for i, entry := range perfCatalogue {
			perfRuleIDs[i] = entry.id
		}
		sort.Strings(perfRuleIDs)
	})
	return perfCatalogue
}

// RuleIDs implements core.Detector.
func (d *PythonPerfScan) RuleIDs() []string {
	_ = perfCatalogueSnapshot()
	return append([]string(nil), perfRuleIDs...)
}

// MetadataFor implements core.Detector.
func (d *PythonPerfScan) MetadataFor(ruleID string) *rules.RuleMetadata { return MetadataForID(ruleID) }

// Run implements core.Detector.
func (d *PythonPerfScan) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil || unit.Source == "" || (unit.Language != core.LanguagePython && unit.Language != core.LangPython) {
		return
	}
	all := perfCatalogueSnapshot()
	hasAllowed := false
	for _, entry := range all {
		if ctx == nil || ctx.Allows(entry.id) {
			hasAllowed = true
			break
		}
	}
	if !hasAllowed {
		return
	}
	facts := buildFacts(unit)
	for _, entry := range all {
		if ctx != nil && !ctx.Allows(entry.id) {
			continue
		}
		start := len(*out)
		entry.fn(unit, facts, out)
		if ctx != nil {
			for i := start; i < len(*out); i++ {
				ctx.ApplyFindingOverrides(&(*out)[i])
			}
		}
	}
}

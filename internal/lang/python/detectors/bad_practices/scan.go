package badpractices

import (
	"sort"
	"sync"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// PythonBadPracticeScan is the unified Python bad-practice detector (one Detector, many rules).
type PythonBadPracticeScan struct {
	core.BaseDetector
}

// NewPythonBadPracticeScan returns the unified BP-PY detector bundle.
func NewPythonBadPracticeScan() *PythonBadPracticeScan {
	return &PythonBadPracticeScan{}
}

// Language implements core.Detector.
func (d *PythonBadPracticeScan) Language() core.LanguageID { return core.LanguagePython }

// Immutable catalogue snapshot after init() registration.
var (
	bpCatalogueOnce sync.Once
	bpCatalogue     []ruleEntry
	bpRuleIDs       []string
)

func bpCatalogueSnapshot() []ruleEntry {
	bpCatalogueOnce.Do(func() {
		bpCatalogue = snapshotRules()
		ids := make([]string, len(bpCatalogue))
		for i, e := range bpCatalogue {
			ids[i] = e.id
		}
		sort.Strings(ids)
		bpRuleIDs = ids
	})
	return bpCatalogue
}

// RuleIDs implements core.Detector.
func (d *PythonBadPracticeScan) RuleIDs() []string {
	_ = bpCatalogueSnapshot()
	return bpRuleIDs
}

// MetadataFor implements core.Detector.
func (d *PythonBadPracticeScan) MetadataFor(ruleID string) *rules.RuleMetadata {
	return MetadataForID(ruleID)
}

// Run implements core.Detector.
func (d *PythonBadPracticeScan) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if unit.Language != core.LanguagePython && unit.Language != core.LangPython {
		return
	}
	if unit.Source == "" {
		return
	}

	all := bpCatalogueSnapshot()
	any := false
	for _, e := range all {
		if ctx == nil || ctx.Allows(e.id) {
			any = true
			break
		}
	}
	if !any {
		return
	}

	facts := buildFacts(unit)
	defer facts.close()

	for _, e := range all {
		if ctx != nil && !ctx.Allows(e.id) {
			continue
		}
		start := len(*out)
		e.fn(unit, facts, out)
		if ctx != nil {
			for i := start; i < len(*out); i++ {
				ctx.ApplyFindingOverrides(&(*out)[i])
			}
		}
	}
}

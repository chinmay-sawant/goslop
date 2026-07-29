package cwe

import (
	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

// CWE89Detector is a single-rule adapter retained for seed unit tests.
// Production registration uses GoCweScan via RegisterRule("CWE-89", ...).
type CWE89Detector struct {
	core.BaseDetector
}

// NewCWE89 returns a CWE-89 detector.
func NewCWE89() *CWE89Detector {
	return &CWE89Detector{}
}

// Language implements core.Detector.
func (d *CWE89Detector) Language() core.LanguageID { return core.LangGo }

// RuleIDs implements core.Detector.
func (d *CWE89Detector) RuleIDs() []string { return []string{"CWE-89"} }

// MetadataFor implements core.Detector.
func (d *CWE89Detector) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == "CWE-89" {
		return &MetaCWE89
	}
	return nil
}

// Run implements core.Detector.
func (d *CWE89Detector) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if ctx != nil && !ctx.Allows("CWE-89") {
		return
	}
	detectCWE89(unit, BuildFacts(unit), out)
}

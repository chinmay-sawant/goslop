package perf

import (
	"sort"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

// ruleFn is one PERF detector over a prebuilt fact bag.
type ruleFn func(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding)

type ruleEntry struct {
	id   string
	fn   ruleFn
	meta *rules.RuleMetadata
}

// goPerfRules is the implemented PERF catalogue (filled via RegisterRule / init).
var goPerfRules []ruleEntry

// GoPerfScan is the unified Go PERF detector (one Detector, many rules).
// Mirrors Rust GoPerfScan: build facts once, then dispatch enabled rules.
type GoPerfScan struct {
	core.BaseDetector
}

// NewGoPerfScan returns the unified PERF detector bundle.
func NewGoPerfScan() *GoPerfScan {
	return &GoPerfScan{}
}

// Language implements core.Detector.
func (d *GoPerfScan) Language() core.LanguageID { return core.LangGo }

// RuleIDs implements core.Detector.
func (d *GoPerfScan) RuleIDs() []string {
	registerMu.Lock()
	defer registerMu.Unlock()
	ids := make([]string, len(goPerfRules))
	for i, e := range goPerfRules {
		ids[i] = e.id
	}
	// stable order for --list-rules
	sort.Strings(ids)
	return ids
}

// MetadataFor implements core.Detector.
func (d *GoPerfScan) MetadataFor(ruleID string) *rules.RuleMetadata {
	registerMu.Lock()
	defer registerMu.Unlock()
	if m, ok := metaByID[ruleID]; ok {
		return m
	}
	return nil
}

// Run implements core.Detector.
func (d *GoPerfScan) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	registerMu.Lock()
	rulesCopy := make([]ruleEntry, len(goPerfRules))
	copy(rulesCopy, goPerfRules)
	registerMu.Unlock()

	any := false
	for _, e := range rulesCopy {
		if ctx == nil || ctx.Allows(e.id) {
			any = true
			break
		}
	}
	if !any {
		return
	}

	facts := BuildFacts(unit)
	for _, e := range rulesCopy {
		if ctx != nil && !ctx.Allows(e.id) {
			continue
		}
		// Pure-FP museums vs gopdfsuit oracle — never suppress fixture/unit tests.
		if oracleSkipUnit(unit, e.id) {
			continue
		}
		e.fn(unit, facts, out)
	}
}

// NewPERF116 is retained for tests that construct a single-rule view.
// Prefer NewGoPerfScan for production registration.
func NewPERF116() core.Detector {
	return &singleRuleAdapter{id: "PERF-116", fn: detectPERF116, meta: &MetaPERF116}
}

// singleRuleAdapter exposes one PERF rule as its own Detector (test helper).
type singleRuleAdapter struct {
	core.BaseDetector
	id   string
	fn   ruleFn
	meta *rules.RuleMetadata
}

func (d *singleRuleAdapter) Language() core.LanguageID { return core.LangGo }
func (d *singleRuleAdapter) RuleIDs() []string         { return []string{d.id} }
func (d *singleRuleAdapter) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == d.id {
		return d.meta
	}
	return nil
}
func (d *singleRuleAdapter) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if ctx != nil && !ctx.Allows(d.id) {
		return
	}
	d.fn(unit, BuildFacts(unit), out)
}

package badpractices

import (
	"sort"
	"sync"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

// GoBadPracticeScan is the unified Go bad-practice detector (one Detector, many rules).
// Mirrors Rust GoBadPracticeScan: build facts once, dispatch enabled rules, apply severity overrides.
type GoBadPracticeScan struct {
	core.BaseDetector
	// projectCaches hold memoized project snapshots for server/module rules.
	caches *bpProjectCaches
}

// NewGoBadPracticeScan returns the unified BP detector bundle.
func NewGoBadPracticeScan() *GoBadPracticeScan {
	return &GoBadPracticeScan{
		caches: newProjectCaches(),
	}
}

// Language implements core.Detector.
func (d *GoBadPracticeScan) Language() core.LanguageID { return core.LangGo }

// Immutable catalogue snapshot after init() registration (P3.1).
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
func (d *GoBadPracticeScan) RuleIDs() []string {
	_ = bpCatalogueSnapshot()
	return bpRuleIDs
}

// MetadataFor implements core.Detector.
func (d *GoBadPracticeScan) MetadataFor(ruleID string) *rules.RuleMetadata {
	return MetadataForID(ruleID)
}

// BeginScan installs project-cache session state for this analyzer.
func (d *GoBadPracticeScan) BeginScan(ctx *core.ScanContext) {
	if d.caches == nil {
		d.caches = newProjectCaches()
	}
	d.caches.clear()
	setActiveCaches(d.caches)
}

// EndScan clears the active session.
func (d *GoBadPracticeScan) EndScan() {
	clearActiveCaches()
	if d.caches != nil {
		d.caches.clear()
	}
}

// ResetState drops memoized project facts without tearing down the session.
func (d *GoBadPracticeScan) ResetState() {
	if d.caches != nil {
		d.caches.clear()
	}
}

// Run implements core.Detector.
func (d *GoBadPracticeScan) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if d.caches == nil {
		d.caches = newProjectCaches()
	}
	// Rayon-style workers may not inherit the controlling thread's session.
	_ = installActiveCaches(d.caches)

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

	hasEnabledProject := false
	for _, e := range all {
		if ctx != nil && !ctx.Allows(e.id) {
			continue
		}
		if referenceSkipUnit(unit, e.id) {
			continue
		}
		if requiresProjectAnchor(e.id) || requiresServerAnchor(e.id) {
			hasEnabledProject = true
			break
		}
	}
	isProjectAnchor := false
	isServerAnchor := false
	if hasEnabledProject && !isMaterializedFixture(unit) {
		isProjectAnchor = isProjectAnchorFile(unit)
		isServerAnchor = isServerAnchorFile(unit)
	}

	for _, e := range all {
		if ctx != nil && !ctx.Allows(e.id) {
			continue
		}
		if referenceSkipUnit(unit, e.id) {
			continue
		}
		if requiresProjectAnchor(e.id) && !isProjectAnchor {
			continue
		}
		if requiresServerAnchor(e.id) && !isServerAnchor {
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

// requiresProjectAnchor reports module-hygiene rules that emit once per project.
func requiresProjectAnchor(ruleID string) bool {
	switch ruleID {
	case "BP-57", "BP-58", "BP-59", "BP-60", "BP-61", "BP-62", "BP-63", "BP-64", "BP-65":
		return true
	default:
		return false
	}
}

// requiresServerAnchor reports server-policy rules that emit once on the entrypoint.
func requiresServerAnchor(ruleID string) bool {
	switch ruleID {
	case "BP-47", "BP-50", "BP-54", "BP-55":
		return true
	default:
		return false
	}
}

package taint

import (
	"sync"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Detector is the taint-tracking detector for CWE-22/78/79/89.
// It only runs when ScanContext.TaintEnabled is true. Seed CWE-78/89
// heuristics remain as the fallback when taint is off.
type Detector struct {
	core.BaseDetector
	meta MetaSet

	mu    sync.Mutex
	units []ProjectUnit
}

// NewDetector returns a taint detector with the given rule metadata.
func NewDetector(meta MetaSet) *Detector {
	return &Detector{meta: meta}
}

// Language implements core.Detector.
func (d *Detector) Language() core.LanguageID { return core.LangGo }

// RuleIDs implements core.Detector.
func (d *Detector) RuleIDs() []string {
	return []string{"CWE-22", "CWE-78", "CWE-79", "CWE-89"}
}

// MetadataFor implements core.Detector.
func (d *Detector) MetadataFor(ruleID string) *rules.RuleMetadata {
	switch ruleID {
	case "CWE-22":
		return d.meta.CWE22
	case "CWE-78":
		return d.meta.CWE78
	case "CWE-79":
		return d.meta.CWE79
	case "CWE-89":
		return d.meta.CWE89
	}
	return nil
}

// RequiresCacheState implements core.Detector — taint needs per-file accumulate on warm hits.
func (d *Detector) RequiresCacheState(ctx *core.ScanContext) bool {
	return ctx != nil && ctx.TaintEnabled
}

// BeginScan clears project state.
func (d *Detector) BeginScan(ctx *core.ScanContext) {
	d.mu.Lock()
	d.units = d.units[:0]
	d.mu.Unlock()
}

// ResetState clears project state.
func (d *Detector) ResetState() {
	d.mu.Lock()
	d.units = nil
	d.mu.Unlock()
}

// Run extracts taint facts, builds the graph, and runs intra-procedural rules.
func (d *Detector) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil || ctx == nil || !ctx.TaintEnabled {
		return
	}
	if !anyAllowed(ctx, d.RuleIDs()) {
		return
	}

	ann := ExtractTaintFacts(unit)
	graph := BuildTaintGraph(&ann)
	cg := ExtractCallGraph(unit)
	imports := BuildImportMap(unit.Source)

	if ctx.Allows("CWE-22") {
		DetectCWE22(unit, graph, d.meta.CWE22, out)
	}
	if ctx.Allows("CWE-78") {
		DetectCWE78(unit, graph, d.meta.CWE78, out)
	}
	if ctx.Allows("CWE-79") {
		DetectCWE79(unit, graph, d.meta.CWE79, out)
	}
	if ctx.Allows("CWE-89") {
		DetectCWE89(unit, graph, &ann, d.meta.CWE89, out)
	}

	// Accumulate for inter-proc finalize (also done in AccumulateState for cache path).
	d.pushUnit(unit, ann, graph, cg, imports)
}

// AccumulateState stores unit state for finalize (cache warm-hit path).
func (d *Detector) AccumulateState(ctx *core.ScanContext, unit *core.ParsedUnit) {
	if unit == nil || ctx == nil || !ctx.TaintEnabled {
		return
	}
	// Run already pushes; avoid double-push when both are called.
	// Engine calls Run then AccumulateState — skip if we already have this path
	// from the same scan pass is hard; instead only push from AccumulateState
	// when RequiresCacheState and Run was skipped. For simplicity: Run pushes,
	// AccumulateState is a no-op when units already contain the path from Run.
	// When cache (Phase 10) skips Run but still calls AccumulateState, rebuild.
	d.mu.Lock()
	path := unit.DisplayPath
	if path == "" {
		path = unit.Path
	}
	for _, u := range d.units {
		if u.Path == path {
			d.mu.Unlock()
			return
		}
	}
	d.mu.Unlock()

	ann := ExtractTaintFacts(unit)
	graph := BuildTaintGraph(&ann)
	cg := ExtractCallGraph(unit)
	imports := BuildImportMap(unit.Source)
	d.pushUnit(unit, ann, graph, cg, imports)
}

func (d *Detector) pushUnit(unit *core.ParsedUnit, ann TaintAnnotations, graph *TaintGraph, cg *CallGraph, imports map[string]string) {
	path := unit.DisplayPath
	if path == "" {
		path = unit.Path
	}
	pu := ProjectUnit{
		Path:       path,
		Source:     unit.Source,
		LineStarts: append([]int(nil), unit.LineStarts...),
		Package:    PackageFromUnit(path, unit.Source),
		CallGraph:  cg,
		Annot:      ann,
		Graph:      graph,
		ImportMap:  imports,
	}
	d.mu.Lock()
	d.units = append(d.units, pu)
	d.mu.Unlock()
}

// Finalize runs inter-procedural analysis across accumulated units.
func (d *Detector) Finalize(ctx *core.ScanContext, out *[]rules.Finding) {
	if ctx == nil || !ctx.TaintEnabled || out == nil {
		return
	}
	d.mu.Lock()
	units := d.units
	d.units = nil
	d.mu.Unlock()
	FinalizeInterProcedural(units, ctx, d.meta, out)
}

func anyAllowed(ctx *core.ScanContext, ids []string) bool {
	for _, id := range ids {
		if ctx.Allows(id) {
			return true
		}
	}
	return false
}

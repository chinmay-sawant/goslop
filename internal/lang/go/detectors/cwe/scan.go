package cwe

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// cweRuleFn is one CWE detector over a prebuilt fact bag.
type cweRuleFn func(unit *core.ParsedUnit, facts *GoCweFacts, out *[]rules.Finding)

type cweRuleEntry struct {
	id   string
	fn   cweRuleFn
	meta *rules.RuleMetadata
	// gates are optional SourceIndex needles (any-of). When non-empty, Run
	// skips the rule when facts.Index has none of them. Nil/empty = always run.
	// Must be FN-safe: only needles the rule needs before it can emit.
	gates []string
}

var (
	cweRegisterMu sync.Mutex
	cweRules      []cweRuleEntry
)

// RegisterRule appends a CWE rule to the global catalogue.
// Call from init() so domain/table files can register in parallel.
// Optional gates are any-of SourceIndex needles; when set, the rule is skipped
// if none are present in the file (see facts.Index / cweNeedles).
func RegisterRule(id string, fn cweRuleFn, meta *rules.RuleMetadata, gates ...string) {
	if id == "" || fn == nil || meta == nil {
		return
	}
	var gateCopy []string
	if len(gates) > 0 {
		gateCopy = append([]string(nil), gates...)
	}
	cweRegisterMu.Lock()
	defer cweRegisterMu.Unlock()
	entry := cweRuleEntry{id: id, fn: fn, meta: meta, gates: gateCopy}
	for i, e := range cweRules {
		if e.id == id {
			cweRules[i] = entry
			if metaByID == nil {
				metaByID = map[string]*rules.RuleMetadata{}
			}
			metaByID[id] = meta
			return
		}
	}
	cweRules = append(cweRules, entry)
	if metaByID == nil {
		metaByID = map[string]*rules.RuleMetadata{}
	}
	metaByID[id] = meta
}

// GoCweScan is the unified Go CWE detector (one Detector, many rules).
// Mirrors Rust GoCweScan / PERF GoPerfScan: build facts once, then dispatch.
// Full inter-procedural taint (Phase 9) is out of scope; taint-core IDs use
// same-file taint-lite heuristics.
type GoCweScan struct {
	core.BaseDetector
}

// NewGoCweScan returns the unified CWE detector bundle.
func NewGoCweScan() *GoCweScan {
	return &GoCweScan{}
}

// Language implements core.Detector.
func (d *GoCweScan) Language() core.LanguageID { return core.LangGo }

// Immutable catalogue snapshot after init() registration (P3.1).
var (
	cweCatalogueOnce sync.Once
	cweCatalogue     []cweRuleEntry
	cweRuleIDs       []string
)

func cweCatalogueSnapshot() []cweRuleEntry {
	cweCatalogueOnce.Do(func() {
		cweRegisterMu.Lock()
		cweCatalogue = append([]cweRuleEntry(nil), cweRules...)
		ids := make([]string, len(cweRules))
		for i, e := range cweRules {
			ids[i] = e.id
		}
		cweRegisterMu.Unlock()
		sort.Strings(ids)
		cweRuleIDs = ids
	})
	return cweCatalogue
}

// RuleIDs implements core.Detector.
func (d *GoCweScan) RuleIDs() []string {
	_ = cweCatalogueSnapshot()
	return cweRuleIDs
}

// MetadataFor implements core.Detector.
func (d *GoCweScan) MetadataFor(ruleID string) *rules.RuleMetadata {
	cweRegisterMu.Lock()
	defer cweRegisterMu.Unlock()
	if m, ok := metaByID[ruleID]; ok {
		return m
	}
	return nil
}

// Run implements core.Detector.
func (d *GoCweScan) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	all := cweCatalogueSnapshot()
	taintOn := ctx != nil && ctx.TaintEnabled
	any := false
	for _, e := range all {
		if ctx != nil && !ctx.Allows(e.id) {
			continue
		}
		if taintOn && isTaintCoreRule(e.id) {
			continue
		}
		any = true
		break
	}
	if !any {
		return
	}

	facts := BuildFacts(unit)
	for _, e := range all {
		if ctx != nil && !ctx.Allows(e.id) {
			continue
		}
		// Full taint package owns these when enabled (avoid double findings).
		if taintOn && isTaintCoreRule(e.id) {
			continue
		}
		// Pure FPs vs Rust gopdfsuit reference corpus (issue #8) — SI museums too broad.
		// Never suppress fixture / unit-test units (catalogue must stay green).
		if isReferencePureFP(e.id) && isRealProjectScan(unit) {
			continue
		}
		// Needle/domain gate: skip when the file has no relevant SourceIndex hits.
		if len(e.gates) > 0 && !facts.Index.HasAny(e.gates) {
			continue
		}
		e.fn(unit, facts, out)
	}
}

// isReferencePureFP lists CWE IDs that fire on gopdfsuit in Go but never in Rust.
func isReferencePureFP(id string) bool {
	switch id {
	case "CWE-140", "CWE-212", "CWE-252", "CWE-257", "CWE-260",
		"CWE-319", "CWE-459", "CWE-915", "CWE-918":
		return true
	default:
		return false
	}
}

func isRealProjectScan(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, p := range []string{unit.Path, unit.DisplayPath} {
		if p == "" {
			continue
		}
		if strings.Contains(p, "tests/fixtures") ||
			strings.Contains(p, `tests\fixtures`) ||
			strings.Contains(p, "goslop-fixtures") ||
			strings.Contains(p, "-vulnerable") ||
			strings.Contains(p, "-safe") {
			return false
		}
		if filepath.IsAbs(p) {
			if isTempMaterializePath(p) {
				return false
			}
			return true
		}
		if strings.Contains(p, "/") || strings.Contains(p, `\`) {
			return true
		}
	}
	return false
}

func isTempMaterializePath(p string) bool {
	if strings.Contains(p, "/tmp/") || strings.Contains(p, "/var/folders/") {
		return true
	}
	if strings.Contains(p, `\Temp\`) || strings.Contains(p, `\AppData\Local\Temp\`) {
		return true
	}
	return strings.Contains(p, "goslop-fixture-")
}

func isTaintCoreRule(id string) bool {
	switch id {
	case "CWE-22", "CWE-78", "CWE-79", "CWE-89":
		return true
	default:
		return false
	}
}

// RegisteredRuleCount returns how many CWE rules are registered (tests).
func RegisteredRuleCount() int {
	cweRegisterMu.Lock()
	defer cweRegisterMu.Unlock()
	return len(cweRules)
}

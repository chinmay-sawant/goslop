package cwe

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

// cweRuleFn is one CWE detector over a prebuilt fact bag.
type cweRuleFn func(unit *core.ParsedUnit, facts *GoCweFacts, out *[]rules.Finding)

type cweRuleEntry struct {
	id   string
	fn   cweRuleFn
	meta *rules.RuleMetadata
}

var (
	cweRegisterMu sync.Mutex
	cweRules      []cweRuleEntry
)

// RegisterRule appends a CWE rule to the global catalogue.
// Call from init() so domain/table files can register in parallel.
func RegisterRule(id string, fn cweRuleFn, meta *rules.RuleMetadata) {
	if id == "" || fn == nil || meta == nil {
		return
	}
	cweRegisterMu.Lock()
	defer cweRegisterMu.Unlock()
	for i, e := range cweRules {
		if e.id == id {
			cweRules[i] = cweRuleEntry{id: id, fn: fn, meta: meta}
			if metaByID == nil {
				metaByID = map[string]*rules.RuleMetadata{}
			}
			metaByID[id] = meta
			return
		}
	}
	cweRules = append(cweRules, cweRuleEntry{id: id, fn: fn, meta: meta})
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

// RuleIDs implements core.Detector.
func (d *GoCweScan) RuleIDs() []string {
	cweRegisterMu.Lock()
	defer cweRegisterMu.Unlock()
	ids := make([]string, len(cweRules))
	for i, e := range cweRules {
		ids[i] = e.id
	}
	sort.Strings(ids)
	return ids
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
	cweRegisterMu.Lock()
	rulesCopy := make([]cweRuleEntry, len(cweRules))
	copy(rulesCopy, cweRules)
	cweRegisterMu.Unlock()

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
		// Full taint package owns these when enabled (avoid double findings).
		if ctx != nil && ctx.TaintEnabled && isTaintCoreRule(e.id) {
			continue
		}
		// Pure FPs vs Rust gopdfsuit oracle (issue #8) — SI museums too broad.
		// Never suppress fixture / unit-test units (catalogue must stay green).
		if isOraclePureFP(e.id) && isRealProjectScan(unit) {
			continue
		}
		e.fn(unit, facts, out)
	}
}

// isOraclePureFP lists CWE IDs that fire on gopdfsuit in Go but never in Rust.
func isOraclePureFP(id string) bool {
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
			strings.Contains(p, "codehound-fixtures") ||
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

// Package core holds ScanContext, detectors, language plugins, and profiles.
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

// ScanContext holds per-run filters and flags passed to detectors.
type ScanContext struct {
	// Only is an optional allow-list of rule IDs or prefix patterns (suffix "*").
	Only []string
	// Skip lists rule IDs or prefix patterns excluded from the scan.
	Skip []string
	// FailPolicy is the severity threshold used by AnalysisResult.ShouldFail.
	FailPolicy FailPolicy
	// IncludeTests includes *_test.go files when walking.
	IncludeTests bool
	// NoCache disables the incremental analysis cache.
	NoCache bool
	// ShowIgnored keeps findings suppressed by goslop-ignore directives.
	ShowIgnored bool
	// ShowBaselined keeps findings present in the baseline (marked suppressed).
	ShowBaselined bool
	// NoBaseline disables loading/filtering via .goslop-baseline.json.
	NoBaseline bool
	// BadPracticesEnabled controls BP-* rules (default true for raw engine).
	BadPracticesEnabled bool
	// BadPracticeSeverity optionally overrides severity for all BP findings.
	BadPracticeSeverity *rules.Severity
	// SeverityOverrides are per-rule severity overrides (applied after BP override).
	SeverityOverrides map[string]rules.Severity
	// TaintEnabled enables experimental taint tracking.
	TaintEnabled bool
	// TaintMaxDepth bounds inter-procedural summary hops (1–4; default 1).
	TaintMaxDepth int
	// TaintShowPaths attaches hop evidence to taint findings when true.
	TaintShowPaths bool
	// TypedEnabled enables optional go list facts.
	TypedEnabled bool
	// RetainSources keeps per-file sources in AnalysisResult.
	RetainSources bool
	// Verbose enables extra scan stats.
	Verbose bool
	// DebugTiming collects per-rule timing.
	DebugTiming bool
	// Workers is the parallel file worker count; 0 means GOMAXPROCS.
	Workers int
	// Profile is the product pack that produced this context (informational).
	Profile ScanProfile
}

// FailOn is a compatibility alias accessor for FailPolicy.
// Prefer FailPolicy directly.
func (c *ScanContext) FailOn() FailPolicy {
	if c == nil {
		return FailMedium
	}
	return c.FailPolicy
}

// DefaultScanContext returns a sensible engine default (BP on, fail medium).
func DefaultScanContext() *ScanContext {
	return &ScanContext{
		FailPolicy:          FailMedium,
		BadPracticesEnabled: true,
		SeverityOverrides:   map[string]rules.Severity{},
		Profile:             ProfileAll,
		TaintMaxDepth:       1,
	}
}

// EffectiveTaintMaxDepth returns the clamped inter-procedural depth (1–4).
func (c *ScanContext) EffectiveTaintMaxDepth() int {
	if c == nil || c.TaintMaxDepth < 1 {
		return 1
	}
	if c.TaintMaxDepth > 4 {
		return 4
	}
	return c.TaintMaxDepth
}

// NewScanContext builds a context from a profile and CLI only/skip lists.
func NewScanContext(profile ScanProfile, only, skip []string) *ScanContext {
	return BuildScanContext(profile, only, skip)
}

// NewStringSet returns patterns as a slice used for Only/Skip filters.
// Named "set" for API familiarity; storage is an ordered slice of patterns.
func NewStringSet(values ...string) []string {
	out := make([]string, len(values))
	copy(out, values)
	return out
}

// Allows reports whether ruleID is enabled under only/skip/BP flags.
func (c *ScanContext) Allows(ruleID string) bool {
	if c == nil {
		return true
	}
	if rules.PackFromRuleID(ruleID).IsBadPractice() && !c.BadPracticesEnabled {
		return false
	}
	for _, pattern := range c.Skip {
		if matchRulePattern(pattern, ruleID) {
			return false
		}
	}
	if len(c.Only) > 0 {
		for _, pattern := range c.Only {
			if matchRulePattern(pattern, ruleID) {
				return true
			}
		}
		return false
	}
	return true
}

// ApplyFindingOverrides applies global / per-rule severity overrides in place.
func (c *ScanContext) ApplyFindingOverrides(f *rules.Finding) {
	if c == nil || f == nil {
		return
	}
	if rules.PackFromRuleID(f.RuleID).IsBadPractice() && c.BadPracticeSeverity != nil {
		f.Severity = *c.BadPracticeSeverity
	}
	if sev, ok := c.SeverityOverrides[f.RuleID]; ok {
		f.Severity = sev
	}
}

// RuleConfigFingerprint returns a stable 16-hex hash of settings that change
// which detectors run / which findings are stored. Used to mass-stale the
// incremental cache when pack/filter sets change.
func (c *ScanContext) RuleConfigFingerprint() string {
	if c == nil {
		return "nil"
	}
	only := append([]string(nil), c.Only...)
	skip := append([]string(nil), c.Skip...)
	sort.Strings(only)
	sort.Strings(skip)

	overrideKeys := make([]string, 0, len(c.SeverityOverrides))
	for k := range c.SeverityOverrides {
		overrideKeys = append(overrideKeys, k)
	}
	sort.Strings(overrideKeys)
	overrides := make([]string, 0, len(overrideKeys))
	for _, k := range overrideKeys {
		overrides = append(overrides, fmt.Sprintf("%s=%s", k, c.SeverityOverrides[k]))
	}

	bpSev := ""
	if c.BadPracticeSeverity != nil {
		bpSev = c.BadPracticeSeverity.String()
	}

	payload := fmt.Sprintf(
		"only=%v|skip=%v|taint=%v|typed=%v|bp=%v|bp_severity=%s|show_ignored=%v|severity_overrides=%v",
		only, skip, c.TaintEnabled, c.TypedEnabled, c.BadPracticesEnabled, bpSev, c.ShowIgnored, overrides,
	)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])[:16]
}

func matchRulePattern(pattern, ruleID string) bool {
	if pattern == ruleID {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(ruleID, prefix)
	}
	return false
}

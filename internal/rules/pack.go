package rules

import "strings"

// RulePack is the product / catalog pack a rule belongs to.
type RulePack int

const (
	// PackBadPractice is bad-practice / style heuristics (BP-*).
	PackBadPractice RulePack = iota
	// PackPerformance is performance heuristics (PERF-*).
	PackPerformance
	// PackSecurity is security / CWE heuristics (CWE-*).
	PackSecurity
	// PackGeneral is uncategorized or language-specific rules.
	PackGeneral
)

// PackFromRuleID classifies a rule id into a pack by prefix.
func PackFromRuleID(ruleID string) RulePack {
	switch {
	case strings.HasPrefix(ruleID, "BP-"):
		return PackBadPractice
	case strings.HasPrefix(ruleID, "PERF-"):
		return PackPerformance
	case strings.HasPrefix(ruleID, "CWE-"):
		return PackSecurity
	default:
		return PackGeneral
	}
}

// FromRuleID is an alias for PackFromRuleID.
func FromRuleID(ruleID string) RulePack { return PackFromRuleID(ruleID) }

// CategoryStr returns a coarse category string for reporters.
func (p RulePack) CategoryStr() string {
	switch p {
	case PackBadPractice:
		return "bad_practice"
	case PackPerformance:
		return "performance"
	case PackSecurity:
		return "security"
	default:
		return "general"
	}
}

// String returns a stable snake_case pack name.
func (p RulePack) String() string { return p.CategoryStr() }

// IsBadPractice reports whether this pack is the bad-practice family.
func (p RulePack) IsBadPractice() bool { return p == PackBadPractice }

// OnlyGlob returns the "all rules in this pack" allow/skip pattern, if any.
func (p RulePack) OnlyGlob() string {
	switch p {
	case PackBadPractice:
		return "BP-*"
	case PackPerformance:
		return "PERF-*"
	case PackSecurity:
		return "CWE-*"
	default:
		return ""
	}
}

// TaintCoreCWERules are the taint-core CWE ids shared by recommended and security packs.
var TaintCoreCWERules = []string{
	"CWE-22", "CWE-78", "CWE-79", "CWE-89", "CWE-90", "CWE-91",
}

// PerfTierSRules are S-tier PERF rule ids for the recommended pack.
// PERF-116 is included so the seed detector participates in recommended scans.
var PerfTierSRules = []string{
	"PERF-1", "PERF-7", "PERF-50", "PERF-58", "PERF-71",
	"PERF-101", "PERF-103", "PERF-116", "PERF-189", "PERF-190",
}

// PerfTierARules are A-tier PERF rule ids for the perf pack.
var PerfTierARules = []string{
	"PERF-11", "PERF-12", "PERF-22", "PERF-31", "PERF-82", "PERF-85",
	"PERF-142", "PERF-143", "PERF-164", "PERF-183", "PERF-210", "PERF-213",
}

// SecurityPackRules are taint core + high-value structural neighbors.
var SecurityPackRules = []string{
	"CWE-22", "CWE-41", "CWE-59", "CWE-78", "CWE-79",
	"CWE-89", "CWE-90", "CWE-91", "CWE-93",
}

// StylePackPatterns is the style-pack allow pattern (all bad-practice rules).
var StylePackPatterns = []string{"BP-*"}

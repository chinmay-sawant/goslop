package rules

import "github.com/chinmay/goslop/internal/cwe"

// Maturity is a catalogue trust / readiness tag for a rule.
type Maturity string

const (
	// MaturityProduction is CI-certified / default-pack eligible.
	MaturityProduction Maturity = "production"
	// MaturityExperimental is usable but not CI-hard-fail certified.
	MaturityExperimental Maturity = "experimental"
	// MaturityFixtureOnly is museum / fixture corpus; opt-in via --profile all or --only.
	MaturityFixtureOnly Maturity = "fixture-only"
	// MaturityQuarantined is reserved or incomplete; available under all only.
	MaturityQuarantined Maturity = "quarantined"
)

// RuleMetadata is catalogue metadata for a single rule id.
type RuleMetadata struct {
	ID          string
	Title       string
	Description string
	Severity    Severity
	CWE         []cwe.CweRef
	Fix         string
	Pack        RulePack
	// Maturity defaults to production when empty; see InferMaturity.
	Maturity Maturity
	// QuarantineReason explains fixture-only / quarantined status when set.
	QuarantineReason string
}

// Meta constructs RuleMetadata and classifies pack from the rule id.
func Meta(id, title, description string, sev Severity, cwes []cwe.CweRef, fix string) RuleMetadata {
	m := RuleMetadata{
		ID:          id,
		Title:       title,
		Description: description,
		Severity:    sev,
		CWE:         cwes,
		Fix:         fix,
		Pack:        PackFromRuleID(id),
	}
	m.Maturity = InferMaturity(id)
	if m.Maturity == MaturityQuarantined {
		m.QuarantineReason = "reserved / incomplete rule id"
	}
	if m.Maturity == MaturityFixtureOnly {
		m.QuarantineReason = "fixture-only museum entry; not in default CI packs"
	}
	return m
}

// EffectiveMaturity returns Maturity or a inferred default.
func (m *RuleMetadata) EffectiveMaturity() Maturity {
	if m == nil {
		return MaturityProduction
	}
	if m.Maturity != "" {
		return m.Maturity
	}
	return InferMaturity(m.ID)
}

// InferMaturity returns a coarse maturity tag from rule id conventions.
// Full tables can override via RuleMetadata.Maturity later.
func InferMaturity(ruleID string) Maturity {
	switch ruleID {
	case "BP-63":
		return MaturityQuarantined
	case "CWE-22", "CWE-78", "CWE-79", "CWE-89", "CWE-90", "CWE-91":
		return MaturityExperimental // taint-core; name-string honesty
	}
	// PERF S-tier: production; other PERF: experimental until pack-certified.
	if PackFromRuleID(ruleID) == PackPerformance {
		for _, id := range PerfTierSRules {
			if id == ruleID {
				return MaturityProduction
			}
		}
		return MaturityExperimental
	}
	// Structural CWE museums (high id / long-tail) stay fixture-only-ish.
	if PackFromRuleID(ruleID) == PackSecurity {
		return MaturityFixtureOnly
	}
	return MaturityProduction
}

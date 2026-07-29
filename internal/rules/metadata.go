package rules

import "github.com/chinmay/codehound/internal/cwe"

// RuleMetadata is catalogue metadata for a single rule id.
type RuleMetadata struct {
	ID          string
	Title       string
	Description string
	Severity    Severity
	CWE         []cwe.CweRef
	Fix         string
	Pack        RulePack
}

// Meta constructs RuleMetadata and classifies pack from the rule id.
func Meta(id, title, description string, sev Severity, cwes []cwe.CweRef, fix string) RuleMetadata {
	return RuleMetadata{
		ID:          id,
		Title:       title,
		Description: description,
		Severity:    sev,
		CWE:         cwes,
		Fix:         fix,
		Pack:        PackFromRuleID(id),
	}
}

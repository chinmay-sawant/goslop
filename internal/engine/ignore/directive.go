// Package ignore parses and applies codehound-ignore suppression directives.
package ignore

// Directive is a rule allow-list or “all rules”.
type Directive struct {
	// RuleIDs is nil when every rule is suppressed.
	RuleIDs []string
}

// All returns a directive that suppresses every rule.
func All() Directive {
	return Directive{RuleIDs: nil}
}

// Rules returns a directive that suppresses only the listed rule IDs.
func Rules(ids ...string) Directive {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != "" {
			out = append(out, id)
		}
	}
	return Directive{RuleIDs: out}
}

// Matches reports whether this directive suppresses ruleID.
func (d Directive) Matches(ruleID string) bool {
	if d.RuleIDs == nil {
		return true
	}
	for _, id := range d.RuleIDs {
		if id == ruleID {
			return true
		}
	}
	return false
}

// IsAll reports whether every rule is suppressed.
func (d Directive) IsAll() bool {
	return d.RuleIDs == nil
}

// Merge unions other into d (all wins).
func (d *Directive) Merge(other Directive) {
	if d.RuleIDs == nil || other.RuleIDs == nil {
		d.RuleIDs = nil
		return
	}
	for _, id := range other.RuleIDs {
		found := false
		for _, existing := range d.RuleIDs {
			if existing == id {
				found = true
				break
			}
		}
		if !found {
			d.RuleIDs = append(d.RuleIDs, id)
		}
	}
}

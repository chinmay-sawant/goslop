package cwe_test

import "testing"

func TestCWEMassAssignmentFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-914", "CWE-915", "CWE-916"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

package cwe_test

import "testing"

func TestCWECodeDynamicFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-214", "CWE-215", "CWE-695", "CWE-749", "CWE-829"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

package cwe_test

import "testing"

func TestInjectionExpansionFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-88", "CWE-90", "CWE-91", "CWE-93", "CWE-94", "CWE-117"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

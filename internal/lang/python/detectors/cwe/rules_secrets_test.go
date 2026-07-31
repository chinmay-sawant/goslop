package cwe_test

import "testing"

func TestCWESecretsFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-256", "CWE-260", "CWE-261", "CWE-312", "CWE-319", "CWE-523", "CWE-547", "CWE-798"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

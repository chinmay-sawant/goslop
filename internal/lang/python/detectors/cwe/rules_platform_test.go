package cwe_test

import "testing"

func TestPlatformFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{
		"CWE-396",
		"CWE-397",
		"CWE-478",
		"CWE-252",
		"CWE-390",
		"CWE-584",
	} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

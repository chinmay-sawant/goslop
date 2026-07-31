package cwe_test

import "testing"

func TestWebConfigFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{
		"CWE-756",
		"CWE-489",
		"CWE-15",
		"CWE-1051",
		"CWE-1052",
		"CWE-1125",
		"CWE-1188",
		"CWE-921",
	} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

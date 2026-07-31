package cwe_test

import "testing"

func TestCWEValidationFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{
		"CWE-1173",
		"CWE-1230",
		"CWE-1236",
		"CWE-1286",
		"CWE-1289",
		"CWE-1333",
		"CWE-1389",
		"CWE-140",
	} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

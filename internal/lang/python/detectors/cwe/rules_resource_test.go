package cwe_test

import "testing"

func TestCWEResourceFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{
		"CWE-434",
		"CWE-427",
		"CWE-379",
		"CWE-459",
		"CWE-772",
		"CWE-770",
		"CWE-708",
		"CWE-477",
	} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

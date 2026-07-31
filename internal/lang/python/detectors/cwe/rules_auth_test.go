package cwe_test

import "testing"

func TestCWEAuthFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{
		"CWE-306",
		"CWE-307",
		"CWE-346",
		"CWE-359",
		"CWE-613",
		"CWE-565",
		"CWE-807",
		"CWE-698",
	} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

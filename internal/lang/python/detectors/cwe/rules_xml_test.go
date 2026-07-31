package cwe_test

import "testing"

func TestCWEXMLFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-112", "CWE-611", "CWE-776"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

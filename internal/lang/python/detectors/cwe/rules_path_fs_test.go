package cwe_test

import "testing"

func TestCWEPathFSFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-41", "CWE-59", "CWE-73", "CWE-250", "CWE-276", "CWE-378", "CWE-426", "CWE-494"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

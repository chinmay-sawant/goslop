package cwe_test

import "testing"

func TestCWECryptoFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-295", "CWE-328", "CWE-335", "CWE-338", "CWE-347", "CWE-1204", "CWE-1240", "CWE-1241", "CWE-1392"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

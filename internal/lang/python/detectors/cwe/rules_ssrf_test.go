package cwe_test

import "testing"

func TestSSRFRedirectAndChannelFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-601", "CWE-605", "CWE-918", "CWE-924", "CWE-940", "CWE-941"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

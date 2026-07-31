package cwe_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe"
)

func TestCoreCWEFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"CWE-22", "CWE-78", "CWE-79", "CWE-89", "CWE-502"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

func TestBatchIntegrityNoDoubleIDs(t *testing.T) {
	t.Parallel()
	ids := cwe.RegisteredRuleIDs()
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate rule id %s", id)
		}
		seen[id] = true
	}
	if len(ids) != 159 {
		t.Fatalf("ids = %v (len %d), want 159", ids, len(ids))
	}
}

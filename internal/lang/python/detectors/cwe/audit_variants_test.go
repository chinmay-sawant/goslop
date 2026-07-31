package cwe_test

import "testing"

func TestCWEFalsePositiveAuditFixtureVariants(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		rule     string
		caseName string
	}{
		{rule: "CWE-89", caseName: "CWE-89-orm-expression"},
		{rule: "CWE-93", caseName: "CWE-93-numeric-header"},
		{rule: "CWE-312", caseName: "CWE-312-test-fixture"},
		{rule: "CWE-396", caseName: "CWE-396-test-collection"},
		{rule: "CWE-798", caseName: "CWE-798-test-fixture"},
		{rule: "CWE-1046", caseName: "CWE-1046-counter-loop"},
		{rule: "CWE-1124", caseName: "CWE-1124-data-layout"},
	} {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			t.Parallel()
			assertCWEFixtureCase(t, tc.rule, tc.caseName)
		})
	}
}

package cwe_test

import "testing"

// Path/injection/info residual FP fixtures for Project_Parva-shaped guards.
// Kept in a dedicated file so concurrent edits to audit_variants_test.go
// (CWE-396 wave) do not collide with appends for these rules.
func TestCWEPathMiscAuditFixtureVariants(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		rule     string
		caseName string
	}{
		{rule: "CWE-93", caseName: "CWE-93-literal-ternary"},
		{rule: "CWE-93", caseName: "CWE-93-link-merge"},
		{rule: "CWE-73", caseName: "CWE-73-cli-main-argv"},
		{rule: "CWE-88", caseName: "CWE-88-offline-tools"},
		{rule: "CWE-22", caseName: "CWE-22-benchmark-report"},
		{rule: "CWE-215", caseName: "CWE-215-release-script"},
	} {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			t.Parallel()
			assertCWEFixtureCase(t, tc.rule, tc.caseName)
		})
	}
}

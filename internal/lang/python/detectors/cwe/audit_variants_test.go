package cwe_test

import "testing"

func TestCWEFalsePositiveAuditFixtureVariants(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		rule     string
		caseName string
	}{
		{rule: "CWE-22", caseName: "CWE-22-no-restricted-join"},
		{rule: "CWE-73", caseName: "CWE-73-main-cli-outdir"},
		{rule: "CWE-88", caseName: "CWE-88-fixed-argv"},
		{rule: "CWE-89", caseName: "CWE-89-orm-expression"},
		{rule: "CWE-89", caseName: "CWE-89-orm-direct"},
		{rule: "CWE-89", caseName: "CWE-89-execute-wrapper"},
		{rule: "CWE-89", caseName: "CWE-89-adjacent-literal"},
		{rule: "CWE-89", caseName: "CWE-89-sqlalchemy-text"},
		{rule: "CWE-89", caseName: "CWE-89-non-sql-execute"},
		{rule: "CWE-89", caseName: "CWE-89-def-execute"},
		{rule: "CWE-90", caseName: "CWE-90-regex-search"},
		{rule: "CWE-91", caseName: "CWE-91-fromstring-parse"},
		{rule: "CWE-93", caseName: "CWE-93-numeric-header"},
		{rule: "CWE-94", caseName: "CWE-94-attr-method-exec"},
		{rule: "CWE-256", caseName: "CWE-256-fixture-password"},
		{rule: "CWE-312", caseName: "CWE-312-benchmark-secret"},
		{rule: "CWE-312", caseName: "CWE-312-test-fixture"},
		{rule: "CWE-328", caseName: "CWE-328-fingerprint"},
		{rule: "CWE-396", caseName: "CWE-396-test-collection"},
		{rule: "CWE-396", caseName: "CWE-396-batch3-reraise-from"},
		{rule: "CWE-396", caseName: "CWE-396-batch3-exc-info"},
		{rule: "CWE-396", caseName: "CWE-396-batch3-set-exception"},
		{rule: "CWE-396", caseName: "CWE-396-batch3-result-error"},
		{rule: "CWE-798", caseName: "CWE-798-test-fixture"},
		{rule: "CWE-798", caseName: "CWE-798-benchmark-secret"},
		{rule: "CWE-798", caseName: "CWE-798-fixture-password"},
		{rule: "CWE-1046", caseName: "CWE-1046-counter-loop"},
		{rule: "CWE-1046", caseName: "CWE-1046-bytearray-loop"},
		{rule: "CWE-1124", caseName: "CWE-1124-data-layout"},
		{rule: "CWE-1124", caseName: "CWE-1124-declaration-scope"},
		{rule: "CWE-924", caseName: "CWE-924-authenticated-route"},
		{rule: "CWE-88", caseName: "CWE-88-constant-argv"},
		{rule: "CWE-93", caseName: "CWE-93-header-read"},
		{rule: "CWE-117", caseName: "CWE-117-constant-log"},
		{rule: "CWE-1341", caseName: "CWE-1341-distinct-handles"},
		{rule: "CWE-215", caseName: "CWE-215-literal-password-word"},
		{rule: "CWE-367", caseName: "CWE-367-different-path"},
		{rule: "CWE-502", caseName: "CWE-502-ruamel"},
	} {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			t.Parallel()
			assertCWEFixtureCase(t, tc.rule, tc.caseName)
		})
	}
}

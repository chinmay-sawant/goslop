package cwe_test

import "testing"

func TestCWEInfoExposureFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{
		"CWE-201",
		"CWE-204",
		"CWE-208",
		"CWE-209",
		"CWE-212",
		"CWE-213",
		"CWE-488",
		"CWE-497",
	} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

func TestCWE208TestModuleAssertSkipped(t *testing.T) {
	t.Parallel()
	assertCWEFixtureCase(t, "CWE-208", "CWE-208-test-assert")
}

func TestCWE208AuthorityEnumSkipped(t *testing.T) {
	t.Parallel()
	assertCWEFixtureCase(t, "CWE-208", "CWE-208-authority-enum")
}

package cwe_test

import "testing"

// Source examples live exclusively in tests/fixtures/python/cwe. This keeps
// unit and integration coverage on exactly the same executable corpus.
func TestCWETierBFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{
		"CWE-66", "CWE-76", "CWE-178", "CWE-179", "CWE-182", "CWE-184", "CWE-186", "CWE-257", "CWE-272", "CWE-279", "CWE-289", "CWE-290", "CWE-323", "CWE-331", "CWE-334",
		"CWE-367", "CWE-403", "CWE-409", "CWE-454", "CWE-472", "CWE-521", "CWE-524", "CWE-538", "CWE-552", "CWE-617", "CWE-641", "CWE-648", "CWE-779", "CWE-836", "CWE-838",
		"CWE-908", "CWE-909", "CWE-910", "CWE-911", "CWE-920", "CWE-939", "CWE-1007", "CWE-1021", "CWE-1046", "CWE-1050", "CWE-1060", "CWE-1067", "CWE-1071", "CWE-1072", "CWE-1084",
		"CWE-1104", "CWE-1106", "CWE-1108", "CWE-1121", "CWE-1123", "CWE-1124", "CWE-1220", "CWE-1265", "CWE-1284", "CWE-1285", "CWE-1287", "CWE-1288", "CWE-1322", "CWE-1339", "CWE-1341",
	} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			assertCWEFixturePair(t, rule)
		})
	}
}

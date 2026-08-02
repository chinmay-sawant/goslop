package perf_test

import "testing"

func TestRuntimeRulesHitAndMiss(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		rule     string
		caseName string
	}{
		{"PERF-PY-5", "PERF-PY-5"},
		{"PERF-PY-7", "PERF-PY-7"},
		{"PERF-PY-17", "PERF-PY-17"},
		{"PERF-PY-18", "PERF-PY-18"},
		{"PERF-PY-22", "PERF-PY-22"},
	} {
		tc := tc
		t.Run(tc.rule, func(t *testing.T) {
			t.Parallel()
			assertPERFFixtureCase(t, tc.rule, tc.caseName)
		})
	}
}

func TestRuntimeConfigRulesSuppressTestsAndLocalConfigurations(t *testing.T) {
	t.Parallel()
	assertPERFFixtureCase(t, "PERF-PY-17", "PERF-PY-17-test-path")
	assertPERFFixtureCase(t, "PERF-PY-17", "PERF-PY-17-local-config")
	assertPERFFixtureCase(t, "PERF-PY-22", "PERF-PY-22-test-path")
	assertPERFFixtureCase(t, "PERF-PY-22", "PERF-PY-22-local-config")
}

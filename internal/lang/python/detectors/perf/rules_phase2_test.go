package perf_test

import "testing"

func TestPhase2Rules(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		rule     string
		caseName string
	}{
		{"materialized sort / database aggregate", "PERF-PY-1", "PERF-PY-1"},
		{"batch create / dependent create", "PERF-PY-3", "PERF-PY-3"},
		{"counter save / atomic counter", "PERF-PY-4", "PERF-PY-4"},
		{"parse then dump / raw payload", "PERF-PY-9", "PERF-PY-9"},
		{"dump for log unrelated create", "PERF-PY-9", "PERF-PY-9-dump-for-log"},
		{"sleep after work / continue after work", "PERF-PY-10", "PERF-PY-10"},
		{"per row mutation / set update", "PERF-PY-11", "PERF-PY-11"},
		{"sqlalchemy per row / non ORM save", "PERF-PY-11", "PERF-PY-11-sqlalchemy"},
		{"unbounded json / limited json", "PERF-PY-12", "PERF-PY-12"},
		{"lookup then create / upsert", "PERF-PY-14", "PERF-PY-14"},
		{"helper lookup then insert", "PERF-PY-14", "PERF-PY-14-helper-lookup"},
		{"top-level lookup after function", "PERF-PY-14", "PERF-PY-14-top-level"},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertPERFFixtureCase(t, tc.rule, tc.caseName)
		})
	}
}

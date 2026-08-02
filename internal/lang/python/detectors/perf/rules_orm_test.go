package perf_test

import "testing"

// The detector Run boundary is the public seam for these source-only rules.
func TestORMRules(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		rule     string
		caseName string
	}{
		{"lookup in loop", "PERF-PY-2", "PERF-PY-2"},
		{"work claim no lock", "PERF-PY-6", "PERF-PY-6"},
		{"lazy sqlalchemy relation", "PERF-PY-8", "PERF-PY-8"},
		{"full hydration projection", "PERF-PY-13", "PERF-PY-13"},
		{"composite filter visible no index", "PERF-PY-15", "PERF-PY-15"},
		{"retention predicate visible no index", "PERF-PY-16", "PERF-PY-16"},
		{"unbounded lock sweep", "PERF-PY-19", "PERF-PY-19"},
		{"sqlalchemy unbounded lock sweep", "PERF-PY-19", "PERF-PY-19-sqlalchemy"},
		{"sort visible no index", "PERF-PY-20", "PERF-PY-20"},
		{"unbounded maintenance delete", "PERF-PY-21", "PERF-PY-21"},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertPERFFixtureCase(t, tc.rule, tc.caseName)
		})
	}
}

func TestIndexRulesSuppressWhenModelDeclarationIsUnavailable(t *testing.T) {
	t.Parallel()
	assertPERFFixtureCase(t, "PERF-PY-15", "PERF-PY-15-no-model")
	assertPERFFixtureCase(t, "PERF-PY-16", "PERF-PY-16-no-model")
	assertPERFFixtureCase(t, "PERF-PY-20", "PERF-PY-20-no-model")
	assertPERFFixtureCase(t, "PERF-PY-15", "PERF-PY-15-external-migration")
}

func TestLookupInSmallExplicitLoopIsSuppressed(t *testing.T) {
	t.Parallel()
	assertPERFFixtureCase(t, "PERF-PY-2", "PERF-PY-2-small-loop")
}

func TestPERFPY6RequiresBoundStatusMutation(t *testing.T) {
	t.Parallel()
	assertPERFFixtureCase(t, "PERF-PY-6", "PERF-PY-6-unbound-mutation")
}

func TestPERFPY20PKSuppressAndUnrelatedIndex(t *testing.T) {
	t.Parallel()
	assertPERFFixtureCase(t, "PERF-PY-20", "PERF-PY-20-pk-filter")
	assertPERFFixtureCase(t, "PERF-PY-20", "PERF-PY-20-unrelated-index")
}

func TestPERFPY14RequiresSharedKey(t *testing.T) {
	t.Parallel()
	assertPERFFixtureCase(t, "PERF-PY-14", "PERF-PY-14-shared-key")
}

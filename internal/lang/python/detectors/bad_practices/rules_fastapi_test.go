package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY29MutableGlobal(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-29", "BP-PY-29")
	assertBPFixtureCase(t, "BP-PY-29", "BP-PY-29-cache-mut")
	assertBPFixtureCase(t, "BP-PY-29", "BP-PY-29-flaskish")
}

func TestBPPY30BlockingIOAsyncRoute(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-30", "BP-PY-30")
	assertBPFixtureCase(t, "BP-PY-30", "BP-PY-30-requests")
	assertBPFixtureCase(t, "BP-PY-30", "BP-PY-30-sync-route")

	vuln := loadBPFixture(t, "BP-PY-30", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-30" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-30 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY31ResponseModel(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-31", "BP-PY-31")
	assertBPFixtureCase(t, "BP-PY-31", "BP-PY-31-response-model")
	assertBPFixtureCase(t, "BP-PY-31", "BP-PY-31-dict")
}

func TestBPPY32FileResponseUserPath(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-32", "BP-PY-32")
	assertBPFixtureCase(t, "BP-PY-32", "BP-PY-32-fstring")
	assertBPFixtureCase(t, "BP-PY-32", "BP-PY-32-confined-basename")
	assertBPFixtureCase(t, "BP-PY-32", "BP-PY-32-constant-dir")

	vuln := loadBPFixture(t, "BP-PY-32", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-32" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-32 severity = %v, want high", f.Severity)
		}
	}
}

func TestFastAPIRulesRegistered(t *testing.T) {
	t.Parallel()
	for _, id := range []string{"BP-PY-29", "BP-PY-30", "BP-PY-31", "BP-PY-32"} {
		safe := loadBPFixture(t, id, false)
		assertRule(t, id, safe.path, safe.body, false)
	}
}

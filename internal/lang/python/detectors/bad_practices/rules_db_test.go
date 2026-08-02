package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY35SQLAlchemyTextFString(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-35", "BP-PY-35")
	assertBPFixtureCase(t, "BP-PY-35", "BP-PY-35-format")

	vuln := loadBPFixture(t, "BP-PY-35", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-35" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-35 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY36SessionNotClosed(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-36", "BP-PY-36")
	assertBPFixtureCase(t, "BP-PY-36", "BP-PY-36-sessionlocal")
	assertBPFixtureCase(t, "BP-PY-36", "BP-PY-36-close")
}

func TestBPPY37ExecutePercentFormat(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-37", "BP-PY-37")
	assertBPFixtureCase(t, "BP-PY-37", "BP-PY-37-fstring")
	assertBPFixtureCase(t, "BP-PY-37", "BP-PY-37-param-placeholder")

	vuln := loadBPFixture(t, "BP-PY-37", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-37" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-37 severity = %v, want high", f.Severity)
		}
	}
}

func TestDBAndTemplateRulesRegistered(t *testing.T) {
	t.Parallel()
	// Registration smoke: benign sources should not fire these rules.
	for _, id := range []string{"BP-PY-33", "BP-PY-35", "BP-PY-37"} {
		safe := loadBPFixture(t, id, false)
		assertRule(t, id, safe.path, safe.body, false)
	}
}

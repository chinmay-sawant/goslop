package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY33JinjaAutoescape(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-33", "BP-PY-33")
	assertBPFixtureCase(t, "BP-PY-33", "BP-PY-33-multiline")
	assertBPFixtureCase(t, "BP-PY-33", "BP-PY-33-true")

	vuln := loadBPFixture(t, "BP-PY-33", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-33" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-33 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY34MarkupSafe(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-34", "BP-PY-34")
	assertBPFixtureCase(t, "BP-PY-34", "BP-PY-34-request")
	assertBPFixtureCase(t, "BP-PY-34", "BP-PY-34-jinja-filter")
	assertBPFixtureCase(t, "BP-PY-34", "BP-PY-34-literal")
}

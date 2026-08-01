package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY38CreateTaskBare(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-38", "BP-PY-38")
	assertBPFixtureCase(t, "BP-PY-38", "BP-PY-38-task-append")
	assertBPFixtureCase(t, "BP-PY-38", "BP-PY-38-ensure-future")
}

func TestBPPY39TimeSleepInAsync(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-39", "BP-PY-39")
	assertBPFixtureCase(t, "BP-PY-39", "BP-PY-39-sync-sleep")

	vuln := loadBPFixture(t, "BP-PY-39", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-39" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-39 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY40ThreadingWithoutJoin(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-40", "BP-PY-40")
	assertBPFixtureCase(t, "BP-PY-40", "BP-PY-40-daemon-start")
}

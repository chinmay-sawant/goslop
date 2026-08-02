package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY46PrintInLibrary(t *testing.T) {
	t.Parallel()
	vuln := loadBPFixture(t, "BP-PY-46", true)
	safeMain := loadBPFixture(t, "BP-PY-46-main-guard", false)
	safeLogging := loadBPFixture(t, "BP-PY-46", false)
	assertRule(t, "BP-PY-46", vuln.path, vuln.body, true)
	assertRule(t, "BP-PY-46", safeMain.path, safeMain.body, false)
	assertRule(t, "BP-PY-46", safeLogging.path, safeLogging.body, false)
	// Test file skip
	testSkip := loadBPFixture(t, "BP-PY-46-test-path", false)
	assertRule(t, "BP-PY-46", testSkip.path, testSkip.body, false)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-46" && f.Severity != rules.SeverityInfo {
			t.Fatalf("BP-PY-46 severity = %v, want info", f.Severity)
		}
	}
}

func TestBPPY47EagerLogFormat(t *testing.T) {
	t.Parallel()
	vuln := loadBPFixture(t, "BP-PY-47", true)
	vulnFormat := loadBPFixture(t, "BP-PY-47-format", true)
	vulnLogging := loadBPFixture(t, "BP-PY-47-logging-module", true)
	safe := loadBPFixture(t, "BP-PY-47", false)
	safeLogging := loadBPFixture(t, "BP-PY-47-logging-module", false)
	assertRule(t, "BP-PY-47", vuln.path, vuln.body, true)
	assertRule(t, "BP-PY-47", vulnFormat.path, vulnFormat.body, true)
	assertRule(t, "BP-PY-47", vulnLogging.path, vulnLogging.body, true)
	assertRule(t, "BP-PY-47", safe.path, safe.body, false)
	assertRule(t, "BP-PY-47", safeLogging.path, safeLogging.body, false)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-47" && f.Severity != rules.SeverityInfo {
			t.Fatalf("BP-PY-47 severity = %v, want info", f.Severity)
		}
	}
}

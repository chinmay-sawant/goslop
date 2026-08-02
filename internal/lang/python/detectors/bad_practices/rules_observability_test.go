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
	// Script-path / CLI decorator / docstring masking variants
	for _, caseName := range []string{
		"BP-PY-46-script-path",
		"BP-PY-46-examples-library-print",
		"BP-PY-46-click-cli",
		"BP-PY-46-typer-cli",
		"BP-PY-46-cyclopts-cli",
		"BP-PY-46-docstring-print",
		"BP-PY-46-main-epilog",
		"BP-PY-46-shebang-script",
		"BP-PY-46-shebang-vestigial-library",
		"BP-PY-46-setup-script",
		"BP-PY-46-string-template",
		"BP-PY-46-cli-module",
		"BP-PY-46-commands-path",
		"BP-PY-46-rich-print",
		"BP-PY-46-console-method",
		"BP-PY-46-demo-script-examples",
		"BP-PY-46-demo-guard-local",
		"BP-PY-46-script-completion",
		"BP-PY-46-standalone-module-print",
	} {
		vulnCase := loadBPFixture(t, caseName, true)
		safeCase := loadBPFixture(t, caseName, false)
		assertRule(t, "BP-PY-46", vulnCase.path, vulnCase.body, true)
		assertRule(t, "BP-PY-46", safeCase.path, safeCase.body, false)
	}
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
	// f-string logs in tests / constant f-strings without interpolation
	for _, caseName := range []string{
		"BP-PY-47-test-path",
		"BP-PY-47-constant-fstring",
	} {
		safeCase := loadBPFixture(t, caseName, false)
		vulnCase := loadBPFixture(t, caseName, true)
		assertRule(t, "BP-PY-47", safeCase.path, safeCase.body, false)
		assertRule(t, "BP-PY-47", vulnCase.path, vulnCase.body, true)
	}
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-47" && f.Severity != rules.SeverityInfo {
			t.Fatalf("BP-PY-47 severity = %v, want info", f.Severity)
		}
	}
}

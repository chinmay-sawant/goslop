package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY46PrintInLibrary(t *testing.T) {
	t.Parallel()
	vuln := "def f():\n    print('debug')\n"
	safeMain := "if __name__ == '__main__':\n    print('cli')\n"
	safeLogging := "import logging\ndef f():\n    logging.info('debug')\n"
	assertRule(t, "BP-PY-46", "lib.py", vuln, true)
	assertRule(t, "BP-PY-46", "lib.py", safeMain, false)
	assertRule(t, "BP-PY-46", "lib.py", safeLogging, false)
	// Test file skip
	assertRule(t, "BP-PY-46", "test_lib.py", vuln, false)
	findings := runBP(t, nil, vuln, "lib.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-46" && f.Severity != rules.SeverityInfo {
			t.Fatalf("BP-PY-46 severity = %v, want info", f.Severity)
		}
	}
}

func TestBPPY47EagerLogFormat(t *testing.T) {
	t.Parallel()
	vuln := "logger.info(f\"user={user}\")\n"
	vulnFormat := "logger.debug(\"user={}\".format(user))\n"
	vulnLogging := "import logging\nlogging.info(f\"x={x}\")\n"
	safe := "logger.info(\"user=%s\", user)\n"
	safeLogging := "import logging\nlogging.info(\"user=%s\", user)\n"
	assertRule(t, "BP-PY-47", "log.py", vuln, true)
	assertRule(t, "BP-PY-47", "log.py", vulnFormat, true)
	assertRule(t, "BP-PY-47", "log.py", vulnLogging, true)
	assertRule(t, "BP-PY-47", "log.py", safe, false)
	assertRule(t, "BP-PY-47", "log.py", safeLogging, false)
	findings := runBP(t, nil, vuln, "log.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-47" && f.Severity != rules.SeverityInfo {
			t.Fatalf("BP-PY-47 severity = %v, want info", f.Severity)
		}
	}
}

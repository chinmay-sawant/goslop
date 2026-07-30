package badpractices_test

import (
	"testing"
)

func TestBPPY43RequirementsWithoutPins(t *testing.T) {
	t.Parallel()
	// Path A: path-gated on requirements*.txt; unit is LanguagePython for unit tests.
	vuln := "requests\nflask\n"
	safe := "requests==2.31.0\nflask>=2.0\n"
	safePinned := "django~=4.2\n"
	assertRule(t, "BP-PY-43", "requirements.txt", vuln, true)
	assertRule(t, "BP-PY-43", "requirements.txt", safe, false)
	assertRule(t, "BP-PY-43", "requirements.txt", safePinned, false)
	// Directives / comments miss
	assertRule(t, "BP-PY-43", "requirements.txt", "-r other.txt\n# comment\n-e .\n", false)
	// Non-requirements path must not fire even with bare package-like lines
	assertRule(t, "BP-PY-43", "app.py", "requests\nflask\n", false)
	// requirements-dev.txt also gated
	assertRule(t, "BP-PY-43", "requirements-dev.txt", "pytest\n", true)
}

func TestBPPY44DeprecatedStdlibImport(t *testing.T) {
	t.Parallel()
	assertRule(t, "BP-PY-44", "legacy.py", "import imp\n", true)
	assertRule(t, "BP-PY-44", "legacy.py", "from asyncore import dispatcher\n", true)
	assertRule(t, "BP-PY-44", "legacy.py", "import cgi\n", true)
	assertRule(t, "BP-PY-44", "modern.py", "import importlib\n", false)
	assertRule(t, "BP-PY-44", "modern.py", "import asyncio\n", false)
}

func TestBPPY45SysPathMutation(t *testing.T) {
	t.Parallel()
	vuln := "import sys\nsys.path.insert(0, './lib')\n"
	vulnAppend := "import sys\nsys.path.append('/opt/pkg')\n"
	assertRule(t, "BP-PY-45", "app.py", vuln, true)
	assertRule(t, "BP-PY-45", "app.py", vulnAppend, true)
	// Test file skip
	assertRule(t, "BP-PY-45", "tests/test_path.py", vuln, false)
	assertRule(t, "BP-PY-45", "test_path.py", vuln, false)
	// Read-only sys.path does not fire
	assertRule(t, "BP-PY-45", "app.py", "import sys\nprint(sys.path)\n", false)
}

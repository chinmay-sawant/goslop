package badpractices_test

import (
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	badpractices "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPRulesRegistered(t *testing.T) {
	t.Parallel()
	d := badpractices.NewPythonBadPracticeScan()
	if d.Language() != core.LanguagePython {
		t.Fatalf("Language() = %v, want LanguagePython", d.Language())
	}
	ids := d.RuleIDs()
	want := []string{
		"BP-PY-1", "BP-PY-2", "BP-PY-4", "BP-PY-6", "BP-PY-7",
		"BP-PY-8", "BP-PY-9", "BP-PY-10", "BP-PY-11", "BP-PY-12", "BP-PY-13",
		"BP-PY-16", "BP-PY-17", "BP-PY-21",
		"BP-PY-22", "BP-PY-23", "BP-PY-24", "BP-PY-25", "BP-PY-26", "BP-PY-27", "BP-PY-28",
		"BP-PY-29", "BP-PY-30", "BP-PY-31", "BP-PY-32",
		"BP-PY-33", "BP-PY-34",
		"BP-PY-35", "BP-PY-36", "BP-PY-37",
		"BP-PY-48", "BP-PY-49", "BP-PY-50",
	}
	if len(ids) < len(want) {
		t.Fatalf("expected >= %d BP-PY rules, got %d: %v", len(want), len(ids), ids)
	}
	for _, id := range want {
		if !contains(ids, id) {
			t.Errorf("missing registered rule %s", id)
		}
		if !strings.HasPrefix(id, "BP-PY-") {
			t.Errorf("id %q must use BP-PY- prefix", id)
		}
		meta := d.MetadataFor(id)
		if meta == nil {
			t.Errorf("MetadataFor(%s) = nil", id)
			continue
		}
		if meta.Pack != rules.PackBadPractice {
			t.Errorf("%s pack = %v, want PackBadPractice", id, meta.Pack)
		}
		if meta.ID != id {
			t.Errorf("meta.ID = %q, want %q", meta.ID, id)
		}
	}
	// Collision guard: no bare BP-<n>
	for _, id := range ids {
		if strings.HasPrefix(id, "BP-") && !strings.HasPrefix(id, "BP-PY-") {
			t.Errorf("registered bare Go-style id %q", id)
		}
	}
	if badpractices.CatalogueSize() < len(want) {
		t.Errorf("catalogue size too small: %d", badpractices.CatalogueSize())
	}
}

func TestFullCatalogueParseable(t *testing.T) {
	t.Parallel()
	// Implemented subset severities must match catalogue when ruleset JSON is reachable.
	if err := badpractices.ValidateImplementedMetadata(); err != nil {
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "cannot find") {
			t.Skipf("catalogue file not reachable from test CWD: %v", err)
		}
		t.Fatal(err)
	}
	d := badpractices.NewPythonBadPracticeScan()
	for _, id := range d.RuleIDs() {
		meta := badpractices.MetadataForID(id)
		if meta == nil {
			t.Fatalf("nil metadata for registered %s", id)
		}
	}
	full, err := badpractices.FullCatalogueSize()
	if err != nil {
		t.Skipf("full catalogue size: %v", err)
	}
	if full < 50 {
		t.Fatalf("full catalogue keys = %d, want >= 50", full)
	}
}

func TestBPPY1BareExcept(t *testing.T) {
	t.Parallel()
	vuln := "def f():\n    try:\n        x()\n    except:\n        pass\n"
	safe := "def f():\n    try:\n        x()\n    except ValueError as e:\n        raise RuntimeError('bad') from e\n"
	assertRule(t, "BP-PY-1", "bare.py", vuln, true)
	assertRule(t, "BP-PY-1", "bare.py", safe, false)
}

func TestBPPY1BroadException(t *testing.T) {
	t.Parallel()
	vuln := "def f():\n    try:\n        x()\n    except Exception:\n        log()\n"
	safe := "def f():\n    try:\n        x()\n    except Exception:\n        raise\n"
	assertRule(t, "BP-PY-1", "broad.py", vuln, true)
	assertRule(t, "BP-PY-1", "broad.py", safe, false)
}

func TestBPPY2ExceptPass(t *testing.T) {
	t.Parallel()
	vuln := "def f():\n    try:\n        x()\n    except ValueError:\n        pass\n"
	safe := "def f():\n    try:\n        x()\n    except ValueError as e:\n        logger.exception('x')\n        raise\n"
	assertRule(t, "BP-PY-2", "pass.py", vuln, true)
	assertRule(t, "BP-PY-2", "pass.py", safe, false)
}

func TestBPPY4MutableDefault(t *testing.T) {
	t.Parallel()
	vuln := "def append_item(x, items=[]):\n    items.append(x)\n    return items\n"
	safe := "def append_item(x, items=None):\n    if items is None:\n        items = []\n    items.append(x)\n    return items\n"
	assertRule(t, "BP-PY-4", "mut.py", vuln, true)
	assertRule(t, "BP-PY-4", "mut.py", safe, false)
	// Severity high
	findings := runBP(t, nil, vuln, "mut.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-4" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-4 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY6AssertValidation(t *testing.T) {
	t.Parallel()
	vuln := "def check(user):\n    assert user.is_admin\n    return True\n"
	assertRule(t, "BP-PY-6", "lib.py", vuln, true)
	// Test file skip
	assertRule(t, "BP-PY-6", "test_auth.py", vuln, false)
	assertRule(t, "BP-PY-6", "auth_test.py", vuln, false)
	assertRule(t, "BP-PY-6", "tests/test_auth.py", vuln, false)
}

func TestBPPY7OpenWithoutWith(t *testing.T) {
	t.Parallel()
	vuln := "def read(path):\n    f = open(path)\n    return f.read()\n"
	safe := "def read(path):\n    with open(path) as f:\n        return f.read()\n"
	assertRule(t, "BP-PY-7", "io.py", vuln, true)
	assertRule(t, "BP-PY-7", "io.py", safe, false)
}

func TestBPPY8SubprocessShell(t *testing.T) {
	t.Parallel()
	vuln := "import subprocess\nsubprocess.run(cmd, shell=True)\n"
	safe := "import subprocess\nsubprocess.run(['ls', '-la'], shell=False)\n"
	assertRule(t, "BP-PY-8", "sh.py", vuln, true)
	assertRule(t, "BP-PY-8", "sh.py", safe, false)
	findings := runBP(t, nil, vuln, "sh.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-8" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-8 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY9OSSystem(t *testing.T) {
	t.Parallel()
	vuln := "import os\nos.system('ls ' + path)\n"
	safe := "import subprocess\nsubprocess.run(['ls', path])\n"
	assertRule(t, "BP-PY-9", "sys.py", vuln, true)
	assertRule(t, "BP-PY-9", "sys.py", safe, false)
	assertRule(t, "BP-PY-9", "pop.py", "import os\nos.popen('ls')\n", true)
}

func TestBPPY10Pickle(t *testing.T) {
	t.Parallel()
	vuln := "import pickle\ndata = pickle.loads(body)\n"
	safe := "import json\ndata = json.loads(body)\n"
	assertRule(t, "BP-PY-10", "pk.py", vuln, true)
	assertRule(t, "BP-PY-10", "pk.py", safe, false)
}

func TestBPPY11YamlLoad(t *testing.T) {
	t.Parallel()
	vuln := "import yaml\ncfg = yaml.load(stream)\n"
	safe := "import yaml\ncfg = yaml.safe_load(stream)\n"
	safe2 := "import yaml\ncfg = yaml.load(stream, Loader=yaml.SafeLoader)\n"
	assertRule(t, "BP-PY-11", "y.py", vuln, true)
	assertRule(t, "BP-PY-11", "y.py", safe, false)
	assertRule(t, "BP-PY-11", "y.py", safe2, false)
}

func TestBPPY12EvalExec(t *testing.T) {
	t.Parallel()
	vuln := "def run(user_code):\n    return eval(user_code)\n"
	safe := "def run():\n    return eval('1+1')\n"
	assertRule(t, "BP-PY-12", "ev.py", vuln, true)
	assertRule(t, "BP-PY-12", "ev.py", safe, false)
	assertRule(t, "BP-PY-12", "ex.py", "exec(payload)\n", true)
}

func TestBPPY13HardcodedSecret(t *testing.T) {
	t.Parallel()
	vuln := "api_key = 'prod_not_a_real_credential_abcdef'\n"
	safePlaceholder := "password = 'changeme'\n"
	safeEnv := "import os\npassword = os.environ['PASSWORD']\n"
	assertRule(t, "BP-PY-13", "sec.py", vuln, true)
	assertRule(t, "BP-PY-13", "sec.py", safePlaceholder, false)
	assertRule(t, "BP-PY-13", "sec.py", safeEnv, false)
}

func TestBPPY16FlaskDebug(t *testing.T) {
	t.Parallel()
	vuln := "from flask import Flask\napp = Flask(__name__)\nif __name__ == '__main__':\n    app.run(debug=True)\n"
	safe := "from flask import Flask\napp = Flask(__name__)\nif __name__ == '__main__':\n    app.run(debug=False)\n"
	assertRule(t, "BP-PY-16", "app.py", vuln, true)
	assertRule(t, "BP-PY-16", "app.py", safe, false)
	// test file skip
	assertRule(t, "BP-PY-16", "test_app.py", vuln, false)
}

func TestBPPY17FlaskSecretKey(t *testing.T) {
	t.Parallel()
	vuln := "from flask import Flask\napp = Flask(__name__)\napp.secret_key = 'super-secret-key-value'\n"
	safe := "from flask import Flask\nimport os\napp = Flask(__name__)\napp.secret_key = os.environ['SECRET_KEY']\n"
	assertRule(t, "BP-PY-17", "app.py", vuln, true)
	assertRule(t, "BP-PY-17", "app.py", safe, false)
}

func TestBPPY21DjangoDebug(t *testing.T) {
	t.Parallel()
	vuln := "DEBUG = True\nINSTALLED_APPS = []\n"
	safe := "DEBUG = False\nINSTALLED_APPS = []\n"
	assertRule(t, "BP-PY-21", "settings.py", vuln, true)
	assertRule(t, "BP-PY-21", "settings.py", safe, false)
	// non-settings path without django markers should not fire
	assertRule(t, "BP-PY-21", "util.py", "DEBUG = True\n", false)
}

func TestNonPythonUnitSkipped(t *testing.T) {
	t.Parallel()
	d := badpractices.NewPythonBadPracticeScan()
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	unit := core.NewParsedUnit(core.LangGo, "x.go", "except:\n pass\n")
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	if len(out) != 0 {
		t.Fatalf("expected no findings for Go unit, got %#v", out)
	}
}

func TestEmptySourceSkipped(t *testing.T) {
	t.Parallel()
	d := badpractices.NewPythonBadPracticeScan()
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	unit := core.NewParsedUnit(core.LanguagePython, "x.py", "")
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	if len(out) != 0 {
		t.Fatalf("expected no findings for empty source, got %#v", out)
	}
}

func assertRule(t *testing.T, rule, path, src string, wantFire bool) {
	t.Helper()
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	ctx.Only = []string{rule}
	findings := runBP(t, ctx, src, path)
	got := hasRule(findings, rule)
	if got != wantFire {
		t.Fatalf("rule %s path %s want fire=%v got=%v findings=%v\nsrc:\n%s", rule, path, wantFire, got, findings, src)
	}
}

func runBP(t *testing.T, ctx *core.ScanContext, src, path string) []rules.Finding {
	t.Helper()
	d := badpractices.NewPythonBadPracticeScan()
	if ctx == nil {
		ctx = core.DefaultScanContext()
		ctx.BadPracticesEnabled = true
	}
	d.BeginScan(ctx)
	defer d.EndScan()
	unit := core.NewParsedUnit(core.LanguagePython, path, src)
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	return out
}

func hasRule(findings []rules.Finding, id string) bool {
	for _, f := range findings {
		if f.RuleID == id {
			return true
		}
	}
	return false
}

func contains(ids []string, want string) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
}

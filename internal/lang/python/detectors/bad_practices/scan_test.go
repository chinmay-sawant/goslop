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
		"BP-PY-1",
		"BP-PY-2",
		"BP-PY-3",
		"BP-PY-4",
		"BP-PY-5",
		"BP-PY-6",
		"BP-PY-7",
		"BP-PY-8",
		"BP-PY-9",
		"BP-PY-10",
		"BP-PY-11",
		"BP-PY-12",
		"BP-PY-13",
		"BP-PY-14",
		"BP-PY-15",
		"BP-PY-16",
		"BP-PY-17",
		"BP-PY-18",
		"BP-PY-19",
		"BP-PY-20",
		"BP-PY-21",
		"BP-PY-22",
		"BP-PY-23",
		"BP-PY-24",
		"BP-PY-25",
		"BP-PY-26",
		"BP-PY-27",
		"BP-PY-28",
		"BP-PY-29",
		"BP-PY-30",
		"BP-PY-31",
		"BP-PY-32",
		"BP-PY-33",
		"BP-PY-34",
		"BP-PY-35",
		"BP-PY-36",
		"BP-PY-37",
		"BP-PY-38",
		"BP-PY-39",
		"BP-PY-40",
		"BP-PY-41",
		"BP-PY-42",
		"BP-PY-43",
		"BP-PY-44",
		"BP-PY-45",
		"BP-PY-46",
		"BP-PY-47",
		"BP-PY-48",
		"BP-PY-49",
		"BP-PY-50",
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
		t.Fatalf("catalogue size too small: %d", badpractices.CatalogueSize())
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
	assertBPFixtureCase(t, "BP-PY-1", "BP-PY-1")
}

func TestBPPY1BroadException(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-1", "BP-PY-1-broad")
}

func TestBPPY2ExceptPass(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-2", "BP-PY-2")
}

func TestBPPY4MutableDefault(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-4", "BP-PY-4")
	vuln := loadBPFixture(t, "BP-PY-4", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-4" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-4 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY6AssertValidation(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-6", "BP-PY-6")
	assertBPFixtureCase(t, "BP-PY-6", "BP-PY-6-invariant")
	assertBPFixtureCase(t, "BP-PY-6", "BP-PY-6-test-path")
	// Extra test-path filename shapes reuse the vulnerable body.
	vuln := loadBPFixture(t, "BP-PY-6", true)
	assertRule(t, "BP-PY-6", "auth_test.py", vuln.body, false)
	assertRule(t, "BP-PY-6", "tests/test_auth.py", vuln.body, false)
}

func TestBPPY7OpenWithoutWith(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-7", "BP-PY-7")
}

func TestBPPY8SubprocessShell(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-8", "BP-PY-8")
	vuln := loadBPFixture(t, "BP-PY-8", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-8" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-8 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY9OSSystem(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-9", "BP-PY-9")
	assertBPFixtureCase(t, "BP-PY-9", "BP-PY-9-popen")
}

func TestBPPY10Pickle(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-10", "BP-PY-10")
	assertBPFixtureCase(t, "BP-PY-10", "BP-PY-10-local-cache")
	assertBPFixtureCase(t, "BP-PY-10", "BP-PY-10-cache-blob")
}

func TestBPPY11YamlLoad(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-11", "BP-PY-11")
	assertBPFixtureCase(t, "BP-PY-11", "BP-PY-11-safe-loader")
}

func TestBPPY12EvalExec(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-12", "BP-PY-12")
	assertBPFixtureCase(t, "BP-PY-12", "BP-PY-12-exec")
}

func TestBPPY13HardcodedSecret(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-13", "BP-PY-13")
	assertBPFixtureCase(t, "BP-PY-13", "BP-PY-13-placeholder")
}

func TestBPPY16FlaskDebug(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-16", "BP-PY-16")
	assertBPFixtureCase(t, "BP-PY-16", "BP-PY-16-test-path")
}

func TestBPPY17FlaskSecretKey(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-17", "BP-PY-17")
}

func TestBPPY21DjangoDebug(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-21", "BP-PY-21")
	assertBPFixtureCase(t, "BP-PY-21", "BP-PY-21-util")
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

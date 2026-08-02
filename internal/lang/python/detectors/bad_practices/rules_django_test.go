package badpractices_test

import (
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	badpractices "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPYDjangoRulesRegistered(t *testing.T) {
	t.Parallel()
	d := badpractices.NewPythonBadPracticeScan()
	want := []string{
		"BP-PY-21",
		"BP-PY-22", "BP-PY-23", "BP-PY-24", "BP-PY-25",
		"BP-PY-26", "BP-PY-27", "BP-PY-28",
	}
	ids := d.RuleIDs()
	for _, id := range want {
		if !contains(ids, id) {
			t.Errorf("missing registered rule %s (have %v)", id, ids)
		}
		meta := d.MetadataFor(id)
		if meta == nil {
			t.Errorf("MetadataFor(%s) = nil", id)
			continue
		}
		if meta.Pack != rules.PackBadPractice {
			t.Errorf("%s pack = %v", id, meta.Pack)
		}
		if !strings.HasPrefix(id, "BP-PY-") {
			t.Errorf("id %q must use BP-PY- prefix", id)
		}
	}
}

func TestBPPY22DjangoSecretKey(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-22", "BP-PY-22")
	assertBPFixtureCase(t, "BP-PY-22", "BP-PY-22-getenv")
	assertBPFixtureCase(t, "BP-PY-22", "BP-PY-22-flask-only")
	assertBPFixtureCase(t, "BP-PY-22", "BP-PY-22-util")

	vuln := loadBPFixture(t, "BP-PY-22", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-22" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-22 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY23AllowedHosts(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-23", "BP-PY-23")
	assertBPFixtureCase(t, "BP-PY-23", "BP-PY-23-empty-debug-false")
	assertBPFixtureCase(t, "BP-PY-23", "BP-PY-23-empty-debug-true")
}

func TestBPPY24RawSQLFormat(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-24", "BP-PY-24")
	assertBPFixtureCase(t, "BP-PY-24", "BP-PY-24-percent")
	assertBPFixtureCase(t, "BP-PY-24", "BP-PY-24-format")
	assertBPFixtureCase(t, "BP-PY-24", "BP-PY-24-execute-params")
}

func TestBPPY25MarkSafe(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-25", "BP-PY-25")
	assertBPFixtureCase(t, "BP-PY-25", "BP-PY-25-request")
}

func TestBPPY26CSRFExempt(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-26", "BP-PY-26")
}

func TestBPPY27MassAssignment(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-27", "BP-PY-27")
	assertBPFixtureCase(t, "BP-PY-27", "BP-PY-27-kwargs")
}

func TestBPPY28NPlusOne(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-28", "BP-PY-28")
	assertBPFixtureCase(t, "BP-PY-28", "BP-PY-28-single-hop")

	vuln := loadBPFixture(t, "BP-PY-28", true)
	findings := runBP(t, onlyCtx("BP-PY-28"), vuln.body, vuln.path)
	found := false
	for _, f := range findings {
		if f.RuleID == "BP-PY-28" {
			found = true
			if !strings.Contains(strings.ToLower(f.Message), "heuristic") &&
				!strings.Contains(strings.ToLower(f.Message), "review") {
				t.Fatalf("BP-PY-28 message should document confidence, got %q", f.Message)
			}
		}
	}
	if !found {
		t.Fatal("expected BP-PY-28 finding")
	}
}

func onlyCtx(rule string) *core.ScanContext {
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	ctx.Only = []string{rule}
	return ctx
}

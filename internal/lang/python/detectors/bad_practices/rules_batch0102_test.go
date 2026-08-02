package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	badpractices "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Batch 01+02 registration surface (BP-PY-3,5,14,15,18,19,20).
func TestBatch0102RulesRegistered(t *testing.T) {
	t.Parallel()
	d := badpractices.NewPythonBadPracticeScan()
	want := []string{
		"BP-PY-3", "BP-PY-5",
		"BP-PY-14", "BP-PY-15",
		"BP-PY-18", "BP-PY-19", "BP-PY-20",
	}
	ids := d.RuleIDs()
	for _, id := range want {
		if !contains(ids, id) {
			t.Errorf("missing registered rule %s", id)
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
}

func TestBPPY3RaiseGenericException(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-3", "BP-PY-3")
	assertBPFixtureCase(t, "BP-PY-3", "BP-PY-3-bare")
	assertBPFixtureCase(t, "BP-PY-3", "BP-PY-3-base-exception")
	assertBPFixtureCase(t, "BP-PY-3", "BP-PY-3-runtimeerror")
	assertBPFixtureCase(t, "BP-PY-3", "BP-PY-3-test-path")
	assertBPFixtureCase(t, "BP-PY-3", "BP-PY-3-tests-helpers")
}

func TestBPPY5WildcardImport(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-5", "BP-PY-5")
	assertBPFixtureCase(t, "BP-PY-5", "BP-PY-5-relative-star")
	assertBPFixtureCase(t, "BP-PY-5", "BP-PY-5-init-reexport")
}

func TestBPPY14RequestsWithoutTimeout(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-14", "BP-PY-14")
	assertBPFixtureCase(t, "BP-PY-14", "BP-PY-14-post")
	assertBPFixtureCase(t, "BP-PY-14", "BP-PY-14-timeout-tuple")
	assertBPFixtureCase(t, "BP-PY-14", "BP-PY-14-session")
	assertBPFixtureCase(t, "BP-PY-14", "BP-PY-14-test-path")
}

func TestBPPY15HttpxAsyncClientNotClosed(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-15", "BP-PY-15")
	assertBPFixtureCase(t, "BP-PY-15", "BP-PY-15-aclose")
}

func TestBPPY18FlaskRouteMissingMethods(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-18", "BP-PY-18")
	assertBPFixtureCase(t, "BP-PY-18", "BP-PY-18-args-only")
	assertBPFixtureCase(t, "BP-PY-18", "BP-PY-18-get-json")
}

func TestBPPY19FlaskJsonifyErrorLeak(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-19", "BP-PY-19")
	assertBPFixtureCase(t, "BP-PY-19", "BP-PY-19-traceback")
}

func TestBPPY20FlaskSendFileUserPath(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-20", "BP-PY-20")
	assertBPFixtureCase(t, "BP-PY-20", "BP-PY-20-from-directory")
	assertBPFixtureCase(t, "BP-PY-20", "BP-PY-20-user-root")
	assertBPFixtureCase(t, "BP-PY-20", "BP-PY-20-literal")

	vuln := loadBPFixture(t, "BP-PY-20", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-20" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-20 severity = %v, want high", f.Severity)
		}
	}
}

func TestBatch0102MetadataSeverities(t *testing.T) {
	t.Parallel()
	// Spot-check catalogue-aligned severities for new rules.
	cases := map[string]rules.Severity{
		"BP-PY-3":  rules.SeverityLow,
		"BP-PY-5":  rules.SeverityLow,
		"BP-PY-14": rules.SeverityMedium,
		"BP-PY-15": rules.SeverityMedium,
		"BP-PY-18": rules.SeverityLow,
		"BP-PY-19": rules.SeverityMedium,
		"BP-PY-20": rules.SeverityHigh,
	}
	for id, want := range cases {
		meta := badpractices.MetadataForID(id)
		if meta == nil {
			t.Fatalf("MetadataForID(%s) = nil", id)
		}
		if meta.Severity != want {
			t.Errorf("%s severity = %v, want %v", id, meta.Severity, want)
		}
	}
}

// Ensure scan helpers used by this file compile against package API.
var _ = core.LanguagePython

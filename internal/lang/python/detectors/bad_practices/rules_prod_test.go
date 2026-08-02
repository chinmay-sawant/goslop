package badpractices_test

import (
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	badpractices "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY48CORSStarWithCredentials(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-48", "BP-PY-48")
	assertBPFixtureCase(t, "BP-PY-48", "BP-PY-48-no-creds")
	assertBPFixtureCase(t, "BP-PY-48", "BP-PY-48-django")
	assertBPFixtureCase(t, "BP-PY-48", "BP-PY-48-flask")
	assertBPFixtureCase(t, "BP-PY-48", "BP-PY-48-test-path")

	vuln := loadBPFixture(t, "BP-PY-48", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-48" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-48 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY49TLSVerificationDisabled(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-49", "BP-PY-49")
	assertBPFixtureCase(t, "BP-PY-49", "BP-PY-49-httpx")
	assertBPFixtureCase(t, "BP-PY-49", "BP-PY-49-unverified-context")
	assertBPFixtureCase(t, "BP-PY-49", "BP-PY-49-cert-none")
	assertBPFixtureCase(t, "BP-PY-49", "BP-PY-49-default")
	assertBPFixtureCase(t, "BP-PY-49", "BP-PY-49-ca-path")
	assertBPFixtureCase(t, "BP-PY-49", "BP-PY-49-test-path")

	// Explicit: this detector must not fire BP-PY-14 (not implemented in this batch).
	missDefault := loadBPFixture(t, "BP-PY-49-default", false)
	findings := runBP(t, nil, missDefault.body, missDefault.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-14" {
			t.Fatalf("unexpected BP-PY-14 from batch-08 detector path; findings=%v", findings)
		}
		if f.RuleID == "BP-PY-49" {
			t.Fatalf("BP-PY-49 must not fire on requests without verify=False; findings=%v", findings)
		}
	}

	hit := loadBPFixture(t, "BP-PY-49", true)
	hitFindings := runBP(t, nil, hit.body, hit.path)
	for _, f := range hitFindings {
		if f.RuleID == "BP-PY-49" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-49 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY50InsecureCookieFlags(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-50", "BP-PY-50")
	assertBPFixtureCase(t, "BP-PY-50", "BP-PY-50-csrf")
	assertBPFixtureCase(t, "BP-PY-50", "BP-PY-50-httponly")
	assertBPFixtureCase(t, "BP-PY-50", "BP-PY-50-flask-config")
	assertBPFixtureCase(t, "BP-PY-50", "BP-PY-50-test-path")

	vuln := loadBPFixture(t, "BP-PY-50", true)
	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-50" && f.Severity != rules.SeverityMedium {
			t.Fatalf("BP-PY-50 severity = %v, want medium", f.Severity)
		}
	}
}

func TestBPProdBatch08Registered(t *testing.T) {
	t.Parallel()
	d := badpractices.NewPythonBadPracticeScan()
	ids := d.RuleIDs()
	for _, id := range []string{"BP-PY-48", "BP-PY-49", "BP-PY-50"} {
		if !contains(ids, id) {
			t.Errorf("missing registered rule %s in %v", id, ids)
		}
		meta := d.MetadataFor(id)
		if meta == nil {
			t.Errorf("MetadataFor(%s) = nil", id)
			continue
		}
		if meta.Pack != rules.PackBadPractice {
			t.Errorf("%s pack = %v, want PackBadPractice", id, meta.Pack)
		}
		if !strings.HasPrefix(id, "BP-PY-") {
			t.Errorf("id %q must use BP-PY- prefix", id)
		}
	}
	// BP-PY-14 is owned by batch-01; batch-08 only guarantees 48–50 are present.
	// If batch-01 has also landed, RuleIDs may include 14 — that is fine.
	// Catalogue metadata validation when JSON reachable.
	if err := badpractices.ValidateImplementedMetadata(); err != nil {
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "cannot find") {
			t.Logf("catalogue not reachable: %v", err)
		} else {
			t.Fatal(err)
		}
	}
	// Language / only filter smoke.
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	ctx.Only = []string{"BP-PY-48", "BP-PY-49", "BP-PY-50"}
	vuln := loadBPFixture(t, "BP-PY-49", true)
	findings := runBP(t, ctx, vuln.body, vuln.path)
	if !hasRule(findings, "BP-PY-49") {
		t.Fatalf("Only filter should allow BP-PY-49; findings=%v", findings)
	}
}

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
	// FastAPI / Starlette middleware hit
	hitFastAPI := `
from fastapi.middleware.cors import CORSMiddleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
)
`
	// Explicit origins + credentials — miss
	missExplicit := `
from fastapi.middleware.cors import CORSMiddleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["https://app.example"],
    allow_credentials=True,
)
`
	// Star without credentials true — miss
	missNoCreds := `
from fastapi.middleware.cors import CORSMiddleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=False,
)
`
	// Django cors-headers pair
	hitDjango := `
CORS_ALLOW_ALL_ORIGINS = True
CORS_ALLOW_CREDENTIALS = True
`
	missDjangoOnlyStar := `
CORS_ALLOW_ALL_ORIGINS = True
CORS_ALLOW_CREDENTIALS = False
`
	// flask-cors
	hitFlask := `
from flask_cors import CORS
CORS(app, origins="*", supports_credentials=True)
`
	assertRule(t, "BP-PY-48", "main.py", hitFastAPI, true)
	assertRule(t, "BP-PY-48", "main.py", missExplicit, false)
	assertRule(t, "BP-PY-48", "main.py", missNoCreds, false)
	assertRule(t, "BP-PY-48", "settings.py", hitDjango, true)
	assertRule(t, "BP-PY-48", "settings.py", missDjangoOnlyStar, false)
	assertRule(t, "BP-PY-48", "app.py", hitFlask, true)
	// test file skip
	assertRule(t, "BP-PY-48", "test_cors.py", hitFastAPI, false)

	findings := runBP(t, nil, hitFastAPI, "main.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-48" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-48 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY49TLSVerificationDisabled(t *testing.T) {
	t.Parallel()
	hitRequests := "import requests\nrequests.get(url, verify=False)\n"
	hitHttpx := "import httpx\nhttpx.get(url, verify=False)\n"
	hitUnverified := "import ssl\nctx = ssl._create_unverified_context()\n"
	hitCertNone := "import ssl\nctx = ssl.create_default_context()\nctx.verify_mode = ssl.CERT_NONE\n"
	missDefault := "import requests\nrequests.get(url, timeout=5)\n"
	missVerifyTrue := "import requests\nrequests.get(url, verify=True)\n"
	missCAPath := "import requests\nrequests.get(url, verify='/etc/ssl/certs/ca.pem')\n"

	assertRule(t, "BP-PY-49", "client.py", hitRequests, true)
	assertRule(t, "BP-PY-49", "client.py", hitHttpx, true)
	assertRule(t, "BP-PY-49", "ssl_util.py", hitUnverified, true)
	assertRule(t, "BP-PY-49", "ssl_util.py", hitCertNone, true)
	assertRule(t, "BP-PY-49", "client.py", missDefault, false)
	assertRule(t, "BP-PY-49", "client.py", missVerifyTrue, false)
	assertRule(t, "BP-PY-49", "client.py", missCAPath, false)
	// test file skip
	assertRule(t, "BP-PY-49", "test_client.py", hitRequests, false)

	// Explicit: this detector must not fire BP-PY-14 (not implemented in this batch).
	findings := runBP(t, nil, missDefault, "client.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-14" {
			t.Fatalf("unexpected BP-PY-14 from batch-08 detector path; findings=%v", findings)
		}
		if f.RuleID == "BP-PY-49" {
			t.Fatalf("BP-PY-49 must not fire on requests without verify=False; findings=%v", findings)
		}
	}

	hitFindings := runBP(t, nil, hitRequests, "client.py")
	for _, f := range hitFindings {
		if f.RuleID == "BP-PY-49" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-49 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY50InsecureCookieFlags(t *testing.T) {
	t.Parallel()
	hitSession := "SESSION_COOKIE_SECURE = False\n"
	hitCSRF := "CSRF_COOKIE_SECURE = False\n"
	hitHttpOnly := "SESSION_COOKIE_HTTPONLY = False\n"
	hitFlaskConfig := "app.config['SESSION_COOKIE_SECURE'] = False\n"
	missSecureTrue := "SESSION_COOKIE_SECURE = True\n"
	missCSRFTrue := "CSRF_COOKIE_SECURE = True\n"
	missHttpOnlyTrue := "SESSION_COOKIE_HTTPONLY = True\n"

	assertRule(t, "BP-PY-50", "settings.py", hitSession, true)
	assertRule(t, "BP-PY-50", "settings.py", hitCSRF, true)
	assertRule(t, "BP-PY-50", "settings.py", hitHttpOnly, true)
	assertRule(t, "BP-PY-50", "app.py", hitFlaskConfig, true)
	assertRule(t, "BP-PY-50", "settings.py", missSecureTrue, false)
	assertRule(t, "BP-PY-50", "settings.py", missCSRFTrue, false)
	assertRule(t, "BP-PY-50", "settings.py", missHttpOnlyTrue, false)
	// test file skip
	assertRule(t, "BP-PY-50", "test_settings.py", hitSession, false)

	findings := runBP(t, nil, hitSession, "settings.py")
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
	src := "import requests\nrequests.get(url, verify=False)\n"
	findings := runBP(t, ctx, src, "client.py")
	if !hasRule(findings, "BP-PY-49") {
		t.Fatalf("Only filter should allow BP-PY-49; findings=%v", findings)
	}
}

package badpractices_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY41TestWithoutAssert(t *testing.T) {
	t.Parallel()
	vuln := loadBPFixture(t, "BP-PY-41", true)
	safe := loadBPFixture(t, "BP-PY-41", false)
	assertRule(t, "BP-PY-41", vuln.path, vuln.body, true)
	assertRule(t, "BP-PY-41", safe.path, safe.body, false)

	// unittest style assert counts
	unittest := loadBPFixture(t, "BP-PY-41-unittest-assert", false)
	assertRule(t, "BP-PY-41", unittest.path, unittest.body, false)

	// _assert_* helper call counts as assertion (audit FPs 54–74, 88–89)
	helper := loadBPFixture(t, "BP-PY-41-assert-helper", false)
	assertRule(t, "BP-PY-41", helper.path, helper.body, false)

	// same-file helper that raises AssertionError (audit FP 75 pattern)
	ae := loadBPFixture(t, "BP-PY-41-assertionerror-helper", false)
	assertRule(t, "BP-PY-41", ae.path, ae.body, false)

	// pytest-regressions check_func (pictex)
	checkFunc := loadBPFixture(t, "BP-PY-41-check-func", false)
	assertRule(t, "BP-PY-41", checkFunc.path, checkFunc.body, false)

	// pytest-benchmark fixture (httptap)
	bench := loadBPFixture(t, "BP-PY-41-benchmark", false)
	assertRule(t, "BP-PY-41", bench.path, bench.body, false)

	// bare-assert helper delegation (safer / niquests _inner_*)
	bare := loadBPFixture(t, "BP-PY-41-bare-assert-helper", false)
	assertRule(t, "BP-PY-41", bare.path, bare.body, false)

	// triple-quoted sample string must not abort body scan (whatsapp-wrapped / rendercv)
	strBody := loadBPFixture(t, "BP-PY-41-string-body-scan", false)
	assertRule(t, "BP-PY-41", strBody.path, strBody.body, false)
	strVuln := loadBPFixture(t, "BP-PY-41-string-body-scan", true)
	assertRule(t, "BP-PY-41", strVuln.path, strVuln.body, true)

	// explicit raise AssertionError verifies the rejection outcome (httptap)
	raiseAssert := loadBPFixture(t, "BP-PY-41-raise-assertion", false)
	assertRule(t, "BP-PY-41", raiseAssert.path, raiseAssert.body, false)
	raiseAssertVuln := loadBPFixture(t, "BP-PY-41-raise-assertion", true)
	assertRule(t, "BP-PY-41", raiseAssertVuln.path, raiseAssertVuln.body, true)

	// top-level pytest.fail after a retry loop is a bare fallback, not an
	// assertion (httpmorph test_proxy TPs must keep firing)
	retryVuln := loadBPFixture(t, "BP-PY-41-retry-fallback", true)
	assertRule(t, "BP-PY-41", retryVuln.path, retryVuln.body, true)
	// pytest.fail inside an except suite verifies the guarded call (logxide)
	retrySafe := loadBPFixture(t, "BP-PY-41-retry-fallback", false)
	assertRule(t, "BP-PY-41", retrySafe.path, retrySafe.body, false)

	// noxfile session runners and pytest.raises-bearing helpers are not placeholders
	nox := loadBPFixture(t, "BP-PY-41-noxfile", false)
	assertRule(t, "BP-PY-41", nox.path, nox.body, false)
	noxVuln := loadBPFixture(t, "BP-PY-41-noxfile", true)
	assertRule(t, "BP-PY-41", noxVuln.path, noxVuln.body, true)
	raisesHelper := loadBPFixture(t, "BP-PY-41-pytest-raises-helper", false)
	assertRule(t, "BP-PY-41", raisesHelper.path, raisesHelper.body, false)
	raisesHelperVuln := loadBPFixture(t, "BP-PY-41-pytest-raises-helper", true)
	assertRule(t, "BP-PY-41", raisesHelperVuln.path, raisesHelperVuln.body, true)

	// niquests unicode_/redirect pure-call smokes; stream close; server wait;
	// rendercv never_crashes / silently_ignores names
	for _, caseName := range []string{
		"BP-PY-41-http-encoding-smoke",
		"BP-PY-41-stream-close",
		"BP-PY-41-server-wait-smoke",
		"BP-PY-41-never-crashes",
		// Project_Parva residual FPs: route handlers, check=True, validate_*
		"BP-PY-41-route-handler",
		"BP-PY-41-subprocess-check",
		"BP-PY-41-validate-helper",
	} {
		safeCase := loadBPFixture(t, caseName, false)
		vulnCase := loadBPFixture(t, caseName, true)
		assertRule(t, "BP-PY-41", safeCase.path, safeCase.body, false)
		assertRule(t, "BP-PY-41", vulnCase.path, vulnCase.body, true)
	}

	findings := runBP(t, nil, vuln.body, vuln.path)
	for _, f := range findings {
		if f.RuleID == "BP-PY-41" && f.Severity != rules.SeverityInfo {
			t.Fatalf("BP-PY-41 severity = %v, want info", f.Severity)
		}
	}
}

func TestBPPY41SiblingHelpersPy(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	helpers := loadBPFixtureFile(t, "BP-PY-41-sibling-helpers-helpers.txt")
	if err := os.WriteFile(filepath.Join(dir, "helpers.py"), []byte(helpers.body), 0o644); err != nil {
		t.Fatal(err)
	}
	call := loadBPFixtureFile(t, "BP-PY-41-sibling-helpers-call.txt")
	assertRule(t, "BP-PY-41", filepath.Join(dir, filepath.Base(call.path)), call.body, false)

	// Production-only call in the same layout must still fire.
	vuln := loadBPFixture(t, "BP-PY-41-sibling-helpers", true)
	assertRule(t, "BP-PY-41", filepath.Join(dir, filepath.Base(vuln.path)), vuln.body, true)
}

func TestBPPY42TryExceptInsteadOfRaises(t *testing.T) {
	t.Parallel()
	vuln := loadBPFixture(t, "BP-PY-42-assertionerror", true)
	// Broad except Exception also hits
	vuln2 := loadBPFixture(t, "BP-PY-42", true)
	safe := loadBPFixture(t, "BP-PY-42", false)
	safe2 := loadBPFixture(t, "BP-PY-42-assertraises", false)
	assertRule(t, "BP-PY-42", vuln.path, vuln.body, true)
	assertRule(t, "BP-PY-42", vuln2.path, vuln2.body, true)
	assertRule(t, "BP-PY-42", safe.path, safe.body, false)
	assertRule(t, "BP-PY-42", safe2.path, safe2.body, false)
}

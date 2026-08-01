package badpractices_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY41TestWithoutAssert(t *testing.T) {
	t.Parallel()
	vuln := "def test_api():\n    client.get('/x')\n"
	safe := "def test_api():\n    r = client.get('/x')\n    assert r.status_code == 200\n"
	assertRule(t, "BP-PY-41", "test_api.py", vuln, true)
	assertRule(t, "BP-PY-41", "test_api.py", safe, false)
	// unittest style assert counts
	assertRule(t, "BP-PY-41", "test_api.py",
		"def test_api():\n    r = client.get('/x')\n    self.assertEqual(r.status_code, 200)\n", false)
	// _assert_* helper call counts as assertion (audit FPs 54–74, 88–89)
	assertRule(t, "BP-PY-41", "test_compliance.py",
		"class T:\n    def _assert_verapdf(self, name):\n        self.assertEqual(0, 0)\n    def test_retail_passes_pdfa4(self):\n        self._assert_verapdf('x')\n", false)
	// same-file helper that raises AssertionError (audit FP 75 pattern)
	assertRule(t, "BP-PY-41", "test_doc.py",
		"def find_object_with(data, marker, offsets):\n    raise AssertionError('missing')\n"+
			"def test_catalog_and_pages_present():\n    find_object_with(b'', b'/Type', {})\n", false)
	findings := runBP(t, nil, vuln, "test_api.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-41" && f.Severity != rules.SeverityInfo {
			t.Fatalf("BP-PY-41 severity = %v, want info", f.Severity)
		}
	}
}

func TestBPPY41SiblingHelpersPy(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	helpers := "def find_object_with(data, marker, offsets):\n    raise AssertionError('missing')\n"
	if err := os.WriteFile(filepath.Join(dir, "helpers.py"), []byte(helpers), 0o644); err != nil {
		t.Fatal(err)
	}
	testSrc := "def test_catalog_and_pages_present():\n    find_object_with(b'', b'/Type', {})\n"
	assertRule(t, "BP-PY-41", filepath.Join(dir, "test_doc.py"), testSrc, false)
	// Production-only call in the same layout must still fire.
	assertRule(t, "BP-PY-41", filepath.Join(dir, "test_api.py"),
		"def test_api():\n    client.get('/x')\n", true)
}

func TestBPPY42TryExceptInsteadOfRaises(t *testing.T) {
	t.Parallel()
	vuln := "def test_foo():\n    try:\n        boom()\n    except AssertionError:\n        pass\n"
	// Broad except Exception also hits
	vuln2 := "def test_foo():\n    try:\n        boom()\n    except Exception:\n        pass\n"
	safe := "import pytest\ndef test_foo():\n    with pytest.raises(ValueError):\n        boom()\n"
	safe2 := "def test_foo():\n    with self.assertRaises(ValueError):\n        boom()\n"
	assertRule(t, "BP-PY-42", "test_foo.py", vuln, true)
	assertRule(t, "BP-PY-42", "test_foo.py", vuln2, true)
	assertRule(t, "BP-PY-42", "test_foo.py", safe, false)
	assertRule(t, "BP-PY-42", "test_foo.py", safe2, false)
}

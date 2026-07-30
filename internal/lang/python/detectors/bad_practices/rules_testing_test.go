package badpractices_test

import (
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
	findings := runBP(t, nil, vuln, "test_api.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-41" && f.Severity != rules.SeverityInfo {
			t.Fatalf("BP-PY-41 severity = %v, want info", f.Severity)
		}
	}
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

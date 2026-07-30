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
	assertRule(t, "BP-PY-3", "app.py", "def f():\n    raise Exception('bad')\n", true)
	assertRule(t, "BP-PY-3", "app.py", "def f():\n    raise Exception\n", true)
	assertRule(t, "BP-PY-3", "app.py", "def f():\n    raise BaseException('x')\n", true)
	assertRule(t, "BP-PY-3", "app.py", "def f():\n    raise ValueError('bad')\n", false)
	assertRule(t, "BP-PY-3", "app.py", "def f():\n    raise RuntimeError('bad')\n", false)
	// test file skip
	assertRule(t, "BP-PY-3", "test_app.py", "def f():\n    raise Exception('bad')\n", false)
	assertRule(t, "BP-PY-3", "tests/helpers.py", "def f():\n    raise Exception('bad')\n", false)
}

func TestBPPY5WildcardImport(t *testing.T) {
	t.Parallel()
	assertRule(t, "BP-PY-5", "app.py", "from os.path import *\n", true)
	assertRule(t, "BP-PY-5", "app.py", "from .models import *\n", true)
	assertRule(t, "BP-PY-5", "app.py", "from os.path import join, exists\n", false)
	// Policy: allow re-export in package __init__.py
	assertRule(t, "BP-PY-5", "pkg/__init__.py", "from .models import *\n", false)
}

func TestBPPY14RequestsWithoutTimeout(t *testing.T) {
	t.Parallel()
	assertRule(t, "BP-PY-14", "client.py", "import requests\nrequests.get('https://example.com')\n", true)
	assertRule(t, "BP-PY-14", "client.py", "import requests\nrequests.post(url, json=body)\n", true)
	assertRule(t, "BP-PY-14", "client.py", "import requests\nrequests.get(url, timeout=5)\n", false)
	assertRule(t, "BP-PY-14", "client.py", "import requests\nrequests.get(url, timeout=(1, 3))\n", false)
	// session helper name
	assertRule(t, "BP-PY-14", "client.py", "session.get(url)\n", true)
	assertRule(t, "BP-PY-14", "client.py", "session.get(url, timeout=2)\n", false)
	// test file skip
	assertRule(t, "BP-PY-14", "test_client.py", "import requests\nrequests.get(url)\n", false)
}

func TestBPPY15HttpxAsyncClientNotClosed(t *testing.T) {
	t.Parallel()
	hit := "import httpx\nclient = httpx.AsyncClient()\n"
	assertRule(t, "BP-PY-15", "http.py", hit, true)

	missWith := "import httpx\nasync def f(url):\n    async with httpx.AsyncClient() as client:\n        await client.get(url)\n"
	assertRule(t, "BP-PY-15", "http.py", missWith, false)

	missAclose := "import httpx\nasync def f():\n    client = httpx.AsyncClient()\n    try:\n        await client.get(url)\n    finally:\n        await client.aclose()\n"
	assertRule(t, "BP-PY-15", "http.py", missAclose, false)
}

func TestBPPY18FlaskRouteMissingMethods(t *testing.T) {
	t.Parallel()
	hit := "from flask import Flask, request\napp = Flask(__name__)\n\n@app.route('/login')\ndef login():\n    user = request.form['user']\n    return user\n"
	assertRule(t, "BP-PY-18", "app.py", hit, true)

	missMethods := "from flask import Flask, request\napp = Flask(__name__)\n\n@app.route('/login', methods=['POST'])\ndef login():\n    user = request.form['user']\n    return user\n"
	assertRule(t, "BP-PY-18", "app.py", missMethods, false)

	missArgsOnly := "from flask import Flask, request\napp = Flask(__name__)\n\n@app.route('/search')\ndef search():\n    q = request.args.get('q')\n    return q\n"
	assertRule(t, "BP-PY-18", "app.py", missArgsOnly, false)

	// get_json without methods
	hitJSON := "from flask import Flask, request\napp = Flask(__name__)\n\n@app.route('/api')\ndef api():\n    data = request.get_json()\n    return data\n"
	assertRule(t, "BP-PY-18", "app.py", hitJSON, true)
}

func TestBPPY19FlaskJsonifyErrorLeak(t *testing.T) {
	t.Parallel()
	hit := "from flask import Flask, jsonify\napp = Flask(__name__)\n\n@app.errorhandler(Exception)\ndef handle(e):\n    return jsonify(error=str(e)), 500\n"
	assertRule(t, "BP-PY-19", "app.py", hit, true)

	hitTB := "from flask import Flask\nimport traceback\napp = Flask(__name__)\n\n@app.errorhandler(Exception)\ndef handle(e):\n    return traceback.format_exc()\n"
	assertRule(t, "BP-PY-19", "app.py", hitTB, true)

	miss := "from flask import Flask, jsonify\napp = Flask(__name__)\n\n@app.errorhandler(Exception)\ndef handle(e):\n    return jsonify(error='internal'), 500\n"
	assertRule(t, "BP-PY-19", "app.py", miss, false)
}

func TestBPPY20FlaskSendFileUserPath(t *testing.T) {
	t.Parallel()
	hit := "from flask import send_file, request\n\n@app.route('/dl')\ndef dl():\n    return send_file(request.args['path'])\n"
	assertRule(t, "BP-PY-20", "app.py", hit, true)

	hitDir := "from flask import send_from_directory, request\n\ndef dl():\n    return send_from_directory('/var/data', request.args.get('f'))\n"
	assertRule(t, "BP-PY-20", "app.py", hitDir, true)

	missLit := "from flask import send_file\n\ndef dl():\n    return send_file('/var/data/report.pdf')\n"
	assertRule(t, "BP-PY-20", "app.py", missLit, false)

	missSafe := "from flask import send_from_directory, request, safe_join\nROOT = '/var/data'\n\ndef dl():\n    return send_from_directory(ROOT, safe_join(ROOT, request.args.get('f')))\n"
	assertRule(t, "BP-PY-20", "app.py", missSafe, false)

	// severity high
	findings := runBP(t, nil, hit, "app.py")
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
			continue
		}
		if meta.Severity != want {
			t.Errorf("%s severity = %v, want %v", id, meta.Severity, want)
		}
	}
}

// Ensure scan helpers used by this file compile against package API.
var _ = core.LanguagePython

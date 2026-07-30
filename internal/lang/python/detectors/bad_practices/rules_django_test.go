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
	vuln := "SECRET_KEY = 'django-insecure-abc123-not-for-prod'\nINSTALLED_APPS = []\n"
	safeEnv := "import os\nSECRET_KEY = os.environ['SECRET_KEY']\nINSTALLED_APPS = []\n"
	safeGetenv := "import os\nSECRET_KEY = os.getenv('SECRET_KEY')\nINSTALLED_APPS = []\n"
	assertRule(t, "BP-PY-22", "settings.py", vuln, true)
	assertRule(t, "BP-PY-22", "settings.py", safeEnv, false)
	assertRule(t, "BP-PY-22", "settings.py", safeGetenv, false)
	// Flask-only module should not steal Flask cases (BP-PY-17 owns those).
	flaskOnly := "from flask import Flask\napp = Flask(__name__)\napp.secret_key = 'super-secret-key-value'\n"
	assertRule(t, "BP-PY-22", "app.py", flaskOnly, false)
	// Non-settings without django markers.
	assertRule(t, "BP-PY-22", "util.py", "SECRET_KEY = 'abc1234567890'\n", false)
	// Severity high
	findings := runBP(t, nil, vuln, "settings.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-22" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-22 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY23AllowedHosts(t *testing.T) {
	t.Parallel()
	star := "DEBUG = False\nALLOWED_HOSTS = ['*']\nINSTALLED_APPS = []\n"
	emptyDebugFalse := "DEBUG = False\nALLOWED_HOSTS = []\nINSTALLED_APPS = []\n"
	ok := "DEBUG = False\nALLOWED_HOSTS = ['example.com']\nINSTALLED_APPS = []\n"
	// Empty with DEBUG True should not fire (catalogue: empty when DEBUG False).
	emptyDebugTrue := "DEBUG = True\nALLOWED_HOSTS = []\nINSTALLED_APPS = []\n"
	assertRule(t, "BP-PY-23", "settings.py", star, true)
	assertRule(t, "BP-PY-23", "settings.py", emptyDebugFalse, true)
	assertRule(t, "BP-PY-23", "settings.py", ok, false)
	assertRule(t, "BP-PY-23", "settings.py", emptyDebugTrue, false)
}

func TestBPPY24RawSQLFormat(t *testing.T) {
	t.Parallel()
	fstr := "from django.contrib.auth.models import User\nUser.objects.raw(f\"SELECT * FROM auth_user WHERE id = {uid}\")\n"
	pct := "cursor.execute(\"SELECT * FROM t WHERE id = %s\" % (uid,))\n"
	safe := "from django.contrib.auth.models import User\nUser.objects.raw(\"SELECT * FROM auth_user WHERE id = %s\", [uid])\n"
	safeExec := "cursor.execute(\"SELECT * FROM t WHERE id = %s\", [uid])\n"
	fmt := "User.objects.raw(\"SELECT * FROM t WHERE id = {}\".format(uid))\n"
	assertRule(t, "BP-PY-24", "views.py", fstr, true)
	assertRule(t, "BP-PY-24", "views.py", pct, true)
	assertRule(t, "BP-PY-24", "views.py", fmt, true)
	assertRule(t, "BP-PY-24", "views.py", safe, false)
	assertRule(t, "BP-PY-24", "views.py", safeExec, false)
}

func TestBPPY25MarkSafe(t *testing.T) {
	t.Parallel()
	vuln := "from django.utils.safestring import mark_safe\nhtml = mark_safe(user_input)\n"
	vulnReq := "from django.utils.safestring import mark_safe\nhtml = mark_safe(request.GET['q'])\n"
	safe := "from django.utils.safestring import mark_safe\nhtml = mark_safe(\"<br>\")\n"
	assertRule(t, "BP-PY-25", "views.py", vuln, true)
	assertRule(t, "BP-PY-25", "views.py", vulnReq, true)
	assertRule(t, "BP-PY-25", "views.py", safe, false)
}

func TestBPPY26CSRFExempt(t *testing.T) {
	t.Parallel()
	vuln := "from django.views.decorators.csrf import csrf_exempt\n\n@csrf_exempt\ndef pay(request):\n    amount = request.POST['amount']\n    return amount\n"
	safe := "def pay(request):\n    amount = request.POST['amount']\n    return amount\n"
	assertRule(t, "BP-PY-26", "views.py", vuln, true)
	assertRule(t, "BP-PY-26", "views.py", safe, false)
}

func TestBPPY27MassAssignment(t *testing.T) {
	t.Parallel()
	vuln := "def create(request):\n    u = User(**request.POST)\n    return u\n"
	vulnCreate := "def create(request):\n    return User.objects.create(**request.data)\n"
	safe := "def create(request):\n    return User.objects.create(username=request.POST['u'])\n"
	assertRule(t, "BP-PY-27", "views.py", vuln, true)
	assertRule(t, "BP-PY-27", "views.py", vulnCreate, true)
	assertRule(t, "BP-PY-27", "views.py", safe, false)
}

func TestBPPY28NPlusOne(t *testing.T) {
	t.Parallel()
	vuln := "def list_posts():\n    for item in Post.objects.all():\n        print(item.author.name)\n"
	safe := "def list_posts():\n    for item in Post.objects.select_related('author').all():\n        print(item.author.name)\n"
	// Single-hop access only: conservative heuristic should miss.
	singleHop := "def list_posts():\n    for item in Post.objects.all():\n        print(item.title)\n"
	assertRule(t, "BP-PY-28", "views.py", vuln, true)
	assertRule(t, "BP-PY-28", "views.py", safe, false)
	assertRule(t, "BP-PY-28", "views.py", singleHop, false)
	// Message should mention heuristic/review-only.
	findings := runBP(t, onlyCtx("BP-PY-28"), vuln, "views.py")
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

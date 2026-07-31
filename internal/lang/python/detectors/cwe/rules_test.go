package cwe_test

import (
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func runScan(src string) []rules.Finding {
	d := cwe.NewPyCweScan()
	unit := core.NewParsedUnit(core.LanguagePython, "sample.py", src)
	var out []rules.Finding
	d.Run(nil, unit, &out)
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

func TestCWE502HitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		hit  bool
	}{
		{
			name: "pickle.loads vulnerable",
			src:  "import pickle\n\ndef load(data):\n    return pickle.loads(data)\n",
			hit:  true,
		},
		{
			name: "pickle.load vulnerable",
			src:  "import pickle\n\ndef load(f):\n    return pickle.load(f)\n",
			hit:  true,
		},
		{
			name: "yaml.load without SafeLoader",
			src:  "import yaml\n\ndef load(data):\n    return yaml.load(data)\n",
			hit:  true,
		},
		{
			name: "yaml.unsafe_load",
			src:  "import yaml\n\ndef load(data):\n    return yaml.unsafe_load(data)\n",
			hit:  true,
		},
		{
			name: "yaml.safe_load ok",
			src:  "import yaml\n\ndef load(data):\n    return yaml.safe_load(data)\n",
			hit:  false,
		},
		{
			name: "yaml.load SafeLoader ok",
			src:  "import yaml\n\ndef load(data):\n    return yaml.load(data, Loader=yaml.SafeLoader)\n",
			hit:  false,
		},
		{
			name: "json.loads ok",
			src:  "import json\n\ndef load(data):\n    return json.loads(data)\n",
			hit:  false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			findings := runScan(tc.src)
			got := hasRule(findings, "CWE-502")
			if got != tc.hit {
				t.Fatalf("CWE-502 hit=%v want=%v findings=%v\nsrc:\n%s", got, tc.hit, summarize(findings), tc.src)
			}
			if tc.hit {
				for _, f := range findings {
					if f.RuleID == "CWE-502" && f.Severity != rules.SeverityHigh {
						t.Fatalf("severity = %v, want High", f.Severity)
					}
					if f.RuleID == "CWE-502" && f.Line < 1 {
						t.Fatalf("line = %d, want >= 1", f.Line)
					}
				}
			}
		})
	}
}

func TestCWE78HitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		hit  bool
	}{
		{
			name: "os.system dynamic",
			src:  "import os\n\ndef run(cmd):\n    os.system(cmd)\n",
			hit:  true,
		},
		{
			name: "os.system f-string",
			src:  "import os\n\ndef run(user):\n    os.system(f'ls {user}')\n",
			hit:  true,
		},
		{
			name: "subprocess shell=True dynamic",
			src:  "import subprocess\n\ndef run(cmd):\n    subprocess.run(cmd, shell=True)\n",
			hit:  true,
		},
		{
			name: "subprocess.call shell=True concat",
			src:  "import subprocess\n\ndef run(user):\n    subprocess.call('echo ' + user, shell=True)\n",
			hit:  true,
		},
		{
			name: "list argv shell=False safe",
			src:  "import subprocess\n\ndef run():\n    subprocess.run(['ls', '-l'], shell=False)\n",
			hit:  false,
		},
		{
			name: "list argv default shell safe",
			src:  "import subprocess\n\ndef run():\n    subprocess.run(['ls', '-l'])\n",
			hit:  false,
		},
		{
			name: "os.system pure literal safe",
			src:  "import os\n\ndef run():\n    os.system('ls -l')\n",
			hit:  false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			findings := runScan(tc.src)
			got := hasRule(findings, "CWE-78")
			if got != tc.hit {
				t.Fatalf("CWE-78 hit=%v want=%v findings=%v\nsrc:\n%s", got, tc.hit, summarize(findings), tc.src)
			}
		})
	}
}

func TestCWE89HitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		hit  bool
	}{
		{
			name: "f-string execute",
			src:  "def q(cur, user_id):\n    cur.execute(f\"SELECT * FROM users WHERE id = {user_id}\")\n",
			hit:  true,
		},
		{
			name: "percent format execute",
			src:  "def q(cur, name):\n    cur.execute(\"SELECT * FROM users WHERE name = '%s'\" % name)\n",
			hit:  true,
		},
		{
			name: "format method execute",
			src:  "def q(cur, name):\n    cur.execute(\"SELECT * FROM users WHERE name = '{}'\".format(name))\n",
			hit:  true,
		},
		{
			name: "parameterized ? safe",
			src:  "def q(cur, user_id):\n    cur.execute(\"SELECT * FROM users WHERE id = ?\", (user_id,))\n",
			hit:  false,
		},
		{
			name: "parameterized %s bound safe",
			src:  "def q(cur, user_id):\n    cur.execute(\"SELECT * FROM users WHERE id = %s\", (user_id,))\n",
			hit:  false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			findings := runScan(tc.src)
			got := hasRule(findings, "CWE-89")
			if got != tc.hit {
				t.Fatalf("CWE-89 hit=%v want=%v findings=%v\nsrc:\n%s", got, tc.hit, summarize(findings), tc.src)
			}
		})
	}
}

func TestCWE22HitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		hit  bool
	}{
		{
			name: "open join user path",
			src:  "import os\n\ndef read(root, user_path):\n    return open(os.path.join(root, user_path))\n",
			hit:  true,
		},
		{
			name: "Path div dynamic",
			src: "from pathlib import Path\n\ndef read(root, request_path):\n" +
				"    p = Path(root) / request_path\n    return open(p)\n",
			hit: true,
		},
		{
			name: "literal open safe",
			src:  "def read():\n    return open('/etc/hosts')\n",
			hit:  false,
		},
		{
			name: "basename confinement safe",
			src: "import os\n\ndef read(root, user_path):\n" +
				"    name = os.path.basename(user_path)\n    return open(os.path.join(root, name))\n",
			hit: false,
		},
		{
			name: "resolve startswith safe",
			src: "from pathlib import Path\n\ndef read(root, user_path):\n" +
				"    base = Path(root).resolve()\n    full = (base / user_path).resolve()\n" +
				"    if str(full).startswith(str(base)):\n        return open(full)\n    raise ValueError('escape')\n",
			hit: false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			findings := runScan(tc.src)
			got := hasRule(findings, "CWE-22")
			if got != tc.hit {
				t.Fatalf("CWE-22 hit=%v want=%v findings=%v\nsrc:\n%s", got, tc.hit, summarize(findings), tc.src)
			}
		})
	}
}

func TestCWE79HitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		hit  bool
	}{
		{
			name: "mark_safe dynamic",
			src:  "from django.utils.safestring import mark_safe\n\ndef render(user):\n    return mark_safe(user)\n",
			hit:  true,
		},
		{
			name: "Markup dynamic",
			src:  "from markupsafe import Markup\n\ndef render(req):\n    return Markup(req.args.get('q'))\n",
			hit:  true,
		},
		{
			name: "render_template_string f-string",
			src:  "from flask import render_template_string\n\ndef render(name):\n    return render_template_string(f'<h1>{name}</h1>')\n",
			hit:  true,
		},
		{
			name: "plain render_template safe",
			src:  "from flask import render_template\n\ndef render(name):\n    return render_template('x.html', name=name)\n",
			hit:  false,
		},
		{
			name: "mark_safe pure literal safe",
			src:  "from django.utils.safestring import mark_safe\n\ndef render():\n    return mark_safe('<b>ok</b>')\n",
			hit:  false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			findings := runScan(tc.src)
			got := hasRule(findings, "CWE-79")
			if got != tc.hit {
				t.Fatalf("CWE-79 hit=%v want=%v findings=%v\nsrc:\n%s", got, tc.hit, summarize(findings), tc.src)
			}
		})
	}
}

func TestBatchIntegrityNoDoubleIDs(t *testing.T) {
	t.Parallel()
	ids := cwe.RegisteredRuleIDs()
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate rule id %s", id)
		}
		seen[id] = true
	}
	if len(ids) != 159 {
		t.Fatalf("ids = %v (len %d), want 159", ids, len(ids))
	}
}

func summarize(findings []rules.Finding) string {
	if len(findings) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteString("[")
	for i, f := range findings {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(f.RuleID)
		b.WriteString("@")
		b.WriteString(strings.ReplaceAll(f.Message, " ", "_"))
	}
	b.WriteString("]")
	return b.String()
}

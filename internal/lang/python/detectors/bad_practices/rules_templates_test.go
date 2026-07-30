package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY33JinjaAutoescape(t *testing.T) {
	t.Parallel()
	vuln := `from jinja2 import Environment
env = Environment(autoescape=False)
`
	vulnMulti := `from jinja2 import Environment
env = Environment(
    loader=loader,
    autoescape=False,
)
`
	safeTrue := `from jinja2 import Environment
env = Environment(autoescape=True)
`
	safeSelect := `from jinja2 import Environment, select_autoescape
env = Environment(autoescape=select_autoescape(["html", "xml"]))
`
	assertRule(t, "BP-PY-33", "tpl.py", vuln, true)
	assertRule(t, "BP-PY-33", "tpl.py", vulnMulti, true)
	assertRule(t, "BP-PY-33", "tpl.py", safeTrue, false)
	assertRule(t, "BP-PY-33", "tpl.py", safeSelect, false)
	findings := runBP(t, nil, vuln, "tpl.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-33" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-33 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY34MarkupSafe(t *testing.T) {
	t.Parallel()
	vulnName := `from markupsafe import Markup
def render(user_html):
    return Markup(user_html)
`
	vulnReq := `from flask import request
from markupsafe import Markup
def render():
    return Markup(request.args["x"])
`
	vulnSafeFilter := `template = "{{ x|safe }}"
`
	safeLit := `from markupsafe import Markup
def render():
    return Markup("<br/>")
`
	assertRule(t, "BP-PY-34", "tpl.py", vulnName, true)
	assertRule(t, "BP-PY-34", "tpl.py", vulnReq, true)
	assertRule(t, "BP-PY-34", "tpl.py", vulnSafeFilter, true)
	assertRule(t, "BP-PY-34", "tpl.py", safeLit, false)
}

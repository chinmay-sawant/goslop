package ignore_test

import (
	"strings"
	"testing"

	"github.com/chinmay/goslop/internal/engine/ignore"
	"github.com/chinmay/goslop/internal/rules"
)

func TestParseInlineNextLine(t *testing.T) {
	src := `package main

// goslop-ignore: CWE-78
exec.Command("sh", "-c", cmd)
`
	m := ignore.ParseInlineIgnores(src)
	// line of exec.Command should be suppressed
	var found bool
	for ln, d := range m {
		if d.Matches("CWE-78") {
			found = true
			// Should target the code line after the comment
			if ln < 3 {
				t.Fatalf("unexpected line %d for ignore", ln)
			}
		}
	}
	if !found {
		t.Fatalf("no CWE-78 ignore: %+v", m)
	}
}

func TestApplySuppressesFinding(t *testing.T) {
	src := `package main

import "os/exec"

func handler() {
	// goslop-ignore: CWE-78
	exec.Command("sh", "-c", "x")
}
`
	// Find line of exec.Command
	lines := strings.Split(src, "\n")
	execLine := 0
	for i, ln := range lines {
		if strings.Contains(ln, "exec.Command") {
			execLine = i + 1
			break
		}
	}
	findings := []rules.Finding{
		{RuleID: "CWE-78", File: "x.go", Line: execLine, Column: 2, Message: "cmd injection", Severity: rules.SeverityHigh},
		{RuleID: "PERF-1", File: "x.go", Line: execLine, Column: 2, Message: "other", Severity: rules.SeverityLow},
	}
	out, n := ignore.Apply(src, findings, ignore.ApplyOptions{})
	if n != 1 {
		t.Fatalf("suppressed count: %d", n)
	}
	if len(out) != 1 || out[0].RuleID != "PERF-1" {
		t.Fatalf("out: %+v", out)
	}
}

func TestApplyShowIgnored(t *testing.T) {
	src := `package main
// goslop-ignore: CWE-78
exec.Command("x")
`
	findings := []rules.Finding{
		{RuleID: "CWE-78", File: "x.go", Line: 3, Column: 1, Message: "bad", Severity: rules.SeverityHigh},
	}
	out, n := ignore.Apply(src, findings, ignore.ApplyOptions{ShowIgnored: true})
	if n != 1 || len(out) != 1 {
		t.Fatalf("n=%d out=%+v", n, out)
	}
	if !out[0].Suppressed || out[0].Severity != rules.SeverityInfo {
		t.Fatalf("show-ignored: %+v", out[0])
	}
}

func TestFileIgnoreAll(t *testing.T) {
	src := `// goslop-ignore-file
package main
func f() {}
`
	d, ok := ignore.ParseFileIgnore(src)
	if !ok || !d.IsAll() {
		t.Fatalf("file ignore: ok=%v d=%+v", ok, d)
	}
	findings := []rules.Finding{{RuleID: "CWE-78", Line: 3, Message: "x", Severity: rules.SeverityHigh}}
	out, n := ignore.Apply(src, findings, ignore.ApplyOptions{})
	if n != 1 || len(out) != 0 {
		t.Fatalf("n=%d out=%+v", n, out)
	}
}

func TestIgnoreNotInString(t *testing.T) {
	src := `package main
const s = "// goslop-ignore: CWE-78"
exec.Command("x")
`
	m := ignore.ParseInlineIgnores(src)
	for _, d := range m {
		if d.Matches("CWE-78") {
			t.Fatal("should not parse ignore inside string")
		}
	}
}

func TestFixtureStyleSuppressedInline(t *testing.T) {
	// Mirrors tests/fixtures/go/baseline/suppressed_inline.txt
	src := `package main

import (
	"net/http"
	"os/exec"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	// goslop-ignore: CWE-78
	exec.Command("sh", "-c", cmd).Run()
}
`
	findings := []rules.Finding{
		{RuleID: "CWE-78", File: "suppressed_inline.go", Line: 12, Column: 2, Message: "command injection", Severity: rules.SeverityHigh},
	}
	// Detect actual line of exec.Command
	for i, ln := range strings.Split(src, "\n") {
		if strings.Contains(ln, "exec.Command") {
			findings[0].Line = i + 1
		}
	}
	out, n := ignore.Apply(src, findings, ignore.ApplyOptions{})
	if n != 1 || len(out) != 0 {
		t.Fatalf("fixture-style suppress: n=%d out=%+v line=%d", n, out, findings[0].Line)
	}
}

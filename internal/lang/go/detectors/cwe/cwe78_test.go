package cwe_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/go/detectors/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestCWE78Vulnerable(t *testing.T) {
	src := `package sample

import (
	"net/http"
	"os/exec"
	"path/filepath"
)

func RunCommand(w http.ResponseWriter, r *http.Request) {
	cmdName := r.URL.Query().Get("cmd")
	_ = exec.Command("sh", "-c", filepath.Clean(cmdName))
}
`
	findings := runCWE78(t, src, "CWE-78-taint-vulnerable.go")
	if !hasRule(findings, "CWE-78") {
		t.Fatalf("expected CWE-78 finding, got %#v", findings)
	}
}

func TestCWE78Safe(t *testing.T) {
	src := `package sample

import (
	"net/http"
	"os/exec"
)

func RunCommand(w http.ResponseWriter, r *http.Request) {
	_ = r
	_ = exec.Command("sh", "-c", "status")
}
`
	findings := runCWE78(t, src, "CWE-78-taint-safe.go")
	if hasRule(findings, "CWE-78") {
		t.Fatalf("safe fixture should not emit CWE-78, got %#v", findings)
	}
}

func TestCWE78RespectsAllows(t *testing.T) {
	src := `package sample
import ("net/http"; "os/exec")
func H(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("c")
	_ = exec.Command("bash", "-c", cmd)
}
`
	d := cwe.NewCWE78()
	unit := core.NewParsedUnit(core.LangGo, "x.go", src)
	ctx := core.DefaultScanContext()
	ctx.Skip = []string{"CWE-78"}
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	if len(out) != 0 {
		t.Fatalf("skip should suppress findings: %#v", out)
	}
}

func runCWE78(t *testing.T, src, path string) []rules.Finding {
	t.Helper()
	d := cwe.NewCWE78()
	unit := core.NewParsedUnit(core.LangGo, path, src)
	ctx := core.DefaultScanContext()
	var out []rules.Finding
	d.BeginScan(ctx)
	d.Run(ctx, unit, &out)
	d.Finalize(ctx, &out)
	d.EndScan()
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

package taint_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/fixture"
	"github.com/chinmay-sawant/goslop/internal/lang/go/detectors/cwe"
	"github.com/chinmay-sawant/goslop/internal/lang/go/detectors/cwe/taint"
	"github.com/chinmay-sawant/goslop/internal/lang/go/goparse"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Quarantined honest FNs: known unsupported shapes (see documents/taint.md).
// These vulnerable fixtures must stay silent until a stronger contract ships.
var quarantinedVulnerable = map[string]string{
	"IP-010-vulnerable.txt": "cross-goroutine channel handoff not modeled (G5 v0)",
}

func taintFixtureDir(t *testing.T) string {
	t.Helper()
	// Walk up from this package to repo root.
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		candidate := filepath.Join(dir, "tests", "fixtures", "go", "taint")
		if st, err := os.Stat(candidate); err == nil && st.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find tests/fixtures/go/taint")
		}
		dir = parent
	}
}

func newTaintDetector() *taint.Detector {
	return taint.NewDetector(taint.MetaSet{
		CWE22: &cwe.MetaCWE22,
		CWE78: &cwe.MetaCWE78,
		CWE79: &cwe.MetaCWE79,
		CWE89: &cwe.MetaCWE89,
	})
}

func taintCtx() *core.ScanContext {
	ctx := core.DefaultScanContext()
	ctx.TaintEnabled = true
	ctx.TaintMaxDepth = 4
	ctx.Only = []string{"CWE-22", "CWE-78", "CWE-79", "CWE-89"}
	return ctx
}

func runTaintSource(t *testing.T, path, source string) []rules.Finding {
	t.Helper()
	tree, err := goparse.Parse([]byte(source))
	if tree == nil || tree.File == nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	defer tree.Close()
	unit := core.NewParsedUnitWithTree(core.LangGo, path, source, tree)
	unit.DisplayPath = path
	d := newTaintDetector()
	ctx := taintCtx()
	var out []rules.Finding
	d.BeginScan(ctx)
	d.Run(ctx, unit, &out)
	d.Finalize(ctx, &out)
	d.EndScan()
	return out
}

func TestTaintFixtures(t *testing.T) {
	dir := taintFixtureDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		name := e.Name()
		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatal(err)
			}
			fx, err := fixture.ParseFixture(string(raw), name)
			if err != nil {
				t.Fatal(err)
			}
			findings := runTaintSource(t, fx.Filename, fx.Source)
			wantVuln := strings.Contains(name, "vulnerable")
			if reason, q := quarantinedVulnerable[name]; q {
				if len(findings) != 0 {
					t.Fatalf("quarantined FN fixture must stay silent (%s); got %d findings: %#v", reason, len(findings), findings)
				}
				t.Logf("quarantined honest FN: %s", reason)
				return
			}
			if wantVuln {
				if len(findings) == 0 {
					t.Fatalf("expected finding(s) on vulnerable fixture, got none")
				}
			} else {
				if len(findings) != 0 {
					t.Fatalf("safe fixture should be silent, got %#v", findings)
				}
			}
		})
	}
}

func TestTaintProjects(t *testing.T) {
	// Locate taint_projects
	dir := taintFixtureDir(t)
	projects := filepath.Join(filepath.Dir(dir), "taint_projects")
	if st, err := os.Stat(projects); err != nil || !st.IsDir() {
		t.Skip("taint_projects not found")
	}
	entries, err := os.ReadDir(projects)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		t.Run(name, func(t *testing.T) {
			root := filepath.Join(projects, name)
			var units []*core.ParsedUnit
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() || !strings.HasSuffix(path, ".go") {
					return nil
				}
				src, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				tree, err := goparse.Parse(src)
				if tree == nil || tree.File == nil {
					return err
				}
				// Keep tree alive for scan — store as unit tree.
				rel, _ := filepath.Rel(root, path)
				unit := core.NewParsedUnitWithTree(core.LangGo, rel, string(src), tree)
				unit.DisplayPath = rel
				units = append(units, unit)
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				for _, u := range units {
					if u.Tree != nil {
						if c, ok := u.Tree.(interface{ Close() }); ok {
							c.Close()
						}
					}
				}
			}()

			d := newTaintDetector()
			ctx := taintCtx()
			var out []rules.Finding
			d.BeginScan(ctx)
			for _, u := range units {
				d.Run(ctx, u, &out)
			}
			d.Finalize(ctx, &out)
			d.EndScan()

			wantVuln := strings.Contains(name, "vulnerable")
			if wantVuln {
				if len(out) == 0 {
					t.Fatalf("expected inter-proc findings, got none")
				}
			} else if len(out) != 0 {
				t.Fatalf("safe project should be silent, got %#v", out)
			}
		})
	}
}

func TestRequiresCacheState(t *testing.T) {
	d := newTaintDetector()
	ctx := core.DefaultScanContext()
	if d.RequiresCacheState(ctx) {
		t.Fatal("requires_cache_state false when taint off")
	}
	ctx.TaintEnabled = true
	if !d.RequiresCacheState(ctx) {
		t.Fatal("requires_cache_state true when taint on")
	}
}

func TestIntraCWE22Direct(t *testing.T) {
	src := `package sample
import ("net/http"; "os"; "path/filepath")
func ServeFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	_ = os.Open(filepath.Clean(path))
}
`
	findings := runTaintSource(t, "x.go", src)
	if !hasRule(findings, "CWE-22") {
		t.Fatalf("expected CWE-22, got %#v", findings)
	}
}

func TestIntraCWE22SafeBase(t *testing.T) {
	src := `package sample
import ("net/http"; "os"; "path/filepath")
func ServeFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	safe := filepath.Base(path)
	_ = os.Open(safe)
}
`
	findings := runTaintSource(t, "x.go", src)
	if hasRule(findings, "CWE-22") {
		t.Fatalf("Base should sanitize, got %#v", findings)
	}
}

func hasRule(findings []rules.Finding, id string) bool {
	for _, f := range findings {
		if f.RuleID == id {
			return true
		}
	}
	return false
}

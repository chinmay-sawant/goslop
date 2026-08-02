package detectors_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/fixture"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Regression: PERF owns cost-shaped ORM/N+1; BP owns blocking-async / timeout hygiene.
func TestPythonPERFvsBPNeighborhoodOwnership(t *testing.T) {
	t.Parallel()

	t.Run("orm n+1 is PERF not BP", func(t *testing.T) {
		t.Parallel()
		fx := loadNeighborhoodFixture(t, "perf", "PERF-PY-2", true)
		findings := runAllPython(t, fx.body, fx.path, nil)
		assertHas(t, findings, "PERF-PY-2")
		assertLacksPrefix(t, findings, "BP-PY-")
	})

	t.Run("blocking sleep in async route is BP-PY-30", func(t *testing.T) {
		t.Parallel()
		fx := loadNeighborhoodFixture(t, "bp", "BP-PY-30", true)
		findings := runAllPython(t, fx.body, fx.path, nil)
		assertHas(t, findings, "BP-PY-30")
		// No PERF rule should claim the blocking-async smell.
		for _, f := range findings {
			if f.RuleID == "PERF-PY-10" {
				t.Fatalf("PERF-PY-10 should not own FastAPI blocking sleep; got %#v", findings)
			}
		}
	})
}

type neighborhoodFixture struct {
	path string
	body string
}

func loadNeighborhoodFixture(t *testing.T, pack, caseName string, vulnerable bool) neighborhoodFixture {
	t.Helper()
	suf := "safe"
	if vulnerable {
		suf = "vulnerable"
	}
	txtPath := filepath.Join(repoPythonFixturesRoot(t), pack, caseName+"-"+suf+".txt")
	data, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("read fixture %s: %v", txtPath, err)
	}
	fx, err := fixture.ParseFixture(string(data), txtPath)
	if err != nil {
		t.Fatalf("parse fixture %s: %v", txtPath, err)
	}
	path := fx.Filename
	if path == "" {
		path = caseName + "-" + suf + ".py"
	}
	return neighborhoodFixture{path: path, body: fx.Source}
}

func repoPythonFixturesRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", ".."))
	fx := filepath.Join(root, "tests", "fixtures", "python")
	if _, err := os.Stat(fx); err != nil {
		t.Fatalf("fixtures root %s: %v", fx, err)
	}
	return fx
}

func runAllPython(t *testing.T, src, path string, only []string) []rules.Finding {
	t.Helper()
	ctx := core.BuildScanContext(core.ProfileAll, only, nil)
	unit := core.NewParsedUnit(core.LanguagePython, path, src)
	var out []rules.Finding
	for _, det := range detectors.All() {
		det.Run(ctx, unit, &out)
	}
	return out
}

func assertHas(t *testing.T, findings []rules.Finding, id string) {
	t.Helper()
	for _, f := range findings {
		if f.RuleID == id {
			return
		}
	}
	t.Fatalf("missing %s in %#v", id, findings)
}

func assertLacksPrefix(t *testing.T, findings []rules.Finding, prefix string) {
	t.Helper()
	for _, f := range findings {
		if len(f.RuleID) >= len(prefix) && f.RuleID[:len(prefix)] == prefix {
			t.Fatalf("unexpected %s finding in cost-only snippet: %#v", prefix, findings)
		}
	}
}

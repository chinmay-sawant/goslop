package detectors_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Regression: PERF owns cost-shaped ORM/N+1; BP owns blocking-async / timeout hygiene.
func TestPythonPERFvsBPNeighborhoodOwnership(t *testing.T) {
	t.Parallel()

	t.Run("orm n+1 is PERF not BP", func(t *testing.T) {
		t.Parallel()
		src := "def reserve(items):\n    for item in items:\n        sku = Sku.objects.get(code=item['code'])\n"
		findings := runAllPython(t, src, "svc.py", nil)
		assertHas(t, findings, "PERF-PY-2")
		assertLacksPrefix(t, findings, "BP-PY-")
	})

	t.Run("blocking sleep in async route is BP-PY-30", func(t *testing.T) {
		t.Parallel()
		src := "from fastapi import FastAPI\nimport time\napp = FastAPI()\n\n@app.get('/slow')\nasync def slow():\n    time.sleep(1)\n    return {'ok': True}\n"
		findings := runAllPython(t, src, "main.py", nil)
		assertHas(t, findings, "BP-PY-30")
		// No PERF rule should claim the blocking-async smell.
		for _, f := range findings {
			if f.RuleID == "PERF-PY-10" {
				t.Fatalf("PERF-PY-10 should not own FastAPI blocking sleep; got %#v", findings)
			}
		}
	})
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

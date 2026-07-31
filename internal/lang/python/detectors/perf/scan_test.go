package perf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	perf "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/perf"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestPERFPYCatalogueRegistersAllSeededRules(t *testing.T) {
	detector := perf.NewPythonPerfScan()
	if detector.Language() != core.LanguagePython {
		t.Fatalf("Language = %v, want Python", detector.Language())
	}
	ids := detector.RuleIDs()
	if len(ids) != 22 {
		t.Fatalf("RuleIDs length = %d, want 22: %v", len(ids), ids)
	}
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		seen[id] = true
		meta := detector.MetadataFor(id)
		if meta == nil || meta.ID != id || meta.Pack != rules.PackPerformance {
			t.Fatalf("metadata for %s = %#v, want performance metadata", id, meta)
		}
	}
	for i := 1; i <= 22; i++ {
		want := fmt.Sprintf("PERF-PY-%d", i)
		if !seen[want] {
			t.Fatalf("missing registered %s", want)
		}
	}
	if perf.CatalogueSize() != 22 {
		t.Fatalf("CatalogueSize = %d, want 22", perf.CatalogueSize())
	}
	for _, id := range ids {
		if strings.HasPrefix(id, "PERF-") && !strings.HasPrefix(id, "PERF-PY-") {
			t.Fatalf("bare Go PERF id registered by Python scan: %q", id)
		}
	}
}

func TestPERFPYSkipsNonPythonAndEmptyUnits(t *testing.T) {
	detector := perf.NewPythonPerfScan()
	for _, unit := range []*core.ParsedUnit{
		core.NewParsedUnit(core.LangGo, "x.go", "items = db.all()"),
		core.NewParsedUnit(core.LanguagePython, "x.py", ""),
	} {
		var findings []rules.Finding
		detector.Run(core.DefaultScanContext(), unit, &findings)
		if len(findings) != 0 {
			t.Fatalf("unexpected findings for %s: %#v", unit.Path, findings)
		}
	}
}

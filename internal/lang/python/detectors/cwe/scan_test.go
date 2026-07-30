package cwe_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestPyCweScanLanguageAndCatalogue(t *testing.T) {
	t.Parallel()
	d := cwe.NewPyCweScan()
	if d.Language() != core.LanguagePython {
		t.Fatalf("Language() = %v, want LanguagePython", d.Language())
	}
	ids := d.RuleIDs()
	want := map[string]bool{
		"CWE-22": true, "CWE-78": true, "CWE-79": true, "CWE-89": true, "CWE-502": true,
	}
	if len(ids) != len(want) {
		t.Fatalf("RuleIDs() = %v (len %d), want 5 priority IDs", ids, len(ids))
	}
	for _, id := range ids {
		if !want[id] {
			t.Fatalf("unexpected rule id %q", id)
		}
		meta := d.MetadataFor(id)
		if meta == nil {
			t.Fatalf("MetadataFor(%s) = nil", id)
		}
		if meta.ID != id {
			t.Fatalf("meta.ID = %q, want %q", meta.ID, id)
		}
		if meta.Severity != rules.SeverityHigh {
			t.Fatalf("%s severity = %v, want High", id, meta.Severity)
		}
		if meta.Pack != rules.PackSecurity {
			t.Fatalf("%s pack = %v, want PackSecurity", id, meta.Pack)
		}
	}
	if cwe.RegisteredRuleCount() != 5 {
		t.Fatalf("RegisteredRuleCount = %d, want 5", cwe.RegisteredRuleCount())
	}
}

func TestEmptySourceEmitsZeroFindings(t *testing.T) {
	t.Parallel()
	d := cwe.NewPyCweScan()
	var out []rules.Finding
	unit := core.NewParsedUnit(core.LanguagePython, "empty.py", "")
	d.Run(nil, unit, &out)
	if len(out) != 0 {
		t.Fatalf("empty source findings = %d, want 0", len(out))
	}
	d.Run(nil, nil, &out)
	if len(out) != 0 {
		t.Fatalf("nil unit findings = %d, want 0", len(out))
	}
}

func TestAllowsSkipsDisallowedRules(t *testing.T) {
	t.Parallel()
	d := cwe.NewPyCweScan()
	src := "import pickle\npickle.loads(data)\n"
	unit := core.NewParsedUnit(core.LanguagePython, "x.py", src)
	// Only list that matches nothing → Allows false for all registered rules.
	ctx := &core.ScanContext{Only: []string{"__none__"}}
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	if len(out) != 0 {
		t.Fatalf("deny-all findings = %v, want 0", out)
	}
}

func TestNoNeedleFileEmitsZero(t *testing.T) {
	t.Parallel()
	d := cwe.NewPyCweScan()
	src := "def hello():\n    return 1 + 2\n"
	unit := core.NewParsedUnit(core.LanguagePython, "clean.py", src)
	var out []rules.Finding
	d.Run(nil, unit, &out)
	if len(out) != 0 {
		t.Fatalf("no-sink findings = %v, want 0", out)
	}
}

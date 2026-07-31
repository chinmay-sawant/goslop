package cwe_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
	"github.com/chinmay-sawant/goslop/tests/integration"
)

func TestPyCweScanLanguageAndCatalogue(t *testing.T) {
	t.Parallel()
	d := cwe.NewPyCweScan()
	if d.Language() != core.LanguagePython {
		t.Fatalf("Language() = %v, want LanguagePython", d.Language())
	}
	const wantRuleCount = 159
	ids := d.RuleIDs()
	if len(ids) != wantRuleCount {
		t.Fatalf("RuleIDs() = %v (len %d), want %d supported IDs", ids, len(ids), wantRuleCount)
	}
	for _, id := range ids {
		meta := d.MetadataFor(id)
		if meta == nil {
			t.Fatalf("MetadataFor(%s) = nil", id)
		}
		if meta.ID != id {
			t.Fatalf("meta.ID = %q, want %q", meta.ID, id)
		}
		if meta.Severity != rules.SeverityHigh && meta.Severity != rules.SeverityMedium {
			t.Fatalf("%s severity = %v, want High or Medium", id, meta.Severity)
		}
		if meta.Pack != rules.PackSecurity {
			t.Fatalf("%s pack = %v, want PackSecurity", id, meta.Pack)
		}
	}
	if cwe.RegisteredRuleCount() != wantRuleCount {
		t.Fatalf("RegisteredRuleCount = %d, want %d", cwe.RegisteredRuleCount(), wantRuleCount)
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
	opts := integration.DefaultMatrixOptions()
	opts.Only = []string{"__none__"}
	opts.Languages = []core.LanguageID{core.LanguagePython}
	result, err := integration.MaterializeAndScanOpts(integration.PythonCWEFixtureRel("CWE-502", true), opts)
	if err != nil {
		t.Fatalf("scan fixture with deny-all filter: %v", err)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("deny-all findings = %v, want 0", result.Findings)
	}
}

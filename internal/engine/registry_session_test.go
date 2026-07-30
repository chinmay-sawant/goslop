package engine

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

type catalogueParityDetector struct {
	core.BaseDetector
	ruleIDs []string
}

func (d *catalogueParityDetector) Language() core.LanguageID { return core.LangGo }

func (d *catalogueParityDetector) RuleIDs() []string { return d.ruleIDs }

func (*catalogueParityDetector) Run(*core.ScanContext, *core.ParsedUnit, *[]rules.Finding) {}

type catalogueParityPlugin struct {
	core.BasePlugin
	catalogue []core.Detector
	session   []core.Detector
}

func (*catalogueParityPlugin) ID() core.LanguageID { return core.LangGo }

func (*catalogueParityPlugin) Extensions() []string { return []string{"go"} }

func (p *catalogueParityPlugin) Detectors() []core.Detector { return p.catalogue }

func (p *catalogueParityPlugin) NewDetectors() []core.Detector { return p.session }

func TestNewScanSessionRejectsDivergentRuleMultiplicity(t *testing.T) {
	plugin := &catalogueParityPlugin{
		catalogue: []core.Detector{&catalogueParityDetector{ruleIDs: []string{"TEST-A", "TEST-B"}}},
		session:   []core.Detector{&catalogueParityDetector{ruleIDs: []string{"TEST-A", "TEST-A"}}},
	}
	registry, err := NewRegistry([]core.LanguagePlugin{plugin})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.newScanSession(); err == nil {
		t.Fatal("expected divergent session catalogue to be rejected")
	}
}

func TestNewScanSessionAcceptsFreshMatchingCatalogue(t *testing.T) {
	plugin := &catalogueParityPlugin{
		catalogue: []core.Detector{&catalogueParityDetector{ruleIDs: []string{"TEST-A", "TEST-B"}}},
		session:   []core.Detector{&catalogueParityDetector{ruleIDs: []string{"TEST-A", "TEST-B"}}},
	}
	registry, err := NewRegistry([]core.LanguagePlugin{plugin})
	if err != nil {
		t.Fatal(err)
	}
	session, err := registry.newScanSession()
	if err != nil {
		t.Fatal(err)
	}
	if session.DetectorCount() != 1 || session.Detector(0) == plugin.catalogue[0] {
		t.Fatalf("unexpected session detectors: %#v", session.Detectors())
	}
}

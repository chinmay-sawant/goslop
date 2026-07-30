package engine_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/lang/python"
)

func TestDefaultRegistryEmptyOrRegistered(t *testing.T) {
	reg := engine.DefaultRegistry()
	if reg == nil {
		t.Fatal("DefaultRegistry must not be nil")
	}
	// Production may register a Go plugin later; MVP allows empty.
	if p, ok := reg.Plugin(core.LangGo); ok && p == nil {
		t.Fatal("Plugin ok but nil")
	}
}

func TestDefaultRegistryGoOnly(t *testing.T) {
	reg := engine.DefaultRegistry()
	if reg == nil {
		t.Fatal("DefaultRegistry must not be nil")
	}
	if _, ok := reg.Plugin(core.LanguageGo); !ok {
		t.Fatal("DefaultRegistry must register Go")
	}
	if _, ok := reg.Plugin(core.LanguagePython); ok {
		t.Fatal("DefaultRegistry must not register Python (production default is Go-only)")
	}
	plugins := reg.Plugins()
	for _, p := range plugins {
		if p == nil {
			t.Fatal("nil plugin in DefaultRegistry")
		}
		if p.ID() != core.LanguageGo {
			t.Fatalf("unexpected plugin language %v in DefaultRegistry", p.ID())
		}
	}
	ext := reg.ExtensionMap()
	if _, ok := ext["py"]; ok {
		t.Fatal("DefaultRegistry must not claim .py extension")
	}
	if lang, ok := ext["go"]; !ok || lang != core.LanguageGo {
		t.Fatalf("ExtensionMap go = (%v, %v)", lang, ok)
	}
}

func TestRegistryIndexesDetectorsByLanguageID(t *testing.T) {
	// Go-only default: all detectors keyed under LanguageGo.
	reg := engine.DefaultRegistry()
	goDet := reg.DetectorsForLanguage(core.LanguageGo)
	if len(goDet) == 0 {
		t.Fatal("expected Go detectors in DefaultRegistry")
	}
	for _, d := range goDet {
		if d.Language() != core.LanguageGo {
			t.Fatalf("detector language %v, want go", d.Language())
		}
	}
	if n := len(reg.DetectorsForLanguage(core.LanguagePython)); n != 0 {
		t.Fatalf("python detector indices should be empty without plugin, got %d", n)
	}

	// Custom multi-language registry: python stub adds zero detectors under LanguagePython.
	multi, err := engine.NewRegistry([]core.LanguagePlugin{
		// Reuse only python here to avoid re-registering Go detectors twice if DefaultRegistry shared.
		python.NewPlugin(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if n := len(multi.DetectorsForLanguage(core.LanguagePython)); n != 0 {
		t.Fatalf("python stub detectors = %d, want 0", n)
	}
	if _, ok := multi.Plugin(core.LanguagePython); !ok {
		t.Fatal("python plugin missing from multi registry")
	}
}

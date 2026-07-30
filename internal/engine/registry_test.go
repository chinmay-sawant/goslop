package engine_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
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

func TestDefaultRegistryIsGoOnly(t *testing.T) {
	reg := engine.DefaultRegistry()
	if _, ok := reg.Plugin(core.LanguagePython); ok {
		t.Fatal("DefaultRegistry must remain Go-only (Python is opt-in)")
	}
	if p := reg.PluginForPath("mod.py"); p != nil {
		t.Fatalf("PluginForPath(.py) = %v, want nil on default registry", p.ID())
	}
}

func TestNewRegistryWithLanguagesIncludesPython(t *testing.T) {
	reg, err := engine.NewRegistryWithLanguages(core.LanguageGo, core.LanguagePython)
	if err != nil {
		t.Fatal(err)
	}
	if p := reg.PluginForPath("a.go"); p == nil || p.ID() != core.LanguageGo {
		t.Fatal("expected Go for .go")
	}
	if p := reg.PluginForPath("b.py"); p == nil || p.ID() != core.LanguagePython {
		t.Fatal("expected Python for .py")
	}
}

func TestNewRegistryWithLanguagesUnknown(t *testing.T) {
	const unknown core.LanguageID = 99
	if _, err := engine.NewRegistryWithLanguages(unknown); err == nil {
		t.Fatal("expected error for unknown language id")
	}
}

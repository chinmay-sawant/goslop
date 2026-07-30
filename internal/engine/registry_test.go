package engine_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	golang "github.com/chinmay-sawant/goslop/internal/lang/go"
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
	if p := reg.PluginForPath("mod.py"); p != nil {
		t.Fatalf("PluginForPath(.py) = %v, want nil on default registry", p.ID())
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

func TestRegistryIndexesDetectorsByLanguageID(t *testing.T) {
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

	multi, err := engine.NewRegistry([]core.LanguagePlugin{
		python.NewPlugin(),
	})
	if err != nil {
		t.Fatal(err)
	}
	pyDet := multi.DetectorsForLanguage(core.LanguagePython)
	if len(pyDet) == 0 {
		t.Fatal("python plugin should register BP detectors")
	}
	for _, d := range pyDet {
		if d.Language() != core.LanguagePython {
			t.Fatalf("detector language %v, want python", d.Language())
		}
	}
	if _, ok := multi.Plugin(core.LanguagePython); !ok {
		t.Fatal("python plugin missing from multi registry")
	}
}

func TestRegistryForLanguagesFiltersPlugins(t *testing.T) {
	base, err := engine.NewRegistry([]core.LanguagePlugin{
		golang.NewPlugin(),
		python.NewPlugin(),
	})
	if err != nil {
		t.Fatal(err)
	}

	goOnly, err := engine.RegistryForLanguages(base, []core.LanguageID{core.LanguageGo})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := goOnly.Plugin(core.LanguageGo); !ok {
		t.Fatal("expected Go plugin")
	}
	if _, ok := goOnly.Plugin(core.LanguagePython); ok {
		t.Fatal("python plugin should be filtered out")
	}
	if _, ok := goOnly.ExtensionMap()["py"]; ok {
		t.Fatal("py extension should not be in go-only map")
	}
	if _, ok := goOnly.ExtensionMap()["go"]; !ok {
		t.Fatal("go extension missing")
	}

	pyOnly, err := engine.RegistryForLanguages(base, []core.LanguageID{core.LanguagePython})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := pyOnly.Plugin(core.LanguageGo); ok {
		t.Fatal("go should be filtered out")
	}
	if _, ok := pyOnly.Plugin(core.LanguagePython); !ok {
		t.Fatal("expected python plugin")
	}
	if pyOnly.DetectorCount() == 0 {
		t.Fatal("python-only registry should include BP detectors")
	}
}

func TestRegistryForLanguagesSkipsMissingPlugin(t *testing.T) {
	// Default production registry has Go only. Enabling python must not crash
	// when the plugin is not in the base registry.
	base := engine.DefaultRegistry()
	reg, err := engine.RegistryForLanguages(base, []core.LanguageID{core.LanguagePython})
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.Plugins()) != 0 {
		t.Fatalf("expected empty registry when python plugin absent, got %d plugins", len(reg.Plugins()))
	}
	// Enabling go+python keeps go and skips missing python from DefaultRegistry.
	reg, err = engine.RegistryForLanguages(base, []core.LanguageID{core.LanguageGo, core.LanguagePython})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.Plugin(core.LanguageGo); !ok {
		t.Fatal("expected Go plugin")
	}
	if _, ok := reg.Plugin(core.LanguagePython); ok {
		t.Fatal("python must not appear without a registered plugin on base")
	}
}

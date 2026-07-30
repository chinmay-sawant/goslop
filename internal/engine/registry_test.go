package engine_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	golang "github.com/chinmay-sawant/goslop/internal/lang/go"
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

// stubPythonPlugin is a zero-detector LanguagePlugin used only in registry tests.
type stubPythonPlugin struct {
	core.BasePlugin
}

func (stubPythonPlugin) ID() core.LanguageID        { return core.LanguagePython }
func (stubPythonPlugin) Extensions() []string       { return []string{"py"} }
func (stubPythonPlugin) Detectors() []core.Detector { return nil }

func (p stubPythonPlugin) ParseSource(path, source string) (*core.ParseResult, error) {
	return &core.ParseResult{
		Unit:    core.NewParsedUnit(core.LanguagePython, path, source),
		Quality: core.ParseQualitySourceOnly,
	}, nil
}

func TestRegistryForLanguagesFiltersPlugins(t *testing.T) {
	base, err := engine.NewRegistry([]core.LanguagePlugin{
		golang.NewPlugin(),
		stubPythonPlugin{},
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
	// Zero detectors is fine for a stub.
	if pyOnly.DetectorCount() != 0 {
		t.Fatalf("detector count=%d", pyOnly.DetectorCount())
	}
}

func TestRegistryForLanguagesSkipsMissingPlugin(t *testing.T) {
	// Default production registry has Go only. Enabling python must not crash.
	base := engine.DefaultRegistry()
	reg, err := engine.RegistryForLanguages(base, []core.LanguageID{core.LanguagePython})
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.Plugins()) != 0 {
		t.Fatalf("expected empty registry when python plugin absent, got %d plugins", len(reg.Plugins()))
	}
	// Enabling go+python keeps go and skips missing python.
	reg, err = engine.RegistryForLanguages(base, []core.LanguageID{core.LanguageGo, core.LanguagePython})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.Plugin(core.LanguageGo); !ok {
		t.Fatal("expected Go plugin")
	}
	if _, ok := reg.Plugin(core.LanguagePython); ok {
		t.Fatal("python must not appear without a registered plugin")
	}
}

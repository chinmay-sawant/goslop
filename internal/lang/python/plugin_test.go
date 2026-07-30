package python_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/lang/python"
)

func TestPythonPluginIdentityAndParse(t *testing.T) {
	t.Parallel()
	p := python.NewPlugin()
	if p.ID() != core.LanguagePython {
		t.Fatalf("ID() = %v, want python", p.ID())
	}
	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != "py" {
		t.Fatalf("Extensions() = %v, want [py]", exts)
	}
	if len(p.Detectors()) != 0 {
		t.Fatalf("Detectors() must be empty, got %d", len(p.Detectors()))
	}
	if len(p.NewDetectors()) != 0 {
		t.Fatalf("NewDetectors() must be empty, got %d", len(p.NewDetectors()))
	}

	result, err := p.ParseSource("sample.py", "x = 1\n")
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || result.Unit == nil {
		t.Fatal("expected source-only unit")
	}
	if result.Unit.Language != core.LanguagePython {
		t.Fatalf("unit.Language = %v, want python (BasePlugin default is Go)", result.Unit.Language)
	}
	if result.Quality != core.ParseQualitySourceOnly {
		t.Fatalf("quality = %v, want source-only", result.Quality)
	}
	if result.Unit.Tree != nil {
		t.Fatal("stub must not attach an AST tree")
	}
}

func TestPythonPluginRegisterInCustomRegistryOnly(t *testing.T) {
	// DefaultRegistry must remain Go-only even if this package is imported.
	def := engine.DefaultRegistry()
	if _, ok := def.Plugin(core.LanguagePython); ok {
		t.Fatal("DefaultRegistry must not include Python plugin")
	}
	if _, ok := def.Plugin(core.LanguageGo); !ok {
		t.Fatal("DefaultRegistry must include Go plugin")
	}

	reg, err := engine.NewRegistry([]core.LanguagePlugin{python.NewPlugin()})
	if err != nil {
		t.Fatal(err)
	}
	p, ok := reg.Plugin(core.LanguagePython)
	if !ok || p == nil {
		t.Fatal("custom registry should register Python stub")
	}
	extMap := reg.ExtensionMap()
	if lang, ok := extMap["py"]; !ok || lang != core.LanguagePython {
		t.Fatalf("ExtensionMap py = (%v, %v), want python", lang, ok)
	}
	if n := len(reg.DetectorsForLanguage(core.LanguagePython)); n != 0 {
		t.Fatalf("python detectors = %d, want 0", n)
	}
}

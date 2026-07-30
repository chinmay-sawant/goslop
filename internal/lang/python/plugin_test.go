package python_test

import (
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/lang/python"
)

func TestPythonPluginIdentityAndParse(t *testing.T) {
	t.Parallel()
	p := python.NewPlugin()
	if p.ID() != core.LanguagePython {
		t.Fatalf("ID() = %v, want LanguagePython", p.ID())
	}
	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != "py" {
		t.Fatalf("Extensions() = %v, want [py]", exts)
	}
	if len(p.Detectors()) == 0 {
		t.Fatal("Detectors() empty; expected BP-PY detector bundle")
	}
	if len(p.NewDetectors()) == 0 {
		t.Fatal("NewDetectors() empty; expected BP-PY detector bundle")
	}

	const src = "def hello():\n    print('hi')\n"
	result, err := p.ParseSource("sample.py", src)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || result.Unit == nil {
		t.Fatal("expected parse result with unit")
	}
	if result.Unit.Language != core.LanguagePython {
		t.Fatalf("unit.Language = %v, want LanguagePython (not Go)", result.Unit.Language)
	}
	if result.Quality != core.ParseQualitySourceOnly {
		t.Fatalf("quality = %v, want source-only", result.Quality)
	}
	if result.Unit.Tree != nil {
		t.Fatal("stub must not attach an AST tree")
	}
	if result.Unit.Source != src {
		t.Fatal("source mismatch")
	}
	if result.Unit.Path != "sample.py" {
		t.Fatalf("path = %q", result.Unit.Path)
	}
}

func TestPythonPluginHasBPPYRules(t *testing.T) {
	t.Parallel()
	p := python.NewPlugin()
	dets := p.NewDetectors()
	if len(dets) == 0 {
		t.Fatal("expected at least one Python detector")
	}
	var foundBP bool
	for _, d := range dets {
		if d.Language() != core.LanguagePython {
			t.Errorf("detector language = %v, want Python", d.Language())
		}
		for _, id := range d.RuleIDs() {
			if strings.HasPrefix(id, "BP-PY-") {
				foundBP = true
				meta := d.MetadataFor(id)
				if meta == nil {
					t.Errorf("MetadataFor(%s) = nil", id)
				}
			}
			if strings.HasPrefix(id, "BP-") && !strings.HasPrefix(id, "BP-PY-") {
				t.Errorf("bare Go BP id on Python detector: %s", id)
			}
		}
	}
	if !foundBP {
		t.Fatal("expected at least one BP-PY-* rule id on Python detectors")
	}
}

func TestRegistryWithPythonResolvesPy(t *testing.T) {
	reg, err := engine.NewRegistryWithLanguages(core.LanguageGo, core.LanguagePython)
	if err != nil {
		t.Fatal(err)
	}
	if p := reg.PluginForPath("pkg/main.go"); p == nil || p.ID() != core.LanguageGo {
		t.Fatalf("PluginForPath(main.go) = %v, want Go", pluginID(p))
	}
	if p := reg.PluginForPath("app/hello.py"); p == nil || p.ID() != core.LanguagePython {
		t.Fatalf("PluginForPath(hello.py) = %v, want Python", pluginID(p))
	}
	py, ok := reg.Plugin(core.LanguagePython)
	if !ok || py == nil {
		t.Fatal("Python plugin missing from multi-language registry")
	}
}

func TestDefaultRegistryRemainsGoOnly(t *testing.T) {
	reg := engine.DefaultRegistry()
	if _, ok := reg.Plugin(core.LanguagePython); ok {
		t.Fatal("DefaultRegistry must not register Python by default")
	}
	if p := reg.PluginForPath("hello.py"); p != nil {
		t.Fatalf("DefaultRegistry PluginForPath(.py) = %v, want nil", p.ID())
	}
	if p, ok := reg.Plugin(core.LanguageGo); !ok || p == nil {
		t.Fatal("DefaultRegistry must still include Go")
	}
}

func TestRegisterHook(t *testing.T) {
	reg, err := engine.NewRegistry(nil)
	if err != nil {
		t.Fatal(err)
	}
	python.Register(reg)
	if p, ok := reg.Plugin(core.LanguagePython); !ok || p == nil {
		t.Fatal("Register did not add Python plugin")
	}
	if p := reg.PluginForPath("x.py"); p == nil || p.ID() != core.LanguagePython {
		t.Fatal("Register: PluginForPath(.py) failed")
	}
}

func TestPythonPluginRegisterInCustomRegistryOnly(t *testing.T) {
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
		t.Fatal("custom registry should register Python plugin")
	}
	extMap := reg.ExtensionMap()
	if lang, ok := extMap["py"]; !ok || lang != core.LanguagePython {
		t.Fatalf("ExtensionMap py = (%v, %v), want python", lang, ok)
	}
	if n := len(reg.DetectorsForLanguage(core.LanguagePython)); n == 0 {
		t.Fatal("python detectors = 0, want BP-PY detector(s)")
	}
}

func pluginID(p core.LanguagePlugin) any {
	if p == nil {
		return nil
	}
	return p.ID()
}

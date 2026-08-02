package python_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestPythonSlimCorpusExpectedFindings(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..", "fixtures", "python", "corpus")
	apps, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	reg, err := engine.NewRegistryWithLanguages(core.LanguagePython)
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	for _, app := range apps {
		if !app.IsDir() {
			continue
		}
		app := app
		t.Run(app.Name(), func(t *testing.T) {
			t.Parallel()
			dir := filepath.Join(root, app.Name())
			raw, err := os.ReadFile(filepath.Join(dir, "expected.json"))
			if err != nil {
				t.Fatalf("expected.json: %v", err)
			}
			var want map[string][]string
			if err := json.Unmarshal(raw, &want); err != nil {
				t.Fatalf("parse expected: %v", err)
			}
			for file, ruleIDs := range want {
				src, err := os.ReadFile(filepath.Join(dir, file))
				if err != nil {
					t.Fatalf("read %s: %v", file, err)
				}
				ctx := core.BuildScanContext(core.ProfileAll, ruleIDs, nil)
				unit := core.NewParsedUnit(core.LanguagePython, file, string(src))
				var findings []rules.Finding
				for _, det := range reg.DetectorsForLanguage(core.LanguagePython) {
					det.Run(ctx, unit, &findings)
				}
				got := map[string]bool{}
				for _, f := range findings {
					got[f.RuleID] = true
				}
				for _, id := range ruleIDs {
					if !got[id] {
						t.Fatalf("%s/%s missing %s; got %#v", app.Name(), file, id, findings)
					}
				}
			}
		})
	}
}

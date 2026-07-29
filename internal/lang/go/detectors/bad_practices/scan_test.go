package badpractices_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/fixture"
	badpractices "github.com/chinmay/codehound/internal/lang/go/detectors/bad_practices"
	"github.com/chinmay/codehound/internal/rules"
)

func TestBPRulesRegistered(t *testing.T) {
	d := badpractices.NewGoBadPracticeScan()
	ids := d.RuleIDs()
	if len(ids) < 40 {
		t.Fatalf("expected many BP rules registered, got %d: %v", len(ids), ids)
	}
	for _, want := range []string{"BP-1", "BP-5", "BP-3", "BP-10", "BP-11", "BP-47", "BP-50"} {
		if !contains(ids, want) {
			t.Errorf("missing registered rule %s", want)
		}
	}
	if badpractices.CatalogueSize() < 100 {
		t.Errorf("catalogue size too small: %d", badpractices.CatalogueSize())
	}
}

func TestRecommendedProfileKeepsBPOff(t *testing.T) {
	ctx := core.BuildScanContext(core.ProfileRecommended, nil, nil)
	if ctx.BadPracticesEnabled {
		t.Fatal("recommended profile must keep BadPracticesEnabled=false")
	}
	if ctx.Allows("BP-1") {
		t.Fatal("recommended profile must not allow BP-1")
	}
	style := core.BuildScanContext(core.ProfileStyle, nil, nil)
	if !style.BadPracticesEnabled || !style.Allows("BP-1") {
		t.Fatal("style profile should enable BP rules")
	}
}

func TestBPSeverityOverride(t *testing.T) {
	sev := rules.SeverityHigh
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	ctx.BadPracticeSeverity = &sev
	src := `package sample
import "os"
func WriteFile(path string, data []byte) {
	_ = os.WriteFile(path, data, 0644)
}
`
	findings := runBP(t, ctx, src, "BP-1-override.go")
	if !hasRule(findings, "BP-1") {
		t.Fatalf("expected BP-1, got %#v", findings)
	}
	for _, f := range findings {
		if f.RuleID == "BP-1" && f.Severity != rules.SeverityHigh {
			t.Fatalf("expected severity override high, got %v", f.Severity)
		}
	}
}

func TestBP1Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-1-vulnerable.txt", "BP-1", true)
}

func TestBP1Safe(t *testing.T) {
	runFixtureRule(t, "BP-1-safe.txt", "BP-1", false)
}

func TestBP5Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-5-vulnerable.txt", "BP-5", true)
}

func TestBP5Safe(t *testing.T) {
	runFixtureRule(t, "BP-5-safe.txt", "BP-5", false)
}

func TestBP2Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-2-vulnerable.txt", "BP-2", true)
}

func TestBP3Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-3-vulnerable.txt", "BP-3", true)
}

func TestBP3Safe(t *testing.T) {
	runFixtureRule(t, "BP-3-safe.txt", "BP-3", false)
}

func TestBP4Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-4-vulnerable.txt", "BP-4", true)
}

func TestBP6Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-6-vulnerable.txt", "BP-6", true)
}

func TestBP7Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-7-vulnerable.txt", "BP-7", true)
}

func TestBP10Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-10-vulnerable.txt", "BP-10", true)
}

func TestBP10Safe(t *testing.T) {
	runFixtureRule(t, "BP-10-safe.txt", "BP-10", false)
}

func TestBP11Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-11-vulnerable.txt", "BP-11", true)
}

func TestBP13Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-13-vulnerable.txt", "BP-13", true)
}

func TestBP16Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-16-vulnerable.txt", "BP-16", true)
}

func TestBP17Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-17-vulnerable.txt", "BP-17", true)
}

func TestBP26Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-26-vulnerable.txt", "BP-26", true)
}

func TestBP46Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-46-vulnerable.txt", "BP-46", true)
}

func TestBP48Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-48-vulnerable.txt", "BP-48", true)
}

func TestBP49Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-49-vulnerable.txt", "BP-49", true)
}

func TestBP80Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-80-vulnerable.txt", "BP-80", true)
}

func TestBP95Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-95-vulnerable.txt", "BP-95", true)
}

func TestBP98Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-98-vulnerable.txt", "BP-98", true)
}

func TestBP100Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-100-vulnerable.txt", "BP-100", true)
}

func TestBP110Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-110-vulnerable.txt", "BP-110", true)
}

func TestBP154Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-154-vulnerable.txt", "BP-154", true)
}

func TestBP160Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-160-vulnerable.txt", "BP-160", true)
}

func TestProjectBP47Vulnerable(t *testing.T) {
	runProjectRule(t, "BP-47-vulnerable", "BP-47", true)
}

func TestProjectBP47Safe(t *testing.T) {
	runProjectRule(t, "BP-47-safe", "BP-47", false)
}

func TestProjectBP50Vulnerable(t *testing.T) {
	runProjectRule(t, "BP-50-vulnerable", "BP-50", true)
}

func TestProjectBP50Safe(t *testing.T) {
	runProjectRule(t, "BP-50-safe", "BP-50", false)
}

func runFixtureRule(t *testing.T, name, rule string, wantFire bool) {
	t.Helper()
	data, path := readFixture(t, "bad_practices", name)
	fx, err := fixture.ParseFixture(string(data), path)
	if err != nil {
		t.Fatal(err)
	}
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	ctx.Only = []string{rule}
	findings := runBP(t, ctx, fx.Source, fx.Filename)
	got := hasRule(findings, rule)
	if got != wantFire {
		t.Fatalf("fixture %s rule %s want fire=%v got=%v findings=%v", name, rule, wantFire, got, findings)
	}
}

func runProjectRule(t *testing.T, project, rule string, wantFire bool) {
	t.Helper()
	root := findProjectFixture(t, project)
	mainPath := filepath.Join(root, "main.go")
	data, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatal(err)
	}
	d := badpractices.NewGoBadPracticeScan()
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	ctx.Only = []string{rule}
	d.BeginScan(ctx)
	defer d.EndScan()
	unit := core.NewParsedUnit(core.LangGo, mainPath, string(data))
	unit.DisplayPath = "main.go"
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	got := hasRule(out, rule)
	if got != wantFire {
		t.Fatalf("project %s rule %s want fire=%v got=%v findings=%v", project, rule, wantFire, got, out)
	}
}

func runBP(t *testing.T, ctx *core.ScanContext, src, path string) []rules.Finding {
	t.Helper()
	d := badpractices.NewGoBadPracticeScan()
	if ctx == nil {
		ctx = core.DefaultScanContext()
		ctx.BadPracticesEnabled = true
	}
	d.BeginScan(ctx)
	defer d.EndScan()
	unit := core.NewParsedUnit(core.LangGo, path, src)
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	return out
}

func hasRule(findings []rules.Finding, id string) bool {
	for _, f := range findings {
		if f.RuleID == id {
			return true
		}
	}
	return false
}

func contains(ids []string, want string) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
}

func readFixture(t *testing.T, sub, name string) ([]byte, string) {
	t.Helper()
	cwd, _ := os.Getwd()
	var (
		data []byte
		err  error
		path string
	)
	candidates := []string{
		filepath.Join(cwd, "..", "..", "..", "..", "..", "tests", "fixtures", "go", sub, name),
		filepath.Join(cwd, "tests", "fixtures", "go", sub, name),
	}
	for d := cwd; d != string(filepath.Separator) && d != "."; d = filepath.Dir(d) {
		candidates = append(candidates, filepath.Join(d, "tests", "fixtures", "go", sub, name))
		if filepath.Dir(d) == d {
			break
		}
	}
	for _, c := range candidates {
		data, err = os.ReadFile(c)
		if err == nil {
			path = c
			break
		}
	}
	if err != nil {
		t.Fatalf("read fixture %s: %v (cwd=%s)", name, err, cwd)
	}
	return data, path
}

func findProjectFixture(t *testing.T, project string) string {
	t.Helper()
	cwd, _ := os.Getwd()
	candidates := []string{
		filepath.Join(cwd, "..", "..", "..", "..", "..", "tests", "fixtures", "go", "bad_practices_projects", project),
		filepath.Join(cwd, "tests", "fixtures", "go", "bad_practices_projects", project),
	}
	for d := cwd; d != string(filepath.Separator) && d != "."; d = filepath.Dir(d) {
		candidates = append(candidates, filepath.Join(d, "tests", "fixtures", "go", "bad_practices_projects", project))
		if filepath.Dir(d) == d {
			break
		}
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c
		}
	}
	t.Fatalf("project fixture %s not found (cwd=%s); tried %s", project, cwd, strings.Join(candidates, ", "))
	return ""
}

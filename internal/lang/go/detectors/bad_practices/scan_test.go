package badpractices_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/fixture"
	badpractices "github.com/chinmay/goslop/internal/lang/go/detectors/bad_practices"
	"github.com/chinmay/goslop/internal/rules"
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

// TestBP28BP29UnbalancedBraceNeedleDoesNotHang covers a regression where the
// text scanners for BP-28/BP-29 matched "interface {" inside a string literal
// whose braces never balanced, then set start=abs and looped forever.
// Scanning rules_api.go itself used to hang make run when SCAN_PATH was self.
func TestBP28BP29UnbalancedBraceNeedleDoesNotHang(t *testing.T) {
	// Mimic detector source: needle appears in a string with no matching '}'.
	src := "package p\n" +
		"func f() {\n" +
		"\t_ = \"interface {\"\n" +
		"\t_ = \"interface{\"\n" +
		"\tx := 1\n" +
		"\t_ = x\n" +
		"}\n"
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	ctx.Only = []string{"BP-28", "BP-29"}
	done := make(chan []rules.Finding, 1)
	go func() {
		done <- runBP(t, ctx, src, "bp28_bp29_string_needle.go")
	}()
	select {
	case <-done:
		// Success: terminated (no findings required).
	case <-time.After(2 * time.Second):
		t.Fatal("BP-28/BP-29 hung on unbalanced \"interface {\" string needle")
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

func TestBP15IndirectSafe(t *testing.T) {
	runFixtureRule(t, "BP-15-indirect-safe.txt", "BP-15", false)
}
func TestBP15Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-15-vulnerable.txt", "BP-15", true)
}
func TestBP15IndirectVulnerable(t *testing.T) {
	runFixtureRule(t, "BP-15-indirect-vulnerable.txt", "BP-15", true)
}
func TestBP18Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-18-vulnerable.txt", "BP-18", true)
}
func TestBP18Safe(t *testing.T) {
	runFixtureRule(t, "BP-18-safe.txt", "BP-18", false)
}
func TestBP23Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-23-vulnerable.txt", "BP-23", true)
}
func TestBP23Safe(t *testing.T) {
	runFixtureRule(t, "BP-23-safe.txt", "BP-23", false)
}
func TestBP33Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-33-vulnerable.txt", "BP-33", true)
}
func TestBP33Safe(t *testing.T) {
	runFixtureRule(t, "BP-33-safe.txt", "BP-33", false)
}
func TestBP37Safe(t *testing.T) {
	runFixtureRule(t, "BP-37-safe.txt", "BP-37", false)
}
func TestBP37TypeSwitchVulnerable(t *testing.T) {
	runFixtureRule(t, "BP-37-type-switch-vulnerable.txt", "BP-37", true)
}
func TestBP37Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-37-vulnerable.txt", "BP-37", true)
}
func TestBP41Vulnerable(t *testing.T) {
	// Nested package path; materialize so packageDoc snapshot can see the file.
	data, path := readFixture(t, "bad_practices", "BP-41-vulnerable.txt")
	fx, err := fixture.ParseFixture(string(data), path)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	srcPath, err := fixture.MaterializeFixture(root, fx)
	if err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true
	ctx.Only = []string{"BP-41"}
	d := badpractices.NewGoBadPracticeScan()
	d.BeginScan(ctx)
	defer d.EndScan()
	unit := core.NewParsedUnit(core.LangGo, srcPath, string(src))
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	if !hasRule(out, "BP-41") {
		t.Fatalf("expected BP-41 on nested package fixture, got %#v path=%s", out, srcPath)
	}
}
func TestBP44Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-44-vulnerable.txt", "BP-44", true)
}
func TestBP44Safe(t *testing.T) {
	runFixtureRule(t, "BP-44-safe.txt", "BP-44", false)
}
func TestBP45Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-45-vulnerable.txt", "BP-45", true)
}
func TestBP53Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-53-vulnerable.txt", "BP-53", true)
}
func TestBP53Safe(t *testing.T) {
	runFixtureRule(t, "BP-53-safe.txt", "BP-53", false)
}
func TestBP73Safe(t *testing.T) {
	runFixtureRule(t, "BP-73-safe.txt", "BP-73", false)
}
func TestBP73Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-73-vulnerable.txt", "BP-73", true)
}
func TestBP75VariantVulnerable(t *testing.T) {
	runFixtureRule(t, "BP-75-variant-vulnerable.txt", "BP-75", true)
}
func TestBP75Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-75-vulnerable.txt", "BP-75", true)
}
func TestBP76Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-76-vulnerable.txt", "BP-76", true)
}
func TestBP76VariantVulnerable(t *testing.T) {
	runFixtureRule(t, "BP-76-variant-vulnerable.txt", "BP-76", true)
}
func TestBP76Safe(t *testing.T) {
	runFixtureRule(t, "BP-76-safe.txt", "BP-76", false)
}
func TestBP79VariantSafe(t *testing.T) {
	runFixtureRule(t, "BP-79-variant-safe.txt", "BP-79", false)
}
func TestBP79Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-79-vulnerable.txt", "BP-79", true)
}
func TestBP81VariantSafe(t *testing.T) {
	runFixtureRule(t, "BP-81-variant-safe.txt", "BP-81", false)
}
func TestBP81Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-81-vulnerable.txt", "BP-81", true)
}
func TestBP84Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-84-vulnerable.txt", "BP-84", true)
}
func TestBP84VariantVulnerable(t *testing.T) {
	runFixtureRule(t, "BP-84-variant-vulnerable.txt", "BP-84", true)
}
func TestBP84Safe(t *testing.T) {
	runFixtureRule(t, "BP-84-safe.txt", "BP-84", false)
}
func TestBP87Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-87-vulnerable.txt", "BP-87", true)
}
func TestBP87VariantVulnerable(t *testing.T) {
	runFixtureRule(t, "BP-87-variant-vulnerable.txt", "BP-87", true)
}
func TestBP87Safe(t *testing.T) {
	runFixtureRule(t, "BP-87-safe.txt", "BP-87", false)
}
func TestBP88Safe(t *testing.T) {
	runFixtureRule(t, "BP-88-safe.txt", "BP-88", false)
}
func TestBP88VariantSafe(t *testing.T) {
	runFixtureRule(t, "BP-88-variant-safe.txt", "BP-88", false)
}
func TestBP88Vulnerable(t *testing.T) {
	runFixtureRule(t, "BP-88-vulnerable.txt", "BP-88", true)
}

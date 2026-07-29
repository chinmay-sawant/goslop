package cwe_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/fixture"
	"github.com/chinmay/goslop/internal/lang/go/detectors/cwe"
	"github.com/chinmay/goslop/internal/rules"
)

func TestRegistryCoverage(t *testing.T) {
	d := cwe.NewGoCweScan()
	ids := d.RuleIDs()
	if got := len(ids); got < 175 {
		t.Fatalf("expected >= 175 registered CWE rules, got %d", got)
	}
	// Every registry ID must be present.
	want := expectedRegistryIDs(t)
	have := map[string]bool{}
	for _, id := range ids {
		have[id] = true
	}
	var missing []string
	for _, id := range want {
		if !have[id] {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("missing %d registry rules: %v", len(missing), missing)
	}
}

func TestMetadataForAllRegistered(t *testing.T) {
	d := cwe.NewGoCweScan()
	for _, id := range d.RuleIDs() {
		m := d.MetadataFor(id)
		if m == nil {
			t.Errorf("nil metadata for %s", id)
			continue
		}
		if m.ID != id {
			t.Errorf("metadata id %q != %q", m.ID, id)
		}
		if m.Title == "" {
			t.Errorf("empty title for %s", id)
		}
	}
}

func TestSeedCWE78StillWorks(t *testing.T) {
	src := `package sample
import (
	"net/http"
	"os/exec"
	"path/filepath"
)
func RunCommand(w http.ResponseWriter, r *http.Request) {
	cmdName := r.URL.Query().Get("cmd")
	_ = exec.Command("sh", "-c", filepath.Clean(cmdName))
}
`
	findings := runCWE(t, src, "CWE-78-taint-vulnerable.go")
	if !hasRule(findings, "CWE-78") {
		t.Fatalf("expected CWE-78, got %#v", findings)
	}
}

func TestSeedCWE89StillWorks(t *testing.T) {
	src := `package sample
import (
	"database/sql"
	"fmt"
	"net/http"
)
func Q(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	name := r.URL.Query().Get("name")
	q := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
	_, _ = db.Query(q)
}
`
	findings := runCWE(t, src, "CWE-89-vulnerable.go")
	if !hasRule(findings, "CWE-89") {
		t.Fatalf("expected CWE-89, got %#v", findings)
	}
}

// Sample structural fixture matrix across domains.
func TestStructuralFixtureSample(t *testing.T) {
	cases := []struct {
		id   int
		fire bool
	}{
		{338, true},
		{338, false}, // handled below via kind
		{325, true},
		{366, true},
		{619, true},
		{601, true},
		{798, true},
		{1327, true},
		{1204, true},
	}
	// explicit list of vulnerable/safe pairs
	pairs := []int{325, 334, 335, 338, 342, 343, 366, 368, 434, 601, 619, 798, 820, 821, 918, 1204, 1240, 1327, 1333}
	for _, n := range pairs {
		n := n
		rule := fmt.Sprintf("CWE-%d", n)
		t.Run(rule+"/vulnerable", func(t *testing.T) {
			runCWEFixture(t, fmt.Sprintf("CWE-%d-vulnerable.txt", n), rule, true)
		})
		t.Run(rule+"/safe", func(t *testing.T) {
			runCWEFixture(t, fmt.Sprintf("CWE-%d-safe.txt", n), rule, false)
		})
	}
	_ = cases
}

func TestTaintLiteFixtures(t *testing.T) {
	// Prefer tests/fixtures/go/taint when present, else stdlib.
	pairs := []struct {
		rule string
		vuln string
		safe string
	}{
		{"CWE-22", "CWE-22-vulnerable.txt", "CWE-22-safe.txt"},
		{"CWE-78", "CWE-78-vulnerable.txt", "CWE-78-safe.txt"},
		{"CWE-79", "CWE-79-vulnerable.txt", "CWE-79-safe.txt"},
		{"CWE-89", "CWE-89-vulnerable.txt", "CWE-89-safe.txt"},
	}
	for _, p := range pairs {
		p := p
		t.Run(p.rule+"/vulnerable", func(t *testing.T) {
			runCWEFixtureAny(t, p.vuln, p.rule, true)
		})
		t.Run(p.rule+"/safe", func(t *testing.T) {
			runCWEFixtureAny(t, p.safe, p.rule, false)
		})
	}
}

func expectedRegistryIDs(t *testing.T) []string {
	t.Helper()
	// Hard-coded from registry TOMLs count; verify count == 175 via package.
	// Build from RuleIDs after registration — circular. Instead list known seed + count.
	// Parse registry dir.
	cwd, _ := os.Getwd()
	var regDir string
	for d := cwd; d != string(filepath.Separator) && d != "."; d = filepath.Dir(d) {
		cand := filepath.Join(d, "internal", "lang", "go", "detectors", "cwe", "registry")
		if st, err := os.Stat(cand); err == nil && st.IsDir() {
			regDir = cand
			break
		}
		if filepath.Dir(d) == d {
			break
		}
	}
	if regDir == "" {
		// package dir is .../cwe
		regDir = filepath.Join(cwd, "registry")
	}
	entries, err := os.ReadDir(regDir)
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	var ids []string
	seen := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(regDir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "cwe ") || strings.HasPrefix(line, "cwe=") {
				// cwe = N
				parts := strings.Split(line, "=")
				if len(parts) != 2 {
					continue
				}
				num := strings.TrimSpace(parts[1])
				id := "CWE-" + num
				if !seen[id] {
					seen[id] = true
					ids = append(ids, id)
				}
			}
		}
	}
	if len(ids) != 175 {
		t.Fatalf("registry parse got %d ids, want 175", len(ids))
	}
	return ids
}

func runCWEFixture(t *testing.T, name, rule string, wantFire bool) {
	t.Helper()
	runCWEFixtureIn(t, name, rule, wantFire, "stdlib")
}

func runCWEFixtureAny(t *testing.T, name, rule string, wantFire bool) {
	t.Helper()
	// try taint then stdlib
	if tryCWEFixture(t, name, rule, wantFire, "taint") {
		return
	}
	runCWEFixtureIn(t, name, rule, wantFire, "stdlib")
}

func tryCWEFixture(t *testing.T, name, rule string, wantFire bool, sub string) bool {
	t.Helper()
	path := findFixture(name, sub)
	if path == "" {
		return false
	}
	runCWEFixturePath(t, path, rule, wantFire)
	return true
}

func runCWEFixtureIn(t *testing.T, name, rule string, wantFire bool, sub string) {
	t.Helper()
	path := findFixture(name, sub)
	if path == "" {
		t.Fatalf("fixture %s not found under %s", name, sub)
	}
	runCWEFixturePath(t, path, rule, wantFire)
}

func runCWEFixturePath(t *testing.T, path, rule string, wantFire bool) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	fx, err := fixture.ParseFixture(string(data), path)
	if err != nil {
		t.Fatal(err)
	}
	findings := runCWE(t, fx.Source, fx.Filename)
	got := hasRule(findings, rule)
	if got != wantFire {
		t.Fatalf("fixture %s rule %s want fire=%v got=%v findings=%v", path, rule, wantFire, got, summarize(findings))
	}
}

func findFixture(name, sub string) string {
	cwd, _ := os.Getwd()
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
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

func runCWE(t *testing.T, src, path string) []rules.Finding {
	t.Helper()
	d := cwe.NewGoCweScan()
	unit := core.NewParsedUnit(core.LangGo, path, src)
	ctx := core.DefaultScanContext()
	var out []rules.Finding
	d.BeginScan(ctx)
	d.Run(ctx, unit, &out)
	d.Finalize(ctx, &out)
	d.EndScan()
	return out
}

func summarize(findings []rules.Finding) []string {
	var ids []string
	for _, f := range findings {
		ids = append(ids, f.RuleID)
	}
	return ids
}

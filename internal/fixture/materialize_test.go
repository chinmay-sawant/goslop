package fixture_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay/goslop/internal/fixture"
)

func fixturesRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join("..", "..", "tests", "fixtures")
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("fixtures root %s: %v", root, err)
	}
	return root
}

func TestMaterializeFixture_CWE22(t *testing.T) {
	root := t.TempDir()
	txt := filepath.Join(fixturesRoot(t), "go", "taint", "CWE-22-vulnerable.txt")
	out, err := fixture.MaterializeFixtureFile(txt, root)
	if err != nil {
		t.Fatalf("MaterializeFixtureFile: %v", err)
	}
	wantSuffix := filepath.Join("go", "CWE-22-taint-vulnerable.go")
	if !strings.HasSuffix(out, wantSuffix) {
		t.Errorf("out path %q does not end with %q", out, wantSuffix)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "func ServeFile") {
		t.Errorf("materialized body missing ServeFile")
	}
}

func TestMaterializeTree_TaintSubset(t *testing.T) {
	outRoot := t.TempDir()
	src := filepath.Join(fixturesRoot(t), "go", "taint")
	if err := fixture.MaterializeTree(src, outRoot); err != nil {
		t.Fatalf("MaterializeTree: %v", err)
	}
	// All go fixtures under taint should land in outRoot/go/
	goDir := filepath.Join(outRoot, "go")
	entries, err := os.ReadDir(goDir)
	if err != nil {
		t.Fatalf("read go dir: %v", err)
	}
	if len(entries) < 10 {
		t.Errorf("expected many materialized go files, got %d", len(entries))
	}
	// Spot-check a couple of known names
	for _, name := range []string{"CWE-22-taint-vulnerable.go", "CWE-22-taint-safe.go"} {
		p := filepath.Join(goDir, name)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}
}

func TestMaterializeTree_IncludesPythonTagged(t *testing.T) {
	outRoot := t.TempDir()
	src := filepath.Join(fixturesRoot(t), "python")
	if err := fixture.MaterializeTree(src, outRoot); err != nil {
		t.Fatalf("MaterializeTree: %v", err)
	}
	py := filepath.Join(outRoot, "python", "sample.py")
	if _, err := os.Stat(py); err != nil {
		t.Fatalf("python fixture not materialized under lang tag: %v", err)
	}
}

func TestSanitizeFilename(t *testing.T) {
	if _, err := fixture.SanitizeFilename("/abs/path.go"); err == nil {
		t.Error("absolute path should fail")
	}
	if _, err := fixture.SanitizeFilename("../escape.go"); err == nil {
		t.Error(".. should fail")
	}
	if _, err := fixture.SanitizeFilename("a/../../x.go"); err == nil {
		t.Error("nested .. should fail")
	}
	got, err := fixture.SanitizeFilename("sub/dir/file.go")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Clean("sub/dir/file.go") {
		t.Errorf("got %q", got)
	}
}

func TestMaterializeFixture_RejectsEscape(t *testing.T) {
	root := t.TempDir()
	_, err := fixture.MaterializeFixture(root, fixture.TextFixture{
		Language: fixture.LangGo,
		Filename: "../escape.go",
		Source:   "package p\n",
	})
	if err == nil {
		t.Fatal("expected escape error")
	}
}

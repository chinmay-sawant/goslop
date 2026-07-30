package badpractices

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay/goslop/internal/core"
)

func TestPackageTypeFactsAreScopedAndMemoizedPerScan(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "types.go")
	src := "package sample\n\ntype Reader interface { Read() }\n\ntype impl struct{}\nfunc (impl) Read() {}\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	unit := core.NewParsedUnit(core.LangGo, path, src)

	firstScan := newProjectCaches()
	first := packageTypeFactsForUnit(unit, firstScan)
	if again := packageTypeFactsForUnit(unit, firstScan); first != again {
		t.Fatal("package facts should be reused within one scan")
	}
	secondScan := newProjectCaches()
	if next := packageTypeFactsForUnit(unit, secondScan); next == first {
		t.Fatal("package facts must not cross scan sessions")
	}
}

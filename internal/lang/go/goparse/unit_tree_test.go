package goparse_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/go/goparse"
)

func TestTreeForUnit_ReusesAttachedTree(t *testing.T) {
	src := "package p\nfunc F() {}\n"
	first, err := goparse.Parse([]byte(src))
	if err != nil || first == nil || first.File == nil {
		t.Fatalf("parse: %v tree=%v", err, first)
	}
	unit := core.NewParsedUnitWithTree(core.LangGo, "p.go", src, first)

	got := goparse.TreeForUnit(unit)
	if got != first {
		t.Fatalf("TreeForUnit did not reuse unit.Tree: got %p want %p", got, first)
	}
}

func TestTreeForUnit_ParsesOnceAndPins(t *testing.T) {
	src := "package p\nfunc F() {}\n"
	unit := core.NewParsedUnit(core.LangGo, "p.go", src)
	if unit.Tree != nil {
		t.Fatal("expected nil Tree before first resolve")
	}

	a := goparse.TreeForUnit(unit)
	if a == nil || a.File == nil {
		t.Fatal("first TreeForUnit returned nil")
	}
	if unit.Tree != a {
		t.Fatal("TreeForUnit did not pin tree on unit")
	}

	b := goparse.TreeForUnit(unit)
	if b != a {
		t.Fatalf("second TreeForUnit re-parsed: %p vs %p", b, a)
	}
}

func TestAsTree(t *testing.T) {
	if goparse.AsTree(nil) != nil {
		t.Fatal("AsTree(nil)")
	}
	if goparse.AsTree("x") != nil {
		t.Fatal("AsTree(string)")
	}
	src := "package p\n"
	tr, err := goparse.Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if goparse.AsTree(tr) != tr {
		t.Fatal("AsTree usable tree")
	}
}

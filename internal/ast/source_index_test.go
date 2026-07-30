package ast_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/ast"
)

func TestSourceIndex_Has(t *testing.T) {
	needles := []string{"alpha", "beta", "gamma", "sync.Mutex", "select {"}
	idx := ast.Build("use sync.Mutex; select { default: }", needles)

	if !idx.Has("sync.Mutex") {
		t.Error("expected sync.Mutex")
	}
	if !idx.Has("select {") {
		t.Error("expected select {")
	}
	if idx.Has("alpha") {
		t.Error("alpha should be absent")
	}
	if idx.Has("missing") {
		t.Error("unknown needle should be false")
	}
	if idx.Len() != len(needles) {
		t.Errorf("len = %d", idx.Len())
	}
}

func TestSourceIndex_HasAny(t *testing.T) {
	idx := ast.Build("gamma only", []string{"alpha", "beta", "gamma"})
	if !idx.HasAny([]string{"alpha", "gamma"}) {
		t.Error("expected has_any true")
	}
	if idx.HasAny([]string{"alpha", "beta"}) {
		t.Error("expected has_any false")
	}
}

func TestSourceIndex_Empty(t *testing.T) {
	idx := ast.Build("src", nil)
	if !idx.Empty() {
		t.Error("empty needles should Empty()")
	}
	if idx.Has("x") {
		t.Error("no needles")
	}
}

package ast_test

import (
	"testing"

	"github.com/chinmay/codehound/internal/ast"
)

func TestSourceIndexStillWorksWithoutTreeSitter(t *testing.T) {
	idx := ast.Build("hello fmt.Sprintf(world)", []string{"fmt.Sprintf(", "missing"})
	if !idx.Has("fmt.Sprintf(") {
		t.Fatal("expected needle present")
	}
	if idx.Has("missing") {
		t.Fatal("expected needle absent")
	}
}

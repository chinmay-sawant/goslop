package ast_test

import (
	"testing"

	"github.com/chinmay/codehound/internal/ast"
)

func TestComputeLineStartsAndLineCol(t *testing.T) {
	src := "package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"
	starts := ast.ComputeLineStarts(src)
	// line 1 @ 0, line 2 after first \n, etc.
	if starts[0] != 0 {
		t.Fatalf("starts[0]=%d", starts[0])
	}

	// offset 0 → (1,1)
	line, col := ast.LineColWithStarts(starts, 0)
	if line != 1 || col != 1 {
		t.Errorf("offset 0 → (%d,%d), want (1,1)", line, col)
	}

	// first char of "func" is line 3
	idx := 0
	for i, b := range []byte(src) {
		if b == 'f' && i > 0 && src[i:i+4] == "func" {
			idx = i
			break
		}
	}
	line, col = ast.LineColWithStarts(starts, idx)
	if line != 3 || col != 1 {
		t.Errorf("func offset → (%d,%d), want (3,1)", line, col)
	}

	// exact line start via binary search Ok path
	line, col = ast.LineColWithStarts(starts, starts[2])
	if line != 3 || col != 1 {
		t.Errorf("line start 3 → (%d,%d)", line, col)
	}
}

func TestLineColEmptyStarts(t *testing.T) {
	line, col := ast.LineColWithStarts(nil, 10)
	if line != 1 || col != 1 {
		t.Errorf("got (%d,%d)", line, col)
	}
}

func TestLineColAt(t *testing.T) {
	starts := ast.ComputeLineStarts("a\nb")
	lc := ast.LineColAt(starts, 2) // 'b'
	if !lc.Valid() || lc.Line != 2 || lc.Column != 1 {
		t.Errorf("got %+v", lc)
	}
}

package tsparse_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay/codehound/internal/ast"
	"github.com/chinmay/codehound/internal/fixture"
	"github.com/chinmay/codehound/internal/lang/go/tsparse"
)

func TestParse_Simple(t *testing.T) {
	src := []byte("package main\n\nfunc main() {}\n")
	tree, err := tsparse.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	defer tree.Close()

	root := tree.RootNode()
	if root == nil {
		t.Fatal("nil root")
	}
	if root.Kind() != "source_file" {
		t.Errorf("kind = %q, want source_file", root.Kind())
	}

	var fnCount int
	ast.WalkNodes(root, []string{"function_declaration"}, func(n *tsparse.Node) {
		fnCount++
		name := n.ChildByFieldName("name")
		if name == nil || name.Utf8Text(src) != "main" {
			t.Errorf("expected function name main, got %v", name)
		}
	})
	if fnCount != 1 {
		t.Errorf("function_declaration count = %d", fnCount)
	}
}

func TestParse_MaterializedCWE22(t *testing.T) {
	txt := filepath.Join("..", "..", "..", "..", "tests", "fixtures", "go", "taint", "CWE-22-vulnerable.txt")
	data, err := os.ReadFile(txt)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	fx, err := fixture.ParseFixture(string(data), txt)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte(fx.Source)
	tree, err := tsparse.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	defer tree.Close()

	root := tree.RootNode()
	if root.HasError() {
		t.Logf("tree has error nodes (still usable): %s", root.ToSexp())
	}

	var calls []string
	ast.WalkCalls(root, func(n *tsparse.Node) {
		calls = append(calls, n.Utf8Text(src))
	})
	joined := strings.Join(calls, "\n")
	if !strings.Contains(joined, "os.Open") && !strings.Contains(joined, "Open") {
		// call text is the whole call_expression; should include os.Open(...)
		t.Errorf("expected Open call in %v", calls)
	}
	if !strings.Contains(joined, "filepath.Clean") && !strings.Contains(joined, "Clean") {
		t.Errorf("expected Clean call in %v", calls)
	}

	// Line/col for first function
	var fnLine int
	ast.WalkNodes(root, []string{"function_declaration"}, func(n *tsparse.Node) {
		if fnLine == 0 {
			line, col := tree.LineCol(int(n.StartByte()))
			fnLine = line
			if line < 1 || col < 1 {
				t.Errorf("invalid line/col (%d,%d)", line, col)
			}
		}
	})
	if fnLine == 0 {
		t.Error("no function_declaration found")
	}
}

func TestParser_Reuse(t *testing.T) {
	p, err := tsparse.NewParser()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	for _, src := range []string{
		"package a\n",
		"package b\nfunc F() {}\n",
	} {
		tree, err := p.Parse([]byte(src))
		if err != nil {
			t.Fatal(err)
		}
		if tree.RootNode().Kind() != "source_file" {
			t.Errorf("kind = %s", tree.RootNode().Kind())
		}
		tree.Close()
	}
}

func TestLineColMatchesAst(t *testing.T) {
	src := []byte("package p\n\nconst X = 1\n")
	tree, err := tsparse.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	defer tree.Close()

	// offset of 'c' in const
	off := strings.Index(string(src), "const")
	line, col := tree.LineCol(off)
	wantLine, wantCol := ast.LineColWithStarts(tree.LineStarts, off)
	if line != wantLine || col != wantCol {
		t.Errorf("LineCol = (%d,%d), ast = (%d,%d)", line, col, wantLine, wantCol)
	}
	if line != 3 || col != 1 {
		t.Errorf("expected (3,1), got (%d,%d)", line, col)
	}
}

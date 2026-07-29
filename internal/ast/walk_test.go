package ast_test

import (
	"testing"

	"github.com/chinmay/codehound/internal/ast"
	"github.com/chinmay/codehound/internal/lang/go/tsparse"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

func TestWalkCalls(t *testing.T) {
	src := []byte(`package sample

import "fmt"

func Hello() {
	fmt.Println("hi")
	fmt.Printf("%d", 1)
}
`)
	tree, err := tsparse.Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	defer tree.Close()

	var calls []string
	ast.WalkCalls(tree.RootNode(), func(n *sitter.Node) {
		calls = append(calls, n.Utf8Text(src))
	})
	if len(calls) < 2 {
		t.Fatalf("expected >=2 call_expression, got %d: %v", len(calls), calls)
	}
}

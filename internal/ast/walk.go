package ast

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// Walk visits every node under root (pre-order, iterative cursor).
// The visitor may return false to skip descending into the current node's children.
func Walk(root *sitter.Node, visit func(n *sitter.Node) bool) {
	if root == nil {
		return
	}
	cursor := root.Walk()
	defer cursor.Close()

	for {
		n := cursor.Node()
		descend := true
		if n != nil {
			descend = visit(n)
		}
		if descend && cursor.GotoFirstChild() {
			continue
		}
		for !cursor.GotoNextSibling() {
			if !cursor.GotoParent() {
				return
			}
		}
	}
}

// WalkNodes visits every node whose kind is in kinds.
func WalkNodes(root *sitter.Node, kinds []string, f func(n *sitter.Node)) {
	if root == nil || len(kinds) == 0 {
		return
	}
	set := make(map[string]struct{}, len(kinds))
	for _, k := range kinds {
		set[k] = struct{}{}
	}
	Walk(root, func(n *sitter.Node) bool {
		if _, ok := set[n.Kind()]; ok {
			f(n)
		}
		return true
	})
}

// WalkCalls visits call nodes (Go: call_expression; other langs: call).
func WalkCalls(root *sitter.Node, f func(n *sitter.Node)) {
	WalkNodes(root, []string{"call_expression", "call"}, f)
}

// WalkKinds visits nodes matching any of the given kinds (map form for hot loops).
func WalkKinds(root *sitter.Node, kinds map[string]struct{}, f func(n *sitter.Node)) {
	if root == nil || len(kinds) == 0 {
		return
	}
	Walk(root, func(n *sitter.Node) bool {
		if _, ok := kinds[n.Kind()]; ok {
			f(n)
		}
		return true
	})
}

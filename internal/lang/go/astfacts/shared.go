// Package astfacts builds a once-per-file AST fact bag shared by Go detectors
// (PERF, BP, …) so each pack does not re-Inspect the same *ast.File.
package astfacts

import (
	goast "go/ast"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/go/goparse"
)

// Call is one CallExpr site.
type Call struct {
	Callee string
	Args   []string
	Start  int
	End    int
	Text   string
}

// Assign is one AssignStmt.
type Assign struct {
	Start int
	End   int
	Text  string
}

// FuncDecl is one function declaration (not FuncLit).
type FuncDecl struct {
	Name   string
	Start  int
	End    int
	BodyS  int
	BodyE  int
	Params string
	IsMain bool
}

// ValueSpec is a var/const name with optional type text (PERF VarKinds).
type ValueSpec struct {
	Names    []string
	TypeText string
	Start    int
}

// Span is a byte range with optional source text (defer/go bodies).
type Span struct {
	Start int
	End   int
	Text  string
}

// Shared is the fused AST extract for one Go file.
type Shared struct {
	ForRanges      [][2]int
	DeferSpans     []Span
	GoSpans        []Span
	FuncLitRanges  [][2]int
	TypeAssertions [][2]int
	Calls          []Call
	Assigns        []Assign
	FuncDecls      []FuncDecl
	ValueSpecs     []ValueSpec
}

// Ensure returns a Shared bag for unit, building it once via unit.FactCache.
func Ensure(unit *core.ParsedUnit) *Shared {
	if unit == nil {
		return nil
	}
	if s, ok := unit.FactCache.(*Shared); ok && s != nil {
		return s
	}
	s := build(unit)
	unit.FactCache = s
	return s
}

func build(unit *core.ParsedUnit) *Shared {
	s := &Shared{}
	if unit == nil || unit.Source == "" {
		return s
	}
	tree := goparse.TreeForUnit(unit)
	if tree == nil || tree.File == nil {
		return s
	}

	// Pass 1: spans + decls + value specs (needed before loop enclosure).
	goast.Inspect(tree.File, func(n goast.Node) bool {
		if n == nil {
			return false
		}
		switch x := n.(type) {
		case *goast.ForStmt, *goast.RangeStmt:
			s.ForRanges = append(s.ForRanges, [2]int{tree.Offset(x.Pos()), tree.Offset(x.End())})
		case *goast.FuncLit:
			s.FuncLitRanges = append(s.FuncLitRanges, [2]int{tree.Offset(x.Pos()), tree.Offset(x.End())})
		case *goast.DeferStmt:
			s.DeferSpans = append(s.DeferSpans, Span{
				Start: tree.Offset(x.Pos()), End: tree.Offset(x.End()), Text: tree.NodeText(x),
			})
		case *goast.GoStmt:
			s.GoSpans = append(s.GoSpans, Span{
				Start: tree.Offset(x.Pos()), End: tree.Offset(x.End()), Text: tree.NodeText(x),
			})
		case *goast.TypeAssertExpr:
			s.TypeAssertions = append(s.TypeAssertions, [2]int{tree.Offset(x.Pos()), tree.Offset(x.End())})
		case *goast.FuncDecl:
			name := ""
			if x.Name != nil {
				name = x.Name.Name
			}
			start, end := tree.Offset(x.Pos()), tree.Offset(x.End())
			bodyS, bodyE := start, end
			if x.Body != nil {
				bodyS = tree.Offset(x.Body.Pos())
				bodyE = tree.Offset(x.Body.End())
			}
			params := ""
			if x.Type != nil && x.Type.Params != nil {
				params = tree.NodeText(x.Type.Params)
			}
			s.FuncDecls = append(s.FuncDecls, FuncDecl{
				Name: name, Start: start, End: end, BodyS: bodyS, BodyE: bodyE,
				Params: params, IsMain: name == "main",
			})
		case *goast.ValueSpec:
			var names []string
			for _, id := range x.Names {
				if id != nil && id.Name != "" && id.Name != "_" {
					names = append(names, id.Name)
				}
			}
			typeText := ""
			if x.Type != nil {
				typeText = tree.NodeText(x.Type)
			}
			if len(names) > 0 {
				s.ValueSpecs = append(s.ValueSpecs, ValueSpec{
					Names: names, TypeText: typeText, Start: tree.Offset(x.Pos()),
				})
			}
		}
		return true
	})

	// Pass 2: calls + assignments.
	goast.Inspect(tree.File, func(n goast.Node) bool {
		if n == nil {
			return false
		}
		switch x := n.(type) {
		case *goast.CallExpr:
			start, end := tree.Offset(x.Pos()), tree.Offset(x.End())
			callee := strings.TrimSpace(tree.NodeText(x.Fun))
			args := make([]string, 0, len(x.Args))
			for _, a := range x.Args {
				args = append(args, strings.TrimSpace(tree.NodeText(a)))
			}
			s.Calls = append(s.Calls, Call{
				Callee: callee, Args: args, Start: start, End: end,
				Text: tree.NodeText(x),
			})
		case *goast.AssignStmt:
			s.Assigns = append(s.Assigns, Assign{
				Start: tree.Offset(x.Pos()),
				End:   tree.Offset(x.End()),
				Text:  tree.NodeText(x),
			})
		}
		return true
	})
	return s
}

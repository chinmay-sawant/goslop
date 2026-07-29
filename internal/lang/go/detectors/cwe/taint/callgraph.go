package taint

import (
	goast "go/ast"
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/lang/go/goparse"
)

// ExtractCallGraph builds a per-file call graph from the unit AST.
func ExtractCallGraph(unit *core.ParsedUnit) *CallGraph {
	cg := NewCallGraph()
	if unit == nil {
		return cg
	}
	tree := unitTree(unit)
	if tree == nil || tree.File == nil {
		return cg
	}
	walkCallGraph(tree.File, tree, cg, nil)
	return cg
}

func walkCallGraph(n goast.Node, tree *goparse.Tree, cg *CallGraph, parents []goast.Node) {
	if n == nil {
		return
	}
	switch x := n.(type) {
	case *goast.FuncDecl:
		if x.Name != nil {
			name := x.Name.Name
			pc := 0
			if x.Type != nil && x.Type.Params != nil {
				pc = len(x.Type.Params.List)
			}
			if x.Recv == nil || len(x.Recv.List) == 0 {
				cg.AddDeclaration(name, FunctionDecl{Name: name, ParamCount: pc, IsMethod: false})
			} else {
				recv := tree.NodeText(x.Recv)
				identity := name
				if nrm := NormalizeReceiverType(recv); nrm != "" {
					identity = nrm + "." + name
				}
				cg.AddDeclaration(identity, FunctionDecl{
					Name: name, ParamCount: pc, IsMethod: true, ReceiverType: recv,
				})
			}
		}
	case *goast.CallExpr:
		recordCallSite(x, tree, cg, parents)
	}
	// Full-slice append so child walks cannot clobber sibling parent stacks.
	next := append(parents[:len(parents):len(parents)], n)
	forEachChild(n, func(child goast.Node) {
		walkCallGraph(child, tree, cg, next)
	})
}

func recordCallSite(node *goast.CallExpr, tree *goparse.Tree, cg *CallGraph, parents []goast.Node) {
	if node == nil || node.Fun == nil {
		return
	}
	callee := strings.TrimSpace(tree.NodeText(node.Fun))
	if callee == "" {
		return
	}
	caller := enclosingFunctionName(parents, tree)
	_, isMethod := node.Fun.(*goast.SelectorExpr)
	args := argumentTexts(node, tree)
	lhs := resultVariableOfCall(parents, tree)
	returns := callResultIsReturned(parents, tree, lhs)
	cg.AddSite(CallSite{
		Caller:        caller,
		Callee:        callee,
		ByteRange:     nodeRange(tree, node),
		Arguments:     args,
		AssignmentLHS: lhs,
		ReturnsResult: returns,
		IsMethodCall:  isMethod,
		IsClosure:     strings.HasPrefix(callee, "func(") || strings.HasPrefix(callee, "func "),
	})
}

func enclosingFunctionName(parents []goast.Node, tree *goparse.Tree) string {
	for i := len(parents) - 1; i >= 0; i-- {
		switch p := parents[i].(type) {
		case *goast.FuncDecl:
			if p.Name == nil {
				return "<anonymous>"
			}
			name := p.Name.Name
			if p.Recv == nil || len(p.Recv.List) == 0 {
				return name
			}
			recv := NormalizeReceiverType(tree.NodeText(p.Recv))
			if recv != "" {
				return recv + "." + name
			}
			return name
		case *goast.FuncLit:
			return "<anonymous>"
		}
	}
	return "<package>"
}

func callResultIsReturned(parents []goast.Node, tree *goparse.Tree, assignmentLHS string) bool {
	for i := len(parents) - 1; i >= 0; i-- {
		switch p := parents[i].(type) {
		case *goast.ReturnStmt:
			return true
		case *goast.FuncDecl:
			if assignmentLHS == "" {
				return false
			}
			body := ""
			if p.Body != nil {
				body = tree.NodeText(p.Body)
			}
			return assignmentReturnedInBody(body, assignmentLHS)
		case *goast.FuncLit:
			if assignmentLHS == "" {
				return false
			}
			body := ""
			if p.Body != nil {
				body = tree.NodeText(p.Body)
			}
			return assignmentReturnedInBody(body, assignmentLHS)
		}
	}
	return false
}

func assignmentReturnedInBody(body, assignmentLHS string) bool {
	for _, name := range strings.Split(assignmentLHS, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		for _, line := range strings.Split(body, "\n") {
			rest, ok := strings.CutPrefix(strings.TrimSpace(line), "return")
			if !ok {
				continue
			}
			rest = strings.TrimSpace(rest)
			// first identifier
			end := 0
			for end < len(rest) {
				c := rest[end]
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
					break
				}
				end++
			}
			if rest[:end] == name {
				return true
			}
		}
	}
	return false
}

// BuildImportMap maps import aliases to paths from source text.
func BuildImportMap(source string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(source, "\n") {
		line = strings.TrimSpace(line)
		// import "path" or import alias "path"
		if !strings.Contains(line, `"`) && !strings.Contains(line, "`") {
			continue
		}
		for _, quote := range []string{`"`, "`"} {
			// find quoted path
			i := strings.Index(line, quote)
			if i < 0 {
				continue
			}
			j := strings.LastIndex(line, quote)
			if j <= i {
				continue
			}
			path := line[i+1 : j]
			if !strings.Contains(path, "/") && path != "fmt" && path != "os" && path != "io" &&
				path != "net" && path != "flag" && path != "html" && path != "log" && path != "sync" {
				// still accept short paths
			}
			prefix := strings.TrimSpace(line[:i])
			// strip "import"
			prefix = strings.TrimPrefix(prefix, "import")
			prefix = strings.TrimSpace(prefix)
			alias := ""
			if prefix == "" {
				// default alias = last path segment
				if k := strings.LastIndex(path, "/"); k >= 0 {
					alias = path[k+1:]
				} else {
					alias = path
				}
			} else {
				fields := strings.Fields(prefix)
				if len(fields) > 0 {
					alias = fields[len(fields)-1]
				}
			}
			if alias != "" && alias != "." && alias != "_" {
				out[alias] = path
			}
		}
	}
	return out
}

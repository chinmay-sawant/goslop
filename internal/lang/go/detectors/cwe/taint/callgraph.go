package taint

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// ExtractCallGraph builds a per-file call graph from the unit AST.
func ExtractCallGraph(unit *core.ParsedUnit) *CallGraph {
	cg := NewCallGraph()
	if unit == nil {
		return cg
	}
	root := unitRoot(unit)
	if root == nil {
		return cg
	}
	src := []byte(unit.Source)
	walkCallGraph(root, src, cg)
	return cg
}

func walkCallGraph(node *sitter.Node, src []byte, cg *CallGraph) {
	if node == nil {
		return
	}
	switch node.Kind() {
	case "function_declaration":
		if nameN := node.ChildByFieldName("name"); nameN != nil {
			name := strings.TrimSpace(nameN.Utf8Text(src))
			pc := 0
			if p := node.ChildByFieldName("parameters"); p != nil {
				pc = int(p.NamedChildCount())
			}
			cg.AddDeclaration(name, FunctionDecl{Name: name, ParamCount: pc, IsMethod: false})
		}
	case "method_declaration":
		if nameN := node.ChildByFieldName("name"); nameN != nil {
			name := strings.TrimSpace(nameN.Utf8Text(src))
			pc := 0
			if p := node.ChildByFieldName("parameters"); p != nil {
				pc = int(p.NamedChildCount())
			}
			recv := ""
			if r := node.ChildByFieldName("receiver"); r != nil {
				recv = r.Utf8Text(src)
			}
			identity := name
			if n := NormalizeReceiverType(recv); n != "" {
				identity = n + "." + name
			}
			cg.AddDeclaration(identity, FunctionDecl{
				Name: name, ParamCount: pc, IsMethod: true, ReceiverType: recv,
			})
		}
	case "call_expression":
		recordCallSite(node, src, cg)
	}
	for i := uint(0); i < node.ChildCount(); i++ {
		walkCallGraph(node.Child(i), src, cg)
	}
}

func recordCallSite(node *sitter.Node, src []byte, cg *CallGraph) {
	fn := node.ChildByFieldName("function")
	if fn == nil {
		return
	}
	callee := strings.TrimSpace(fn.Utf8Text(src))
	if callee == "" {
		return
	}
	caller := enclosingFunctionName(node, src)
	isMethod := fn.Kind() == "selector_expression"
	args := argumentTexts(node, src)
	lhs := resultVariableOfCall(node, src)
	returns := callResultIsReturned(node, src, lhs)
	cg.AddSite(CallSite{
		Caller:        caller,
		Callee:        callee,
		ByteRange:     ByteRange{int(node.StartByte()), int(node.EndByte())},
		Arguments:     args,
		AssignmentLHS: lhs,
		ReturnsResult: returns,
		IsMethodCall:  isMethod,
		IsClosure:     strings.HasPrefix(callee, "func(") || strings.HasPrefix(callee, "func "),
	})
}

func enclosingFunctionName(node *sitter.Node, src []byte) string {
	for p := node.Parent(); p != nil; p = p.Parent() {
		switch p.Kind() {
		case "function_declaration":
			if n := p.ChildByFieldName("name"); n != nil {
				return strings.TrimSpace(n.Utf8Text(src))
			}
		case "method_declaration":
			name := "<anonymous>"
			if n := p.ChildByFieldName("name"); n != nil {
				name = strings.TrimSpace(n.Utf8Text(src))
			}
			recv := ""
			if r := p.ChildByFieldName("receiver"); r != nil {
				recv = NormalizeReceiverType(r.Utf8Text(src))
			}
			if recv != "" {
				return recv + "." + name
			}
			return name
		case "func_literal":
			return "<anonymous>"
		}
	}
	return "<package>"
}

func callResultIsReturned(node *sitter.Node, src []byte, assignmentLHS string) bool {
	for p := node.Parent(); p != nil; p = p.Parent() {
		switch p.Kind() {
		case "return_statement":
			return true
		case "function_declaration", "method_declaration", "func_literal":
			if assignmentLHS == "" {
				return false
			}
			body := p.Utf8Text(src)
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

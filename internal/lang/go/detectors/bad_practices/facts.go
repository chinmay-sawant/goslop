package badpractices

import (
	"strings"

	"github.com/chinmay/codehound/internal/ast"
	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/tsparse"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// Shared needle table for BP fast-paths (Rust parity intent).
var bpNeedles = []string{
	"_", ":=", "recover()", ".Close(", "defer ", "go ", "go func",
	"panic(", "time.After", "time.Sleep", "context.Background", "context.TODO",
	"sync.Mutex", "sync.WaitGroup", ".Add(", "http.Server", "log.Fatal", "os.Exit",
	"json.Unmarshal", "errors.As", "os.Open", "os.Create", "http.Get", "http.Do",
	".Query(", ".QueryContext(", "ShouldBind", "BodyParser", "Run:", "RunE:",
	"signal.Notify", ".Shutdown(", "ListenAndServe", "rate.NewLimiter",
}

// bpFacts is the fused fact bag for BP detectors.
type bpFacts struct {
	Source      string
	Index       ast.SourceIndex
	tree        *tsparse.Tree
	ownedTree   bool
	// AST-derived ranges
	assignNodes []nodeSpan // assignment_statement / short_var_declaration
	callNodes   []callSpan
	deferNodes  []nodeSpan
	goNodes     []nodeSpan
	forRanges   [][2]int
	funcDecls   []funcSpan
}

type nodeSpan struct {
	start int
	end   int
	text  string
}

type callSpan struct {
	start  int
	end    int
	callee string
	text   string
	parent string // parent kind when known
}

type funcSpan struct {
	name  string
	start int
	end   int
	bodyS int
	bodyE int
	isMain bool
	params string
}

func buildFacts(unit *core.ParsedUnit) *bpFacts {
	f := &bpFacts{}
	if unit == nil {
		return f
	}
	f.Source = unit.Source
	f.Index = ast.Build(unit.Source, bpNeedles)

	src := []byte(unit.Source)
	var tree *tsparse.Tree
	if t, ok := unit.Tree.(*tsparse.Tree); ok && t != nil {
		tree = t
	} else if unit.Source != "" {
		t, err := tsparse.Parse(src)
		if err == nil && t != nil {
			tree = t
			f.ownedTree = true
		}
	}
	if tree == nil {
		return f
	}
	f.tree = tree
	root := tree.RootNode()
	if root == nil {
		return f
	}

	kinds := map[string]struct{}{
		"assignment_statement":  {},
		"short_var_declaration": {},
		"call_expression":       {},
		"defer_statement":       {},
		"go_statement":          {},
		"for_statement":         {},
		"function_declaration":  {},
		"method_declaration":    {},
	}
	ast.WalkKinds(root, kinds, func(n *sitter.Node) {
		start := int(n.StartByte())
		end := int(n.EndByte())
		if end > len(src) {
			end = len(src)
		}
		if start < 0 {
			start = 0
		}
		text := string(src[start:end])
		switch n.Kind() {
		case "assignment_statement", "short_var_declaration":
			f.assignNodes = append(f.assignNodes, nodeSpan{start: start, end: end, text: text})
		case "call_expression":
			callee := ""
			if fn := n.ChildByFieldName("function"); fn != nil {
				fs, fe := int(fn.StartByte()), int(fn.EndByte())
				if fe <= len(src) && fs >= 0 {
					callee = string(src[fs:fe])
				}
			}
			parentKind := ""
			if p := n.Parent(); p != nil {
				parentKind = p.Kind()
			}
			f.callNodes = append(f.callNodes, callSpan{
				start: start, end: end, callee: callee, text: text, parent: parentKind,
			})
		case "defer_statement":
			f.deferNodes = append(f.deferNodes, nodeSpan{start: start, end: end, text: text})
		case "go_statement":
			f.goNodes = append(f.goNodes, nodeSpan{start: start, end: end, text: text})
		case "for_statement":
			f.forRanges = append(f.forRanges, [2]int{start, end})
		case "function_declaration", "method_declaration":
			name := ""
			if nm := n.ChildByFieldName("name"); nm != nil {
				ns, ne := int(nm.StartByte()), int(nm.EndByte())
				if ne <= len(src) && ns >= 0 {
					name = string(src[ns:ne])
				}
			}
			params := ""
			if pl := n.ChildByFieldName("parameters"); pl != nil {
				ps, pe := int(pl.StartByte()), int(pl.EndByte())
				if pe <= len(src) && ps >= 0 {
					params = string(src[ps:pe])
				}
			}
			bodyS, bodyE := start, end
			if body := n.ChildByFieldName("body"); body != nil {
				bodyS, bodyE = int(body.StartByte()), int(body.EndByte())
			}
			f.funcDecls = append(f.funcDecls, funcSpan{
				name: name, start: start, end: end, bodyS: bodyS, bodyE: bodyE,
				isMain: name == "main", params: params,
			})
		}
	})
	return f
}

func (f *bpFacts) close() {
	if f != nil && f.ownedTree && f.tree != nil {
		f.tree.Close()
		f.tree = nil
	}
}

func (f *bpFacts) insideLoop(byteOffset int) bool {
	for _, r := range f.forRanges {
		if byteOffset >= r[0] && byteOffset < r[1] {
			return true
		}
	}
	return false
}

func (f *bpFacts) enclosingFunc(byteOffset int) *funcSpan {
	var best *funcSpan
	for i := range f.funcDecls {
		fn := &f.funcDecls[i]
		if byteOffset >= fn.start && byteOffset < fn.end {
			if best == nil || (fn.end-fn.start) < (best.end-best.start) {
				best = fn
			}
		}
	}
	return best
}

func (f *bpFacts) has(needle string) bool {
	if f == nil {
		return false
	}
	if f.Index.Has(needle) {
		return true
	}
	return strings.Contains(f.Source, needle)
}

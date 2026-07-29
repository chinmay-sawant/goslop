package badpractices

import (
	goast "go/ast"
	"strings"

	"github.com/chinmay/goslop/internal/ast"
	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/lang/go/goparse"
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
	tree        *goparse.Tree
	ownedTree   bool
	assignNodes []nodeSpan
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
	parent string
}

type funcSpan struct {
	name   string
	start  int
	end    int
	bodyS  int
	bodyE  int
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
	var tree *goparse.Tree
	if t, ok := unit.Tree.(*goparse.Tree); ok && t != nil && t.File != nil {
		tree = t
	} else if unit.Source != "" {
		t, err := goparse.Parse(src)
		if err == nil && t != nil && t.File != nil {
			tree = t
			f.ownedTree = true
		} else if t != nil && t.File != nil {
			tree = t
			f.ownedTree = true
		}
	}
	if tree == nil || tree.File == nil {
		return f
	}
	f.tree = tree

	goast.Inspect(tree.File, func(n goast.Node) bool {
		if n == nil {
			return true
		}
		start := tree.Offset(n.Pos())
		end := tree.Offset(n.End())
		text := tree.NodeText(n)
		switch x := n.(type) {
		case *goast.AssignStmt:
			f.assignNodes = append(f.assignNodes, nodeSpan{start: start, end: end, text: text})
		case *goast.CallExpr:
			callee := strings.TrimSpace(tree.NodeText(x.Fun))
			f.callNodes = append(f.callNodes, callSpan{
				start: start, end: end, callee: callee, text: text,
			})
		case *goast.DeferStmt:
			f.deferNodes = append(f.deferNodes, nodeSpan{start: start, end: end, text: text})
		case *goast.GoStmt:
			f.goNodes = append(f.goNodes, nodeSpan{start: start, end: end, text: text})
		case *goast.ForStmt, *goast.RangeStmt:
			f.forRanges = append(f.forRanges, [2]int{start, end})
		case *goast.FuncDecl:
			name := ""
			if x.Name != nil {
				name = x.Name.Name
			}
			bodyS, bodyE := start, end
			if x.Body != nil {
				bodyS = tree.Offset(x.Body.Pos())
				bodyE = tree.Offset(x.Body.End())
			}
			params := ""
			if x.Type != nil && x.Type.Params != nil {
				params = tree.NodeText(x.Type.Params)
			}
			f.funcDecls = append(f.funcDecls, funcSpan{
				name:   name,
				start:  start,
				end:    end,
				bodyS:  bodyS,
				bodyE:  bodyE,
				isMain: name == "main",
				params: params,
			})
		}
		return true
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

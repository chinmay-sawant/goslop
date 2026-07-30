package badpractices

import (
	"strings"

	"github.com/chinmay/goslop/internal/ast"
	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/lang/go/astfacts"
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

// codeLine is one source line with 0-based index and byte offset (P1.1).
type codeLine struct {
	idx  int
	text string
	byte int
}

// bpFacts is the fused fact bag for BP detectors.
type bpFacts struct {
	Source string
	Index  ast.SourceIndex
	tree   *goparse.Tree
	// projectFacts is owned by the current detector scan. Rules must obtain
	// cross-file facts through it rather than a process-global cache.
	projectFacts *bpProjectCaches
	// lines is built once per file (codeLines / stripLineComment cost).
	lines       []codeLine
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

func buildFacts(unit *core.ParsedUnit, projectFacts *bpProjectCaches) *bpFacts {
	f := &bpFacts{projectFacts: projectFacts}
	if unit == nil {
		return f
	}
	f.Source = unit.Source
	// Once per file: shared by all BP rules that walk source lines.
	f.lines = buildCodeLines(unit.Source)
	f.Index = ast.Build(unit.Source, bpNeedles)

	// Shared AST walk with PERF (one Inspect pair per file via unit.FactCache).
	tree := goparse.TreeForUnit(unit)
	if tree != nil {
		f.tree = tree
	}
	shared := astfacts.Ensure(unit)
	if shared == nil {
		return f
	}
	f.forRanges = shared.ForRanges
	for _, a := range shared.Assigns {
		f.assignNodes = append(f.assignNodes, nodeSpan{start: a.Start, end: a.End, text: a.Text})
	}
	for _, c := range shared.Calls {
		f.callNodes = append(f.callNodes, callSpan{
			start: c.Start, end: c.End, callee: c.Callee, text: c.Text,
		})
	}
	for _, r := range shared.DeferSpans {
		f.deferNodes = append(f.deferNodes, nodeSpan{start: r.Start, end: r.End, text: r.Text})
	}
	for _, r := range shared.GoSpans {
		f.goNodes = append(f.goNodes, nodeSpan{start: r.Start, end: r.End, text: r.Text})
	}
	for _, fd := range shared.FuncDecls {
		f.funcDecls = append(f.funcDecls, funcSpan{
			name: fd.Name, start: fd.Start, end: fd.End,
			bodyS: fd.BodyS, bodyE: fd.BodyE, isMain: fd.IsMain, params: fd.Params,
		})
	}
	return f
}

func (f *bpFacts) projectSnapshot(unit *core.ParsedUnit) *ProjectSnapshot {
	return projectSnapshot(unit, f.projectFacts)
}

func (f *bpFacts) packageDocSnapshot(unit *core.ParsedUnit) *PackageDocSnapshot {
	return packageDocSnapshotForUnit(unit, f.projectFacts)
}

func (f *bpFacts) packageTypeFacts(unit *core.ParsedUnit) *packageTypeFacts {
	return packageTypeFactsForUnit(unit, f.projectFacts)
}

// close releases any BP-owned resources. The shared unit AST is not closed here
// (analyzer closeUnitTree owns unit.Tree lifetime for cross-pack reuse).
func (f *bpFacts) close() {
	if f != nil {
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

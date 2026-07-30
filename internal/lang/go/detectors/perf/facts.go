package perf

import (
	goast "go/ast"
	"go/token"
	"strings"
	"unicode"

	"github.com/chinmay/goslop/internal/ast"
	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/lang/go/astfacts"
	"github.com/chinmay/goslop/internal/lang/go/goparse"
)

// CallFact is one call site extracted for PERF rules.
type CallFact struct {
	Callee        string
	Arguments     []string
	StartByte     int
	EnclosingLoop *int
}

// AssignmentFact is one assignment / short-var declaration.
type AssignmentFact struct {
	Name          string
	Expr          string
	Text          string
	StartByte     int
	EnclosingLoop *int
}

// ConversionFact is a type conversion expression.
type ConversionFact struct {
	Text      string
	StartByte int
	InLoop    bool
}

// VarKind is a coarse declaration-site type for suppressing FPs.
type VarKind int

const (
	VarUnknown VarKind = iota
	VarNumeric
	VarString
	VarBytes
)

// GoPerfFacts is the fused fact bag built once per file for PERF detectors.
type GoPerfFacts struct {
	Calls                 []CallFact
	Assignments           []AssignmentFact
	Conversions           []ConversionFact
	VarKinds              map[string]VarKind
	DeferStarts           [][2]int
	GoStarts              [][2]int
	ForRanges             [][2]int
	FunctionLiteralRanges [][2]int
	TypeAssertions        [][2]int
	SourceIndex           ast.SourceIndex
}

// BuildFacts builds PERF facts from the shared AST walk (astfacts) plus
// pack-specific conversion / var-kind post-processing.
func BuildFacts(unit *core.ParsedUnit) *GoPerfFacts {
	facts := &GoPerfFacts{VarKinds: map[string]VarKind{}}
	if unit == nil || unit.Source == "" {
		return facts
	}

	// Index first so Run can gate rules; walk uses shared bag with BP.
	facts.SourceIndex = ast.Build(unit.Source, perfNeedles)

	shared := astfacts.Ensure(unit)
	if shared == nil {
		return facts
	}

	facts.ForRanges = shared.ForRanges
	facts.FunctionLiteralRanges = shared.FuncLitRanges
	facts.TypeAssertions = shared.TypeAssertions
	for _, r := range shared.DeferSpans {
		facts.DeferStarts = append(facts.DeferStarts, [2]int{r.Start, r.End})
	}
	for _, r := range shared.GoSpans {
		facts.GoStarts = append(facts.GoStarts, [2]int{r.Start, r.End})
	}

	for _, vs := range shared.ValueSpecs {
		kind := classifyTypeText(vs.TypeText)
		for _, name := range vs.Names {
			facts.VarKinds[name] = kind
		}
	}

	for _, c := range shared.Calls {
		loop := enclosingLoop(facts.ForRanges, c.Start)
		facts.Calls = append(facts.Calls, CallFact{
			Callee:        c.Callee,
			Arguments:     append([]string(nil), c.Args...),
			StartByte:     c.Start,
			EnclosingLoop: loop,
		})
		text := strings.TrimSpace(c.Text)
		if isStringBytesConversion(text) {
			facts.Conversions = append(facts.Conversions, ConversionFact{
				Text:      text,
				StartByte: c.Start,
				InLoop:    loop != nil,
			})
		}
	}

	for _, a := range shared.Assigns {
		loop := enclosingLoop(facts.ForRanges, a.Start)
		recordAssignFromText(a.Text, a.Start, facts, loop)
	}

	return facts
}

func classifyTypeText(typeText string) VarKind {
	t := strings.TrimSpace(typeText)
	switch {
	case strings.Contains(t, "[]byte") || strings.Contains(t, "[]uint8"):
		return VarBytes
	case t == "string":
		return VarString
	case t == "int" || t == "int64" || t == "int32" || t == "uint" || t == "uint64" ||
		t == "float64" || t == "float32" || t == "byte" || t == "rune" || t == "time.Duration":
		return VarNumeric
	default:
		return VarUnknown
	}
}

func recordAssignFromText(text string, start int, facts *GoPerfFacts, loop *int) {
	lhs, rhs, ok := splitAssignment(text)
	if !ok {
		return
	}
	isShort := strings.Contains(text, ":=")
	names := extractIdents(lhs)
	if len(names) == 0 {
		if n := firstAssignName(lhs); n != "" {
			names = []string{n}
		}
	}
	for _, name := range names {
		facts.Assignments = append(facts.Assignments, AssignmentFact{
			Name:          name,
			Expr:          rhs,
			Text:          text,
			StartByte:     start,
			EnclosingLoop: loop,
		})
		if isShort {
			if _, exists := facts.VarKinds[name]; !exists {
				if k := classifyInit(rhs); k != VarUnknown {
					facts.VarKinds[name] = k
				}
			}
		}
	}
}

func firstAssignName(lhs string) string {
	lhs = strings.TrimSpace(lhs)
	if i := strings.IndexAny(lhs, ", "); i >= 0 {
		lhs = strings.TrimSpace(lhs[:i])
	}
	return lhs
}

func enclosingLoop(ranges [][2]int, byteOff int) *int {
	// Prefer innermost: last range that contains byteOff.
	best := -1
	for i := range ranges {
		if byteOff >= ranges[i][0] && byteOff < ranges[i][1] {
			best = ranges[i][0]
		}
	}
	if best < 0 {
		return nil
	}
	s := best
	return &s
}

func recordCallAST(tree *goparse.Tree, ce *goast.CallExpr, facts *GoPerfFacts, loop *int) {
	if ce == nil {
		return
	}
	callee := strings.TrimSpace(tree.NodeText(ce.Fun))
	if callee == "" {
		return
	}
	args := make([]string, 0, len(ce.Args))
	for _, a := range ce.Args {
		args = append(args, strings.TrimSpace(tree.NodeText(a)))
	}
	facts.Calls = append(facts.Calls, CallFact{
		Callee:        callee,
		Arguments:     args,
		StartByte:     tree.Offset(ce.Pos()),
		EnclosingLoop: loop,
	})
}

func recordAssignAST(tree *goparse.Tree, as *goast.AssignStmt, facts *GoPerfFacts, loop *int) {
	if as == nil {
		return
	}
	text := tree.NodeText(as)
	lhs, rhs, ok := splitAssignment(text)
	if !ok && len(as.Lhs) > 0 {
		lhs = tree.NodeText(as.Lhs[0])
		if len(as.Rhs) > 0 {
			rhs = tree.NodeText(as.Rhs[0])
		}
		ok = true
	}
	if !ok {
		return
	}
	start := tree.Offset(as.Pos())
	isShort := as.Tok == token.DEFINE
	names := extractIdents(lhs)
	if len(names) == 0 {
		for _, l := range as.Lhs {
			names = append(names, extractIdents(tree.NodeText(l))...)
		}
	}
	for _, name := range names {
		facts.Assignments = append(facts.Assignments, AssignmentFact{
			Name:          name,
			Expr:          rhs,
			Text:          text,
			StartByte:     start,
			EnclosingLoop: loop,
		})
		if isShort {
			if _, exists := facts.VarKinds[name]; !exists {
				if k := classifyInit(rhs); k != VarUnknown {
					facts.VarKinds[name] = k
				}
			}
		}
	}
}

func isStringBytesConversion(text string) bool {
	t := strings.TrimSpace(text)
	if strings.HasPrefix(t, "[]byte(") || strings.HasPrefix(t, "[]uint8(") {
		return true
	}
	// string(x) but not string("literal")
	if strings.HasPrefix(t, "string(") && !strings.HasPrefix(t, "string(\"") && !strings.HasPrefix(t, "string(`") {
		return true
	}
	return false
}

func collectVarSpecAST(tree *goparse.Tree, vs *goast.ValueSpec, facts *GoPerfFacts) {
	if vs == nil {
		return
	}
	kind := VarUnknown
	if vs.Type != nil {
		typeText := tree.NodeText(vs.Type)
		switch {
		case strings.Contains(typeText, "[]byte") || strings.Contains(typeText, "[]uint8"):
			kind = VarBytes
		case typeText == "string":
			kind = VarString
		case typeText == "int" || typeText == "int64" || typeText == "float64" ||
			typeText == "float32" || typeText == "uint" || typeText == "uint64" ||
			typeText == "time.Duration":
			kind = VarNumeric
		}
	}
	if kind == VarUnknown {
		return
	}
	for _, name := range vs.Names {
		if name == nil {
			continue
		}
		if _, ok := facts.VarKinds[name.Name]; !ok {
			facts.VarKinds[name.Name] = kind
		}
	}
}

func splitAssignment(text string) (lhs, rhs string, ok bool) {
	text = strings.TrimSpace(text)
	if i := strings.Index(text, ":="); i > 0 {
		return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+2:]), true
	}
	for _, op := range []string{"+=", "-=", "*=", "/=", "%="} {
		if i := strings.Index(text, op); i > 0 {
			return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+len(op):]), true
		}
	}
	if i := strings.Index(text, "="); i > 0 {
		if i+1 < len(text) && text[i+1] == '=' {
			return "", "", false
		}
		return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+1:]), true
	}
	return "", "", false
}

func extractIdents(lhs string) []string {
	parts := strings.Split(lhs, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		name := lhsPrimaryName(p)
		if name == "" || name == "_" {
			continue
		}
		out = append(out, name)
	}
	return out
}

func lhsPrimaryName(lhs string) string {
	s := strings.TrimSpace(lhs)
	if s == "" {
		return ""
	}
	for {
		i := strings.Index(s, "[")
		if i < 0 {
			break
		}
		j := strings.Index(s[i:], "]")
		if j < 0 {
			break
		}
		s = s[:i] + s[i+j+1:]
	}
	s = strings.TrimSpace(s)
	if k := strings.LastIndex(s, "."); k >= 0 {
		s = s[k+1:]
	}
	if i := strings.IndexAny(s, " \t"); i >= 0 {
		s = s[:i]
	}
	if !isSimpleIdent(s) {
		return ""
	}
	return s
}

func isSimpleIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
			continue
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func classifyInit(rhs string) VarKind {
	r := strings.TrimSpace(rhs)
	if r == "" {
		return VarUnknown
	}
	if strings.HasPrefix(r, "[]byte") || strings.HasPrefix(r, "[]uint8") {
		return VarBytes
	}
	if (r[0] == '"' || r[0] == '`') && !strings.Contains(r, "+") {
		return VarString
	}
	if isNumericLiteral(r) {
		return VarNumeric
	}
	return VarUnknown
}

func isNumericLiteral(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	dot := 0
	for i, r := range s {
		if r == '.' {
			dot++
			if dot > 1 {
				return false
			}
			continue
		}
		if r == '_' {
			continue
		}
		if i == 0 && (r == '+' || r == '-') {
			continue
		}
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// perfNeedles backs SourceIndex + rule gates (zz_gates.go). Keep gate tokens here.
var perfNeedles = []string{
	"fmt.Sprintf(",
	"fmt.Fprintf(",
	"fmt.Errorf(",
	"for ",
	"defer ",
	"regexp.MustCompile",
	"regexp.Compile",
	"regexp.MatchString",
	"time.Parse",
	"json.Marshal",
	"json.Unmarshal",
	"json.NewEncoder",
	"json.NewDecoder",
	"[]byte(",
	"string(",
	"http.Client{",
	"&http.Client{",
	"bytes.Buffer{",
	"new(bytes.Buffer)",
	"reflect.ValueOf(",
	"io.ReadAll(",
	"os.ReadFile(",
	"ioutil.ReadFile(",
	"bytes.NewReader(",
	"bytes.NewBuffer(",
	"sync.Mutex",
	"sync.RWMutex",
	"strings.Trim",
	"TrimSpace",
	"TrimPrefix",
	"strings.Split",
	"strings.Fields",
	"strings.HasPrefix",
	"HasPrefix",
	"http.Server{",
	"&http.Server{",
	"ListenAndServe",
}

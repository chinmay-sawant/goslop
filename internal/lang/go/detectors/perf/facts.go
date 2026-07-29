package perf

import (
	"strings"
	"unicode"

	"github.com/chinmay/codehound/internal/ast"
	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/tsparse"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// CallFact is one call_expression site extracted for PERF rules.
type CallFact struct {
	Callee        string
	Arguments     []string
	StartByte     int
	EnclosingLoop *int // start byte of nearest for_statement, if any
}

// AssignmentFact is one assignment / short-var declaration.
type AssignmentFact struct {
	Name          string
	Expr          string
	Text          string
	StartByte     int
	EnclosingLoop *int
}

// ConversionFact is a type conversion expression ([]byte(s) / string(b)).
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

// BuildFacts walks the unit AST (parsing on demand when unit.Tree is unset).
func BuildFacts(unit *core.ParsedUnit) *GoPerfFacts {
	facts := &GoPerfFacts{
		VarKinds: map[string]VarKind{},
	}
	if unit == nil || unit.Source == "" {
		return facts
	}

	src := []byte(unit.Source)
	var tree *tsparse.Tree
	owned := false
	if t, ok := unit.Tree.(*tsparse.Tree); ok && t != nil {
		tree = t
	} else {
		t, err := tsparse.Parse(src)
		if err != nil || t == nil {
			return facts
		}
		tree = t
		owned = true
	}
	if owned {
		defer tree.Close()
	}

	root := tree.RootNode()
	if root == nil {
		return facts
	}

	kinds := map[string]struct{}{
		"call_expression":            {},
		"assignment_statement":       {},
		"short_var_declaration":      {},
		"defer_statement":            {},
		"go_statement":               {},
		"for_statement":              {},
		"func_literal":               {},
		"type_assertion_expression":  {},
		"type_conversion_expression": {},
		"conversion_expression":      {},
	}
	ast.WalkKinds(root, kinds, func(n *sitter.Node) {
		switch n.Kind() {
		case "call_expression":
			recordCall(n, facts, src)
		case "assignment_statement", "short_var_declaration":
			recordAssignment(n, facts, src)
		case "defer_statement":
			facts.DeferStarts = append(facts.DeferStarts, [2]int{int(n.StartByte()), int(n.EndByte())})
		case "go_statement":
			facts.GoStarts = append(facts.GoStarts, [2]int{int(n.StartByte()), int(n.EndByte())})
		case "for_statement":
			facts.ForRanges = append(facts.ForRanges, [2]int{int(n.StartByte()), int(n.EndByte())})
		case "func_literal":
			facts.FunctionLiteralRanges = append(facts.FunctionLiteralRanges, [2]int{int(n.StartByte()), int(n.EndByte())})
		case "type_assertion_expression":
			facts.TypeAssertions = append(facts.TypeAssertions, [2]int{int(n.StartByte()), int(n.EndByte())})
		case "type_conversion_expression", "conversion_expression":
			recordConversion(n, facts, src)
		}
	})

	// var_spec kinds only when shapes used by PERF-2/32 appear.
	needKinds := strings.Contains(unit.Source, "+=") ||
		strings.Contains(unit.Source, "[]byte(") ||
		strings.Contains(unit.Source, "[]uint8(") ||
		strings.Contains(unit.Source, "string(")
	if needKinds {
		ast.WalkNodes(root, []string{"var_spec"}, func(n *sitter.Node) {
			collectVarSpecKinds(n, facts, src)
		})
	}

	facts.SourceIndex = ast.Build(unit.Source, perfNeedles)
	return facts
}

func recordCall(n *sitter.Node, facts *GoPerfFacts, src []byte) {
	fn := n.ChildByFieldName("function")
	if fn == nil {
		return
	}
	callee := strings.TrimSpace(fn.Utf8Text(src))
	if callee == "" {
		return
	}
	var args []string
	if argNode := n.ChildByFieldName("arguments"); argNode != nil {
		cursor := argNode.Walk()
		for _, child := range argNode.NamedChildren(cursor) {
			args = append(args, strings.TrimSpace(child.Utf8Text(src)))
		}
		cursor.Close()
	}
	facts.Calls = append(facts.Calls, CallFact{
		Callee:        callee,
		Arguments:     args,
		StartByte:     int(n.StartByte()),
		EnclosingLoop: enclosingLoopStart(n),
	})
}

func recordAssignment(n *sitter.Node, facts *GoPerfFacts, src []byte) {
	text := n.Utf8Text(src)
	lhs, rhs, ok := splitAssignment(text)
	if !ok {
		return
	}
	isShort := strings.Contains(text, ":=")
	loop := enclosingLoopStart(n)
	for _, name := range extractIdents(lhs) {
		facts.Assignments = append(facts.Assignments, AssignmentFact{
			Name:          name,
			Expr:          rhs,
			Text:          text,
			StartByte:     int(n.StartByte()),
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

func recordConversion(n *sitter.Node, facts *GoPerfFacts, src []byte) {
	text := strings.TrimSpace(n.Utf8Text(src))
	facts.Conversions = append(facts.Conversions, ConversionFact{
		Text:      text,
		StartByte: int(n.StartByte()),
		InLoop:    enclosingLoopStart(n) != nil,
	})
}

func enclosingLoopStart(n *sitter.Node) *int {
	cur := n
	for cur != nil {
		p := cur.Parent()
		if p == nil {
			return nil
		}
		if p.Kind() == "for_statement" {
			s := int(p.StartByte())
			return &s
		}
		cur = p
	}
	return nil
}

func splitAssignment(text string) (lhs, rhs string, ok bool) {
	text = strings.TrimSpace(text)
	if i := strings.Index(text, ":="); i > 0 {
		return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+2:]), true
	}
	// Prefer compound ops before bare '='.
	for _, op := range []string{"+=", "-=", "*=", "/=", "%="} {
		if i := strings.Index(text, op); i > 0 {
			return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+len(op):]), true
		}
	}
	if i := strings.Index(text, "="); i > 0 {
		// skip ==
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

// lhsPrimaryName extracts a usable identifier from an assignment LHS that may
// be an index or selector expression (tables[k], font.UsedChars, lines[i].joined).
func lhsPrimaryName(lhs string) string {
	s := strings.TrimSpace(lhs)
	if s == "" {
		return ""
	}
	// Drop index expressions: lines[li].joined → lines.joined, tables[k] → tables
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
	// Prefer the rightmost selector field for field assigns.
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
	// numeric / duration-ish literals
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
	// strip trailing type suffixes like 0.0 or 1_000
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

func collectVarSpecKinds(n *sitter.Node, facts *GoPerfFacts, src []byte) {
	// var name T = expr  OR  var name = expr
	// tree-sitter-go var_spec may have multiple names; collect identifier children.
	cursor := n.Walk()
	var names []string
	for _, ch := range n.NamedChildren(cursor) {
		if ch.Kind() == "identifier" {
			names = append(names, ch.Utf8Text(src))
		}
	}
	cursor.Close()
	if len(names) == 0 {
		return
	}
	typeNode := n.ChildByFieldName("type")
	typeText := ""
	if typeNode != nil {
		typeText = typeNode.Utf8Text(src)
	}
	kind := VarUnknown
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
	if kind == VarUnknown {
		return
	}
	for _, name := range names {
		if _, ok := facts.VarKinds[name]; !ok {
			facts.VarKinds[name] = kind
		}
	}
}

// Core needles for whole-file guards (subset of Rust NEEDLES).
var perfNeedles = []string{
	"fmt.Sprintf(",
	"fmt.Fprintf(",
	"for ",
	"regexp.MustCompile",
	"regexp.Compile",
	"regexp.MatchString",
	"time.Parse",
	"json.Marshal",
	"json.Unmarshal",
	"[]byte(",
	"string(",
	"http.ResponseWriter",
	"*gin.Context",
	"echo.Context",
	"*fiber.Ctx",
	"gin.HandlerFunc",
	"http.HandlerFunc",
	"defer ",
	"strings.Index",
	"strings.Builder",
	"bytes.Buffer",
}

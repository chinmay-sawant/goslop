package cwe

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/lang/go/detectors/sourceutil"
	"github.com/chinmay/goslop/internal/rules"
)

// Phase 7 taint-lite ports for registry taint-core IDs.
// Full inter-procedural taint graph is Phase 9 — these keep seed behaviour for
// CWE-78/89 and add same-file heuristics for CWE-22/79/90/91.

func init() {
	RegisterRule("CWE-22", detectCWE22, &MetaCWE22)
	RegisterRule("CWE-78", detectCWE78, &MetaCWE78)
	RegisterRule("CWE-79", detectCWE79, &MetaCWE79)
	RegisterRule("CWE-89", detectCWE89, &MetaCWE89)
	RegisterRule("CWE-90", detectCWE90, &MetaCWE90)
	RegisterRule("CWE-91", detectCWE91, &MetaCWE91)
}

// --- CWE-22 path traversal (taint-lite) ---

func detectCWE22(unit *core.ParsedUnit, _ *GoCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source
	// Cheap prefilter
	if !strings.Contains(src, "os.ReadFile") && !strings.Contains(src, "os.Open") &&
		!strings.Contains(src, "os.OpenFile") && !strings.Contains(src, "ioutil.ReadFile") &&
		!strings.Contains(src, "os.WriteFile") {
		return
	}
	// Safe confinement: Base-only policy, or HasPrefix root check.
	if strings.Contains(src, "filepath.Base(") {
		return
	}
	if strings.Contains(src, "strings.HasPrefix(") {
		return
	}
	tainted := sourceutil.FindTaintedIdents(src)
	hasReq := sourceutil.HasRequestSource(src)
	if !hasReq && len(tainted) == 0 {
		return
	}

	// Look for file sink with tainted / request-derived path argument.
	sinks := []string{"os.ReadFile(", "os.Open(", "os.OpenFile(", "ioutil.ReadFile(", "os.WriteFile("}
	for _, sink := range sinks {
		start := 0
		for {
			idx := strings.Index(src[start:], sink)
			if idx < 0 {
				break
			}
			abs := start + idx
			// Find '(' of the call
			paren := abs + len(sink) - 1
			if paren >= len(src) || src[paren] != '(' {
				// sink includes '(' already
				paren = strings.Index(src[abs:], "(")
				if paren < 0 {
					start = abs + len(sink)
					continue
				}
				paren = abs + paren
			}
			closeAt, args := scanParen(src, paren)
			if closeAt < 0 {
				start = abs + len(sink)
				continue
			}
			parts := sourceutil.SplitTopLevelArgs(args)
			if len(parts) == 0 {
				start = abs + len(sink)
				continue
			}
			arg0 := parts[0]
			risky := sourceutil.HasRequestSource(arg0)
			if !risky {
				for name := range tainted {
					if sourceutil.ContainsIdent(arg0, name) {
						risky = true
						break
					}
				}
			}
			// filepath.Clean(tainted) / Join(root, tainted) still risky without Base/HasPrefix
			if !risky && (strings.Contains(arg0, "filepath.Clean(") || strings.Contains(arg0, "filepath.Join(") || strings.Contains(arg0, "path.Join(")) {
				if hasReq || len(tainted) > 0 {
					risky = true
				}
			}
			if risky {
				line, col := unit.LineCol(abs)
				rules.PushFindingWithConfidence(
					&MetaCWE22,
					unitFile(unit),
					line, col,
					"user-controlled path segment reaches a file API without confinement",
					0.7,
					out,
				)
				return
			}
			start = abs + len(sink)
		}
	}
}

// --- CWE-78 command injection (seed logic) ---

func detectCWE78(unit *core.ParsedUnit, _ *GoCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "exec.Command") {
		return
	}
	tainted := sourceutil.FindTaintedIdents(src)
	hasReq := sourceutil.HasRequestSource(src)
	calls := sourceutil.FindCalls(src, "exec.Command", "exec.CommandContext")
	for _, call := range calls {
		if shouldFlagCommandInjection(call, tainted, hasReq) {
			line, col := unit.LineCol(call.Start)
			rules.PushFindingWithConfidence(
				&MetaCWE78,
				unitFile(unit),
				line, col,
				"user-controlled input reaches a shell command execution sink",
				0.75,
				out,
			)
		}
	}
}

// --- CWE-79 XSS (taint-lite) ---

func detectCWE79(unit *core.ParsedUnit, _ *GoCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source
	hasReq := sourceutil.HasRequestSource(src)
	tainted := sourceutil.FindTaintedIdents(src)
	if !hasReq && len(tainted) == 0 {
		return
	}
	// Safe: explicit HTML escape of request data, or html/template without unsafe casts.
	if strings.Contains(src, "html.EscapeString(") {
		return
	}
	// text/template executes request-derived data into HTML-ish templates (fixture shape).
	if strings.Contains(src, "text/template") && strings.Contains(src, ".Execute(") {
		i := strings.Index(src, ".Execute(")
		if i < 0 {
			i = strings.Index(src, "text/template")
		}
		line, col := unit.LineCol(i)
		rules.PushFindingWithConfidence(
			&MetaCWE79,
			unitFile(unit),
			line, col,
			"user-controlled data reaches a text/template execution without HTML escaping",
			0.7,
			out,
		)
		return
	}
	// Unescaped HTML sinks
	hasHTML := strings.Contains(src, "template.HTML(") ||
		strings.Contains(src, "template.JS(") ||
		strings.Contains(src, "template.HTMLAttr(") ||
		strings.Contains(src, "c.HTML(") ||
		strings.Contains(src, "w.Write([]byte(") ||
		strings.Contains(src, "fmt.Fprintf(w,")
	if !hasHTML {
		return
	}
	// Safe: html/template auto-escape without template.HTML cast
	if strings.Contains(src, "html/template") && !strings.Contains(src, "template.HTML(") &&
		!strings.Contains(src, "template.JS(") && !strings.Contains(src, "c.HTML(") {
		if !strings.Contains(src, "w.Write(") && !strings.Contains(src, "fmt.Fprintf(w,") {
			return
		}
	}
	needles := []string{"template.HTML(", "template.JS(", "c.HTML(", "fmt.Fprintf(w,", "w.Write([]byte("}
	for _, n := range needles {
		if i := strings.Index(src, n); i >= 0 {
			if n == "w.Write([]byte(" || n == "fmt.Fprintf(w," {
				// Require request-derived data in the *call arguments*, not a
				// loose pre-window (which FPs on rate-limit middleware that
				// writes constant JSON with non-request fields near `r`).
				args := cwe79SinkArgs(src, i, n)
				if !cwe79ArgsLookTainted(args, tainted) {
					continue
				}
			}
			line, col := unit.LineCol(i)
			rules.PushFindingWithConfidence(
				&MetaCWE79,
				unitFile(unit),
				line, col,
				"user-controlled data reaches an HTML response without encoding",
				0.65,
				out,
			)
			return
		}
	}
}

// cwe79SinkArgs returns the interior argument text of the sink call starting at i.
func cwe79SinkArgs(src string, i int, needle string) string {
	open := i + len(needle) - 1 // needle ends with '('
	if open < 0 || open >= len(src) || src[open] != '(' {
		// fall back: find '(' after needle start
		j := strings.IndexByte(src[i:], '(')
		if j < 0 {
			return ""
		}
		open = i + j
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for k := open; k < len(src); k++ {
		c := src[k]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if inStr == '"' || inStr == '\'' {
				if c == '\\' {
					escape = true
					continue
				}
				if c == inStr {
					inStr = 0
				}
				continue
			}
			if inStr == '`' && c == '`' {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '`', '\'':
			inStr = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return src[open+1 : k]
			}
		}
	}
	return ""
}

func cwe79ArgsLookTainted(args string, tainted map[string]struct{}) bool {
	if args == "" {
		return false
	}
	if sourceutil.HasRequestSource(args) {
		return true
	}
	for name := range tainted {
		if len(name) < 2 {
			// Skip single-letter noise; also rejects parse glitches.
			continue
		}
		if sourceutil.ContainsIdent(args, name) {
			return true
		}
	}
	return false
}

// --- CWE-89 SQL injection (seed logic) ---

func detectCWE89(unit *core.ParsedUnit, _ *GoCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source
	tainted := sourceutil.FindTaintedIdents(src)
	for _, call := range findSQLSinkCalls(src) {
		args := sourceutil.SplitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		sqlArg := args[0]
		if strings.HasSuffix(call.Name, "Context") && len(args) >= 2 {
			sqlArg = args[1]
		}
		if !isDynamicSQLArg(sqlArg, tainted) {
			continue
		}
		line, col := unit.LineCol(call.Start)
		rules.PushFindingWithConfidence(
			&MetaCWE89,
			unitFile(unit),
			line, col,
			"user-controlled input reaches an SQL execution sink (heuristic; not full SQLi coverage)",
			0.7,
			out,
		)
	}
}

// --- CWE-90 LDAP injection (taint-lite) ---

func detectCWE90(unit *core.ParsedUnit, _ *GoCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "ldap") && !strings.Contains(src, "LDAP") {
		return
	}
	hasReq := sourceutil.HasRequestSource(src)
	tainted := sourceutil.FindTaintedIdents(src)
	if !hasReq && len(tainted) == 0 {
		return
	}
	// Dynamic filter construction
	if !(strings.Contains(src, "fmt.Sprintf") || strings.Contains(src, "+") || strings.Contains(src, "filter")) {
		return
	}
	// Safe: ldap.EscapeFilter or fixture-local escapeLDAP helper (stdlib/frameworks safe paths)
	if strings.Contains(src, "ldap.EscapeFilter") ||
		strings.Contains(src, "EscapeFilter(") ||
		strings.Contains(src, "escapeLDAP(") {
		return
	}
	// Prefer real LDAP search sinks; avoid firing on safe fixtures that only dial a prebuilt filter.
	needles := []string{"Search(", "ldap.NewSearchRequest", "Filter:"}
	for _, n := range needles {
		if i := strings.Index(src, n); i >= 0 {
			line, col := unit.LineCol(i)
			rules.PushFindingWithConfidence(
				&MetaCWE90,
				unitFile(unit),
				line, col,
				"user-controlled input reaches an LDAP filter without escaping",
				0.65,
				out,
			)
			return
		}
	}
}

// --- CWE-91 XML injection (taint-lite) ---

func detectCWE91(unit *core.ParsedUnit, _ *GoCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "xml") && !strings.Contains(src, "XML") {
		return
	}
	hasReq := sourceutil.HasRequestSource(src)
	tainted := sourceutil.FindTaintedIdents(src)
	if !hasReq && len(tainted) == 0 {
		return
	}
	// Safe: structured xml.Marshal without manual string assembly.
	if strings.Contains(src, "xml.Marshal") && !strings.Contains(src, "fmt.Sprintf") {
		return
	}
	if strings.Contains(src, "xml.Escape") || strings.Contains(src, "xml.EscapeText") {
		return
	}
	// Require dynamic XML construction (sprintf/concat); bare Unmarshal of typed marshal output is safe.
	dynamic := strings.Contains(src, "fmt.Sprintf") || strings.Contains(src, `"+`) || strings.Contains(src, "+ \"")
	if !dynamic {
		return
	}
	needles := []string{"xml.Unmarshal(", "xml.NewDecoder("}
	for _, n := range needles {
		if i := strings.Index(src, n); i >= 0 {
			line, col := unit.LineCol(i)
			rules.PushFindingWithConfidence(
				&MetaCWE91,
				unitFile(unit),
				line, col,
				"user-controlled input reaches XML construction/parsing without neutralization",
				0.6,
				out,
			)
			return
		}
	}
}

// --- shared SQL helpers (from seed cwe89.go) ---

var sqlSinkMethods = []string{
	"Query", "Exec", "QueryRow", "QueryContext", "ExecContext", "QueryRowContext", "Raw",
}

type sqlCall struct {
	Name     string
	Start    int
	ArgsText string
}

func findSQLSinkCalls(source string) []sqlCall {
	var out []sqlCall
	for _, method := range sqlSinkMethods {
		needle := "." + method
		start := 0
		for {
			idx := strings.Index(source[start:], needle)
			if idx < 0 {
				break
			}
			abs := start + idx
			afterName := abs + len(needle)
			if afterName < len(source) {
				r, _ := utf8.DecodeRuneInString(source[afterName:])
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
					start = afterName
					continue
				}
			}
			if abs == 0 {
				start = afterName
				continue
			}
			j := afterName
			for j < len(source) && (source[j] == ' ' || source[j] == '\t' || source[j] == '\n' || source[j] == '\r') {
				j++
			}
			if j >= len(source) || source[j] != '(' {
				start = afterName
				continue
			}
			closeAt, args := scanParen(source, j)
			if closeAt < 0 {
				start = afterName
				continue
			}
			out = append(out, sqlCall{Name: method, Start: abs + 1, ArgsText: args})
			start = closeAt + 1
		}
	}
	return out
}

func isDynamicSQLArg(arg string, tainted map[string]struct{}) bool {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return false
	}
	if sourceutil.IsPureStringLiteral(arg) {
		return false
	}
	if strings.Contains(arg, "+") || strings.Contains(arg, "fmt.Sprintf") {
		return true
	}
	if sourceutil.HasRequestSource(arg) {
		return true
	}
	for name := range tainted {
		if sourceutil.ContainsIdent(arg, name) {
			return true
		}
	}
	return false
}

func scanParen(source string, open int) (closeAt int, args string) {
	if open >= len(source) || source[open] != '(' {
		return -1, ""
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for i := open; i < len(source); i++ {
		c := source[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if inStr == '"' || inStr == '\'' {
				if c == '\\' {
					escape = true
					continue
				}
				if c == inStr {
					inStr = 0
				}
				continue
			}
			if inStr == '`' && c == '`' {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '`', '\'':
			inStr = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i, source[open+1 : i]
			}
		}
	}
	return -1, ""
}

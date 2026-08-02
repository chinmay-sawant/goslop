package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-PY-35", detectBPPY35)
	RegisterRule("BP-PY-36", detectBPPY36)
	RegisterRule("BP-PY-37", detectBPPY37)
}

var (
	sessionAssignRe = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(SessionLocal|Session|sessionmaker\s*\([^)]*\)\s*)\s*\(`)
	sessionCloseRe  = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*\.\s*close\s*\(`)
	executeCallRe   = regexp.MustCompile(`\.\s*execute\s*\(`)
)

func hasSQLAlchemyImport(src string) bool {
	return strings.Contains(src, "sqlalchemy") || strings.Contains(src, "from sqlalchemy") ||
		strings.Contains(src, "import sqlalchemy") || strings.Contains(src, "SessionLocal") ||
		strings.Contains(src, "sessionmaker")
}

// looksDynamicSQL reports f-string / .format / % / + construction of SQL text.
func looksDynamicSQL(arg string) bool {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return false
	}
	if isFStringArg(arg) {
		return true
	}
	if strings.Contains(arg, ".format(") {
		return true
	}
	// "..." % (...)  or '...' % name
	if isPercentFormatted(arg) {
		return true
	}
	// Concatenation of string fragments
	if strings.Contains(arg, " + ") {
		return true
	}
	return false
}

// sqlFStringIsParamsOnly reports an f-string SQL literal whose every
// interpolation is a bound-placeholder token: a .param() call directly or a
// local alias assigned from one (placeholder = store.param()). Real values
// travel through the second execute argument, so the SQL text is static.
// Mirrors cwe.sqlFStringIsParamsOnly without importing the cwe package.
func sqlFStringIsParamsOnly(source string, callStart int, expr string) bool {
	parts := bpFStringInterpolations(expr)
	if len(parts) == 0 {
		return false
	}
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t == "" {
			continue
		}
		if isParamPlaceholderCall(t) {
			continue
		}
		if !isSimpleIdent(t) || !identAssignedFromParamCall(source, callStart, t, map[string]bool{}) {
			return false
		}
	}
	return true
}

// isParamPlaceholderCall reports a `<obj>.param(...)` call expression.
func isParamPlaceholderCall(t string) bool {
	idx := strings.LastIndex(t, ".param")
	if idx < 0 {
		return false
	}
	open := idx + len(".param")
	if open >= len(t) || t[open] != '(' {
		return false
	}
	closeAt := matchingParenClose(t, open)
	return closeAt >= 0 && strings.TrimSpace(t[closeAt+1:]) == ""
}

// identAssignedFromParamCall resolves ident through local assignments to a
// .param() placeholder call (placeholder = self.store.param()).
func identAssignedFromParamCall(source string, callStart int, ident string, seen map[string]bool) bool {
	if seen[ident] {
		return false
	}
	seen[ident] = true
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	lines := strings.Split(source[:callStart], "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") {
			// A function parameter is runtime-controlled, not a placeholder.
			return false
		}
		lhs, rhs, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		if strings.TrimSpace(lhs) != ident {
			continue
		}
		rhs = stripPyComment(strings.TrimSpace(rhs))
		if isParamPlaceholderCall(rhs) {
			return true
		}
		if isSimpleIdent(rhs) {
			return identAssignedFromParamCall(source, callStart, rhs, seen)
		}
		return false
	}
	return false
}

// matchingParenClose returns the index of the ')' matching source[open] == '('.
func matchingParenClose(source string, open int) int {
	if open < 0 || open >= len(source) || source[open] != '(' {
		return -1
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
			if c == '\\' {
				escape = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'':
			inStr = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// bpFStringInterpolations returns the {…} expression bodies of an f-string
// literal expression, skipping escaped braces ({{, }}).
func bpFStringInterpolations(expr string) []string {
	t := strings.TrimSpace(expr)
	i := 0
	for i < len(t) {
		c := t[i]
		if c == 'f' || c == 'F' || c == 'r' || c == 'R' || c == 'b' || c == 'B' || c == 'u' || c == 'U' {
			i++
			continue
		}
		break
	}
	if i >= len(t) {
		return nil
	}
	quote := t[i]
	if quote != '"' && quote != '\'' {
		return nil
	}
	var body string
	if i+2 < len(t) && t[i+1] == quote && t[i+2] == quote {
		after := t[i+3:]
		end := strings.Index(after, string([]byte{quote, quote, quote}))
		if end < 0 {
			return nil
		}
		body = after[:end]
	} else {
		escaped := false
		for j := i + 1; j < len(t); j++ {
			if escaped {
				escaped = false
				continue
			}
			if t[j] == '\\' {
				escaped = true
				continue
			}
			if t[j] == quote {
				body = t[i+1 : j]
				break
			}
		}
		if body == "" {
			return nil
		}
	}
	return bpCollectFStringParts(body)
}

func bpCollectFStringParts(body string) []string {
	var out []string
	depth := 0
	var cur strings.Builder
	for k := 0; k < len(body); k++ {
		c := body[k]
		if c == '{' {
			if k+1 < len(body) && body[k+1] == '{' {
				k++
				continue
			}
			if depth == 0 {
				cur.Reset()
			}
			depth++
			continue
		}
		if c == '}' {
			if k+1 < len(body) && body[k+1] == '}' {
				k++
				continue
			}
			depth--
			if depth == 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
			continue
		}
		if depth > 0 {
			cur.WriteByte(c)
		}
	}
	return out
}

func isPercentFormatted(arg string) bool {
	// Detect string % expr at top level (not SQL bind "%s" alone inside a plain string).
	// Patterns: "x %s" % y   or   'x' % (y,)
	depth := 0
	inStr := byte(0)
	escape := false
	triple := false
	for i := 0; i < len(arg); i++ {
		c := arg[i]
		if inStr != 0 {
			i = consumeQuotedCharacter(arg, i, &inStr, &escape, &triple)
			continue
		}
		if c == '"' || c == '\'' {
			inStr = c
			triple = i+2 < len(arg) && arg[i+1] == c && arg[i+2] == c
			if triple {
				i += 2
			}
			continue
		}
		switch c {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '%':
			if i+1 < len(arg) && arg[i+1] == '%' {
				i++
				continue
			}
			if depth == 0 && isStringFormattingOperator(arg, i) {
				return true
			}
		}
	}
	return false
}

func consumeQuotedCharacter(arg string, index int, quote *byte, escaped *bool, triple *bool) int {
	c := arg[index]
	if *escaped {
		*escaped = false
		return index
	}
	if c == '\\' && !*triple {
		*escaped = true
		return index
	}
	if !*triple {
		if c == *quote {
			*quote = 0
		}
		return index
	}
	if c == *quote && index+2 < len(arg) && arg[index+1] == *quote && arg[index+2] == *quote {
		*quote = 0
		*triple = false
		return index + 2
	}
	return index
}

func isStringFormattingOperator(arg string, operatorIdx int) bool {
	// We only see % outside strings, so the left side should end with a string.
	left := strings.TrimSpace(arg[:operatorIdx])
	right := strings.TrimSpace(arg[operatorIdx+1:])
	if left == "" || right == "" {
		return false
	}
	return strings.HasSuffix(left, `"`) || strings.HasSuffix(left, `'`) ||
		strings.HasSuffix(left, `"""`) || strings.HasSuffix(left, `'''`)
}

// BP-PY-35: sqlalchemy.text with f-string / format / % / concat.
// Policy: fire on text(...); BP-PY-37 covers bare .execute(f"...") without text().
func detectBPPY35(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-35")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "text(") && !strings.Contains(src, "text (") {
		return
	}
	// Require sqlalchemy context to reduce false positives on other text() helpers.
	if !hasSQLAlchemyImport(src) && !strings.Contains(src, "sa.text") &&
		!strings.Contains(src, "sqlalchemy.text") {
		return
	}

	start := 0
	for {
		// Find text(
		idx := indexOfIdent(src[start:], "text")
		if idx < 0 {
			return
		}
		abs := start + idx
		// Must be text( or sa.text / sqlalchemy.text — check prefix
		if abs > 0 && isIdentByte(src[abs-1]) && !isSQLAlchemyTextCall(src, abs) {
			start = abs + 4
			continue
		}
		end := abs + len("text")
		if end >= len(src) {
			return
		}
		// optional whitespace then (
		rest := src[end:]
		sp := 0
		for sp < len(rest) && (rest[sp] == ' ' || rest[sp] == '\t') {
			sp++
		}
		if sp >= len(rest) || rest[sp] != '(' {
			start = abs + 4
			continue
		}
		openAbs := end + sp
		arg, ok := firstCallArg(src, openAbs)
		if !ok {
			start = abs + 4
			continue
		}
		// Also: text("...".format(...)) — first arg may be the whole call chain as one expression
		// firstCallArg returns first top-level arg; "SELECT {}".format(x) is one arg with .format inside.
		if looksDynamicSQL(arg) {
			pushAt(unit, meta, abs, "sqlalchemy.text builds SQL with f-string/format; use bound params (:name)", out)
		}
		start = abs + 4
	}
}

func isSQLAlchemyTextCall(src string, textStart int) bool {
	return (textStart >= len("sqlalchemy.") && src[textStart-len("sqlalchemy."):textStart] == "sqlalchemy.") ||
		(textStart >= len("sa.") && src[textStart-len("sa."):textStart] == "sa.")
}

// BP-PY-36: Session / SessionLocal created without with or .close().
// Full CFG proof of all exit paths is out of scope for v0.
func detectBPPY36(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-36")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "Session") && !strings.Contains(src, "SessionLocal") &&
		!strings.Contains(src, "sessionmaker") {
		return
	}

	lines := codeLinesFacts(facts, src)
	// Skip files that only use with Session... as session
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		// with SessionLocal() as session: / with Session() as s:
		if strings.HasPrefix(t, "with ") || strings.HasPrefix(t, "async with ") {
			continue
		}
		name := sessionAssignmentName(t)
		if name == "" || strings.Contains(t, "with ") {
			continue
		}
		if sessionUnclosedInScope(lines, i, name) {
			pushAt(unit, meta, line.byte, "SQLAlchemy Session created without with/close; use a context manager or session.close()", out)
		}
	}
}

func sessionAssignmentName(line string) string {
	if match := sessionAssignRe.FindStringSubmatchIndex(line); match != nil {
		return line[match[2]:match[3]]
	}
	if !strings.Contains(line, "SessionLocal(") && !strings.Contains(line, "Session(") {
		return ""
	}
	equals := strings.Index(line, "=")
	if equals <= 0 {
		return ""
	}
	name := strings.TrimSpace(line[:equals])
	right := strings.TrimSpace(line[equals+1:])
	if !isSimpleIdent(name) {
		return ""
	}
	if !strings.HasPrefix(right, "SessionLocal(") && !strings.HasPrefix(right, "Session(") &&
		!strings.Contains(right, "sessionmaker(") {
		return ""
	}
	if strings.HasPrefix(right, "sessionmaker") && !strings.Contains(right, ")(") {
		return ""
	}
	return name
}

// sessionUnclosedInScope reports whether `name` has no .close() and was not created via with
// within the same function (indent window) as assignLineIdx.
func sessionUnclosedInScope(lines []codeLine, assignLineIdx int, name string) bool {
	// Find enclosing function
	defIdx := -1
	for i := assignLineIdx; i >= 0; i-- {
		t := strings.TrimSpace(lines[i].text)
		if strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") {
			defIdx = i
			break
		}
	}
	var bodyStart, bodyEnd int
	if defIdx >= 0 {
		bodyStart, bodyEnd = functionBodyRange(lines, defIdx)
	} else {
		// Module level: scan whole file
		bodyStart, bodyEnd = 0, len(lines)
	}
	// If assignment is on a with line, skip (already handled)
	at := strings.TrimSpace(lines[assignLineIdx].text)
	if strings.Contains(at, "with ") {
		return false
	}
	// Look for with ... name or name.close() in body
	closeNeedle := name + ".close("
	for j := bodyStart; j < bodyEnd; j++ {
		bt := lines[j].text
		if strings.Contains(bt, closeNeedle) || sessionCloseRe.MatchString(bt) && strings.Contains(bt, name+".close") {
			return false
		}
		// with SessionLocal() as name
		if strings.Contains(bt, " as "+name) && (strings.Contains(bt, "with ") || strings.Contains(bt, "with\t")) {
			return false
		}
		// contextlib.closing
		if strings.Contains(bt, "closing(") && strings.Contains(bt, name) {
			return false
		}
	}
	// Also check remaining lines after assignment in same indent scope if no def found
	if defIdx < 0 {
		for j := assignLineIdx; j < len(lines); j++ {
			if strings.Contains(lines[j].text, closeNeedle) {
				return false
			}
		}
	}
	return true
}

// BP-PY-37: DB-API / driver .execute with f-string / % / .format SQL.
// Overlap policy: session.execute(text(f"...")) may also fire BP-PY-35 on text(;
// this rule targets the .execute( call when the first arg itself is dynamic SQL.
func detectBPPY37(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-37")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, ".execute(") && !strings.Contains(src, ".execute (") {
		return
	}

	start := 0
	for {
		loc := executeCallRe.FindStringIndex(src[start:])
		if loc == nil {
			return
		}
		abs := start + loc[0]
		// open paren is last char of match
		openAbs := abs + strings.Index(src[abs:start+loc[1]], "(")
		if openAbs < abs {
			// fallback
			openAbs = strings.Index(src[abs:], "(")
			if openAbs < 0 {
				start = abs + 8
				continue
			}
			openAbs += abs
		}
		arg, ok := firstCallArg(src, openAbs)
		if !ok {
			start = abs + 8
			continue
		}
		arg = strings.TrimSpace(arg)
		// Miss: static SQL string with bound params as second arg
		// execute("SELECT ... %s", (uid,)) — first arg is plain string, not percent-formatted expression
		if isPlainStringLiteral(arg) {
			start = abs + 8
			continue
		}
		// Miss: text("SELECT ... :id") or text with binds — static text() is OK; dynamic text handled by 35
		// If first arg is text(...), only flag here if the text( arg is dynamic AND we want double fire.
		// Policy: if arg starts with text(, let BP-PY-35 handle; skip 37 to reduce duplicate noise on same call.
		if isTextCallArg(arg) {
			start = abs + 8
			continue
		}
		if looksDynamicSQL(arg) {
			// Miss: f"... {store.param()} ..." where every interpolation is a
			// bound-placeholder token (or an alias assigned from .param()); values
			// are passed as the second execute argument (Project_Parva shape).
			if isFStringArg(arg) && sqlFStringIsParamsOnly(src, abs, arg) {
				start = abs + 8
				continue
			}
			pushAt(unit, meta, abs, "cursor/connection execute builds SQL with f-string or % format; use bound parameters", out)
		}
		start = abs + 8
	}
}

func isTextCallArg(arg string) bool {
	arg = strings.TrimSpace(arg)
	// text(...) or sa.text(...) or sqlalchemy.text(...)
	if strings.HasPrefix(arg, "text(") {
		return true
	}
	if strings.HasPrefix(arg, "sa.text(") || strings.HasPrefix(arg, "sqlalchemy.text(") {
		return true
	}
	return false
}

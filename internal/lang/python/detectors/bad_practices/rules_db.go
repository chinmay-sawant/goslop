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
			if escape {
				escape = false
				continue
			}
			if c == '\\' && !triple {
				escape = true
				continue
			}
			if triple {
				if c == inStr && i+2 < len(arg) && arg[i+1] == inStr && arg[i+2] == inStr {
					inStr = 0
					triple = false
					i += 2
				}
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		if c == '"' || c == '\'' {
			if i+2 < len(arg) && arg[i+1] == c && arg[i+2] == c {
				inStr = c
				triple = true
				i += 2
				continue
			}
			inStr = c
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
			if depth == 0 {
				// Python modulo formatting: string-lit % expr
				// Ensure not %% and not mid-ident
				if i+1 < len(arg) && arg[i+1] == '%' {
					i++
					continue
				}
				// Look left for a string end (we only see % outside strings, so left should be closed string)
				left := strings.TrimSpace(arg[:i])
				right := strings.TrimSpace(arg[i+1:])
				if left != "" && right != "" && (strings.HasSuffix(left, `"`) || strings.HasSuffix(left, `'`) ||
					strings.HasSuffix(left, `"""`) || strings.HasSuffix(left, `'''`)) {
					return true
				}
			}
		}
	}
	return false
}

// BP-PY-35: sqlalchemy.text with f-string / format / % / concat.
// Policy: fire on text(...); BP-PY-37 covers bare .execute(f"...") without text().
func detectBPPY35(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
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
		if abs > 0 && isIdentByte(src[abs-1]) {
			// Could be sqlalchemy.text — check
			if abs >= 11 && src[abs-11:abs] == "sqlalchemy." {
				// ok
			} else if abs >= 3 && src[abs-3:abs] == "sa." {
				// ok
			} else {
				start = abs + 4
				continue
			}
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
		arg, _, ok := firstCallArg(src, openAbs)
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
		// Match assignment session = SessionLocal( or Session(
		m := sessionAssignRe.FindStringSubmatchIndex(t)
		if m == nil {
			// Also: session = SessionLocal() without nested sessionmaker complexity
			if !strings.Contains(t, "SessionLocal(") && !strings.Contains(t, "Session(") {
				continue
			}
			// Manual parse: name = SessionLocal( or name = Session(
			eq := strings.Index(t, "=")
			if eq <= 0 {
				continue
			}
			lhs := strings.TrimSpace(t[:eq])
			rhs := strings.TrimSpace(t[eq+1:])
			if !isSimpleIdent(lhs) {
				continue
			}
			if !strings.HasPrefix(rhs, "SessionLocal(") && !strings.HasPrefix(rhs, "Session(") &&
				!strings.Contains(rhs, "sessionmaker(") {
				continue
			}
			// sessionmaker()() double call
			if strings.HasPrefix(rhs, "sessionmaker") && !strings.Contains(rhs, ")(") {
				// sessionmaker(...) alone is factory, not session
				if !strings.Contains(rhs, ")(") {
					continue
				}
			}
			if sessionUnclosedInScope(lines, i, lhs) {
				pushAt(unit, meta, line.byte, "SQLAlchemy Session created without with/close; use a context manager or session.close()", out)
			}
			continue
		}
		name := t[m[2]:m[3]]
		// Ensure this line is not inside a with header on same line
		if strings.Contains(t, "with ") {
			continue
		}
		if sessionUnclosedInScope(lines, i, name) {
			pushAt(unit, meta, line.byte, "SQLAlchemy Session created without with/close; use a context manager or session.close()", out)
		}
	}
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
func detectBPPY37(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
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
		arg, _, ok := firstCallArg(src, openAbs)
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

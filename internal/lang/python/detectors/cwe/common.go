package cwe

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/chinmay-sawant/goslop/internal/core"
)

// unitFile returns the display path for findings.
func unitFile(unit *core.ParsedUnit) string {
	if unit == nil {
		return ""
	}
	if unit.DisplayPath != "" {
		return unit.DisplayPath
	}
	return unit.Path
}

// callSite is a lightweight function/method call match in source text.
type callSite struct {
	Name     string // callee text matched (e.g. "pickle.loads")
	Start    int    // byte offset of Name
	End      int    // past-the-end of matching call
	ArgsText string // interior of top-level argument list
}

// findCalls finds call-like sites for each name in names (exact callee strings).
func findCalls(source string, names ...string) []callSite {
	var out []callSite
	code := pythonCodeMask(source)
	for _, name := range names {
		start := 0
		for {
			idx := strings.Index(code[start:], name)
			if idx < 0 {
				break
			}
			abs := start + idx
			if !identBoundaryOK(code, abs, abs+len(name)) {
				start = abs + len(name)
				continue
			}
			after := abs + len(name)
			j := skipWS(code, after)
			if j >= len(code) || code[j] != '(' {
				start = after
				continue
			}
			closeAt, args := scanCallArgs(source, j)
			if closeAt < 0 {
				start = after
				continue
			}
			out = append(out, callSite{
				Name:     name,
				Start:    abs,
				End:      closeAt + 1,
				ArgsText: args,
			})
			start = closeAt + 1
		}
	}
	return out
}

// pythonCodeMask keeps byte offsets stable while blanking comments and string
// literals. Source-pattern rules can use it to avoid interpreting examples in
// docstrings, comments, and quoted data as executable Python.
func pythonCodeMask(source string) string {
	masked := []byte(source)
	inString := byte(0)
	triple := false
	escaped := false
	inComment := false
	for i := 0; i < len(masked); i++ {
		c := masked[i]
		if inComment {
			if c == '\n' {
				inComment = false
			} else {
				masked[i] = ' '
			}
			continue
		}
		if inString != 0 {
			masked[i] = ' '
			if triple {
				if c == inString && i+2 < len(masked) && masked[i+1] == inString && masked[i+2] == inString {
					masked[i+1], masked[i+2] = ' ', ' '
					i += 2
					inString = 0
					triple = false
				}
				continue
			}
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == inString {
				inString = 0
			}
			continue
		}
		switch c {
		case '#':
			masked[i] = ' '
			inComment = true
		case '\'', '"':
			inString = c
			if i+2 < len(masked) && masked[i+1] == c && masked[i+2] == c {
				masked[i], masked[i+1], masked[i+2] = ' ', ' ', ' '
				i += 2
				triple = true
			} else {
				masked[i] = ' '
			}
		}
	}
	return string(masked)
}

func firstCodeMatchStart(source string, pattern *regexp.Regexp) int {
	if pattern == nil {
		return -1
	}
	masked := pythonCodeMask(source)
	for _, match := range pattern.FindAllStringIndex(source, -1) {
		if strings.TrimSpace(masked[match[0]:match[1]]) != "" {
			return match[0]
		}
	}
	return -1
}

func identBoundaryOK(source string, start, end int) bool {
	// Method-form needles like ".execute" intentionally start with '.' so the
	// left side may be an identifier (cursor.execute). Only enforce the left
	// boundary for names that do not begin with '.'.
	if start > 0 && (start >= len(source) || source[start] != '.') {
		r, _ := utf8.DecodeLastRuneInString(source[:start])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' {
			return false
		}
	}
	if end < len(source) {
		r, _ := utf8.DecodeRuneInString(source[end:])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return false
		}
	}
	return true
}

func skipWS(s string, i int) int {
	for i < len(s) {
		switch s[i] {
		case ' ', '\t', '\n', '\r':
			i++
		default:
			return i
		}
	}
	return i
}

// scanCallArgs starts at '(' and returns index of matching ')' and interior args.
func scanCallArgs(source string, open int) (closeAt int, args string) {
	if open >= len(source) || source[open] != '(' {
		return -1, ""
	}
	depth := 0
	inStr := byte(0)
	escape := false
	// Python triple-quote tracking: """ or '''
	triple := 0 // 0 none, 3-ish when opening triple
	for i := open; i < len(source); i++ {
		c := source[i]
		if inStr != 0 {
			if triple > 0 {
				// inside triple-quoted string
				if c == inStr && i+2 < len(source) && source[i+1] == inStr && source[i+2] == inStr {
					inStr = 0
					triple = 0
					i += 2
				}
				continue
			}
			if escape {
				escape = false
				continue
			}
			if c == '\\' && (inStr == '"' || inStr == '\'') {
				escape = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		// not in string
		if c == '"' || c == '\'' {
			// check for triple quote
			if i+2 < len(source) && source[i+1] == c && source[i+2] == c {
				inStr = c
				triple = 3
				i += 2
				continue
			}
			inStr = c
			continue
		}
		switch c {
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

// splitTopLevelArgs splits a call argument list on top-level commas.
func splitTopLevelArgs(args string) []string {
	var out []string
	depth := 0
	inStr := byte(0)
	escape := false
	triple := 0
	start := 0
	for i := 0; i < len(args); i++ {
		c := args[i]
		if inStr != 0 {
			if triple > 0 {
				if c == inStr && i+2 < len(args) && args[i+1] == inStr && args[i+2] == inStr {
					inStr = 0
					triple = 0
					i += 2
				}
				continue
			}
			if escape {
				escape = false
				continue
			}
			if c == '\\' && (inStr == '"' || inStr == '\'') {
				escape = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		if c == '"' || c == '\'' {
			if i+2 < len(args) && args[i+1] == c && args[i+2] == c {
				inStr = c
				triple = 3
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
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(args[start:i]))
				start = i + 1
			}
		}
	}
	rest := strings.TrimSpace(args[start:])
	if rest != "" || len(out) > 0 {
		out = append(out, rest)
	}
	return out
}

// isPureStringLiteral reports whether expr is exactly one Python string literal
// (optional string prefix b/r/u/f/fr… then ' or " or triple quotes). Pure f-strings
// are treated as dynamic (not pure) because they embed expressions.
func isPureStringLiteral(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	// Strip string prefixes (order-insensitive common set).
	i := 0
	hasF := false
	for i < len(t) {
		c := t[i]
		if c == 'f' || c == 'F' {
			hasF = true
			i++
			continue
		}
		if c == 'b' || c == 'B' || c == 'r' || c == 'R' || c == 'u' || c == 'U' {
			i++
			continue
		}
		break
	}
	if hasF {
		// f-strings embed expressions → dynamic for injection purposes
		return false
	}
	if i >= len(t) {
		return false
	}
	quote := t[i]
	if quote != '"' && quote != '\'' {
		return false
	}
	// triple?
	if i+2 < len(t) && t[i+1] == quote && t[i+2] == quote {
		end := strings.Index(t[i+3:], string([]byte{quote, quote, quote}))
		if end < 0 {
			return false
		}
		return strings.TrimSpace(t[i+3+end+3:]) == ""
	}
	// single-quoted string with escapes
	escape := false
	for j := i + 1; j < len(t); j++ {
		if escape {
			escape = false
			continue
		}
		if t[j] == '\\' {
			escape = true
			continue
		}
		if t[j] == quote {
			return strings.TrimSpace(t[j+1:]) == ""
		}
	}
	return false
}

// isDynamicExpr reports whether expr looks non-constant (variable, format, concat, f-string).
func isDynamicExpr(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if isPureStringLiteral(t) {
		return false
	}
	// list/tuple of pure string literals is still "static argv"
	if looksStaticStringList(t) {
		return false
	}
	return true
}

// looksStaticStringList is true for ["a","b"] / ('a', 'b') of pure literals only.
func looksStaticStringList(expr string) bool {
	t := strings.TrimSpace(expr)
	if len(t) < 2 {
		return false
	}
	open, close := t[0], byte(0)
	switch open {
	case '[':
		close = ']'
	case '(':
		close = ')'
	default:
		return false
	}
	if t[len(t)-1] != close {
		return false
	}
	inner := strings.TrimSpace(t[1 : len(t)-1])
	if inner == "" {
		return true
	}
	parts := splitTopLevelArgs(inner)
	if len(parts) == 0 {
		return false
	}
	for _, p := range parts {
		if !isPureStringLiteral(p) {
			return false
		}
	}
	return true
}

// hasKwargTrue reports whether argsText contains an exact name=True keyword.
func hasKwargTrue(argsText, name string) bool {
	return hasBooleanKwarg(argsText, name, "True")
}

// hasKwargFalse reports an exact name=False keyword.
func hasKwargFalse(argsText, name string) bool {
	return hasBooleanKwarg(argsText, name, "False")
}

func hasBooleanKwarg(argsText, name, value string) bool {
	for _, arg := range splitTopLevelArgs(argsText) {
		key, candidate, ok := strings.Cut(arg, "=")
		if !ok || strings.TrimSpace(key) != name {
			continue
		}
		return strings.TrimSpace(candidate) == value
	}
	return false
}

func compactWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// containsIdent reports whole-identifier presence of ident in text.
func containsIdent(text, ident string) bool {
	if ident == "" {
		return false
	}
	start := 0
	for {
		idx := strings.Index(text[start:], ident)
		if idx < 0 {
			return false
		}
		abs := start + idx
		okLeft := true
		if abs > 0 {
			r, _ := utf8.DecodeLastRuneInString(text[:abs])
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				okLeft = false
			}
		}
		okRight := true
		end := abs + len(ident)
		if end < len(text) {
			r, _ := utf8.DecodeRuneInString(text[end:])
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				okRight = false
			}
		}
		if okLeft && okRight {
			return true
		}
		start = abs + len(ident)
	}
}

// looksSQLFormatted reports f-string / % / .format style string construction in expr.
func looksSQLFormatted(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	// f"SELECT ... {x}" or F'...'
	if isFStringExpr(t) {
		return true
	}
	// "... %s ..." % var  or  "...".format(...)
	if strings.Contains(t, ".format(") {
		return true
	}
	// binary % formatting: lhs % rhs (not modulo of numbers only)
	if idx := indexTopLevelPercent(t); idx > 0 {
		lhs := strings.TrimSpace(t[:idx])
		// string % something
		if isPureStringLiteral(lhs) || isFStringExpr(lhs) || strings.HasPrefix(strings.TrimSpace(lhs), "\"") || strings.HasPrefix(strings.TrimSpace(lhs), "'") {
			return true
		}
		// parenthesized concat etc.
		if strings.Contains(lhs, "SELECT") || strings.Contains(lhs, "INSERT") ||
			strings.Contains(lhs, "UPDATE") || strings.Contains(lhs, "DELETE") ||
			strings.Contains(lhs, "select") || strings.Contains(lhs, "insert") {
			return true
		}
	}
	// "SELECT " + var
	if strings.Contains(t, "+") && !isPureStringLiteral(t) {
		if strings.Contains(strings.ToUpper(t), "SELECT") ||
			strings.Contains(strings.ToUpper(t), "INSERT") ||
			strings.Contains(strings.ToUpper(t), "UPDATE") ||
			strings.Contains(strings.ToUpper(t), "DELETE") ||
			strings.Contains(t, "\"") || strings.Contains(t, "'") {
			// dynamic concat of SQL-ish pieces
			if !isPureStringLiteral(t) {
				return true
			}
		}
	}
	return false
}

func isFStringExpr(expr string) bool {
	t := strings.TrimSpace(expr)
	// leading prefixes may include f/F among b/r/u
	i := 0
	hasF := false
	for i < len(t) {
		c := t[i]
		if c == 'f' || c == 'F' {
			hasF = true
			i++
			continue
		}
		if c == 'b' || c == 'B' || c == 'r' || c == 'R' || c == 'u' || c == 'U' {
			i++
			continue
		}
		break
	}
	if !hasF || i >= len(t) {
		return false
	}
	return t[i] == '"' || t[i] == '\''
}

// indexTopLevelPercent finds first '%' not inside strings/parens at depth 0.
func indexTopLevelPercent(s string) int {
	depth := 0
	inStr := byte(0)
	escape := false
	triple := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			if triple > 0 {
				if c == inStr && i+2 < len(s) && s[i+1] == inStr && s[i+2] == inStr {
					inStr = 0
					triple = 0
					i += 2
				}
				continue
			}
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
		if c == '"' || c == '\'' {
			if i+2 < len(s) && s[i+1] == c && s[i+2] == c {
				inStr = c
				triple = 3
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
				// skip %=
				if i+1 < len(s) && s[i+1] == '=' {
					continue
				}
				return i
			}
		}
	}
	return -1
}

// hasSafePathConfinement reports file-level safe patterns for CWE-22.
// Documented suppressions: basename-only policy, or resolve + startswith root check.
func hasSafePathConfinement(src string) bool {
	if strings.Contains(src, "os.path.basename(") || strings.Contains(src, "pathlib.Path(...).name") {
		return true
	}
	// resolve + startswith / is_relative_to confinement
	hasResolve := strings.Contains(src, ".resolve(") || strings.Contains(src, "os.path.realpath(") || strings.Contains(src, "os.path.abspath(")
	hasPrefix := strings.Contains(src, "startswith(") || strings.Contains(src, "is_relative_to(") || strings.Contains(src, "commonpath(")
	return hasResolve && hasPrefix
}

package badpractices

import (
	"path/filepath"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func fileDisplayPath(unit *core.ParsedUnit) string {
	if unit == nil {
		return ""
	}
	if unit.DisplayPath != "" {
		return unit.DisplayPath
	}
	return unit.Path
}

// isPythonTestFile reports test modules we should skip for validation-style rules.
// Matches test_*.py, *_test.py, and paths under tests/ or test/.
func isPythonTestFile(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, p := range []string{fileDisplayPath(unit), unit.Path} {
		if p == "" {
			continue
		}
		base := filepath.Base(p)
		if strings.HasPrefix(base, "test_") && strings.HasSuffix(base, ".py") {
			return true
		}
		if strings.HasSuffix(base, "_test.py") {
			return true
		}
		// Normalize separators.
		norm := filepath.ToSlash(p)
		if strings.Contains(norm, "/tests/") || strings.Contains(norm, "/test/") ||
			strings.HasPrefix(norm, "tests/") || strings.HasPrefix(norm, "test/") {
			return true
		}
	}
	return false
}

func pushAt(unit *core.ParsedUnit, meta *rules.RuleMetadata, byteOffset int, message string, out *[]rules.Finding) {
	if unit == nil || meta == nil || out == nil {
		return
	}
	line, col := unit.LineCol(byteOffset)
	rules.PushFinding(meta, fileDisplayPath(unit), line, col, message, out)
}

// stripPyComment removes a trailing # comment outside of string literals.
func stripPyComment(line string) string {
	inStr := byte(0)
	escape := false
	triple := false
	for i := 0; i < len(line); i++ {
		c := line[i]
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
				if c == inStr && i+2 < len(line) && line[i+1] == inStr && line[i+2] == inStr {
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
			// Triple-quoted?
			if i+2 < len(line) && line[i+1] == c && line[i+2] == c {
				inStr = c
				triple = true
				i += 2
				continue
			}
			inStr = c
			continue
		}
		if c == '#' {
			return strings.TrimRight(line[:i], " \t")
		}
	}
	return line
}

// codeLine is one source line with 0-based index and byte offset.
type codeLine struct {
	idx  int
	text string // comment-stripped
	raw  string
	byte int
}

func buildCodeLines(source string) []codeLine {
	if source == "" {
		return nil
	}
	n := 1
	for i := 0; i < len(source); i++ {
		if source[i] == '\n' {
			n++
		}
	}
	out := make([]codeLine, 0, n)
	byteOff := 0
	lineIdx := 0
	for {
		nl := strings.IndexByte(source[byteOff:], '\n')
		var raw string
		if nl < 0 {
			raw = source[byteOff:]
			out = append(out, codeLine{idx: lineIdx, text: stripPyComment(raw), raw: raw, byte: byteOff})
			break
		}
		raw = source[byteOff : byteOff+nl]
		out = append(out, codeLine{idx: lineIdx, text: stripPyComment(raw), raw: raw, byte: byteOff})
		byteOff += nl + 1
		lineIdx++
		if byteOff >= len(source) {
			out = append(out, codeLine{idx: lineIdx, text: "", raw: "", byte: byteOff})
			break
		}
	}
	return out
}

func codeLinesFacts(facts *bpFacts, source string) []codeLine {
	if facts != nil && facts.lines != nil && (source == "" || source == facts.Source) {
		return facts.lines
	}
	return buildCodeLines(source)
}

func isIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// isSimpleIdent reports whether s is a non-empty Python identifier (no dots).
func isSimpleIdent(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isIdentByte(s[i]) {
			return false
		}
	}
	c := s[0]
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// indexOfIdent finds needle as a whole identifier/attribute fragment (not mid-ident).
func indexOfIdent(source, needle string) int {
	start := 0
	for {
		idx := strings.Index(source[start:], needle)
		if idx < 0 {
			return -1
		}
		abs := start + idx
		if abs > 0 {
			prev := source[abs-1]
			if isIdentByte(prev) {
				start = abs + len(needle)
				continue
			}
		}
		end := abs + len(needle)
		if end < len(source) && isIdentByte(source[end]) {
			start = end
			continue
		}
		return abs
	}
}

// findAllIdent returns all absolute offsets of needle as whole identifiers.
func findAllIdent(source, needle string) []int {
	var out []int
	start := 0
	for {
		idx := indexOfIdent(source[start:], needle)
		if idx < 0 {
			return out
		}
		abs := start + idx
		out = append(out, abs)
		start = abs + len(needle)
		if start >= len(source) {
			return out
		}
	}
}

// indentWidth returns leading space/tab count (tabs count as 1 for relative compare).
func indentWidth(line string) int {
	n := 0
	for _, r := range line {
		if r == ' ' || r == '\t' {
			n++
			continue
		}
		break
	}
	return n
}

// isStringLiteral reports whether s is a Python string/bytes literal (possibly with prefix).
func isStringLiteral(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// Strip common prefixes: r, u, b, f, fr, rf, br, rb (case-insensitive for b/u/r).
	for {
		if len(s) < 2 {
			break
		}
		c0 := s[0]
		if c0 == 'r' || c0 == 'R' || c0 == 'u' || c0 == 'U' || c0 == 'b' || c0 == 'B' || c0 == 'f' || c0 == 'F' {
			s = s[1:]
			continue
		}
		break
	}
	if len(s) < 2 {
		return false
	}
	q := s[0]
	if q != '"' && q != '\'' {
		return false
	}
	// Triple quotes.
	if len(s) >= 6 && s[1] == q && s[2] == q {
		return strings.HasSuffix(s, string([]byte{q, q, q}))
	}
	return s[len(s)-1] == q
}

// firstCallArg extracts the first top-level argument text of a call starting at
// openParen (index of '('). Returns arg, end offset of ')', ok.
func firstCallArg(source string, openParen int) (arg string, close int, ok bool) {
	if openParen < 0 || openParen >= len(source) || source[openParen] != '(' {
		return "", -1, false
	}
	depth := 0
	inStr := byte(0)
	escape := false
	start := openParen + 1
	for i := openParen; i < len(source); i++ {
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
				// Full arg list region.
				inner := source[start:i]
				// Split first top-level comma.
				arg = splitFirstArg(inner)
				return strings.TrimSpace(arg), i, true
			}
		case ',':
			if depth == 1 {
				arg = strings.TrimSpace(source[start:i])
				// Still need to find close for completeness — continue scan? For first arg only we have it.
				// Find matching close.
				for j := i; j < len(source); j++ {
					// reuse simplified: just return with close unknown; scan for close
				}
				// Find close paren at depth 0 from here.
				close = findMatchingClose(source, openParen)
				return arg, close, true
			}
		}
	}
	return "", -1, false
}

func splitFirstArg(inner string) string {
	depth := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(inner); i++ {
		c := inner[i]
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
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				return strings.TrimSpace(inner[:i])
			}
		}
	}
	return strings.TrimSpace(inner)
}

func findMatchingClose(source string, openParen int) int {
	depth := 0
	inStr := byte(0)
	escape := false
	for i := openParen; i < len(source); i++ {
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

// callArgsRegion returns the full argument list text inside (... ) for a call.
func callArgsRegion(source string, openParen int) (inner string, close int, ok bool) {
	close = findMatchingClose(source, openParen)
	if close < 0 {
		return "", -1, false
	}
	return source[openParen+1 : close], close, true
}

// lineContainsWithOpen reports whether the line uses open as a with context.
func lineContainsWithOpen(line string) bool {
	t := strings.TrimSpace(line)
	// with open(...) or with path.open(...)
	if !strings.HasPrefix(t, "with ") && !strings.Contains(t, " with ") {
		// also: async with
		if !strings.HasPrefix(t, "async with ") {
			// may still be "with open" after leading content? rare
			if !strings.Contains(t, "with open(") && !strings.Contains(t, "with Path.open(") {
				// check "with <expr>.open("
				if !containsWithOpenCall(t) {
					return false
				}
			}
		}
	}
	return strings.Contains(t, "open(")
}

func containsWithOpenCall(t string) bool {
	// crude: "with " then later "open("
	idx := strings.Index(t, "with ")
	if idx < 0 {
		idx = strings.Index(t, "async with ")
	}
	if idx < 0 {
		return false
	}
	return strings.Contains(t[idx:], "open(")
}

func looksLikePlaceholderSecret(val string) bool {
	v := strings.ToLower(strings.TrimSpace(val))
	// strip quotes
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			v = v[1 : len(v)-1]
		}
	}
	if v == "" {
		return true
	}
	// Exact / whole-value placeholders only (conservative FP policy).
	exact := map[string]struct{}{
		"changeme": {}, "change-me": {}, "change_me": {},
		"xxx": {}, "xxxxx": {}, "todo": {}, "fix": {}, "placeholder": {},
		"secret": {}, "password": {}, "api_key": {}, "apikey": {}, "token": {},
		"your_secret": {}, "yoursecret": {}, "your_password": {}, "your_token": {},
		"example": {}, "sample": {}, "dummy": {}, "replace_me": {}, "insert_here": {},
		"notasecret": {}, "not-a-secret": {}, "not_a_secret": {},
		"***": {}, "...": {}, "none": {}, "null": {}, "n/a": {}, "na": {},
		"secret_key": {}, "supersecret": {},
	}
	if _, ok := exact[v]; ok {
		return true
	}
	// Prefix patterns for documented placeholders.
	prefixes := []string{"changeme", "your_", "replace", "todo", "example_", "dummy_", "placeholder"}
	for _, p := range prefixes {
		if strings.HasPrefix(v, p) {
			return true
		}
	}
	return false
}

package taint

import (
	"strings"
	"unicode"
)

// SplitAssignment splits lhs := rhs / lhs = rhs / compound assignment.
func SplitAssignment(text string) (lhs, rhs string, ok bool) {
	text = strings.TrimSpace(text)
	if i := strings.Index(text, ":="); i > 0 {
		return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+2:]), true
	}
	for _, op := range []string{"+=", "-=", "*=", "/=", "%=", "<<=", ">>=", "&^=", "&=", "|=", "^="} {
		if i := strings.Index(text, op); i > 0 {
			return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+len(op):]), true
		}
	}
	if i := strings.Index(text, "="); i > 0 {
		if i+1 < len(text) && text[i+1] == '=' {
			return "", "", false
		}
		// skip !=, <=, >=
		if i > 0 {
			prev := text[i-1]
			if prev == '!' || prev == '<' || prev == '>' {
				return "", "", false
			}
		}
		return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+1:]), true
	}
	return "", "", false
}

// ExtractLHSNames returns comma-separated LHS names (keeps field keys).
func ExtractLHSNames(lhs string) []string {
	parts := strings.Split(lhs, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		name := strings.TrimSpace(p)
		if name == "" || name == "_" {
			continue
		}
		out = append(out, name)
	}
	return out
}

// ResultVariableAtReturnIndex picks the retIdx-th LHS binding (or last).
func ResultVariableAtReturnIndex(lhs string, retIdx int) string {
	vars := strings.Split(lhs, ",")
	for i := range vars {
		vars[i] = strings.TrimSpace(vars[i])
	}
	var pick string
	if retIdx < len(vars) {
		pick = vars[retIdx]
	} else if len(vars) > 0 {
		pick = vars[len(vars)-1]
	}
	if pick != "" && IsIdent(pick) {
		return pick
	}
	return ""
}

// ReferencedNames returns identifiers and field-access chains from an expression.
func ReferencedNames(expr string) []string {
	out := referencedIdentifiers(expr)
	// Also collect ident.field chains.
	bytes := []byte(expr)
	start := 0
	for start < len(bytes) {
		for start < len(bytes) && !(isIdentStart(bytes[start])) {
			start++
		}
		if start >= len(bytes) {
			break
		}
		end := start
		for end < len(bytes) && (isIdentCont(bytes[end]) || bytes[end] == '.') {
			end++
		}
		for end > start && bytes[end-1] == '.' {
			end--
		}
		token := expr[start:end]
		if strings.Contains(token, ".") && token != "" && len(token) < 256 && !containsStr(out, token) {
			out = append(out, token)
		}
		if end <= start {
			start++
		} else {
			start = end
		}
	}
	return out
}

func referencedIdentifiers(expr string) []string {
	var out []string
	var tokenStart = -1
	var quote byte
	escaped := false

	push := func(end int) {
		if tokenStart < 0 {
			return
		}
		token := expr[tokenStart:end]
		tokenStart = -1
		if token == "" || len(token) >= 256 {
			return
		}
		if isAllDigits(token) || isGoKeyword(token) {
			return
		}
		out = append(out, token)
	}

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if quote != 0 {
			switch quote {
			case '`':
				if ch == '`' {
					quote = 0
				}
			case '"', '\'':
				if escaped {
					escaped = false
				} else if ch == '\\' {
					escaped = true
				} else if ch == quote {
					quote = 0
				}
			}
			continue
		}
		switch ch {
		case '"', '\'', '`':
			push(i)
			quote = ch
			escaped = false
		default:
			if isIdentCont(ch) || (tokenStart < 0 && isIdentStart(ch)) {
				if tokenStart < 0 && isIdentStart(ch) {
					tokenStart = i
				} else if tokenStart < 0 {
					// digit-leading: skip
				}
			} else {
				push(i)
			}
		}
	}
	push(len(expr))
	return out
}

func isIdentStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isIdentCont(b byte) bool {
	return isIdentStart(b) || (b >= '0' && b <= '9')
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func containsStr(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

func isGoKeyword(token string) bool {
	switch token {
	case "break", "case", "chan", "const", "continue", "default", "defer", "else",
		"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
		"map", "package", "range", "return", "select", "struct", "switch", "type",
		"var", "string", "int", "bool", "true", "false", "nil", "byte", "error",
		"float64", "float32", "int64", "int32", "uint", "uint64", "rune":
		return true
	}
	return false
}

// IsPureStringLiteral reports whether expr is exactly one Go string literal.
func IsPureStringLiteral(expr string) bool {
	t := strings.TrimSpace(expr)
	if strings.HasPrefix(t, `"`) {
		end := endOfInterpretedString(t)
		return end > 0 && strings.TrimSpace(t[end:]) == ""
	}
	if strings.HasPrefix(t, "`") {
		for i := 1; i < len(t); i++ {
			if t[i] == '`' {
				return strings.TrimSpace(t[i+1:]) == ""
			}
		}
	}
	return false
}

func endOfInterpretedString(s string) int {
	if len(s) == 0 || s[0] != '"' {
		return -1
	}
	i := 1
	for i < len(s) {
		switch s[i] {
		case '\\':
			i += 2
		case '"':
			return i + 1
		default:
			i++
		}
	}
	return -1
}

// NthTopLevelArg returns the n-th top-level argument from call body after '('.
func NthTopLevelArg(body string, n int) (string, bool) {
	depth := 0
	start := 0
	idx := 0
	i := 0
	bytes := []byte(body)
	for i < len(bytes) {
		switch bytes[i] {
		case '"':
			i++
			for i < len(bytes) {
				if bytes[i] == '\\' {
					i += 2
					continue
				}
				if bytes[i] == '"' {
					i++
					break
				}
				i++
			}
			continue
		case '`':
			i++
			for i < len(bytes) {
				if bytes[i] == '`' {
					i++
					break
				}
				i++
			}
			continue
		case '(', '[', '{':
			depth++
			i++
		case ')', ']', '}':
			if depth == 0 {
				if idx != n {
					return "", false
				}
				return strings.TrimSpace(body[start:i]), true
			}
			depth--
			i++
		case ',':
			if depth == 0 {
				if idx == n {
					return strings.TrimSpace(body[start:i]), true
				}
				idx++
				start = i + 1
			}
			i++
		default:
			i++
		}
	}
	if idx == n {
		return strings.TrimSpace(body[start:]), true
	}
	return "", false
}

// ChannelFromReceiveRHS parses `<-ch` → "ch" for simple identifiers.
func ChannelFromReceiveRHS(rhs string) string {
	trimmed := strings.TrimSpace(rhs)
	rest, ok := strings.CutPrefix(trimmed, "<-")
	if !ok {
		return ""
	}
	rest = strings.TrimSpace(rest)
	if rest == "" || !IsIdent(rest) {
		return ""
	}
	return rest
}

// CompactWhitespace removes spaces/tabs/newlines.
func CompactWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

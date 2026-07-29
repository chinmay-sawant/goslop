// Package sourceutil provides lightweight source-text helpers for seed detectors.
package sourceutil

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// RequestSourcePatterns are common request-derived input shapes in Go handlers.
var RequestSourcePatterns = []string{
	"r.URL.Query().Get(",
	"r.FormValue(",
	"r.PostFormValue(",
	"r.Header.Get(",
	"r.PathValue(",
	"r.Cookie(",
	"r.Form.Get(",
	"r.PostForm.Get(",
	"c.Query(",
	"c.Param(",
	"c.PostForm(",
	"c.GetHeader(",
	"c.Cookie(",
	"ctx.Query(",
	"ctx.Param(",
	"mux.Vars(",
}

// assignment is a simple same-line assignment/short-decl.
type assignment struct {
	lhs []string
	rhs string
}

// FindTaintedIdents returns identifiers assigned from request-derived expressions.
// It is a same-file heuristic (not full taint). A few fixed-point passes propagate
// through intermediate locals (e.g. q := fmt.Sprintf(..., name)).
func FindTaintedIdents(source string) map[string]struct{} {
	out := make(map[string]struct{})
	assigns := parseAssignments(source)
	// Seed from direct request sources.
	for _, a := range assigns {
		if looksRequestDerived(a.rhs) {
			for _, name := range a.lhs {
				out[name] = struct{}{}
			}
		}
	}
	// Propagate through locals that reference tainted idents (bounded).
	for pass := 0; pass < 4; pass++ {
		changed := false
		for _, a := range assigns {
			if !rhsTouchesTaint(a.rhs, out) {
				continue
			}
			for _, name := range a.lhs {
				if _, ok := out[name]; ok {
					continue
				}
				out[name] = struct{}{}
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return out
}

func parseAssignments(source string) []assignment {
	var out []assignment
	lines := strings.Split(source, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		// short var decl: name := <expr>
		if idx := strings.Index(trimmed, ":="); idx > 0 {
			lhs := strings.TrimSpace(trimmed[:idx])
			rhs := trimmed[idx+2:]
			out = append(out, assignment{lhs: splitIdents(lhs), rhs: rhs})
			continue
		}
		// assignment: name = <expr>
		if idx := strings.Index(trimmed, "="); idx > 0 {
			// skip ==, !=, <=, >=, :=
			if idx+1 < len(trimmed) {
				prev := trimmed[idx-1]
				next := trimmed[idx+1]
				if prev == '!' || prev == '<' || prev == '>' || prev == '=' || next == '=' {
					continue
				}
			}
			lhs := strings.TrimSpace(trimmed[:idx])
			rhs := trimmed[idx+1:]
			out = append(out, assignment{lhs: splitIdents(lhs), rhs: rhs})
		}
	}
	return out
}

func rhsTouchesTaint(rhs string, tainted map[string]struct{}) bool {
	if looksRequestDerived(rhs) {
		return true
	}
	for name := range tainted {
		if ContainsIdent(rhs, name) {
			return true
		}
	}
	return false
}

func looksRequestDerived(expr string) bool {
	for _, p := range RequestSourcePatterns {
		if strings.Contains(expr, p) {
			return true
		}
	}
	return false
}

// HasRequestSource reports whether source contains any request-derived pattern.
func HasRequestSource(source string) bool {
	for _, p := range RequestSourcePatterns {
		if strings.Contains(source, p) {
			return true
		}
	}
	return false
}

func splitIdents(lhs string) []string {
	// handle "a, b := ..." and "_ = ..."
	parts := strings.Split(lhs, ",")
	var out []string
	for _, p := range parts {
		name := strings.TrimSpace(p)
		if name == "" || name == "_" {
			continue
		}
		// strip type annotations rarely present on short decl lhs
		if i := strings.IndexAny(name, " \t"); i >= 0 {
			name = name[:i]
		}
		if isIdent(name) {
			out = append(out, name)
		}
	}
	return out
}

func isIdent(s string) bool {
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

// CallSite is a lightweight function call match in source text.
type CallSite struct {
	// Name is the callee text matched (e.g. "exec.Command").
	Name string
	// Start is the byte offset of Name in source.
	Start int
	// End is past-the-end of the matching call (closing paren) when found.
	End int
	// ArgsText is the interior of the top-level argument list.
	ArgsText string
}

// FindCalls finds top-level-ish calls to each name in names.
// names should be exact callee strings like "exec.Command" or "strings.Index".
func FindCalls(source string, names ...string) []CallSite {
	var out []CallSite
	for _, name := range names {
		start := 0
		for {
			idx := strings.Index(source[start:], name)
			if idx < 0 {
				break
			}
			abs := start + idx
			// ensure not a longer identifier prefix/suffix
			if abs > 0 {
				r, _ := utf8.DecodeLastRuneInString(source[:abs])
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' {
					start = abs + len(name)
					continue
				}
			}
			after := abs + len(name)
			if after < len(source) {
				r, _ := utf8.DecodeRuneInString(source[after:])
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
					start = after
					continue
				}
			}
			// skip whitespace then require '('
			j := after
			for j < len(source) && (source[j] == ' ' || source[j] == '\t' || source[j] == '\n' || source[j] == '\r') {
				j++
			}
			if j >= len(source) || source[j] != '(' {
				start = after
				continue
			}
			closeAt, args := scanCallArgs(source, j)
			if closeAt < 0 {
				start = after
				continue
			}
			out = append(out, CallSite{
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

// scanCallArgs starts at '(' and returns index of matching ')' and interior args text.
func scanCallArgs(source string, open int) (closeAt int, args string) {
	if open >= len(source) || source[open] != '(' {
		return -1, ""
	}
	depth := 0
	inStr := byte(0) // '"', '`', or '\''
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
			// raw string
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

// SplitTopLevelArgs splits a call argument list on top-level commas.
func SplitTopLevelArgs(args string) []string {
	var out []string
	depth := 0
	inStr := byte(0)
	escape := false
	start := 0
	for i := 0; i < len(args); i++ {
		c := args[i]
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

// IsPureStringLiteral reports whether expr is exactly one Go string literal.
func IsPureStringLiteral(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if t[0] == '"' {
		end := endOfInterpretedString(t)
		return end > 0 && strings.TrimSpace(t[end:]) == ""
	}
	if t[0] == '`' {
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
	escape := false
	for i := 1; i < len(s); i++ {
		if escape {
			escape = false
			continue
		}
		if s[i] == '\\' {
			escape = true
			continue
		}
		if s[i] == '"' {
			return i + 1
		}
	}
	return -1
}

// ContainsIdent reports whether text contains ident as a whole identifier.
func ContainsIdent(text, ident string) bool {
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

// CompactWhitespace removes whitespace for shell-arg shape checks.
func CompactWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

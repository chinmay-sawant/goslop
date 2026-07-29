package badpractices

import (
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-46", detectBP46)
	RegisterRule("BP-47", detectBP47)
	RegisterRule("BP-48", detectBP48)
	RegisterRule("BP-49", detectBP49)
	RegisterRule("BP-50", detectBP50)
	RegisterRule("BP-51", detectBP51)
	RegisterRule("BP-52", detectBP52)
	RegisterRule("BP-53", detectBP53)
	RegisterRule("BP-54", detectBP54)
	RegisterRule("BP-55", detectBP55)
	RegisterRule("BP-56", detectBP56)
}

func detectBP46(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-46")
	if isTestFile(unit) || !strings.Contains(unit.Source, "http.Server") {
		return
	}
	msg := "http.Server should set both ReadTimeout and WriteTimeout"
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "http.Server{")
		if idx < 0 {
			// also &http.Server{
			idx = strings.Index(src[start:], "http.Server{")
			if idx < 0 {
				// try without type name after composite
				break
			}
		}
		abs := start + idx
		end := abs
		depth := 0
		for i := abs; i < len(src); i++ {
			if src[i] == '{' {
				depth++
			} else if src[i] == '}' {
				depth--
				if depth == 0 {
					end = i + 1
					break
				}
			}
		}
		if end > abs {
			lit := src[abs:end]
			if !strings.Contains(lit, "ReadTimeout:") || !strings.Contains(lit, "WriteTimeout:") {
				pushAt(unit, meta, abs, msg, out)
			}
		}
		start = abs + len("http.Server{")
	}
	// also handle &http.Server{ pattern already covered
	// and multi-line with http.Server{\n
	if !strings.Contains(src, "http.Server{") && strings.Contains(src, "http.Server") {
		// var s http.Server; s.Addr= — skip
	}
}

func detectBP47(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-47")
	if isMaterializedFixture(unit) {
		return
	}
	snap := projectSnapshot(unit)
	if !snap.HasServerStart || snap.HasShutdown {
		return
	}
	pushAt(unit, meta, 0, "server startup should include a graceful shutdown path", out)
}

func detectBP48(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-48")
	if isTestFile(unit) {
		return
	}
	msg := "library code should return errors instead of exiting the process"
	callees := []string{"log.Fatal", "log.Fatalf", "log.Fatalln", "os.Exit"}
	if len(facts.callNodes) > 0 {
		for _, c := range facts.callNodes {
			match := false
			for _, name := range callees {
				if c.callee == name {
					match = true
					break
				}
			}
			if !match {
				continue
			}
			fn := facts.enclosingFunc(c.start)
			if fn != nil && (fn.isMain || fn.name == "TestMain") {
				continue
			}
			if name, ok := enclosingFuncName(unit.Source, c.start); ok && (name == "main" || name == "TestMain") {
				continue
			}
			pushAt(unit, meta, c.start, msg, out)
		}
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		for _, name := range callees {
			if strings.Contains(t, name+"(") {
				if n, ok := enclosingFuncName(unit.Source, line.byte); ok && (n == "main" || n == "TestMain") {
					continue
				}
				pushAt(unit, meta, line.byte, msg, out)
			}
		}
	}
}

func detectBP49(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-49")
	if isTestFile(unit) {
		return
	}
	msg := "deferred cleanup drops an error; wrap it in a deferred function and check the result"
	for _, d := range facts.deferNodes {
		text := d.text
		if strings.Contains(text, "func()") {
			continue
		}
		if strings.Contains(text, ".Close()") || strings.Contains(text, ".Flush()") || strings.Contains(text, ".Sync()") {
			pushAt(unit, meta, d.start, msg, out)
		}
	}
	if len(facts.deferNodes) == 0 {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if strings.HasPrefix(t, "defer ") && !strings.Contains(t, "func()") {
				if strings.Contains(t, ".Close()") || strings.Contains(t, ".Flush()") || strings.Contains(t, ".Sync()") {
					pushAt(unit, meta, line.byte, msg, out)
				}
			}
		}
	}
}

func detectBP50(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-50")
	if isMaterializedFixture(unit) {
		return
	}
	snap := projectSnapshot(unit)
	if !snap.HasServerStart || snap.HasSignalHandling {
		return
	}
	pushAt(unit, meta, 0, "long-running server should handle SIGTERM or SIGINT", out)
}

func detectBP51(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-51")
	if isTestFile(unit) || packageName(unit.Source) == "main" {
		return
	}
	if !strings.Contains(unit.Source, "recover()") {
		return
	}
	if strings.Contains(unit.Source, "panic(") || strings.Contains(unit.Source, "log.") || strings.Contains(unit.Source, "fmt.") {
		return
	}
	if pos := strings.Index(unit.Source, "recover()"); pos >= 0 {
		pushAt(unit, meta, pos, "library recover paths should re-panic or convert the panic into an explicit error contract", out)
	}
}

func detectBP52(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-52")
	if isTestFile(unit) || !strings.Contains(unit.Source, "make(") || !strings.Contains(unit.Source, "*") {
		return
	}
	// Rust parity (production_hardening.rs): every make(...) call whose text
	// contains '*' and whose enclosing func/func-literal lacks an overflow guard.
	// Guards are scope-local (not file-wide) and include "/' (naive, matches Rust).
	src := unit.Source
	for _, call := range findMakeCalls(src) {
		if !strings.Contains(call.text, "*") {
			continue
		}
		scope := enclosingFuncOrLiteralScope(src, call.start)
		if scope == "" {
			// Package-level make without enclosing func — Rust skips these.
			continue
		}
		if bp52HasOverflowGuard(scope) {
			continue
		}
		pushAt(unit, meta, call.start, "multiplication used in an allocation path without an obvious overflow guard", out)
	}
}

type makeCall struct {
	start int
	text  string
}

func findMakeCalls(src string) []makeCall {
	var out []makeCall
	start := 0
	for {
		idx := strings.Index(src[start:], "make(")
		if idx < 0 {
			break
		}
		abs := start + idx
		// avoid longer idents like remake(
		if abs > 0 {
			prev := src[abs-1]
			if (prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z') || (prev >= '0' && prev <= '9') || prev == '_' {
				start = abs + 4
				continue
			}
		}
		open := abs + 4 // "make" is 4 chars; '(' follows
		if open >= len(src) || src[open] != '(' {
			start = abs + 4
			continue
		}
		closeAt := scanBalancedCallEnd(src, open)
		if closeAt < 0 {
			start = abs + 4
			continue
		}
		out = append(out, makeCall{start: abs, text: src[abs : closeAt+1]})
		start = closeAt + 1
	}
	return out
}

// scanBalancedCallEnd starts at '(' and returns index of matching ')'.
func scanBalancedCallEnd(src string, open int) int {
	if open >= len(src) || src[open] != '(' {
		return -1
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for i := open; i < len(src); i++ {
		c := src[i]
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
				return i
			}
		}
	}
	return -1
}

// enclosingFuncOrLiteralScope returns the source span of the innermost
// func/func-literal whose body braces enclose pos. Empty if none.
func enclosingFuncOrLiteralScope(src string, pos int) string {
	type cand struct {
		start int
		brace int
	}
	var cands []cand
	// Find "func " / "func(" occurrences before pos.
	search := 0
	for search < pos {
		idx := strings.Index(src[search:], "func")
		if idx < 0 {
			break
		}
		abs := search + idx
		if abs >= pos {
			break
		}
		// word boundary after "func"
		end := abs + 4
		if end < len(src) {
			n := src[end]
			if (n >= 'a' && n <= 'z') || (n >= 'A' && n <= 'Z') || (n >= '0' && n <= '9') || n == '_' {
				search = end
				continue
			}
		}
		if abs > 0 {
			p := src[abs-1]
			if (p >= 'a' && p <= 'z') || (p >= 'A' && p <= 'Z') || (p >= '0' && p <= '9') || p == '_' {
				search = end
				continue
			}
		}
		// find opening brace of this func
		brace := -1
		inStr := byte(0)
		escape := false
		for i := end; i < len(src) && i < pos+1; i++ {
			c := src[i]
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
			case '{':
				brace = i
			}
			if brace >= 0 {
				break
			}
			// signature without body (interface) — skip
			if c == '\n' && !strings.Contains(src[end:i], "(") {
				// still could be multi-line; continue
			}
		}
		if brace >= 0 && brace < pos {
			cands = append(cands, cand{start: abs, brace: brace})
		}
		search = end
	}
	// Innermost cand whose braces enclose pos.
	for i := len(cands) - 1; i >= 0; i-- {
		c := cands[i]
		end := matchBraceEnd(src, c.brace)
		if end >= pos {
			return src[c.start : end+1]
		}
	}
	return ""
}

func matchBraceEnd(src string, openBrace int) int {
	if openBrace < 0 || openBrace >= len(src) || src[openBrace] != '{' {
		return -1
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for i := openBrace; i < len(src); i++ {
		c := src[i]
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
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func bp52HasOverflowGuard(scope string) bool {
	// Match Rust has_guard set exactly.
	return strings.Contains(scope, "MaxInt") ||
		strings.Contains(scope, "MaxUint") ||
		strings.Contains(scope, "overflow") ||
		strings.Contains(scope, "/") ||
		strings.Contains(scope, "bits.Mul") ||
		strings.Contains(scope, "checkedMul") ||
		strings.Contains(scope, "checked_mul")
}

func detectBP53(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-53")
	src := unit.Source
	if isTestFile(unit) || !strings.Contains(src, "gob.Register(") {
		return
	}
	msg := "gob.Register uses a type that does not line up with the nearby Encode/Decode payloads"
	registered := collectCallTargets(src, "gob.Register(")
	if len(registered) == 0 {
		return
	}
	encoded := collectCallTargets(src, ".Encode(")
	decoded := collectCallTargets(src, ".Decode(")
	if len(encoded) == 0 && len(decoded) == 0 {
		if pos := strings.Index(src, "gob.Register("); pos >= 0 {
			pushAt(unit, meta, pos, msg, out)
		}
		return
	}
	knownTypes := collectLocalTypeHints(src)
	matched := false
	for _, value := range append(append([]string{}, encoded...), decoded...) {
		normalized := normalizeIdentifier(value)
		ty, ok := knownTypes[normalized]
		if !ok {
			// Encode may pass a composite literal type name directly
			ty = normalizeTypeName(value)
		}
		for _, candidate := range registered {
			if normalizeTypeName(candidate) == ty {
				matched = true
				break
			}
		}
		if matched {
			break
		}
	}
	if !matched {
		if pos := strings.Index(src, "gob.Register("); pos >= 0 {
			pushAt(unit, meta, pos, msg, out)
		}
	}
}

func collectCallTargets(source, needle string) []string {
	var values []string
	start := 0
	for {
		offset := strings.Index(source[start:], needle)
		if offset < 0 {
			break
		}
		idx := start + offset + len(needle)
		end := strings.IndexByte(source[idx:], ')')
		if end < 0 {
			break
		}
		values = append(values, strings.TrimSpace(source[idx:idx+end]))
		start = idx + end + 1
	}
	return values
}

func collectLocalTypeHints(source string) map[string]string {
	types := map[string]string{}
	for _, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, ":=") {
			parts := strings.SplitN(trimmed, ":=", 2)
			if len(parts) != 2 {
				continue
			}
			name := strings.TrimSpace(parts[0])
			rhs := strings.TrimSpace(parts[1])
			if i := strings.IndexByte(rhs, '{'); i >= 0 {
				ident := strings.TrimSpace(rhs[:i])
				ident = strings.TrimPrefix(strings.TrimPrefix(ident, "&"), "*")
				if ident != "" && isSimpleIdent(ident) {
					types[name] = ident
				}
			}
		} else if strings.HasPrefix(trimmed, "var ") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "var "))
			fields := strings.Fields(rest)
			if len(fields) >= 2 {
				name := fields[0]
				ty := strings.TrimPrefix(fields[1], "*")
				types[name] = ty
			}
		}
	}
	return types
}

func normalizeIdentifier(value string) string {
	v := strings.TrimSpace(value)
	v = strings.TrimPrefix(v, "&")
	v = strings.TrimPrefix(v, "*")
	v = strings.Trim(v, "()")
	return v
}

func normalizeTypeName(value string) string {
	value = strings.TrimSpace(value)
	if i := strings.IndexByte(value, '{'); i >= 0 {
		return normalizeIdentifier(value[:i])
	}
	return normalizeIdentifier(value)
}

func detectBP54(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-54")
	if isMaterializedFixture(unit) {
		return
	}
	snap := projectSnapshot(unit)
	if !snap.HasServerStart || !snap.HasPublicRoute || snap.HasRateLimiting {
		return
	}
	pushAt(unit, meta, 0, "public HTTP handlers should enforce a rate-limiting guard", out)
}

func detectBP55(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-55")
	if isMaterializedFixture(unit) {
		return
	}
	snap := projectSnapshot(unit)
	if !snap.HasServerStart || !snap.HasPublicRoute || !snap.HasLogging || snap.HasRequestID {
		return
	}
	pushAt(unit, meta, 0, "request-handling code logs traffic without a visible request-id propagation path", out)
}

func detectBP56(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-56")
	deprecated := []string{
		`"io/ioutil"`, `"golang.org/x/net/context"`, `"github.com/golang/protobuf"`,
	}
	for _, d := range deprecated {
		if strings.Contains(unit.Source, d) {
			if pos := strings.Index(unit.Source, d); pos >= 0 {
				pushAt(unit, meta, pos, "deprecated package import; migrate to the supported replacement", out)
			}
		}
	}
}

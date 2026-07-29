package badpractices

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("BP-26", detectBP26)
	RegisterRule("BP-27", detectBP27)
	RegisterRule("BP-28", detectBP28)
	RegisterRule("BP-29", detectBP29)
	RegisterRule("BP-30", detectBP30)
	RegisterRule("BP-31", detectBP31)
	RegisterRule("BP-32", detectBP32)
	RegisterRule("BP-34", detectBP34)
	RegisterRule("BP-36", detectBP36)
	RegisterRule("BP-37", detectBP37)
	RegisterRule("BP-38", detectBP38)
	RegisterRule("BP-39", detectBP39)
	RegisterRule("BP-41", detectBP41)
	RegisterRule("BP-42", detectBP42)
	RegisterRule("BP-43", detectBP43)
	RegisterRule("BP-44", detectBP44)
	RegisterRule("BP-66", detectBP66)
	RegisterRule("BP-76", detectBP76)
	RegisterRule("BP-85", detectBP85)
	RegisterRule("BP-91", detectBP91)
	RegisterRule("BP-138", detectBP138)
	RegisterRule("BP-141", detectBP141)
	RegisterRule("BP-164", detectBP164)
}

func detectBP26(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-26")
	if !strings.Contains(unit.Source, "context.Context") {
		return
	}
	msg := "context.Context should be the first parameter"
	for _, fn := range facts.funcDecls {
		p := fn.params
		if !strings.Contains(p, "context.Context") {
			continue
		}
		// strip outer parens
		inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(p, "("), ")"))
		if inner == "" {
			continue
		}
		// first param should contain context.Context
		first := strings.Split(inner, ",")[0]
		if !strings.Contains(first, "context.Context") {
			pushAt(unit, meta, fn.start, msg, out)
		}
	}
	if len(facts.funcDecls) == 0 {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if !strings.HasPrefix(t, "func ") || !strings.Contains(t, "context.Context") {
				continue
			}
			// find param list
			open := strings.Index(t, "(")
			close := strings.LastIndex(t, ")")
			if open < 0 || close <= open {
				continue
			}
			// may be method: func (r *T) Name(
			rest := t[open+1:]
			if strings.HasPrefix(strings.TrimSpace(rest), "*") || strings.Contains(rest[:min(20, len(rest))], " ") && !strings.Contains(rest[:min(30, len(rest))], "context") {
				// method receiver — find second (
				if second := strings.Index(t[open+1:], "("); second >= 0 {
					open = open + 1 + second
					rest = t[open+1:]
					close = strings.Index(rest, ")")
					if close < 0 {
						continue
					}
					inner := rest[:close]
					first := strings.Split(inner, ",")[0]
					if strings.Contains(inner, "context.Context") && !strings.Contains(first, "context.Context") {
						pushAt(unit, meta, line.byte, msg, out)
					}
				}
				continue
			}
			close = strings.Index(rest, ")")
			if close < 0 {
				continue
			}
			inner := rest[:close]
			if !strings.Contains(inner, "context.Context") {
				continue
			}
			first := strings.Split(inner, ",")[0]
			if !strings.Contains(first, "context.Context") {
				pushAt(unit, meta, line.byte, msg, out)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func detectBP27(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-27")
	// exported func returns unexported type
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "func ") {
			continue
		}
		// func Name(...) unexported
		nameStart := len("func ")
		rest := t[nameStart:]
		if strings.HasPrefix(rest, "(") {
			// method
			continue
		}
		end := 0
		for end < len(rest) && (unicode.IsLetter(rune(rest[end])) || unicode.IsDigit(rune(rest[end])) || rest[end] == '_') {
			end++
		}
		if end == 0 {
			continue
		}
		name := rest[:end]
		if name == "" || !unicode.IsUpper(rune(name[0])) {
			continue
		}
		// look at return types after )
		if i := strings.LastIndex(t, ")"); i >= 0 {
			rets := strings.TrimSpace(t[i+1:])
			rets = strings.TrimPrefix(rets, "(")
			rets = strings.TrimSuffix(rets, "{")
			rets = strings.TrimSpace(strings.TrimSuffix(rets, ")"))
			for _, r := range strings.Split(rets, ",") {
				r = strings.TrimSpace(r)
				r = strings.TrimPrefix(r, "*")
				if r == "" || r == "error" || r == "string" || r == "int" || r == "bool" || r == "byte" {
					continue
				}
				if unicode.IsLetter(rune(r[0])) && unicode.IsLower(rune(r[0])) && !strings.Contains(r, ".") {
					pushAt(unit, meta, line.byte, "exported function returns an unexported type", out)
					break
				}
			}
		}
	}
}

func detectBP28(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-28")
	// single-method interface
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "interface {")
		if idx < 0 {
			idx = strings.Index(src[start:], "interface{")
			if idx < 0 {
				break
			}
		}
		abs := start + idx
		// count methods roughly
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
		block := src[abs:end]
		// method lines: name(
		methods := 0
		for _, line := range strings.Split(block, "\n") {
			t := strings.TrimSpace(line)
			if t == "" || strings.HasPrefix(t, "//") || t == "interface {" || t == "interface{" || t == "}" {
				continue
			}
			if strings.Contains(t, "(") {
				methods++
			}
		}
		if methods == 1 {
			pushAt(unit, meta, abs, "single-method interface; consider accepting a function type instead", out)
		}
		start = end
	}
}

func detectBP29(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-29")
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "interface {")
		if idx < 0 {
			break
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
		block := src[abs:end]
		methods := 0
		for _, line := range strings.Split(block, "\n") {
			t := strings.TrimSpace(line)
			if strings.Contains(t, "(") && !strings.HasPrefix(t, "interface") {
				methods++
			}
		}
		// Rust: more than five methods.
		if methods > 5 {
			pushAt(unit, meta, abs, "interface declares more than five methods and is likely too broad", out)
		}
		start = end
	}
}

// detectBP30: exported interface with no evident same-package implementation.
// Uses the whole package directory (Rust package_method_sets parity).
func detectBP30(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	if !strings.Contains(unit.Source, "interface") {
		return
	}
	meta := MetadataForID("BP-30")
	src := unit.Source
	pkgFacts := packageTypeFactsForUnit(unit)

	// Only report interfaces declared in this file (per-file findings).
	for _, m := range reExportedIface.FindAllStringSubmatchIndex(src, -1) {
		ifaceName := src[m[2]:m[3]]
		methods, ok := pkgFacts.interfaces[ifaceName]
		if !ok || len(methods) == 0 {
			// re-parse local body if package scan missed embedding noise
			continue
		}
		hasImpl := false
		for typ := range pkgFacts.methods {
			if typ == ifaceName {
				continue
			}
			if typeImplements(pkgFacts, typ, methods) {
				hasImpl = true
				break
			}
		}
		if !hasImpl {
			pushAt(unit, meta, m[0], "exported interface has no evident same-package implementation", out)
		}
	}
}

// detectBP31: New* constructor returns a concrete type even though a fitting
// package interface already exists (Rust api_design::detect_bp_31).
func detectBP31(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	if !strings.Contains(unit.Source, "func New") {
		return
	}
	meta := MetadataForID("BP-31")
	pkgFacts := packageTypeFactsForUnit(unit)
	if len(pkgFacts.interfaces) == 0 {
		return
	}
	lines := strings.Split(unit.Source, "\n")
	byteOff := 0
	for _, line := range lines {
		t := strings.TrimSpace(line)
		// package-level constructors only (not methods)
		if strings.HasPrefix(t, "func New") || strings.HasPrefix(t, "func Must") {
			// func NewX(...) *Type / Type
			if strings.HasPrefix(t, "func (") {
				byteOff += len(line) + 1
				continue
			}
			name := firstIdent(strings.TrimPrefix(t, "func "))
			if name == "" || !unicode.IsUpper(rune(name[0])) {
				byteOff += len(line) + 1
				continue
			}
			if !strings.HasPrefix(name, "New") && !strings.HasPrefix(name, "Must") {
				byteOff += len(line) + 1
				continue
			}
			returned := firstResultTypeFromFuncLine(t)
			if returned == "" || !unicode.IsUpper(rune(returned[0])) {
				byteOff += len(line) + 1
				continue
			}
			methods := pkgFacts.methods[returned]
			if len(methods) == 0 {
				byteOff += len(line) + 1
				continue
			}
			exposes := false
			for ifaceName, ifaceMethods := range pkgFacts.interfaces {
				if ifaceName == returned || len(ifaceMethods) == 0 {
					continue
				}
				if typeImplements(pkgFacts, returned, ifaceMethods) {
					exposes = true
					break
				}
			}
			if exposes {
				pushAt(unit, meta, byteOff,
					"constructor returns a concrete type even though the package already exposes a fitting interface",
					out)
			}
		}
		byteOff += len(line) + 1
	}
}

// firstResultTypeFromFuncLine extracts the first named result type from a single-line
// func signature (best-effort; multi-line signatures are skipped).
func firstResultTypeFromFuncLine(line string) string {
	// func Name(...) results {
	closeParams := strings.Index(line, ")")
	if closeParams < 0 {
		return ""
	}
	rest := strings.TrimSpace(line[closeParams+1:])
	if rest == "" || strings.HasPrefix(rest, "{") {
		return ""
	}
	// strip trailing {
	if i := strings.Index(rest, "{"); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	// (T, error) or *T or (a T, err error)
	rest = strings.TrimSpace(rest)
	if strings.HasPrefix(rest, "(") {
		inner := strings.TrimSuffix(strings.TrimPrefix(rest, "("), ")")
		parts := strings.Split(inner, ",")
		if len(parts) == 0 {
			return ""
		}
		return normalizeResultType(parts[0])
	}
	// single result
	if i := strings.IndexAny(rest, " \t"); i > 0 {
		rest = rest[:i]
	}
	return normalizeResultType(rest)
}

func normalizeResultType(s string) string {
	s = strings.TrimSpace(s)
	// named result: "err error" / "v *Type"
	fields := strings.Fields(s)
	if len(fields) >= 2 {
		s = fields[len(fields)-1]
	}
	s = strings.TrimPrefix(s, "*")
	// drop package qualifier
	if i := strings.LastIndex(s, "."); i >= 0 {
		s = s[i+1:]
	}
	// only plain identifiers
	if s == "" || !unicode.IsLetter(rune(s[0])) {
		return ""
	}
	return s
}

func detectBP32(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-32")
	// type X string with Error() method — string alias error
	if strings.Contains(unit.Source, "type ") && strings.Contains(unit.Source, " string") && strings.Contains(unit.Source, "Error() string") {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if strings.HasPrefix(t, "type ") && strings.HasSuffix(t, " string") {
				pushAt(unit, meta, line.byte, "string alias used as error type; prefer a struct implementing error", out)
				return
			}
		}
	}
}

func detectBP34(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-34")
	// fmt.Errorf without %w when wrapping err — require bare `err` arg, not errCount.
	emitted := 0
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.Contains(t, "fmt.Errorf(") || strings.Contains(t, "%w") {
			continue
		}
		if !bp34WrapsErr(t) {
			continue
		}
		pushAt(unit, meta, line.byte, "error wrapping without %w loses the error chain", out)
		emitted++
		if emitted >= 2 {
			return
		}
	}
}

// bp34WrapsErr reports whether fmt.Errorf call arguments include the bare
// identifier `err` (word boundary), not errCount / errno / etc.
func bp34WrapsErr(line string) bool {
	// After fmt.Errorf(, scan for , err or (err as standalone ident.
	idx := strings.Index(line, "fmt.Errorf(")
	if idx < 0 {
		return false
	}
	args := line[idx+len("fmt.Errorf("):]
	// Strip trailing comments.
	if i := strings.Index(args, "//"); i >= 0 {
		args = args[:i]
	}
	for i := 0; i < len(args); {
		// find "err"
		j := strings.Index(args[i:], "err")
		if j < 0 {
			return false
		}
		j += i
		// word start
		if j > 0 {
			prev := args[j-1]
			if (prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z') ||
				(prev >= '0' && prev <= '9') || prev == '_' {
				i = j + 3
				continue
			}
		}
		// word end
		end := j + 3
		if end < len(args) {
			next := args[end]
			if (next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') ||
				(next >= '0' && next <= '9') || next == '_' {
				i = j + 3
				continue
			}
		}
		return true
	}
	return false
}

func detectBP36(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust: init body contains call_expression / go / defer.
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-36")
	pos := strings.Index(unit.Source, "func init(")
	if pos < 0 {
		return
	}
	body := unit.Source[pos:]
	if end := strings.Index(body, "\nfunc "); end > 0 {
		body = body[:end]
	}
	// Extract brace body
	open := strings.Index(body, "{")
	if open < 0 {
		return
	}
	depth := 0
	end := open
	for i := open; i < len(body); i++ {
		if body[i] == '{' {
			depth++
		} else if body[i] == '}' {
			depth--
			if depth == 0 {
				end = i
				break
			}
		}
	}
	inner := body[open+1 : end]
	// Side effect: any call/conversion (contains '(' not in string — crude) or go/defer.
	if strings.Contains(inner, "go ") || strings.Contains(inner, "defer ") ||
		strings.Contains(inner, "(") {
		pushAt(unit, meta, pos, "init() performs side effects beyond simple package setup", out)
	}
}

func detectBP37(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust parity: package-level var that is later written (assignment/index).
	// Oracle-tight: only arrays written via name[i]=, or pointer/map caches.
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-37")
	globals := packageLevelVarNames(unit.Source)
	if len(globals) == 0 {
		return
	}
	written := packageLevelVarsWritten(unit.Source, globals)
	if len(written) == 0 {
		return
	}
	// Tight oracle-aligned filter: only (1) fixed arrays written via index, or
	// (2) pointer-to-struct package caches. Map/slice globals over-fire vs Rust.
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "var ") {
			continue
		}
		for name := range written {
			if !strings.Contains(t, name) || strings.HasPrefix(name, "Err") {
				continue
			}
			if strings.Contains(t, "sync.Pool") || strings.Contains(t, "map[") {
				continue
			}
			// skip slice literals: var x = []T{...}
			if strings.Contains(t, "[]") {
				continue
			}
			isArray := false
			// var name [N]T
			rest := strings.TrimSpace(strings.TrimPrefix(t, "var "))
			if strings.HasPrefix(rest, name) {
				after := strings.TrimSpace(rest[len(name):])
				if strings.HasPrefix(after, "[") && !strings.HasPrefix(after, "[]") {
					isArray = true
				}
			}
			isPtrCache := strings.Contains(t, " = &") || strings.Contains(t, "=&")
			if !isArray && !isPtrCache {
				continue
			}
			pushAt(unit, meta, line.byte, "package-level mutable global state makes behavior harder to reason about", out)
			return // one per file
		}
	}
}

func packageLevelVarNames(source string) map[string]struct{} {
	out := map[string]struct{}{}
	inFunc := 0
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		// crude brace depth for top-level only
		open := strings.Count(t, "{")
		closeN := strings.Count(t, "}")
		if inFunc == 0 && strings.HasPrefix(t, "var ") {
			// var name T or var ( ... )
			rest := strings.TrimSpace(strings.TrimPrefix(t, "var "))
			if strings.HasPrefix(rest, "(") {
				// multi-line block handled poorly — skip paren form for now
			} else {
				name := firstIdent(rest)
				if name != "" && !strings.HasPrefix(name, "Err") {
					out[name] = struct{}{}
				}
			}
		}
		inFunc += open - closeN
		if inFunc < 0 {
			inFunc = 0
		}
	}
	return out
}

func packageLevelVarsWritten(source string, globals map[string]struct{}) map[string]struct{} {
	written := map[string]struct{}{}
	// After first function starts, look for clear writes on the global or its fields.
	// Do NOT treat name.Method() as a write.
	inFunc := false
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "func ") {
			inFunc = true
		}
		if !inFunc {
			continue
		}
		for name := range globals {
			if strings.Contains(t, name+" :=") || strings.Contains(t, name+":=") {
				continue // shadowing short decl
			}
			// name = / name += / name[
			if strings.Contains(t, name+" =") || strings.Contains(t, name+"=") ||
				strings.Contains(t, name+"[") || strings.Contains(t, name+" +=") ||
				strings.Contains(t, name+"+=") {
				// Exclude method calls: name.Foo( without assignment after
				if isMethodCallOnly(t, name) {
					continue
				}
				written[name] = struct{}{}
				continue
			}
			// field write: name.field = or name.field[
			if strings.Contains(t, name+".") {
				rest := t[strings.Index(t, name+".")+len(name)+1:]
				if strings.Contains(rest, "=") || strings.Contains(rest, "[") {
					// still exclude name.Method()
					dot := strings.Index(t, name+".")
					after := t[dot+len(name)+1:]
					identEnd := 0
					for identEnd < len(after) && (isIdentByteBP(after[identEnd])) {
						identEnd++
					}
					if identEnd < len(after) && after[identEnd] == '(' {
						continue
					}
					written[name] = struct{}{}
				}
			}
		}
	}
	return written
}

func isMethodCallOnly(line, name string) bool {
	// true if the only uses of name are name.Method( forms without assignment
	if strings.Contains(line, name+" =") || strings.Contains(line, name+"=") ||
		strings.Contains(line, name+"[") {
		return false
	}
	return strings.Contains(line, name+".")
}

func isIdentByteBP(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func firstIdent(s string) string {
	s = strings.TrimSpace(s)
	i := 0
	for i < len(s) {
		r := rune(s[i])
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return ""
			}
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			break
		}
		i++
	}
	return s[:i]
}

func detectBP38(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Unexported helper with no same-file callers.
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-38")
	src := unit.Source
	helpers := unexportedHelpers(src)
	if len(helpers) == 0 {
		return
	}
	calls := localCallNames(src)
	for name, byteOff := range helpers {
		if _, used := calls[name]; used {
			continue
		}
		// skip common generated / interface method noise
		if name == "init" || name == "main" {
			continue
		}
		pushAt(unit, meta, byteOff, "unexported helper has no same-file callers", out)
	}
}

func unexportedHelpers(source string) map[string]int {
	// Package-scope funcs and methods (Rust collects both function_declaration
	// and method_declaration). Name filter: helper/must/build/parse*.
	out := map[string]int{}
	byteOff := 0
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "func ") {
			rest := strings.TrimPrefix(t, "func ")
			// methods: func (r *T) name(
			if strings.HasPrefix(rest, "(") {
				if close := strings.Index(rest, ")"); close >= 0 {
					rest = strings.TrimSpace(rest[close+1:])
				} else {
					byteOff += len(line) + 1
					continue
				}
			}
			name := firstIdent(rest)
			if name != "" && unicode.IsLower(rune(name[0])) && looksLikeHelperName(name) {
				if name != "init" && name != "main" &&
					!strings.HasPrefix(name, "Test") &&
					!strings.HasPrefix(name, "Benchmark") &&
					!strings.HasPrefix(name, "Example") {
					// first declaration wins for multi-line
					if _, ok := out[name]; !ok {
						out[name] = byteOff
					}
				}
			}
		}
		byteOff += len(line) + 1
	}
	return out
}

func looksLikeHelperName(name string) bool {
	// Rust parity: only helper/must/build/parse prefixes.
	lower := strings.ToLower(name)
	return lower == "helper" ||
		strings.HasPrefix(lower, "helper") ||
		strings.HasPrefix(lower, "must") ||
		strings.HasPrefix(lower, "build") ||
		strings.HasPrefix(lower, "parse")
}

func localCallNames(source string) map[string]struct{} {
	out := map[string]struct{}{}
	// crude: bare name( not preceded by ., excluding func/method declarations.
	for i := 0; i < len(source); i++ {
		if source[i] != '(' {
			continue
		}
		// walk back over ident
		j := i - 1
		for j >= 0 && isIdentByte(source[j]) {
			j--
		}
		name := source[j+1 : i]
		if name == "" {
			continue
		}
		// skip method/selector call: x.Name(
		if j >= 0 && source[j] == '.' {
			continue
		}
		// skip function/method declarations:
		//   func name(
		//   func (recv T) name(
		k := j
		for k >= 0 && (source[k] == ' ' || source[k] == '\t') {
			k--
		}
		if k >= 0 && source[k] == ')' {
			// method declaration receiver — skip back over "(...)" then look for func
			depth := 1
			k--
			for k >= 0 && depth > 0 {
				switch source[k] {
				case ')':
					depth++
				case '(':
					depth--
				}
				k--
			}
			for k >= 0 && (source[k] == ' ' || source[k] == '\t') {
				k--
			}
		}
		if k >= 3 && source[k-3:k+1] == "func" {
			if k-3 == 0 || !isIdentByte(source[k-4]) {
				continue // declaration, not a call
			}
		}
		if unicode.IsLower(rune(name[0])) {
			out[name] = struct{}{}
		}
	}
	return out
}

func detectBP39(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust: exported API (funcs + methods on exported receivers) need a doc
	// comment that starts with the function/method name. Emit all hits (no
	// early return) for real-repos parity.
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-39")
	lines := strings.Split(unit.Source, "\n")
	byteOff := 0
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if !strings.HasPrefix(t, "func ") {
			byteOff += len(line) + 1
			continue
		}
		rest := strings.TrimPrefix(t, "func ")
		name := ""
		if strings.HasPrefix(rest, "(") {
			// method: func (r *Type) Name
			closeParen := strings.Index(rest, ")")
			if closeParen < 0 {
				byteOff += len(line) + 1
				continue
			}
			recv := strings.TrimSpace(rest[:closeParen+1]) // includes leading (
			recvInner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(recv, "("), ")"))
			fields := strings.Fields(recvInner)
			if len(fields) == 0 {
				byteOff += len(line) + 1
				continue
			}
			recvType := strings.TrimPrefix(fields[len(fields)-1], "*")
			// Rust is_exported_api: only methods on exported receivers.
			if recvType == "" || !unicode.IsUpper(rune(recvType[0])) {
				byteOff += len(line) + 1
				continue
			}
			after := strings.TrimSpace(rest[closeParen+1:])
			name = firstIdent(after)
		} else {
			name = firstIdent(rest)
		}
		if name == "" || !unicode.IsUpper(rune(name[0])) {
			byteOff += len(line) + 1
			continue
		}
		// Collect contiguous // comments immediately above (blank lines stop after docs start).
		var docs []string
		for j := i - 1; j >= 0; j-- {
			pt := strings.TrimSpace(lines[j])
			if pt == "" {
				break // Rust breaks on empty (after reading upward)
			}
			if strings.HasPrefix(pt, "//") {
				docs = append([]string{strings.TrimSpace(strings.TrimPrefix(pt, "//"))}, docs...)
				continue
			}
			break
		}
		ok := len(docs) > 0 && strings.HasPrefix(docs[0], name)
		if !ok {
			pushAt(unit, meta, byteOff, "exported API should have a doc comment that starts with its name", out)
		}
		byteOff += len(line) + 1
	}
}

// detectBP42: import alias used only once (Rust count_word_occurrences <= 2).
func detectBP42(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-42")
	src := unit.Source
	// Single-line: import alias "path"  (not import "path")
	reSingle := regexp.MustCompile(`(?m)^import\s+([A-Za-z_]\w*)\s+"[^"]+"`)
	// Inside import ( ): alias "path"  (not bare "path")
	reBlock := regexp.MustCompile(`(?m)^\s*([A-Za-z_]\w*)\s+"[^"]+"`)
	inImport := false
	byteOff := 0
	for _, line := range strings.Split(src, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "import ") || t == "import (" || strings.HasPrefix(t, "import(") {
			if strings.Contains(t, "(") {
				inImport = true
				// also handle: import ( alias "path" ) rare one-liners
			} else if m := reSingle.FindStringSubmatchIndex(line); m != nil {
				alias := line[m[2]:m[3]]
				if isUsefulImportAlias(alias) && countWordOccurrences(src, alias) <= 2 {
					pushAt(unit, meta, byteOff+m[2],
						"import alias is only used once and likely adds indirection without value", out)
				}
			}
			byteOff += len(line) + 1
			continue
		}
		if inImport {
			if t == ")" {
				inImport = false
				byteOff += len(line) + 1
				continue
			}
			// skip bare "path" imports (no alias)
			if strings.HasPrefix(t, "\"") || strings.HasPrefix(t, ". \"") || strings.HasPrefix(t, "_ \"") {
				byteOff += len(line) + 1
				continue
			}
			if m := reBlock.FindStringSubmatchIndex(line); m != nil {
				alias := line[m[2]:m[3]]
				if isUsefulImportAlias(alias) && countWordOccurrences(src, alias) <= 2 {
					pushAt(unit, meta, byteOff+m[2],
						"import alias is only used once and likely adds indirection without value", out)
				}
			}
		}
		byteOff += len(line) + 1
	}
}

func isUsefulImportAlias(alias string) bool {
	if alias == "" || alias == "_" || alias == "." {
		return false
	}
	// never treat the import keyword / common false parses as an alias
	switch alias {
	case "import", "package", "var", "const", "type", "func", "go", "defer", "map", "chan", "interface", "struct":
		return false
	}
	return true
}

func detectBP41(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust parity: only the package anchor file reports missing package doc.
	if isTestFile(unit) || isMaterializedFixture(unit) {
		return
	}
	meta := MetadataForID("BP-41")
	pkg := packageName(unit.Source)
	if pkg == "" {
		return
	}
	snap := packageDocSnapshotForUnit(unit)
	anchor, ok := snap.Anchors[pkg]
	if !ok {
		return
	}
	// Compare absolute paths when possible.
	unitPath := unit.Path
	if abs, err := filepath.Abs(unitPath); err == nil {
		unitPath = abs
	}
	if abs, err := filepath.Abs(anchor); err == nil {
		anchor = abs
	}
	if unitPath != anchor {
		return
	}
	if _, documented := snap.DocumentedPackages[pkg]; documented {
		return
	}
	pushAt(unit, meta, 0, "package is missing a package-level doc comment", out)
}

func detectBP43(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-43")
	if isTestFile(unit) {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, ". \"") || strings.Contains(t, "\t. \"") || strings.Contains(t, " . \"") {
			pushAt(unit, meta, line.byte, "dot import outside tests pollutes the namespace", out)
		}
		if t == `. "` || strings.HasPrefix(t, `. "`) {
			pushAt(unit, meta, line.byte, "dot import outside tests pollutes the namespace", out)
		}
	}
	// import ( . "fmt" )
	if strings.Contains(unit.Source, `. "`) {
		if pos := strings.Index(unit.Source, `. "`); pos >= 0 {
			// avoid false positive in strings
			pushAt(unit, meta, pos, "dot import outside tests pollutes the namespace", out)
		}
	}
}

func detectBP44(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-44")
	// blank import without comment
	lines := strings.Split(unit.Source, "\n")
	byteOff := 0
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "_ \"") || strings.HasPrefix(t, "_ \"") {
			// comment on same line or previous?
			if !strings.Contains(line, "//") {
				prevOK := false
				if i > 0 && strings.Contains(lines[i-1], "//") {
					prevOK = true
				}
				if !prevOK {
					pushAt(unit, meta, byteOff, "blank import without a justifying comment", out)
				}
			}
		}
		byteOff += len(line) + 1
	}
}

func detectBP66(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-66")
	// err == SomeSentinel when wrapping used
	if strings.Contains(unit.Source, "fmt.Errorf") && strings.Contains(unit.Source, "%w") {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if (strings.Contains(t, "err ==") || strings.Contains(t, "err !=")) && strings.Contains(t, "Err") {
				if !strings.Contains(t, "errors.Is") {
					pushAt(unit, meta, line.byte, "wrapped error compared directly; use errors.Is", out)
				}
			}
		}
	}
}

func detectBP76(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-76")
	// for k, v := range map used to build ordered output
	if strings.Contains(unit.Source, "range ") && strings.Contains(unit.Source, "map[") {
		if strings.Contains(unit.Source, "sort.") || strings.Contains(unit.Source, "json.Marshal") || strings.Contains(unit.Source, "fmt.Print") {
			// weak
			if !strings.Contains(unit.Source, "sort.") {
				for _, line := range codeLines(unit.Source) {
					t := strings.TrimSpace(line.text)
					if strings.HasPrefix(t, "for ") && strings.Contains(t, " range ") {
						// check if ranging over map var
						pushAt(unit, meta, line.byte, "map range used where ordered output may be expected", out)
						return
					}
				}
			}
		}
	}
}

func detectBP85(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-85")
	// r.Context().Value without ok check - type assert
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, ".Value(") && strings.Contains(t, ".(") && !strings.Contains(t, "ok") {
			pushAt(unit, meta, line.byte, "unchecked type assertion on request context value", out)
		}
	}
}

func detectBP91(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-91")
	// chan struct{} for data-bearing? opposite: make(chan T) for notification
	// flag make(chan bool) used as signal
	if strings.Contains(unit.Source, "make(chan bool") || strings.Contains(unit.Source, "make(chan int") {
		if strings.Contains(unit.Source, "<-") && !strings.Contains(unit.Source, "make(chan struct{}") {
			if pos := strings.Index(unit.Source, "make(chan "); pos >= 0 {
				pushAt(unit, meta, pos, "notification channel carries data; prefer chan struct{}", out)
			}
		}
	}
}

func detectBP138(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-138")
	// GORM hook with http/external
	if strings.Contains(unit.Source, "BeforeCreate") || strings.Contains(unit.Source, "AfterCreate") ||
		strings.Contains(unit.Source, "BeforeSave") || strings.Contains(unit.Source, "AfterSave") {
		if strings.Contains(unit.Source, "http.") || strings.Contains(unit.Source, "smtp") || strings.Contains(unit.Source, "Dial") {
			pos := strings.Index(unit.Source, "Before")
			if pos < 0 {
				pos = strings.Index(unit.Source, "After")
			}
			if pos >= 0 {
				pushAt(unit, meta, pos, "external I/O inside a GORM hook; keep hooks local and fast", out)
			}
		}
	}
}

func detectBP141(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-141")
	if !strings.Contains(unit.Source, "Named") && !strings.Contains(unit.Source, "sqlx") {
		return
	}
	// :name placeholders without matching db tags — hard; skip unless obvious
	if strings.Contains(unit.Source, "NamedExec") || strings.Contains(unit.Source, "NamedQuery") {
		if !strings.Contains(unit.Source, "`db:") {
			if pos := strings.Index(unit.Source, "Named"); pos >= 0 {
				pushAt(unit, meta, pos, "sqlx named query may lack matching struct db tags", out)
			}
		}
	}
}

func detectBP164(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-164")
	// functional option mutates package default
	if strings.Contains(unit.Source, "func ") && strings.Contains(unit.Source, "Option") {
		// defaultOptions or Default
		if strings.Contains(unit.Source, "default") && strings.Contains(unit.Source, "func With") {
			for _, line := range codeLines(unit.Source) {
				t := strings.TrimSpace(line.text)
				if strings.HasPrefix(t, "func With") && strings.Contains(unit.Source, "default") {
					// if assigns to package var
					if strings.Contains(unit.Source, "default") && strings.Contains(t, "func With") {
						// look for assignment to package-level in option body — weak
						if strings.Contains(unit.Source, "defaultConfig") || strings.Contains(unit.Source, "defaultOptions") {
							pushAt(unit, meta, line.byte, "functional option mutates a package-level default", out)
							return
						}
					}
				}
			}
		}
	}
}

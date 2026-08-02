package cwe

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// False-positive guard helpers shared by CWE-117 / 1341 / 367 / 88 / 93 / 215 / 502.

var (
	fStringInterpRE    = regexp.MustCompile(`\{([^{][^}]*)\}`)
	upperSnakeRE       = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	closeCallRE        = regexp.MustCompile(`(?is)(\w+)\.close\s*\(\s*\)`)
	toctouExistsCallRE = regexp.MustCompile(`(?is)os\.path\.(?:exists|lexists)\s*\(\s*(\w+)\s*\)`)
	toctouUseCallRE    = regexp.MustCompile(`(?is)(?:open|os\.remove|os\.unlink)\s*\(\s*(\w+)\b`)
)

// logMessageHasCRLFCapableValue is true when a formatted log argument can carry
// attacker-controlled CR/LF. Pure constants, len(...), numerics, loop counters,
// and internally generated names cannot.
func logMessageHasCRLFCapableValue(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if isFStringExpr(t) {
		matches := fStringInterpRE.FindAllStringSubmatch(t, -1)
		if len(matches) == 0 {
			return false
		}
		for _, m := range matches {
			if !logInterpLooksCRLFSafe(m[1], ctx) {
				return true
			}
		}
		return false
	}
	if strings.Contains(t, ".format(") {
		open := strings.Index(t, ".format(")
		inner, ok := callArgsRegion(t, open+len(".format"))
		if !ok {
			return true
		}
		for _, arg := range splitTopLevelArgs(inner) {
			if !logInterpLooksCRLFSafe(arg, ctx) {
				return true
			}
		}
		return false
	}
	if idx := indexTopLevelPercent(t); idx > 0 {
		rhs := strings.TrimSpace(t[idx+1:])
		if strings.HasPrefix(rhs, "(") && strings.HasSuffix(rhs, ")") {
			rhs = strings.TrimSpace(rhs[1 : len(rhs)-1])
		}
		for _, arg := range splitTopLevelArgs(rhs) {
			if !logInterpLooksCRLFSafe(arg, ctx) {
				return true
			}
		}
		return false
	}
	if strings.Contains(t, "+") {
		for _, part := range splitTopLevelConcat(t) {
			part = strings.TrimSpace(part)
			if isPureStringLiteral(part) || logInterpLooksCRLFSafe(part, ctx) {
				continue
			}
			return true
		}
		return false
	}
	return isDynamicExpr(t)
}

func logInterpLooksCRLFSafe(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return true
	}
	// Numeric format spec ({elapsed:.2f}, {n:03d}) proves a numeric value.
	if fstringNumericSpec(t) {
		return true
	}
	// Drop f-string conversion/format specs: {elapsed:.2f} / {n!r}
	if i := strings.IndexAny(t, ":!"); i >= 0 {
		t = strings.TrimSpace(t[:i])
	}
	if t == "" {
		return true
	}
	if isPureStringLiteral(t) || isNumericLiteral(t) {
		return true
	}
	compact := compactWhitespace(t)
	if strings.HasPrefix(compact, "len(") && strings.HasSuffix(compact, ")") {
		return true
	}
	if headerValueIsInternalNumeric(compact) {
		return true
	}
	if looksConstantStringRepeat(compact) {
		return true
	}
	if strings.HasPrefix(compact, "str(") && strings.HasSuffix(compact, ")") {
		inner := strings.TrimSpace(compact[4 : len(compact)-1])
		return logInterpLooksCRLFSafe(inner, ctx)
	}
	if strings.HasPrefix(compact, "int(") || strings.HasPrefix(compact, "float(") ||
		strings.HasPrefix(compact, "round(") {
		return true
	}
	// Known non-input calls: datetime.now(), time.time(), os.path.getsize(...).
	for _, p := range internalCallPrefixes {
		if strings.HasPrefix(compact, p) {
			return true
		}
	}
	// Numeric arithmetic over internals: {i + 1}, {i % 10}.
	if parts, ok := splitTopLevelBinary(compact); ok {
		allSafe := true
		for _, part := range parts {
			if part == "" {
				continue
			}
			if !logInterpLooksCRLFSafe(part, ctx) {
				allSafe = false
				break
			}
		}
		if allSafe {
			return true
		}
	}
	// Bare identifiers resolve against same-file evidence (loop counters,
	// numeric/counter assignments, internally generated *_name values, and
	// int/float-annotated parameters).
	if bareIdentRE.MatchString(t) {
		if ctx == nil {
			return false
		}
		return ctx.identLooksInternal(t) || ctx.numericAnnotated[t]
	}
	return false
}

// fstringNumericSpec reports a format spec that can only apply to a numeric
// value (d/f/x/o/b/e/g/% conversions).
func fstringNumericSpec(t string) bool {
	colon := strings.LastIndex(t, ":")
	if colon < 0 {
		return false
	}
	spec := t[colon+1:]
	if j := strings.IndexAny(spec, ":!"); j >= 0 {
		spec = spec[:j]
	}
	if spec == "" || strings.ContainsAny(spec, "sS") {
		return false
	}
	return strings.ContainsAny(spec, "dDxXoObBeEgGfF%")
}

func looksConstantStringRepeat(expr string) bool {
	t := strings.TrimSpace(expr)
	if !strings.Contains(t, "*") {
		return false
	}
	parts := strings.SplitN(t, "*", 2)
	if len(parts) != 2 {
		return false
	}
	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])
	return (isPureStringLiteral(left) && isNumericLiteral(right)) ||
		(isNumericLiteral(left) && isPureStringLiteral(right))
}

func isNumericLiteral(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if strings.HasPrefix(t, "-") || strings.HasPrefix(t, "+") {
		t = strings.TrimSpace(t[1:])
	}
	if t == "" {
		return false
	}
	dot := 0
	for _, r := range t {
		if r == '.' {
			dot++
			if dot > 1 {
				return false
			}
			continue
		}
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// callArgsRegion returns the interior of a call whose '(' is at openIdx.
func callArgsRegion(src string, openIdx int) (string, bool) {
	if openIdx < 0 || openIdx >= len(src) || src[openIdx] != '(' {
		return "", false
	}
	closeAt, args := scanCallArgs(src, openIdx)
	if closeAt < 0 {
		return "", false
	}
	return args, true
}

// sameHandleDoubleCloseStart reports the first double-close of the same
// receiver within one function body. Lifecycle hooks (__exit__/__del__) and
// distinct receivers are ignored.
func sameHandleDoubleCloseStart(masked string) int {
	if masked == "" {
		return -1
	}
	matches := closeCallRE.FindAllStringSubmatchIndex(masked, -1)
	for i := 0; i < len(matches); i++ {
		recv1 := masked[matches[i][2]:matches[i][3]]
		end1 := matches[i][1]
		for j := i + 1; j < len(matches); j++ {
			start2 := matches[j][0]
			if start2-end1 > closedResourceWindow {
				break
			}
			recv2 := masked[matches[j][2]:matches[j][3]]
			if recv1 != recv2 {
				continue
			}
			between := masked[end1:start2]
			if strings.Contains(between, "\ndef ") || strings.Contains(between, "\nasync def ") ||
				strings.Contains(between, "\nclass ") {
				continue
			}
			// Mutually exclusive branches release the handle once per path
			// (niquests async_session AsyncResponse vs Response close FPs).
			if closeSitesInMutuallyExclusiveBranches(between) {
				continue
			}
			// First close is on an early-exit arm (raise/return/continue/break)
			// so the second close is a different path — one release per path
			// (niquests wasi adapter retry: close then raise vs close then continue).
			if closeFollowedByPathExit(between) {
				continue
			}
			if closeSiteInLifecycleHook(masked, matches[i][0]) || closeSiteInLifecycleHook(masked, start2) {
				continue
			}
			return matches[i][0]
		}
	}
	return -1
}

// closeSitesInMutuallyExclusiveBranches reports that the text between two
// close() calls is only an else/elif arm boundary (same handle, one path).
func closeSitesInMutuallyExclusiveBranches(between string) bool {
	// Typical shape:
	//   if ...:
	//       await resp.close()
	//   else:
	//       resp.close()
	trimmed := strings.TrimSpace(between)
	if trimmed == "" {
		return false
	}
	// Must contain else/elif and no other close calls in between.
	if !strings.Contains(between, "else:") && !strings.Contains(between, "elif ") {
		return false
	}
	if closeCallRE.MatchString(between) {
		return false
	}
	// Reject if substantial non-branch logic sits between closes.
	for _, line := range strings.Split(between, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		if t == "else:" || strings.HasPrefix(t, "else:") || strings.HasPrefix(t, "elif ") {
			continue
		}
		// Allow blank structural noise only.
		return false
	}
	return true
}

// closeFollowedByPathExit reports that the first significant statement after
// a close() exits the path (raise/return/continue/break), so a later close of
// the same name is on a different control-flow path (retry-loop one-close-per-
// iteration, or if/raise vs fall-through close).
func closeFollowedByPathExit(between string) bool {
	for _, line := range strings.Split(between, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		if strings.HasPrefix(t, "raise") || strings.HasPrefix(t, "return") ||
			t == "continue" || strings.HasPrefix(t, "continue ") ||
			t == "break" || strings.HasPrefix(t, "break ") {
			return true
		}
		// First non-exit statement means the path continues toward the second close.
		return false
	}
	return false
}

func closeSiteInLifecycleHook(src string, offset int) bool {
	if offset < 0 || offset > len(src) {
		return false
	}
	head := src[:offset]
	defAt := strings.LastIndex(head, "\ndef ")
	asyncAt := strings.LastIndex(head, "\nasync def ")
	at := defAt
	prefix := "\ndef "
	if asyncAt > defAt {
		at = asyncAt
		prefix = "\nasync def "
	}
	if at < 0 {
		// module-top def without leading newline
		if strings.HasPrefix(strings.TrimLeft(head, " \t"), "def ") {
			at = strings.Index(head, "def ")
			prefix = "def "
		} else {
			return false
		}
	} else {
		at++ // skip newline for \ndef
		prefix = strings.TrimPrefix(prefix, "\n")
	}
	rest := head[at:]
	if !strings.HasPrefix(rest, prefix) && !strings.HasPrefix(rest, "def ") && !strings.HasPrefix(rest, "async def ") {
		rest = src[at:]
	}
	lineEnd := strings.IndexByte(rest, '\n')
	if lineEnd < 0 {
		lineEnd = len(rest)
	}
	header := rest[:lineEnd]
	return strings.Contains(header, "def __exit__") || strings.Contains(header, "def __del__") ||
		strings.Contains(header, "async def __aexit__")
}

func toctouSamePathStart(facts *PyCweFacts, source string) int {
	if source == "" {
		return -1
	}
	var masked string
	if facts != nil {
		masked = facts.codeMask(source, fragStartHint(facts, source))
	} else {
		masked = pythonCodeMask(source)
	}
	exists := toctouExistsCallRE.FindAllStringSubmatchIndex(masked, -1)
	uses := toctouUseCallRE.FindAllStringSubmatchIndex(masked, -1)
	for _, ex := range exists {
		pathName := masked[ex[2]:ex[3]]
		exEnd := ex[1]
		for _, use := range uses {
			if use[0] < exEnd {
				continue
			}
			const toctouWindowBytes = 300
			if use[0]-exEnd > toctouWindowBytes {
				break
			}
			useName := masked[use[2]:use[3]]
			if useName == pathName {
				return ex[0]
			}
		}
	}
	return -1
}

// splitAssignmentEq splits a true assignment (=), rejecting ==, !=, <=, >=, :=.
func splitAssignmentEq(line string) (string, string, bool) {
	inStr := byte(0)
	esc := false
	triple := false
	depth := 0
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inStr != 0 {
			i, inStr, esc, triple = advanceInString(line, i, inStr, esc, triple)
			continue
		}
		switch c {
		case '\'', '"':
			inStr = c
			if i+2 < len(line) && line[i+1] == c && line[i+2] == c {
				triple = true
				i += 2
			}
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '=':
			if depth != 0 {
				continue
			}
			if i > 0 {
				prev := line[i-1]
				if prev == '=' || prev == '!' || prev == '<' || prev == '>' || prev == ':' {
					continue
				}
			}
			if i+1 < len(line) && line[i+1] == '=' {
				continue
			}
			return line[:i], line[i+1:], true
		}
	}
	return "", "", false
}

// advanceInString steps one byte inside a string literal scanner state.
func advanceInString(line string, i int, inStr byte, esc, triple bool) (int, byte, bool, bool) {
	c := line[i]
	if triple {
		if c == inStr && i+2 < len(line) && line[i+1] == inStr && line[i+2] == inStr {
			return i + 2, 0, false, false
		}
		return i, inStr, esc, triple
	}
	if esc {
		return i, inStr, false, triple
	}
	if c == '\\' {
		return i, inStr, true, triple
	}
	if c == inStr {
		return i, 0, false, triple
	}
	return i, inStr, esc, triple
}

// sensitiveIdentOutsideLiterals reports a sensitive identifier that is not only
// present as an English word inside a string/f-string literal.
func sensitiveIdentOutsideLiterals(args string) bool {
	if args == "" {
		return false
	}
	masked := pythonCodeMask(args)
	// Keep f-string interpolations visible: Mask blanks the whole literal,
	// including {expr}. Re-inject {...} bodies from the original.
	if isFStringExpr(strings.TrimSpace(args)) || strings.Contains(args, "f\"") || strings.Contains(args, "f'") ||
		strings.Contains(args, "F\"") || strings.Contains(args, "F'") {
		var b strings.Builder
		b.Grow(len(args))
		for _, m := range fStringInterpRE.FindAllStringSubmatchIndex(args, -1) {
			b.WriteByte(' ')
			b.WriteString(args[m[2]:m[3]])
			b.WriteByte(' ')
		}
		masked += b.String()
	}
	return sensitiveValueRE.MatchString(masked)
}

func argvSegmentLooksTrusted(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return true
	}
	if t == "sys.executable" {
		return true
	}
	if upperSnakeRE.MatchString(t) {
		return true
	}
	// Variadic spread of a resolved-internal collection: [prog, *args].
	if strings.HasPrefix(t, "*") {
		ident := strings.TrimSpace(t[1:])
		return bareIdentRE.MatchString(ident) && ctx != nil && ctx.identLooksInternal(ident)
	}
	if strings.HasPrefix(t, "str(") && strings.HasSuffix(t, ")") {
		inner := strings.TrimSpace(t[4 : len(t)-1])
		if upperSnakeRE.MatchString(inner) || inner == "sys.executable" {
			return true
		}
		if bareIdentRE.MatchString(inner) && ctx != nil && ctx.identLooksInternal(inner) {
			return true
		}
		if looksTrustedPathConcat(inner, ctx) {
			return true
		}
		// str(X.relative_to(C)) / str(X.resolve()) — path method on a constant
		// or internal receiver.
		if dot := strings.IndexByte(inner, '.'); dot > 0 {
			recv := strings.TrimSpace(inner[:dot])
			rest := strings.TrimSpace(inner[dot:])
			if strings.HasPrefix(rest, ".") && strings.Contains(rest, "(") && strings.HasSuffix(rest, ")") {
				if upperSnakeRE.MatchString(recv) || (ctx != nil && ctx.identLooksInternal(recv)) {
					return true
				}
			}
		}
	}
	// Path normalization of a constant path: str(X.relative_to(C)).replace(a, b).
	if strings.HasSuffix(t, ")") {
		if stripped, ok := stripTrailingReplaceCall(t); ok {
			return argvSegmentLooksTrusted(stripped, ctx)
		}
	}
	if bareIdentRE.MatchString(t) && ctx != nil && ctx.identLooksInternal(t) {
		return true
	}
	// f-string argv elements whose interpolations are all internal/numeric
	// (f"--remote-debugging-port={port}") carry no option text.
	if isFStringExpr(t) && fstringInterpsSafe(t, ctx) {
		return true
	}
	if looksTrustedArgvConcat(t, ctx) {
		return true
	}
	return false
}

// stripTrailingReplaceCall removes a trailing .replace(<lit>, <lit>) call.
func stripTrailingReplaceCall(expr string) (string, bool) {
	t := strings.TrimSpace(expr)
	idx := strings.LastIndex(t, ".replace(")
	if idx < 0 {
		return "", false
	}
	open := idx + len(".replace")
	closeAt, inner := scanBracketRegion(t, open)
	if closeAt < 0 || strings.TrimSpace(t[closeAt+1:]) != "" {
		return "", false
	}
	args := splitTopLevelArgs(inner)
	if len(args) != 2 || !isPureStringLiteral(strings.TrimSpace(args[0])) || !isPureStringLiteral(strings.TrimSpace(args[1])) {
		return "", false
	}
	return strings.TrimSpace(t[:idx]), true
}

// looksTrustedPathConcat reports pathlib-style a / b / c paths assembled from
// constants and internal identifiers.
func looksTrustedPathConcat(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	parts := splitTopLevelWord(t, "/")
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return false
		}
		if isPureStringLiteral(p) || upperSnakeRE.MatchString(p) {
			continue
		}
		if bareIdentRE.MatchString(p) && ctx != nil && ctx.identLooksInternal(p) {
			continue
		}
		if compact := compactWhitespace(p); strings.HasPrefix(compact, "Path(") || strings.HasPrefix(compact, "pathlib.Path(") {
			continue
		}
		// Parenthesized literal ternary: ("playwright" if cond else "playwright.cmd")
		if strings.HasPrefix(p, "(") && strings.HasSuffix(p, ")") {
			inner := strings.TrimSpace(p[1 : len(p)-1])
			if ternaryOfLiterals(inner) {
				continue
			}
		}
		return false
	}
	return true
}

// ternaryOfLiterals reports a ?-free Python conditional expr of literals:
// A if B else C where A and C are string literals.
func ternaryOfLiterals(expr string) bool {
	idx := strings.Index(expr, " if ")
	if idx < 0 {
		return isPureStringLiteral(strings.TrimSpace(expr))
	}
	elseIdx := strings.Index(expr, " else ")
	if elseIdx < 0 {
		return false
	}
	return isPureStringLiteral(strings.TrimSpace(expr[:idx])) &&
		isPureStringLiteral(strings.TrimSpace(expr[elseIdx+len(" else "):]))
}

func looksTrustedArgvConcat(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	if !strings.Contains(t, "+") {
		return false
	}
	parts := splitTopLevelConcat(t)
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if isPureStringLiteral(p) || isPureBytesLiteral(p) || argvSegmentLooksTrusted(p, ctx) {
			continue
		}
		return false
	}
	return true
}

func sourceUsesRuamelYAML(src string) bool {
	return strings.Contains(src, "ruamel.yaml") ||
		strings.Contains(src, "from ruamel") ||
		strings.Contains(src, "import ruamel")
}

// yamlLoadLooksLikeRuamel is true for ruamel.YAML().load-style calls (no
// PyYAML Loader= kwarg). PyYAML yaml.load(..., Loader=...) stays reportable.
func yamlLoadLooksLikeRuamel(src, args string) bool {
	compact := compactWhitespace(args)
	if strings.Contains(compact, "Loader=") {
		return false
	}
	return sourceUsesRuamelYAML(src) || strings.Contains(src, "YAML()")
}

// ---------------------------------------------------------------------------
// Per-file identifier evidence for CWE-88 / CWE-93 / CWE-117 guardrails.
// ---------------------------------------------------------------------------

var (
	pyDefHeaderRE = regexp.MustCompile(`(?m)^[ \t]*(?:async\s+)?def\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`)
	pyAugAssignRE = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*(\+=|-=|\*=|/=|//=|%=|^=|\*\*=|<<=|>>=|&=|\|=)`)
	pyWithAsRE    = regexp.MustCompile(`as\s+([A-Za-z_][A-Za-z0-9_]*)\b`)
	bareIdentRE   = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

// internalCallPrefixes are calls that never read user/request/env input:
// numerics, clock reads, path/tempfile/socket helpers, and uuid generation.
var internalCallPrefixes = []string{
	"len(", "max(", "min(", "sum(", "abs(", "round(", "int(", "float(",
	"time.time(", "time.perf_counter(", "time.monotonic(", "time.process_time(",
	"datetime.now(", "datetime.utcnow(",
	"os.path.getsize(", "os.path.exists(", "os.path.isfile(", "os.path.isdir(",
	"os.path.join(", "os.getcwd(", "os.path.abspath(",
	"tempfile.mkdtemp(", "tempfile.mkstemp(",
	"pathlib.Path(", "Path(", "shutil.which(", "socket.gethostname(",
	"uuid.uuid4(", "uuid.uuid1(", "uuid.uuid5(",
}

// pySigParam is one function parameter of a same-file def header.
type pySigParam struct {
	name       string
	hasDefault bool
	litDefault bool
	numeric    bool // annotation names an int/float type
}

type pySig struct {
	params []pySigParam
	self   bool // first parameter is self/cls (method form shifts call args)
}

type pyCallSite struct {
	args   []string
	method bool // called as obj.name(...) — positional args shift past self
}

// pyFileCtx is the fixed-point identifier evidence for one scanned file.
type pyFileCtx struct {
	internal         map[string]bool
	loopTargets      map[string]bool
	numericAnnotated map[string]bool
}

func (c *pyFileCtx) identLooksInternal(ident string) bool {
	if c == nil {
		return false
	}
	if c.internal[ident] || c.loopTargets[ident] {
		return true
	}
	// *_name identifiers are internally generated names (thread names, pool
	// names, config-record fields), not external input.
	return strings.HasSuffix(ident, "_name")
}

// buildPyFileCtx resolves same-file identifiers to constant/internal values:
// loop counters, numeric/time/tempfile/socket assignments, literal bindings,
// and parameters whose same-file call sites pass only literals.
func buildPyFileCtx(facts *PyCweFacts, source string) *pyFileCtx {
	ctx := &pyFileCtx{
		internal:         map[string]bool{},
		loopTargets:      map[string]bool{},
		numericAnnotated: map[string]bool{},
	}
	if source == "" {
		return ctx
	}
	lines := maskedLinesOf(facts, source)
	sigs := map[string]*pySig{}
	callArgsByParam := map[string]map[int][]string{}
	collectPySignatures(lines, source, sigs)
	collectPyCalls(lines, source, sigs, callArgsByParam)
	for _, s := range sigs {
		for _, p := range s.params {
			if p.numeric {
				ctx.numericAnnotated[p.name] = true
			}
		}
	}
	for iter := 0; iter < 8; iter++ {
		changed := false
		for i := 0; i < len(lines); i++ {
			line := lines[i]
			orig := origLineAt(source, line)
			// Merge continuation lines so multi-line list/tuple/call
			// assignments (components = ["a", "b", ...]) classify correctly.
			consumed := 0
			for consumed < 7 && i+consumed+1 < len(lines) && pyBracketsUnbalanced(orig) {
				consumed++
				orig += "\n" + origLineAt(source, lines[i+consumed])
			}
			i += consumed
			if classifyPyLine(line, orig, ctx) {
				changed = true
			}
		}
		for fn, byIdx := range callArgsByParam {
			s := sigs[fn]
			if s == nil {
				continue
			}
			for idx, argTexts := range byIdx {
				if idx >= len(s.params) {
					continue
				}
				p := s.params[idx]
				if p.name == "" || ctx.internal[p.name] {
					continue
				}
				allOK := true
				for _, at := range argTexts {
					if !callArgLooksInternal(at, ctx) {
						allOK = false
						break
					}
				}
				if allOK {
					ctx.internal[p.name] = true
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	return ctx
}

func maskedLinesOf(facts *PyCweFacts, source string) []pyMaskedLine {
	if facts != nil {
		if lines := facts.MaskedLines(); len(lines) > 0 {
			return lines
		}
	}
	return buildMaskedPythonLines(pythonCodeMask(source))
}

// origLineAt recovers the unmasked line bytes aligned with a masked line span.
func origLineAt(source string, line pyMaskedLine) string {
	if line.start < 0 || line.start > len(source) {
		return ""
	}
	end := line.start + len(line.text)
	if end > len(source) {
		end = len(source)
	}
	return source[line.start:end]
}

// classifyPyLine feeds one line of evidence into ctx and reports additions.
func classifyPyLine(line pyMaskedLine, orig string, ctx *pyFileCtx) bool {
	changed := false
	masked := line.text
	trimmed := strings.TrimSpace(masked)

	// Loop targets over internal iterables (range/enumerate/zip, module
	// constants, literal collections, or already-internal collections).
	if strings.HasPrefix(trimmed, "for ") {
		if targets, iterable, ok := pyForTargets(trimmed); ok && pyIterableLooksInternal(iterable, ctx) {
			for _, t := range targets {
				if t != "" && !ctx.loopTargets[t] {
					ctx.loopTargets[t] = true
					changed = true
				}
			}
		}
	}

	// tempfile context manager targets: with tempfile.X() as <name>:
	if strings.Contains(masked, "tempfile.") && strings.Contains(masked, " as ") {
		if m := pyWithAsRE.FindStringSubmatch(trimmed); len(m) == 2 {
			if !ctx.internal[m[1]] {
				ctx.internal[m[1]] = true
				changed = true
			}
		}
	}

	// Augmented counters: attempt += 1 → numeric counter.
	if m := pyAugAssignRE.FindStringSubmatch(trimmed); len(m) == 3 {
		if !ctx.internal[m[1]] {
			ctx.internal[m[1]] = true
			changed = true
		}
	}

	// Plain assignments: X = <rhs>.
	lhs, _, ok := splitAssignmentEq(masked)
	if ok && bareIdentRE.MatchString(strings.TrimSpace(lhs)) {
		_, origRhs, origOK := splitAssignmentEq(orig)
		if origOK && rhsLooksInternal(origRhs, ctx) && !ctx.internal[strings.TrimSpace(lhs)] {
			ctx.internal[strings.TrimSpace(lhs)] = true
			changed = true
		}
	}
	return changed
}

// pyForTargets extracts loop targets and the iterable expression text.
func pyForTargets(line string) ([]string, string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "for ") {
		return nil, "", false
	}
	rest := strings.TrimSpace(trimmed[len("for "):])
	inIdx := strings.Index(rest, " in ")
	if inIdx < 0 {
		return nil, "", false
	}
	targets := splitTopLevelArgs(rest[:inIdx])
	iterable := strings.TrimSpace(rest[inIdx+len(" in "):])
	iterable = strings.TrimSuffix(strings.TrimSpace(iterable), ":")
	return targets, iterable, true
}

func pyIterableLooksInternal(iterable string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(iterable)
	if t == "" {
		return true // masked string literal iterable
	}
	compact := compactWhitespace(t)
	if strings.HasPrefix(compact, "range(") || strings.HasPrefix(compact, "enumerate(") ||
		strings.HasPrefix(compact, "zip(") || strings.HasPrefix(compact, "reversed(") {
		return true
	}
	if upperSnakeRE.MatchString(t) {
		return true
	}
	if looksStaticStringListLenient(t) || isPureStringLiteral(t) {
		return true
	}
	if bareIdentRE.MatchString(t) {
		return ctx != nil && ctx.internal[t]
	}
	return false
}

// looksStaticStringListLenient accepts literal lists with a trailing comma:
// ["a", "b",] — common for multi-line module constants.
func looksStaticStringListLenient(expr string) bool {
	t := strings.TrimSpace(expr)
	if looksStaticStringList(t) {
		return true
	}
	if len(t) >= 2 && (t[0] == '[' || t[0] == '(') {
		closeB := byte(']')
		if t[0] == '(' {
			closeB = ')'
		}
		if t[len(t)-1] != closeB {
			return false
		}
		idx := strings.LastIndex(t[:len(t)-1], ",")
		if idx < 0 {
			return false
		}
		return looksStaticStringList(t[:idx] + t[idx+1:])
	}
	return false
}

// rhsLooksInternal classifies an assignment right-hand side as constant or
// internally generated.
func rhsLooksInternal(rhs string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(rhs)
	if t == "" {
		return false
	}
	if isPureStringLiteral(t) || isNumericLiteral(t) || looksStaticStringListLenient(t) {
		return true
	}
	if t == "sys.executable" || upperSnakeRE.MatchString(t) {
		return true
	}
	if bareIdentRE.MatchString(t) {
		return ctx != nil && ctx.identLooksInternal(t)
	}
	if looksConstantStringRepeat(t) {
		return true
	}
	if isFStringExpr(t) && fstringInterpsSafe(t, ctx) {
		return true
	}
	compact := compactWhitespace(t)
	for _, p := range internalCallPrefixes {
		if strings.HasPrefix(compact, p) {
			return true
		}
	}
	if strings.Contains(compact, "getsockname()") {
		return true
	}
	if looksOrChainInternal(t, ctx) {
		return true
	}
	return looksPathConcatInternal(t, ctx)
}

// looksOrChainInternal reports X or Y or Z chains whose parts are all
// literals or internal identifiers (fallback-constant assignments).
func looksOrChainInternal(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	parts := splitTopLevelWord(t, "or")
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if isPureStringLiteral(p) || isNumericLiteral(p) {
			continue
		}
		if bareIdentRE.MatchString(p) && ctx != nil && ctx.identLooksInternal(p) {
			continue
		}
		return false
	}
	return true
}

// looksPathConcatInternal reports pathlib-style a / b / c chains whose parts
// are literals, module constants, or internal identifiers.
func looksPathConcatInternal(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	parts := splitTopLevelWord(t, "/")
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return false
		}
		if isPureStringLiteral(p) || isNumericLiteral(p) || upperSnakeRE.MatchString(p) {
			continue
		}
		if bareIdentRE.MatchString(p) && ctx != nil && ctx.identLooksInternal(p) {
			continue
		}
		if compact := compactWhitespace(p); strings.HasPrefix(compact, "Path(") || strings.HasPrefix(compact, "pathlib.Path(") {
			continue
		}
		if strings.HasPrefix(p, "(") && strings.HasSuffix(p, ")") {
			if ternaryOfLiterals(strings.TrimSpace(p[1 : len(p)-1])) {
				continue
			}
		}
		return false
	}
	return true
}

// callArgLooksInternal reports a call-site argument that is a literal, a
// resolved internal identifier, or a path concat (fixed point over ctx.internal).
func callArgLooksInternal(arg string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(arg)
	if t == "" {
		return true
	}
	if isPureStringLiteral(t) || isNumericLiteral(t) {
		return true
	}
	if looksPathConcatInternal(t, ctx) {
		return true
	}
	return bareIdentRE.MatchString(t) && ctx != nil && ctx.identLooksInternal(t)
}

// collectPySignatures parses def headers into parameter evidence.
func collectPySignatures(lines []pyMaskedLine, source string, sigs map[string]*pySig) {
	for _, line := range lines {
		m := pyDefHeaderRE.FindStringSubmatchIndex(line.text)
		if m == nil {
			continue
		}
		name := line.text[m[2]:m[3]]
		orig := origLineAt(source, line)
		params := parsePySigParams(orig)
		s := &pySig{params: params}
		if len(params) > 0 && (params[0].name == "self" || params[0].name == "cls") {
			s.self = true
		}
		sigs[name] = s
	}
}

func parsePySigParams(defLine string) []pySigParam {
	open := strings.Index(defLine, "(")
	if open < 0 {
		return nil
	}
	closeAt, _ := scanCallArgs(defLine, open)
	if closeAt < 0 {
		return nil
	}
	args := splitTopLevelArgs(defLine[open+1 : closeAt])
	var out []pySigParam
	for _, a := range args {
		a = strings.TrimSpace(a)
		if a == "" || strings.HasPrefix(a, "**") {
			continue
		}
		a = strings.TrimPrefix(a, "*")
		p := pySigParam{}
		if i := strings.IndexAny(a, ":="); i >= 0 {
			head := a[:i]
			body := a[i:]
			p.name = strings.TrimSpace(head)
			if body != "" && body[0] == ':' {
				rest := strings.TrimSpace(body[1:])
				annotation := rest
				if j := strings.IndexByte(rest, '='); j >= 0 {
					annotation = strings.TrimSpace(rest[:j])
					p.hasDefault = true
					p.litDefault = pyDefaultLooksLiteral(strings.TrimSpace(rest[j+1:]))
				}
				c := compactWhitespace(annotation)
				if strings.Contains(c, "int") || strings.Contains(c, "float") {
					p.numeric = true
				}
			} else {
				p.hasDefault = true
				p.litDefault = pyDefaultLooksLiteral(strings.TrimSpace(body[1:]))
			}
		} else {
			p.name = a
		}
		out = append(out, p)
	}
	return out
}

func pyDefaultLooksLiteral(def string) bool {
	if def == "" {
		return false
	}
	return isPureStringLiteral(def) || isNumericLiteral(def)
}

// collectPyCalls records same-file call sites per def name and maps positional
// arguments onto parameter indexes. Both bare F(...) and method obj.F(...)
// forms are matched (a preceding '.' is an allowed left boundary).
func collectPyCalls(lines []pyMaskedLine, source string, sigs map[string]*pySig, out map[string]map[int][]string) {
	masked := pythonCodeMask(source)
	for fn := range sigs {
		byIdx := map[int][]string{}
		start := 0
		for {
			idx := strings.Index(masked[start:], fn)
			if idx < 0 {
				break
			}
			abs := start + idx
			if !callNameBoundaryOK(masked, abs, abs+len(fn)) {
				start = abs + len(fn)
				continue
			}
			// Skip the def header itself.
			prefix := strings.TrimRight(masked[:abs], " \t")
			if strings.HasSuffix(prefix, "def") {
				start = abs + len(fn)
				continue
			}
			after := abs + len(fn)
			j := skipWS(masked, after)
			if j >= len(masked) || masked[j] != '(' {
				start = after
				continue
			}
			closeAt, _ := scanCallArgs(source, j)
			if closeAt < 0 {
				start = after
				continue
			}
			method := abs > 0 && masked[abs-1] == '.'
			s := sigs[fn]
			assignPyCallArgs(s, method, source[j+1:closeAt], byIdx)
			start = closeAt + 1
		}
		if len(byIdx) > 0 {
			out[fn] = byIdx
		}
		collectThreadStyleCalls(masked, source, fn, byIdx)
		if len(byIdx) > 0 {
			out[fn] = byIdx
		}
	}
}

// callNameBoundaryOK accepts an identifier-boundary match where the char left
// of the name is not an identifier character (a preceding '.' is allowed).
func callNameBoundaryOK(masked string, start, end int) bool {
	if start > 0 {
		r, _ := utf8.DecodeLastRuneInString(masked[:start])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return false
		}
	}
	if end < len(masked) {
		r, _ := utf8.DecodeRuneInString(masked[end:])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return false
		}
	}
	return true
}

// collectThreadStyleCalls records threading.Thread(target=fn, args=[...]) and
// executor.submit(fn, ...) call shapes as call-site evidence for fn.
func collectThreadStyleCalls(masked, source, fn string, byIdx map[int][]string) {
	// Thread/Process constructor form: target=<fn>, args=[...] or args=(...).
	search := 0
	for {
		idx := strings.Index(masked[search:], "target="+fn)
		if idx < 0 {
			break
		}
		abs := search + idx + len("target=")
		if !callNameBoundaryOK(masked, abs, abs+len(fn)) {
			search = abs + len(fn)
			continue
		}
		lineEnd := strings.IndexByte(masked[abs:], '\n')
		window := len(masked)
		if lineEnd >= 0 {
			window = abs + lineEnd
		}
		if argsIdx := strings.Index(masked[abs:window], "args="); argsIdx >= 0 {
			valueStart := abs + argsIdx + len("args=")
			valueStart = skipWS(masked, valueStart)
			if valueStart < len(masked) && (masked[valueStart] == '[' || masked[valueStart] == '(') {
				closeAt, inner := scanBracketRegion(masked, valueStart)
				if closeAt >= 0 {
					args := splitTopLevelArgs(inner)
					for k, a := range args {
						byIdx[k] = append(byIdx[k], strings.TrimSpace(a))
					}
				}
			}
		}
		search = abs + len(fn)
	}
	// Executor form: .submit(<fn>, ...) — remaining positional args.
	for _, callName := range []string{".submit(", ".apply_async(", ".apply("} {
		search = 0
		for {
			idx := strings.Index(masked[search:], callName)
			if idx < 0 {
				break
			}
			abs := search + idx
			open := abs + len(callName) - 1
			closeAt, inner := scanCallArgs(source, open)
			if closeAt < 0 {
				search = abs + len(callName)
				continue
			}
			args := splitTopLevelArgs(inner)
			if len(args) > 1 && strings.TrimSpace(args[0]) == fn {
				for k, a := range args[1:] {
					byIdx[k] = append(byIdx[k], strings.TrimSpace(a))
				}
			}
			search = closeAt + 1
		}
	}
}

func assignPyCallArgs(s *pySig, method bool, argsText string, byIdx map[int][]string) {
	args := splitTopLevelArgs(argsText)
	positional := []string{}
	kw := map[string]string{}
	for _, a := range args {
		if key, val, ok := strings.Cut(a, "="); ok && bareIdentRE.MatchString(strings.TrimSpace(key)) {
			kw[strings.TrimSpace(key)] = val
		} else {
			positional = append(positional, a)
		}
	}
	if s == nil || len(s.params) == 0 {
		return
	}
	base := 0
	if method && s.self {
		base = 1
	}
	for i, arg := range positional {
		idx := base + i
		if idx >= 0 && idx < len(s.params) {
			byIdx[idx] = append(byIdx[idx], arg)
		}
	}
	for name, val := range kw {
		for idx, p := range s.params {
			if p.name == name {
				byIdx[idx] = append(byIdx[idx], val)
				break
			}
		}
	}
}

// splitTopLevelWord splits on a word/operator token at depth 0 outside strings.
func splitTopLevelWord(expr, token string) []string {
	if token == "" {
		return []string{expr}
	}
	var parts []string
	depth := 0
	inStr := byte(0)
	esc := false
	triple := false
	start := 0
	for i := 0; i < len(expr); {
		if inStr != 0 {
			var t bool
			i, inStr, esc, t = advanceInString(expr, i, inStr, esc, triple)
			triple = t
			i++
			continue
		}
		c := expr[i]
		if c == '\'' || c == '"' {
			inStr = c
			if i+2 < len(expr) && expr[i+1] == c && expr[i+2] == c {
				triple = true
				i += 2
			}
			i++
			continue
		}
		switch c {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		default:
			if depth == 0 && strings.HasPrefix(expr[i:], token) {
				parts = append(parts, strings.TrimSpace(expr[start:i]))
				i += len(token)
				start = i
				continue
			}
		}
		i++
	}
	parts = append(parts, strings.TrimSpace(expr[start:]))
	return parts
}

// splitTopLevelBinary splits on any of the given operator characters at depth 0.
func splitTopLevelBinary(expr string) ([]string, bool) {
	parts := []string{}
	depth := 0
	inStr := byte(0)
	esc := false
	triple := false
	start := 0
	found := false
	for i := 0; i < len(expr); {
		if inStr != 0 {
			var t bool
			i, inStr, esc, t = advanceInString(expr, i, inStr, esc, triple)
			triple = t
			i++
			continue
		}
		c := expr[i]
		if c == '\'' || c == '"' {
			inStr = c
			if i+2 < len(expr) && expr[i+1] == c && expr[i+2] == c {
				triple = true
				i += 2
			}
			i++
			continue
		}
		switch c {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '+', '-', '*', '/', '%', '&', '|', '^':
			if depth == 0 {
				found = true
				parts = append(parts, strings.TrimSpace(expr[start:i]))
				start = i + 1
			}
		}
		i++
	}
	if !found {
		return nil, false
	}
	parts = append(parts, strings.TrimSpace(expr[start:]))
	return parts, true
}

// fstringInterpsSafe reports an f-string whose interpolations are all
// CRLF-free / internal expressions.
func fstringInterpsSafe(expr string, ctx *pyFileCtx) bool {
	matches := fStringInterpRE.FindAllStringSubmatch(expr, -1)
	if len(matches) == 0 {
		return true
	}
	for _, m := range matches {
		if !logInterpLooksCRLFSafe(m[1], ctx) {
			return false
		}
	}
	return true
}

// pyBracketsUnbalanced reports a text with unclosed ( [ { brackets outside
// string literals — i.e. a statement that continues on the next line.
func pyBracketsUnbalanced(text string) bool {
	depth := 0
	inStr := byte(0)
	esc := false
	triple := false
	for i := 0; i < len(text); i++ {
		if inStr != 0 {
			var t bool
			i, inStr, esc, t = advanceInString(text, i, inStr, esc, triple)
			triple = t
			continue
		}
		c := text[i]
		if c == '"' || c == '\'' {
			inStr = c
			if i+2 < len(text) && text[i+1] == c && text[i+2] == c {
				triple = true
				i += 2
			}
			continue
		}
		switch c {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		}
	}
	return depth > 0
}

// scanBracketRegion returns the matching close bracket index and the interior
// of a (...) or [...] region starting at open (string-aware).
func scanBracketRegion(text string, open int) (int, string) {
	if open < 0 || open >= len(text) || (text[open] != '(' && text[open] != '[') {
		return -1, ""
	}
	depth := 1 // the opening bracket at open
	inStr := byte(0)
	esc := false
	triple := false
	for i := open + 1; i < len(text); i++ {
		if inStr != 0 {
			var t bool
			i, inStr, esc, t = advanceInString(text, i, inStr, esc, triple)
			triple = t
			continue
		}
		c := text[i]
		if c == '"' || c == '\'' {
			inStr = c
			if i+2 < len(text) && text[i+1] == c && text[i+2] == c {
				triple = true
				i += 2
			}
			continue
		}
		switch c {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i, text[open+1 : i]
			}
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i, text[open+1 : i]
			}
		}
	}
	return -1, ""
}

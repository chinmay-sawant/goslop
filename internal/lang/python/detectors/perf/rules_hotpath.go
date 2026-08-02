package perf

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("PERF-PY-23", detectPYPERF23)
	RegisterRule("PERF-PY-24", detectPYPERF24)
	RegisterRule("PERF-PY-25", detectPYPERF25)
	RegisterRule("PERF-PY-26", detectPYPERF26)
	RegisterRule("PERF-PY-27", detectPYPERF27)
	RegisterRule("PERF-PY-28", detectPYPERF28)
	RegisterRule("PERF-PY-29", detectPYPERF29)
	RegisterRule("PERF-PY-30", detectPYPERF30)
}

var (
	encodeInLoopRE = regexp.MustCompile(
		`\b(encode_value|encode_dict|encode_object|serialize|to_json|json\.dumps)\s*\(`,
	)
	measureHelpers = []string{"wrap_text", "text_width", "measure_", "compute_height", "line_count", "row_height"}
	decodeHotRE    = regexp.MustCompile(
		`\b(decode_[A-Za-z0-9_]+|parse_[A-Za-z0-9_]+|Image\.open|zlib\.decompress)\s*\(`,
	)
	fromFileRE = regexp.MustCompile(
		`\b(from_file|load_config|load_[A-Za-z0-9_]+|read_bytes|read_text)\s*\(|\.read_bytes\s*\(|\.read_text\s*\(`,
	)
	executorRE = regexp.MustCompile(
		`\b(ThreadPoolExecutor|ProcessPoolExecutor)\s*\(`,
	)
	bytearrayEstimateRE = regexp.MustCompile(
		`bytearray\s*\(\s*(?:estimate|sum\s*\(|total)`,
	)
)

// PERF-PY-23: polymorphic encode/serialize per loop item.
// Seed/migrate helpers (function-name structural skip below) are not hot paths.
func detectPYPERF23(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !encodeInLoopRE.MatchString(line.text) {
			continue
		}
		if !facts.lineInLoop(i) {
			continue
		}
		if strings.Contains(line.text, "encode_record") || strings.Contains(line.text, "encode_cell") ||
			strings.Contains(line.text, "encode_leaf") {
			continue
		}
		// One-shot schema seed / migration loops (billing seed_plans).
		if start, _ := functionWindow(facts.lines, i); start >= 0 && start < len(facts.lines) {
			header := facts.lines[start].trim
			if perfFunctionIsSeedOrMigrate(header) {
				continue
			}
		}
		pushLine(unit, "PERF-PY-23", line, "encode", "polymorphic encode/serialize runs inside a hot loop; specialize fixed-schema leaves or encode outside the loop", out)
	}
}

func perfFunctionIsSeedOrMigrate(header string) bool {
	h := strings.ToLower(header)
	// def seed_plans( / def migrate( / async def seed_...
	if strings.Contains(h, "def seed_") || strings.Contains(h, "def seed(") {
		return true
	}
	if strings.Contains(h, "def migrate") || strings.Contains(h, "def _seed") {
		return true
	}
	return false
}

// PERF-PY-24: same pure measure helper used twice in one function.
func detectPYPERF24(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		trimmed := line.trim
		if !strings.HasPrefix(trimmed, "def ") && !strings.HasPrefix(trimmed, "async def ") {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if end <= start {
			continue
		}
		counts := map[string]int{}
		var first *codeLine
		for j := start; j < end; j++ {
			text := facts.lines[j].text
			for _, helper := range measureHelpers {
				if !strings.Contains(text, helper+"(") {
					continue
				}
				counts[helper]++
				if counts[helper] == 1 && first == nil {
					l := facts.lines[j]
					first = &l
				}
				if counts[helper] >= 2 && first != nil &&
					windowHas(facts.lines, start, end, "draw_", "paint", "render", "emit") {
					pushLine(unit, "PERF-PY-24", *first, helper, "pure measure helper is invoked more than once for the same logical unit; return measure+payload from one pass", out)
					return
				}
			}
		}
		if windowHas(facts.lines, start, end, "row_height") && windowHas(facts.lines, start, end, "wrap_text") {
			pushLine(unit, "PERF-PY-24", facts.lines[i], "row_height", "row height and wrap_text both run for the same unit; measure once and reuse lines", out)
			return
		}
	}
}

var perfHeavyCtorRE = regexp.MustCompile(`\b([A-Z][A-Za-z0-9_]*)\s*\(`)

// PERF-PY-25: constructor or lambda-per-item in loops.
// Skips common exception/builtin constructors; requires lambda or alloc-like context.
// Sort/min/max key= lambdas, light attribute lambdas, and construct-then-return
// early-exit loop bodies are not per-element allocation smells.
func detectPYPERF25(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		text := line.text
		hasLambda := strings.Contains(text, "lambda:") || strings.Contains(text, "lambda ")
		ctor := ""
		if m := perfHeavyCtorRE.FindStringSubmatch(text); len(m) == 2 {
			ctor = m[1]
		}
		if ctor != "" && perfLightweightCtor(ctor) {
			ctor = ""
		}
		heavyLambda := hasLambda && !perfLambdaIsSortKey(text) && perfLambdaLooksHeavy(text)
		// Ctor/lambda needle before inLoop so non-matching lines skip membership.
		if ctor == "" && !heavyLambda {
			continue
		}
		if !facts.lineInLoop(i) {
			continue
		}
		if strings.Contains(text, "encode_cell") || strings.Contains(text, "encode_record") {
			continue
		}
		// Require alloc/schema context for bare constructors; heavy lambdas remain signal.
		if ctor != "" && !heavyLambda && !perfAllocContext(text) {
			continue
		}
		if heavyLambda && ctor == "" && perfLambdaFollowedByReturn(facts.lines, i) {
			continue
		}
		pushLine(unit, "PERF-PY-25", line, "alloc", "heavy object or lambda is constructed per homogeneous loop element; use a fixed-schema path", out)
	}
}

func perfLambdaIsSortKey(text string) bool {
	compact := strings.ReplaceAll(text, " ", "")
	return strings.Contains(compact, "key=lambda") ||
		strings.Contains(compact, "key=lambda:")
}

func perfLambdaLooksHeavy(text string) bool {
	idx := strings.Index(text, "lambda")
	if idx < 0 {
		return false
	}
	rest := text[idx:]
	colon := strings.IndexByte(rest, ':')
	if colon < 0 {
		return true
	}
	body := strings.TrimSpace(rest[colon+1:])
	if body == "" {
		return false
	}
	// Strip trailing call-arg punctuation from being inside sorted(...).
	body = strings.TrimRight(body, ",)")
	body = strings.TrimSpace(body)
	// Attribute/name-only bodies are closures, not heavy allocs.
	for _, r := range body {
		if r == '(' || r == '[' || r == '{' {
			return true
		}
	}
	return false
}

func perfLambdaFollowedByReturn(lines []codeLine, idx int) bool {
	if idx < 0 || idx >= len(lines) {
		return false
	}
	indent := indentWidth(lines[idx].raw)
	for j := idx + 1; j < len(lines) && j <= idx+4; j++ {
		t := lines[j].trim
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind < indent {
			return false
		}
		if ind == indent && strings.HasPrefix(t, "return ") {
			return true
		}
		if ind == indent {
			return false
		}
	}
	return false
}

func perfLightweightCtor(name string) bool {
	switch name {
	case "Exception", "ValueError", "TypeError", "KeyError", "RuntimeError", "HTTPException",
		"PermissionError", "OSError", "IOError", "StopIteration", "AssertionError",
		"Dict", "List", "Set", "Tuple", "Optional", "Any":
		return true
	default:
		return false
	}
}

func perfAllocContext(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "parent=") || strings.Contains(lower, ".meta") ||
		strings.Contains(lower, "encode") || strings.Contains(lower, "serialize") ||
		strings.Contains(lower, "buffer") || strings.Contains(lower, "schema") ||
		strings.Contains(lower, "payload") || strings.Contains(lower, "workitem") ||
		strings.Contains(lower, "record") || strings.Contains(lower, "document")
}

// PERF-PY-26: expensive decode on render/job path without cache.
// parse_* recursive-descent / CLI helpers only fire on explicit handle_job /
// handle_request windows; decode_/Image.open/zlib still fire in loops.
// Lightweight wire-field codecs (decode_field_value, header latin-1) are not
// expensive blob decodes — niquests WASI adapter trailer/header loops.
func detectPYPERF26(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !decodeHotRE.MatchString(line.text) {
			continue
		}
		if perfLineIsLightweightDecode(line.text) {
			continue
		}
		if strings.Contains(facts.Source, "_BLOB_CACHE") || strings.Contains(facts.Source, "_DECODE_CACHE") ||
			strings.Contains(facts.Source, "lru_cache") {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "_CACHE", "lru_cache", "cache.get", "cached =") {
			continue
		}
		hotHandler := windowHas(facts.lines, start, end, "handle_job", "handle_request")
		inLoop := facts.lineInLoop(i)
		isParseOnly := perfLineIsParseCall(line.text) && !perfLineIsDecodeOrImage(line.text)
		if isParseOnly {
			if !hotHandler {
				continue
			}
		} else if !hotHandler && !inLoop {
			continue
		}
		pushLine(unit, "PERF-PY-26", line, "decode", "expensive decode/parse runs on a hot path without a visible cache", out)
	}
}

var (
	decodeHotParseRE = regexp.MustCompile(`\bparse_[A-Za-z0-9_]+\s*\(`)
	decodeHotBlobRE  = regexp.MustCompile(`\b(decode_[A-Za-z0-9_]+|Image\.open|zlib\.decompress)\s*\(`)
	// HTTP header/trailer field codecs — O(header size) latin-1, not blob/image work.
	lightweightDecodeRE = regexp.MustCompile(
		`\bdecode_(?:field(?:_value)?|header(?:_value)?|trailer(?:_value)?|name|token)\s*\(`,
	)
)

func perfLineIsParseCall(text string) bool {
	return decodeHotParseRE.MatchString(text)
}

func perfLineIsDecodeOrImage(text string) bool {
	return decodeHotBlobRE.MatchString(text)
}

// perfLineIsLightweightDecode reports header/field wire codecs that match
// decode_* by name but are not expensive binary expansion.
func perfLineIsLightweightDecode(text string) bool {
	return lightweightDecodeRE.MatchString(text)
}

// PERF-PY-27: from_file / read_bytes of the same invariant path inside a batch loop.
// Requires the load site to be in-loop or in a function directly called from a loop,
// and that the path expression is not the loop variable or a callee parameter
// (once-per-distinct-path loads are not "repeated same path").
//
// Seed/migrate helpers are skipped (not hot request paths). Unique-path /
// loop-derived receivers still suppress once-per-input ETL shapes.
//
// Intermediate path bindings derived from the loop variable
// (validation_path = triad_paths(rule.festival_id)["validation"]; path.read_text())
// are once-per-distinct-path, not repeated loads of one invariant file.
func detectPYPERF27(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !fromFileRE.MatchString(line.text) {
			continue
		}
		if strings.Contains(facts.Source, "_CONFIG_CACHE") || strings.Contains(facts.Source, "_FILE_CACHE") ||
			strings.Contains(facts.Source, "lru_cache") {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "_CACHE", "lru_cache", "cache.get", "cached =") {
			continue
		}
		if !(strings.Contains(line.text, "read_text") || strings.Contains(line.text, "read_bytes") ||
			strings.Contains(line.text, "from_file") || windowHas(facts.lines, start, end, "CONFIG_PATH", "Path(")) {
			continue
		}
		inLoop := facts.lineInLoop(i)
		fromLoop := functionCalledFromLoop(facts, start)
		if !inLoop && !fromLoop {
			continue
		}
		if inLoop {
			if loopVar, ok := enclosingForLoopVar(facts.lines, i); ok {
				if perfLoadLineUsesName(line.text, loopVar) {
					continue
				}
				// Receiver assigned from an expression involving the loop var
				// (transitively) earlier in the same loop body.
				if perfLoadReceiverDerivedFromLoopVar(facts.lines, i, loopVar) {
					continue
				}
			}
		} else if fromLoop {
			params := pythonDefParamNamesMulti(facts.lines, start)
			if perfLoadLineUsesAnyName(line.text, params) {
				continue
			}
		}
		pushLine(unit, "PERF-PY-27", line, "load", "same path is loaded/parsed without a visible cache; reuse immutable parse results", out)
	}
}

// perfLoadReceiverDerivedFromLoopVar is true when the load receiver (or a
// subscript base) was assigned above in the same loop body from an expression
// that mentions the loop variable — e.g. paths = triad_paths(rule.id).
// Transitive: txt_path = out_base.with_suffix(".txt") where out_base = tmp / img.stem.
func perfLoadReceiverDerivedFromLoopVar(lines []codeLine, idx int, loopVar string) bool {
	if loopVar == "" || idx <= 0 || idx >= len(lines) {
		return false
	}
	recv := perfLoadReceiverName(lines[idx].text)
	if recv == "" {
		return false
	}
	// Collect names derived from loopVar within the enclosing loop body.
	derived := map[string]bool{loopVar: true}
	// Walk forward from the loop header to idx so derivation order is correct.
	loopStart := -1
	indent := indentWidth(lines[idx].raw)
	for j := idx - 1; j >= 0; j-- {
		t := lines[j].trim
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind < indent && isLoopTrim(t) {
			loopStart = j
			break
		}
		if ind < indent && isLoopTrim(t) {
			break
		}
	}
	if loopStart < 0 {
		// Fallback: single-pass reverse one-level check.
		return perfAssignMentionsAny(lines, idx, recv, derived)
	}
	for j := loopStart + 1; j < idx; j++ {
		lhs, rhs, ok := splitSimpleAssign(lines[j].text)
		if !ok {
			continue
		}
		if perfLoadLineUsesAnyName(rhs, mapKeys(derived)) {
			derived[lhs] = true
		}
	}
	return derived[recv]
}

func mapKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func perfAssignMentionsAny(lines []codeLine, idx int, recv string, seeds map[string]bool) bool {
	indent := indentWidth(lines[idx].raw)
	for j := idx - 1; j >= 0; j-- {
		t := lines[j].trim
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind < indent && isLoopTrim(t) {
			break
		}
		lhs, rhs, ok := splitSimpleAssign(lines[j].text)
		if !ok || lhs != recv {
			continue
		}
		return perfLoadLineUsesAnyName(rhs, mapKeys(seeds))
	}
	return false
}

// pythonDefParamNamesMulti joins a multi-line def header so typed params on
// following lines (overrides_path: Path,) are visible to from-loop suppressions.
func pythonDefParamNamesMulti(lines []codeLine, start int) []string {
	if start < 0 || start >= len(lines) {
		return nil
	}
	header := lines[start].text
	// Accumulate until balanced parens / def line ends with ):
	depth := 0
	sawOpen := false
	for j := start; j < len(lines) && j <= start+12; j++ {
		if j > start {
			header += " " + lines[j].text
		}
		for _, r := range lines[j].text {
			switch r {
			case '(':
				depth++
				sawOpen = true
			case ')':
				if depth > 0 {
					depth--
				}
			}
		}
		if sawOpen && depth == 0 {
			break
		}
	}
	return pythonDefParamNames(header)
}

// perfLoadReceiverName extracts the identifier receiving .read_text/.read_bytes
// (handles name.read_text, name["k"].read_text, name['k'].read_text).
func perfLoadReceiverName(line string) string {
	for _, meth := range []string{".read_text", ".read_bytes", ".from_file"} {
		at := strings.Index(line, meth)
		if at <= 0 {
			continue
		}
		left := strings.TrimSpace(line[:at])
		// Strip trailing subscript: name["validation"] or name['validation']
		for {
			if i := strings.LastIndex(left, "["); i >= 0 && strings.HasSuffix(left, "]") {
				left = strings.TrimSpace(left[:i])
				continue
			}
			break
		}
		// Trailing attribute chain: take rightmost bare ident.
		if dot := strings.LastIndex(left, "."); dot >= 0 {
			left = strings.TrimSpace(left[dot+1:])
		}
		left = strings.TrimSpace(left)
		if re := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`); re.MatchString(left) {
			return left
		}
	}
	// json.loads(path.read_text(...)) — still match path via nested call.
	if at := strings.Index(line, ".read_text"); at > 0 {
		// find identifier immediately left of .read_text skipping ]
		i := at - 1
		for i >= 0 && (line[i] == ']' || line[i] == '"' || line[i] == '\'' || line[i] == ' ') {
			// walk back through ["validation"]
			if line[i] == ']' {
				// find matching [
				depth := 1
				i--
				for i >= 0 && depth > 0 {
					if line[i] == ']' {
						depth++
					} else if line[i] == '[' {
						depth--
					}
					i--
				}
				continue
			}
			i--
		}
		end := i + 1
		for i >= 0 {
			c := line[i]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
				i--
				continue
			}
			break
		}
		name := strings.TrimSpace(line[i+1 : end])
		if name != "" {
			return name
		}
	}
	return ""
}

func splitSimpleAssign(line string) (lhs, rhs string, ok bool) {
	// Prefer '=' not '==', '!=', '<=', '>='.
	for i := 0; i < len(line); i++ {
		if line[i] != '=' {
			continue
		}
		if i > 0 {
			prev := line[i-1]
			if prev == '=' || prev == '!' || prev == '<' || prev == '>' {
				continue
			}
		}
		if i+1 < len(line) && line[i+1] == '=' {
			continue
		}
		lhs = strings.TrimSpace(line[:i])
		rhs = strings.TrimSpace(line[i+1:])
		// Only bare identifier LHS (not attributes).
		if regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(lhs) {
			return lhs, rhs, true
		}
		return "", "", false
	}
	return "", "", false
}

func functionCalledFromLoop(facts *pyPerfFacts, defLineIdx int) bool {
	if facts == nil || defLineIdx < 0 || defLineIdx >= len(facts.lines) {
		return false
	}
	lines := facts.lines
	header := lines[defLineIdx].trim
	var name string
	switch {
	case strings.HasPrefix(header, "async def "):
		name = strings.TrimSpace(strings.TrimPrefix(header, "async def "))
	case strings.HasPrefix(header, "def "):
		name = strings.TrimSpace(strings.TrimPrefix(header, "def "))
	default:
		return false
	}
	if at := strings.IndexByte(name, '('); at >= 0 {
		name = name[:at]
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	needle := name + "("
	for i, line := range lines {
		if i == defLineIdx {
			continue
		}
		if strings.Contains(line.text, needle) && facts.lineInLoop(i) {
			return true
		}
	}
	return false
}

func enclosingForLoopVar(lines []codeLine, idx int) (string, bool) {
	if idx < 0 || idx >= len(lines) {
		return "", false
	}
	indent := indentWidth(lines[idx].raw)
	for i := idx - 1; i >= 0; i-- {
		t := lines[i].trim
		if t == "" {
			continue
		}
		if indentWidth(lines[i].raw) >= indent {
			continue
		}
		if !isLoopTrim(t) {
			continue
		}
		if loopVar, _, ok := perfLoopBinding(lines[i].text); ok {
			return loopVar, true
		}
		return "", false
	}
	return "", false
}

func pythonDefParamNames(header string) []string {
	trimmed := strings.TrimSpace(header)
	trimmed = strings.TrimPrefix(trimmed, "async def ")
	trimmed = strings.TrimPrefix(trimmed, "def ")
	open := strings.IndexByte(trimmed, '(')
	closeParen := strings.LastIndexByte(trimmed, ')')
	if open < 0 || closeParen <= open {
		return nil
	}
	raw := trimmed[open+1 : closeParen]
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "*" || part == "/" {
			continue
		}
		part = strings.TrimPrefix(part, "**")
		part = strings.TrimPrefix(part, "*")
		part = strings.TrimSpace(part)
		if at := strings.IndexAny(part, ":="); at >= 0 {
			part = strings.TrimSpace(part[:at])
		}
		if part == "" || part == "self" || part == "cls" {
			continue
		}
		out = append(out, part)
	}
	return out
}

func perfLoadLineUsesAnyName(line string, names []string) bool {
	for _, name := range names {
		if perfLoadLineUsesName(line, name) {
			return true
		}
	}
	return false
}

func perfLoadLineUsesName(line, name string) bool {
	if name == "" {
		return false
	}
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`)
	return re.MatchString(line)
}

// PERF-PY-28: ThreadPoolExecutor constructed inside the unit of work.
func detectPYPERF28(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !executorRE.MatchString(line.text) {
			continue
		}
		if indentWidth(line.raw) == 0 {
			continue
		}
		start, _ := functionWindow(facts.lines, i)
		trimmed := facts.lines[start].trim
		if (strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ")) && i > start {
			pushLine(unit, "PERF-PY-28", line, "Executor", "executor/pool is created per unit of work; reuse a process-lifetime pool", out)
		}
	}
}

// PERF-PY-29: materialize all bodies then allocate full bytearray(estimate).
func detectPYPERF29(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		hit := bytearrayEstimateRE.MatchString(line.text) ||
			strings.Contains(line.text, "bytearray(estimate)") ||
			(strings.Contains(line.text, "sum(len(") && strings.Contains(line.text, "bodies"))
		if !hit {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "bodies.append", "chunks.append") &&
			!windowHas(facts.lines, start, end, "chunks.clear", "bodies.clear") {
			// Safe fixture clears chunks[i] = b"" — still has clear intent via clear()
			// Vulnerable has neither clear nor empty-body assign pattern alone is weak;
			// require no .clear() call.
			pushLine(unit, "PERF-PY-29", line, "bytearray", "full intermediate bodies are retained while allocating the full output buffer", out)
		}
	}
}

// PERF-PY-30: str/bytes + inside nested loops.
func detectPYPERF30(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		text := line.trim
		isConcat := strings.Contains(text, "+=") ||
			(strings.Contains(text, " = ") && strings.Contains(text, " + "))
		if !isConcat {
			continue
		}
		if !strings.Contains(text, "b\"") && !strings.Contains(text, "b'") &&
			!strings.Contains(text, ".encode(") && !strings.Contains(text, "+ b") {
			continue
		}
		if nestedForDepth(facts.lines, i) >= 2 {
			pushLine(unit, "PERF-PY-30", line, "+=", "string/bytes fragments are assembled with + inside nested loops; append and join once per page/row", out)
		}
	}
}

func nestedForDepth(lines []codeLine, idx int) int {
	if idx < 0 || idx >= len(lines) {
		return 0
	}
	targetIndent := indentWidth(lines[idx].raw)
	depth := 0
	for i := idx - 1; i >= 0; i-- {
		trimmed := lines[i].trim
		if trimmed == "" {
			continue
		}
		ind := indentWidth(lines[i].raw)
		if ind >= targetIndent {
			continue
		}
		if strings.HasPrefix(trimmed, "for ") && strings.HasSuffix(trimmed, ":") {
			depth++
			targetIndent = ind
			continue
		}
		if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") || strings.HasPrefix(trimmed, "class ") {
			break
		}
		if strings.HasSuffix(trimmed, ":") {
			targetIndent = ind
		}
	}
	return depth
}

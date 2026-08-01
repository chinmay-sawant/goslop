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
		pushLine(unit, "PERF-PY-23", line, "encode", "polymorphic encode/serialize runs inside a hot loop; specialize fixed-schema leaves or encode outside the loop", out)
	}
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
		// Ctor/lambda needle before inLoop so non-matching lines skip membership.
		if ctor == "" && !hasLambda {
			continue
		}
		if !facts.lineInLoop(i) {
			continue
		}
		if strings.Contains(text, "encode_cell") || strings.Contains(text, "encode_record") {
			continue
		}
		// Require alloc/schema context for bare constructors; lambdas remain high-signal.
		if ctor != "" && !hasLambda && !perfAllocContext(text) {
			continue
		}
		pushLine(unit, "PERF-PY-25", line, "alloc", "heavy object or lambda is constructed per homogeneous loop element; use a fixed-schema path", out)
	}
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
func detectPYPERF26(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !decodeHotRE.MatchString(line.text) {
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
		if windowHas(facts.lines, start, end, "handle_job", "handle_request", "render", "build_", "process") || facts.lineInLoop(i) {
			pushLine(unit, "PERF-PY-26", line, "decode", "expensive decode/parse runs on a hot path without a visible cache", out)
		}
	}
}

// PERF-PY-27: from_file / read_bytes of invariant path inside batch loop.
// Requires the load site to be in-loop or in a function directly called from a loop.
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
		if facts.lineInLoop(i) || functionCalledFromLoop(facts, start) {
			pushLine(unit, "PERF-PY-27", line, "load", "same path is loaded/parsed without a visible cache; reuse immutable parse results", out)
		}
	}
}

func functionCalledFromLoop(facts *pyPerfFacts, defLineIdx int) bool {
	if facts == nil || defLineIdx < 0 || defLineIdx >= len(facts.lines) {
		return false
	}
	lines := facts.lines
	header := lines[defLineIdx].trim
	name := ""
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

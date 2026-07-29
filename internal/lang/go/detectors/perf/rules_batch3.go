package perf

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

func init() {
	RegisterRule("PERF-112", detectPERF112, &MetaPERF112)
	RegisterRule("PERF-113", detectPERF113, &MetaPERF113)
	RegisterRule("PERF-114", detectPERF114, &MetaPERF114)
	RegisterRule("PERF-115", detectPERF115, &MetaPERF115)
	RegisterRule("PERF-117", detectPERF117, &MetaPERF117)
	RegisterRule("PERF-118", detectPERF118, &MetaPERF118)
	RegisterRule("PERF-119", detectPERF119, &MetaPERF119)
	RegisterRule("PERF-120", detectPERF120, &MetaPERF120)
	RegisterRule("PERF-121", detectPERF121, &MetaPERF121)
	RegisterRule("PERF-122", detectPERF122, &MetaPERF122)
	RegisterRule("PERF-123", detectPERF123, &MetaPERF123)
	RegisterRule("PERF-124", detectPERF124, &MetaPERF124)
	RegisterRule("PERF-125", detectPERF125, &MetaPERF125)
	RegisterRule("PERF-126", detectPERF126, &MetaPERF126)
	RegisterRule("PERF-127", detectPERF127, &MetaPERF127)
	RegisterRule("PERF-128", detectPERF128, &MetaPERF128)
	RegisterRule("PERF-129", detectPERF129, &MetaPERF129)
	RegisterRule("PERF-130", detectPERF130, &MetaPERF130)
	RegisterRule("PERF-131", detectPERF131, &MetaPERF131)
	RegisterRule("PERF-132", detectPERF132, &MetaPERF132)
	RegisterRule("PERF-133", detectPERF133, &MetaPERF133)
	RegisterRule("PERF-134", detectPERF134, &MetaPERF134)
	RegisterRule("PERF-135", detectPERF135, &MetaPERF135)
	RegisterRule("PERF-137", detectPERF137, &MetaPERF137)
	RegisterRule("PERF-138", detectPERF138, &MetaPERF138)
	RegisterRule("PERF-139", detectPERF139, &MetaPERF139)
	RegisterRule("PERF-140", detectPERF140, &MetaPERF140)
	RegisterRule("PERF-141", detectPERF141, &MetaPERF141)
	RegisterRule("PERF-142", detectPERF142, &MetaPERF142)
	RegisterRule("PERF-143", detectPERF143, &MetaPERF143)
	RegisterRule("PERF-144", detectPERF144, &MetaPERF144)
	RegisterRule("PERF-145", detectPERF145, &MetaPERF145)
	RegisterRule("PERF-146", detectPERF146, &MetaPERF146)
	RegisterRule("PERF-147", detectPERF147, &MetaPERF147)
	RegisterRule("PERF-148", detectPERF148, &MetaPERF148)
	RegisterRule("PERF-149", detectPERF149, &MetaPERF149)
	RegisterRule("PERF-150", detectPERF150, &MetaPERF150)
	RegisterRule("PERF-151", detectPERF151, &MetaPERF151)
	RegisterRule("PERF-152", detectPERF152, &MetaPERF152)
	RegisterRule("PERF-153", detectPERF153, &MetaPERF153)
	RegisterRule("PERF-154", detectPERF154, &MetaPERF154)
	RegisterRule("PERF-155", detectPERF155, &MetaPERF155)
	RegisterRule("PERF-156", detectPERF156, &MetaPERF156)
	RegisterRule("PERF-157", detectPERF157, &MetaPERF157)
	RegisterRule("PERF-158", detectPERF158, &MetaPERF158)
	RegisterRule("PERF-159", detectPERF159, &MetaPERF159)
	RegisterRule("PERF-160", detectPERF160, &MetaPERF160)
	RegisterRule("PERF-161", detectPERF161, &MetaPERF161)
	RegisterRule("PERF-162", detectPERF162, &MetaPERF162)
	RegisterRule("PERF-163", detectPERF163, &MetaPERF163)
}

// ---- helpers (batch3-local) ----

func b3windowAround(src string, start, before, after int) string {
	if start < 0 {
		start = 0
	}
	lo := start - before
	if lo < 0 {
		lo = 0
	}
	hi := start + after
	if hi > len(src) {
		hi = len(src)
	}
	return src[lo:hi]
}

func b3lineWindow(src string, start int) string {
	if start < 0 {
		start = 0
	}
	if start > len(src) {
		start = len(src)
	}
	lo := strings.LastIndexByte(src[:start], '\n') + 1
	hi := strings.IndexByte(src[start:], '\n')
	if hi < 0 {
		return src[lo:]
	}
	return src[lo : start+hi]
}

func b3countSubstring(s, sub string) int {
	if sub == "" {
		return 0
	}
	return strings.Count(s, sub)
}

func b3hasHandlerSig(src string) bool {
	return strings.Contains(src, "http.ResponseWriter") ||
		strings.Contains(src, "*http.Request") ||
		strings.Contains(src, "*gin.Context") ||
		strings.Contains(src, "echo.Context") ||
		strings.Contains(src, "*fiber.Ctx")
}

var reMakeZero = regexp.MustCompile(`make\(\s*(\[[^\]]+\]\w*|map\[[^\]]+\]\w*)\s*,\s*0(\s*,|\s*\))`)
var reIife = regexp.MustCompile(`func\s*\(\s*\)\s*\{[^}]*\}\s*\(\s*\)`)
var reLargeArray = regexp.MustCompile(`var\s+\w+\s+\[\d{3,}\]`)
var reArrayDecl = regexp.MustCompile(`\[\d{4,}\]`)

// detectPERF112: strings.ToLower/ToUpper used in a comparison (Rust: EqualFold).
func detectPERF112(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	// Rust: each ToLower/ToUpper call with ==/!= in a nearby window fires.
	file := unitFile(unit)
	emitted := 0
	for _, call := range facts.Calls {
		if call.Callee != "strings.ToLower" && call.Callee != "strings.ToUpper" {
			continue
		}
		start := call.StartByte - 8
		if start < 0 {
			start = 0
		}
		end := call.StartByte + 96
		if end > len(unit.Source) {
			end = len(unit.Source)
		}
		window := unit.Source[start:end]
		if !strings.Contains(window, "==") && !strings.Contains(window, "!=") {
			continue
		}
		lineN, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF112, file, lineN, col,
			"case conversion before comparison allocates; use strings.EqualFold", out)
		emitted++
		if emitted >= 12 {
			return
		}
	}
}

// detectPERF113: select with a single case and no default.
// Rust parity: scan for "select {", take text until the first '}', count
// substring "case " / presence of "default:". Nested braces in case bodies
// intentionally truncate the window (matches tree-sitter-free Rust heuristic).
func detectPERF113(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "select {") {
		return
	}
	file := unitFile(unit)
	searchFrom := 0
	for {
		rel := strings.Index(src[searchFrom:], "select {")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		endRel := strings.Index(src[start:], "}")
		end := len(src)
		if endRel >= 0 {
			end = start + endRel
		}
		block := src[start:end]
		if strings.Count(block, "case ") == 1 && !strings.Contains(block, "default:") {
			lineN, col := unit.LineCol(start)
			rules.PushFinding(&MetaPERF113, file, lineN, col,
				"single-case select should be a direct channel send or receive", out)
		}
		searchFrom = end + 1
		if searchFrom >= len(src) {
			return
		}
	}
}

func b3isIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func b3extractBraceBlock(src string, open int) (string, bool) {
	if open < 0 || open >= len(src) || src[open] != '{' {
		return "", false
	}
	depth := 0
	for i := open; i < len(src); i++ {
		switch src[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return src[open : i+1], true
			}
		}
	}
	return "", false
}

// detectPERF114: manual for-range copy instead of copy()
// Rust parity: for i, v := range src { dst[i] = v } (both bindings non-blank).
func detectPERF114(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "range") || !strings.Contains(src, "] =") && !strings.Contains(src, "]=") {
		return
	}
	file := unitFile(unit)
	lines := strings.Split(src, "\n")
	for li, ln := range lines {
		t := strings.TrimSpace(ln)
		if !strings.HasPrefix(t, "for ") || !strings.Contains(t, "range") {
			continue
		}
		// Require two-binding range: for idx, val := range …
		head := t
		rangeIdx := strings.Index(head, "range")
		if rangeIdx < 0 {
			continue
		}
		bindings := strings.TrimSpace(head[len("for "):rangeIdx])
		if i := strings.Index(bindings, ":="); i >= 0 {
			bindings = strings.TrimSpace(bindings[:i])
		}
		if i := strings.Index(bindings, "="); i >= 0 {
			bindings = strings.TrimSpace(bindings[:i])
		}
		parts := strings.Split(bindings, ",")
		if len(parts) != 2 {
			continue
		}
		idx := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if idx == "" || idx == "_" || val == "" || val == "_" {
			continue
		}
		if !b3isSimpleIdent(idx) || !b3isSimpleIdent(val) {
			continue
		}
		// Collect single-statement body until closing brace.
		var bodyLines []string
		for j := li + 1; j < len(lines) && j < li+8; j++ {
			body := strings.TrimSpace(lines[j])
			if body == "}" || strings.HasPrefix(body, "}") {
				break
			}
			if body == "" || strings.HasPrefix(body, "//") {
				continue
			}
			bodyLines = append(bodyLines, body)
		}
		if len(bodyLines) != 1 {
			continue
		}
		body := bodyLines[0]
		// dst[idx] = val  (pure copy)
		eq := strings.Index(body, "=")
		if eq < 0 {
			continue
		}
		lhs := strings.TrimSpace(body[:eq])
		rhs := strings.TrimSpace(body[eq+1:])
		if rhs != val && !strings.HasPrefix(rhs, val+" ") {
			// allow simple ident RHS only
			if !b3isSimpleIdent(rhs) {
				continue
			}
		}
		open := strings.Index(lhs, "[")
		close := strings.LastIndex(lhs, "]")
		if open < 0 || close <= open {
			continue
		}
		dst := strings.TrimSpace(lhs[:open])
		indexExpr := strings.TrimSpace(lhs[open+1 : close])
		if !b3isSimpleIdent(dst) || indexExpr != idx {
			continue
		}
		if rhs != val {
			continue
		}
		// Suppress interface-slice destinations.
		if destinationIsInterfaceSlice(src, dst) {
			continue
		}
		off := 0
		for k := 0; k < li; k++ {
			off += len(lines[k]) + 1
		}
		lineN, col := unit.LineCol(off)
		rules.PushFinding(&MetaPERF114, file, lineN, col,
			"manual for-range copy should be the copy() builtin (3-10x faster, handles overlap)", out)
	}
}

func b3isSimpleIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
				return false
			}
			continue
		}
		if r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

func destinationIsInterfaceSlice(src, dst string) bool {
	markers := []string{
		dst + " := make([]interface{}",
		dst + " = make([]interface{}",
		dst + " := make([]any,",
		dst + " = make([]any,",
		"var " + dst + " []interface{}",
		"var " + dst + " []any",
		dst + " := []interface{}",
		dst + " = []interface{}",
		dst + " := []any{",
		dst + " = []any{",
	}
	for _, m := range markers {
		if strings.Contains(src, m) {
			return true
		}
	}
	return false
}

// detectPERF115: strings.Compare == 0 / != 0
func detectPERF115(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	for _, call := range facts.Calls {
		if call.Callee != "strings.Compare" {
			continue
		}
		line := b3lineWindow(unit.Source, call.StartByte)
		if strings.Contains(line, "== 0") || strings.Contains(line, "!= 0") ||
			strings.Contains(line, "==0") || strings.Contains(line, "!=0") {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF115, unitFile(unit), lineN, col,
				"strings.Compare used for equality; prefer == / !=", out)
			return
		}
	}
}

// detectPERF117: bytes.Compare == 0 / != 0
func detectPERF117(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	for _, call := range facts.Calls {
		if call.Callee != "bytes.Compare" {
			continue
		}
		line := b3lineWindow(unit.Source, call.StartByte)
		if strings.Contains(line, "== 0") || strings.Contains(line, "!= 0") ||
			strings.Contains(line, "==0") || strings.Contains(line, "!=0") {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF117, unitFile(unit), lineN, col,
				"bytes.Compare used for equality; prefer bytes.Equal", out)
			return
		}
	}
}

// detectPERF118: http.NewRequest("GET"|"HEAD"|"POST", ...)
func detectPERF118(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "http.NewRequest" && call.Callee != "http.NewRequestWithContext" {
			continue
		}
		if len(call.Arguments) == 0 {
			continue
		}
		method := strings.Trim(call.Arguments[0], `"`)
		// for NewRequestWithContext method is arg1
		if call.Callee == "http.NewRequestWithContext" && len(call.Arguments) > 1 {
			method = strings.Trim(call.Arguments[1], `"`)
		}
		switch method {
		case "GET", "HEAD", "POST":
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF118, file, lineN, col,
				"http.NewRequest for simple method; consider http.Get/Head/Post", out)
			return
		}
	}
}

// detectPERF119: 2+ consecutive appends to same slice (Rust call-fact parity).
func detectPERF119(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	b3detectConsecutiveAppendsFacts(unit, facts, out, 2, &MetaPERF119,
		"consecutive append calls to the same slice can be merged into one variadic append")
}

// detectPERF128: 3+ consecutive appends
func detectPERF128(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	b3detectConsecutiveAppendsFacts(unit, facts, out, 3, &MetaPERF128,
		"three or more independent append calls can be combined into one variadic append")
}

func b3detectConsecutiveAppendsFacts(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding, min int, meta *rules.RuleMetadata, msg string) {
	file := unitFile(unit)
	var appends []CallFact
	for _, c := range facts.Calls {
		if c.Callee == "append" {
			appends = append(appends, c)
		}
	}
	if len(appends) < min {
		// No call-fact coverage — last-resort line scan (fixtures / broken AST).
		b3detectConsecutiveAppendsLines(unit, out, min, meta, msg)
		return
	}
	// sort by start
	for i := 0; i < len(appends); i++ {
		for j := i + 1; j < len(appends); j++ {
			if appends[j].StartByte < appends[i].StartByte {
				appends[i], appends[j] = appends[j], appends[i]
			}
		}
	}
	src := unit.Source
	// Rust parity (strings_bytes.rs / maps_and_slices.rs): first same-target
	// window of min appends with no intervening_read between consecutive pairs.
	// No maxGap — intermediate window starts at prev.start+64 (empty when close).
	for i := 0; i+min-1 < len(appends); i++ {
		target := ""
		if len(appends[i].Arguments) > 0 {
			target = strings.TrimSpace(appends[i].Arguments[0])
		}
		if target == "" {
			continue
		}
		ok := true
		for k := 1; k < min; k++ {
			a := appends[i+k]
			if len(a.Arguments) == 0 || strings.TrimSpace(a.Arguments[0]) != target {
				ok = false
				break
			}
			prev := appends[i+k-1]
			winStart := prev.StartByte + 64
			if winStart > a.StartByte {
				winStart = a.StartByte
			}
			if winStart < a.StartByte {
				win := src[winStart:a.StartByte]
				if b3interveningRead(win, target) {
					ok = false
					break
				}
			}
		}
		if !ok {
			continue
		}
		lineN, col := unit.LineCol(appends[i].StartByte)
		rules.PushFinding(meta, file, lineN, col, msg, out)
		return
	}
	// Call facts present but no matching run — do not line-fallback (avoids
	// signature.go-style FPs vs Rust when intervening_read already filtered).
}

func b3interveningRead(window, target string) bool {
	if window == "" {
		return false
	}
	// Exact Rust marker set from intervening_read in strings_bytes.rs.
	markers := []string{"(", target, ")", "len(", "range ", "copy("}
	for _, m := range markers {
		if m != "" && strings.Contains(window, m) {
			return true
		}
	}
	return false
}

func b3detectConsecutiveAppendsLines(unit *core.ParsedUnit, out *[]rules.Finding, min int, meta *rules.RuleMetadata, msg string) {
	src := unit.Source
	if !strings.Contains(src, "append(") {
		return
	}
	lines := strings.Split(src, "\n")
	var runVar string
	runCount := 0
	runStartLine := 0
	file := unitFile(unit)
	for li, ln := range lines {
		t := strings.TrimSpace(ln)
		if strings.Contains(t, "= append(") {
			lhs := strings.TrimSpace(strings.Split(t, "=")[0])
			lhs = strings.TrimSuffix(lhs, ":")
			if runVar == lhs {
				runCount++
			} else {
				runVar = lhs
				runCount = 1
				runStartLine = li
			}
			if runCount >= min {
				off := 0
				for k := 0; k < runStartLine; k++ {
					off += len(lines[k]) + 1
				}
				lineN, col := unit.LineCol(off)
				rules.PushFinding(meta, file, lineN, col, msg, out)
				return
			}
			continue
		}
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		runVar = ""
		runCount = 0
	}
}

// detectPERF120: time.Now().Sub(
func detectPERF120(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "time.Now()") {
		return
	}
	file := unitFile(unit)
	idx := strings.Index(src, "time.Now()")
	for idx >= 0 {
		win := b3windowAround(src, idx, 0, 48)
		if strings.Contains(win, ".Sub(") {
			lineN, col := unit.LineCol(idx)
			rules.PushFinding(&MetaPERF120, file, lineN, col,
				"time.Now().Sub(t); prefer time.Since(t)", out)
			return
		}
		next := strings.Index(src[idx+1:], "time.Now()")
		if next < 0 {
			break
		}
		idx = idx + 1 + next
	}
}

// detectPERF121: two consecutive same-shape struct literals where the second
// builds field-by-field from the first. Direct conversion T(x) would suffice.
// Rust parity: different type names, identical field sets, ≤2 newlines between
// literals, and every keyed field in the later literal reads binding.Field.
func detectPERF121(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "{") {
		return
	}
	file := unitFile(unit)
	lits := b3collectStructLiterals(src)
	if len(lits) < 2 {
		return
	}
	for i := 0; i+1 < len(lits); i++ {
		a, b := lits[i], lits[i+1]
		if a.typeName == b.typeName {
			continue
		}
		if !b3stringSlicesEqual(a.fields, b.fields) {
			continue
		}
		// Adjacent: at most 2 newlines between closing of a and opening of b.
		if a.end > b.start || a.end > len(src) || b.start > len(src) {
			continue
		}
		between := src[a.end:b.start]
		if strings.Count(between, "\n") > 2 {
			continue
		}
		binding, ok := b3literalBinding(src, a.start)
		if !ok {
			continue
		}
		if !b3literalFieldsReadFrom(src, b, binding) {
			continue
		}
		lineN, col := unit.LineCol(a.start)
		rules.PushFinding(&MetaPERF121, file, lineN, col,
			"struct literal copies another literal of the same shape; use a direct type conversion (T(x))", out)
		return
	}
}

type b3structLiteral struct {
	typeName string
	fields   []string
	start    int
	end      int
}

func b3stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func b3collectStructLiterals(source string) []b3structLiteral {
	var out []b3structLiteral
	for i := 0; i < len(source); i++ {
		if source[i] != '{' {
			continue
		}
		if i == 0 {
			continue
		}
		pre := strings.TrimRight(source[:i], " \t\n\r")
		// Preceding token must be a simple type name.
		j := len(pre)
		for j > 0 {
			r := rune(pre[j-1])
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				j--
				continue
			}
			break
		}
		name := pre[j:]
		if !b3isSimpleIdent(name) {
			continue
		}
		// Matching '}' (first close; fine for non-nested fixture literals).
		closeRel := strings.IndexByte(source[i:], '}')
		if closeRel < 0 {
			continue
		}
		body := source[i+1 : i+closeRel]
		fields := b3parseFieldList(body)
		if len(fields) == 0 {
			continue
		}
		out = append(out, b3structLiteral{
			typeName: name,
			fields:   fields,
			start:    i,
			end:      i + closeRel + 1,
		})
	}
	return out
}

func b3parseFieldList(body string) []string {
	var fields []string
	var current strings.Builder
	depth := 0
	for _, c := range body {
		switch c {
		case '(', '[', '{':
			depth++
			current.WriteRune(c)
		case ')', ']', '}':
			depth--
			current.WriteRune(c)
		case ',':
			if depth == 0 {
				if name, ok := b3fieldName(current.String()); ok {
					fields = append(fields, name)
				}
				current.Reset()
				continue
			}
			current.WriteRune(c)
		default:
			current.WriteRune(c)
		}
	}
	if name, ok := b3fieldName(current.String()); ok {
		fields = append(fields, name)
	}
	return fields
}

func b3fieldName(text string) (string, bool) {
	text = strings.TrimSpace(text)
	name, _, ok := strings.Cut(text, ":")
	if !ok {
		return "", false
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "", false
	}
	return name, true
}

// b3literalBinding returns the local name receiving a struct literal when it is
// the complete RHS of a short declaration (name := Type{...}).
func b3literalBinding(source string, literalStart int) (string, bool) {
	if literalStart <= 0 || literalStart > len(source) {
		return "", false
	}
	// Walk left over TypeName before '{'.
	pre := strings.TrimRight(source[:literalStart], " \t\n\r")
	// Drop type name.
	j := len(pre)
	for j > 0 {
		r := rune(pre[j-1])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			j--
			continue
		}
		break
	}
	pre = strings.TrimRight(pre[:j], " \t\n\r")
	// Expect ":="
	if !strings.HasSuffix(pre, ":=") {
		return "", false
	}
	pre = strings.TrimRight(pre[:len(pre)-2], " \t\n\r")
	// Binding is last identifier on the line.
	if nl := strings.LastIndexByte(pre, '\n'); nl >= 0 {
		pre = pre[nl+1:]
	}
	pre = strings.TrimSpace(pre)
	if !b3isSimpleIdent(pre) {
		return "", false
	}
	return pre, true
}

// b3literalFieldsReadFrom requires every keyed field in the target literal to
// read the corresponding field from binding (binding.Field).
func b3literalFieldsReadFrom(source string, lit b3structLiteral, binding string) bool {
	if lit.start+1 >= lit.end-1 || lit.end > len(source) {
		return false
	}
	body := source[lit.start+1 : lit.end-1]
	keyed := b3parseKeyedFieldValues(body)
	for _, field := range lit.fields {
		want := binding + "." + field
		found := false
		for _, kv := range keyed {
			if kv.name == field && kv.value == want {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return len(lit.fields) > 0
}

type b3keyedField struct {
	name  string
	value string
}

func b3parseKeyedFieldValues(body string) []b3keyedField {
	var fields []b3keyedField
	var current strings.Builder
	depth := 0
	var quote rune
	runes := []rune(body)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		if quote != 0 {
			current.WriteRune(ch)
			if ch == '\\' && i+1 < len(runes) {
				i++
				current.WriteRune(runes[i])
				continue
			}
			if ch == quote {
				quote = 0
			}
			continue
		}
		switch ch {
		case '"', '`':
			quote = ch
			current.WriteRune(ch)
		case '/':
			if i+1 < len(runes) && runes[i+1] == '/' {
				// Line comment — skip to newline.
				i++
				for i+1 < len(runes) && runes[i+1] != '\n' {
					i++
				}
				continue
			}
			if i+1 < len(runes) && runes[i+1] == '*' {
				i++
				for i+1 < len(runes) {
					i++
					if runes[i-1] == '*' && runes[i] == '/' {
						break
					}
				}
				continue
			}
			current.WriteRune(ch)
		case '(', '[', '{':
			depth++
			current.WriteRune(ch)
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
			current.WriteRune(ch)
		case ',':
			if depth == 0 {
				if kf, ok := b3keyedFieldValue(current.String()); ok {
					fields = append(fields, kf)
				}
				current.Reset()
				continue
			}
			current.WriteRune(ch)
		default:
			current.WriteRune(ch)
		}
	}
	if kf, ok := b3keyedFieldValue(current.String()); ok {
		fields = append(fields, kf)
	}
	return fields
}

func b3keyedFieldValue(field string) (b3keyedField, bool) {
	name, value, ok := strings.Cut(field, ":")
	if !ok {
		return b3keyedField{}, false
	}
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	if name == "" || value == "" {
		return b3keyedField{}, false
	}
	return b3keyedField{name: name, value: value}, true
}

// detectPERF122: HasPrefix + slice with len (Rust: window -64..+256, needle ":]")
func detectPERF122(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	src := unit.Source
	for _, call := range facts.Calls {
		if call.Callee != "strings.HasPrefix" {
			continue
		}
		start := call.StartByte
		if start > 64 {
			start -= 64
		} else {
			start = 0
		}
		end := call.StartByte + 256
		if end > len(src) {
			end = len(src)
		}
		win := src[start:end]
		if strings.Contains(win, "len(") && strings.Contains(win, ":]") {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF122, file, lineN, col,
				"strings.HasPrefix + slice should be strings.TrimPrefix", out)
		}
	}
}

// detectPERF123: make(T, 0) or make(T, 0, 0) — multi-fire; allow make(T, 0, cap>0).
func detectPERF123(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	// Prefer call facts when make is recorded.
	for _, call := range facts.Calls {
		if call.Callee != "make" {
			continue
		}
		args := call.Arguments
		if len(args) < 2 {
			continue
		}
		if strings.TrimSpace(args[1]) != "0" {
			continue
		}
		// Allow make(T, 0, cap) where cap != 0.
		if len(args) == 3 && strings.TrimSpace(args[2]) != "0" {
			continue
		}
		lineN, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF123, file, lineN, col,
			"make with explicit zero length/capacity is redundant; omit the zero argument", out)
	}
	// Text fallback if facts omit make.
	if len(facts.Calls) > 0 {
		return
	}
	src := unit.Source
	idx := strings.Index(src, "make(")
	for idx >= 0 {
		end := idx + 5
		depth := 1
		for end < len(src) && depth > 0 {
			switch src[end] {
			case '(':
				depth++
			case ')':
				depth--
			}
			end++
		}
		call := src[idx:end]
		// make(T, 0) or make(T, 0, 0)
		inner := strings.TrimSuffix(strings.TrimPrefix(call, "make("), ")")
		parts := splitTopLevelArgs(inner)
		if len(parts) >= 2 && strings.TrimSpace(parts[1]) == "0" {
			if len(parts) == 3 && strings.TrimSpace(parts[2]) != "0" {
				// capacity form — skip
			} else {
				lineN, col := unit.LineCol(idx)
				rules.PushFinding(&MetaPERF123, file, lineN, col,
					"make with explicit zero length/capacity is redundant; omit the zero argument", out)
			}
		}
		next := strings.Index(src[idx+1:], "make(")
		if next < 0 {
			break
		}
		idx = idx + 1 + next
	}
}

func splitTopLevelArgs(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	parts = append(parts, strings.TrimSpace(s[start:]))
	return parts
}

// detectPERF124 / 147: strings.Replace(..., -1)
func detectPERF124(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	b3detectReplaceNeg1(unit, facts, out, &MetaPERF124)
}

func detectPERF147(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	b3detectReplaceNeg1(unit, facts, out, &MetaPERF147)
}

func b3detectReplaceNeg1(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding, meta *rules.RuleMetadata) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "strings.Replace" && call.Callee != "bytes.Replace" {
			continue
		}
		for _, a := range call.Arguments {
			if strings.TrimSpace(a) == "-1" {
				lineN, col := unit.LineCol(call.StartByte)
				rules.PushFinding(meta, file, lineN, col,
					"Replace with -1; prefer ReplaceAll", out)
				return
			}
		}
	}
}

// detectPERF125: if s != nil { s = append(s, ...) }
func detectPERF125(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "append(") {
		return
	}
	if !strings.Contains(src, "!= nil") && !strings.Contains(src, "!=nil") {
		return
	}
	file := unitFile(unit)
	lines := strings.Split(src, "\n")
	for li, ln := range lines {
		t := strings.TrimSpace(ln)
		if !strings.Contains(t, "!= nil") && !strings.Contains(t, "!=nil") {
			continue
		}
		name := ""
		if strings.HasPrefix(t, "if ") {
			rest := strings.TrimPrefix(t, "if ")
			rest = strings.TrimSuffix(rest, "{")
			rest = strings.TrimSpace(rest)
			if strings.HasSuffix(rest, "!= nil") || strings.HasSuffix(rest, "!=nil") {
				fields := strings.Fields(rest)
				if len(fields) >= 1 {
					name = fields[0]
				}
			}
		}
		if name == "" {
			continue
		}
		for j := li + 1; j < len(lines) && j < li+6; j++ {
			b := lines[j]
			if strings.Contains(b, "append("+name+",") || strings.Contains(b, "append( "+name+",") {
				off := 0
				for k := 0; k < li; k++ {
					off += len(lines[k]) + 1
				}
				lineN, col := unit.LineCol(off)
				rules.PushFinding(&MetaPERF125, file, lineN, col,
					"redundant nil check before append", out)
				return
			}
			if strings.TrimSpace(b) == "}" {
				break
			}
		}
	}
}

// detectPERF126: CanonicalHeaderKey on common headers
func detectPERF126(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	canonical := map[string]bool{
		"Accept": true, "Accept-Encoding": true, "Accept-Language": true,
		"Authorization": true, "Cache-Control": true, "Connection": true,
		"Content-Length": true, "Content-Type": true, "Cookie": true,
		"Date": true, "Host": true, "Location": true, "Origin": true,
		"Referer": true, "Server": true, "Set-Cookie": true,
		"User-Agent": true, "X-Request-Id": true, "X-Forwarded-For": true,
	}
	for _, call := range facts.Calls {
		if call.Callee != "http.CanonicalHeaderKey" || len(call.Arguments) == 0 {
			continue
		}
		arg := strings.Trim(call.Arguments[0], `"`)
		if canonical[arg] {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF126, file, lineN, col,
				"CanonicalHeaderKey on already-canonical header name", out)
			return
		}
	}
}

// detectPERF127: log.*(fmt.Sprintf(...)) with no format verbs in sprintf
func detectPERF127(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "fmt.Sprintf") || !strings.Contains(src, "log.") {
		return
	}
	file := unitFile(unit)
	// log.X(fmt.Sprintf(...))
	idx := strings.Index(src, "log.")
	for idx >= 0 {
		win := b3windowAround(src, idx, 0, 120)
		if strings.Contains(win, "fmt.Sprintf") {
			// crude: if sprintf format has no %
			if m := regexp.MustCompile(`fmt\.Sprintf\(\s*"([^"]*)"`).FindStringSubmatch(win); m != nil {
				if !strings.Contains(m[1], "%") {
					lineN, col := unit.LineCol(idx)
					rules.PushFinding(&MetaPERF127, file, lineN, col,
						"fmt.Sprintf wrapping a static string inside log call", out)
					return
				}
			} else if strings.Contains(win, "fmt.Sprintf") {
				// still flag log.Printf(fmt.Sprintf(...))
				lineN, col := unit.LineCol(idx)
				rules.PushFinding(&MetaPERF127, file, lineN, col,
					"fmt.Sprintf wrapping inside log call", out)
				return
			}
		}
		next := strings.Index(src[idx+1:], "log.")
		if next < 0 {
			break
		}
		idx = idx + 1 + next
	}
}

// detectPERF129: for _, unusedVal := range with unused value
func detectPERF129(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "range") {
		return
	}
	file := unitFile(unit)
	// for _, name := range
	re := regexp.MustCompile(`for\s+(_|\w+)\s*,\s*(\w+)\s*:=\s*range\s+`)
	for _, m := range re.FindAllStringSubmatchIndex(src, -1) {
		valName := src[m[4]:m[5]]
		if valName == "_" {
			// underscore value is already optimal? Actually for _, _ is rare; for i, _ is PERF-156 territory
			// for _, v where v unused — key is _
			// if value is _ then value is not copied? In Go for _, v := range still? for _, _ is invalid-ish
			// Actually for i, _ := range still iterates values in some versions? No - with blank identifier value is not copied in modern Go for some cases
			// For PERF-129: value is named but unused
			continue
		}
		// find loop body
		brace := strings.Index(src[m[1]:], "{")
		if brace < 0 {
			continue
		}
		open := m[1] + brace
		block, ok := b3extractBraceBlock(src, open)
		if !ok {
			continue
		}
		body := block[1 : len(block)-1]
		// if valName appears as a word in body, used
		if b3identUsed(body, valName) {
			continue
		}
		lineN, col := unit.LineCol(m[0])
		rules.PushFinding(&MetaPERF129, file, lineN, col,
			"range binds value that is never used; use for i := range", out)
		return
	}
}

func b3identUsed(body, name string) bool {
	// word-boundary search
	for i := 0; i < len(body); {
		j := strings.Index(body[i:], name)
		if j < 0 {
			return false
		}
		pos := i + j
		beforeOK := pos == 0 || !b3isIdentByte(body[pos-1])
		after := pos + len(name)
		afterOK := after >= len(body) || !b3isIdentByte(body[after])
		if beforeOK && afterOK {
			return true
		}
		i = pos + len(name)
	}
	return false
}

// detectPERF130: immediately-invoked func literal whose body is a single call.
func detectPERF130(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "func()") {
		return
	}
	file := unitFile(unit)
	searchFrom := 0
	for {
		rel := strings.Index(src[searchFrom:], "func()")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		// Reject go/for/arg-list positions (Rust prev-char check).
		if start > 0 {
			prev := ' '
			for i := start - 1; i >= 0; i-- {
				c := rune(src[i])
				if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
					continue
				}
				prev = c
				break
			}
			if prev == 'o' || prev == 'f' || prev == ',' || prev == '(' {
				searchFrom = start + len("func()")
				continue
			}
		}
		windowEnd := start + 96
		if windowEnd > len(src) {
			windowEnd = len(src)
		}
		window := src[start:windowEnd]
		closeIdx := strings.Index(window, "}")
		if closeIdx < 0 {
			searchFrom = start + len("func()")
			continue
		}
		afterClose := strings.TrimLeft(window[closeIdx+1:], " \t\n\r")
		if !strings.HasPrefix(afterClose, "(") {
			searchFrom = start + len("func()")
			continue
		}
		bodyStart := strings.Index(window, "{")
		if bodyStart < 0 || bodyStart >= closeIdx {
			searchFrom = start + len("func()")
			continue
		}
		body := strings.TrimSpace(window[bodyStart+1 : closeIdx])
		// Single call expression: Foo(...) or pkg.Foo(...) optionally with trailing newline.
		if !b3isSingleCallExpression(body) {
			searchFrom = start + len("func()")
			continue
		}
		lineN, col := unit.LineCol(start)
		rules.PushFinding(&MetaPERF130, file, lineN, col,
			"unnecessary func() { f(args) }() wrapper; inline the call", out)
		return
	}
}

func b3isSingleCallExpression(body string) bool {
	body = strings.TrimSpace(body)
	if body == "" {
		return false
	}
	// Rust rejects `;` and control-flow keywords; multi-line without `;` is OK.
	if strings.Contains(body, ";") {
		return false
	}
	if strings.Contains(body, "if ") || strings.Contains(body, "for ") ||
		strings.Contains(body, "switch ") || strings.Contains(body, "select ") ||
		strings.Contains(body, "var ") || strings.Contains(body, "return ") {
		return false
	}
	open := strings.Index(body, "(")
	if open <= 0 {
		return false
	}
	// Prefix before first `(` must be a bare ident or method chain — no spaces
	// (rejects `_ = f.Close()` assignments).
	prefix := strings.TrimSpace(body[:open])
	if prefix == "" {
		return false
	}
	depth := 0
	lastWasDot := false
	for _, c := range prefix {
		switch {
		case c == '(':
			depth++
		case c == ')':
			depth--
		case c == '.' && depth == 0:
			lastWasDot = true
		case (c == ' ' || c == '\t') && depth == 0:
			return false
		}
	}
	if depth != 0 || strings.HasSuffix(prefix, ".") {
		return false
	}
	_ = lastWasDot
	return true
}

// detectPERF131: mu.Lock(); counter++; mu.Unlock()
func detectPERF131(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, ".Lock()") || !strings.Contains(src, ".Unlock()") {
		return
	}
	file := unitFile(unit)
	// between Lock and Unlock only simple ++/-- or += 1
	idx := strings.Index(src, ".Lock()")
	for idx >= 0 {
		after := src[idx+len(".Lock()"):]
		u := strings.Index(after, ".Unlock()")
		if u > 0 && u < 200 {
			body := after[:u]
			// strip whitespace
			compact := strings.Map(func(r rune) rune {
				if unicode.IsSpace(r) {
					return -1
				}
				return r
			}, body)
			// allow counter++ / counter-- / counter += 1 / counter = counter + 1
			if strings.Contains(compact, "++") || strings.Contains(compact, "--") ||
				strings.Contains(compact, "+=") || strings.Contains(compact, "-=") {
				// reject if method calls other than unlock
				if !strings.Contains(body, "(") {
					lineN, col := unit.LineCol(idx)
					rules.PushFinding(&MetaPERF131, file, lineN, col,
						"mutex guards a simple counter; prefer sync/atomic", out)
					return
				}
			}
		}
		next := strings.Index(src[idx+1:], ".Lock()")
		if next < 0 {
			break
		}
		idx = idx + 1 + next
	}
}

// detectPERF132: go func() { ... db.Query / http without ctx }
// detectPERF132: go func() that does cancellable I/O without accepting a context.
// Rust parity: parent has ctx context.Context; the go statement is go func()
// (no params); body performs I/O. Do not scan comments for "ctx" — that caused
// FNs when the vulnerable fixture documents missing context in a comment.
func detectPERF132(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !strings.Contains(src, "go func()") {
		return
	}
	// Parent must have something to forward; without a ctx param the warning is moot.
	if !strings.Contains(src, "ctx context.Context") {
		return
	}
	file := unitFile(unit)
	for _, gr := range facts.GoStarts {
		if gr[0] < 0 || gr[1] > len(src) || gr[0] >= gr[1] {
			continue
		}
		block := src[gr[0]:gr[1]]
		goText := block
		if len(goText) > 256 {
			goText = goText[:256]
		}
		// Accept only parameterless go func(); reject go func(ctx context.Context).
		if !strings.Contains(goText, "go func()") {
			continue
		}
		bodyStart := strings.Index(goText, "{")
		if bodyStart < 0 {
			continue
		}
		// Prefer full go-statement body for I/O (comments ignored for IO needles).
		fullBodyStart := strings.Index(block, "{")
		if fullBodyStart < 0 {
			continue
		}
		body := block[fullBodyStart+1:]
		if !b3bodyHasIO(body) {
			continue
		}
		lineN, col := unit.LineCol(gr[0])
		rules.PushFinding(&MetaPERF132, file, lineN, col,
			"go func() body makes I/O calls but the goroutine doesn't accept a context; cancellation cannot propagate", out)
		return
	}
}

func b3bodyHasIO(body string) bool {
	// Common packages whose calls take a context as the first argument.
	for _, p := range []string{
		"http.", "db.", "sql.", "redis.", "rdb.", "client.", "store.", "queue.", "kafka.",
	} {
		if strings.Contains(body, p) {
			return true
		}
	}
	return false
}

// detectPERF133: sort.Slice inside loop
func detectPERF133(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "sort.Slice" && call.Callee != "sort.SliceStable" {
			continue
		}
		if !IsInLoop(call) {
			continue
		}
		lineN, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF133, file, lineN, col,
			"sort.Slice closure allocated inside loop", out)
		return
	}
}

// detectPERF134: manual Read/Write loop
func detectPERF134(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "for ") {
		return
	}
	file := unitFile(unit)
	// for { ... Read ... Write ... }
	for _, fr := range facts.ForRanges {
		block := src[fr[0]:fr[1]]
		if (strings.Contains(block, ".Read(") || strings.Contains(block, "Read(")) &&
			(strings.Contains(block, ".Write(") || strings.Contains(block, "Write(")) {
			// not io.Copy
			if strings.Contains(block, "io.Copy") {
				continue
			}
			lineN, col := unit.LineCol(fr[0])
			rules.PushFinding(&MetaPERF134, file, lineN, col,
				"manual Read/Write loop; prefer io.Copy", out)
			return
		}
	}
}

// detectPERF135: gob.NewEncoder/Decoder inside loop
func detectPERF135(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "gob.NewEncoder" && call.Callee != "gob.NewDecoder" {
			continue
		}
		if !IsInLoop(call) {
			continue
		}
		lineN, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF135, file, lineN, col,
			"gob encoder/decoder created inside loop; reuse one instance", out)
		return
	}
}

// detectPERF137: runtime.Caller on a hot path / request handler.
// Rust parity: only fire when the 1 KiB window before the call looks like a
// request handler (or the call is in a loop). Whole-file handler presence
// must not flag package-level caches like `var tag = func(){ runtime.Caller(0) }()`.
func detectPERF137(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	src := unit.Source
	for _, call := range facts.Calls {
		if call.Callee != "runtime.Caller" {
			continue
		}
		if IsInLoop(call) || IsHandlerShaped(src, call.StartByte) {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF137, file, lineN, col,
				"runtime.Caller on hot path", out)
			return
		}
	}
}

// detectPERF138: runtime.Stack on hot path
func detectPERF138(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	src := unit.Source
	for _, call := range facts.Calls {
		if call.Callee != "runtime.Stack" {
			continue
		}
		if IsInLoop(call) || IsHotPath(src, call.StartByte, IsInLoop(call)) || b3hasHandlerSig(src) {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF138, file, lineN, col,
				"runtime.Stack on hot path", out)
			return
		}
	}
}

// detectPERF139: go/func closure capturing outer vars in handler
func detectPERF139(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b3hasHandlerSig(src) {
		return
	}
	file := unitFile(unit)
	// go func() { ... use outer ... } without params
	for _, gr := range facts.GoStarts {
		block := src[gr[0]:gr[1]]
		if !strings.Contains(block, "func()") && !strings.Contains(block, "func ()") {
			continue
		}
		// heuristic: uses w or secret-like outer without being a parameter
		if strings.Contains(block, "func()") || strings.Contains(block, "func ()") {
			// if func has no params and references identifiers from outer
			if strings.Contains(block, "w.") || strings.Contains(block, "w,") ||
				strings.Contains(block, "secret") || strings.Contains(block, "r.") {
				lineN, col := unit.LineCol(gr[0])
				rules.PushFinding(&MetaPERF139, file, lineN, col,
					"closure captures outer variables causing escapes", out)
				return
			}
		}
	}
}

// detectPERF140: debug.SetGCPercent(-1) or low values
func detectPERF140(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "debug.SetGCPercent" {
			continue
		}
		if len(call.Arguments) == 0 {
			continue
		}
		arg := strings.TrimSpace(call.Arguments[0])
		// flag -1 or small numbers
		if arg == "-1" || arg == "0" || arg == "1" || arg == "10" || arg == "25" {
			// skip if GOMEMLIMIT nearby
			if strings.Contains(unit.Source, "GOMEMLIMIT") || strings.Contains(unit.Source, "SetMemoryLimit") {
				continue
			}
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF140, file, lineN, col,
				"debug.SetGCPercent with aggressive value", out)
			return
		}
	}
}

// detectPERF141: r.URL.Query() called 2+ times
func detectPERF141(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if b3countSubstring(src, ".URL.Query()") < 2 {
		return
	}
	// if assigned once to var then .Get — safe (one Query call)
	// count call sites
	n := b3countSubstring(src, ".URL.Query()")
	// also q := r.URL.Query() counts as 1
	if n < 2 {
		return
	}
	file := unitFile(unit)
	idx := strings.Index(src, ".URL.Query()")
	lineN, col := unit.LineCol(idx)
	rules.PushFinding(&MetaPERF141, file, lineN, col,
		"URL.Query() called repeatedly; cache the result", out)
}

// detectPERF142: io.ReadAll(r.Body) without MaxBytesReader (Rust handler-gated).
func detectPERF142(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	// Rust: file_has_handler || http.ResponseWriter (ResponseWriter ⊂ handler).
	if !FileHasHandler(src) && !strings.Contains(src, "http.ResponseWriter") {
		return
	}
	if strings.Contains(src, "MaxBytesReader") {
		return
	}
	file := unitFile(unit)
	// Exact body-read pairs from Rust detect_perf_142.
	bodyReads := [][2]string{
		{"io.ReadAll(", "r.Body"},
		{"io.ReadAll(", "c.Request.Body"},
		{"io.ReadAll(", "req.Body"},
		{"io.ReadAll(", "ctx.Request.Body"},
		{"ioutil.ReadAll(", "r.Body"},
		{"ioutil.ReadAll(", "c.Request.Body"},
	}
	var pos int = -1
	for _, pair := range bodyReads {
		if strings.Contains(src, pair[0]) && strings.Contains(src, pair[1]) {
			pos = strings.Index(src, pair[0])
			break
		}
	}
	if pos < 0 {
		_ = facts
		return
	}
	lineN, col := unit.LineCol(pos)
	rules.PushFinding(&MetaPERF142, file, lineN, col,
		"request body is read without http.MaxBytesReader; cap the body size to prevent OOM", out)
}

func b3min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// detectPERF143: http.HandleFunc without TimeoutHandler
func detectPERF143(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if strings.Contains(src, "TimeoutHandler") {
		return
	}
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee == "http.HandleFunc" || call.Callee == "http.Handle" ||
			call.Callee == "mux.HandleFunc" || call.Callee == "mux.Handle" {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF143, file, lineN, col,
				"route registered without http.TimeoutHandler", out)
			return
		}
	}
}

// detectPERF144: w.Write(body) without Content-Length
func detectPERF144(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b3hasHandlerSig(src) {
		return
	}
	if strings.Contains(src, "Content-Length") || strings.Contains(src, "ContentLength") {
		return
	}
	file := unitFile(unit)
	// Write of a known []byte variable
	if strings.Contains(src, "w.Write(") && (strings.Contains(src, "[]byte(") || strings.Contains(src, ":= []byte") || strings.Contains(src, "body :=")) {
		idx := strings.Index(src, "w.Write(")
		if idx >= 0 {
			lineN, col := unit.LineCol(idx)
			rules.PushFinding(&MetaPERF144, file, lineN, col,
				"response Write without Content-Length for known body size", out)
		}
	}
}

// detectPERF145: r.WithContext in middleware
func detectPERF145(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	file := unitFile(unit)
	for _, call := range facts.Calls {
		// method call WithContext may appear as r.WithContext
		if !strings.HasSuffix(call.Callee, ".WithContext") && call.Callee != "WithContext" {
			// facts store full callee as r.WithContext
			if !strings.Contains(call.Callee, "WithContext") {
				continue
			}
		}
		// middleware-ish
		name, _ := EnclosingFunctionName(src, call.StartByte)
		lower := strings.ToLower(name)
		if strings.Contains(lower, "middleware") || strings.Contains(src, "http.Handler") ||
			IsHotPath(src, call.StartByte, IsInLoop(call)) {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF145, file, lineN, col,
				"r.WithContext allocates a new request on hot middleware path", out)
			return
		}
	}
	// text fallback
	if strings.Contains(src, ".WithContext(") && (strings.Contains(strings.ToLower(src), "middleware") || strings.Contains(src, "http.Handler")) {
		idx := strings.Index(src, ".WithContext(")
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF145, file, lineN, col,
			"r.WithContext allocates a new request on hot middleware path", out)
	}
}

// detectPERF146: fmt.Sprintf("%s", x)
func detectPERF146(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "fmt.Sprintf" || len(call.Arguments) < 2 {
			continue
		}
		fmtStr := strings.Trim(call.Arguments[0], `"`)
		if fmtStr == "%s" || fmtStr == "%v" {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF146, file, lineN, col,
				`fmt.Sprintf("%s", s) is unnecessary; use s directly`, out)
			return
		}
	}
}

// detectPERF148: unbuffered channel with send but no receive (exact Rust heuristic).
func detectPERF148(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "make(chan ") {
		return
	}
	if !strings.Contains(src, " <- ") {
		return
	}
	// Rust only treats these as receives (not `<-sem` / `:= <-` generically).
	hasReceive := strings.Contains(src, "<-ch") ||
		strings.Contains(src, "<- ch") ||
		strings.Contains(src, "for v := range ch") ||
		strings.Contains(src, "for range ch")
	if hasReceive {
		return
	}
	file := unitFile(unit)
	searchFrom := 0
	for {
		rel := strings.Index(src[searchFrom:], "make(chan ")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		stmtEndRel := strings.IndexAny(src[start:], "\n;")
		stmtEnd := len(src)
		if stmtEndRel >= 0 {
			stmtEnd = start + stmtEndRel
		}
		makeStmt := src[start:stmtEnd]
		isUnbuffered := true
		if strings.Contains(makeStmt, ", ") {
			// Capacity after last ", " — only digit runs count; variable caps
			// are treated as unbuffered (Rust parity).
			p := strings.LastIndex(makeStmt, ", ")
			after := makeStmt[p+2:]
			capStr := ""
			for _, r := range after {
				if r >= '0' && r <= '9' {
					capStr += string(r)
				} else {
					break
				}
			}
			if capStr != "0" && capStr != "" {
				isUnbuffered = false
			}
		}
		if isUnbuffered {
			lineN, col := unit.LineCol(start)
			rules.PushFinding(&MetaPERF148, file, lineN, col,
				"unbuffered channel created with no receive in this function; the sender may block forever if the receiver exits early", out)
			return
		}
		searchFrom = stmtEnd
		if searchFrom >= len(src) {
			return
		}
	}
}

// detectPERF149: conn.Read/Write without deadline
func detectPERF149(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if strings.Contains(src, "SetReadDeadline") || strings.Contains(src, "SetWriteDeadline") ||
		strings.Contains(src, "SetDeadline") {
		return
	}
	file := unitFile(unit)
	for _, call := range facts.Calls {
		// conn.Read / conn.Write
		if strings.HasSuffix(call.Callee, ".Read") || strings.HasSuffix(call.Callee, ".Write") {
			// likely net.Conn if file imports net or has net.Conn param
			if strings.Contains(src, "net.Conn") || strings.Contains(src, "\"net\"") {
				lineN, col := unit.LineCol(call.StartByte)
				rules.PushFinding(&MetaPERF149, file, lineN, col,
					"net.Conn I/O without deadline", out)
				return
			}
		}
	}
}

// detectPERF150: large stack arrays
func detectPERF150(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	file := unitFile(unit)
	if loc := reLargeArray.FindStringIndex(src); loc != nil {
		// extract size
		m := reLargeArray.FindString(src)
		// [4096]
		if reArrayDecl.MatchString(m) || reArrayDecl.MatchString(src[loc[0]:b3min(len(src), loc[0]+30)]) {
			lineN, col := unit.LineCol(loc[0])
			rules.PushFinding(&MetaPERF150, file, lineN, col,
				"large stack array may prevent inlining / bloat frames", out)
			return
		}
	}
	// var x [4096]byte
	re := regexp.MustCompile(`\[\d{4,}\]byte`)
	if loc := re.FindStringIndex(src); loc != nil {
		lineN, col := unit.LineCol(loc[0])
		rules.PushFinding(&MetaPERF150, file, lineN, col,
			"large stack array may prevent inlining / bloat frames", out)
	}
}

// detectPERF151: non-inlinable handler (Rust: loop+switch OR >50 func lines, + closure).
// Gate with FileHasHandler (not bare *http.Request) so CLI tools stay silent.
func detectPERF151(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !FileHasHandler(src) {
		return
	}
	hasLoop := strings.Contains(src, "for ")
	hasSwitch := strings.Contains(src, "switch ")
	hasClosure := strings.Contains(src, "func(") || strings.Contains(src, "go ")
	// Count lines from first func until a blank line (Rust heuristic).
	funcLines := 0
	seenFunc := false
	for _, line := range strings.Split(src, "\n") {
		if !seenFunc {
			if strings.Contains(line, "func ") {
				seenFunc = true
			} else {
				continue
			}
		}
		if strings.TrimSpace(line) == "" {
			break
		}
		funcLines++
	}
	complex := (hasLoop && hasSwitch) || funcLines > 50
	if !(complex && hasClosure) {
		return
	}
	file := unitFile(unit)
	idx := strings.Index(src, "func ")
	lineN, col := unit.LineCol(idx)
	rules.PushFinding(&MetaPERF151, file, lineN, col,
		"non-inlinable handler function: too complex for the Go compiler to inline; reduce body size or split into smaller functions", out)
}

// detectPERF152: for k,v := range src { dst.Set(k,v) }
func detectPERF152(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "range") || !strings.Contains(src, "http.Header") {
		return
	}
	file := unitFile(unit)
	for _, fr := range facts.ForRanges {
		block := src[fr[0]:fr[1]]
		if (strings.Contains(block, ".Set(") || strings.Contains(block, ".Add(")) &&
			strings.Contains(block, "range") {
			lineN, col := unit.LineCol(fr[0])
			rules.PushFinding(&MetaPERF152, file, lineN, col,
				"manual header copy loop; prefer Header.Clone", out)
			return
		}
	}
}

// detectPERF153: cookie.String() 2+ times
func detectPERF153(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if b3countSubstring(src, ".String()") < 2 {
		return
	}
	// look for cookie.String specifically
	n := 0
	for _, call := range facts.Calls {
		if strings.HasSuffix(call.Callee, ".String") || call.Callee == "String" {
			// context: cookie
			win := b3windowAround(unit.Source, call.StartByte, 40, 20)
			if strings.Contains(strings.ToLower(win), "cookie") || strings.Contains(call.Callee, "cookie") {
				n++
			}
		}
	}
	// text: cookie.String()
	if b3countSubstring(src, "cookie.String()") >= 2 || n >= 2 {
		idx := strings.Index(src, "cookie.String()")
		if idx < 0 {
			idx = strings.Index(src, ".String()")
		}
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF153, unitFile(unit), lineN, col,
			"cookie.String() called repeatedly; cache the result", out)
	}
}

// detectPERF154: http.Handle(..., http.HandlerFunc(h))
func detectPERF154(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "http.HandlerFunc(") {
		return
	}
	file := unitFile(unit)
	// Handle with HandlerFunc conversion of named func (not a literal)
	idx := strings.Index(src, "http.HandlerFunc(")
	for idx >= 0 {
		win := b3windowAround(src, idx, 40, 40)
		if strings.Contains(win, "http.Handle(") || strings.Contains(win, "Handle(") {
			// if argument is identifier not func literal
			rest := src[idx+len("http.HandlerFunc("):]
			arg := strings.TrimSpace(rest)
			if !strings.HasPrefix(arg, "func") {
				lineN, col := unit.LineCol(idx)
				rules.PushFinding(&MetaPERF154, file, lineN, col,
					"unnecessary http.HandlerFunc conversion; use HandleFunc", out)
				return
			}
		}
		next := strings.Index(src[idx+1:], "http.HandlerFunc(")
		if next < 0 {
			break
		}
		idx = idx + 1 + next
	}
}

// detectPERF155: if r.Method == in handler
func detectPERF155(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !b3hasHandlerSig(src) {
		return
	}
	if !strings.Contains(src, ".Method") {
		return
	}
	file := unitFile(unit)
	if strings.Contains(src, "r.Method ==") || strings.Contains(src, "r.Method==") ||
		strings.Contains(src, "switch r.Method") {
		idx := strings.Index(src, "r.Method")
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF155, file, lineN, col,
			"manual method check; prefer method-aware routing", out)
	}
}

// detectPERF156: for i, _ := range
func detectPERF156(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	file := unitFile(unit)
	re := regexp.MustCompile(`for\s+\w+\s*,\s*_\s*:=\s*range\s+`)
	if loc := re.FindStringIndex(src); loc != nil {
		lineN, col := unit.LineCol(loc[0])
		rules.PushFinding(&MetaPERF156, file, lineN, col,
			"range with blank value still decodes/copies; use for i := range", out)
		return
	}
}

// detectPERF157: fmt.Sprint("literal") or single string arg
func detectPERF157(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "fmt.Sprint" && call.Callee != "fmt.Sprintln" {
			continue
		}
		if len(call.Arguments) != 1 {
			continue
		}
		a := strings.TrimSpace(call.Arguments[0])
		if strings.HasPrefix(a, `"`) || strings.HasPrefix(a, "`") || isSimpleIdent(a) {
			// if ident, may not be string — still flag for fixtures with literal
			if strings.HasPrefix(a, `"`) || strings.HasPrefix(a, "`") {
				lineN, col := unit.LineCol(call.StartByte)
				rules.PushFinding(&MetaPERF157, file, lineN, col,
					"fmt.Sprint on a single string is unnecessary", out)
				return
			}
		}
	}
}

// detectPERF158: sort.Slice with simple < on basic slice
func detectPERF158(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "sort.Slice" {
			continue
		}
		// look at surrounding for simple comparator xs[i] < xs[j]
		win := b3windowAround(src, call.StartByte, 0, 120)
		if (strings.Contains(win, " < ") || strings.Contains(win, " > ")) &&
			strings.Contains(win, "func(") {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF158, file, lineN, col,
				"sort.Slice on basic types; prefer sort.Ints/Strings or slices.Sort", out)
			return
		}
	}
}

// detectPERF159: buffer then json.NewDecoder(bytes.NewReader(body))
func detectPERF159(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	file := unitFile(unit)
	if strings.Contains(src, "json.NewDecoder") &&
		(strings.Contains(src, "bytes.NewReader") || strings.Contains(src, "bytes.NewBuffer") ||
			strings.Contains(src, "json.Unmarshal")) {
		// ReadAll + decode path
		if strings.Contains(src, "ReadAll") || strings.Contains(src, "readAll") || strings.Contains(src, "ioutil.ReadAll") {
			idx := strings.Index(src, "json.NewDecoder")
			if idx < 0 {
				idx = strings.Index(src, "json.Unmarshal")
			}
			lineN, col := unit.LineCol(idx)
			rules.PushFinding(&MetaPERF159, file, lineN, col,
				"buffering body then decoding; stream with json.NewDecoder(r).Decode", out)
			return
		}
	}
	// classic ReadAll + Unmarshal
	if (strings.Contains(src, "io.ReadAll") || strings.Contains(src, "ioutil.ReadAll")) &&
		strings.Contains(src, "json.Unmarshal") {
		idx := strings.Index(src, "json.Unmarshal")
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF159, file, lineN, col,
			"ReadAll + json.Unmarshal; prefer streaming decoder", out)
	}
}

// detectPERF160: sql.Open in handler
// detectPERF160: sql.Open inside a request-handling file (not package-level var).
// Rust parity: file has a handler; skip `var ... = sql.Open` at package scope.
func detectPERF160(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !FileHasHandler(src) && !strings.Contains(src, "http.ResponseWriter") {
		return
	}
	if !strings.Contains(src, "sql.Open(") {
		return
	}
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "sql.Open" {
			continue
		}
		// Package-level `var db = sql.Open(...)` is the safe pattern.
		preStart := call.StartByte - 16
		if preStart < 0 {
			preStart = 0
		}
		pre := src[preStart:call.StartByte]
		if strings.Contains(pre, "var ") && !strings.Contains(pre, "func ") {
			continue
		}
		lineN, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF160, file, lineN, col,
			"sql.Open in a request handler; open the *sql.DB once at startup and reuse it across requests", out)
		return
	}
}

// detectPERF161: for rows.Next() without rows.Err()
func detectPERF161(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !strings.Contains(src, ".Next()") {
		return
	}
	if strings.Contains(src, ".Err()") {
		return
	}
	file := unitFile(unit)
	// for rows.Next()
	re := regexp.MustCompile(`for\s+\w+\.Next\s*\(\s*\)`)
	if loc := re.FindStringIndex(src); loc != nil {
		lineN, col := unit.LineCol(loc[0])
		rules.PushFinding(&MetaPERF161, file, lineN, col,
			"rows.Next loop without rows.Err() check", out)
		return
	}
	// also: for rows.Next()
	if strings.Contains(src, "rows.Next()") && strings.Contains(src, "for ") {
		idx := strings.Index(src, "rows.Next()")
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF161, file, lineN, col,
			"rows.Next loop without rows.Err() check", out)
	}
}

// detectPERF162: db.Ping inside a request handler (not a health-check route).
// Rust parity: file has a handler; skip calls under func Health / Healthz.
func detectPERF162(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !FileHasHandler(src) && !strings.Contains(src, "http.ResponseWriter") {
		return
	}
	if !strings.Contains(src, ".Ping(") && !strings.Contains(src, ".PingContext(") {
		return
	}
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "db.Ping" && call.Callee != "db.PingContext" &&
			!strings.HasSuffix(call.Callee, ".Ping") && !strings.HasSuffix(call.Callee, ".PingContext") {
			continue
		}
		// Health-check endpoints are the intended home for db.Ping.
		preStart := call.StartByte - 256
		if preStart < 0 {
			preStart = 0
		}
		pre := src[preStart:call.StartByte]
		if strings.Contains(pre, "func Health") ||
			strings.Contains(pre, "func (h *Health") ||
			strings.Contains(pre, "func Healthz") {
			continue
		}
		lineN, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF162, file, lineN, col,
			"db.Ping in a request handler; add a dedicated health-check endpoint or a periodic background ping instead", out)
		return
	}
}

// detectPERF163: db.Query + if rows.Next() (not for) single row
func detectPERF163(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !strings.Contains(src, ".Query(") && !strings.Contains(src, ".QueryContext(") {
		return
	}
	if strings.Contains(src, "QueryRow") {
		return
	}
	// if rows.Next() without for
	if strings.Contains(src, "if rows.Next()") || strings.Contains(src, "if rows.Next (") {
		file := unitFile(unit)
		idx := strings.Index(src, ".Query(")
		if idx < 0 {
			idx = strings.Index(src, ".QueryContext(")
		}
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF163, file, lineN, col,
			"db.Query used for single-row result; prefer QueryRow", out)
		return
	}
}

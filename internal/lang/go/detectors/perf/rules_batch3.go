package perf

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
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

// detectPERF112: strings.ToLower/ToUpper both sides of == / !=
func detectPERF112(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "strings.ToLower") && !strings.Contains(src, "strings.ToUpper") {
		return
	}
	file := unitFile(unit)
	// scan for ToLower/ToUpper == ToLower/ToUpper patterns on a line
	for i := 0; i < len(src); {
		idx := strings.Index(src[i:], "strings.To")
		if idx < 0 {
			break
		}
		pos := i + idx
		line := b3lineWindow(src, pos)
		hasFold := (strings.Contains(line, "strings.ToLower") || strings.Contains(line, "strings.ToUpper")) &&
			(strings.Contains(line, "==") || strings.Contains(line, "!="))
		// both sides
		nLower := strings.Count(line, "strings.ToLower") + strings.Count(line, "strings.ToUpper")
		if hasFold && nLower >= 2 {
			lineN, col := unit.LineCol(pos)
			rules.PushFinding(&MetaPERF112, file, lineN, col,
				"case-insensitive compare via ToLower/ToUpper; prefer strings.EqualFold", out)
			return
		}
		i = pos + 1
	}
}

// detectPERF113: select with a single case and no default
func detectPERF113(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "select") {
		return
	}
	file := unitFile(unit)
	// crude scan of select blocks
	for i := 0; i < len(src); {
		idx := strings.Index(src[i:], "select")
		if idx < 0 {
			break
		}
		pos := i + idx
		// word boundary
		if pos > 0 && (b3isIdentByte(src[pos-1]) || src[pos-1] == '.') {
			i = pos + 6
			continue
		}
		rest := src[pos+6:]
		// skip whitespace
		rest = strings.TrimLeft(rest, " \t\n\r")
		if !strings.HasPrefix(rest, "{") {
			i = pos + 6
			continue
		}
		block, ok := b3extractBraceBlock(src, pos+6+strings.Index(src[pos+6:], "{"))
		if !ok {
			i = pos + 6
			continue
		}
		// count case / default at top level of select body
		body := block[1 : len(block)-1]
		cases := 0
		hasDefault := false
		for _, ln := range strings.Split(body, "\n") {
			t := strings.TrimSpace(ln)
			if strings.HasPrefix(t, "case ") || strings.HasPrefix(t, "case\t") {
				cases++
			}
			if t == "default:" || strings.HasPrefix(t, "default:") {
				hasDefault = true
			}
		}
		if cases == 1 && !hasDefault {
			lineN, col := unit.LineCol(pos)
			rules.PushFinding(&MetaPERF113, file, lineN, col,
				"single-case select adds overhead; use a direct channel op", out)
			return
		}
		i = pos + 6
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
func detectPERF114(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "range") || !strings.Contains(src, "for ") {
		return
	}
	file := unitFile(unit)
	// line-based: for i, v := range X { Y[i] = v } (no conversion call on RHS)
	lines := strings.Split(src, "\n")
	for li, ln := range lines {
		t := strings.TrimSpace(ln)
		if !strings.HasPrefix(t, "for ") || !strings.Contains(t, "range") {
			continue
		}
		// for i, v := range src  OR  for i := range src with body dst[i]=src[i]
		hasVal := strings.Contains(t, ",")
		for j := li + 1; j < len(lines) && j < li+4; j++ {
			body := strings.TrimSpace(lines[j])
			if body == "}" {
				break
			}
			if !(strings.Contains(body, "[") && (strings.Contains(body, "] =") || strings.Contains(body, "]="))) {
				continue
			}
			eq := strings.Index(body, "] =")
			if eq < 0 {
				eq = strings.Index(body, "]=")
			}
			if eq < 0 {
				continue
			}
			rhs := body[eq:]
			// reject conversion / function call on RHS (except pure index)
			if strings.Contains(rhs, "(") {
				continue
			}
			// require simple assignment form
			if !hasVal && !strings.Contains(rhs, "[") {
				continue
			}
			off := 0
			for k := 0; k < li; k++ {
				off += len(lines[k]) + 1
			}
			lineN, col := unit.LineCol(off)
			rules.PushFinding(&MetaPERF114, file, lineN, col,
				"manual for-range copy; use the copy() builtin", out)
			return
		}
	}
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

// detectPERF119: 2+ consecutive appends to same slice
func detectPERF119(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	b3detectConsecutiveAppends(unit, facts, out, 2, &MetaPERF119,
		"multiple sequential appends; merge into one variadic append")
}

// detectPERF128: 3+ consecutive appends
func detectPERF128(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	b3detectConsecutiveAppends(unit, facts, out, 3, &MetaPERF128,
		"3+ sequential appends; merge into one variadic append")
}

func b3detectConsecutiveAppends(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding, min int, meta *rules.RuleMetadata, msg string) {
	_ = facts
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
		// match: name = append(name, ...)
		if strings.Contains(t, "= append(") {
			lhs := strings.TrimSpace(strings.Split(t, "=")[0])
			lhs = strings.TrimSuffix(lhs, ":")
			// single-value append heuristic: not multi-arg with commas after first
			// still count multi-arg as one append statement
			if runVar == lhs {
				runCount++
			} else {
				runVar = lhs
				runCount = 1
				runStartLine = li
			}
			if runCount >= min {
				// byte offset of run start
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
		// intervening non-append resets
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

// detectPERF121: struct literal field-by-field copy Options{ Host: src.Host, Port: src.Port }
func detectPERF121(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	// look for Type{ Field: other.Field, Field2: other.Field2 }
	// at least 2 fields from same source
	file := unitFile(unit)
	// simple: consecutive Field: x.Field patterns inside composite literal
	re := regexp.MustCompile(`(\w+)\s*\{\s*(\w+)\s*:\s*(\w+)\.(\w+)\s*,\s*(\w+)\s*:\s*(\w+)\.(\w+)`)
	if m := re.FindStringSubmatchIndex(src); m != nil {
		// same source ident and field names match keys
		sub := re.FindStringSubmatch(src)
		if sub != nil && sub[2] == sub[4] && sub[5] == sub[7] && sub[3] == sub[6] {
			lineN, col := unit.LineCol(m[0])
			rules.PushFinding(&MetaPERF121, file, lineN, col,
				"manual struct field copy; consider direct type conversion when layouts match", out)
			return
		}
	}
}

// detectPERF122: HasPrefix + slice with len
func detectPERF122(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "strings.HasPrefix" && call.Callee != "bytes.HasPrefix" {
			continue
		}
		win := b3windowAround(unit.Source, call.StartByte, 0, 320)
		if strings.Contains(win, "len(") && strings.Contains(win, "]:") {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF122, file, lineN, col,
				"HasPrefix + slice; prefer TrimPrefix", out)
			return
		}
	}
}

// detectPERF123: make([]T, 0) or make(map[K]V, 0)
func detectPERF123(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "make(") {
		return
	}
	file := unitFile(unit)
	// make([]T, 0) or make([]T, 0, n) or make(map[K]V, 0)
	if loc := reMakeZero.FindStringIndex(src); loc != nil {
		// make([]T, 0, cap) is actually useful for capacity — detection notes say flag make(T,0) and make(T,0,N)
		// But make([]T, 0, n) is a common capacity hint and often intentional.
		// Fixture is make([]int, 0) — flag that.
		snippet := src[loc[0]:loc[1]]
		// flag make(T, 0) and make(T, 0) only when no capacity, OR make(map, 0)
		if strings.Contains(snippet, "map[") || strings.HasSuffix(strings.TrimSpace(snippet), ", 0)") || strings.Contains(snippet, ", 0)") {
			// if make([]T, 0, N) — still per notes, but skip capacity form to reduce FP noise except pure , 0)
			if strings.Count(snippet, ",") == 1 {
				lineN, col := unit.LineCol(loc[0])
				rules.PushFinding(&MetaPERF123, file, lineN, col,
					"redundant zero argument to make", out)
				return
			}
		}
	}
	// also catch make([]int, 0) with flexible spacing
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
		// one comma and ends with , 0)
		if strings.Count(call, ",") == 1 && (strings.HasSuffix(call, ", 0)") || strings.HasSuffix(call, ",0)")) {
			if strings.Contains(call, "[]") || strings.Contains(call, "map[") {
				lineN, col := unit.LineCol(idx)
				rules.PushFinding(&MetaPERF123, file, lineN, col,
					"redundant zero argument to make", out)
				return
			}
		}
		next := strings.Index(src[idx+1:], "make(")
		if next < 0 {
			break
		}
		idx = idx + 1 + next
	}
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

// detectPERF130: immediately-invoked func literal func(){ ... }()
func detectPERF130(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "func()") && !strings.Contains(src, "func ()") {
		return
	}
	file := unitFile(unit)
	// look for function literals that are immediately invoked and not assigned
	// Use FunctionLiteralRanges + following ()
	for _, fl := range facts.FunctionLiteralRanges {
		end := fl[1]
		if end >= len(src) {
			continue
		}
		// after literal, optional whitespace then (
		rest := strings.TrimLeft(src[end:], " \t\n\r")
		if !strings.HasPrefix(rest, "(") {
			continue
		}
		// skip if this is go func or defer func
		before := b3windowAround(src, fl[0], 16, 0)
		if strings.Contains(before, "go ") || strings.Contains(before, "defer ") {
			continue
		}
		// skip if assigned: = func
		if strings.Contains(before, "=") {
			continue
		}
		// body should be simple (single statement-ish)
		block := src[fl[0]:fl[1]]
		// count statements roughly by semicolons/newlines inside
		lineN, col := unit.LineCol(fl[0])
		rules.PushFinding(&MetaPERF130, file, lineN, col,
			"unnecessary immediately-invoked function wrapper", out)
		_ = block
		return
	}
	// text fallback
	if loc := reIife.FindStringIndex(src); loc != nil {
		before := b3windowAround(src, loc[0], 16, 0)
		if strings.Contains(before, "go ") || strings.Contains(before, "defer ") || strings.Contains(before, "=") {
			return
		}
		lineN, col := unit.LineCol(loc[0])
		rules.PushFinding(&MetaPERF130, file, lineN, col,
			"unnecessary immediately-invoked function wrapper", out)
	}
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
func detectPERF132(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !strings.Contains(src, "go ") {
		return
	}
	file := unitFile(unit)
	for _, gr := range facts.GoStarts {
		block := src[gr[0]:gr[1]]
		// has I/O
		hasIO := strings.Contains(block, "db.Query") || strings.Contains(block, "db.Exec") ||
			strings.Contains(block, "http.") || strings.Contains(block, ".Do(") ||
			strings.Contains(block, "Query(") || strings.Contains(block, "Exec(")
		if !hasIO {
			continue
		}
		// context used?
		if strings.Contains(block, "ctx") || strings.Contains(block, "Context") {
			continue
		}
		// go func() without params taking context
		lineN, col := unit.LineCol(gr[0])
		rules.PushFinding(&MetaPERF132, file, lineN, col,
			"goroutine performs I/O without context propagation", out)
		return
	}
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

// detectPERF137: runtime.Caller on hot path
func detectPERF137(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	src := unit.Source
	for _, call := range facts.Calls {
		if call.Callee != "runtime.Caller" {
			continue
		}
		if IsInLoop(call) || IsHotPath(src, call.StartByte, IsInLoop(call)) || b3hasHandlerSig(src) {
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

// detectPERF142: io.ReadAll(r.Body) / json decoder without MaxBytesReader
func detectPERF142(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !strings.Contains(src, "r.Body") && !strings.Contains(src, "req.Body") {
		return
	}
	if strings.Contains(src, "MaxBytesReader") {
		return
	}
	file := unitFile(unit)
	for _, call := range facts.Calls {
		switch call.Callee {
		case "io.ReadAll", "ioutil.ReadAll", "json.NewDecoder":
		default:
			// also method ReadAll on body
			continue
		}
		// check args involve Body
		joined := strings.Join(call.Arguments, ",")
		if strings.Contains(joined, ".Body") || strings.Contains(unit.Source[call.StartByte:b3min(len(unit.Source), call.StartByte+40)], "Body") {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF142, file, lineN, col,
				"reading request body without http.MaxBytesReader", out)
			return
		}
	}
	// text: io.ReadAll(r.Body)
	if strings.Contains(src, "ReadAll(r.Body)") || strings.Contains(src, "ReadAll(req.Body)") ||
		strings.Contains(src, "NewDecoder(r.Body)") {
		idx := strings.Index(src, "ReadAll(r.Body)")
		if idx < 0 {
			idx = strings.Index(src, "NewDecoder(r.Body)")
		}
		if idx < 0 {
			idx = strings.Index(src, "r.Body")
		}
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF142, file, lineN, col,
			"reading request body without http.MaxBytesReader", out)
	}
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

// detectPERF148: unbuffered channel send without guaranteed receiver
func detectPERF148(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "make(chan") {
		return
	}
	// unbuffered: make(chan T) without capacity arg
	if !regexp.MustCompile(`make\s*\(\s*chan\s+[^,)]+\s*\)`).MatchString(src) {
		return
	}
	if !strings.Contains(src, "<-") {
		return
	}
	file := unitFile(unit)
	sendOnly := false
	for _, ln := range strings.Split(src, "\n") {
		t := strings.TrimSpace(ln)
		if !strings.Contains(t, "<-") {
			continue
		}
		// receive forms: <-ch, x := <-ch, case <-ch
		if strings.HasPrefix(t, "<-") || strings.Contains(t, ":= <-") ||
			strings.Contains(t, ", <-") || strings.Contains(t, "case <-") {
			return
		}
		// send form: ch <- value
		if regexp.MustCompile(`\w+\s*<-\s*`).MatchString(t) {
			sendOnly = true
		}
	}
	if !sendOnly {
		return
	}
	idx := strings.Index(src, "<-")
	lineN, col := unit.LineCol(idx)
	rules.PushFinding(&MetaPERF148, file, lineN, col,
		"send on unbuffered channel without guaranteed receiver", out)
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

// detectPERF151: complex handler (switch+for+go)
func detectPERF151(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !b3hasHandlerSig(src) {
		return
	}
	file := unitFile(unit)
	// count complexity tokens in functions with Handler in name or handler sig
	// simple: function with switch + for + go
	if strings.Contains(src, "switch ") && strings.Contains(src, "for ") && strings.Contains(src, "go ") {
		idx := strings.Index(src, "func ")
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF151, file, lineN, col,
			"hot-path handler likely exceeds inlining budget", out)
		return
	}
	// long function: >40 lines
	if strings.Count(src, "\n") > 45 && b3hasHandlerSig(src) {
		idx := strings.Index(src, "func ")
		lineN, col := unit.LineCol(idx)
		rules.PushFinding(&MetaPERF151, file, lineN, col,
			"hot-path function is complex/long; may not inline", out)
	}
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
func detectPERF160(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "sql.Open" {
			continue
		}
		if IsInLoop(call) || b3hasHandlerSig(src) || IsHotPath(src, call.StartByte, IsInLoop(call)) {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF160, file, lineN, col,
				"sql.Open inside handler/loop; reuse a shared *sql.DB", out)
			return
		}
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

// detectPERF162: db.Ping in handler
func detectPERF162(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "db.Ping" && call.Callee != "db.PingContext" &&
			!strings.HasSuffix(call.Callee, ".Ping") && !strings.HasSuffix(call.Callee, ".PingContext") {
			continue
		}
		if b3hasHandlerSig(src) || IsHotPath(src, call.StartByte, IsInLoop(call)) {
			lineN, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF162, file, lineN, col,
				"db.Ping inside request handler", out)
			return
		}
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

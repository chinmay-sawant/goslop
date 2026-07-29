package perf

import (
	"strings"
	"unicode"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("PERF-164", detectPERF164, &MetaPERF164)
	RegisterRule("PERF-165", detectPERF165, &MetaPERF165)
	RegisterRule("PERF-166", detectPERF166, &MetaPERF166)
	RegisterRule("PERF-167", detectPERF167, &MetaPERF167)
	RegisterRule("PERF-168", detectPERF168, &MetaPERF168)
	RegisterRule("PERF-169", detectPERF169, &MetaPERF169)
	RegisterRule("PERF-170", detectPERF170, &MetaPERF170)
	RegisterRule("PERF-171", detectPERF171, &MetaPERF171)
	RegisterRule("PERF-172", detectPERF172, &MetaPERF172)
	RegisterRule("PERF-173", detectPERF173, &MetaPERF173)
	RegisterRule("PERF-174", detectPERF174, &MetaPERF174)
	RegisterRule("PERF-175", detectPERF175, &MetaPERF175)
	RegisterRule("PERF-176", detectPERF176, &MetaPERF176)
	RegisterRule("PERF-177", detectPERF177, &MetaPERF177)
	RegisterRule("PERF-178", detectPERF178, &MetaPERF178)
	RegisterRule("PERF-179", detectPERF179, &MetaPERF179)
	RegisterRule("PERF-180", detectPERF180, &MetaPERF180)
	RegisterRule("PERF-181", detectPERF181, &MetaPERF181)
	RegisterRule("PERF-182", detectPERF182, &MetaPERF182)
	RegisterRule("PERF-183", detectPERF183, &MetaPERF183)
	RegisterRule("PERF-184", detectPERF184, &MetaPERF184)
	RegisterRule("PERF-185", detectPERF185, &MetaPERF185)
	RegisterRule("PERF-186", detectPERF186, &MetaPERF186)
	RegisterRule("PERF-187", detectPERF187, &MetaPERF187)
	RegisterRule("PERF-188", detectPERF188, &MetaPERF188)
	RegisterRule("PERF-189", detectPERF189, &MetaPERF189)
	RegisterRule("PERF-190", detectPERF190, &MetaPERF190)
	RegisterRule("PERF-191", detectPERF191, &MetaPERF191)
	RegisterRule("PERF-192", detectPERF192, &MetaPERF192)
	RegisterRule("PERF-193", detectPERF193, &MetaPERF193)
	RegisterRule("PERF-194", detectPERF194, &MetaPERF194)
	RegisterRule("PERF-195", detectPERF195, &MetaPERF195)
	RegisterRule("PERF-196", detectPERF196, &MetaPERF196)
	RegisterRule("PERF-197", detectPERF197, &MetaPERF197)
	RegisterRule("PERF-198", detectPERF198, &MetaPERF198)
	RegisterRule("PERF-199", detectPERF199, &MetaPERF199)
	RegisterRule("PERF-200", detectPERF200, &MetaPERF200)
	RegisterRule("PERF-201", detectPERF201, &MetaPERF201)
	RegisterRule("PERF-202", detectPERF202, &MetaPERF202)
	RegisterRule("PERF-203", detectPERF203, &MetaPERF203)
	RegisterRule("PERF-204", detectPERF204, &MetaPERF204)
	RegisterRule("PERF-205", detectPERF205, &MetaPERF205)
	RegisterRule("PERF-206", detectPERF206, &MetaPERF206)
	RegisterRule("PERF-207", detectPERF207, &MetaPERF207)
	RegisterRule("PERF-209", detectPERF209, &MetaPERF209)
	RegisterRule("PERF-210", detectPERF210, &MetaPERF210)
	RegisterRule("PERF-211", detectPERF211, &MetaPERF211)
	RegisterRule("PERF-212", detectPERF212, &MetaPERF212)
	RegisterRule("PERF-213", detectPERF213, &MetaPERF213)
	RegisterRule("PERF-214", detectPERF214, &MetaPERF214)
}

// --- batch4-local helpers (b4_ prefix avoids collisions with other batches) ---

func b4_fileHasHandler(source string) bool {
	return IsHandlerShaped(source, len(source))
}

func b4_methodName(callee string) string {
	if i := strings.LastIndex(callee, "."); i >= 0 {
		return callee[i+1:]
	}
	return callee
}

func b4_isLargeStructLiteral(literal string) bool {
	inner := strings.TrimSpace(literal)
	inner = strings.TrimPrefix(inner, "{")
	inner = strings.TrimSuffix(inner, "}")
	depth := 0
	fields := 0
	hasComplex := false
	var current strings.Builder
	for _, c := range inner {
		switch c {
		case '(', '[', '{':
			depth++
			current.WriteRune(c)
		case ')', ']', '}':
			depth--
			current.WriteRune(c)
		case ',':
			if depth == 0 {
				fields++
				trimmed := strings.TrimSpace(current.String())
				if strings.Contains(trimmed, "[") || strings.Contains(trimmed, "map[") {
					hasComplex = true
				}
				current.Reset()
			} else {
				current.WriteRune(c)
			}
		default:
			current.WriteRune(c)
		}
	}
	if strings.TrimSpace(current.String()) != "" {
		fields++
		trimmed := strings.TrimSpace(current.String())
		if strings.Contains(trimmed, "[") || strings.Contains(trimmed, "map[") {
			hasComplex = true
		}
	}
	return fields >= 4 || hasComplex
}

func b4_hasLargeStringLiteral(window string) bool {
	bytes := []byte(window)
	inString := false
	start := 0
	total := 0
	for i, b := range bytes {
		if b == '"' && (i == 0 || bytes[i-1] != '\\') {
			if inString {
				length := i - start - 1
				if length > 64 {
					return true
				}
				total += length
				if total > 64 {
					return true
				}
			} else {
				start = i
			}
			inString = !inString
		}
	}
	return false
}

func b4_enclosingFunctionBody(source string, startByte int) string {
	if startByte > len(source) {
		startByte = len(source)
	}
	if startByte < 0 {
		startByte = 0
	}
	funcStart := strings.LastIndex(source[:startByte], "func ")
	if funcStart < 0 {
		return source
	}
	rel := strings.Index(source[funcStart:], "{")
	if rel < 0 {
		return source
	}
	bodyStart := funcStart + rel + 1
	depth := 1
	for i := bodyStart; i < len(source); i++ {
		switch source[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return source[bodyStart:i]
			}
		}
	}
	return source[bodyStart:]
}

func b4_mapSizeHintAvailable(source string, startByte int) bool {
	body := b4_enclosingFunctionBody(source, startByte)
	return strings.Contains(body, " range ") ||
		strings.Contains(body, "\trange ") ||
		strings.Contains(body, "range ") ||
		strings.Contains(body, "len(")
}

type b4_packageLevelCache struct {
	name      string
	byteOff   int
	isSyncMap bool
}

func b4_packageLevelCaches(source string) []b4_packageLevelCache {
	var out []b4_packageLevelCache
	depth := 0
	inVarBlock := false
	byteOff := 0
	for _, line := range strings.SplitAfter(source, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		indent := len(line) - len(trimmed)
		trimmedByte := byteOff + indent
		if depth == 0 {
			if inVarBlock {
				if strings.HasPrefix(trimmed, ")") {
					inVarBlock = false
				} else if c, ok := b4_parseCacheDeclLine(trimmed, trimmedByte, false); ok {
					out = append(out, c)
				}
			} else if strings.HasPrefix(trimmed, "var (") {
				inVarBlock = true
			} else if c, ok := b4_parseCacheDeclLine(trimmed, trimmedByte, true); ok {
				out = append(out, c)
			}
		}
		depth += strings.Count(line, "{") - strings.Count(line, "}")
		byteOff += len(line)
	}
	return out
}

func b4_parseCacheDeclLine(line string, byteOff int, expectVarPrefix bool) (b4_packageLevelCache, bool) {
	rest := line
	if expectVarPrefix {
		if !strings.HasPrefix(line, "var ") {
			return b4_packageLevelCache{}, false
		}
		rest = strings.TrimPrefix(line, "var ")
	}
	parts := strings.Fields(rest)
	if len(parts) == 0 {
		return b4_packageLevelCache{}, false
	}
	name := strings.TrimSuffix(parts[0], ",")
	if !isSimpleIdent(name) {
		return b4_packageLevelCache{}, false
	}
	isSyncMap := strings.Contains(rest, "sync.Map")
	isPlainMap := strings.Contains(rest, "map[")
	if !isSyncMap && !isPlainMap {
		return b4_packageLevelCache{}, false
	}
	return b4_packageLevelCache{name: name, byteOff: byteOff, isSyncMap: isSyncMap}, true
}

func b4_cacheHasEvictionBound(source, name string) bool {
	lower := strings.ToLower(source)
	nameLower := strings.ToLower(name)
	patterns := []string{
		"len(" + name + ") >",
		"len(" + name + ") >=",
		"cap(" + name + ") >",
		"cap(" + name + ") >=",
		"clear(" + name + ")",
		"delete(" + name + ",",
		name + ".delete(",
		name + ".loadanddelete(",
	}
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	timeMarkers := []string{
		"ttl", "expire", "expires", "expiry", "evict", "eviction",
		"time.now(", "time.since(", ".before(", ".after(", "ticker", "timer",
	}
	for _, line := range strings.Split(source, "\n") {
		ll := strings.ToLower(line)
		if !strings.Contains(ll, nameLower) {
			continue
		}
		for _, m := range timeMarkers {
			if strings.Contains(ll, m) {
				return true
			}
		}
	}
	return false
}

func b4_cacheIsWritten(source string, facts *GoPerfFacts, name string, isSyncMap bool) bool {
	if isSyncMap {
		prefix := name + "."
		for _, call := range facts.Calls {
			if !strings.HasPrefix(call.Callee, prefix) {
				continue
			}
			switch b4_methodName(call.Callee) {
			case "Store", "Swap", "LoadOrStore", "Delete", "LoadAndDelete":
				return true
			}
		}
		return false
	}
	for _, a := range facts.Assignments {
		if strings.HasPrefix(strings.TrimLeft(a.Text, " \t"), name+"[") {
			return true
		}
	}
	return strings.Contains(source, name+"[")
}

func b4_volatileCacheKey(key string) bool {
	lower := strings.ToLower(key)
	return strings.Contains(key, "&") ||
		strings.Contains(lower, "requestid") ||
		strings.Contains(lower, "reqid") ||
		strings.Contains(lower, "trace") ||
		strings.Contains(lower, "session") ||
		strings.Contains(lower, " idx") ||
		strings.Contains(lower, "index") ||
		strings.Contains(lower, "timestamp") ||
		strings.Contains(lower, "time.now") ||
		strings.Contains(key, "%p")
}

func b4_cacheKeyIsVolatile(unit *core.ParsedUnit, facts *GoPerfFacts, callStart int, key string) bool {
	if b4_volatileCacheKey(key) {
		return true
	}
	if !isSimpleIdent(key) {
		return false
	}
	body := b4_enclosingFunctionBody(unit.Source, callStart)
	funcStart := strings.LastIndex(unit.Source[:callStart], "func ")
	if funcStart < 0 {
		return false
	}
	rel := strings.Index(unit.Source[funcStart:], "{")
	if rel < 0 {
		return false
	}
	scopeStart := funcStart + rel
	scopeEnd := scopeStart + len(body) + 1
	for i := len(facts.Assignments) - 1; i >= 0; i-- {
		a := facts.Assignments[i]
		if a.Name != key {
			continue
		}
		if a.StartByte >= scopeStart && a.StartByte < callStart && a.StartByte < scopeEnd {
			if b4_volatileCacheKey(a.Expr) {
				return true
			}
		}
	}
	return false
}

func b4_indexAny(source string, needles ...string) int {
	best := -1
	for _, n := range needles {
		if i := strings.Index(source, n); i >= 0 {
			if best < 0 || i < best {
				best = i
			}
		}
	}
	return best
}

// --- detectors ---

func detectPERF164(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !b4_fileHasHandler(source) && !strings.Contains(source, "http.ResponseWriter") {
		return
	}
	hasRequestVar := strings.Contains(source, "r *http.Request") ||
		strings.Contains(source, "r.Context()") ||
		strings.Contains(source, "req.Context()") ||
		strings.Contains(source, "r.Header") ||
		strings.Contains(source, "r.URL") ||
		strings.Contains(source, "r.Method") ||
		strings.Contains(source, "r.Body") ||
		strings.Contains(source, "c.Request.Context()") ||
		strings.Contains(source, "ctx.Request.Context()") ||
		strings.Contains(source, "c.Request.Body")
	if !hasRequestVar {
		return
	}
	if !(strings.Contains(source, "db.Query(") || strings.Contains(source, "db.Exec(") ||
		strings.Contains(source, "db.Prepare(") || strings.Contains(source, "db.Begin(")) {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "db.Query", "db.Exec", "db.Prepare", "db.Begin":
		default:
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF164, file, line, col,
			"db.* call without Context in a request handler; use the Context variant for cancellation propagation", out)
	}
}

func detectPERF165(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "rows.Scan(") {
		return
	}
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], "rows.Scan(")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		afterStart := start + len("rows.Scan(")
		afterEnd := afterStart + 384
		if afterEnd > len(source) {
			afterEnd = len(source)
		}
		window := source[afterStart:afterEnd]
		if close := strings.Index(window, ")"); close >= 0 {
			firstArg := window[:close]
			if strings.Contains(firstArg, "&sql.Null") ||
				strings.Contains(firstArg, "&*string") ||
				strings.Contains(firstArg, "&*int") {
				searchFrom = afterStart
				continue
			}
		}
		funcEnd := len(source)
		if i := strings.Index(source[afterStart:], "}\n\n"); i >= 0 {
			funcEnd = afterStart + i
		} else if i := strings.Index(source[afterStart:], "}\nfunc"); i >= 0 {
			funcEnd = afterStart + i
		} else if i := strings.LastIndex(source[afterStart:], "}"); i >= 0 {
			funcEnd = afterStart + i
		}
		afterBlock := source[afterStart:funcEnd]
		hasParse := strings.Contains(afterBlock, "strconv.Parse") ||
			strings.Contains(afterBlock, "time.Parse") ||
			strings.Contains(afterBlock, "uuid.Parse") ||
			strings.Contains(afterBlock, ".Parse(")
		hasCustom := strings.Contains(afterBlock, "MyID(") ||
			strings.Contains(afterBlock, "MyType(") ||
			strings.Contains(afterBlock, "UUID(") ||
			strings.Contains(afterBlock, "Timestamp(") ||
			strings.Contains(afterBlock, "MyTime(")
		if hasParse && hasCustom {
			line, col := unit.LineCol(start)
			rules.PushFinding(&MetaPERF165, file, line, col,
				"rows.Scan into a primitive type followed by manual conversion to a custom type; implement sql.Scanner on the custom type", out)
		}
		searchFrom = afterStart
	}
}

func detectPERF166(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "rows.Scan(") {
		return
	}
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], "rows.Scan(")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		afterStart := start + len("rows.Scan(")
		afterEnd := afterStart + 384
		if afterEnd > len(source) {
			afterEnd = len(source)
		}
		window := source[afterStart:afterEnd]
		close := strings.Index(window, ")")
		if close < 0 {
			searchFrom = afterStart
			continue
		}
		firstArg := strings.TrimSpace(window[:close])
		if !strings.HasPrefix(firstArg, "&") {
			searchFrom = afterStart
			continue
		}
		varName := strings.TrimSpace(strings.TrimPrefix(firstArg, "&"))
		if !isSimpleIdent(varName) || strings.HasPrefix(varName, "Null") {
			searchFrom = afterStart
			continue
		}
		blockEnd := afterStart + 512
		if blockEnd > len(source) {
			blockEnd = len(source)
		}
		if !strings.Contains(source[afterStart:blockEnd], "if "+varName+" != nil") {
			searchFrom = afterStart
			continue
		}
		line, col := unit.LineCol(start)
		rules.PushFinding(&MetaPERF166, file, line, col,
			"rows.Scan into a pointer followed by a null check; use sql.NullString / sql.NullInt64 instead", out)
		searchFrom = afterStart
	}
}

func detectPERF167(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		if !strings.HasSuffix(call.Callee, ".Add") {
			continue
		}
		windowStart := call.StartByte - 2048
		if windowStart < 0 {
			windowStart = 0
		}
		window := source[windowStart:call.StartByte]
		goIdx := strings.LastIndex(window, "go func()")
		if goIdx < 0 {
			continue
		}
		after := window[goIdx:]
		depth := strings.Count(after, "{") - strings.Count(after, "}")
		if depth <= 0 {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF167, file, line, col,
			"wg.Add inside a goroutine body; the Add must happen before the goroutine is started", out)
	}
}

func detectPERF168(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], "<- ")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		if start == 0 || source[start-1] != ' ' {
			searchFrom = start + len("<- ")
			continue
		}
		arrowEnd := start + len("<- ")
		openRel := strings.Index(source[arrowEnd:], "{")
		if openRel < 0 {
			searchFrom = arrowEnd
			continue
		}
		pre := source[arrowEnd : arrowEnd+openRel]
		if strings.HasPrefix(strings.TrimLeft(pre, " \t"), "&") {
			searchFrom = arrowEnd + openRel
			continue
		}
		litStart := arrowEnd + openRel
		closeRel := strings.Index(source[litStart:], "}")
		if closeRel < 0 {
			searchFrom = litStart
			continue
		}
		literal := source[litStart : litStart+closeRel+1]
		if b4_isLargeStructLiteral(literal) {
			line, col := unit.LineCol(start)
			rules.PushFinding(&MetaPERF168, file, line, col,
				"large struct sent by value over a channel; pass a pointer instead", out)
		}
		searchFrom = litStart + closeRel + 1
	}
}

func detectPERF169(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "sync/atomic") {
		return
	}
	for _, call := range facts.Calls {
		if !strings.HasSuffix(call.Callee, ".Store") || !IsInLoop(call) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF169, file, line, col,
			"atomic.Value.Store inside a loop allocates an interface{} per call; use atomic.Pointer[T] (Go 1.19+) for frequent updates", out)
		return
	}
}

func detectPERF170(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "sync.Once") || !strings.Contains(source, ".Do(") {
		return
	}
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], ".Do(")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		preStart := start - 16
		if preStart < 0 {
			preStart = 0
		}
		if !strings.Contains(strings.ToLower(source[preStart:start]), "once") {
			searchFrom = start + len(".Do(")
			continue
		}
		funcStart := start - 1024
		if funcStart < 0 {
			funcStart = 0
		}
		funcWindow := source[funcStart:start]
		if !(strings.Contains(funcWindow, "http.ResponseWriter") ||
			strings.Contains(funcWindow, "gin.Context") ||
			strings.Contains(funcWindow, "echo.Context") ||
			strings.Contains(funcWindow, "*fiber.Ctx")) {
			searchFrom = start + len(".Do(")
			continue
		}
		line, col := unit.LineCol(start)
		rules.PushFinding(&MetaPERF170, file, line, col,
			"sync.Once.Do in a request handler; use a sync/atomic.Bool or hoist the once out of the request path", out)
		searchFrom = start + len(".Do(")
	}
}

func detectPERF171(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "make" || len(call.Arguments) < 2 {
			continue
		}
		chanType := strings.TrimSpace(call.Arguments[0])
		isMutexShape := chanType == "chan struct{}" ||
			chanType == "chan bool" ||
			chanType == "chan struct{ }" ||
			strings.HasPrefix(chanType, "chan struct{},") ||
			strings.HasPrefix(chanType, "chan bool,") ||
			strings.Contains(chanType, "chan struct{},") ||
			strings.Contains(chanType, "chan bool,")
		if !isMutexShape || strings.TrimSpace(call.Arguments[1]) != "1" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF171, file, line, col,
			"make(chan T, 1) used as a mutex; use sync.Mutex instead of a channel", out)
	}
}

func detectPERF172(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !b4_fileHasHandler(source) && !strings.Contains(source, "http.ResponseWriter") {
		return
	}
	if !strings.Contains(source, ".Wait()") || !strings.Contains(source, "go ") {
		return
	}
	waitPos := strings.LastIndex(source, ".Wait()")
	if waitPos < 0 {
		return
	}
	tail := source[waitPos+7:]
	if !(strings.Contains(tail, ".JSON(") || strings.Contains(tail, ".Write(") ||
		strings.Contains(tail, ".WriteHeader(") || strings.Contains(tail, ".String(") ||
		strings.Contains(tail, ".HTML(")) {
		return
	}
	windowStart := waitPos - 2048
	if windowStart < 0 {
		windowStart = 0
	}
	window := source[windowStart:waitPos]
	if strings.Contains(window, "ctx := c.Request.Context()") ||
		strings.Contains(window, "ctx := r.Context()") ||
		strings.Contains(window, "ctx, cancel := context.WithCancel") ||
		strings.Contains(window, "ctx := c.Copy()") ||
		strings.Contains(window, "case <-ctx.Done()") {
		return
	}
	if gfp := strings.LastIndex(source[:waitPos], "go func"); gfp >= 0 {
		goBody := source[gfp:waitPos]
		hasWork := false
		for _, l := range strings.Split(goBody, "\n") {
			t := strings.TrimSpace(l)
			if strings.Contains(t, "(") &&
				!strings.Contains(t, "wg.Done") &&
				!strings.Contains(t, "wg.Add") &&
				!strings.Contains(t, "go func") &&
				!strings.Contains(t, "defer func") &&
				t != "}()" && t != "})" && !strings.HasPrefix(t, "}(") {
				hasWork = true
				break
			}
		}
		if hasWork {
			return
		}
	}
	line, col := unit.LineCol(waitPos)
	rules.PushFinding(&MetaPERF172, file, line, col,
		"wg.Wait in a request handler blocks the serving goroutine; use context cancellation or errgroup instead", out)
}

func detectPERF173(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		if call.Callee != "time.Tick" {
			continue
		}
		if strings.Contains(source[call.StartByte:], ".Stop()") {
			continue
		}
		preStart := call.StartByte - 32
		if preStart < 0 {
			preStart = 0
		}
		if strings.Contains(source[preStart:call.StartByte], "range ") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF173, file, line, col,
			"time.Tick returns an unstoppable ticker; use time.NewTicker and call ticker.Stop() when done", out)
	}
}

func detectPERF174(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		if call.Callee != "close" || len(call.Arguments) == 0 {
			continue
		}
		ch := strings.TrimSpace(strings.TrimSuffix(call.Arguments[0], ","))
		if ch == "" {
			continue
		}
		funcStart := strings.LastIndex(source[:call.StartByte], "func ")
		if funcStart < 0 {
			funcStart = 0
		}
		bodyStartRel := strings.Index(source[funcStart:], "{")
		if bodyStartRel < 0 {
			continue
		}
		bodyStart := funcStart + bodyStartRel + 1
		afterBody := source[bodyStart:]
		depth := 1
		bodyEnd := len(afterBody)
		for i, c := range afterBody {
			switch c {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					bodyEnd = i
					goto done
				}
			}
		}
	done:
		body := afterBody[:bodyEnd]
		if strings.Contains(body, "<-"+ch) || strings.Contains(body, "<- "+ch) {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF174, file, line, col,
				"close() called on a channel that is also received on in the same function; only the sender should close a channel", out)
		}
	}
}

func detectPERF175(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "for {") || !strings.Contains(source, "<-") {
		return
	}
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], "for {")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		end := start + 1024
		if end > len(source) {
			end = len(source)
		}
		bodyWindow := source[start:end]
		if !strings.Contains(bodyWindow, "<-") ||
			strings.Contains(bodyWindow, "select ") ||
			strings.Contains(bodyWindow, "time.Sleep(") {
			searchFrom = start + len("for {")
			continue
		}
		openBrace := strings.Index(bodyWindow, "{")
		if openBrace < 0 {
			searchFrom = start + len("for {")
			continue
		}
		if strings.Contains(bodyWindow[openBrace+1:], "<-") {
			line, col := unit.LineCol(start)
			rules.PushFinding(&MetaPERF175, file, line, col,
				"for { ... <-ch ... } spins on a buffered channel; use a select with default or wait on a done channel", out)
			return
		}
		searchFrom = start + len("for {")
	}
}

func detectPERF176(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "io.Copy" || !IsInLoop(call) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF176, file, line, col,
			"io.Copy inside a loop allocates a 32 KiB buffer per call; use io.CopyBuffer with a pooled buffer", out)
	}
}

func detectPERF177(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, ".Readdir(") {
		return
	}
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], ".Readdir(")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		line, col := unit.LineCol(start)
		rules.PushFinding(&MetaPERF177, file, line, col,
			"(*os.File).Readdir returns []os.FileInfo; prefer os.ReadDir for []os.DirEntry", out)
		searchFrom = start + len(".Readdir(")
	}
}

func detectPERF178(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	var formats []CallFact
	for _, c := range facts.Calls {
		if strings.HasSuffix(c.Callee, ".Format") && !strings.HasSuffix(c.Callee, "AppendFormat") {
			formats = append(formats, c)
		}
	}
	if len(formats) < 2 {
		return
	}
	for i := 0; i+1 < len(formats); i++ {
		a, b := formats[i], formats[i+1]
		if a.Callee != b.Callee || b.StartByte-a.StartByte > 1024 {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(&MetaPERF178, file, line, col,
			"time.Format called repeatedly with the same format string; use time.AppendFormat to write into a pooled buffer", out)
		return
	}
}

func detectPERF179(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !strings.Contains(unit.Source, "strings.Replace(") {
		return
	}
	type pair struct {
		key string
		pos int
	}
	var keys []pair
	for _, call := range facts.Calls {
		if call.Callee != "strings.Replace" || len(call.Arguments) < 3 {
			continue
		}
		keys = append(keys, pair{key: call.Arguments[1] + "\x01" + call.Arguments[2], pos: call.StartByte})
	}
	for i := 0; i+1 < len(keys); i++ {
		if keys[i].key != keys[i+1].key || keys[i+1].pos-keys[i].pos > 2048 {
			continue
		}
		line, col := unit.LineCol(keys[i].pos)
		rules.PushFinding(&MetaPERF179, file, line, col,
			"strings.Replace with the same old/new pair called repeatedly; build a strings.Replacer once and reuse it", out)
		return
	}
}

func detectPERF180(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !strings.Contains(unit.Source, "csv.NewReader") {
		return
	}
	for _, call := range facts.Calls {
		if !strings.HasSuffix(call.Callee, ".Read") || !IsInLoop(call) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF180, file, line, col,
			"csv.Reader.Read called inside a loop; reuse a single reader and consider ReadAll for bulk parsing", out)
		return
	}
}

func detectPERF181(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "json.NewDecoder") {
		return
	}
	if !strings.Contains(source, "int") && !strings.Contains(source, "int64") && !strings.Contains(source, "int32") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "json.NewDecoder" {
			continue
		}
		after := call.StartByte + 256
		if after > len(source) {
			after = len(source)
		}
		if strings.Contains(source[call.StartByte:after], ".UseNumber()") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF181, file, line, col,
			"json.NewDecoder without .UseNumber() silently loses precision for large integers", out)
	}
}

func detectPERF182(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "bufio.NewWriter") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "bufio.NewWriter" || len(call.Arguments) != 1 {
			continue
		}
		afterEnd := call.StartByte + 512
		if afterEnd > len(source) {
			afterEnd = len(source)
		}
		window := source[call.StartByte:afterEnd]
		if !strings.Contains(window, ".Write(") && !strings.Contains(window, ".WriteString(") {
			continue
		}
		if !b4_hasLargeStringLiteral(window) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF182, file, line, col,
			"bufio.NewWriter without an explicit buffer size; the default 4 KiB buffer thrashes on large writes", out)
	}
}

func detectPERF183(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		switch call.Callee {
		case "context.WithTimeout", "context.WithDeadline", "context.WithCancel":
		default:
			continue
		}
		if !IsInLoop(call) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF183, file, line, col,
			"context.WithTimeout inside a loop; create the context once outside the loop and derive per-iteration values via context.WithValue", out)
	}
}

func detectPERF184(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "mime.TypeByExtension" {
			continue
		}
		if !IsInLoop(call) && !IsHandlerShaped(unit.Source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF184, file, line, col,
			"mime.TypeByExtension walks the mime.types table; cache the result for the extensions you handle", out)
	}
}

func detectPERF185(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "http.DetectContentType" || !IsHandlerShaped(unit.Source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF185, file, line, col,
			"http.DetectContentType in a request handler; parse the Content-Type header or cache the result for the bodies you serve", out)
	}
}

func detectPERF186(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "strings.Fields" {
			continue
		}
		if !IsInLoop(call) && !IsHandlerShaped(unit.Source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF186, file, line, col,
			"strings.Fields in a hot path; use strings.IndexByte to walk whitespace and slice once per token", out)
	}
}

func detectPERF187(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "template.HTMLEscaper" && call.Callee != "HTMLEscaper" {
			continue
		}
		if !IsInLoop(call) && !IsHandlerShaped(unit.Source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF187, file, line, col,
			"template.HTMLEscaper in a hot path; pre-escape at write time or use template.HTML when the input is trusted", out)
	}
}

func detectPERF188(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "fmt.Sscanf" {
			continue
		}
		if !IsInLoop(call) && !IsHandlerShaped(unit.Source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF188, file, line, col,
			"fmt.Sscanf in a hot path; use strconv.ParseInt / strconv.ParseFloat for the common conversions", out)
	}
}

func detectPERF189(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "io.Copy(io.Discard,") {
		return
	}
	drainPos := strings.Index(source, "io.Copy(io.Discard,")
	closePos := strings.Index(source, ".Body.Close()")
	if closePos > 0 && closePos < drainPos {
		line, col := unit.LineCol(closePos)
		rules.PushFinding(&MetaPERF189, file, line, col,
			"Body.Close called before io.Copy(io.Discard, body); drain BEFORE close to allow keep-alive connection reuse", out)
	}
}

func detectPERF190(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], "http.Client{")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		lineStart := strings.LastIndex(source[:start], "\n")
		if lineStart < 0 {
			lineStart = 0
		} else {
			lineStart++
		}
		if strings.HasPrefix(strings.TrimLeft(source[lineStart:start], " \t"), "var ") {
			searchFrom = start + len("http.Client{")
			continue
		}
		end := start + 192
		if end > len(source) {
			end = len(source)
		}
		if strings.Contains(source[start:end], "Timeout:") {
			searchFrom = start + len("http.Client{")
			continue
		}
		line, col := unit.LineCol(start)
		rules.PushFinding(&MetaPERF190, file, line, col,
			"http.Client is missing Timeout; requests can hang indefinitely", out)
		searchFrom = start + len("http.Client{")
	}
}

func detectPERF191(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if strings.Contains(source, "proto.") || strings.Contains(source, "protobuf") {
		return
	}
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], "[]*")
		if rel < 0 {
			return
		}
		pos := searchFrom + rel
		afterEnd := pos + 3 + 128
		if afterEnd > len(source) {
			afterEnd = len(source)
		}
		after := source[pos+3 : afterEnd]
		typeName := ""
		for _, r := range after {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				typeName += string(r)
			} else {
				break
			}
		}
		if typeName == "" {
			searchFrom = pos + 3
			continue
		}
		pattern := "type " + typeName + " struct"
		structStart := strings.Index(source, pattern)
		if structStart >= 0 {
			openRel := strings.Index(source[structStart:], "{")
			if openRel >= 0 {
				bodyStart := structStart + openRel + 1
				closeRel := strings.Index(source[bodyStart:], "}")
				if closeRel >= 0 {
					body := source[bodyStart : bodyStart+closeRel]
					fieldCount := 0
					for _, l := range strings.Split(body, "\n") {
						t := strings.TrimSpace(l)
						if t == "" || strings.HasPrefix(t, "//") ||
							strings.HasPrefix(t, "{") || strings.HasPrefix(t, "}") ||
							strings.HasPrefix(t, "`") {
							continue
						}
						fieldCount++
					}
					if fieldCount > 0 && fieldCount <= 2 {
						line, col := unit.LineCol(pos)
						rules.PushFinding(&MetaPERF191, file, line, col,
							"slice of pointers to a small struct; use []T (value type) to avoid per-element heap allocations", out)
						return
					}
				}
			}
		}
		searchFrom = pos + 3
	}
}

func detectPERF192(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		if call.Callee != "make" || len(call.Arguments) != 1 {
			continue
		}
		if !strings.HasPrefix(call.Arguments[0], "map[") {
			continue
		}
		if !b4_mapSizeHintAvailable(source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF192, file, line, col,
			"make(map[K]V) without a size hint; pass len(src) to avoid map growth", out)
		return
	}
}

func detectPERF193(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		if call.Callee != "time.After" || !IsInLoop(call) {
			continue
		}
		end := call.StartByte + 1024
		if end > len(source) {
			end = len(source)
		}
		if strings.Contains(source[call.StartByte:end], ".Reset(") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF193, file, line, col,
			"time.After inside a loop without a reusable timer; hoist a *time.Timer and call t.Reset(d) per iteration", out)
	}
}

func detectPERF194(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "time.Sleep" || !IsHandlerShaped(unit.Source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF194, file, line, col,
			"time.Sleep in a request handler; use a channel / context cancellation or a longer-lived ticker to poll", out)
	}
}

func detectPERF195(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		switch call.Callee {
		case "log.Fatal", "log.Fatalf", "log.Fatalln", "log.Panic", "log.Panicf", "log.Panicln":
		default:
			continue
		}
		windowStart := call.StartByte - 2048
		if windowStart < 0 {
			windowStart = 0
		}
		if !strings.Contains(source[windowStart:call.StartByte], "go func()") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF195, file, line, col,
			"log.Fatal / log.Panic inside a goroutine; return the error and let the caller decide whether to terminate the process", out)
	}
}

func detectPERF196(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !b4_fileHasHandler(source) && !strings.Contains(source, "http.ResponseWriter") {
		return
	}
	triggers := []string{"jwt.Parse(", "jwt.ParseWithClaims(", "session.Get(", "sessions.Get(", "cookie.Get("}
	any := false
	for _, t := range triggers {
		if strings.Contains(source, t) {
			any = true
			break
		}
	}
	if !any {
		return
	}
	for _, trigger := range triggers {
		rel := strings.Index(source, trigger)
		if rel < 0 {
			continue
		}
		preStart := rel - 512
		if preStart < 0 {
			preStart = 0
		}
		pre := source[preStart:rel]
		if strings.Contains(pre, "func AuthMiddleware") ||
			strings.Contains(pre, "func SessionMiddleware") ||
			strings.Contains(pre, "func Middleware") ||
			strings.Contains(pre, "func (h *Handler)") ||
			strings.Contains(pre, "func Authenticate") {
			continue
		}
		line, col := unit.LineCol(rel)
		rules.PushFinding(&MetaPERF196, file, line, col,
			"session / JWT parse in a request handler; cache the parsed session for the duration of the request", out)
		return
	}
}

func detectPERF197(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	var reads []CallFact
	for _, c := range facts.Calls {
		if c.Callee == "io.ReadAll" || c.Callee == "ioutil.ReadAll" {
			reads = append(reads, c)
		}
	}
	if len(reads) < 2 {
		return
	}
	for i := 0; i+1 < len(reads); i++ {
		a, b := reads[i], reads[i+1]
		if len(a.Arguments) == 0 || len(b.Arguments) == 0 {
			continue
		}
		if a.Arguments[0] == b.Arguments[0] &&
			(strings.Contains(a.Arguments[0], "Body") || strings.Contains(a.Arguments[0], "body")) {
			line, col := unit.LineCol(b.StartByte)
			rules.PushFinding(&MetaPERF197, file, line, col,
				"io.ReadAll(c.Request.Body) called twice; the second read returns EOF, cache the body or read into a buffer", out)
			return
		}
	}
}

func detectPERF198(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "strings.Contains" || len(call.Arguments) == 0 {
			continue
		}
		first := call.Arguments[0]
		if !strings.Contains(first, "Content-Type") && !strings.Contains(first, "contentType") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF198, file, line, col,
			"Content-Type checks should parse or compare the media type instead of using strings.Contains", out)
	}
}

func detectPERF199(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	hasSession := strings.Contains(source, "session.Get(") ||
		strings.Contains(source, "sessions.Get(") ||
		strings.Contains(source, "c.Cookie(") ||
		strings.Contains(source, "r.Cookie(") ||
		strings.Contains(source, "cookie.Get(") ||
		strings.Contains(source, "rdb.Get(") ||
		strings.Contains(source, "redis.Get(")
	if !hasSession {
		return
	}
	if !b4_fileHasHandler(source) && !strings.Contains(source, "http.ResponseWriter") {
		return
	}
	if strings.Contains(source, ".Use(") || strings.Contains(source, "Group.Use(") {
		return
	}
	triggers := []string{"c.Cookie(", "r.Cookie(", "session.Get(", "sessions.Get(", "cookie.Get(", "rdb.Get(", "redis.Get("}
	for _, trigger := range triggers {
		pos := strings.Index(source, trigger)
		if pos < 0 || !IsHandlerShaped(source, pos) {
			continue
		}
		funcStart := strings.LastIndex(source[:pos], "func ")
		if funcStart < 0 {
			funcStart = 0
		}
		funcEnd := pos + 1024
		if funcEnd > len(source) {
			funcEnd = len(source)
		}
		funcWindow := source[funcStart:funcEnd]
		if strings.Contains(funcWindow, "c.Next()") ||
			strings.Contains(funcWindow, "gin.HandlerFunc") ||
			strings.Contains(funcWindow, "Middleware") ||
			strings.Contains(funcWindow, "AuthMiddleware") {
			continue
		}
		line, col := unit.LineCol(pos)
		rules.PushFinding(&MetaPERF199, file, line, col,
			"session lookup in a route handler; move the lookup to a middleware that sets the request context", out)
		return
	}
}

func detectPERF200(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, ".Use(") {
		return
	}
	authPos := b4_indexAny(source, "Auth", "auth.", "RequireAuth", "Authenticate", "JWT", "RateLimit")
	corsPos := b4_indexAny(source, "CORS", "cors.", "Cache")
	if authPos >= 0 && corsPos >= 0 && authPos < corsPos {
		line, col := unit.LineCol(corsPos)
		rules.PushFinding(&MetaPERF200, file, line, col,
			"expensive middleware (Auth) registered before cheap preflight (CORS); move CORS first to short-circuit preflight requests", out)
	}
}

func detectPERF201(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "OPTIONS") || !strings.Contains(source, "Access-Control-") {
		return
	}
	if strings.Contains(source, "github.com/gin-contrib/cors") || strings.Contains(source, "cors.New(") {
		return
	}
	pos := strings.Index(source, "OPTIONS")
	if pos < 0 {
		return
	}
	line, col := unit.LineCol(pos)
	rules.PushFinding(&MetaPERF201, file, line, col,
		"custom CORS preflight handler; use a community package (cors, gin-contrib/cors) for the standard preflight", out)
}

func detectPERF202(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee == "json.MarshalIndent" {
			if !IsHandlerShaped(unit.Source, call.StartByte) {
				continue
			}
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF202, file, line, col,
				"json.MarshalIndent in a request handler; use json.Marshal for compact output in production", out)
			continue
		}
		if strings.HasSuffix(call.Callee, ".SetIndent") && IsHandlerShaped(unit.Source, call.StartByte) {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF202, file, line, col,
				"json.Encoder.SetIndent in a request handler; indentation doubles the response size and slows down marshalling", out)
		}
	}
}

func detectPERF203(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	var stringsCalls []CallFact
	for _, c := range facts.Calls {
		if strings.HasSuffix(c.Callee, ".String") {
			stringsCalls = append(stringsCalls, c)
		}
	}
	if len(stringsCalls) < 2 {
		return
	}
	for i := 0; i+1 < len(stringsCalls); i++ {
		a, b := stringsCalls[i], stringsCalls[i+1]
		if a.Callee != b.Callee || b.StartByte-a.StartByte > 1024 {
			continue
		}
		callee := a.Callee
		if strings.Contains(strings.ToLower(callee), "ip.") || strings.HasPrefix(callee, "ip.") {
			line, col := unit.LineCol(a.StartByte)
			rules.PushFinding(&MetaPERF203, file, line, col,
				"ip.String() called repeatedly on the same IP; cache the result or write directly to a buffer", out)
			return
		}
	}
}

func detectPERF204(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, ".Updates(") {
		return
	}
	for _, call := range facts.Calls {
		if !strings.HasSuffix(call.Callee, ".Updates") {
			continue
		}
		if len(call.Arguments) == 0 || !strings.Contains(call.Arguments[0], "map[") {
			continue
		}
		windowStart := call.StartByte - 256
		if windowStart < 0 {
			windowStart = 0
		}
		window := source[windowStart:call.StartByte]
		if idx := strings.LastIndex(window, ".Select("); idx >= 0 {
			before := window[:idx]
			stmtStart := 0
			if i := strings.LastIndex(before, "\n"); i >= 0 {
				stmtStart = i + 1
			}
			if j := strings.LastIndex(before, ";"); j >= 0 && j+1 > stmtStart {
				stmtStart = j + 1
			}
			if stmtStart > 0 {
				continue
			}
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF204, file, line, col,
			"db.Updates(map) without a preceding .Select; GORM will UPDATE every column", out)
	}
}

func detectPERF205(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	hasCount := strings.Contains(source, "db.Count(") || strings.Contains(source, ".Count(&")
	if !hasCount || !strings.Contains(source, ".Offset(") ||
		!strings.Contains(source, ".Limit(") || !strings.Contains(source, ".Find(") {
		return
	}
	countPos := strings.Index(source, ".Count(")
	findPos := strings.Index(source, ".Find(")
	if countPos <= 0 || findPos <= countPos || findPos-countPos > 2048 {
		return
	}
	line, col := unit.LineCol(countPos)
	rules.PushFinding(&MetaPERF205, file, line, col,
		"db.Count + db.Offset.Limit.Find pattern; use keyset pagination (where id > last_id) for large tables", out)
}

func detectPERF206(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "Unsafe(") {
		return
	}
	for _, call := range facts.Calls {
		callee := call.Callee
		if !strings.HasSuffix(callee, ".Where") &&
			!strings.HasSuffix(callee, ".Find") &&
			!strings.HasSuffix(callee, ".First") {
			continue
		}
		if len(call.Arguments) == 0 {
			continue
		}
		first := call.Arguments[0]
		if strings.HasPrefix(first, "\"") {
			continue
		}
		if strings.Contains(first, "+ \"") ||
			strings.Contains(first, "\" +") ||
			strings.Contains(first, "fmt.Sprintf(") ||
			(!strings.HasPrefix(first, "\"") && strings.Contains(first, "\"")) {
			if strings.Contains(callee, "Unsafe") {
				line, col := unit.LineCol(call.StartByte)
				rules.PushFinding(&MetaPERF206, file, line, col,
					"sqlx.Unsafe used with a non-literal query; use a static string for the query when in Unsafe mode", out)
				return
			}
		}
	}
}

func detectPERF207(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "c.SendFile(") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "c.SendFile" && call.Callee != "SendFile" {
			continue
		}
		if !IsHandlerShaped(source, call.StartByte) {
			continue
		}
		ws := call.StartByte - 512
		if ws < 0 {
			ws = 0
		}
		we := call.StartByte + 512
		if we > len(source) {
			we = len(source)
		}
		window := source[ws:we]
		if strings.Contains(window, "Cache-Control") ||
			strings.Contains(window, "ETag") ||
			strings.Contains(window, "Last-Modified") ||
			strings.Contains(window, "CacheControl") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF207, file, line, col,
			"c.SendFile without Cache-Control / ETag / Last-Modified headers; set cache headers to allow downstream caching", out)
	}
}

func detectPERF209(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "cobra.Command") {
		return
	}
	if !strings.Contains(source, "PersistentPreRunE") && !strings.Contains(source, "PersistentPostRunE") {
		return
	}
	for _, marker := range []string{"PersistentPreRunE", "PersistentPostRunE"} {
		from := 0
		for {
			rel := strings.Index(source[from:], marker)
			if rel < 0 {
				break
			}
			start := from + rel
			pre := source[:start]
			lastNL := strings.LastIndex(pre, "\n")
			if lastNL < 0 {
				lastNL = -1
			}
			between := source[lastNL+1 : start]
			onlyWS := true
			for _, r := range between {
				if r != ' ' && r != '\t' {
					onlyWS = false
					break
				}
			}
			if !onlyWS {
				from = start + len(marker)
				continue
			}
			lineStart := lastNL + 1
			if lineStart < 0 {
				lineStart = 0
			}
			if strings.HasPrefix(strings.TrimLeft(source[lineStart:start], " \t"), "//") {
				from = start + len(marker)
				continue
			}
			line, col := unit.LineCol(start)
			rules.PushFinding(&MetaPERF209, file, line, col,
				"PersistentPreRunE / PersistentPostRunE runs for every subcommand; use a sync.Once or pre-build the dependency", out)
			from = start + len(marker)
		}
	}
}

func detectPERF210(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		callee := call.Callee
		if !strings.HasSuffix(callee, ".Keys") && callee != "Keys" {
			continue
		}
		if !IsHandlerShaped(unit.Source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF210, file, line, col,
			"redis KEYS command in a request handler; use SCAN for incremental iteration to avoid blocking the Redis server", out)
	}
}

func detectPERF211(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	hasNot := strings.Contains(source, "db.Not(") || strings.Contains(source, ".Not(")
	hasNotIn := strings.Contains(source, "NOT IN") || strings.Contains(source, "not in")
	hasNotLike := strings.Contains(source, "NOT LIKE") || strings.Contains(source, "not like")
	if !hasNot && !hasNotIn && !hasNotLike {
		return
	}
	for _, call := range facts.Calls {
		callee := call.Callee
		if strings.HasSuffix(callee, ".Not") && len(call.Arguments) > 0 {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF211, file, line, col,
				"db.Not(...) defeats index lookups; rewrite as a positive WHERE clause", out)
			continue
		}
		if strings.HasSuffix(callee, ".Where") {
			for _, arg := range call.Arguments {
				up := strings.ToUpper(arg)
				if strings.Contains(up, "NOT IN") || strings.Contains(up, "NOT LIKE") {
					line, col := unit.LineCol(call.StartByte)
					rules.PushFinding(&MetaPERF211, file, line, col,
						"NOT IN / NOT LIKE defeats index lookups; use a positive WHERE clause", out)
					break
				}
			}
		}
	}
}

func detectPERF212(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		callee := call.Callee
		if callee != "db.Find" && !strings.HasSuffix(callee, ".Find") {
			continue
		}
		if len(call.Arguments) == 0 {
			continue
		}
		first := strings.TrimLeft(call.Arguments[0], " \t")
		if !strings.HasPrefix(first, "&") {
			continue
		}
		afterAmp := strings.TrimLeft(strings.TrimPrefix(first, "&"), " \t")
		ident := ""
		for _, r := range afterAmp {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				ident += string(r)
			} else {
				break
			}
		}
		if ident == "" {
			continue
		}
		decls := []string{"var " + ident + " []", ident + " := []", ident + " := make([]"}
		ok := false
		for _, d := range decls {
			if strings.Contains(source, d) {
				ok = true
				break
			}
		}
		if !ok {
			continue
		}
		stmtStart := strings.LastIndex(source[:call.StartByte], "\n")
		if stmtStart < 0 {
			stmtStart = 0
		} else {
			stmtStart++
		}
		stmt := source[stmtStart:call.StartByte]
		combined := stmt + callee
		if strings.Contains(combined, "Limit(") || strings.Contains(combined, "Preload(") ||
			strings.Contains(combined, "Joins(") || strings.Contains(combined, "Select(") ||
			strings.Contains(combined, "Where(") || strings.Contains(combined, "Not(") ||
			strings.Contains(combined, "Order(") || strings.Contains(combined, "Group(") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF212, file, line, col,
			"db.Find(&slice) without a preceding .Limit; bound the result set on tables that can grow unbounded", out)
	}
}

func detectPERF213(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, cache := range b4_packageLevelCaches(source) {
		if !strings.Contains(strings.ToLower(cache.name), "cache") {
			continue
		}
		if !b4_cacheIsWritten(source, facts, cache.name, cache.isSyncMap) {
			continue
		}
		if b4_cacheHasEvictionBound(source, cache.name) {
			continue
		}
		line, col := unit.LineCol(cache.byteOff)
		rules.PushFinding(&MetaPERF213, file, line, col,
			"package-level cache has writes but no eviction or entry bound in the same compilation unit", out)
		return
	}
}

func detectPERF214(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !strings.Contains(unit.Source, "sync.Map") {
		return
	}
	for _, call := range facts.Calls {
		method := b4_methodName(call.Callee)
		if method != "Load" && method != "Store" && method != "LoadOrStore" {
			continue
		}
		if len(call.Arguments) == 0 {
			continue
		}
		if b4_cacheKeyIsVolatile(unit, facts, call.StartByte, call.Arguments[0]) {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF214, file, line, col,
				"cache key includes request-scoped or volatile fields, which collapses cache hit rate", out)
			return
		}
	}
}

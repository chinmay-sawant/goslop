package perf

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("PERF-61", detectPERF61, &MetaPERF61)
	RegisterRule("PERF-62", detectPERF62, &MetaPERF62)
	RegisterRule("PERF-63", detectPERF63, &MetaPERF63)
	RegisterRule("PERF-64", detectPERF64, &MetaPERF64)
	RegisterRule("PERF-65", detectPERF65, &MetaPERF65)
	RegisterRule("PERF-66", detectPERF66, &MetaPERF66)
	RegisterRule("PERF-67", detectPERF67, &MetaPERF67)
	RegisterRule("PERF-68", detectPERF68, &MetaPERF68)
	RegisterRule("PERF-69", detectPERF69, &MetaPERF69)
	RegisterRule("PERF-70", detectPERF70, &MetaPERF70)
	RegisterRule("PERF-71", detectPERF71, &MetaPERF71)
	RegisterRule("PERF-72", detectPERF72, &MetaPERF72)
	RegisterRule("PERF-73", detectPERF73, &MetaPERF73)
	RegisterRule("PERF-74", detectPERF74, &MetaPERF74)
	RegisterRule("PERF-75", detectPERF75, &MetaPERF75)
	RegisterRule("PERF-76", detectPERF76, &MetaPERF76)
	RegisterRule("PERF-77", detectPERF77, &MetaPERF77)
	RegisterRule("PERF-78", detectPERF78, &MetaPERF78)
	RegisterRule("PERF-79", detectPERF79, &MetaPERF79)
	RegisterRule("PERF-80", detectPERF80, &MetaPERF80)
	RegisterRule("PERF-81", detectPERF81, &MetaPERF81)
	RegisterRule("PERF-82", detectPERF82, &MetaPERF82)
	RegisterRule("PERF-83", detectPERF83, &MetaPERF83)
	RegisterRule("PERF-84", detectPERF84, &MetaPERF84)
	RegisterRule("PERF-85", detectPERF85, &MetaPERF85)
	RegisterRule("PERF-86", detectPERF86, &MetaPERF86)
	RegisterRule("PERF-87", detectPERF87, &MetaPERF87)
	RegisterRule("PERF-88", detectPERF88, &MetaPERF88)
	RegisterRule("PERF-89", detectPERF89, &MetaPERF89)
	RegisterRule("PERF-90", detectPERF90, &MetaPERF90)
	RegisterRule("PERF-91", detectPERF91, &MetaPERF91)
	RegisterRule("PERF-92", detectPERF92, &MetaPERF92)
	RegisterRule("PERF-93", detectPERF93, &MetaPERF93)
	RegisterRule("PERF-94", detectPERF94, &MetaPERF94)
	RegisterRule("PERF-95", detectPERF95, &MetaPERF95)
	RegisterRule("PERF-96", detectPERF96, &MetaPERF96)
	RegisterRule("PERF-97", detectPERF97, &MetaPERF97)
	RegisterRule("PERF-98", detectPERF98, &MetaPERF98)
	RegisterRule("PERF-99", detectPERF99, &MetaPERF99)
	RegisterRule("PERF-100", detectPERF100, &MetaPERF100)
	RegisterRule("PERF-101", detectPERF101, &MetaPERF101)
	RegisterRule("PERF-102", detectPERF102, &MetaPERF102)
	RegisterRule("PERF-103", detectPERF103, &MetaPERF103)
	RegisterRule("PERF-105", detectPERF105, &MetaPERF105)
	RegisterRule("PERF-106", detectPERF106, &MetaPERF106)
	RegisterRule("PERF-107", detectPERF107, &MetaPERF107)
	RegisterRule("PERF-108", detectPERF108, &MetaPERF108)
	RegisterRule("PERF-109", detectPERF109, &MetaPERF109)
	RegisterRule("PERF-110", detectPERF110, &MetaPERF110)
	RegisterRule("PERF-111", detectPERF111, &MetaPERF111)
}

func b2HasAny(src string, needles []string) bool {
	for _, n := range needles {
		if n != "" && strings.Contains(src, n) {
			return true
		}
	}
	return false
}

func b2FirstPos(src string, needles []string) int {
	best := -1
	for _, n := range needles {
		if n == "" {
			continue
		}
		if i := strings.Index(src, n); i >= 0 && (best < 0 || i < best) {
			best = i
		}
	}
	if best < 0 {
		return 0
	}
	return best
}

func b2TopCommas(s string) int {
	depth, count := 0, 0
	for _, c := range s {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				count++
			}
		}
	}
	return count
}

func b2CallInLoop(facts *GoPerfFacts, needles []string) (int, bool) {
	for _, c := range facts.Calls {
		if !IsInLoop(c) {
			continue
		}
		for _, n := range needles {
			if strings.Contains(c.Callee, n) {
				return c.StartByte, true
			}
		}
	}
	return 0, false
}

func b2Emit(unit *core.ParsedUnit, meta *rules.RuleMetadata, pos int, msg string, out *[]rules.Finding) {
	file := unitFile(unit)
	line, col := unit.LineCol(pos)
	rules.PushFinding(meta, file, line, col, msg, out)
}

func b2IsRequestPath(src string, facts *GoPerfFacts) bool {
	if IsRequestPath(facts.SourceIndex) {
		return true
	}
	return b2HasAny(src, []string{
		"gin.HandlerFunc", "echo.HandlerFunc", "http.HandlerFunc",
		"http.ResponseWriter", "*gin.Context", "gin.Context",
		"echo.Context", "*fiber.Ctx", "fiber.Ctx",
	})
}

func b2EnclosingFuncBody(source string, startByte int) string {
	if startByte > len(source) {
		startByte = len(source)
	}
	if startByte < 0 {
		startByte = 0
	}
	head := source[:startByte]
	funcKw := strings.LastIndex(head, "func ")
	if funcKw < 0 {
		return source[startByte:]
	}
	rel := strings.Index(source[funcKw:startByte], "{")
	if rel < 0 {
		return source[startByte:]
	}
	bodyOpen := funcKw + rel
	depth := 0
	for i := bodyOpen; i < len(source); i++ {
		switch source[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				if startByte >= bodyOpen && startByte <= i {
					return source[bodyOpen : i+1]
				}
				break
			}
		}
	}
	return source[bodyOpen:startByte]
}

func b2IsIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func b2WordBoundaryContains(s, word string) bool {
	if word == "" {
		return false
	}
	from := 0
	for {
		i := strings.Index(s[from:], word)
		if i < 0 {
			return false
		}
		abs := from + i
		leftOK := abs == 0 || !b2IsIdentByte(s[abs-1])
		right := abs + len(word)
		rightOK := right >= len(s) || !b2IsIdentByte(s[right])
		if leftOK && rightOK {
			return true
		}
		from = abs + 1
	}
}

func b2WriteHeaderFollowedByReturn(src string, start int) bool {
	rest := src[start:]
	depth := 0
	end := -1
	for i, r := range rest {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				end = i
			}
		}
		if end >= 0 {
			break
		}
	}
	if end < 0 {
		return false
	}
	after := strings.TrimSpace(rest[end+1:])
	after = strings.TrimPrefix(after, ";")
	after = strings.TrimSpace(after)
	return strings.HasPrefix(after, "return")
}

func b2MethodName(callee string) string {
	if i := strings.LastIndex(callee, "."); i >= 0 {
		return callee[i+1:]
	}
	return callee
}

func detectPERF61(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	trig := []string{"gin.Static(", "router.Static(", "r.Static(", "c.File("}
	if !b2HasAny(src, trig) {
		return
	}
	if b2HasAny(src, []string{"Cache-Control", "cacheControl", "MaxAge", "Max-Age", `c.Header("ETag"`}) {
		return
	}
	b2Emit(unit, &MetaPERF61, b2FirstPos(src, trig),
		"static file served without Cache-Control / ETag headers; configure cache headers or front with a CDN", out)
}

func detectPERF62(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2IsRequestPath(src, facts) || !strings.Contains(src, "c.Param(") {
		return
	}
	if !(strings.Contains(src, "regexp.MustCompile(") ||
		strings.Contains(src, "regexp.Compile(") ||
		strings.Contains(src, "json.Unmarshal(")) {
		return
	}
	b2Emit(unit, &MetaPERF62, strings.Index(src, "c.Param("),
		"complex c.Param parsing (regex / json.Unmarshal) lives in middleware; move to the route handler that needs it", out)
}

func detectPERF63(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2IsRequestPath(src, facts) || !strings.Contains(src, "binding.Validator.Engine()") {
		return
	}
	if strings.Contains(src, "var engine = binding.Validator.Engine()") ||
		strings.Contains(src, "once.Do(func()") ||
		strings.Contains(src, "sync.Once") {
		return
	}
	b2Emit(unit, &MetaPERF63, strings.Index(src, "binding.Validator.Engine()"),
		"binding.Validator.Engine() is invoked per request; cache the engine at startup", out)
}

func detectPERF64(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "go func()") || strings.Contains(src, "c.Copy()") {
		return
	}
	goPos := strings.Index(src, "go func()")
	if goPos < 0 {
		return
	}
	rest := src[goPos:]
	brace := strings.Index(rest, "{")
	if brace < 0 {
		return
	}
	depth := 0
	bodyEnd := -1
	for i := brace; i < len(rest); i++ {
		switch rest[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				bodyEnd = i
			}
		}
		if bodyEnd >= 0 {
			break
		}
	}
	if bodyEnd < 0 {
		return
	}
	body := rest[brace : bodyEnd+1]
	cMethods := []string{"c.JSON(", "c.AbortWithStatus(", "c.String(", "c.HTML(", "c.Request.", "c.Writer."}
	if !b2HasAny(body, cMethods) {
		return
	}
	b2Emit(unit, &MetaPERF64, goPos,
		"go func(){} uses *gin.Context; call c.Copy() before passing the context to a goroutine", out)
}

func detectPERF65(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	middleware := strings.Contains(src, "r.Use(") ||
		strings.Contains(src, "RouterGroup.Use(") ||
		strings.Contains(src, "routerGroup.Use(") ||
		strings.Contains(src, "engine.Use(")
	if !middleware || !strings.Contains(src, "c.ShouldBind(") {
		return
	}
	b2Emit(unit, &MetaPERF65, strings.Index(src, "c.ShouldBind("),
		"c.ShouldBind runs in middleware registered via .Use(); it parses the body for every route in the chain", out)
}

func detectPERF66(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, ".Use(") {
		return
	}
	searchFrom := 0
	for {
		rel := strings.Index(src[searchFrom:], ".Use(")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		after := start + len(".Use(")
		depth := 1
		closeOff := -1
		for i := 0; i < len(src[after:]); i++ {
			switch src[after+i] {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					closeOff = i
				}
			}
			if closeOff >= 0 {
				break
			}
		}
		if closeOff < 0 {
			return
		}
		if b2TopCommas(src[after:after+closeOff])+1 > 5 {
			b2Emit(unit, &MetaPERF66, start,
				"more than 5 middlewares are passed to a single .Use(...) call; consider splitting into nested groups", out)
			return
		}
		searchFrom = after + closeOff + 1
	}
}

func detectPERF67(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "gin.New()") {
		return
	}
	if b2HasAny(src, []string{
		"gin.Recovery()", "gin.RecoveryWithWriter(",
		"gin.CustomRecovery(", "gin.CustomRecoveryWithWriter(",
	}) {
		return
	}
	b2Emit(unit, &MetaPERF67, strings.Index(src, "gin.New()"),
		"router is created with gin.New() but no gin.Recovery() middleware is installed", out)
}

func detectPERF68(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "gin.Logger") {
		return
	}
	if b2HasAny(src, []string{"Output: io.Discard", "// logger disabled", "LoggerConfig{Output:"}) {
		return
	}
	b2Emit(unit, &MetaPERF68, strings.Index(src, "gin.Logger"),
		"gin.Logger() performs synchronous I/O on the request path; use an async logger or disable in production", out)
}

func detectPERF69(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	trig := []string{"c.Writer.Write(", "c.Stream("}
	if !b2HasAny(src, trig) {
		return
	}
	if strings.Contains(src, "c.Writer.Flush()") || strings.Contains(src, "c.Writer.FlushHeaders()") {
		return
	}
	b2Emit(unit, &MetaPERF69, b2FirstPos(src, trig),
		"c.Writer.Write / c.Stream is used without c.Writer.Flush(); streaming clients see higher time-to-first-byte", out)
}

func detectPERF70(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2IsRequestPath(src, facts) || !strings.Contains(src, "go func()") {
		return
	}
	if b2HasAny(src, []string{
		"sync.WaitGroup", "wg.Add(", "done := make(chan",
		"ctx, cancel := context.WithCancel",
		"ctx, cancel := context.WithTimeout",
		"ctx, cancel := context.WithDeadline",
		"c.Request.Context()", "sync.Once", "errgroup",
		"sem := make(chan", "semaphore", "workerPool", "workerCount",
	}) {
		return
	}
	b2Emit(unit, &MetaPERF70, strings.Index(src, "go func()"),
		"go func(){} in a Gin handler has no WaitGroup / done channel / context cancellation tied to the request", out)
}

func detectPERF71(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if strings.Contains(src, "Preload(") || strings.Contains(src, ".Joins(") {
		return
	}
	if pos, ok := b2CallInLoop(facts, []string{"db.Find", "db.First", "db.Take", "db.Where"}); ok {
		b2Emit(unit, &MetaPERF71, pos,
			"GORM query inside a loop body suggests an N+1 access pattern; use Preload or batch the fetch", out)
	}
}

func detectPERF72(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2IsRequestPath(src, facts) {
		return
	}
	trig := []string{"db.Transaction(", "db.Begin(", "tx := db.Begin(", "tx, err := db.Begin("}
	if !b2HasAny(src, trig) {
		return
	}
	b2Emit(unit, &MetaPERF72, b2FirstPos(src, trig),
		"GORM transaction opened inside a request handler; collapse to a single statement or hoist the work", out)
}

func detectPERF73(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if strings.Contains(src, "Preload(") || strings.Contains(src, ".Joins(") {
		return
	}
	if !b2HasAny(src, []string{"db.Find(", "db.First(", "db.Take("}) {
		return
	}
	relations := []string{".Orders", ".Items", ".Author", ".Comments", ".Profile", ".Children", ".Posts", ".Tags", ".Addresses"}
	if !b2HasAny(src, relations) {
		return
	}
	needle := "db.Find("
	for _, n := range []string{"db.Find(", "db.First(", "db.Take("} {
		if strings.Contains(src, n) {
			needle = n
			break
		}
	}
	b2Emit(unit, &MetaPERF73, strings.Index(src, needle),
		"GORM relation field accessed without Preload/Joins; the relation will not be loaded", out)
}

func detectPERF74(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2IsRequestPath(src, facts) {
		return
	}
	if strings.Contains(src, ".Select(") {
		return
	}
	trig := []string{"db.Find(", "db.First(", "db.Take(", "db.Where("}
	if !b2HasAny(src, trig) {
		return
	}
	b2Emit(unit, &MetaPERF74, b2FirstPos(src, trig),
		"GORM query reads all columns; project only the fields the handler returns with Select", out)
}

func detectPERF75(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2IsRequestPath(src, facts) {
		return
	}
	if !strings.Contains(src, "db.Session(&gorm.Session{") {
		return
	}
	if b2HasAny(src, []string{"var sessionOpts =", "var defaultSession =", "sessionOnce"}) {
		return
	}
	b2Emit(unit, &MetaPERF75, strings.Index(src, "db.Session(&gorm.Session{"),
		"GORM session is constructed per request; hoist Session options to package scope or sync.Once", out)
}

func detectPERF76(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if strings.Contains(src, "CreateInBatches") {
		return
	}
	if pos, ok := b2CallInLoop(facts, []string{"db.Create"}); ok {
		b2Emit(unit, &MetaPERF76, pos,
			"db.Create is called inside a loop; batch with CreateInBatches or hoist the create out", out)
	}
}

func detectPERF77(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "db.Save(") {
		return
	}
	if strings.Contains(src, "db.Create(") || strings.Contains(src, ".Updates(") {
		return
	}
	b2Emit(unit, &MetaPERF77, strings.Index(src, "db.Save("),
		"db.Save in an update-only path; use db.Updates with the changed fields to avoid full-row writes", out)
}

func detectPERF78(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	trig := []string{"db.Raw(", "db.Exec("}
	if !b2HasAny(src, trig) {
		return
	}
	needle := "db.Raw("
	if !strings.Contains(src, needle) {
		needle = "db.Exec("
	}
	afterIdx := strings.Index(src, needle)
	after := src[afterIdx:]
	if !b2HasAny(after, []string{"WHERE", "ORDER BY", "JOIN", "where ", "order by ", "join "}) {
		return
	}
	if b2HasAny(src, []string{"// index-backed", "/* index */", "USING INDEX", "use index"}) {
		return
	}
	b2Emit(unit, &MetaPERF78, afterIdx,
		"Raw/Exec query with WHERE/JOIN/ORDER BY; confirm an index backs the clause", out)
}

func detectPERF79(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	trig := []string{"sql.Open(", "gorm.Open(", "gorm.New(", "postgres.Open(", "mysql.Open("}
	if !b2HasAny(src, trig) {
		return
	}
	if b2HasAny(src, []string{"SetMaxOpenConns(", "SetMaxIdleConns(", "SetConnMaxLifetime("}) {
		return
	}
	b2Emit(unit, &MetaPERF79, b2FirstPos(src, trig),
		"database handle opened without SetMaxOpenConns / SetMaxIdleConns / SetConnMaxLifetime", out)
}

func detectPERF80(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	trig := []string{"Pluck(", "Distinct("}
	if !b2HasAny(src, trig) {
		return
	}
	if strings.Contains(src, ".Limit(") {
		return
	}
	b2Emit(unit, &MetaPERF80, b2FirstPos(src, trig),
		"Pluck/Distinct query has no Limit; bound the result set with Limit, batching, or streaming", out)
}

func detectPERF81(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	trig := []string{"db.Select(", "db.Queryx(", "db.QueryxContext("}
	if !b2HasAny(src, trig) {
		return
	}
	if !strings.Contains(src, "IN (?)") {
		return
	}
	if b2HasAny(src, []string{"chunk", "Chunked", "batchIDs"}) {
		return
	}
	b2Emit(unit, &MetaPERF81, b2FirstPos(src, trig),
		"sqlx IN (?) expands an unbounded slice into one query; chunk the input first", out)
}

func detectPERF82(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	if pos, ok := b2CallInLoop(facts, []string{"rows.StructScan"}); ok {
		b2Emit(unit, &MetaPERF82, pos,
			"rows.StructScan inside a for rows.Next() loop; pre-allocate the destination slice", out)
	}
}

func detectPERF83(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	if pos, ok := b2CallInLoop(facts, []string{"rows.MapScan"}); ok {
		b2Emit(unit, &MetaPERF83, pos,
			"rows.MapScan inside a for rows.Next() loop on a hot path; switch to StructScan with a typed destination", out)
	}
}

func detectPERF84(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2IsRequestPath(src, facts) {
		return
	}
	trig := []string{"db.Beginx(", "db.MustBegin(", "tx, err := db.Beginx", "tx := db.MustBegin"}
	if !b2HasAny(src, trig) {
		return
	}
	b2Emit(unit, &MetaPERF84, b2FirstPos(src, trig),
		"sqlx transaction opened inside handler; collapse to a single statement or shorter transaction", out)
}

func detectPERF85(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	if pos, ok := b2CallInLoop(facts, []string{"sqlx.Named", "sqlx.In"}); ok {
		b2Emit(unit, &MetaPERF85, pos,
			"sqlx.Named / sqlx.In inside a loop with a stable query shape; precompile the query once", out)
	}
}

func detectPERF86(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "echo.Context") {
		return
	}
	needle := ""
	if strings.Contains(src, "c.JSON(") {
		needle = "c.JSON("
	} else if strings.Contains(src, "c.JSONP(") {
		needle = "c.JSONP("
	} else {
		return
	}
	if strings.Contains(src, "json.NewEncoder(") || strings.Contains(src, "c.Stream(") {
		return
	}
	b2Emit(unit, &MetaPERF86, strings.Index(src, needle),
		"Echo c.JSON allocates an encoder per response; pool the encoder or stream with json.NewEncoder", out)
}

func detectPERF87(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "echo.Context") {
		return
	}
	needle := ""
	if strings.Contains(src, "c.BindWith(") {
		needle = "c.BindWith("
	} else if strings.Contains(src, "c.Bind(") {
		needle = "c.Bind("
	} else {
		return
	}
	if b2HasAny(src, []string{"echo.Binder", "NewBinder()", "DefaultBinder{}"}) {
		return
	}
	b2Emit(unit, &MetaPERF87, strings.Index(src, needle),
		"Echo default binder runs full validation per request; use a custom binder for trusted paths", out)
}

func detectPERF88(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	trig := []string{"e.Static(", "e.File(", "echo.Static(", "Static("}
	if !b2HasAny(src, trig) {
		return
	}
	if b2HasAny(src, []string{"Cache-Control", "cacheControl", "SetCache", "MaxAge"}) {
		return
	}
	b2Emit(unit, &MetaPERF88, b2FirstPos(src, trig),
		"Static handler is missing cache headers; set explicit Cache-Control for static assets", out)
}

func detectPERF89(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "echo.HandlerFunc") {
		return
	}
	trig := []string{"make([]", "make(map[", "json.Unmarshal(", "&MyConfig{}"}
	if !b2HasAny(src, trig) {
		return
	}
	if strings.Contains(src, "sync.Once") || strings.Contains(src, "var once ") {
		return
	}
	b2Emit(unit, &MetaPERF89, b2FirstPos(src, trig),
		"Echo middleware allocates per request; move construction to package scope or sync.Once", out)
}

func detectPERF90(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "echo.HandlerFunc") {
		return
	}
	if !strings.Contains(src, "c.Set(") {
		return
	}
	if b2HasAny(src, []string{`c.Set("user_id",`, `c.Set("request_id",`, `c.Set("trace_id",`}) {
		return
	}
	b2Emit(unit, &MetaPERF90, strings.Index(src, "c.Set("),
		"c.Set in Echo middleware stores a value; prefer small scalars (ids, request ids) over large blobs", out)
}

var b2FiberMarkers = []string{"*fiber.Ctx", "fiber.Ctx", "fiber.App", "fiber.New(", "fiber.Config", "fiber.Handler"}

func detectPERF91(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2HasAny(src, b2FiberMarkers) {
		return
	}
	if strings.Contains(src, "sync.Pool") || strings.Contains(src, "bytePool") {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "c.Request.Body", "c.Request.BodyStream", "c.Response.BodyWriter", "bytes.NewReader":
			b2Emit(unit, &MetaPERF91, call.StartByte,
				"Fiber handler allocates a per-request buffer without using a sync.Pool; reuse buffers across requests", out)
			return
		}
	}
}

func detectPERF92(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2HasAny(src, b2FiberMarkers) {
		return
	}
	for _, gr := range facts.GoStarts {
		text := src[gr[0]:gr[1]]
		if strings.Contains(text, "c.UserContext()") || strings.Contains(text, "c.Context()") {
			continue
		}
		if b2WordBoundaryContains(text, "c") {
			b2Emit(unit, &MetaPERF92, gr[0],
				"Fiber *fiber.Ctx is captured inside a goroutine; the ctx is reused per request and will race — use c.UserContext()", out)
			return
		}
	}
}

func detectPERF93(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2HasAny(src, b2FiberMarkers) {
		return
	}
	if strings.Contains(src, "encoderPool") || strings.Contains(src, "jsonPool") {
		return
	}
	if !strings.Contains(src, "c.JSON(") && !strings.Contains(src, "json.NewEncoder(") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee == "c.JSON" || call.Callee == "json.NewEncoder" {
			b2Emit(unit, &MetaPERF93, call.StartByte,
				"JSON response is allocated per request in a Fiber handler; reuse a pooled encoder", out)
			return
		}
	}
}

func detectPERF94(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2HasAny(src, b2FiberMarkers) {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "io.ReadAll":
			if len(call.Arguments) == 0 {
				continue
			}
			t := call.Arguments[0]
			if strings.Contains(t, "RequestBodyStream") ||
				strings.Contains(t, "BodyStream") ||
				strings.Contains(t, "c.Request.Body") {
				b2Emit(unit, &MetaPERF94, call.StartByte,
					"io.ReadAll on a Fiber body stream triggers an extra copy; use c.PostBody() for zero-copy reads", out)
				return
			}
		case "c.Body":
			b2Emit(unit, &MetaPERF94, call.StartByte,
				"c.Body() copies the request body; use c.PostBody() for zero-copy access in Fiber handlers", out)
			return
		}
	}
}

func detectPERF95(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !b2HasAny(src, []string{"fiber.New(", "fiber.App", "app.Use(", "app.Group("}) {
		return
	}
	useCount := 0
	var first *CallFact
	for i := range facts.Calls {
		if facts.Calls[i].Callee == "app.Use" {
			useCount++
			if first == nil {
				first = &facts.Calls[i]
			}
		}
	}
	if useCount < 2 || first == nil {
		return
	}
	b2Emit(unit, &MetaPERF95, first.StartByte,
		"Fiber app registers multiple app.Use middlewares; group them by route to keep the per-request chain small", out)
}

func detectPERF96(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	grpcMarkers := []string{"RecvMsg(", "SendMsg(", "grpc.ClientStream", "grpc.ServerStream", "google.golang.org/grpc"}
	if !b2HasAny(src, grpcMarkers) || !strings.Contains(src, "RecvMsg(") {
		return
	}
	if strings.Contains(src, "msg.Reset()") || strings.Contains(src, "m.Reset()") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "stream.RecvMsg" || !IsInLoop(call) || call.EnclosingLoop == nil {
			continue
		}
		loopStart := *call.EnclosingLoop
		hasAlloc := false
		for _, a := range facts.Assignments {
			if a.EnclosingLoop == nil || *a.EnclosingLoop != loopStart {
				continue
			}
			if strings.Contains(a.Expr, "New") || (strings.Contains(a.Expr, "&") && strings.Contains(a.Expr, "{")) {
				hasAlloc = true
				break
			}
		}
		if hasAlloc {
			b2Emit(unit, &MetaPERF96, call.StartByte,
				"gRPC client allocates a new message inside the Recv loop; reuse a single message struct across iterations", out)
			return
		}
	}
}

func detectPERF97(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !strings.Contains(src, "proto.Marshal") && !strings.Contains(src, "protojson.Marshal") {
		return
	}
	if b2HasAny(src, []string{"bytesPool", "bufPool", "MarshalBuffer"}) {
		return
	}
	if strings.Contains(src, "MarshalOptions{") && (strings.Contains(src, "Pool") || strings.Contains(src, "pool.Get")) {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "proto.Marshal" && call.Callee != "protojson.Marshal" {
			continue
		}
		if !IsInLoop(call) {
			continue
		}
		b2Emit(unit, &MetaPERF97, call.StartByte,
			"proto.Marshal is called inside a loop; reuse a MarshalOptions/buffer pool to avoid repeated allocations", out)
		return
	}
}

var b2RedisLoopTriggers = []string{
	"rdb.Set", "rdb.Get", "rdb.Del", "rdb.Incr", "rdb.Decr",
	"rdb.HSet", "rdb.HGet", "rdb.HDel", "rdb.LPush", "rdb.RPush",
	"rdb.LPop", "rdb.RPop", "rdb.SAdd", "rdb.SRem", "rdb.ZAdd", "rdb.ZRem", "rdb.Expire",
}

func detectPERF98(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	redisMarkers := []string{
		"redis.Client", "*redis.Client", "redis.UniversalClient",
		"github.com/redis/go-redis", "github.com/go-redis/redis",
	}
	if !b2HasAny(src, redisMarkers) {
		return
	}
	if b2HasAny(src, []string{".Pipeline()", ".Pipelined(", ".TxPipeline()", ".TxPipelined("}) {
		return
	}
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		for _, t := range b2RedisLoopTriggers {
			if call.Callee == t {
				b2Emit(unit, &MetaPERF98, call.StartByte,
					"go-redis client is called inside a loop without a pipeline; batch the calls with rdb.Pipeline() to amortise round-trips", out)
				return
			}
		}
	}
}

var b2HighCardLabels = []string{
	"user_id", "userId", "userid", "request_id", "requestId", "requestid",
	"uuid", "UUID", "trace_id", "traceId", "span_id", "spanId",
	"session_id", "sessionId", "email", "ip", "client_ip", "clientIp",
	"remote_addr", "remoteAddr", "user", "username", "account", "account_id", "accountId",
	"tenant_id", "tenantId", "order_id", "orderId", "path",
}

func detectPERF99(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	promMarkers := []string{
		"prometheus.NewCounterVec", "prometheus.NewCounter(",
		"prometheus.NewGaugeVec", "prometheus.NewGauge(",
		"prometheus.NewHistogramVec", "prometheus.NewHistogram(",
		"prometheus.NewSummaryVec", "github.com/prometheus/client_golang",
	}
	if !b2HasAny(src, promMarkers) {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "prometheus.NewCounterVec", "prometheus.NewGaugeVec",
			"prometheus.NewHistogramVec", "prometheus.NewSummaryVec":
		default:
			continue
		}
		for _, arg := range call.Arguments {
			for _, lab := range b2HighCardLabels {
				if strings.Contains(arg, lab) {
					b2Emit(unit, &MetaPERF99, call.StartByte,
						"Prometheus metric registers a high-cardinality label (user ID / UUID / path); time series storage will explode — bound the label space", out)
					return
				}
			}
		}
	}
}

func detectPERF100(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	cobraMarkers := []string{"cobra.Command{", "&cobra.Command{", "github.com/spf13/cobra"}
	if !b2HasAny(src, cobraMarkers) {
		return
	}
	flagSfx := []string{".String(", ".Bool(", ".Int(", ".Int64(", ".Duration(", ".Float64(", ".StringSlice(", ".StringArray("}
	count := 0
	first := -1
	searchFrom := 0
	for {
		rel := -1
		kind := ""
		for _, k := range []string{".Flags().", ".PersistentFlags()."} {
			if i := strings.Index(src[searchFrom:], k); i >= 0 {
				if rel < 0 || i < rel {
					rel = i
					kind = k
				}
			}
		}
		if rel < 0 {
			break
		}
		start := searchFrom + rel
		window := src[start:]
		if len(window) > 48 {
			window = window[:48]
		}
		matched := false
		for _, sfx := range flagSfx {
			if strings.Contains(window, sfx) {
				matched = true
				break
			}
		}
		if matched {
			if first < 0 {
				first = start
			}
			count++
		}
		searchFrom = start + len(kind)
	}
	if count < 4 || first < 0 {
		return
	}
	b2Emit(unit, &MetaPERF100, first,
		"cobra.Command registers many flags inline; defer heavy init to PersistentPreRunE or a sync.Once to keep CLI startup fast", out)
}

func detectPERF101(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	searchFrom := 0
	for {
		rel := strings.Index(src[searchFrom:], "http.Server{")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		end := start + 256
		if end > len(src) {
			end = len(src)
		}
		window := src[start:end]
		if strings.Contains(window, "ReadTimeout:") ||
			strings.Contains(window, "ReadHeaderTimeout:") ||
			strings.Contains(window, "WriteTimeout:") ||
			strings.Contains(window, "IdleTimeout:") {
			searchFrom = start + len("http.Server{")
			continue
		}
		b2Emit(unit, &MetaPERF101, start,
			"http.Server is missing ReadTimeout, WriteTimeout, and IdleTimeout settings", out)
		searchFrom = start + len("http.Server{")
	}
}

func detectPERF102(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !strings.Contains(src, ".WriteHeader(") {
		return
	}
	type site struct {
		start    int
		receiver string
		returns  bool
		funcKey  string
	}
	var sites []site
	for _, call := range facts.Calls {
		if !strings.HasSuffix(call.Callee, ".WriteHeader") {
			continue
		}
		recv := call.Callee[:len(call.Callee)-len(".WriteHeader")]
		name, _ := EnclosingFunctionName(src, call.StartByte)
		sites = append(sites, site{
			start:    call.StartByte,
			receiver: recv,
			returns:  b2WriteHeaderFollowedByReturn(src, call.StartByte),
			funcKey:  name + "|" + recv,
		})
	}
	if len(sites) == 0 {
		return
	}
	byKey := map[string][]site{}
	for _, s := range sites {
		byKey[s.funcKey] = append(byKey[s.funcKey], s)
	}
	for _, group := range byKey {
		if len(group) < 2 {
			continue
		}
		allReturn := true
		for i := 0; i < len(group)-1; i++ {
			if !group[i].returns {
				allReturn = false
				break
			}
		}
		if allReturn {
			continue
		}
		b2Emit(unit, &MetaPERF102, group[0].start,
			"w.WriteHeader called multiple times; only the first call takes effect", out)
		return
	}
}

func detectPERF103(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	for _, call := range facts.Calls {
		switch call.Callee {
		case "client.Do", "http.Get", "http.Post", "http.PostForm", "client.Get", "client.Post":
		default:
			continue
		}
		scope := b2EnclosingFuncBody(src, call.StartByte)
		if strings.Contains(scope, ".Body.Close()") || strings.Contains(scope, ".Body.Close(") {
			continue
		}
		b2Emit(unit, &MetaPERF103, call.StartByte,
			"http response body is not closed in this function; defer resp.Body.Close() after the call", out)
	}
}

func detectPERF105(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	for _, call := range facts.Calls {
		if call.Callee != "runtime.SetFinalizer" {
			continue
		}
		b2Emit(unit, &MetaPERF105, call.StartByte,
			"runtime.SetFinalizer adds GC overhead; prefer explicit Close/Release methods", out)
	}
}

func detectPERF106(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	// Path 1: write-heavy sync.Map usage in file.
	if strings.Contains(src, "sync.Map") {
		writes, reads := 0, 0
		for _, call := range facts.Calls {
			m := b2MethodName(call.Callee)
			switch m {
			case "Store", "Swap", "LoadAndDelete", "Delete", "CompareAndSwap", "CompareAndDelete":
				writes++
			case "Load", "LoadOrStore", "Range":
				reads++
			}
		}
		if reads > 0 && writes > reads {
			pos := strings.Index(src, "sync.Map")
			b2Emit(unit, &MetaPERF106, pos,
				"sync.Map is write-heavy; use a plain map guarded by a sync.Mutex instead", out)
			return
		}
	}
	// Path 2: package-level cache without eviction bounds (Rust parity).
	lines := strings.Split(src, "\n")
	depth := 0
	byteOff := 0
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if depth == 0 && strings.HasPrefix(t, "var ") {
			rest := strings.TrimSpace(strings.TrimPrefix(t, "var "))
			if !strings.HasPrefix(rest, "(") {
				name := ""
				for i, r := range rest {
					if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (i > 0 && r >= '0' && r <= '9') {
						name += string(r)
						continue
					}
					break
				}
				isCache := strings.Contains(t, "sync.Map") ||
					(strings.Contains(strings.ToLower(name), "cache") &&
						(strings.Contains(t, "sync.Map") || strings.Contains(t, "map[")))
				if name != "" && isCache {
					lower := strings.ToLower(src)
					hasEvict := strings.Contains(src, "len("+name+")") ||
						strings.Contains(src, "delete("+name+",") ||
						strings.Contains(lower, name+".delete(") ||
						strings.Contains(lower, name+".loadanddelete(") ||
						strings.Contains(src, "clear("+name+")")
					if !hasEvict {
						hasRead := strings.Contains(src, name+".Load")
						hasWrite := strings.Contains(src, name+".Store")
						if !strings.Contains(t, "sync.Map") {
							hasRead = strings.Contains(src, name+"[")
							hasWrite = strings.Contains(src, name+"[")
						}
						if hasRead && hasWrite {
							b2Emit(unit, &MetaPERF106, byteOff,
								"package-level cache without eviction bounds; it will grow unbounded under concurrent load", out)
							return
						}
					}
				}
			}
		}
		depth += strings.Count(line, "{") - strings.Count(line, "}")
		if depth < 0 {
			depth = 0
		}
		byteOff += len(line) + 1
	}
}


func detectPERF107(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		if call.Callee != "binary.Read" && call.Callee != "binary.Write" {
			continue
		}
		b2Emit(unit, &MetaPERF107, call.StartByte,
			"encoding/binary Read/Write inside a loop uses reflection; reuse a pre-encoded buffer or hand-roll the byte order", out)
	}
}

func detectPERF108(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	for _, call := range facts.Calls {
		if call.Callee != "sort.Search" || !IsInLoop(call) {
			continue
		}
		b2Emit(unit, &MetaPERF108, call.StartByte,
			"sort.Search inside a loop; hoist the search or use a different data structure", out)
	}
}

func detectPERF109(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	if !strings.Contains(src, "for ") {
		return
	}
	markers := []string{
		"fmt.Sprintf(", "fmt.Sprint(", "strings.Join(", "strings.ToLower(", "strings.ToUpper(",
		"strconv.Itoa(", "strconv.FormatInt(", "strconv.FormatUint(", "filepath.Join(",
	}
	for _, fr := range facts.ForRanges {
		start, end := fr[0], fr[1]
		if end > len(src) {
			end = len(src)
		}
		if end <= start {
			end = start + 1024
			if end > len(src) {
				end = len(src)
			}
		}
		rangeText := src[start:end]
		for _, marker := range markers {
			if !strings.Contains(rangeText, marker) {
				continue
			}
			if !(strings.Contains(rangeText, "[") && strings.Contains(rangeText, "]")) {
				continue
			}
			b2Emit(unit, &MetaPERF109, start,
				"expensive map-key computation inside the loop; cache or simplify the key", out)
			return
		}
		if strings.Count(rangeText, "fmt.Sprintf(") >= 2 &&
			strings.Contains(rangeText, "[") && strings.Contains(rangeText, "]") {
			b2Emit(unit, &MetaPERF109, start,
				"map key is recomputed multiple times in the loop; cache it per iteration or hoist if invariant", out)
			return
		}
	}
}

func detectPERF110(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	src := unit.Source
	if !strings.Contains(src, "sync.Pool") {
		return
	}
	searchFrom := 0
	for {
		rel := strings.Index(src[searchFrom:], "sync.Pool{")
		if rel < 0 {
			return
		}
		start := searchFrom + rel
		depth := 0
		end := -1
		for i := start; i < len(src); i++ {
			switch src[i] {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					end = i
				}
			}
			if end >= 0 {
				break
			}
		}
		if end < 0 {
			return
		}
		text := src[start : end+1]
		newIdx := strings.Index(text, "New:")
		if newIdx < 0 {
			searchFrom = end + 1
			continue
		}
		after := text[newIdx:]
		open := strings.Index(after, "func()")
		if open < 0 {
			searchFrom = end + 1
			continue
		}
		sig := after[open:]
		sigEnd := strings.Index(sig, "{")
		if sigEnd < 0 {
			searchFrom = end + 1
			continue
		}
		signature := sig[:sigEnd]
		if strings.Contains(signature, "*") {
			searchFrom = end + 1
			continue
		}
		returnType := strings.TrimSpace(strings.TrimPrefix(signature, "func()"))
		returnType = strings.TrimPrefix(returnType, "*")
		returnType = strings.TrimSpace(returnType)
		if returnType == "" || returnType == "_" {
			searchFrom = end + 1
			continue
		}
		if retIdx := strings.Index(after, "return"); retIdx >= 0 {
			afterRet := strings.TrimSpace(after[retIdx+len("return"):])
			if strings.HasPrefix(afterRet, "&") || strings.HasPrefix(afterRet, "new(") {
				searchFrom = end + 1
				continue
			}
		}
		b2Emit(unit, &MetaPERF110, start,
			"sync.Pool New returns a value type; return a pointer (e.g. *Foo) to avoid boxing on Put", out)
		return
	}
}

func detectPERF111(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	src := unit.Source
	for _, fr := range facts.ForRanges {
		start, end := fr[0], fr[1]
		if end > len(src) {
			end = len(src)
		}
		if end <= start {
			continue
		}
		rangeText := src[start:end]
		if !strings.Contains(rangeText, "range []rune(") {
			continue
		}
		b2Emit(unit, &MetaPERF111, start,
			"range over []rune(s) allocates a rune slice; range over the string directly", out)
	}
}

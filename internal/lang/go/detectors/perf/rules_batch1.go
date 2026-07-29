package perf

import (
	"strconv"
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("PERF-9", detectPERF9, &MetaPERF9)
	RegisterRule("PERF-10", detectPERF10, &MetaPERF10)
	RegisterRule("PERF-11", detectPERF11, &MetaPERF11)
	RegisterRule("PERF-12", detectPERF12, &MetaPERF12)
	RegisterRule("PERF-13", detectPERF13, &MetaPERF13)
	RegisterRule("PERF-14", detectPERF14, &MetaPERF14)
	RegisterRule("PERF-15", detectPERF15, &MetaPERF15)
	RegisterRule("PERF-16", detectPERF16, &MetaPERF16)
	RegisterRule("PERF-17", detectPERF17, &MetaPERF17)
	RegisterRule("PERF-18", detectPERF18, &MetaPERF18)
	RegisterRule("PERF-19", detectPERF19, &MetaPERF19)
	RegisterRule("PERF-20", detectPERF20, &MetaPERF20)
	RegisterRule("PERF-21", detectPERF21, &MetaPERF21)
	RegisterRule("PERF-22", detectPERF22, &MetaPERF22)
	RegisterRule("PERF-23", detectPERF23, &MetaPERF23)
	RegisterRule("PERF-24", detectPERF24, &MetaPERF24)
	RegisterRule("PERF-25", detectPERF25, &MetaPERF25)
	RegisterRule("PERF-26", detectPERF26, &MetaPERF26)
	RegisterRule("PERF-27", detectPERF27, &MetaPERF27)
	RegisterRule("PERF-28", detectPERF28, &MetaPERF28)
	RegisterRule("PERF-29", detectPERF29, &MetaPERF29)
	RegisterRule("PERF-30", detectPERF30, &MetaPERF30)
	RegisterRule("PERF-31", detectPERF31, &MetaPERF31)
	RegisterRule("PERF-33", detectPERF33, &MetaPERF33)
	RegisterRule("PERF-34", detectPERF34, &MetaPERF34)
	RegisterRule("PERF-35", detectPERF35, &MetaPERF35)
	RegisterRule("PERF-36", detectPERF36, &MetaPERF36)
	RegisterRule("PERF-37", detectPERF37, &MetaPERF37)
	RegisterRule("PERF-38", detectPERF38, &MetaPERF38)
	RegisterRule("PERF-39", detectPERF39, &MetaPERF39)
	RegisterRule("PERF-40", detectPERF40, &MetaPERF40)
	RegisterRule("PERF-41", detectPERF41, &MetaPERF41)
	RegisterRule("PERF-42", detectPERF42, &MetaPERF42)
	RegisterRule("PERF-43", detectPERF43, &MetaPERF43)
	RegisterRule("PERF-44", detectPERF44, &MetaPERF44)
	RegisterRule("PERF-45", detectPERF45, &MetaPERF45)
	RegisterRule("PERF-46", detectPERF46, &MetaPERF46)
	RegisterRule("PERF-47", detectPERF47, &MetaPERF47)
	RegisterRule("PERF-48", detectPERF48, &MetaPERF48)
	RegisterRule("PERF-49", detectPERF49, &MetaPERF49)
	RegisterRule("PERF-51", detectPERF51, &MetaPERF51)
	RegisterRule("PERF-52", detectPERF52, &MetaPERF52)
	RegisterRule("PERF-53", detectPERF53, &MetaPERF53)
	RegisterRule("PERF-54", detectPERF54, &MetaPERF54)
	RegisterRule("PERF-55", detectPERF55, &MetaPERF55)
	RegisterRule("PERF-56", detectPERF56, &MetaPERF56)
	RegisterRule("PERF-57", detectPERF57, &MetaPERF57)
	RegisterRule("PERF-58", detectPERF58, &MetaPERF58)
	RegisterRule("PERF-59", detectPERF59, &MetaPERF59)
	RegisterRule("PERF-60", detectPERF60, &MetaPERF60)
}

// batch1RequestPath expands IsRequestPath with source-level tokens when the
// seed SourceIndex needle table is incomplete for batch-1 rules.
func batch1RequestPath(unit *core.ParsedUnit, facts *GoPerfFacts) bool {
	if IsRequestPath(facts.SourceIndex) {
		return true
	}
	s := unit.Source
	return strings.Contains(s, "c.JSON(") ||
		strings.Contains(s, "c.String(") ||
		strings.Contains(s, "c.HTML(") ||
		strings.Contains(s, "c.Bind(") ||
		strings.Contains(s, "c.ShouldBind") ||
		strings.Contains(s, "func ServeHTTP") ||
		strings.Contains(s, "gin.HandlerFunc") ||
		strings.Contains(s, "*gin.Context") ||
		strings.Contains(s, "http.ResponseWriter") ||
		strings.Contains(s, "http.HandlerFunc") ||
		strings.Contains(s, "echo.Context") ||
		strings.Contains(s, "*fiber.Ctx")
}

func batch1InLoopByte(facts *GoPerfFacts, startByte int) bool {
	for _, fr := range facts.ForRanges {
		if fr[0] <= startByte && startByte <= fr[1] {
			return true
		}
	}
	return false
}

func batch1FirstPos(source string, needles []string) int {
	best := -1
	for _, n := range needles {
		if i := strings.Index(source, n); i >= 0 && (best < 0 || i < best) {
			best = i
		}
	}
	if best < 0 {
		return 0
	}
	return best
}

// --- PERF-9..16 parsing / loop allocs ---

func detectPERF9(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		if call.Callee != "url.Parse" && call.Callee != "url.ParseRequestURI" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF9, file, line, col, "URL is parsed inside a loop body", out)
	}
}

func detectPERF10(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	requestPath := batch1RequestPath(unit, facts) ||
		strings.Contains(source, "c.HTML(") ||
		strings.Contains(source, "c *gin.Context") ||
		strings.Contains(source, "c echo.Context")
	if !requestPath {
		return
	}
	triggers := []string{
		"template.New(",
		"template.ParseFiles(",
		"template.Must(template.Parse",
		"html/template.New(",
		"html/template.ParseFiles(",
	}
	found := false
	for _, t := range triggers {
		if strings.Contains(source, t) {
			found = true
			break
		}
	}
	if !found {
		return
	}
	if strings.Contains(source, "template.Must(parseTemplates(") ||
		strings.Contains(source, "var indexTmpl =") ||
		strings.Contains(source, "sync.Once") {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "template.New", "template.ParseFiles", "html/template.New", "html/template.ParseFiles":
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(&MetaPERF10, file, line, col, "template is parsed on the request path", out)
			return
		}
	}
	start := batch1FirstPos(source, triggers)
	line, col := unit.LineCol(start)
	rules.PushFinding(&MetaPERF10, file, line, col, "template is parsed on the request path", out)
}

func detectPERF11(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	requestPath := batch1RequestPath(unit, facts) || strings.Contains(source, "func (")
	if !requestPath {
		return
	}
	if !strings.Contains(source, "http.Client{") && !strings.Contains(source, "&http.Client{") {
		return
	}
	if strings.Contains(source, "var defaultClient =") ||
		strings.Contains(source, "var httpClient =") ||
		strings.Contains(source, "sync.Once") {
		return
	}
	for _, a := range facts.Assignments {
		if !strings.Contains(a.Expr, "http.Client{") {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(&MetaPERF11, file, line, col, "http.Client is allocated on the request path", out)
		return
	}
}

func detectPERF12(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	requestPath := batch1RequestPath(unit, facts) || strings.Contains(source, "func (")
	if !requestPath {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "db.Prepare" && call.Callee != "db.PrepareContext" {
			continue
		}
		if !IsInLoop(call) &&
			(strings.Contains(source, "sync.Once") ||
				strings.Contains(source, "var stmtOnce") ||
				strings.Contains(source, "StmtOnce")) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF12, file, line, col, "prepared statement is created on the request path", out)
		return
	}
}

func detectPERF13(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if strings.Contains(source, "time.NewTicker(") || strings.Contains(source, "time.NewTimer(") {
		return
	}
	for _, call := range facts.Calls {
		if !IsInLoop(call) || call.Callee != "time.After" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF13, file, line, col, "time.After is allocated inside a loop body", out)
	}
}

func detectPERF14(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		switch call.Callee {
		case "filepath.Glob", "os.ReadDir", "ioutil.ReadDir":
		default:
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF14, file, line, col, "directory scan is performed inside a loop body", out)
	}
}

func detectPERF15(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		switch call.Callee {
		case "strconv.Itoa", "strconv.FormatInt", "strconv.FormatUint", "strconv.FormatFloat":
		default:
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF15, file, line, col, "strconv formatting is performed inside a loop body", out)
	}
}

func detectPERF16(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "bytes.Buffer{") && !strings.Contains(source, "new(bytes.Buffer)") {
		return
	}
	// Prefer assignment facts (buf := bytes.Buffer{}).
	for _, a := range facts.Assignments {
		if !IsAssignmentInLoop(a) {
			continue
		}
		if !strings.Contains(a.Expr, "bytes.Buffer{") && !strings.Contains(a.Expr, "new(bytes.Buffer)") {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(&MetaPERF16, file, line, col, "bytes.Buffer is allocated inside a loop body", out)
		return
	}
	// Fallback: bare composite literal sites via source + for-ranges.
	for _, needle := range []string{"bytes.Buffer{}", "new(bytes.Buffer)"} {
		search := 0
		for {
			rel := strings.Index(source[search:], needle)
			if rel < 0 {
				break
			}
			pos := search + rel
			if batch1InLoopByte(facts, pos) {
				line, col := unit.LineCol(pos)
				rules.PushFinding(&MetaPERF16, file, line, col, "bytes.Buffer is allocated inside a loop body", out)
				return
			}
			search = pos + len(needle)
		}
	}
}

// --- PERF-17..25 request path / crypto ---

func detectPERF17(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !batch1RequestPath(unit, facts) {
		return
	}
	for _, a := range facts.Assignments {
		if !IsAssignmentInLoop(a) {
			continue
		}
		if !strings.Contains(a.Expr, "strings.Join(") {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(&MetaPERF17, file, line, col, "strings.Join is invoked inside a loop on a request path", out)
	}
}

func detectPERF18(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, a := range facts.Assignments {
		dst, src, ok := parseAppendSpread(a.Expr)
		if !ok {
			continue
		}
		if !strings.Contains(source, "0, len("+src+")") {
			continue
		}
		if !strings.Contains(source, "append("+dst+", "+src+"...)") {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(
			&MetaPERF18, file, line, col,
			"large slice is copied via append(dst, src...) where ownership transfer or reslicing may suffice", out,
		)
		return
	}
}

func parseAppendSpread(expr string) (dst, src string, ok bool) {
	rest := strings.TrimSpace(expr)
	if !strings.HasPrefix(rest, "append(") || !strings.HasSuffix(rest, ")") {
		return "", "", false
	}
	rest = strings.TrimPrefix(rest, "append(")
	rest = strings.TrimSuffix(rest, ")")
	comma := strings.Index(rest, ",")
	if comma < 0 {
		return "", "", false
	}
	dst = strings.TrimSpace(rest[:comma])
	after := strings.TrimSpace(rest[comma+1:])
	if !strings.HasSuffix(after, "...") {
		return "", "", false
	}
	src = strings.TrimSpace(strings.TrimSuffix(after, "..."))
	if dst == "" || src == "" || !isSimpleIdent(dst) || !isSimpleIdent(src) {
		return "", "", false
	}
	return dst, src, true
}

func detectPERF19(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "for _, record := range records") {
		return
	}
	if !strings.Contains(source, "processRecord(record)") {
		return
	}
	if strings.Contains(source, "for _, record := range &records") ||
		strings.Contains(source, "for _, record := range recordsPtr") {
		return
	}
	start := strings.Index(source, "for _, record := range records")
	if start < 0 {
		start = 0
	}
	line, col := unit.LineCol(start)
	rules.PushFinding(
		&MetaPERF19, file, line, col,
		"range over a slice of large structs copies each element by value", out,
	)
}

func detectPERF20(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	if !strings.Contains(source, "reflect.ValueOf(") &&
		!strings.Contains(source, "reflect.TypeOf(") &&
		!strings.Contains(source, "reflect.New(") {
		return
	}
	if strings.Contains(source, "// reflection initialised at startup") {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "reflect.ValueOf", "reflect.TypeOf", "reflect.New":
		default:
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF20, file, line, col,
			"reflect is invoked on a request path; cache reflect.Type or Value at startup", out,
		)
		return
	}
}

func detectPERF21(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !batch1RequestPath(unit, facts) {
		return
	}
	if !strings.Contains(unit.Source, "io.ReadAll(") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "io.ReadAll" || len(call.Arguments) == 0 {
			continue
		}
		arg := call.Arguments[0]
		if strings.Contains(arg, "c.Request.Body") ||
			strings.Contains(arg, "r.Body") ||
			strings.Contains(arg, "req.Body") ||
			strings.Contains(arg, "ctx.Request.Body") {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(
				&MetaPERF21, file, line, col,
				"io.ReadAll fully buffers a request body on a request path", out,
			)
			return
		}
	}
}

func detectPERF22(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	if !strings.Contains(source, "os.ReadFile(") && !strings.Contains(source, "ioutil.ReadFile(") {
		return
	}
	if strings.Contains(source, "sync.Once") ||
		strings.Contains(source, "loadOnce") ||
		strings.Contains(source, "readOnce") ||
		strings.Contains(source, "fileOnce") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "os.ReadFile" && call.Callee != "ioutil.ReadFile" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF22, file, line, col,
			"os.ReadFile is invoked on a request path; load the file once at startup", out,
		)
		return
	}
}

func detectPERF23(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !batch1RequestPath(unit, facts) {
		return
	}
	for _, a := range facts.Assignments {
		if !strings.Contains(a.Text, "bytes.NewReader(") {
			continue
		}
		if !IsAssignmentInLoop(a) && !strings.Contains(a.Text, ":=") {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(
			&MetaPERF23, file, line, col,
			"bytes.NewReader is allocated per request; reuse a pooled buffer instead", out,
		)
		return
	}
}

func detectPERF24(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	triggers := map[string]struct{}{
		"sha256.New": {}, "sha1.New": {}, "md5.New": {},
		"hmac.New": {}, "blake2b.New256": {}, "blake2s.New256": {},
	}
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		if _, ok := triggers[call.Callee]; !ok {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF24, file, line, col, "crypto hasher is allocated inside a loop body", out)
	}
}

func detectPERF25(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	indexTriggers := []string{
		"rsa.GenerateKey(", "rsa.GenerateMultiPrimeKey(",
		"ecdsa.GenerateKey(", "ed25519.GenerateKey(",
	}
	found := false
	for _, t := range indexTriggers {
		if strings.Contains(source, t) {
			found = true
			break
		}
	}
	if !found {
		return
	}
	if strings.Contains(source, "var (") &&
		(strings.Contains(source, "// gen once") || strings.Contains(source, "sync.Once")) {
		return
	}
	callees := map[string]struct{}{
		"rsa.GenerateKey": {}, "rsa.GenerateMultiPrimeKey": {},
		"ecdsa.GenerateKey": {}, "ed25519.GenerateKey": {},
	}
	onRequest := batch1RequestPath(unit, facts)
	inLoop := false
	for _, c := range facts.Calls {
		if IsInLoop(c) {
			if _, ok := callees[c.Callee]; ok {
				inLoop = true
				break
			}
		}
	}
	if !onRequest && !inLoop {
		return
	}
	for _, call := range facts.Calls {
		if _, ok := callees[call.Callee]; !ok {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF25, file, line, col,
			"asymmetric key pair is generated on a request path or in a loop", out,
		)
		return
	}
}

// --- PERF-26..45 general perf ---

func detectPERF26(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		switch call.Callee {
		case "base64.StdEncoding.EncodeToString", "base64.StdEncoding.DecodeString",
			"base64.URLEncoding.EncodeToString", "base64.URLEncoding.DecodeString",
			"base64.RawStdEncoding.EncodeToString", "base64.RawStdEncoding.DecodeString",
			"base64.NewEncoder", "base64.NewDecoder":
		default:
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF26, file, line, col,
			"base64 encoding or decoding is performed inside a loop body", out,
		)
	}
}

func detectPERF27(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if strings.Contains(source, "sync.Pool") {
		return
	}
	for _, a := range facts.Assignments {
		expr := a.Expr
		isBuffer := strings.Contains(expr, "bytes.Buffer{") ||
			strings.Contains(expr, "new(bytes.Buffer)") ||
			strings.Contains(expr, "strings.Builder{")
		largeMake := largeMakeByteSlice(expr)
		if !isBuffer && !largeMake {
			continue
		}
		inLoop := IsAssignmentInLoop(a)
		if largeMake && !inLoop {
			continue
		}
		if !IsHotPath(source, a.StartByte, inLoop) {
			continue
		}
		msg := "bytes.Buffer is allocated on a hot path; pool it via sync.Pool"
		if strings.Contains(expr, "strings.Builder") {
			msg = "strings.Builder is allocated on a hot path; pool it via sync.Pool or hoist + Reset"
		} else if largeMake {
			msg = "large []byte is make'd inside a loop; pool and reuse or hoist the buffer"
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(&MetaPERF27, file, line, col, msg, out)
		return
	}
}

func largeMakeByteSlice(expr string) bool {
	expr = strings.TrimSpace(expr)
	var rest string
	switch {
	case strings.HasPrefix(expr, "make([]byte,"):
		rest = strings.TrimPrefix(expr, "make([]byte,")
	case strings.HasPrefix(expr, "make([]uint8,"):
		rest = strings.TrimPrefix(expr, "make([]uint8,")
	default:
		return false
	}
	rest = strings.TrimSpace(rest)
	var nums []uint64
	cur := ""
	for _, r := range rest {
		if r >= '0' && r <= '9' {
			cur += string(r)
			continue
		}
		if cur != "" {
			if n, err := strconv.ParseUint(cur, 10, 64); err == nil {
				nums = append(nums, n)
			}
			cur = ""
		}
	}
	if cur != "" {
		if n, err := strconv.ParseUint(cur, 10, 64); err == nil {
			nums = append(nums, n)
		}
	}
	for _, n := range nums {
		if n >= 4096 {
			return true
		}
	}
	return false
}

func detectPERF28(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	if !strings.Contains(source, "sync.Mutex") && !strings.Contains(source, "sync.RWMutex") {
		return
	}
	if strings.Contains(source, "var mu sync.Mutex\n") ||
		strings.Contains(source, "var mu sync.Mutex =") ||
		strings.Contains(source, "var (\n") ||
		strings.Contains(source, "var rwMu sync.RWMutex\n") {
		return
	}
	inStruct := strings.Contains(source, "struct {") &&
		(strings.Contains(source, "\tmu sync.Mutex") ||
			strings.Contains(source, "mu sync.Mutex\n") ||
			strings.Contains(source, "rwMu sync.RWMutex"))
	literal := strings.Contains(source, "sync.Mutex{") || strings.Contains(source, "sync.RWMutex{")
	if !inStruct && !literal {
		return
	}
	start := strings.Index(source, "sync.Mutex")
	if start < 0 {
		start = strings.Index(source, "sync.RWMutex")
	}
	if start < 0 {
		start = 0
	}
	line, col := unit.LineCol(start)
	rules.PushFinding(
		&MetaPERF28, file, line, col,
		"sync.Mutex is allocated per request or per record; share a single mutex or use atomics", out,
	)
}

func detectPERF29(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if strings.Contains(source, "errgroup.WithContext") ||
		strings.Contains(source, "sem := make(chan struct{}") ||
		strings.Contains(source, "sem <- struct{}{}") ||
		strings.Contains(source, "workerCount") ||
		strings.Contains(source, "workerPool") ||
		strings.Contains(source, "semaphore") ||
		strings.Contains(source, "sync.WaitGroup") ||
		strings.Contains(source, "wg.Add(") ||
		strings.Contains(source, "c.Request.Context()") ||
		strings.Contains(source, "ctx, cancel := context.WithCancel") ||
		strings.Contains(source, "ctx, cancel := context.WithTimeout") {
		return
	}
	for _, gs := range facts.GoStarts {
		startByte, endByte := gs[0], gs[1]
		if endByte > len(source) {
			endByte = len(source)
		}
		if startByte < 0 || startByte >= len(source) {
			continue
		}
		text := source[startByte:endByte]
		if !strings.Contains(text, "go func") {
			continue
		}
		inLoop := batch1InLoopByte(facts, startByte)
		onRequest := batch1RequestPath(unit, facts)
		if !inLoop && !onRequest {
			continue
		}
		line, col := unit.LineCol(startByte)
		rules.PushFinding(
			&MetaPERF29, file, line, col,
			"goroutine is spawned without a bounded worker pool or semaphore", out,
		)
	}
}

func detectPERF30(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !batch1RequestPath(unit, facts) {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "context.Background" && call.Callee != "context.TODO" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF30, file, line, col,
			"context.Background / TODO detaches the goroutine from the request context", out,
		)
		return
	}
}

func detectPERF31(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	hasResourceDefer := strings.Contains(source, ".Close()") ||
		strings.Contains(source, "cancel()") ||
		strings.Contains(source, ".Stop()") ||
		(strings.Contains(source, "defer") &&
			(strings.Contains(source, ".Put(") ||
				strings.Contains(source, ".Unlock()") ||
				strings.Contains(source, ".RUnlock()") ||
				strings.Contains(source, "bufPool") ||
				strings.Contains(source, "sync.Pool")))
	if hasResourceDefer {
		return
	}
	for _, dr := range facts.DeferStarts {
		startByte := dr[0]
		end := startByte + 80
		if end > len(source) {
			end = len(source)
		}
		if startByte < 0 || startByte >= len(source) {
			continue
		}
		snip := source[startByte:end]
		if strings.Contains(snip, ".Put(") ||
			strings.Contains(snip, ".Close(") ||
			strings.Contains(snip, "cancel(") ||
			strings.Contains(snip, ".Unlock(") ||
			strings.Contains(snip, ".RUnlock(") ||
			strings.Contains(snip, ".Stop(") {
			continue
		}
		line, col := unit.LineCol(startByte)
		rules.PushFinding(
			&MetaPERF31, file, line, col,
			"defer is used in a hot handler function; consider explicit cleanup", out,
		)
	}
}

func detectPERF33(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	if !strings.Contains(source, "for _, item := range items") {
		return
	}
	if strings.Contains(source, "for i := 0; i < len(items);") || strings.Contains(source, "break") {
		return
	}
	start := strings.Index(source, "for _, item := range items")
	if start < 0 {
		start = 0
	}
	line, col := unit.LineCol(start)
	rules.PushFinding(
		&MetaPERF33, file, line, col,
		"range over a large slice on a request path; consider indexed scan or early break", out,
	)
}

func detectPERF34(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if strings.Contains(source, "make([]") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "append" || !IsInLoop(call) {
			continue
		}
		if !strings.Contains(source, "for _, v := range m") &&
			!strings.Contains(source, "for k, v := range m") &&
			!strings.Contains(source, "for k := range m") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF34, file, line, col,
			"append inside a for-range over a map grows the slice without preallocation", out,
		)
		return
	}
}

func detectPERF35(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) && !strings.Contains(source, "for ") {
		return
	}
	inLoopPresent := false
	for _, c := range facts.Calls {
		if IsInLoop(c) {
			inLoopPresent = true
			break
		}
	}
	for _, call := range facts.Calls {
		if call.Callee != "fmt.Sprintf" && call.Callee != "fmt.Errorf" {
			continue
		}
		if !inLoopPresent && !batch1RequestPath(unit, facts) {
			continue
		}
		if len(call.Arguments) < 2 {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF35, file, line, col,
			"fmt.Sprintf / Errorf boxes arguments through interface{} on a hot path", out,
		)
		return
	}
}

func detectPERF36(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "go func()") {
		return
	}
	if strings.Contains(source, "v := v") || strings.Contains(source, "go 1.22") {
		return
	}
	for _, gs := range facts.GoStarts {
		startByte := gs[0]
		if !batch1InLoopByte(facts, startByte) {
			continue
		}
		line, col := unit.LineCol(startByte)
		rules.PushFinding(
			&MetaPERF36, file, line, col,
			"goroutine captures a loop variable by reference; copy it per iteration", out,
		)
		// fire once per match site like Rust walk
	}
}

func detectPERF37(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	hasUnpreallocated := strings.Contains(source, "var out []int") ||
		strings.Contains(source, "out := []int{}") ||
		strings.Contains(source, "results := []int{}") ||
		strings.Contains(source, "var results []int") ||
		strings.Contains(source, "var out []string") ||
		strings.Contains(source, "out := []string{}") ||
		strings.Contains(source, "out := []byte{}") ||
		strings.Contains(source, "var out []byte")
	if !hasUnpreallocated {
		return
	}
	if strings.Contains(source, "make([]") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "append" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF37, file, line, col,
			"slice is grown by append on a request path without a capacity hint", out,
		)
		return
	}
}

func detectPERF38(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "make(chan") {
		return
	}
	if strings.Contains(source, "_test.go") {
		return
	}
	search := 0
	for {
		rel := strings.Index(source[search:], "make(chan")
		if rel < 0 {
			return
		}
		start := search + rel
		after := strings.TrimLeft(source[start+len("make(chan"):], " \t\n")
		// Buffered: make(chan T, N)
		if commaRel := strings.Index(after, ","); commaRel >= 0 {
			closeRel := strings.Index(after, ")")
			if closeRel > commaRel {
				between := strings.TrimSpace(after[:commaRel])
				if between != "" {
					search = start + len("make(chan")
					continue
				}
			}
		}
		// Unbuffered empty-struct signal: make(chan struct{})
		if strings.HasPrefix(after, "struct{}") {
			rest := strings.TrimLeft(after[len("struct{}"):], " \t\n")
			if strings.HasPrefix(rest, ")") {
				search = start + len("make(chan")
				continue
			}
		}
		line, col := unit.LineCol(start)
		rules.PushFinding(
			&MetaPERF38, file, line, col,
			"unbuffered channel blocks producers; consider a buffered channel or a worker pool", out,
		)
		return
	}
}

func detectPERF39(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "default:") {
		return
	}
	if strings.Contains(source, "time.Sleep(") {
		return
	}
	if strings.Contains(source, "!timer.Stop()") || strings.Contains(source, "if !t.Stop()") {
		return
	}
	for _, fr := range facts.ForRanges {
		start, end := fr[0], fr[1]
		if start < 0 || end > len(source) || start >= end {
			continue
		}
		text := source[start:end]
		if !strings.Contains(text, "select") || !strings.Contains(text, "default:") {
			continue
		}
		header := strings.TrimLeft(strings.SplitN(text, "\n", 2)[0], " \t")
		if !strings.HasPrefix(header, "for {") {
			continue
		}
		line, col := unit.LineCol(start)
		rules.PushFinding(
			&MetaPERF39, file, line, col,
			"select with default branch inside a for-loop is a busy-wait; add a backoff or use channels", out,
		)
	}
}

func detectPERF40(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	// Group by enclosing function name (approximation of body range).
	byFunc := map[string][]int{}
	for _, call := range facts.Calls {
		if call.Callee != "time.Now" {
			continue
		}
		name, ok := EnclosingFunctionName(source, call.StartByte)
		if !ok {
			name = ""
		}
		byFunc[name] = append(byFunc[name], call.StartByte)
	}
	for _, sites := range byFunc {
		if len(sites) < 2 {
			continue
		}
		line, col := unit.LineCol(sites[0])
		rules.PushFinding(
			&MetaPERF40, file, line, col,
			"time.Now is called repeatedly in the same function body", out,
		)
		return
	}
}

func detectPERF41(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !batch1RequestPath(unit, facts) {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "log.Println", "log.Printf", "log.Print", "log.Fatal", "log.Fatalf", "log.Fatalln":
		default:
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF41, file, line, col,
			"standard log is used in a request path; prefer a structured leveled logger", out,
		)
		return
	}
}

func detectPERF42(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	inLoopPresent := false
	for _, c := range facts.Calls {
		if IsInLoop(c) {
			inLoopPresent = true
			break
		}
	}
	if !batch1RequestPath(unit, facts) && !inLoopPresent {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "fmt.Errorf" || len(call.Arguments) == 0 {
			continue
		}
		first := call.Arguments[0]
		if !strings.HasPrefix(first, `"`) || !strings.HasSuffix(first, `"`) || len(first) < 2 {
			continue
		}
		literal := first[1 : len(first)-1]
		if strings.Contains(literal, "%") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF42, file, line, col,
			"fmt.Errorf with a static string allocates a Sprintf; use errors.New instead", out,
		)
		return
	}
}

func detectPERF43(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	onRequest := batch1RequestPath(unit, facts)
	if !onRequest && !strings.Contains(source, "for ") {
		return
	}
	for _, dr := range facts.DeferStarts {
		start, end := dr[0], dr[1]
		if start < 0 || end > len(source) || start >= end {
			continue
		}
		text := source[start:end]
		if !strings.Contains(text, "recover()") {
			continue
		}
		if !onRequest && !batch1InLoopByte(facts, start) {
			continue
		}
		line, col := unit.LineCol(start)
		rules.PushFinding(
			&MetaPERF43, file, line, col,
			"defer-recover runs in a hot path; add the recover at a higher boundary", out,
		)
	}
}

func detectPERF44(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) && !strings.Contains(source, "for ") {
		return
	}
	type assertSite struct {
		start int
		lhs   string
		fname string
	}
	var sites []assertSite
	for _, ta := range facts.TypeAssertions {
		start, end := ta[0], ta[1]
		if start < 0 || end > len(source) || start >= end {
			continue
		}
		text := source[start:end]
		lhsPart, _, ok := strings.Cut(text, ".(")
		if !ok {
			continue
		}
		lhs := strings.TrimSpace(lhsPart)
		fname, _ := EnclosingFunctionName(source, start)
		sites = append(sites, assertSite{start: start, lhs: lhs, fname: fname})
	}
	for i := 0; i+1 < len(sites); i++ {
		a, b := sites[i], sites[i+1]
		if a.lhs != "" && a.lhs == b.lhs && a.fname != "" && a.fname == b.fname {
			line, col := unit.LineCol(a.start)
			rules.PushFinding(
				&MetaPERF44, file, line, col,
				"the same type assertion is repeated; cache the result in a local variable", out,
			)
			return
		}
	}
}

func detectPERF45(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "for _, v := range") && !strings.Contains(source, "for i := 0;") {
		return
	}
	if strings.Contains(source, "make([]") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "append" || !IsInLoop(call) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF45, file, line, col,
			"append inside a loop without a capacity hint causes repeated reallocation", out,
		)
		return
	}
}

func detectPERF46(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	if !batch1RequestPath(unit, facts) {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "strings.TrimSpace", "strings.Trim", "strings.TrimPrefix",
			"strings.TrimSuffix", "strings.TrimLeft", "strings.TrimRight":
		default:
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF46, file, line, col,
			"string trimming allocates on a request path; check the need first", out,
		)
		return
	}
}

func detectPERF47(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		switch call.Callee {
		case "strings.Split", "strings.SplitN", "strings.SplitAfter":
		default:
			continue
		}
		if !IsInLoop(call) {
			continue
		}
		// Suppress for _, x := range strings.Split(...) — range iterable.
		if isRangeIterableApprox(source, call.StartByte) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF47, file, line, col,
			"strings.Split allocates a slice; consider a streaming scanner", out,
		)
		return
	}
}

// isRangeIterableApprox is true when the call sits on a for-range header line.
func isRangeIterableApprox(source string, startByte int) bool {
	if startByte < 0 {
		startByte = 0
	}
	if startByte > len(source) {
		startByte = len(source)
	}
	lineStart := strings.LastIndex(source[:startByte], "\n") + 1
	line := source[lineStart:startByte]
	return strings.Contains(line, "range ")
}

func detectPERF48(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) && !strings.Contains(source, "for ") {
		return
	}
	for _, call := range facts.Calls {
		switch call.Callee {
		case "bytes.Equal", "strings.EqualFold", "bytes.Compare":
		default:
			continue
		}
		if strings.Contains(source, "if len(a) != len(b) { return false }") ||
			strings.Contains(source, "len(prefix)") {
			return
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF48, file, line, col,
			"byte / string equality on a hot path; add a length or prefix precheck", out,
		)
		return
	}
}

func detectPERF49(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) && !strings.Contains(source, "for ") {
		return
	}
	if !strings.Contains(source, "copy(buf, payload)") && !strings.Contains(source, "copy(dst, src)") {
		return
	}
	if strings.Contains(source, "if len(payload) > len(buf) { return }") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "copy" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF49, file, line, col,
			"copy(dst, src) is invoked without explicit length validation", out,
		)
		return
	}
}

// --- PERF-51..60 gin / runtime ---

func detectPERF51(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) || !strings.Contains(source, "unsafe.Pointer") {
		return
	}
	if strings.Contains(source, "// benchmark justifies unsafe.Pointer") ||
		strings.Contains(source, "// nolint:unsafe-ptr") {
		return
	}
	start := strings.Index(source, "unsafe.Pointer")
	if start < 0 {
		start = 0
	}
	line, col := unit.LineCol(start)
	rules.PushFinding(
		&MetaPERF51, file, line, col,
		"unsafe.Pointer is used in a request handler; prefer safe alternatives unless a benchmark justifies the pattern", out,
	)
}

func detectPERF52(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee != "runtime.GC" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF52, file, line, col,
			"runtime.GC() forces a stop-the-world GC; remove unless required for tests or controlled shutdown", out,
		)
		return
	}
}

func detectPERF53(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	triggers := []string{"rand.Intn(", "rand.Float64(", "rand.Read("}
	found := false
	for _, t := range triggers {
		if strings.Contains(source, t) {
			found = true
			break
		}
	}
	if !found || strings.Contains(source, "rand.NewSource(") || strings.Contains(source, "rand.New(") {
		return
	}
	start := batch1FirstPos(source, triggers)
	line, col := unit.LineCol(start)
	rules.PushFinding(
		&MetaPERF53, file, line, col,
		"package-level math/rand on a request path contends on a global mutex; use a per-goroutine rand.Source", out,
	)
}

func detectPERF54(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	bytePos := strings.Index(source, "strings.Builder{}")
	if bytePos < 0 {
		return
	}
	if strings.Contains(source, "Reset()") ||
		strings.Contains(source, "var builderPool =") ||
		strings.Contains(source, "sync.Pool") ||
		strings.Contains(source, ".Reset()") {
		return
	}
	inLoop := batch1InLoopByte(facts, bytePos)
	if !IsHotPath(source, bytePos, inLoop) {
		return
	}
	line, col := unit.LineCol(bytePos)
	rules.PushFinding(
		&MetaPERF54, file, line, col,
		"strings.Builder is allocated on a hot path; pool or hoist the builder and call Reset", out,
	)
}

func detectPERF55(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if strings.Contains(source, "bufio.NewScanner(") && !strings.Contains(source, ".Buffer(") {
		start := strings.Index(source, "bufio.NewScanner(")
		if start < 0 {
			start = 0
		}
		line, col := unit.LineCol(start)
		rules.PushFinding(
			&MetaPERF55, file, line, col,
			"bufio.NewScanner is used without an explicit Buffer sizing; large inputs will silently fail at 64KiB", out,
		)
	}
}

func detectPERF56(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee == "c.JSON" && IsInLoop(call) {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(
				&MetaPERF56, file, line, col,
				"c.JSON is called inside a loop body; marshal once and stream or batch the response", out,
			)
			return
		}
	}
}

func detectPERF57(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	if !strings.Contains(source, "*gin.Context") && !strings.Contains(source, "gin.HandlerFunc") {
		return
	}
	if !strings.Contains(source, "c.Next()") {
		return
	}
	trig := []string{"io.ReadAll(", "json.Unmarshal("}
	found := false
	for _, t := range trig {
		if strings.Contains(source, t) {
			found = true
			break
		}
	}
	if !found {
		return
	}
	start := batch1FirstPos(source, trig)
	line, col := unit.LineCol(start)
	rules.PushFinding(
		&MetaPERF57, file, line, col,
		"heavy work in a Gin middleware (io.ReadAll / json.Unmarshal) runs for every request", out,
	)
}

func detectPERF58(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	buffered := []string{
		"io.ReadAll(c.Request.Body)",
		"ioutil.ReadAll(c.Request.Body)",
		"io.ReadAll(r.Body)",
		"ioutil.ReadAll(r.Body)",
		"c.Request.Body.Read(",
	}
	found := false
	for _, t := range buffered {
		if strings.Contains(source, t) {
			found = true
			break
		}
	}
	if !found {
		return
	}
	if strings.Contains(source, "defer c.Request.Body.Close()") ||
		strings.Contains(source, "defer body.Close()") ||
		strings.Contains(source, "io.Copy(io.Discard,") {
		return
	}
	start := batch1FirstPos(source, buffered)
	line, col := unit.LineCol(start)
	rules.PushFinding(
		&MetaPERF58, file, line, col,
		"c.Request.Body is read in a buffered way without deferring Close or draining via io.Copy; the connection may be retained", out,
	)
}

func detectPERF59(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if call.Callee == "c.ShouldBindJSON" {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(
				&MetaPERF59, file, line, col,
				"c.ShouldBindJSON is called per request; consider sharing a pre-validated DTO or per-route binder", out,
			)
			return
		}
	}
}

func detectPERF60(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !batch1RequestPath(unit, facts) {
		return
	}
	trig := []string{
		"render.JSON{", "render.HTML{", "render.IndentedJSON{",
		"render.Redirect{", "render.XML{", "render.YAML{", "render.String{",
	}
	found := false
	for _, t := range trig {
		if strings.Contains(source, t) {
			found = true
			break
		}
	}
	if !found {
		return
	}
	start := batch1FirstPos(source, trig)
	line, col := unit.LineCol(start)
	rules.PushFinding(
		&MetaPERF60, file, line, col,
		"render.Render is allocated directly in a Gin handler; use c.JSON / c.HTML which manage a renderer pool", out,
	)
}

package badpractices

import (
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-101", detectBP101)
	RegisterRule("BP-102", detectBP102)
	RegisterRule("BP-104", detectBP104)
	RegisterRule("BP-105", detectBP105)
	RegisterRule("BP-107", detectBP107)
	RegisterRule("BP-109", detectBP109)
	RegisterRule("BP-110", detectBP110)
	RegisterRule("BP-111", detectBP111)
	RegisterRule("BP-116", detectBP116)
	RegisterRule("BP-117", detectBP117)
	RegisterRule("BP-119", detectBP119)
	RegisterRule("BP-120", detectBP120)
	RegisterRule("BP-122", detectBP122)
	RegisterRule("BP-146", detectBP146)
	RegisterRule("BP-147", detectBP147)
	RegisterRule("BP-149", detectBP149)
	RegisterRule("BP-151", detectBP151)
	RegisterRule("BP-155", detectBP155)
	RegisterRule("BP-156", detectBP156)
	RegisterRule("BP-158", detectBP158)
	RegisterRule("BP-159", detectBP159)
	RegisterRule("BP-160", detectBP160)
}

func detectBP101(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-101")
	src := unit.Source
	if !strings.Contains(src, "WriteHeader") {
		return
	}
	if !strings.Contains(src, "http.ResponseWriter") || !strings.Contains(src, "http.Request") {
		return
	}
	// Collect ResponseWriter parameter names from handler-like signatures.
	writers := httpResponseWriterNames(src)
	if len(writers) == 0 {
		// default common names
		writers = []string{"w", "writer", "rw"}
	}
	// Find body-write positions targeting a writer, and WriteHeader positions.
	var bodyPos, headerPos int = -1, -1
	for _, line := range codeLinesFacts(facts, src) {
		t := strings.TrimSpace(line.text)
		// WriteHeader on a known writer
		for _, w := range writers {
			if strings.Contains(t, w+".WriteHeader(") {
				if headerPos < 0 {
					headerPos = line.byte + strings.Index(line.text, "WriteHeader")
				}
			}
			if strings.Contains(t, w+".Write(") {
				if bodyPos < 0 {
					bodyPos = line.byte
				}
			}
		}
		// fmt.Fprint*/io.WriteString targeting writer
		for _, helper := range []string{"fmt.Fprint(", "fmt.Fprintf(", "fmt.Fprintln(", "io.WriteString("} {
			if !strings.Contains(t, helper) {
				continue
			}
			for _, w := range writers {
				// first arg is writer: helper(w, ...
				if strings.Contains(t, helper+w+",") || strings.Contains(t, helper+w+")") {
					if bodyPos < 0 {
						bodyPos = line.byte
					}
				}
			}
		}
	}
	if bodyPos >= 0 && headerPos >= 0 && bodyPos < headerPos {
		pushAt(unit, meta, headerPos, "response body is written before WriteHeader; set the status before the first body write", out)
	}
}

func httpResponseWriterNames(src string) []string {
	var names []string
	seen := map[string]bool{}
	for _, line := range codeLines(src) {
		t := line.text
		if !strings.Contains(t, "http.ResponseWriter") {
			continue
		}
		// extract identifiers before http.ResponseWriter
		// e.g. func serve(w http.ResponseWriter, r *http.Request)
		// or func(writer http.ResponseWriter, request *http.Request)
		idx := strings.Index(t, "http.ResponseWriter")
		before := t[:idx]
		// last identifier before the type
		fields := strings.FieldsFunc(before, func(r rune) bool {
			return r == '(' || r == ',' || r == ' ' || r == '\t'
		})
		if len(fields) == 0 {
			continue
		}
		name := fields[len(fields)-1]
		if name == "" || name == "func" || name == "type" {
			continue
		}
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}

func detectBP102(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-102")
	src := unit.Source
	if !strings.Contains(src, "http.ResponseWriter") || !strings.Contains(src, "http.Request") {
		return
	}
	// Handler-shaped: ResponseWriter + Request params. Look for err != nil branches
	// that bare-return without writing a response in that branch.
	lines := codeLinesFacts(facts, src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		// if err != nil / if loadErr != nil
		if !strings.Contains(t, "!= nil") || !strings.HasPrefix(t, "if ") {
			continue
		}
		// Scan the if body for bare return without response action.
		depth := 0
		started := false
		hasReturn := false
		hasResponse := false
		for j := i; j < len(lines); j++ {
			lt := lines[j].text
			for _, ch := range lt {
				if ch == '{' {
					depth++
					started = true
				} else if ch == '}' {
					if depth > 0 {
						depth--
					}
				}
			}
			if !started {
				// single-line if without braces is rare for handlers
				continue
			}
			bt := strings.TrimSpace(lt)
			if bt == "return" || strings.HasPrefix(bt, "return ") && !strings.Contains(bt, "http.") {
				// bare return or return of non-http value
				if bt == "return" {
					hasReturn = true
				}
			}
			if strings.Contains(bt, "http.Error") || strings.Contains(bt, "WriteHeader") ||
				strings.Contains(bt, ".Write(") || strings.Contains(bt, "fmt.Fprint") {
				hasResponse = true
			}
			if started && depth == 0 {
				break
			}
		}
		if hasReturn && !hasResponse {
			pushAt(unit, meta, line.byte, "HTTP error path returns without writing an error response or status", out)
			return
		}
	}
}

func detectBP104(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-104")
	// duplicate Handle/HandleFunc patterns
	patterns := map[string]int{}
	for _, line := range codeLinesFacts(facts, unit.Source) {
		t := strings.TrimSpace(line.text)
		for _, meth := range []string{"HandleFunc(", "Handle(", ".GET(", ".POST("} {
			if i := strings.Index(t, meth); i >= 0 {
				// extract first string literal arg
				rest := t[i+len(meth):]
				if strings.HasPrefix(rest, `"`) {
					end := strings.Index(rest[1:], `"`)
					if end >= 0 {
						pat := rest[:end+2]
						if prev, ok := patterns[pat]; ok {
							pushAt(unit, meta, line.byte, "duplicate HTTP mux pattern registration", out)
							_ = prev
							return
						}
						patterns[pat] = line.byte
					}
				}
			}
		}
	}
}

func detectBP105(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-105")
	if !strings.Contains(unit.Source, "http.Cookie") && !strings.Contains(unit.Source, "&http.Cookie{") {
		return
	}
	src := unit.Source
	if strings.Contains(src, "http.Cookie{") || strings.Contains(src, "&http.Cookie{") {
		// extract cookie literal roughly
		if !strings.Contains(src, "HttpOnly") || !strings.Contains(src, "Secure") {
			if pos := strings.Index(src, "Cookie{"); pos >= 0 {
				pushAt(unit, meta, pos, "sensitive cookie missing Secure and/or HttpOnly flags", out)
			}
		}
	}
}

func detectBP107(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-107")
	// middleware func(http.Handler) http.Handler without next.ServeHTTP
	if strings.Contains(unit.Source, "http.Handler") && strings.Contains(unit.Source, "func(") {
		if strings.Contains(unit.Source, "http.HandlerFunc") || strings.Contains(unit.Source, "http.Handler)") {
			if !strings.Contains(unit.Source, "ServeHTTP") && !strings.Contains(unit.Source, "next.ServeHTTP") &&
				!strings.Contains(unit.Source, "next(") {
				// only if looks like middleware
				if strings.Contains(unit.Source, "next http.Handler") || strings.Contains(unit.Source, "http.Handler)") {
					if pos := strings.Index(unit.Source, "http.Handler"); pos >= 0 {
						pushAt(unit, meta, pos, "HTTP middleware never invokes the next handler", out)
					}
				}
			}
		}
	}
}

func detectBP109(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-109")
	src := unit.Source
	if !strings.Contains(src, "gin.") && !strings.Contains(src, "*gin.Context") &&
		!strings.Contains(src, "github.com/gin-gonic/gin") {
		return
	}
	// Any receiver .JSON(error status) without Abort/return in the same block.
	// Supports renamed receivers (ctx *gin.Context) and aliased imports.
	for i, line := range codeLinesFacts(facts, src) {
		t := strings.TrimSpace(line.text)
		if !strings.Contains(t, ".JSON(") {
			continue
		}
		// Error status: 4xx/5xx literals or Status* error constants
		isErrorStatus := strings.Contains(t, "http.StatusBadRequest") ||
			strings.Contains(t, "http.StatusInternalServerError") ||
			strings.Contains(t, "http.StatusNotFound") ||
			strings.Contains(t, "http.StatusUnauthorized") ||
			strings.Contains(t, "http.StatusForbidden") ||
			containsHTTPErrorStatusLiteral(t)
		if !isErrorStatus {
			continue
		}
		// Extract receiver before .JSON
		jsonIdx := strings.Index(t, ".JSON(")
		recv := strings.TrimSpace(t[:jsonIdx])
		// may be "ctx.JSON" or "if ... { ctx.JSON" — take last identifier
		recv = lastIdent(recv)
		if recv == "" {
			continue
		}
		// Check same block for Abort or return after this statement.
		if ginErrorTerminated(src, i, recv) {
			continue
		}
		pushAt(unit, meta, line.byte, "Gin error response without Abort; handler may continue", out)
		return
	}
}

func containsHTTPErrorStatusLiteral(t string) bool {
	// Match bare 4xx/5xx as first arg: .JSON(500, or .JSON(404,
	idx := strings.Index(t, ".JSON(")
	if idx < 0 {
		return false
	}
	rest := t[idx+len(".JSON("):]
	rest = strings.TrimSpace(rest)
	// leading digits
	n := 0
	for n < len(rest) && rest[n] >= '0' && rest[n] <= '9' {
		n++
	}
	if n == 0 {
		return false
	}
	// parse roughly
	code := 0
	for i := 0; i < n; i++ {
		code = code*10 + int(rest[i]-'0')
	}
	return code >= 400 && code <= 599
}

func lastIdent(s string) string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_')
	})
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}

func ginErrorTerminated(src string, lineIdx int, recv string) bool {
	lines := codeLines(src)
	// Look at remaining lines in the same if-block after the JSON call.
	depth := 0
	started := false
	for j := lineIdx; j < len(lines); j++ {
		lt := lines[j].text
		for _, ch := range lt {
			if ch == '{' {
				depth++
				started = true
			} else if ch == '}' {
				if depth > 0 {
					depth--
				}
			}
		}
		if j == lineIdx {
			// same line may have return after JSON — rare
			if strings.Contains(lt, "return") {
				return true
			}
			continue
		}
		bt := strings.TrimSpace(lt)
		if strings.HasPrefix(bt, "return") {
			return true
		}
		if strings.Contains(bt, recv+".Abort") || strings.Contains(bt, "AbortWithStatus") {
			return true
		}
		if started && depth == 0 {
			// left the block without abort/return
			return false
		}
	}
	// If no braces (JSON alone in if), check next non-empty is return/Abort.
	for j := lineIdx + 1; j < len(lines) && j < lineIdx+4; j++ {
		bt := strings.TrimSpace(lines[j].text)
		if bt == "" || bt == "}" {
			if bt == "}" {
				return false
			}
			continue
		}
		if strings.HasPrefix(bt, "return") || strings.Contains(bt, ".Abort") {
			return true
		}
		return false
	}
	return false
}

func detectBP110(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-110")
	if !strings.Contains(unit.Source, "ShouldBind") && !strings.Contains(unit.Source, "BindJSON") && !strings.Contains(unit.Source, "Bind(") {
		return
	}
	msg := "Gin bind error is discarded; check the error and abort the request"
	bindNames := []string{"ShouldBindJSON", "ShouldBind", "BindJSON", "MustBindWith", "ShouldBindQuery"}
	for _, line := range codeLinesFacts(facts, unit.Source) {
		t := strings.TrimSpace(line.text)
		for _, b := range bindNames {
			if strings.Contains(t, b+"(") {
				// bare statement: c.ShouldBindJSON(&x)
				if !strings.Contains(t, ":=") && !strings.Contains(t, "err") && !strings.HasPrefix(t, "if ") {
					pushAt(unit, meta, line.byte, msg, out)
				}
			}
		}
	}
	_ = facts
}

func detectBP111(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-111")
	src := unit.Source
	if !strings.Contains(src, "go ") && !strings.Contains(src, "go\t") {
		return
	}
	if !strings.Contains(src, "gin.Context") && !strings.Contains(src, "*gin.Context") &&
		!strings.Contains(src, "github.com/gin-gonic/gin") {
		return
	}
	// Context parameter names from *gin.Context signatures.
	ctxNames := ginContextParamNames(src)
	if len(ctxNames) == 0 {
		ctxNames = []string{"c", "ctx"}
	}
	// For each go statement / go func, if body references a context param
	// without an intervening .Copy() binding used instead, fire.
	lines := codeLinesFacts(facts, src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "go ") && !strings.HasPrefix(t, "go\t") && !strings.HasPrefix(t, "go func") {
			// also "go func()" on same style
			if !(strings.HasPrefix(t, "go") && strings.Contains(t, "func")) {
				continue
			}
		}
		// Collect goroutine body text
		body := t
		depth := 0
		started := false
		for j := i; j < len(lines); j++ {
			lt := lines[j].text
			body += "\n" + lt
			for _, ch := range lt {
				if ch == '{' {
					depth++
					started = true
				} else if ch == '}' {
					if depth > 0 {
						depth--
					}
				}
			}
			if started && depth == 0 {
				break
			}
		}
		// Does body use a context param (c. or c) but not only a copy?
		// Safe pattern: copy := c.Copy(); go func() { copy.Request... }
		usesOriginal := false
		for _, name := range ctxNames {
			// use of name. inside goroutine
			if strings.Contains(body, name+".") {
				// If this name was assigned from .Copy() before the go, skip.
				// Look at preceding lines in the same function for `x := name.Copy()`
				// and body uses x not name.
				usesOriginal = true
				// Check if all uses are of a copy variable
				copyVars := copyVarsFromContext(src, name)
				onlyCopy := false
				if len(copyVars) > 0 {
					onlyCopy = true
					// if body still has name. (not copy.), still original
					// strip copy var uses and see if name remains
					stripped := body
					for _, cv := range copyVars {
						stripped = strings.ReplaceAll(stripped, cv+".", "")
					}
					if strings.Contains(stripped, name+".") {
						onlyCopy = false
					} else {
						// body only uses copy vars
						usesOriginal = false
					}
				}
				_ = onlyCopy
			}
		}
		if usesOriginal {
			pushAt(unit, meta, line.byte, "Gin context used from a goroutine; copy values before leaving the handler", out)
			return
		}
	}
}

func ginContextParamNames(src string) []string {
	var names []string
	seen := map[string]bool{}
	for _, line := range codeLines(src) {
		t := line.text
		if !strings.Contains(t, "gin.Context") {
			continue
		}
		// c *gin.Context or ctx *gin.Context
		idx := strings.Index(t, "gin.Context")
		before := t[:idx]
		before = strings.TrimSuffix(strings.TrimSpace(before), "*")
		fields := strings.FieldsFunc(before, func(r rune) bool {
			return r == '(' || r == ',' || r == ' ' || r == '\t'
		})
		if len(fields) == 0 {
			continue
		}
		name := fields[len(fields)-1]
		if name == "" || name == "func" || name == "*" {
			continue
		}
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}

func copyVarsFromContext(src, ctxName string) []string {
	var vars []string
	needle := ctxName + ".Copy()"
	for _, line := range codeLines(src) {
		t := strings.TrimSpace(line.text)
		if !strings.Contains(t, needle) {
			continue
		}
		lhs, _, ok := strings.Cut(t, ":=")
		if !ok {
			lhs, _, ok = strings.Cut(t, "=")
		}
		if !ok {
			continue
		}
		name := strings.TrimSpace(lhs)
		if i := strings.IndexAny(name, " \t"); i >= 0 {
			name = name[:i]
		}
		if name != "" && name != "_" {
			vars = append(vars, name)
		}
	}
	return vars
}

func detectBP116(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-116")
	src := unit.Source
	if !strings.Contains(src, "echo.") && !strings.Contains(src, "echo.Context") &&
		!strings.Contains(src, "github.com/labstack/echo") {
		return
	}
	// Bare error JSON (expression statement) followed by return of a business error.
	// Safe pattern binds JSON's error: if responseErr := c.JSON(...); responseErr != nil { return responseErr }
	lines := codeLinesFacts(facts, src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.Contains(t, ".JSON(") {
			continue
		}
		isErrorStatus := strings.Contains(t, "http.StatusInternalServerError") ||
			strings.Contains(t, "http.StatusBadRequest") ||
			strings.Contains(t, "http.StatusNotFound") ||
			containsHTTPErrorStatusLiteral(t)
		if !isErrorStatus {
			continue
		}
		// If JSON result is assigned/checked, returning that response error is OK.
		if strings.Contains(t, ":=") || strings.Contains(t, "err =") || strings.HasPrefix(t, "if ") {
			continue
		}
		// Bare expression statement: c.JSON(...) then later return saveErr/dbErr/err
		depth := 0
		started := false
		// Track depth from surrounding if-block: start from previous open if needed.
		// Use lines after the JSON call within the same brace depth.
		baseDepth := 0
		for _, ch := range t {
			if ch == '{' {
				baseDepth++
			} else if ch == '}' {
				baseDepth--
			}
		}
		for j := i + 1; j < len(lines); j++ {
			lt := lines[j].text
			bt := strings.TrimSpace(lt)
			// Count braces on this line for block end
			lineDepthDelta := 0
			for _, ch := range lt {
				if ch == '{' {
					lineDepthDelta++
					started = true
					depth++
				} else if ch == '}' {
					lineDepthDelta--
					depth--
				}
			}
			if strings.HasPrefix(bt, "return ") {
				ret := strings.TrimSpace(strings.TrimPrefix(bt, "return "))
				// raw error return: single identifier ending with err/Err, not nil
				if ret != "nil" && ret != "" && !strings.Contains(ret, "(") && !strings.Contains(ret, " ") {
					if ret == "err" || strings.HasSuffix(ret, "Err") || strings.HasSuffix(ret, "err") {
						// Exclude response/json binding names that look like response write errors
						low := strings.ToLower(ret)
						if strings.Contains(low, "response") || strings.Contains(low, "json") {
							continue
						}
						pushAt(unit, meta, lines[j].byte, "Echo handler writes an error response and also returns the raw error", out)
						return
					}
				}
			}
			// Left the enclosing block (closing brace without nested open)
			if depth < 0 || (lineDepthDelta < 0 && depth <= 0 && bt == "}") {
				break
			}
			_ = started
		}
	}
}

func detectBP117(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-117")
	if !strings.Contains(unit.Source, "echo") {
		return
	}
	for _, line := range codeLinesFacts(facts, unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, ".Bind(") && !strings.Contains(t, "err") && !strings.Contains(t, ":=") && !strings.HasPrefix(t, "if ") {
			pushAt(unit, meta, line.byte, "Echo Bind error is discarded", out)
		}
	}
}

func detectBP119(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-119")
	if !strings.Contains(unit.Source, "fiber") && !strings.Contains(unit.Source, "*fiber.Ctx") {
		return
	}
	if strings.Contains(unit.Source, "go ") && strings.Contains(unit.Source, "*fiber.Ctx") {
		if pos := strings.Index(unit.Source, "go "); pos >= 0 {
			pushAt(unit, meta, pos, "Fiber context used from a goroutine after the handler may return", out)
		}
	}
}

func detectBP120(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-120")
	if !strings.Contains(unit.Source, "BodyParser") {
		return
	}
	for _, line := range codeLinesFacts(facts, unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "BodyParser(") && !strings.Contains(t, "err") && !strings.Contains(t, ":=") && !strings.HasPrefix(t, "if ") {
			pushAt(unit, meta, line.byte, "Fiber BodyParser error is discarded", out)
		}
	}
}

func detectBP122(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-122")
	if !strings.Contains(unit.Source, "chi.") && !strings.Contains(unit.Source, "go-chi") {
		return
	}
	if strings.Contains(unit.Source, "func(") && strings.Contains(unit.Source, "http.Handler") {
		if !strings.Contains(unit.Source, "next.ServeHTTP") && !strings.Contains(unit.Source, "next(") {
			if pos := strings.Index(unit.Source, "http.Handler"); pos >= 0 {
				pushAt(unit, meta, pos, "Chi middleware never invokes the next handler", out)
			}
		}
	}
}

func detectBP146(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-146")
	src := unit.Source
	if !(strings.Contains(src, "log.") || strings.Contains(src, "slog.") || strings.Contains(src, "zap.")) {
		return
	}
	for _, line := range codeLinesFacts(facts, src) {
		t := line.text
		lt := strings.ToLower(t)
		// Must be a logger call line
		if !(strings.Contains(lt, "log.") || strings.Contains(lt, "slog.") || strings.Contains(lt, "zap.") ||
			strings.Contains(lt, "logger.")) {
			continue
		}
		// Redaction / hashing suppresses the finding.
		if containsIdentCI(lt, "redact") || containsIdentCI(lt, "redacted") ||
			containsIdentCI(lt, "mask") || containsIdentCI(lt, "masked") ||
			containsIdentCI(lt, "hash") {
			continue
		}
		if !containsSensitiveLoggedValue(lt) {
			continue
		}
		pushAt(unit, meta, line.byte, "sensitive field may be logged; redact secrets before logging", out)
		return
	}
}

func containsIdentCI(lower, word string) bool {
	// word-boundary-ish match for lowercase text
	idx := 0
	for {
		i := strings.Index(lower[idx:], word)
		if i < 0 {
			return false
		}
		abs := idx + i
		beforeOK := abs == 0 || !isIdentByte(lower[abs-1])
		after := abs + len(word)
		afterOK := after >= len(lower) || !isIdentByte(lower[after])
		if beforeOK && afterOK {
			return true
		}
		idx = abs + len(word)
	}
}

func containsSensitiveLoggedValue(lowerLine string) bool {
	// Sensitive names that look like credentials.
	names := []string{
		"password", "passwd", "secret", "token", "authorization",
		"api_key", "apikey", "private_key", "client_secret", "ssn", "credit_card",
	}
	hasSensitive := false
	for _, n := range names {
		if containsIdentCI(lowerLine, n) || strings.Contains(lowerLine, n) {
			// avoid credential_redacted alone without value path
			hasSensitive = true
			break
		}
	}
	if !hasSensitive {
		return false
	}
	// Boolean presence (%t) without serializing the secret is OK.
	hasFormatValue := strings.Contains(lowerLine, "%s") || strings.Contains(lowerLine, "%v") ||
		strings.Contains(lowerLine, "%q") || strings.Contains(lowerLine, "%x") ||
		strings.Contains(lowerLine, "%d")
	if !hasFormatValue && strings.Contains(lowerLine, "%t") {
		return false
	}
	// slog style: "password", password  (name appears twice) or format verbs
	// credential_redacted alone without a format of the secret: if redacted already handled
	if strings.Contains(lowerLine, "redacted") || strings.Contains(lowerLine, "_redacted") {
		return false
	}
	if hasFormatValue {
		return true
	}
	// structured: key and variable both present e.g. "password", password
	for _, n := range names {
		if strings.Count(lowerLine, n) >= 2 {
			return true
		}
	}
	return false
}

func detectBP147(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-147")
	if packageName(unit.Source) == "main" {
		return
	}
	src := unit.Source
	// Require real stdlib log import + structured logger import, and a log.Print* call.
	hasStdLogImport := strings.Contains(src, `"log"`) // import "log"
	hasStructured := strings.Contains(src, "log/slog") || strings.Contains(src, "go.uber.org/zap") ||
		strings.Contains(src, "logrus")
	if !hasStdLogImport || !hasStructured {
		return
	}
	// Fire only on log.Print / log.Printf / log.Println (not slog which contains "log.")
	for _, line := range codeLinesFacts(facts, src) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "log.Print") || strings.Contains(t, " log.Print") {
			// ensure not slog
			if strings.Contains(t, "slog.") {
				continue
			}
			pushAt(unit, meta, line.byte, "service package mixes standard log with structured logging", out)
			return
		}
		// mid-line: log.Printf(
		if (strings.Contains(t, "log.Print(") || strings.Contains(t, "log.Printf(") ||
			strings.Contains(t, "log.Println(")) && !strings.Contains(t, "slog.") {
			pushAt(unit, meta, line.byte, "service package mixes standard log with structured logging", out)
			return
		}
	}
}

func detectBP149(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-149")
	// err != nil { slog.Error("msg") } without err attribute
	lines := codeLinesFacts(facts, unit.Source)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "err != nil") {
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				n := strings.TrimSpace(lines[j].text)
				if n == "}" {
					break
				}
				if (strings.Contains(n, "slog.Error") || strings.Contains(n, "log.Error") || strings.Contains(n, ".Error(")) &&
					!strings.Contains(n, "err") {
					pushAt(unit, meta, lines[j].byte, "error log omits the error value as an attribute", out)
					return
				}
			}
		}
	}
}

func detectBP151(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-151")
	src := unit.Source
	if !strings.Contains(src, "os.Getenv") && !strings.Contains(src, "os.LookupEnv") {
		return
	}
	if !(strings.Contains(src, "log.") || strings.Contains(src, "slog.") || strings.Contains(src, "zap.")) {
		return
	}
	// Direct os.Getenv("SECRET...") inside a logger call, or logging the env result
	// with a value format (not mere presence).
	for _, line := range codeLinesFacts(facts, src) {
		t := strings.TrimSpace(line.text)
		lt := strings.ToLower(t)
		isLog := strings.Contains(lt, "log.") || strings.Contains(lt, "slog.") || strings.Contains(lt, "zap.")
		if !isLog {
			continue
		}
		// Presence-only is safe: %t or != ""
		if strings.Contains(t, "%t") && !strings.Contains(t, "%s") && !strings.Contains(t, "%v") &&
			!strings.Contains(t, "%q") {
			continue
		}
		// Direct getenv of secret in the log line
		if (strings.Contains(t, "os.Getenv") || strings.Contains(t, "os.LookupEnv")) &&
			secretEnvNameInCall(t) {
			pushAt(unit, meta, line.byte, "secret environment value may be logged", out)
			return
		}
	}
	// Variable form: token := os.Getenv("API_TOKEN"); log.Printf("%s", token)
	// Only fire when the value (not presence) is logged.
	var secretVars []string
	for _, line := range codeLinesFacts(facts, src) {
		t := strings.TrimSpace(line.text)
		if !(strings.Contains(t, "os.Getenv") || strings.Contains(t, "os.LookupEnv")) {
			continue
		}
		if !secretEnvNameInCall(t) {
			continue
		}
		lhs, _, ok := strings.Cut(t, ":=")
		if !ok {
			lhs, _, ok = strings.Cut(t, "=")
		}
		if !ok {
			continue
		}
		name := strings.TrimSpace(lhs)
		if i := strings.IndexAny(name, " \t,"); i >= 0 {
			name = strings.TrimSpace(name[:i])
		}
		if name != "" && name != "_" {
			secretVars = append(secretVars, name)
		}
	}
	for _, line := range codeLinesFacts(facts, src) {
		t := strings.TrimSpace(line.text)
		lt := strings.ToLower(t)
		isLog := strings.Contains(lt, "log.") || strings.Contains(lt, "slog.") || strings.Contains(lt, "zap.")
		if !isLog {
			continue
		}
		// presence checks: token != "" with %t
		if strings.Contains(t, "%t") || strings.Contains(t, "!= \"\"") || strings.Contains(t, "!=\"\"") {
			continue
		}
		for _, v := range secretVars {
			if containsIdentCI(lt, strings.ToLower(v)) {
				// value format present or structured field
				if strings.Contains(t, "%s") || strings.Contains(t, "%v") || strings.Contains(t, "%q") ||
					strings.Contains(t, v) {
					pushAt(unit, meta, line.byte, "secret environment value may be logged", out)
					return
				}
			}
		}
	}
}

func secretEnvNameInCall(t string) bool {
	upper := strings.ToUpper(t)
	return strings.Contains(upper, "SECRET") || strings.Contains(upper, "PASSWORD") ||
		strings.Contains(upper, "TOKEN") || strings.Contains(upper, "API_KEY") ||
		strings.Contains(upper, "PRIVATE_KEY") || strings.Contains(upper, "CREDENTIAL")
}

func detectBP155(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-155")
	src := unit.Source
	if !strings.Contains(src, "json.NewDecoder") {
		return
	}
	if strings.Contains(src, "MaxBytesReader") || strings.Contains(src, "http.MaxBytesReader") {
		return
	}
	// json.NewDecoder(...Body...).Decode without MaxBytesReader
	for _, line := range codeLinesFacts(facts, src) {
		t := strings.TrimSpace(line.text)
		if !strings.Contains(t, "json.NewDecoder") && !strings.Contains(t, ".Decode(") {
			// chained: may be one line json.NewDecoder(request.Body).Decode
		}
		if strings.Contains(t, "json.NewDecoder") && strings.Contains(t, ".Body") {
			pushAt(unit, meta, line.byte, "JSON request body decoded without a size limit", out)
			return
		}
		// Decode line may be separate; check whole source for NewDecoder(...Body)
	}
	// Multi-line: NewDecoder on one line with .Body anywhere in the decoder expression.
	if strings.Contains(src, "json.NewDecoder") {
		// Find NewDecoder(...) spanning and check for .Body inside args
		idx := 0
		for {
			pos := strings.Index(src[idx:], "json.NewDecoder(")
			if pos < 0 {
				break
			}
			abs := idx + pos
			// scan to matching close paren
			depth := 0
			end := abs
			for end < len(src) {
				if src[end] == '(' {
					depth++
				} else if src[end] == ')' {
					depth--
					if depth == 0 {
						end++
						break
					}
				}
				end++
			}
			arg := src[abs:end]
			if strings.Contains(arg, ".Body") {
				pushAt(unit, meta, abs, "JSON request body decoded without a size limit", out)
				return
			}
			idx = abs + 1
		}
	}
}

func detectBP156(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-156")
	// json omitempty on password/secret fields
	for _, line := range codeLinesFacts(facts, unit.Source) {
		t := strings.ToLower(line.text)
		if strings.Contains(t, "json:") && strings.Contains(t, "omitempty") {
			if strings.Contains(t, "password") || strings.Contains(t, "secret") || strings.Contains(t, "token") {
				pushAt(unit, meta, line.byte, "security-sensitive JSON field uses omitempty", out)
			}
		}
	}
}

func detectBP158(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-158")
	if !strings.Contains(unit.Source, "grpc") && !strings.Contains(unit.Source, "status.Error") {
		// grpc service methods often return (resp, error)
		if !strings.Contains(unit.Source, "context.Context") {
			return
		}
	}
	if strings.Contains(unit.Source, "return nil, err") && strings.Contains(unit.Source, "context.Context") {
		if !strings.Contains(unit.Source, "status.Error") && !strings.Contains(unit.Source, "status.Errorf") {
			// only if looks like grpc (Request/Response types)
			if strings.Contains(unit.Source, "Request") && strings.Contains(unit.Source, "Response") {
				if pos := strings.Index(unit.Source, "return nil, err"); pos >= 0 {
					pushAt(unit, meta, pos, "gRPC handler returns a raw error; wrap with status.Error", out)
				}
			}
		}
	}
}

func detectBP159(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-159")
	src := unit.Source
	if !strings.Contains(src, "flag.") {
		return
	}
	parsePos := strings.Index(src, "flag.Parse(")
	if parsePos < 0 {
		return
	}
	// Collect flag pointer variable names declared before Parse.
	constructors := []string{
		"flag.Bool(", "flag.Duration(", "flag.Float64(", "flag.Int(", "flag.Int64(",
		"flag.String(", "flag.Uint(", "flag.Uint64(",
	}
	var flagNames []string
	for _, line := range codeLinesFacts(facts, src) {
		if line.byte >= parsePos {
			break
		}
		t := strings.TrimSpace(line.text)
		hasCtor := false
		for _, c := range constructors {
			if strings.Contains(t, c) {
				hasCtor = true
				break
			}
		}
		if !hasCtor {
			continue
		}
		lhs, _, ok := strings.Cut(t, ":=")
		if !ok {
			lhs, _, ok = strings.Cut(t, "=")
		}
		if !ok {
			continue
		}
		name := strings.TrimSpace(lhs)
		if i := strings.IndexAny(name, " \t,"); i >= 0 {
			name = strings.TrimSpace(name[:i])
		}
		if name != "" && name != "_" {
			flagNames = append(flagNames, name)
		}
	}
	if len(flagNames) == 0 {
		return
	}
	// Dereference *name before parsePos.
	for _, line := range codeLinesFacts(facts, src) {
		if line.byte >= parsePos {
			break
		}
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "flag.") || strings.Contains(t, "flag.Bool") ||
			strings.Contains(t, "flag.Int") || strings.Contains(t, "flag.String") {
			continue
		}
		for _, name := range flagNames {
			// *name as expression
			if strings.Contains(t, "*"+name) {
				// avoid *testing or pointer types in signatures
				// match *name as whole use
				if strings.Contains(t, "*"+name+")") || strings.Contains(t, "*"+name+",") ||
					strings.Contains(t, "*"+name+" ") || strings.HasSuffix(t, "*"+name) ||
					strings.Contains(t, "(*"+name) || strings.Contains(t, " "+"*"+name) {
					pushAt(unit, meta, line.byte, "flag value is read before flag.Parse() processes command-line arguments", out)
					return
				}
				// simple: *verbose / *port alone in call args
				if strings.Contains(t, "*"+name) {
					pushAt(unit, meta, line.byte, "flag value is read before flag.Parse() processes command-line arguments", out)
					return
				}
			}
		}
	}
}

func detectBP160(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-160")
	if !strings.Contains(unit.Source, "cobra.Command") && !strings.Contains(unit.Source, "cobra.") {
		return
	}
	if strings.Contains(unit.Source, "Run:") && !strings.Contains(unit.Source, "RunE:") {
		if pos := strings.Index(unit.Source, "Run:"); pos >= 0 {
			pushAt(unit, meta, pos, "Cobra command uses Run instead of RunE; errors are swallowed", out)
		}
	}
}

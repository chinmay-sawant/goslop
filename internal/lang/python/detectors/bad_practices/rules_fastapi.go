package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-PY-29", detectBPPY29)
	RegisterRule("BP-PY-30", detectBPPY30)
	RegisterRule("BP-PY-31", detectBPPY31)
	RegisterRule("BP-PY-32", detectBPPY32)
}

var (
	fastapiRouteDecoRe = regexp.MustCompile(`@(?:app|router)\.(get|post|put|patch|delete|head|options|api_route|websocket)\s*\(`)
	fastapiDepDecoRe   = regexp.MustCompile(`@\w+\.dependency\s*\(`)
	blockingCallRe     = regexp.MustCompile(`\b(?:time\.sleep\s*\(|requests\.(?:get|post|put|patch|delete|request|head|options)\s*\(|requests\.Session\s*\(|subprocess\.(?:run|call|Popen|check_output|check_call)\s*\()`)
	ormReturnRe        = regexp.MustCompile(`(?i)\b(?:session\.query|db\.query|\.query\.(?:get|filter|all|first)|session\.get\s*\(|db\.get\s*\(|models\.\w+)`)
)

func looksFastAPIish(unit *core.ParsedUnit, src string) bool {
	if strings.Contains(src, "fastapi") || strings.Contains(src, "FastAPI") ||
		strings.Contains(src, "APIRouter") || strings.Contains(src, "starlette") ||
		strings.Contains(src, "Starlette") {
		return true
	}
	if strings.Contains(src, "@app.") || strings.Contains(src, "@router.") {
		return true
	}
	p := strings.ToLower(fileDisplayPath(unit))
	return strings.Contains(p, "fastapi") || strings.Contains(p, "api.py") || strings.Contains(p, "routes")
}

// functionBodyRange returns inclusive line indices [start, end) for the body of a
// def/async def that starts at defLineIdx (the line with def/async def).
func functionBodyRange(lines []codeLine, defLineIdx int) (bodyStart, bodyEnd int) {
	if defLineIdx < 0 || defLineIdx >= len(lines) {
		return defLineIdx, defLineIdx
	}
	// Find signature end (may be multi-line) then body starts after the ':' line.
	sigEnd := defLineIdx
	sig := strings.TrimSpace(lines[defLineIdx].text)
	for !signatureComplete(sig) && sigEnd+1 < len(lines) {
		sigEnd++
		sig += " " + strings.TrimSpace(lines[sigEnd].text)
		if sigEnd-defLineIdx > 40 {
			break
		}
	}
	defIndent := indentWidth(lines[defLineIdx].raw)
	bodyStart = sigEnd + 1
	bodyEnd = len(lines)
	for j := bodyStart; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		if indentWidth(lines[j].raw) <= defIndent {
			bodyEnd = j
			break
		}
	}
	return bodyStart, bodyEnd
}

// collectDecoratorsAbove gathers contiguous decorator lines immediately above defLineIdx.
func collectDecoratorsAbove(lines []codeLine, defLineIdx int) string {
	if defLineIdx <= 0 {
		return ""
	}
	var parts []string
	for i := defLineIdx - 1; i >= 0; i-- {
		t := strings.TrimSpace(lines[i].text)
		if t == "" {
			continue
		}
		if strings.HasPrefix(t, "@") {
			parts = append([]string{t}, parts...)
			// Multi-line decorator: keep going if previous non-empty was also @ or continuation.
			continue
		}
		// Allow multi-line decorator args: lines that are indented continuation of @...
		if len(parts) > 0 && indentWidth(lines[i].raw) > 0 && !strings.HasPrefix(t, "def ") &&
			!strings.HasPrefix(t, "async def ") && !strings.HasPrefix(t, "class ") {
			parts = append([]string{t}, parts...)
			continue
		}
		break
	}
	return strings.Join(parts, " ")
}

func isFastAPIRouteOrDep(decoBlob string, defLine string) bool {
	if fastapiRouteDecoRe.MatchString(decoBlob) || fastapiDepDecoRe.MatchString(decoBlob) {
		return true
	}
	if strings.Contains(decoBlob, "@app.") || strings.Contains(decoBlob, "@router.") {
		return true
	}
	// Dependency functions often take Depends(...) in params without route decorator.
	if strings.Contains(defLine, "Depends(") || strings.Contains(decoBlob, "Depends(") {
		return true
	}
	return false
}

// BP-PY-29: mutable global / module-level mutation in FastAPI routes or deps.
func detectBPPY29(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-29")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	if !looksFastAPIish(unit, src) {
		return
	}
	if !facts.hasAny("global ", "nonlocal ") &&
		!strings.Contains(src, "global ") &&
		!strings.Contains(src, "nonlocal ") {
		// Still may have STORE[k] = v style — continue if module-level mutables present.
		if !strings.Contains(src, " = {}") && !strings.Contains(src, " = []") {
			return
		}
	}

	// Module-level names assigned to {} or [] (simple heuristic).
	moduleMutables := map[string]struct{}{}
	lines := codeLinesFacts(facts, src)
	for _, line := range lines {
		if indentWidth(line.raw) != 0 {
			continue
		}
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") ||
			strings.HasPrefix(t, "class ") || strings.HasPrefix(t, "@") {
			continue
		}
		// NAME = {} or NAME = []
		if i := strings.Index(t, "="); i > 0 {
			lhs := strings.TrimSpace(t[:i])
			rhs := strings.TrimSpace(t[i+1:])
			if isSimpleIdent(lhs) && (rhs == "{}" || rhs == "[]" ||
				strings.HasPrefix(rhs, "dict(") || strings.HasPrefix(rhs, "list(")) {
				moduleMutables[lhs] = struct{}{}
			}
		}
	}

	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "def ") && !strings.HasPrefix(t, "async def ") {
			continue
		}
		deco := collectDecoratorsAbove(lines, i)
		if !isFastAPIRouteOrDep(deco, t) && !looksLikeRouteDef(t, deco) {
			// Still scan dependency-like names or any def under FastAPI module with global.
			// Prefer route/dep; for global keyword any function under FastAPI module is signal.
			if !strings.Contains(t, "Depends") && deco == "" {
				// allow later body scan only for global/nonlocal if FastAPI module
			}
		}
		bodyStart, bodyEnd := functionBodyRange(lines, i)
		routeOrDep := isFastAPIRouteOrDep(deco, t) || looksLikeRouteDef(t, deco)
		for j := bodyStart; j < bodyEnd; j++ {
			bt := strings.TrimSpace(lines[j].text)
			if bt == "" {
				continue
			}
			if strings.HasPrefix(bt, "global ") || strings.HasPrefix(bt, "nonlocal ") {
				if routeOrDep || looksFastAPIish(unit, src) {
					pushAt(unit, meta, lines[j].byte, "FastAPI route/dependency mutates global state; prefer request-scoped deps or a proper store", out)
					return // one finding per function is enough for v0
				}
			}
			// Module-level mutable mutation: STORE[k] = or STORE.append(
			if routeOrDep {
				for name := range moduleMutables {
					if strings.Contains(bt, name+"[") || strings.Contains(bt, name+".append") ||
						strings.Contains(bt, name+".update") || strings.Contains(bt, name+".pop") ||
						(strings.HasPrefix(bt, name+" ") && strings.Contains(bt, "=")) ||
						strings.HasPrefix(bt, name+"=") {
						// Avoid flagging local shadowing assignments of same name as pure local.
						// Module mutation of known module mutables is the hit.
						pushAt(unit, meta, lines[j].byte, "FastAPI route/dependency mutates module-level mutable state; prefer request-scoped deps", out)
						return
					}
				}
			}
		}
	}
}

func looksLikeRouteDef(defLine, deco string) bool {
	if strings.Contains(deco, "@app.") || strings.Contains(deco, "@router.") {
		return true
	}
	// Common dep name patterns
	if strings.Contains(defLine, "def get_") || strings.Contains(defLine, "async def get_") {
		return strings.Contains(defLine, "Depends")
	}
	return false
}

func isSimpleIdent(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isIdentByte(s[i]) {
			return false
		}
	}
	if s[0] >= '0' && s[0] <= '9' {
		return false
	}
	return true
}

// BP-PY-30: blocking I/O inside async FastAPI route/dependency bodies.
func detectBPPY30(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-30")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	if !looksFastAPIish(unit, src) {
		return
	}
	if !facts.has("async def ") && !strings.Contains(src, "async def ") {
		return
	}
	// Fast path: need some blocking needle present.
	if !facts.hasAny("time.sleep", "requests.", "subprocess.") &&
		!strings.Contains(src, "time.sleep") &&
		!strings.Contains(src, "requests.") &&
		!strings.Contains(src, "subprocess.") {
		// Optional sqlalchemy sync session in async — still check if present.
		if !strings.Contains(src, "session.query") && !strings.Contains(src, "Session(") {
			return
		}
	}

	lines := codeLinesFacts(facts, src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "async def ") {
			continue
		}
		deco := collectDecoratorsAbove(lines, i)
		// Prefer route-scoped; also flag async def with Depends or under FastAPI with @app/@router.
		if !isFastAPIRouteOrDep(deco, t) && !looksLikeRouteDef(t, deco) {
			// If module is FastAPI and has route decorators somewhere, still only fire on decorated async routes.
			continue
		}
		bodyStart, bodyEnd := functionBodyRange(lines, i)
		for j := bodyStart; j < bodyEnd; j++ {
			bt := lines[j].text
			// Skip await asyncio.sleep — not blocking for this needle if only sleep via await.
			if blockingCallRe.MatchString(bt) {
				// Ensure time.sleep is not part of a comment already stripped.
				// Miss path: await asyncio.sleep is different needle.
				if strings.Contains(bt, "asyncio.sleep") {
					continue
				}
				loc := blockingCallRe.FindStringIndex(bt)
				off := lines[j].byte
				if loc != nil {
					off += loc[0]
				}
				pushAt(unit, meta, off, "blocking I/O in async FastAPI route; use await asyncio.sleep/httpx.AsyncClient or a sync def route", out)
				// One per function body is enough.
				break
			}
			// Sync sqlalchemy-ish in async route.
			if strings.Contains(bt, "session.query") || strings.Contains(bt, "session.execute") ||
				strings.Contains(bt, "db.query(") {
				pushAt(unit, meta, lines[j].byte, "sync ORM/DB call in async FastAPI route blocks the event loop; use async session or run_in_executor", out)
				break
			}
		}
	}
}

// BP-PY-31: route without response_model returning ORM-ish values.
func detectBPPY31(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-31")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	if !looksFastAPIish(unit, src) {
		return
	}
	// Need ORM-ish markers or sqlalchemy/django models import for precision.
	hasORM := strings.Contains(src, "sqlalchemy") || strings.Contains(src, "session.query") ||
		strings.Contains(src, "Session") || strings.Contains(src, "models.") ||
		strings.Contains(src, "django.db") || strings.Contains(src, ".query.")
	if !hasORM {
		return
	}

	lines := codeLinesFacts(facts, src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "def ") && !strings.HasPrefix(t, "async def ") {
			continue
		}
		deco := collectDecoratorsAbove(lines, i)
		if !fastapiRouteDecoRe.MatchString(deco) && !strings.Contains(deco, "@app.") && !strings.Contains(deco, "@router.") {
			continue
		}
		if strings.Contains(deco, "response_model") {
			continue
		}
		bodyStart, bodyEnd := functionBodyRange(lines, i)
		for j := bodyStart; j < bodyEnd; j++ {
			bt := strings.TrimSpace(lines[j].text)
			if !strings.HasPrefix(bt, "return ") {
				continue
			}
			ret := strings.TrimSpace(strings.TrimPrefix(bt, "return "))
			// Miss: dict/list/literal/Pydantic construction HeuristicOut(...)
			if ret == "" || ret == "None" || isPlainStringLiteral(ret) ||
				strings.HasPrefix(ret, "{") || strings.HasPrefix(ret, "[") ||
				strings.HasPrefix(ret, "(") {
				continue
			}
			// ORM-ish return needles
			if ormReturnRe.MatchString(ret) || ormReturnRe.MatchString(bt) {
				pushAt(unit, meta, lines[j].byte, "FastAPI route returns ORM-like object without response_model; prefer a Pydantic response model", out)
				break
			}
			// return session.get / return db_user after query patterns in body
			if strings.Contains(ret, "session.") || strings.Contains(ret, "db.") ||
				strings.Contains(ret, ".query") {
				pushAt(unit, meta, lines[j].byte, "FastAPI route returns ORM-like object without response_model; prefer a Pydantic response model", out)
				break
			}
		}
	}
}

// BP-PY-32: FileResponse with user-influenced path.
func detectBPPY32(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-32")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "FileResponse") {
		return
	}
	// Prefer FastAPI/Starlette context; still flag FileResponse alone if path looks dynamic.
	fastapiish := looksFastAPIish(unit, src) || strings.Contains(src, "FileResponse")

	// Collect path-ish parameter names from route signatures near FileResponse.
	paramNames := collectRoutePathParams(src)

	lines := codeLinesFacts(facts, src)
	for _, line := range lines {
		if !strings.Contains(line.text, "FileResponse") {
			continue
		}
		// Find FileResponse(
		start := 0
		text := line.text
		for {
			idx := strings.Index(text[start:], "FileResponse")
			if idx < 0 {
				break
			}
			absInLine := start + idx
			// Ensure call form
			rest := text[absInLine+len("FileResponse"):]
			rest = strings.TrimLeft(rest, " \t")
			if !strings.HasPrefix(rest, "(") {
				start = absInLine + len("FileResponse")
				continue
			}
			// Work on full source for multi-line args.
			openAbs := line.byte + absInLine + len("FileResponse")
			// Adjust: absInLine points to FileResponse; open paren may have spaces.
			// Recompute open paren in source.
			searchFrom := line.byte + absInLine
			openParen := strings.Index(src[searchFrom:], "(")
			if openParen < 0 {
				break
			}
			openAbs = searchFrom + openParen
			arg, _, ok := firstCallArg(src, openAbs)
			if !ok {
				start = absInLine + len("FileResponse")
				continue
			}
			// Named path= support: if first arg empty-ish, scan full args.
			inner, _, iok := callArgsRegion(src, openAbs)
			pathArg := arg
			if iok {
				if strings.Contains(inner, "path=") {
					pathArg = kwArgValue(inner, "path")
				}
			}
			if pathArg == "" {
				start = absInLine + len("FileResponse")
				continue
			}
			if isPlainStringLiteral(pathArg) {
				start = absInLine + len("FileResponse")
				continue
			}
			// Dynamic?
			if isDynamicPathArg(pathArg, paramNames) || (fastapiish && !isPlainStringLiteral(pathArg)) {
				// Require some dynamism signal for non-fastapiish; for fastapiish any non-literal is enough.
				if isDynamicPathArg(pathArg, paramNames) || isFStringArg(pathArg) ||
					strings.Contains(pathArg, "request.") || strings.Contains(pathArg, "+") ||
					strings.Contains(pathArg, ".format") {
					pushAt(unit, meta, line.byte+absInLine, "FileResponse path comes from user input; confine with realpath+prefix checks", out)
				} else if fastapiish && isSimpleIdent(strings.TrimSpace(pathArg)) {
					// Bare param name: FileResponse(name)
					pushAt(unit, meta, line.byte+absInLine, "FileResponse path comes from user input; confine with realpath+prefix checks", out)
				}
			}
			start = absInLine + len("FileResponse")
		}
	}
}

func collectRoutePathParams(src string) map[string]struct{} {
	out := map[string]struct{}{}
	// {name} in decorator path strings
	re := regexp.MustCompile(`\{([A-Za-z_][A-Za-z0-9_]*)\}`)
	for _, m := range re.FindAllStringSubmatch(src, -1) {
		if len(m) > 1 {
			out[m[1]] = struct{}{}
		}
	}
	return out
}

func isDynamicPathArg(arg string, params map[string]struct{}) bool {
	arg = strings.TrimSpace(arg)
	if isFStringArg(arg) {
		return true
	}
	if strings.Contains(arg, "request.") || strings.Contains(arg, ".format(") || strings.Contains(arg, " + ") {
		return true
	}
	if isSimpleIdent(arg) {
		if _, ok := params[arg]; ok {
			return true
		}
		// Common path param names
		switch arg {
		case "path", "filename", "name", "file", "filepath", "file_path", "file_name":
			return true
		}
	}
	return false
}

func kwArgValue(inner, key string) string {
	// Find key= at top level
	needle := key + "="
	start := 0
	depth := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'':
			inStr = c
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		}
		if depth == 0 && i+len(needle) <= len(inner) && inner[i:i+len(needle)] == needle {
			// Ensure not mid-ident
			if i > 0 && isIdentByte(inner[i-1]) {
				continue
			}
			valStart := i + len(needle)
			// Read until top-level comma
			for j := valStart; j < len(inner); j++ {
				c2 := inner[j]
				if inStr != 0 {
					if escape {
						escape = false
						continue
					}
					if c2 == '\\' {
						escape = true
						continue
					}
					if c2 == inStr {
						inStr = 0
					}
					continue
				}
				if c2 == '"' || c2 == '\'' {
					inStr = c2
					continue
				}
				if c2 == '(' || c2 == '[' || c2 == '{' {
					depth++
					continue
				}
				if c2 == ')' || c2 == ']' || c2 == '}' {
					if depth > 0 {
						depth--
					}
					continue
				}
				if c2 == ',' && depth == 0 {
					return strings.TrimSpace(inner[valStart:j])
				}
			}
			return strings.TrimSpace(inner[valStart:])
		}
		_ = start
	}
	return ""
}

// isPlainStringLiteral is a non-f-string Python string/bytes literal.
func isPlainStringLiteral(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if isFStringArg(s) {
		return false
	}
	return isStringLiteral(s)
}

// isFStringArg reports whether s is (or starts as) an f-string / F-string literal.
func isFStringArg(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// f"..." / f'...' / fr"..." / rf"..." / F"""...
	// Strip b/u/r that may combine with f: rf, fr, etc.
	lower := strings.ToLower(s)
	// Find if 'f' appears in prefix before quote.
	i := 0
	hasF := false
	for i < len(s) {
		c := lower[i]
		if c == 'f' {
			hasF = true
			i++
			continue
		}
		if c == 'r' || c == 'b' || c == 'u' {
			i++
			continue
		}
		break
	}
	if !hasF || i >= len(s) {
		return false
	}
	return s[i] == '"' || s[i] == '\''
}

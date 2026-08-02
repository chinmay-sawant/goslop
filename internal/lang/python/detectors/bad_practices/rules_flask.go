package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	metaByID["BP-PY-18"] = &rules.RuleMetadata{
		ID: "BP-PY-18", Title: "Flask Route Missing Methods Restriction",
		Description: "A mutating endpoint accepts all methods because `@app.route` omits `methods=`.",
		Severity:    rules.SeverityLow, Pack: rules.PackBadPractice,
		Fix: "Pass methods=['POST'] (or the methods you intend) on the route decorator.",
	}
	metaByID["BP-PY-19"] = &rules.RuleMetadata{
		ID: "BP-PY-19", Title: "Flask jsonify Error Leaks Exception",
		Description: "Error handlers return str(e) or traceback details to clients.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Return a generic error message to clients; log the exception server-side.",
	}
	metaByID["BP-PY-20"] = &rules.RuleMetadata{
		ID: "BP-PY-20", Title: "Flask send_file User Path",
		Description: "`send_file` / `send_from_directory` is fed a path derived from request input without a jail.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Serve only under a fixed root via send_from_directory + safe_join; never pass raw request paths.",
	}
	RegisterRule("BP-PY-18", detectBPPY18)
	RegisterRule("BP-PY-19", detectBPPY19)
	RegisterRule("BP-PY-20", detectBPPY20)
}

var (
	// @app.route( / @bp.route( / @blueprint.route(
	flaskRouteDecoRe = regexp.MustCompile(`@\w+\.route\s*\(`)
	// POST-like body access on Flask request
	flaskPostishRe = regexp.MustCompile(`\brequest\.(form|get_json|json|files)\b`)
	// errorhandler decorator or register_error_handler(
	flaskErrorHandlerRe = regexp.MustCompile(`(@\w+\.errorhandler\s*\(|\.register_error_handler\s*\()`)
	// Exception text leaked to client
	excLeakRe = regexp.MustCompile(`\b(str|repr)\s*\(\s*(e|exc|err|error|exception)\s*\)|traceback\.format_exc\s*\(`)
	// send_file / send_from_directory calls
	sendFileCallRe = regexp.MustCompile(`\b(send_file|send_from_directory)\s*\(`)
	// request-derived path fragments
	requestPathRe = regexp.MustCompile(`\brequest\.(args|form|files|view_args|json)\b`)
)

// BP-PY-18: route handlers using request.form/get_json/files without methods= restriction.
// v0 same-function window: decorators immediately above def, body until next peer def/class.
func detectBPPY18(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-18")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.hasAny("request.form", "request.get_json", "request.json", "request.files", ".route(") &&
		!strings.Contains(unit.Source, ".route(") {
		return
	}
	// Need both route and postish signals somewhere.
	if !strings.Contains(unit.Source, ".route(") {
		return
	}
	if !flaskPostishRe.MatchString(unit.Source) {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i].text)
		if !flaskRouteDecoRe.MatchString(t) && !strings.Contains(t, ".route(") {
			// Also accept add_url_rule for completeness? plan mentions it; skip if no methods still harder.
			continue
		}
		// Collect decorator block (possibly multi-line route call + other decorators).
		decoStart := i
		decoText := t
		// If route call is multi-line, accumulate until matching close on decorator lines.
		if strings.Contains(t, ".route(") && !routeDecoComplete(t) {
			for j := i + 1; j < len(lines) && j-i < 20; j++ {
				decoText += " " + strings.TrimSpace(lines[j].text)
				i = j
				if routeDecoComplete(decoText) {
					break
				}
			}
		}
		// methods= present → miss for this route.
		if strings.Contains(decoText, "methods=") || strings.Contains(decoText, "methods =") {
			continue
		}
		// Skip remaining stacked decorators to find def.
		defIdx := i + 1
		for defIdx < len(lines) {
			dt := strings.TrimSpace(lines[defIdx].text)
			if dt == "" {
				defIdx++
				continue
			}
			if strings.HasPrefix(dt, "@") {
				defIdx++
				continue
			}
			break
		}
		if defIdx >= len(lines) {
			continue
		}
		defLine := strings.TrimSpace(lines[defIdx].text)
		if !strings.HasPrefix(defLine, "def ") && !strings.HasPrefix(defLine, "async def ") {
			continue
		}
		defIndent := indentWidth(lines[defIdx].raw)
		// Scan body for POST-like request access.
		bodyHasPostish := false
		for k := defIdx + 1; k < len(lines); k++ {
			raw := lines[k].raw
			bt := strings.TrimSpace(lines[k].text)
			if bt == "" {
				continue
			}
			ind := indentWidth(raw)
			if ind <= defIndent {
				// peer def/class/decorator ends body
				break
			}
			if flaskPostishRe.MatchString(bt) {
				bodyHasPostish = true
				break
			}
		}
		if bodyHasPostish {
			pushAt(unit, meta, lines[decoStart].byte,
				"Flask route uses request form/json/files without methods=; restrict allowed HTTP methods", out)
		}
	}
}

func routeDecoComplete(s string) bool {
	// Rough: decorator line(s) that opened .route( should balance parens.
	open := strings.Index(s, ".route(")
	if open < 0 {
		return true
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for i := open + len(".route"); i < len(s); i++ {
		c := s[i]
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
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return true
			}
		}
	}
	return false
}

// BP-PY-19: errorhandler / register_error_handler body returns str(e) / traceback to client.
func detectBPPY19(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-19")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.hasAny("errorhandler", "register_error_handler", "traceback") &&
		!strings.Contains(unit.Source, "errorhandler") &&
		!strings.Contains(unit.Source, "register_error_handler") {
		return
	}
	if !flaskErrorHandlerRe.MatchString(unit.Source) {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i].text)
		isHandler := false
		// @app.errorhandler(...)
		if strings.Contains(t, ".errorhandler(") || strings.HasPrefix(t, "@") && strings.Contains(t, "errorhandler") {
			isHandler = true
		}
		// app.register_error_handler(...)
		if strings.Contains(t, "register_error_handler(") {
			isHandler = true
		}
		if !isHandler {
			continue
		}
		// Find following def (for decorator form) or skip if register_error_handler with lambda — v0: def only.
		defIdx := i + 1
		if strings.Contains(t, "register_error_handler(") {
			// May be register_error_handler(Exception, handle) — check same line and later handle body by name is hard.
			// v0: only decorator-based errorhandler + same-function window after def.
			// If this line itself contains str(e) etc, flag.
			if excLeakRe.MatchString(t) {
				pushAt(unit, meta, lines[i].byte, "error handler returns exception text to clients; use a generic message", out)
			}
			// Also if next lines form a def used inline is rare; continue to look for nearby def only when decorator.
		}
		// Decorator path: walk to def.
		if !isFlaskErrorHandlerDecorator(t) {
			continue
		}
		for defIdx < len(lines) {
			dt := strings.TrimSpace(lines[defIdx].text)
			if dt == "" || strings.HasPrefix(dt, "@") {
				defIdx++
				continue
			}
			break
		}
		if defIdx >= len(lines) {
			continue
		}
		defLine := strings.TrimSpace(lines[defIdx].text)
		if !strings.HasPrefix(defLine, "def ") && !strings.HasPrefix(defLine, "async def ") {
			continue
		}
		defIndent := indentWidth(lines[defIdx].raw)
		for k := defIdx + 1; k < len(lines); k++ {
			bt := strings.TrimSpace(lines[k].text)
			if bt == "" {
				continue
			}
			ind := indentWidth(lines[k].raw)
			if ind <= defIndent {
				break
			}
			if exceptionTextIsReturned(bt) {
				pushAt(unit, meta, lines[k].byte, "error handler returns exception text to clients; use a generic message", out)
				break
			}
		}
	}
}

func exceptionTextIsReturned(line string) bool {
	if !excLeakRe.MatchString(line) {
		return false
	}
	return strings.Contains(line, "return") || strings.Contains(line, "jsonify") ||
		strings.Contains(line, "make_response") || strings.Contains(line, "format_exc")
}

// BP-PY-20: send_file / send_from_directory with request-derived path without a jail.
// Idiomatic send_from_directory(fixed_root, request_name) is safe (Flask jails via
// safe_join internally). Flag send_file(request_path) and send_from_directory when
// the directory/root argument itself is request-derived.
func detectBPPY20(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-20")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.hasAny("send_file", "send_from_directory") &&
		!strings.Contains(unit.Source, "send_file") &&
		!strings.Contains(unit.Source, "send_from_directory") {
		return
	}
	src := unit.Source
	for _, m := range sendFileCallRe.FindAllStringIndex(src, -1) {
		open := strings.IndexByte(src[m[0]:m[1]], '(')
		if open < 0 {
			continue
		}
		openAbs := m[0] + open
		inner, ok := callArgsRegion(src, openAbs)
		if !ok {
			windowEnd := m[1] + flaskSendFileFallbackBytes
			if windowEnd > len(src) {
				windowEnd = len(src)
			}
			inner = src[m[0]:windowEnd]
		}
		if strings.Contains(inner, "safe_join") {
			continue
		}
		callName := src[m[0]:m[1]]
		isFromDir := strings.Contains(callName, "send_from_directory")
		first, argOK := firstCallArg(src, openAbs)
		if !argOK {
			first = ""
		}
		first = strings.TrimSpace(first)

		if isFromDir {
			// Fixed root + request filename is the recommended Flask pattern → miss.
			if first != "" && (isStringLiteral(first) || isSimpleIdent(first) && !nameAssignedFromRequest(src, m[0], first) && !requestPathRe.MatchString(first)) {
				continue
			}
			// Request-derived directory/root → hit.
			if requestPathRe.MatchString(first) || (isSimpleIdent(first) && nameAssignedFromRequest(src, m[0], first)) {
				pushAt(unit, meta, m[0], "send_from_directory root/directory comes from request; path traversal risk", out)
			}
			continue
		}

		// send_file: hit when the path argument is request-derived.
		if requestPathRe.MatchString(inner) || requestPathRe.MatchString(first) {
			pushAt(unit, meta, m[0], "send_file path from request without safe_join; path traversal risk", out)
			continue
		}
		if first == "" || isStringLiteral(first) {
			continue
		}
		if isSimpleIdent(first) && nameAssignedFromRequest(src, m[0], first) {
			pushAt(unit, meta, m[0], "send_file path from request without safe_join; path traversal risk", out)
		}
	}
}

func nameAssignedFromRequest(src string, before int, name string) bool {
	if before < 0 {
		before = 0
	}
	// Window: last ~1500 bytes before the call.
	start := before - flaskRequestAssignmentBytes
	if start < 0 {
		start = 0
	}
	window := src[start:before]
	// name = request.args... or name = request.form...
	// Also name = request.args.get(...)
	re := regexp.MustCompile(`(?m)\b` + regexp.QuoteMeta(name) + `\s*=\s*request\.(args|form|files|view_args|json)\b`)
	return re.MatchString(window)
}

const (
	flaskSendFileFallbackBytes  = 200
	flaskRequestAssignmentBytes = 1500
)

func isFlaskErrorHandlerDecorator(line string) bool {
	return strings.HasPrefix(line, "@") || strings.Contains(line, ".errorhandler(")
}

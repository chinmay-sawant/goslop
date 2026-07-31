package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// No SourceIndex gates: each detector has several valid call or assignment
	// shapes, so a narrower prefilter could silently skip a true positive.
	RegisterRule("CWE-201", detectCWE201, &MetaCWE201)
	RegisterRule("CWE-204", detectCWE204, &MetaCWE204)
	RegisterRule("CWE-208", detectCWE208, &MetaCWE208)
	RegisterRule("CWE-209", detectCWE209, &MetaCWE209)
	RegisterRule("CWE-212", detectCWE212, &MetaCWE212)
	RegisterRule("CWE-213", detectCWE213, &MetaCWE213)
	RegisterRule("CWE-488", detectCWE488, &MetaCWE488)
	RegisterRule("CWE-497", detectCWE497, &MetaCWE497)
}

var (
	pySensitiveDataAccessRE = regexp.MustCompile(`(?i)(?:\.[ \t]*(?:password|passwd|token|secret|api[_-]?key|credential|ssn|social_security|pan|card_number|salary)\b|\[[ \t]*['"](?:password|passwd|token|secret|api[_-]?key|credential|ssn|social_security|pan|card_number|salary)['"][ \t]*\])`)
	pyTimingCompareRE       = regexp.MustCompile(`(?im)\b[a-z_]*(?:password|passwd|token|secret|digest|signature|api[_-]?key|credential|auth)[a-z0-9_]*\s*==\s*(?:[a-z_][a-z0-9_]*\.)?[a-z_]*(?:password|passwd|token|secret|digest|signature|api[_-]?key|credential|auth)[a-z0-9_]*\b|\b(?:[a-z_][a-z0-9_]*\.)?[a-z_]*(?:password|passwd|token|secret|digest|signature|api[_-]?key|credential|auth)[a-z0-9_]*\s*==\s*[a-z_]*(?:password|passwd|token|secret|digest|signature|api[_-]?key|credential|auth)[a-z0-9_]*\b`)
	pyExceptionStringRE     = regexp.MustCompile(`(?i)str\s*\(\s*(?:e|err|exc|error|exception)\b`)
	pyGlobalSessionStateRE  = regexp.MustCompile(`(?im)^\s*global\s+(?:current_user|current_session|active_user|session_user)\b`)
	pySessionStateAssignRE  = regexp.MustCompile(`(?im)^\s*(?:current_user|current_session|active_user|session_user)\s*=`)
)

func detectCWE201(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range responseCalls(unit.Source) {
		if pySensitiveDataAccessRE.MatchString(call.ArgsText) {
			emitInfoExposure(unit, &MetaCWE201, call.Start, "HTTP response directly includes a sensitive data field", confidence82, out)
			return
		}
	}
}

// CWE-204 is restricted to route handlers and the specific pair of account
// existence/password failure messages. Generic validation errors and internal
// helpers are deliberately out of scope for this same-file heuristic.
func detectCWE204(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, handler := range routeHandlerBodies(unit.Source) {
		missingAccount, wrongPassword, start := authResponseDiscrepancy(handler.body)
		if missingAccount && wrongPassword {
			emitInfoExposure(unit, &MetaCWE204, handler.start+start,
				"route reveals whether an account exists through distinct authentication errors", confidence78, out)
			return
		}
	}
}

func authResponseDiscrepancy(source string) (bool, bool, int) {
	missingAccount := false
	wrongPassword := false
	start := -1
	for _, call := range responseCalls(source) {
		args := strings.ToLower(call.ArgsText)
		if strings.Contains(args, "user not found") || strings.Contains(args, "unknown user") || strings.Contains(args, "account not found") {
			missingAccount = true
			if start < 0 {
				start = call.Start
			}
		}
		if strings.Contains(args, "invalid password") || strings.Contains(args, "incorrect password") || strings.Contains(args, "wrong password") {
			wrongPassword = true
			if start < 0 {
				start = call.Start
			}
		}
	}
	return missingAccount, wrongPassword, start
}

// CWE-208 only reports equality comparisons where both operands are named as
// authentication-sensitive values. hmac.compare_digest and ordinary equality
// comparisons without two such identifiers are not reported.
func detectCWE208(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if start := firstCodeMatchStart(unit.Source, pyTimingCompareRE); start >= 0 {
		emitInfoExposure(unit, &MetaCWE208, start, "security-sensitive values are compared with == instead of a constant-time comparison", confidence82, out)
	}
}

func detectCWE209(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range responseCalls(unit.Source) {
		if pyExceptionStringRE.MatchString(call.ArgsText) || strings.Contains(call.ArgsText, "traceback.format_exc") {
			emitInfoExposure(unit, &MetaCWE209, call.Start, "HTTP response includes exception or traceback details", confidence84, out)
			return
		}
	}
}

func detectCWE212(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"json.dump", "yaml.dump", "csv.DictWriter.writerow"} {
		for _, call := range findCalls(unit.Source, name) {
			if pySensitiveDataAccessRE.MatchString(call.ArgsText) {
				emitInfoExposure(unit, &MetaCWE212, call.Start, "serialization stores a sensitive data field without a redacted export", confidence80, out)
				return
			}
		}
	}
}

func detectCWE213(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"requests.post", "requests.put", "requests.request", "httpx.post", "httpx.put", "httpx.request"} {
		for _, call := range findCalls(unit.Source, name) {
			args := strings.ToLower(call.ArgsText)
			if (strings.Contains(args, "/guest") || strings.Contains(args, "/public") || strings.Contains(args, "/anonymous")) && pySensitiveDataAccessRE.MatchString(call.ArgsText) {
				emitInfoExposure(unit, &MetaCWE213, call.Start, "guest or public request includes a sensitive data field", confidence78, out)
				return
			}
		}
	}
}

// CWE-488 requires a route handler to both declare and assign the global
// request-identity state. A module-global cache or a normal local variable is
// therefore not enough to emit a finding.
func detectCWE488(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, handler := range routeHandlerBodies(unit.Source) {
		code := pythonCodeMask(handler.body)
		global := pyGlobalSessionStateRE.FindStringIndex(code)
		assign := pySessionStateAssignRE.FindStringIndex(code)
		if global != nil && assign != nil {
			emitInfoExposure(unit, &MetaCWE488, handler.start+assign[0], "route stores request identity in module-global session state", confidence84, out)
			return
		}
	}
}

func detectCWE497(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range responseCalls(unit.Source) {
		args := call.ArgsText
		if strings.Contains(args, "os.environ") || strings.Contains(args, "sys.path") || strings.Contains(args, "platform.uname") || strings.Contains(args, "socket.gethostname") || strings.Contains(args, "traceback.format_exc") {
			emitInfoExposure(unit, &MetaCWE497, call.Start, "HTTP response exposes sensitive system diagnostic information", confidence84, out)
			return
		}
	}
}

func responseCalls(source string) []callSite {
	return findCalls(source, "jsonify", "flask.jsonify", "HttpResponse", "Response")
}

func emitInfoExposure(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}

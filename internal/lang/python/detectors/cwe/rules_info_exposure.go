package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// Most info-exposure detectors stay ungated: multi-shape call/assignment
	// evidence makes a narrow prefilter FN-prone. CWE-208 is gated on the
	// sensitive-identifier set that must appear for pyTimingCompareRE to fire.
	RegisterRule("CWE-201", detectCWE201, &MetaCWE201)
	RegisterRule("CWE-204", detectCWE204, &MetaCWE204)
	RegisterRule("CWE-208", detectCWE208, &MetaCWE208,
		"password", "passwd", "token", "secret", "digest", "signature",
		"api_key", "credential", "auth")
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

func detectCWE201(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range responseCalls(facts, unit.Source) {
		if pySensitiveDataAccessRE.MatchString(call.ArgsText) {
			emitInfoExposure(unit, &MetaCWE201, call.Start, "HTTP response directly includes a sensitive data field", confidence82, out)
			return
		}
	}
}

// CWE-204 is restricted to route handlers and the specific pair of account
// existence/password failure messages. Generic validation errors and internal
// helpers are deliberately out of scope for this same-file heuristic.
func detectCWE204(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, handler := range routeHandlerBodies(facts, unit.Source) {
		masked := facts.codeMask(handler.body, handler.start)
		missingAccount, wrongPassword, start := authResponseDiscrepancy(handler.body, masked)
		if missingAccount && wrongPassword {
			emitInfoExposure(unit, &MetaCWE204, handler.start+start,
				"route reveals whether an account exists through distinct authentication errors", confidence78, out)
			return
		}
	}
}

func authResponseDiscrepancy(source, masked string) (bool, bool, int) {
	missingAccount := false
	wrongPassword := false
	start := -1
	for _, call := range findCallsMasked(source, masked, "jsonify", "flask.jsonify", "HttpResponse", "Response") {
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
//
// Test modules are skipped: assert token/authority comparisons and fixture
// constant checks are not deployed authentication sinks.
// Offline release/tool scripts (scripts/release signature manifests) are also
// skipped — those == checks validate build artifacts, not live auth.
//
// Bare "auth" substring matches that only fire on authority/author enums
// (Project_Parva AuthorityTaint) are filtered: not credential material.
//
// Two-pass hot path: only lines containing "==" are checked with the timing
// compare RE (avoids FindAllStringIndex over the whole file).
func detectCWE208(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if isPythonTestModule(unit) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "==") {
		return
	}
	masked := facts.Masked
	if masked == "" {
		masked = pythonCodeMask(src)
	}
	offset := 0
	for _, line := range strings.Split(src, "\n") {
		if strings.Contains(line, "==") {
			if loc := pyTimingCompareRE.FindStringIndex(line); loc != nil {
				absStart := offset + loc[0]
				absEnd := offset + loc[1]
				match := line[loc[0]:loc[1]]
				if !timingCompareIsCredentialMaterial(match) {
					offset += len(line) + 1
					continue
				}
				if absEnd <= len(masked) && strings.TrimSpace(masked[absStart:absEnd]) != "" {
					emitInfoExposure(unit, &MetaCWE208, absStart,
						"security-sensitive values are compared with == instead of a constant-time comparison", confidence82, out)
					return
				}
			}
		}
		offset += len(line) + 1
	}
}

// timingCompareIsCredentialMaterial rejects RE hits driven only by
// authority/author identifiers (AuthorityTaint enum ranking), which are not
// secret values subject to timing attacks.
func timingCompareIsCredentialMaterial(match string) bool {
	parts := strings.SplitN(match, "==", 2)
	if len(parts) != 2 {
		return false
	}
	return timingOperandIsCredential(parts[0]) && timingOperandIsCredential(parts[1])
}

func timingOperandIsCredential(expr string) bool {
	// Prefer the rightmost identifier of an attribute chain
	// (resolution.authority → authority; candidate.password → password).
	ident := timingOperandLeaf(expr)
	if ident == "" {
		return false
	}
	lower := strings.ToLower(ident)
	// authority / author / authorize* are not secret material.
	if lower == "authority" || lower == "author" ||
		strings.HasPrefix(lower, "authorit") || strings.HasPrefix(lower, "authoriz") {
		return false
	}
	for _, stem := range []string{
		"password", "passwd", "token", "secret", "digest", "signature",
		"api_key", "apikey", "credential",
	} {
		if strings.Contains(lower, stem) {
			return true
		}
	}
	// Whole-identifier "auth" / "auth_*" only (not authority — filtered above).
	if lower == "auth" || strings.HasPrefix(lower, "auth_") {
		return true
	}
	return false
}

func timingOperandLeaf(expr string) string {
	s := strings.TrimSpace(expr)
	if s == "" {
		return ""
	}
	// Strip trailing non-ident noise.
	end := len(s)
	for end > 0 {
		c := s[end-1]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			break
		}
		end--
	}
	s = s[:end]
	if i := strings.LastIndexByte(s, '.'); i >= 0 {
		s = s[i+1:]
	}
	// Leading junk (operators) — take trailing identifier run.
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' {
			start = i
			break
		}
	}
	// Find end of identifier from start.
	j := start
	for j < len(s) {
		c := s[j]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			j++
			continue
		}
		break
	}
	// Prefer last identifier segment if multiple words.
	if j < len(s) {
		// walk from right for last ident
		for i := len(s) - 1; i >= 0; i-- {
			c := s[i]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
				end := i + 1
				k := i
				for k >= 0 {
					c2 := s[k]
					if (c2 >= 'a' && c2 <= 'z') || (c2 >= 'A' && c2 <= 'Z') || (c2 >= '0' && c2 <= '9') || c2 == '_' {
						k--
						continue
					}
					break
				}
				return s[k+1 : end]
			}
		}
	}
	if start < j {
		return s[start:j]
	}
	return s
}

func detectCWE209(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range responseCalls(facts, unit.Source) {
		if pyExceptionStringRE.MatchString(call.ArgsText) || strings.Contains(call.ArgsText, "traceback.format_exc") {
			emitInfoExposure(unit, &MetaCWE209, call.Start, "HTTP response includes exception or traceback details", confidence84, out)
			return
		}
	}
}

func detectCWE212(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"json.dump", "yaml.dump", "csv.DictWriter.writerow"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if pySensitiveDataAccessRE.MatchString(call.ArgsText) {
				emitInfoExposure(unit, &MetaCWE212, call.Start, "serialization stores a sensitive data field without a redacted export", confidence80, out)
				return
			}
		}
	}
}

func detectCWE213(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"requests.post", "requests.put", "requests.request", "httpx.post", "httpx.put", "httpx.request"} {
		for _, call := range findCalls(facts, unit.Source, name) {
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
func detectCWE488(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, handler := range routeHandlerBodies(facts, unit.Source) {
		code := facts.codeMask(handler.body, handler.start)
		global := pyGlobalSessionStateRE.FindStringIndex(code)
		assign := pySessionStateAssignRE.FindStringIndex(code)
		if global != nil && assign != nil {
			emitInfoExposure(unit, &MetaCWE488, handler.start+assign[0], "route stores request identity in module-global session state", confidence84, out)
			return
		}
	}
}

func detectCWE497(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range responseCalls(facts, unit.Source) {
		args := call.ArgsText
		if strings.Contains(args, "os.environ") || strings.Contains(args, "sys.path") || strings.Contains(args, "platform.uname") || strings.Contains(args, "socket.gethostname") || strings.Contains(args, "traceback.format_exc") {
			emitInfoExposure(unit, &MetaCWE497, call.Start, "HTTP response exposes sensitive system diagnostic information", confidence84, out)
			return
		}
	}
}

func responseCalls(facts *PyCweFacts, source string) []callSite {
	return findCalls(facts, source, "jsonify", "flask.jsonify", "HttpResponse", "Response")
}

func emitInfoExposure(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}

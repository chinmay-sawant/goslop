package cwe

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-66", detectCWE66, &MetaCWE66)
	RegisterRule("CWE-76", detectCWE76, &MetaCWE76)
	RegisterRule("CWE-178", detectCWE178, &MetaCWE178)
	RegisterRule("CWE-179", detectCWE179, &MetaCWE179)
	RegisterRule("CWE-182", detectCWE182, &MetaCWE182)
	RegisterRule("CWE-184", detectCWE184, &MetaCWE184)
	RegisterRule("CWE-186", detectCWE186, &MetaCWE186)
	RegisterRule("CWE-257", detectCWE257, &MetaCWE257)
	RegisterRule("CWE-272", detectCWE272, &MetaCWE272)
	RegisterRule("CWE-279", detectCWE279, &MetaCWE279)
	RegisterRule("CWE-289", detectCWE289, &MetaCWE289)
	RegisterRule("CWE-290", detectCWE290, &MetaCWE290)
	RegisterRule("CWE-323", detectCWE323, &MetaCWE323)
	RegisterRule("CWE-331", detectCWE331, &MetaCWE331)
	RegisterRule("CWE-334", detectCWE334, &MetaCWE334)
}

var (
	pyTierBVirtualNameRE = regexp.MustCompile(`(?i)\b(?:open|Path)\s*\(\s*["'](?:CON|NUL|COM1)(?:\.[^"']*)?["']`)
	pyTierBStripTagRE    = regexp.MustCompile(`(?is)\.replace\s*\(\s*["']<["']\s*,\s*["']["']\s*\)`)
	pyTierBEarlyDecodeRE = regexp.MustCompile(`(?is)(?:validate_[a-z_]+|validate\s*\()[^\n]*\n[^\n]*(?:unquote|unescape|decode)\s*\(`)
	pyTierBCollapseRE    = regexp.MustCompile(`(?is)re\.sub\s*\([^\n]*request\.[^\n]*\)[\s\S]{0,300}(?:is_admin|authorize|check_permission)\s*\(`)
	pyTierBDenyLoopRE    = regexp.MustCompile(`(?is)for\s+\w+\s+in\s+(?:bad|deny|blocked)[a-z_]*\s*:[\s\S]{0,180}if\s+\w+\s+in\s+request\.`)
	pyTierBWeakLengthRE  = regexp.MustCompile(`(?im)^\s*(?:PASSWORD_)?MIN(?:IMUM)?_?LENGTH\s*=\s*[1-5]\b`)
	pyTierBNonceRE       = regexp.MustCompile(`(?is)nonce\s*=\s*b?["'][^"']*["']\s*(?:\*\s*12)?[\s\S]{0,220}AESGCM\s*\([^)]*\)\.encrypt\s*\(\s*nonce\b`)
)

func detectCWE66(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBVirtualNameRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE66, start, "Windows virtual resource name is used as a file path", 0.78, out)
	}
}

func detectCWE76(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBStripTagRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE76, start, "manual removal of one markup delimiter is used as neutralization", 0.76, out)
	}
}

func detectCWE178(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range pythonFunctions(unit.Source) {
		code := pythonCodeMask(fn.body)
		if strings.Contains(code, "request.args.get(") && strings.Contains(code, "username") && strings.Contains(code, "==") && !strings.Contains(code, "casefold(") && !strings.Contains(code, ".lower(") {
			emitTierBFinding(unit, &MetaCWE178, fn.bodyStart, "request username is compared without case normalization", 0.7, out)
			return
		}
	}
}

func detectCWE179(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBEarlyDecodeRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE179, start, "input is validated before a later decoding transformation", 0.78, out)
	}
}

func detectCWE182(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBCollapseRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE182, start, "request data is collapsed before a security decision", 0.76, out)
	}
}

func detectCWE184(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBDenyLoopRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE184, start, "request input is protected by a manual substring deny-list", 0.76, out)
	}
}

func detectCWE186(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "re.compile") {
		if strings.Contains(call.ArgsText, "{3}") || strings.Contains(call.ArgsText, "{2}") {
			emitTierBFinding(unit, &MetaCWE186, call.Start, "validation regular expression accepts only a fixed narrow identifier shape", 0.7, out)
			return
		}
	}
}

func detectCWE257(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, ".encrypt") {
		if containsPasswordName(call.ArgsText) {
			emitTierBFinding(unit, &MetaCWE257, call.Start, "password is encrypted into a recoverable representation", 0.84, out)
			return
		}
	}
}

func detectCWE272(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "os.setuid", "os.seteuid", "os.setgid") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) > 0 && strings.TrimSpace(args[0]) == "0" {
			emitTierBFinding(unit, &MetaCWE272, call.Start, "process explicitly assumes root privileges", 0.84, out)
			return
		}
	}
}

func detectCWE279(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "os.chmod") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 2 && (strings.TrimSpace(args[1]) == "0o777" || strings.TrimSpace(args[1]) == "0o666") {
			emitTierBFinding(unit, &MetaCWE279, call.Start, "runtime permission change makes a resource world accessible", 0.84, out)
			return
		}
	}
}

func detectCWE289(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range pythonFunctions(unit.Source) {
		code := pythonCodeMask(fn.body)
		if strings.Contains(code, "request.headers.get(") && strings.Contains(strings.ToLower(fn.body), "host") && strings.Contains(code, "==") && !strings.Contains(code, ".lower(") && !strings.Contains(code, "casefold(") {
			emitTierBFinding(unit, &MetaCWE289, fn.bodyStart, "host identity is compared without normalization", 0.72, out)
			return
		}
	}
}

func detectCWE290(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "request.headers.get") {
		if strings.Contains(strings.ToLower(call.ArgsText), "x-forwarded-for") {
			emitTierBFinding(unit, &MetaCWE290, call.Start, "client-provided X-Forwarded-For header is trusted directly", 0.82, out)
			return
		}
	}
}

func detectCWE323(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBNonceRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE323, start, "fixed nonce is reused for AES-GCM encryption", 0.9, out)
	}
}

func detectCWE331(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "secrets.token_bytes", "secrets.token_urlsafe") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) > 0 && shortPositiveLiteral(args[0], 15) {
			emitTierBFinding(unit, &MetaCWE331, call.Start, "security token is generated with fewer than 16 bytes of entropy", 0.82, out)
			return
		}
	}
}

func detectCWE334(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "random.randint") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 2 && strings.TrimSpace(args[0]) == "0" && shortPositiveLiteral(args[1], 9999) {
			emitTierBFinding(unit, &MetaCWE334, call.Start, "random value for an authentication-sized token has a small numeric space", 0.8, out)
			return
		}
	}
}

func containsPasswordName(value string) bool {
	lower := strings.ToLower(value)
	return containsIdent(lower, "password") || containsIdent(lower, "passwd") || containsIdent(lower, "pwd")
}

func shortPositiveLiteral(value string, max int) bool {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return false
	}
	return n > 0 && n <= max
}

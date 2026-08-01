package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-798", detectCWE798, &MetaCWE798,
		"password", "PASSWORD", "passwd", "api_key", "API_KEY", "SECRET_KEY",
		"aws_secret_access_key", "access_token", "auth_token", "private_key", "secret")
	RegisterRule("CWE-256", detectCWE256, &MetaCWE256,
		"password", "PASSWORD", "passwd", "pwd")
	RegisterRule("CWE-260", detectCWE260, &MetaCWE260,
		"password", "PASSWORD", "config[", "DATABASES")
	RegisterRule("CWE-261", detectCWE261, &MetaCWE261,
		"base64.b64encode", "binascii.hexlify", "codecs.encode")
	RegisterRule("CWE-312", detectCWE312, &MetaCWE312,
		"api_key", "API_KEY", "SECRET_KEY", "aws_secret_access_key",
		"access_token", "auth_token", "private_key", "secret")
	RegisterRule("CWE-319", detectCWE319, &MetaCWE319,
		"http://", "ftplib.FTP", "smtplib.SMTP", "requests.", "urllib.request.urlopen")
	RegisterRule("CWE-547", detectCWE547, &MetaCWE547,
		"SECURE_SSL_REDIRECT", "SESSION_COOKIE_SECURE", "CSRF_COOKIE_SECURE",
		"SESSION_COOKIE_HTTPONLY", "SECURE_HSTS_SECONDS", "ALLOWED_HOSTS")
	RegisterRule("CWE-523", detectCWE523, &MetaCWE523,
		"verify=False", "CERT_NONE", "check_hostname", "_create_unverified_context", "requests.")
}

var (
	pyHardcodedCredentialRE = regexp.MustCompile(`(?im)^\s*(?:[a-z_]*password|[a-z_]*passwd|[a-z_]*api[_-]?key|[a-z_]*secret(?:_key)?|[a-z_]*access_token|[a-z_]*auth_token|aws_secret_access_key|private_key)\s*=\s*(?:'[^'\r\n]{3,}'|"[^"\r\n]{3,}")`)
	pyPlaintextPasswordRE   = regexp.MustCompile(`(?im)^\s*(?:[a-z_]*password|[a-z_]*passwd|[a-z_]*pwd)\s*=\s*(?:'[^'\r\n]+'|"[^"\r\n]+")`)
	pyConfigPasswordRE      = regexp.MustCompile(`(?im)(?:['"]password['"]\s*:\s*|config\s*\[\s*['"]password['"]\s*\]\s*=\s*)(?:'[^'\r\n]+'|"[^"\r\n]+")`)
	pyCleartextSecretRE     = regexp.MustCompile(`(?im)^\s*(?:[a-z_]*api[_-]?key|[a-z_]*secret(?:_key)?|aws_secret_access_key|[a-z_]*access_token|[a-z_]*auth_token|private_key)\s*=\s*(?:'[^'\r\n]{3,}'|"[^"\r\n]{3,}")`)
	pyWeakSecuritySettingRE = regexp.MustCompile(`(?im)^\s*(?:SECURE_SSL_REDIRECT|SESSION_COOKIE_SECURE|CSRF_COOKIE_SECURE|SESSION_COOKIE_HTTPONLY)\s*=\s*False\b|^\s*SECURE_HSTS_SECONDS\s*=\s*0\b|^\s*ALLOWED_HOSTS\s*=\s*\[\s*['"]\*['"]\s*\]`)
	pyCheckHostnameFalseRE  = regexp.MustCompile(`(?m)check_hostname\s*=\s*False\b`)
)

// CWE-798 reports direct source literals assigned to credential-shaped names.
// Environment or secret-manager lookups do not match, keeping the heuristic
// focused on credentials that are actually committed to Python source.
func detectCWE798(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if isPythonTestModule(unit) || isPythonBenchmarkFile(unit) {
		return
	}
	if start := firstMatchStart(facts, unit, pyHardcodedCredentialRE); start >= 0 {
		pushSecretFinding(unit, &MetaCWE798, start, "credential is hard-coded in Python source", confidence82, out)
	}
}

// CWE-256 is deliberately narrower than CWE-798: it recognizes direct
// password assignments, not tokens and API keys covered by the other rules.
func detectCWE256(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(facts, unit, pyPlaintextPasswordRE); start >= 0 {
		pushSecretFinding(unit, &MetaCWE256, start, "password is stored as a plaintext source literal", confidence80, out)
	}
}

// CWE-260 recognizes literal password values in common Python configuration
// maps. It intentionally excludes os.environ and secret-provider lookups.
func detectCWE260(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(facts, unit, pyConfigPasswordRE); start >= 0 {
		pushSecretFinding(unit, &MetaCWE260, start, "configuration contains a literal password", confidence80, out)
	}
}

func detectCWE261(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"base64.b64encode", "binascii.hexlify", "codecs.encode"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if passwordArgument(call.ArgsText) {
				pushSecretFinding(unit, &MetaCWE261, call.Start, "password is protected only with a reversible encoding", confidence80, out)
				return
			}
		}
	}
}

func detectCWE312(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if isPythonTestModule(unit) || isPythonBenchmarkFile(unit) {
		return
	}
	if start := firstMatchStart(facts, unit, pyCleartextSecretRE); start >= 0 {
		pushSecretFinding(unit, &MetaCWE312, start, "sensitive key or token is stored as a cleartext source literal", confidence80, out)
	}
}

// CWE-319 requires both a cleartext transport and credential evidence for
// HTTP clients. Plain FTP and SMTP are also reported because their protocol
// transport is cleartext unless the caller explicitly upgrades it with TLS.
func detectCWE319(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"requests.get", "requests.post", "requests.put", "requests.request", "urllib.request.urlopen"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			args := strings.ToLower(call.ArgsText)
			if strings.Contains(args, "http://") && hasSensitiveTransportValue(args) {
				pushSecretFinding(unit, &MetaCWE319, call.Start, "credentials are sent through a cleartext HTTP request", confidence82, out)
				return
			}
		}
	}
	if calls := findCalls(facts, unit.Source, "ftplib.FTP"); len(calls) > 0 {
		pushSecretFinding(unit, &MetaCWE319, calls[0].Start, "FTP transports credentials without TLS protection", confidence76, out)
		return
	}
	if !strings.Contains(facts.Masked, ".starttls(") {
		if calls := findCalls(facts, unit.Source, "smtplib.SMTP"); len(calls) > 0 {
			pushSecretFinding(unit, &MetaCWE319, calls[0].Start, "SMTP transport is used without a same-file STARTTLS upgrade", confidence72, out)
		}
	}
}

func detectCWE547(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(facts, unit, pyWeakSecuritySettingRE); start >= 0 {
		pushSecretFinding(unit, &MetaCWE547, start, "security-relevant setting is hard-coded to an insecure value", confidence80, out)
	}
}

func detectCWE523(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"requests.get", "requests.post", "requests.put", "requests.request"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if hasKwargFalse(call.ArgsText, "verify") {
				pushSecretFinding(unit, &MetaCWE523, call.Start, "TLS certificate verification is disabled for a network request", confidence84, out)
				return
			}
		}
	}
	masked := facts.Masked
	for _, marker := range []string{"ssl._create_unverified_context", "ssl.CERT_NONE"} {
		if start := strings.Index(masked, marker); start >= 0 {
			pushSecretFinding(unit, &MetaCWE523, start, "TLS certificate verification is explicitly disabled", confidence82, out)
			return
		}
	}
	if start := pyCheckHostnameFalseRE.FindStringIndex(masked); start != nil {
		pushSecretFinding(unit, &MetaCWE523, start[0], "TLS hostname verification is explicitly disabled", confidence82, out)
	}
}

func firstMatchStart(facts *PyCweFacts, unit *core.ParsedUnit, pattern *regexp.Regexp) int {
	if unit == nil || pattern == nil {
		return -1
	}
	return firstCodeMatchStart(facts, unit.Source, pattern)
}

func passwordArgument(args string) bool {
	lower := strings.ToLower(args)
	return containsIdent(lower, "password") || containsIdent(lower, "passwd") || containsIdent(lower, "pwd")
}

func hasSensitiveTransportValue(args string) bool {
	return strings.Contains(args, "auth=") || strings.Contains(args, "password") || strings.Contains(args, "passwd") ||
		strings.Contains(args, "token=") || strings.Contains(args, "headers=") || strings.Contains(args, "params=")
}

func pushSecretFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, start int, message string, confidence float32, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	line, col := unit.LineCol(start)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}

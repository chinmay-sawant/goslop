package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// Argument-value evidence stays FN-safe when gated on the call/API tokens
	// that must appear before a finding can fire.
	RegisterRule("CWE-295", detectCWE295, &MetaCWE295,
		"verify=False", "CERT_NONE", "_create_unverified_context", "check_hostname",
		"requests.", "httpx.")
	RegisterRule("CWE-328", detectCWE328, &MetaCWE328,
		"hashlib.md5", "hashlib.sha1", "hashlib.new", "Crypto.Hash.MD5", "Crypto.Hash.SHA1")
	RegisterRule("CWE-335", detectCWE335, &MetaCWE335,
		"random.seed", "numpy.random.seed", "np.random.seed")
	RegisterRule("CWE-338", detectCWE338, &MetaCWE338,
		"random.", "numpy.random", "np.random")
	RegisterRule("CWE-347", detectCWE347, &MetaCWE347,
		"jwt.decode", "verify=False", "verify_signature")
	RegisterRule("CWE-1204", detectCWE1204, &MetaCWE1204,
		"AES.new")
	RegisterRule("CWE-1240", detectCWE1240, &MetaCWE1240,
		"xor")
	RegisterRule("CWE-1241", detectCWE1241, &MetaCWE1241,
		"random.", "numpy.random", "np.random")
	RegisterRule("CWE-1392", detectCWE1392, &MetaCWE1392,
		"password", "passwd", "pwd")
}

var (
	pyFixedSeedRE         = regexp.MustCompile(`(?s)^\s*(?:[-+]?\d+(?:\.\d+)?|True|False|None|[brufBRUF]*['"][^'"]*['"])\s*$`)
	pyFixedIVRE           = regexp.MustCompile(`(?s)^[brBR]*['"][^'"]*['"]\s*(?:\*\s*\d+)?\s*$`)
	pyDefaultCredentialRE = regexp.MustCompile(`(?im)^\s*(?:[a-z_]*password|[a-z_]*passwd|[a-z_]*pwd)\s*=\s*[bruf]*['"](?:admin|password|changeme|default|guest|root|toor)['"]\s*$`)
	pyRiskyCryptoNameRE   = regexp.MustCompile(`^(?:xor(?:_cipher|_encrypt|_decrypt)|(?:custom|homegrown)_(?:encrypt|decrypt|cipher))$`)
)

// CWE-295 covers explicit certificate or hostname validation bypasses. Calls
// are located through findCalls, which masks comments and docstrings first.
func detectCWE295(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"requests.get", "requests.post", "requests.put", "requests.request", "httpx.get", "httpx.post", "httpx.request"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if hasKwargFalse(call.ArgsText, "verify") {
				emitCryptoFinding(unit, &MetaCWE295, call.Start, "TLS certificate verification is disabled for a network request", confidence86, out)
				return
			}
		}
	}
	masked := facts.Masked
	for _, marker := range []string{"ssl._create_unverified_context", "ssl.CERT_NONE"} {
		if start := strings.Index(masked, marker); start >= 0 {
			emitCryptoFinding(unit, &MetaCWE295, start, "TLS certificate validation is explicitly disabled", confidence84, out)
			return
		}
	}
	if match := pyCheckHostnameFalseRE.FindStringIndex(masked); match != nil {
		emitCryptoFinding(unit, &MetaCWE295, match[0], "TLS hostname verification is explicitly disabled", confidence84, out)
	}
}

// CWE-328 reports legacy MD5 and SHA-1 hash construction only when the call
// sits in a security-sensitive context (password / credential / token / auth /
// key-derivation). Non-security fingerprints — checksums, format IDs, bench
// digests, test equality hashes — are intentionally omitted. Comments and
// docstrings are already masked before call discovery.
func detectCWE328(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"hashlib.md5", "hashlib.sha1", "Crypto.Hash.MD5.new", "Crypto.Hash.SHA1.new"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if securityHashContext(facts, unit.Source, call.Start) {
				emitCryptoFinding(unit, &MetaCWE328, call.Start, "weak MD5 or SHA-1 hash algorithm is used", confidence84, out)
				return
			}
		}
	}
	for _, call := range findCalls(facts, unit.Source, "hashlib.new") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) > 0 && weakHashName(args[0]) && securityHashContext(facts, unit.Source, call.Start) {
			emitCryptoFinding(unit, &MetaCWE328, call.Start, "weak MD5 or SHA-1 hash algorithm is selected dynamically", confidence84, out)
			return
		}
	}
}

func weakHashName(expr string) bool {
	t := strings.ToLower(strings.TrimSpace(expr))
	return t == "'md5'" || t == `"md5"` || t == "'sha1'" || t == `"sha1"`
}

// securityHashContext mirrors securityRandomContext: require a security-shaped
// identifier on the call line, or in the enclosing function name/body.
func securityHashContext(facts *PyCweFacts, source string, offset int) bool {
	if offset < 0 || offset > len(source) {
		return false
	}
	start := strings.LastIndex(source[:offset], "\n") + 1
	end := len(source)
	if next := strings.Index(source[offset:], "\n"); next >= 0 {
		end = offset + next
	}
	lineSrc := source[start:end]
	var line string
	if facts != nil {
		line = strings.ToLower(facts.codeMask(lineSrc, start))
	} else {
		line = strings.ToLower(pythonCodeMask(lineSrc))
	}
	if hasSecurityHashIdent(line) {
		return true
	}
	if facts == nil {
		return false
	}
	fn, ok := containingPythonFunction(facts.Functions(), offset)
	if !ok {
		return false
	}
	if hasSecurityHashIdent(strings.ToLower(fn.name)) {
		return true
	}
	return hasSecurityHashIdent(strings.ToLower(facts.codeMask(fn.body, fn.bodyStart)))
}

func hasSecurityHashIdent(text string) bool {
	// Distinctive substrings match compound names (user_password, hash_password).
	for _, frag := range []string{"password", "passwd", "credential", "pbkdf", "key_derivation", "derive_key"} {
		if strings.Contains(text, frag) {
			return true
		}
	}
	// Short tokens stay whole-identifier to avoid "author"/"tokenizer" noise.
	for _, ident := range []string{"token", "auth"} {
		if containsIdent(text, ident) {
			return true
		}
	}
	return false
}

// CWE-335 limits reports to an explicitly constant seed. A seed drawn from
// the runtime or omitted entirely is not reported by this same-file rule.
func detectCWE335(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"random.seed", "numpy.random.seed", "np.random.seed"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) > 0 && pyFixedSeedRE.MatchString(args[0]) {
				emitCryptoFinding(unit, &MetaCWE335, call.Start, "pseudo-random generator is seeded with a fixed value", confidence78, out)
				return
			}
		}
	}
}

// CWE-338 and CWE-1241 require a security-shaped identifier on the same
// executable line as a non-cryptographic random API. This avoids flagging
// ordinary sampling, games, and simulations.
func detectCWE338(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	detectWeakSecurityRandom(unit, facts, &MetaCWE338, "cryptographically weak random API produces a security-sensitive value", confidence82, out)
}

func detectCWE1241(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	detectWeakSecurityRandom(unit, facts, &MetaCWE1241, "predictable random algorithm produces a security-sensitive value", confidence80, out)
}

func detectWeakSecurityRandom(unit *core.ParsedUnit, facts *PyCweFacts, meta *rules.RuleMetadata, message string, confidence float32, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"random.random", "random.randint", "random.randrange", "random.choice", "random.getrandbits", "numpy.random.randint", "numpy.random.choice", "np.random.randint", "np.random.choice"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if securityRandomContext(facts, unit.Source, call.Start) {
				emitCryptoFinding(unit, meta, call.Start, message, confidence, out)
				return
			}
		}
	}
}

func securityRandomContext(facts *PyCweFacts, source string, offset int) bool {
	if offset < 0 || offset > len(source) {
		return false
	}
	start := strings.LastIndex(source[:offset], "\n") + 1
	end := len(source)
	if next := strings.Index(source[offset:], "\n"); next >= 0 {
		end = offset + next
	}
	lineSrc := source[start:end]
	var line string
	if facts != nil {
		line = strings.ToLower(facts.codeMask(lineSrc, start))
	} else {
		line = strings.ToLower(pythonCodeMask(lineSrc))
	}
	for _, ident := range []string{"token", "secret", "password", "session", "nonce", "otp", "csrf", "api_key", "auth"} {
		if containsIdent(line, ident) {
			return true
		}
	}
	return false
}

// CWE-347 recognizes the explicit PyJWT signature-verification bypasses. A
// decode call with ordinary defaults remains safe because verification policy
// can otherwise be established outside this source file.
func detectCWE347(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"jwt.decode", "jwt.api_jwt.decode"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			compact := strings.ToLower(compactWhitespace(call.ArgsText))
			if strings.Contains(compact, "verify=false") || strings.Contains(compact, "'verify_signature':false") || strings.Contains(compact, `"verify_signature":false`) {
				emitCryptoFinding(unit, &MetaCWE347, call.Start, "JWT signature verification is explicitly disabled", confidence86, out)
				return
			}
		}
	}
}

// CWE-1204 detects literal IV material only when passed to a cipher mode that
// takes an IV. IVs obtained from secrets.token_bytes are intentionally safe.
func detectCWE1204(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "AES.new") {
		args := splitTopLevelArgs(call.ArgsText)
		if fixedIVArguments(args) {
			emitCryptoFinding(unit, &MetaCWE1204, call.Start, "cipher uses a fixed literal initialization vector", confidence82, out)
			return
		}
	}
}

func fixedIVArguments(args []string) bool {
	for i, arg := range args {
		t := strings.TrimSpace(arg)
		if key, value, ok := strings.Cut(t, "="); ok && strings.TrimSpace(key) == "iv" {
			return pyFixedIVRE.MatchString(strings.TrimSpace(value))
		}
		if i == 2 && len(args) >= 3 && strings.Contains(strings.Join(args, " "), "MODE_CBC") {
			return pyFixedIVRE.MatchString(t)
		}
	}
	return false
}

// CWE-1240 reports only an overt hand-written XOR cipher implementation. A
// generic XOR operation is common in Python and is deliberately not enough.
func detectCWE1240(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range facts.Functions() {
		if pyRiskyCryptoNameRE.MatchString(fn.name) && strings.Contains(facts.codeMask(fn.body, fn.bodyStart), "^") {
			emitCryptoFinding(unit, &MetaCWE1240, fn.start, "hand-written XOR cipher is a risky cryptographic primitive implementation", confidence76, out)
			return
		}
	}
}

// CWE-1392 reports an explicit common default password assignment. It is kept
// separate from CWE-256 because non-default literal passwords have different
// remediation and are already covered by that storage rule.
func detectCWE1392(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstLiteralMatchStartIfContains(facts, unit, pyDefaultCredentialRE,
		"='", "=\"", "= '", "= \""); start >= 0 {
		emitCryptoFinding(unit, &MetaCWE1392, start, "common default password is assigned in Python source", confidence82, out)
	}
}

func emitCryptoFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}

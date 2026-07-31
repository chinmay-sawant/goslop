package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

var (
	// Direct Django/DRF-style unpacking is a high-signal mass-assignment sink.
	pyMassAssignRequestUnpack = regexp.MustCompile(`(?m)(?:\.objects\.(?:create|update)|\.update|[A-Z][A-Za-z0-9_]*)\s*\(\s*\*\*\s*request\.(?:data|POST|json)\b`)
	// Attribute-map mutation is only a CWE-915 sink when its input is directly
	// request-derived, not for ordinary internal object bookkeeping.
	pyDictUpdateRequest  = regexp.MustCompile(`(?m)\.__dict__\s*\.\s*update\s*\(\s*request\.(?:data|POST|json)\b`)
	pyRequestSetattrLoop = regexp.MustCompile(`(?m)for\s+([A-Za-z_][A-Za-z0-9_]*)\s*,\s*([A-Za-z_][A-Za-z0-9_]*)\s+in\s+request\.(?:data|POST|json)\s*\.\s*items\s*\(\s*\)\s*:`)

	// Dynamically selected globals/locals/object variables must be directly
	// request-controlled. This intentionally excludes generic internal maps.
	pyDynamicNamespaceRequest = regexp.MustCompile(`(?m)(?:globals|locals)\s*\(\s*\)\s*\[[^\]\n]*(?:request\.(?:args|form|data|json|GET|POST|query_params)|request\.get_json\s*\()[^\]\n]*\]`)
	pyDynamicVarsRequest      = regexp.MustCompile(`(?m)vars\s*\([^\n)]*\)\s*\[[^\]\n]*(?:request\.(?:args|form|data|json|GET|POST|query_params)|request\.get_json\s*\()[^\]\n]*\]`)
	pyDynamicSetattrRequest   = regexp.MustCompile(`(?m)setattr\s*\([^\n,]+,\s*(?:request\.(?:args|form|data|json|GET|POST|query_params)|request\.get_json\s*\()[^\n,]*,`)

	pyPasswordIdentifier = regexp.MustCompile(`(?i)\b(?:password|passwd|pwd)\b`)
)

func init() {
	RegisterRule("CWE-915", detectCWE915, &MetaCWE915,
		".objects.create", ".objects.update", ".__dict__.update", "setattr(", "request.data", "request.POST")
	RegisterRule("CWE-914", detectCWE914, &MetaCWE914,
		"globals()", "locals()", "vars(", "setattr(")
	RegisterRule("CWE-916", detectCWE916, &MetaCWE916,
		"hashlib.md5", "hashlib.sha1", "crypt.crypt", "md5_crypt")
}

// detectCWE915 detects direct request payload unpacking and request-key
// setattr loops. It deliberately requires the source itself to establish the
// request origin, avoiding findings on internal maps and allowlisted updates.
func detectCWE915(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source
	if start := pyMassAssignRequestUnpack.FindStringIndex(src); start != nil {
		pushCWE915(unit, start[0], "request data is unpacked directly into object attributes (allowlist writable fields)", out)
		return
	}
	if start := pyDictUpdateRequest.FindStringIndex(src); start != nil {
		pushCWE915(unit, start[0], "request data directly updates an object's __dict__ (allowlist writable fields)", out)
		return
	}
	if start := requestSetattrLoopStart(src); start >= 0 {
		pushCWE915(unit, start, "setattr loop applies request-controlled keys without an allowlist", out)
	}
}

func pushCWE915(unit *core.ParsedUnit, start int, message string, out *[]rules.Finding) {
	line, col := unit.LineCol(start)
	rules.PushFindingWithConfidence(&MetaCWE915, unitFile(unit), line, col, message, 0.8, out)
}

func requestSetattrLoopStart(src string) int {
	match := pyRequestSetattrLoop.FindStringSubmatchIndex(src)
	if match == nil {
		return -1
	}
	key := src[match[2]:match[3]]
	value := src[match[4]:match[5]]
	// Restrict the check to the indented loop body, stopping at the next
	// top-level statement. That keeps an unrelated setattr later in the file
	// from turning a benign request-data iteration into a finding.
	bodyEnd := len(src)
	if next := strings.Index(src[match[1]:], "\n\n"); next >= 0 {
		bodyEnd = match[1] + next
	}
	body := src[match[1]:bodyEnd]
	setattr := regexp.MustCompile(`setattr\s*\([^,]+,\s*` + regexp.QuoteMeta(key) + `\s*,\s*` + regexp.QuoteMeta(value) + `\s*\)`)
	if found := setattr.FindStringIndex(body); found != nil {
		// An explicit membership guard in a conventionally named allowlist is a
		// safe suppression for this otherwise direct mass-assignment pattern.
		guard := regexp.MustCompile(`(?m)if\s+` + regexp.QuoteMeta(key) + `\s+in\s+(?:allowed|writable|permitted)[A-Za-z0-9_]*\s*:`)
		if guard.FindStringIndex(body[:found[0]]) != nil {
			return -1
		}
		return match[1] + found[0]
	}
	return -1
}

// detectCWE914 detects direct request-controlled variable and attribute names.
// It does not flag ordinary dict access or a static attribute name.
func detectCWE914(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, pattern := range []*regexp.Regexp{pyDynamicNamespaceRequest, pyDynamicVarsRequest, pyDynamicSetattrRequest} {
		if start := pattern.FindStringIndex(unit.Source); start != nil {
			line, col := unit.LineCol(start[0])
			rules.PushFindingWithConfidence(
				&MetaCWE914,
				unitFile(unit),
				line, col,
				"request-controlled value selects a variable or object attribute (use a fixed allowlist)",
				0.8,
				out,
			)
			return
		}
	}
}

// detectCWE916 detects fast, general-purpose password hashes. Generic uses of
// MD5/SHA-1 are left alone; a password-like identifier is required at the call
// site to retain a high signal-to-noise ratio.
func detectCWE916(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"hashlib.md5", "hashlib.sha1", "md5_crypt.hash", "passlib.hash.md5_crypt.hash"} {
		for _, call := range findCalls(unit.Source, name) {
			if !pyPasswordIdentifier.MatchString(call.ArgsText) {
				continue
			}
			pushCWE916(unit, call.Start, "fast password hash lacks sufficient computational effort (use Argon2id, bcrypt, scrypt, or PBKDF2)", out)
			return
		}
	}
	for _, call := range findCalls(unit.Source, "crypt.crypt") {
		if !pyPasswordIdentifier.MatchString(call.ArgsText) {
			continue
		}
		pushCWE916(unit, call.Start, "crypt.crypt password hash may use an insufficiently costly scheme (use a modern password-hashing API)", out)
		return
	}
}

func pushCWE916(unit *core.ParsedUnit, start int, message string, out *[]rules.Finding) {
	line, col := unit.LineCol(start)
	rules.PushFindingWithConfidence(&MetaCWE916, unitFile(unit), line, col, message, 0.8, out)
}

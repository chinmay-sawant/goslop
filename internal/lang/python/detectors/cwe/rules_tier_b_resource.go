package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-367", detectCWE367, &MetaCWE367,
		"os.path.exists", "os.path.lexists")
	RegisterRule("CWE-403", detectCWE403, &MetaCWE403,
		"subprocess.", "close_fds")
	RegisterRule("CWE-409", detectCWE409, &MetaCWE409,
		".extractall")
	RegisterRule("CWE-454", detectCWE454, &MetaCWE454,
		"pickle.load", "pickle.loads", "runpy.run_path")
	RegisterRule("CWE-472", detectCWE472, &MetaCWE472,
		"**request.data", "**request.POST", "request.data", "request.POST")
	RegisterRule("CWE-521", detectCWE521, &MetaCWE521,
		"MIN_LENGTH", "MINIMUM_LENGTH", "PASSWORD_MIN")
	RegisterRule("CWE-524", detectCWE524, &MetaCWE524,
		"cache.set")
	RegisterRule("CWE-538", detectCWE538, &MetaCWE538,
		"open(", "/tmp/", "static/", "media/")
	RegisterRule("CWE-552", detectCWE552, &MetaCWE552,
		"send_file", "send_from_directory", "FileResponse")
	RegisterRule("CWE-617", detectCWE617, &MetaCWE617,
		"assert", "request.args", "request.user", "request.headers", "request.cookies")
	RegisterRule("CWE-641", detectCWE641, &MetaCWE641,
		"os.path.join", "request.files", "request.args", "request.form")
	RegisterRule("CWE-648", detectCWE648, &MetaCWE648,
		"ctypes.CDLL", "os.setuid", "os.seteuid")
	RegisterRule("CWE-779", detectCWE779, &MetaCWE779,
		"password", "secret", "token", "api_key", "credit_card", "logging.", "logger.", "log.")
	RegisterRule("CWE-836", detectCWE836, &MetaCWE836,
		"authenticate")
	RegisterRule("CWE-838", detectCWE838, &MetaCWE838,
		"Markup")
}

var (
	pyTierBSensitiveCache    = regexp.MustCompile(`(?is)cache\.set\s*\([^\n]*(?:password|secret|token|ssn|credit_card)`)
	pyTierBWriteSecretRE     = regexp.MustCompile(`(?is)open\s*\(\s*["'](?:/tmp/|static/|media/)[^"']*(?:secret|token|password)[^"']*["']\s*,\s*["']w`)
	pyTierBAssertRE          = regexp.MustCompile(`(?im)^\s*assert\s+request\.(?:user|args|headers|cookies)\b`)
	pyTierBLogSecretRE       = regexp.MustCompile(`(?is)(?:logger|logging|log)\.(?:debug|info|warning|error)\s*\([^\n]*(?:password|secret|token|api_key|credit_card)`)
	pyTierBRequestUnpackRE   = regexp.MustCompile(`(?is)\w+\s*\(\s*\*\*request\.(?:data|POST)\s*\)`)
	pyTierBRequestPathOpenRE = regexp.MustCompile(`(?is)open\s*\(\s*os\.path\.join\s*\([^\n]*request\.(?:files|args|form)[^\n]*\)`)
)

func detectCWE367(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if isPythonTestModule(unit) {
		return
	}
	if start := toctouSamePathStart(facts, unit.Source); start >= 0 {
		emitTierBFinding(unit, &MetaCWE367, start, "filesystem path is checked before a later separate use", confidence84, out)
	}
}

func detectCWE403(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "subprocess.run", "subprocess.Popen", "subprocess.call") {
		if hasBooleanKwarg(call.ArgsText, "close_fds", "False") {
			emitTierBFinding(unit, &MetaCWE403, call.Start, "subprocess inherits open file descriptors", confidence86, out)
			return
		}
	}
}

func detectCWE409(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, ".extractall") {
		emitTierBFinding(unit, &MetaCWE409, call.Start, "archive is extracted without an observable size or member limit", confidence78, out)
		return
	}
}

func detectCWE454(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "pickle.load", "pickle.loads", "runpy.run_path") {
		if strings.Contains(call.ArgsText, "request.") {
			emitTierBFinding(unit, &MetaCWE454, call.Start, "request-controlled data initializes trusted application state", confidence84, out)
			return
		}
	}
}

func detectCWE472(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBRequestUnpackRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE472, start, "request mapping is expanded into a model or trusted constructor", confidence82, out)
	}
}

func detectCWE521(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBWeakLengthRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE521, start, "password minimum length is configured below six characters", confidence80, out)
	}
}

func detectCWE524(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBSensitiveCache); start >= 0 {
		emitTierBFinding(unit, &MetaCWE524, start, "sensitive value is stored in a shared cache", confidence82, out)
	}
}

func detectCWE538(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstLiteralMatchStartIfContains(facts, unit, pyTierBWriteSecretRE,
		"/tmp/", "static/", "media/", "secret", "token", "password"); start >= 0 {
		emitTierBFinding(unit, &MetaCWE538, start, "secret-named file is written in an externally accessible directory", confidence84, out)
	}
}

func detectCWE552(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "send_file", "send_from_directory", "FileResponse") {
		if strings.Contains(call.ArgsText, "request.") {
			emitTierBFinding(unit, &MetaCWE552, call.Start, "request-controlled resource is sent to an external party", confidence76, out)
			return
		}
	}
}

func detectCWE617(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBAssertRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE617, start, "reachable assertion depends on request-controlled state", confidence80, out)
	}
}

func detectCWE641(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBRequestPathOpenRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE641, start, "request filename is used directly as a filesystem resource name", confidence82, out)
	}
}

func detectCWE648(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "ctypes.CDLL", "os.setuid", "os.seteuid") {
		emitTierBFinding(unit, &MetaCWE648, call.Start, "privileged operating-system API is called directly", confidence80, out)
		return
	}
}

func detectCWE779(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBLogSecretRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE779, start, "sensitive value is logged without a visible redaction", confidence84, out)
	}
}

func detectCWE836(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "authenticate") {
		if strings.Contains(strings.ToLower(call.ArgsText), "password_hash=request") {
			emitTierBFinding(unit, &MetaCWE836, call.Start, "client-supplied password hash is passed to authentication", confidence86, out)
			return
		}
	}
}

func detectCWE838(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "Markup") {
		if strings.Contains(call.ArgsText, "urllib.parse.quote(") {
			emitTierBFinding(unit, &MetaCWE838, call.Start, "URL encoding is used where HTML-context output encoding is required", confidence82, out)
			return
		}
	}
}

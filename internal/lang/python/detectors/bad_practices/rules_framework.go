package badpractices

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-PY-16", detectBPPY16)
	RegisterRule("BP-PY-17", detectBPPY17)
	RegisterRule("BP-PY-21", detectBPPY21)
}

var (
	appRunDebugRe   = regexp.MustCompile(`\.run\s*\([^)]*debug\s*=\s*True`)
	debugAssignRe   = regexp.MustCompile(`(?i)\bDEBUG\s*=\s*True\b`)
	appConfigDebug  = regexp.MustCompile(`(?i)(?:app\.config\s*\[\s*['"]DEBUG['"]\s*\]|\.config\s*\[\s*['"]DEBUG['"]\s*\])\s*=\s*True`)
	secretKeyRe     = regexp.MustCompile(`(?i)(?:app\.secret_key|SECRET_KEY)\s*=\s*`)
	flaskSecretHint = regexp.MustCompile(`(?i)\b(secret_key|SECRET_KEY)\b`)
)

// BP-PY-16: Flask DEBUG True in production code.
func detectBPPY16(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-16")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.hasAny("debug=True", "DEBUG", "app.run") &&
		!strings.Contains(unit.Source, "debug=True") &&
		!strings.Contains(unit.Source, "DEBUG") {
		return
	}
	src := unit.Source
	// app.run(debug=True) — may span lines; use a windowed approach via lines join for simple cases.
	if appRunDebugRe.MatchString(src) {
		loc := appRunDebugRe.FindStringIndex(src)
		if loc != nil {
			pushAt(unit, meta, loc[0], "Flask debug=True exposes an interactive debugger; disable in production", out)
		}
	}
	// Multi-line app.run( ... debug=True )
	detectMultilineAppRunDebug(unit, src, meta, out)

	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if appConfigDebug.MatchString(t) {
			pushAt(unit, meta, line.byte, "Flask DEBUG enabled in config; disable in production", out)
			continue
		}
		// DEBUG = True only when file looks Flask-related (or generic module-level DEBUG).
		if debugAssignRe.MatchString(t) && looksFlaskish(unit, src) {
			pushAt(unit, meta, line.byte, "DEBUG = True in non-test code; keep debug off in production", out)
		}
	}
}

func detectMultilineAppRunDebug(unit *core.ParsedUnit, src string, meta *rules.RuleMetadata, out *[]rules.Finding) {
	// Find .run( and scan args for debug=True without requiring single-line match.
	start := 0
	for {
		idx := strings.Index(src[start:], ".run(")
		if idx < 0 {
			return
		}
		abs := start + idx
		open := abs + len(".run")
		if open >= len(src) || src[open] != '(' {
			start = abs + 4
			continue
		}
		inner, _, ok := callArgsRegion(src, open)
		if ok && strings.Contains(inner, "debug=True") {
			// Avoid double-report if single-line regex already matched this region.
			already := false
			for _, f := range *out {
				if f.RuleID == "BP-PY-16" {
					// approximate: same line
					line, _ := unit.LineCol(abs)
					if f.Line == line {
						already = true
						break
					}
				}
			}
			if !already {
				pushAt(unit, meta, abs, "Flask debug=True exposes an interactive debugger; disable in production", out)
			}
		}
		start = abs + 4
	}
}

func looksFlaskish(unit *core.ParsedUnit, src string) bool {
	if strings.Contains(src, "flask") || strings.Contains(src, "Flask") {
		return true
	}
	p := strings.ToLower(fileDisplayPath(unit))
	return strings.Contains(p, "flask") || strings.Contains(p, "app.py") || strings.Contains(p, "wsgi")
}

// BP-PY-17: Flask SECRET_KEY / secret_key hardcoded string.
func detectBPPY17(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-17")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.hasAny("secret_key", "SECRET_KEY") && !flaskSecretHint.MatchString(unit.Source) {
		return
	}
	// Prefer Flask-ish modules to reduce collision with Django SECRET_KEY (BP-PY-22 deferred).
	// Still flag app.secret_key always; SECRET_KEY when flaskish or not clearly django settings.
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		loc := secretKeyRe.FindStringIndex(t)
		if loc == nil {
			continue
		}
		if strings.Contains(t, "os.environ") || strings.Contains(t, "getenv") {
			continue
		}
		eq := strings.Index(t[loc[0]:], "=")
		if eq < 0 {
			continue
		}
		rhs := strings.TrimSpace(t[loc[0]+eq+1:])
		if !isStringLiteral(rhs) {
			continue
		}
		if looksLikePlaceholderSecret(rhs) {
			continue
		}
		val := unwrapStringLiteral(rhs)
		if len(val) < 4 {
			continue
		}
		// app.secret_key always; bare SECRET_KEY when flaskish or not django settings path.
		isAppSecret := strings.Contains(strings.ToLower(t), "secret_key") && strings.Contains(t, ".")
		if isAppSecret || looksFlaskish(unit, unit.Source) || !looksDjangoSettings(unit) {
			// Avoid double-firing BP-PY-13 on same line is OK.
			if strings.Contains(t, "app.secret_key") || strings.Contains(t, "SECRET_KEY") {
				pushAt(unit, meta, line.byte+loc[0], "hardcoded Flask SECRET_KEY; load from environment", out)
			}
		}
	}
}

// BP-PY-21: Django DEBUG = True in settings modules.
func detectBPPY21(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-21")
	if isPythonTestFile(unit) {
		return
	}
	if !looksDjangoSettings(unit) {
		// Also allow explicit django settings content markers.
		if !strings.Contains(unit.Source, "django") && !strings.Contains(unit.Source, "INSTALLED_APPS") &&
			!strings.Contains(unit.Source, "MIDDLEWARE") {
			return
		}
		// Without path or django markers, skip to avoid Flask DEBUG collision.
		if !looksDjangoSettings(unit) {
			return
		}
	}
	// Skip local_settings patterns.
	base := strings.ToLower(filepath.Base(fileDisplayPath(unit)))
	if strings.Contains(base, "local_settings") || strings.Contains(base, "dev_settings") {
		return
	}
	if !facts.has("DEBUG") && !strings.Contains(unit.Source, "DEBUG") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if debugAssignRe.MatchString(t) {
			pushAt(unit, meta, line.byte, "Django DEBUG = True in settings; keep False in production", out)
		}
	}
}

func looksDjangoSettings(unit *core.ParsedUnit) bool {
	p := filepath.ToSlash(strings.ToLower(fileDisplayPath(unit)))
	base := filepath.Base(p)
	if base == "settings.py" {
		return true
	}
	if strings.Contains(p, "/settings/") || strings.Contains(p, "/settings.py") {
		return true
	}
	if strings.HasPrefix(base, "settings") && strings.HasSuffix(base, ".py") {
		return true
	}
	return false
}

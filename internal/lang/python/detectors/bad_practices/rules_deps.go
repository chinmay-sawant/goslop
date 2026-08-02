package badpractices

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	metaByID["BP-PY-43"] = &rules.RuleMetadata{
		ID: "BP-PY-43", Title: "requirements Without Pins",
		Description: "requirements.txt lists direct deps without version pins, risking non-reproducible builds.",
		Severity:    rules.SeverityLow, Pack: rules.PackBadPractice,
		Fix: "Pin versions (==, ~=, or lockfiles) for application dependencies.",
	}
	metaByID["BP-PY-44"] = &rules.RuleMetadata{
		ID: "BP-PY-44", Title: "Import Deprecated stdlib Module",
		Description: "Code imports a deprecated or removed stdlib module (e.g. `imp`, `asyncore`, `cgi` patterns).",
		Severity:    rules.SeverityLow, Pack: rules.PackBadPractice,
		Fix: "Migrate to the recommended replacement (importlib, asyncio, etc.).",
	}
	metaByID["BP-PY-45"] = &rules.RuleMetadata{
		ID: "BP-PY-45", Title: "sys.path Mutation At Runtime",
		Description: "`sys.path.insert/append` is used to fix imports instead of proper packaging.",
		Severity:    rules.SeverityLow, Pack: rules.PackBadPractice,
		Fix: "Install the package (editable install) or set PYTHONPATH; avoid mutating sys.path at runtime.",
	}
	RegisterRule("BP-PY-43", detectBPPY43)
	RegisterRule("BP-PY-44", detectBPPY44)
	RegisterRule("BP-PY-45", detectBPPY45)
}

// isRequirementsPath reports requirements* unit paths (Path A for BP-PY-43).
// Accepts requirements*.txt (real projects) and requirements*.py (fixtures:
// walk only collects language extensions, so matrices materialize as .py).
func isRequirementsPath(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, p := range []string{fileDisplayPath(unit), unit.Path} {
		if p == "" {
			continue
		}
		base := strings.ToLower(filepath.Base(p))
		// requirements.txt, requirements-dev.txt, requirements_test.txt,
		// requirements-vulnerable.py (fixture materialization).
		if base == "requirements.txt" {
			return true
		}
		if strings.HasPrefix(base, "requirements") &&
			(strings.HasSuffix(base, ".txt") || strings.HasSuffix(base, ".py")) {
			return true
		}
		// requirements/base.txt style
		norm := filepath.ToSlash(strings.ToLower(p))
		if strings.Contains(norm, "/requirements/") &&
			(strings.HasSuffix(norm, ".txt") || strings.HasSuffix(norm, ".py")) {
			return true
		}
	}
	return false
}

// BP-PY-43: bare package lines in requirements*.txt without version pins.
// Path-gated: only runs when unit path looks like requirements.txt (not .py source).
func detectBPPY43(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-43")
	if !isRequirementsPath(unit) {
		return
	}
	// Line scan of requirements content (not Python).
	src := unit.Source
	if src == "" {
		return
	}
	lines := codeLinesFacts(facts, src)
	// Prefer raw lines without Python comment stripping for requirements? stripPyComment
	// treats # as comment which is correct for requirements.
	for _, line := range lines {
		raw := strings.TrimSpace(line.raw)
		if raw == "" {
			continue
		}
		// Full-line comment
		if strings.HasPrefix(raw, "#") {
			continue
		}
		// Strip inline comment
		text := strings.TrimSpace(line.text)
		if text == "" {
			continue
		}
		if isRequirementsDirective(text) {
			continue
		}
		if isPinnedRequirement(text) {
			continue
		}
		if isVCSRequirement(text) {
			continue
		}
		// Bare package name (with optional extras [foo]) and no version op.
		if looksBareRequirement(text) {
			pushAt(unit, meta, line.byte,
				"requirements entry has no version pin; prefer ==/~=/lockfiles for reproducible builds",
				out)
		}
	}
}

func isRequirementsDirective(line string) bool {
	// -r, -e, -c, --hash, -i, etc.
	if strings.HasPrefix(line, "-") {
		return true
	}
	lower := strings.ToLower(line)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return true
	}
	return false
}

func isVCSRequirement(line string) bool {
	lower := strings.ToLower(line)
	return strings.Contains(lower, "git+") || strings.Contains(lower, "hg+") ||
		strings.Contains(lower, "svn+") || strings.Contains(lower, "bzr+")
}

func isPinnedRequirement(line string) bool {
	// Version operators / markers that imply constraint.
	// Also wheel/path with == is pinned.
	for _, op := range []string{"==", "~=", ">=", "<=", "!=", "===", ">", "<"} {
		if strings.Contains(line, op) {
			return true
		}
	}
	return false
}

func looksBareRequirement(line string) bool {
	// Package name: letter/digit/._- and optional [extras]
	// Reject lines with spaces that look like options already handled.
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}
	// Environment markers after ;
	if i := strings.Index(line, ";"); i >= 0 {
		line = strings.TrimSpace(line[:i])
	}
	if line == "" {
		return false
	}
	// Must start with identifier-ish package name
	r := []rune(line)
	if !unicode.IsLetter(r[0]) && r[0] != '_' {
		return false
	}
	// No path separators for local paths without pin (./pkg) — treat as directive-like miss.
	if strings.HasPrefix(line, ".") || strings.HasPrefix(line, "/") {
		return false
	}
	return true
}

// Deprecated stdlib modules (PEP 594 removals and long-deprecated).
// Map module → preferred replacement (message).
var deprecatedStdlib = map[string]string{
	"imp":         "importlib",
	"asyncore":    "asyncio",
	"asynchat":    "asyncio",
	"cgi":         "urllib.parse / email / custom form parsing",
	"telnetlib":   "third-party or socket-based clients",
	"uu":          "base64",
	"xdrlib":      "a maintained serialization library",
	"aifc":        "modern audio libraries",
	"audioop":     "modern audio libraries",
	"chunk":       "modern audio/container libraries",
	"msilib":      "platform-specific install tooling",
	"nis":         "platform-native directory services",
	"ossaudiodev": "modern audio libraries",
	"pipes":       "subprocess",
	"sunau":       "modern audio libraries",
	"formatter":   "custom formatting",
	"parser":      "ast",
	"symbol":      "ast",
	"binhex":      "modern encoding utilities",
	"smtpd":       "aiosmtpd or a maintained server",
	"optparse":    "argparse",
}

// BP-PY-44: import of deprecated/removed stdlib modules.
func detectBPPY44(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-44")
	if isRequirementsPath(unit) {
		return
	}
	// Early-out if none of the module names appear.
	if !facts.hasAny("import ", "from ") && !strings.Contains(unit.Source, "import") {
		return
	}
	// Cheap prefilter: any deprecated name in source?
	hasDep := false
	for mod := range deprecatedStdlib {
		if facts.has(mod) || strings.Contains(unit.Source, mod) {
			hasDep = true
			break
		}
	}
	if !hasDep {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		mod, ok := importedDeprecatedModule(t)
		if !ok {
			continue
		}
		repl := deprecatedStdlib[mod]
		msg := "imports deprecated/removed stdlib module " + mod
		if repl != "" {
			msg += "; prefer " + repl
		}
		pushAt(unit, meta, line.byte, msg, out)
	}
}

func importedDeprecatedModule(t string) (string, bool) {
	// import imp / import imp as x / import asyncore, cgi
	// from imp import load_source / from cgi import FieldStorage
	if strings.HasPrefix(t, "import ") {
		rest := strings.TrimSpace(strings.TrimPrefix(t, "import "))
		// Split by comma for multi-import
		parts := strings.Split(rest, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			// strip "as name"
			if i := strings.Index(p, " as "); i >= 0 {
				p = strings.TrimSpace(p[:i])
			}
			// dotted: imp.foo → root imp
			root := p
			if i := strings.Index(p, "."); i >= 0 {
				root = p[:i]
			}
			if _, ok := deprecatedStdlib[root]; ok {
				return root, true
			}
		}
	}
	if strings.HasPrefix(t, "from ") {
		// from imp import ...
		rest := strings.TrimSpace(strings.TrimPrefix(t, "from "))
		// module before " import "
		idx := strings.Index(rest, " import ")
		if idx < 0 {
			// from imp importX invalid — try " import"
			idx = strings.Index(rest, " import")
		}
		if idx < 0 {
			return "", false
		}
		mod := strings.TrimSpace(rest[:idx])
		if i := strings.Index(mod, "."); i >= 0 {
			mod = mod[:i]
		}
		if _, ok := deprecatedStdlib[mod]; ok {
			return mod, true
		}
	}
	return "", false
}

// BP-PY-45: sys.path.insert/append/extend at runtime (skip test files / Sphinx docs conf).
func detectBPPY45(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-45")
	if isRequirementsPath(unit) {
		return
	}
	if isPythonTestFile(unit) {
		return
	}
	display := fileDisplayPath(unit)
	// Skip known bootstrap basenames and Sphinx docs/conf.py.
	base := strings.ToLower(filepath.Base(display))
	switch base {
	case "sitecustomize.py", "usercustomize.py", "conftest.py":
		return
	}
	if isSphinxDocsConfPath(display) || isSphinxDocsConfPath(unit.Path) {
		return
	}
	if !facts.hasAny("sys.path", "sys.path.insert", "sys.path.append", "sys.path.extend") &&
		!strings.Contains(unit.Source, "sys.path") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	type scope struct {
		indent int
		kind   byte // 'd' def/class, 'i' if/for/while/with
	}
	var stack []scope
	prev := ""
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		ind := indentWidth(line.raw)
		for len(stack) > 0 && ind <= stack[len(stack)-1].indent {
			stack = stack[:len(stack)-1]
		}
		inFunc := false
		for _, s := range stack {
			if s.kind == 'd' {
				inFunc = true
				break
			}
		}
		isPathMut := strings.Contains(t, "sys.path.insert(") || strings.Contains(t, "sys.path.append(") ||
			strings.Contains(t, "sys.path.extend(") ||
			strings.Contains(t, "sys.path.insert (") || strings.Contains(t, "sys.path.append (")
		if isPathMut && !isSysPathBootstrap(t, inFunc, prev) {
			off := line.byte
			for _, n := range []string{"sys.path.insert", "sys.path.append", "sys.path.extend"} {
				if i := strings.Index(line.text, n); i >= 0 {
					off = line.byte + i
					break
				}
			}
			pushAt(unit, meta, off,
				"sys.path mutation at runtime; prefer packaging/editable installs over path hacks",
				out)
		}
		// Track scopes after evaluating the line (headers open a new block).
		if strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") || strings.HasPrefix(t, "class ") {
			stack = append(stack, scope{indent: ind, kind: 'd'})
		} else if strings.HasSuffix(t, ":") &&
			(strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ") || strings.HasPrefix(t, "else:") ||
				strings.HasPrefix(t, "for ") || strings.HasPrefix(t, "while ") ||
				strings.HasPrefix(t, "with ") || strings.HasPrefix(t, "async with ") ||
				strings.HasPrefix(t, "try:") || strings.HasPrefix(t, "except") ||
				strings.HasPrefix(t, "finally:") || strings.HasPrefix(t, "match ")) {
			stack = append(stack, scope{indent: ind, kind: 'i'})
		}
		prev = t
	}
}

// isSysPathBootstrap reports import-time in-tree bootstraps that are not runtime library mutations.
// Exempts module-level (a) guarded bootstraps — `if <root> not in sys.path:` on the preceding
// line — and (b) inserts whose argument is built from `Path(__file__)`: the standalone in-tree
// script/tool pattern that imports the uninstalled package. In-function inserts (lazy bootstrap
// inside library code) and unguarded hard-coded vendor inserts stay reportable.
func isSysPathBootstrap(t string, inFunc bool, prev string) bool {
	if inFunc {
		return false
	}
	if strings.Contains(t, "Path(__file__") {
		return true
	}
	p := strings.TrimSpace(prev)
	return strings.HasPrefix(p, "if ") && strings.Contains(p, "not in sys.path")
}

// isSphinxDocsConfPath reports docs/**/conf.py (Sphinx build configuration).
func isSphinxDocsConfPath(p string) bool {
	if p == "" {
		return false
	}
	norm := filepath.ToSlash(strings.ToLower(p))
	base := filepath.Base(norm)
	if base != "conf.py" {
		return false
	}
	return strings.Contains(norm, "/docs/") || strings.HasPrefix(norm, "docs/")
}

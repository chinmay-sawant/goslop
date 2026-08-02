// Package badpractices detects common Python implementation mistakes.
package badpractices

import (
	"path/filepath"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func fileDisplayPath(unit *core.ParsedUnit) string {
	if unit == nil {
		return ""
	}
	if unit.DisplayPath != "" {
		return unit.DisplayPath
	}
	return unit.Path
}

// isPythonNoxOrToxFile reports nox/tox session orchestration modules. Their
// test_* session functions run pytest (or other tools) and are not placeholder
// unit tests (niquests noxfile.py BP-PY-41 FPs).
func isPythonNoxOrToxFile(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, p := range []string{fileDisplayPath(unit), unit.Path} {
		if p == "" {
			continue
		}
		base := filepath.Base(p)
		if base == "noxfile.py" || base == "tox.ini" || base == "tox.py" {
			return true
		}
	}
	return false
}

// isPythonTestFile reports test modules we should skip for validation-style rules.
// Matches test_*.py, *_test.py, and paths under tests/ or test/.
func isPythonTestFile(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, p := range []string{fileDisplayPath(unit), unit.Path} {
		if p == "" {
			continue
		}
		base := filepath.Base(p)
		if base == "tests.py" || base == "conftest.py" || (strings.HasPrefix(base, "test_") && strings.HasSuffix(base, ".py")) {
			return true
		}
		if strings.HasSuffix(base, "_test.py") {
			return true
		}
		// Normalize separators.
		norm := filepath.ToSlash(p)
		if strings.Contains(norm, "/tests/") || strings.Contains(norm, "/test/") ||
			strings.HasPrefix(norm, "tests/") || strings.HasPrefix(norm, "test/") {
			return true
		}
	}
	return false
}

// isPythonTestDirPath reports modules under a tests/ or test/ directory tree.
// Unlike isPythonTestFile, this does not match scripts merely named *_test.py
// (Project_Parva backend/tools/load_test.py).
func isPythonTestDirPath(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, p := range []string{fileDisplayPath(unit), unit.Path} {
		if p == "" {
			continue
		}
		norm := filepath.ToSlash(p)
		if strings.Contains(norm, "/tests/") || strings.Contains(norm, "/test/") ||
			strings.HasPrefix(norm, "tests/") || strings.HasPrefix(norm, "test/") {
			return true
		}
	}
	return false
}

// isPythonOfflineScriptPath reports offline tooling / release / validation
// scripts where broad-except noise is expected batch-job noise (Project_Parva
// tools/, scripts/release, public-benchmark runners). Does NOT match generic
// scripts/ trees (WeThePeople operational scripts keep BP-PY-1).
func isPythonOfflineScriptPath(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, path := range []string{fileDisplayPath(unit), unit.Path} {
		if path == "" {
			continue
		}
		norm := filepath.ToSlash(path)
		// Narrow offline trees only — do not match e.g. tools/benchmark-harness
		// under third-party packages (pdf_oxide TP).
		markers := []string{
			"/backend/tools/",
			"/scripts/release/",
			"/public-benchmark/",
		}
		for _, m := range markers {
			if strings.Contains(norm, m) {
				return true
			}
		}
		if strings.HasPrefix(norm, "backend/tools/") ||
			strings.HasPrefix(norm, "scripts/release/") ||
			strings.HasPrefix(norm, "public-benchmark/") {
			return true
		}
		// Project_Parva top-level tools/ (conformance_runner, validate_schemas).
		if strings.Contains(norm, "/Project_Parva/tools/") || strings.HasPrefix(norm, "tools/conformance") ||
			strings.HasPrefix(norm, "tools/validate") || strings.HasPrefix(norm, "tools/release/") {
			return true
		}
	}
	return false
}

// isPythonBenchmarkFile identifies benchmark harnesses whose print calls are
// intentional result output rather than application-library debugging.
func isPythonBenchmarkFile(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, path := range []string{fileDisplayPath(unit), unit.Path} {
		for _, component := range strings.Split(strings.Trim(filepath.ToSlash(path), "/"), "/") {
			if component == "bench" || component == "benchmarks" {
				return true
			}
			stem := strings.TrimSuffix(component, filepath.Ext(component))
			if stem == "bench" || stem == "benchmark" || stem == "benchmarks" ||
				strings.HasPrefix(stem, "bench_") || strings.HasPrefix(stem, "benchmark_") {
				return true
			}
		}
	}
	return false
}

// demoScriptPathDirNames are path components marking runnable demo/example
// trees (BP-PY-46 script exemption). Exemption is conditional on the module
// being self-running (a top-level print, or a __main__ guard invoking a
// module-defined function); library-style modules under these trees
// (web-function collections imported by hosts, e.g. FuncToWeb examples/)
// keep firing.
var demoScriptPathDirNames = map[string]struct{}{
	"examples": {}, "example": {},
	"demos": {}, "demo": {},
}

// toolScriptPathDirNames are path components marking dev-tool script trees.
// Exemption is narrow: only a single top-level completion print (dev-time
// regeneration scripts, e.g. rendercv scripts/update_*.py). Operational
// scripts that log via print (WeThePeople scripts/) keep firing.
var toolScriptPathDirNames = map[string]struct{}{
	"scripts": {}, "script": {},
	"tools": {}, "tool": {},
}

// cliScriptPathDirNames are path components marking CLI subcommand packages
// (Click/Typer) whose print calls are user-facing presentation.
var cliScriptPathDirNames = map[string]struct{}{
	"commands": {}, // Click/Typer CLI subcommand modules (e.g. package/commands/)
	"cli":      {}, // CLI entry package (e.g. rendercv/cli/app.py)
}

// isPythonScriptModule reports packaging/CLI/demo/tool entry modules where
// print is intentional user-facing output, not library debug logging.
func isPythonScriptModule(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, path := range []string{fileDisplayPath(unit), unit.Path} {
		if path == "" {
			continue
		}
		base := filepath.Base(path)
		if base == "setup.py" || base == "__main__.py" || base == "cli.py" ||
			strings.HasPrefix(base, "tmp_") {
			return true
		}
		for _, component := range strings.Split(strings.Trim(filepath.ToSlash(path), "/"), "/") {
			if _, ok := cliScriptPathDirNames[component]; ok {
				return true
			}
			if _, ok := demoScriptPathDirNames[component]; ok {
				// Self-running demo script: prints are the script's own output
				// (module top level, inside functions the __main__ guard
				// invokes, or a demo file with no guard at all). Importable
				// example modules with a guard keep firing.
				if pythonHasTopLevelPrint(unit.Source) ||
					pythonMainGuardCallsLocal(unit.Source) ||
					!pythonHasMainGuard(unit.Source) {
					return true
				}
			}
			if _, ok := toolScriptPathDirNames[component]; ok {
				// Dev-tool scripts:
				//  - single top-level completion print, or
				//  - shebang + no def/class (pure module-level CLI — pdf_oxide
				//    scripts/* debug runners with multi-print progress), or
				//  - shebang + every print is column-0 module-level, or
				//  - tool-package subdirectory whose prints sit on the
				//    __main__ call path.
				// Operational scripts that print from inside functions
				// (WeThePeople scripts/) keep firing.
				if pythonSingleTopLevelCompletionPrint(unit.Source) ||
					(isPythonShebangScript(unit.Source) && pythonIsPureModuleScript(unit.Source)) ||
					(isPythonShebangScript(unit.Source) && pythonAllPrintsModuleLevel(unit.Source)) ||
					(pythonToolSubdirPath(path) && pythonMainGuardCallsLocal(unit.Source)) {
					return true
				}
			}
		}
	}
	return false
}

// pythonToolSubdirPath reports a path where a scripts/tools component is
// followed by a subdirectory (scripts/<subdir>/<file>.py) rather than the file
// directly (scripts/<file>.py).
func pythonToolSubdirPath(path string) bool {
	parts := strings.Split(strings.Trim(filepath.ToSlash(path), "/"), "/")
	for i, part := range parts {
		if part == "scripts" || part == "tools" || part == "script" || part == "tool" {
			if i+2 < len(parts) {
				return true
			}
		}
	}
	return false
}

// isPythonShebangScript reports files whose first non-space content is a shebang.
// A shebang alone does not mean the file is a CLI entry script — see
// pythonShebangIsEntryScript.
func isPythonShebangScript(src string) bool {
	s := strings.TrimLeft(src, " \t")
	return strings.HasPrefix(s, "#!")
}

// pythonShebangIsEntryScript reports shebang files that also look like CLI entry
// scripts (__main__ guard or argparse/click/typer/cyclopts). A bare shebang is
// not enough: importable libraries often keep a vestigial #! line.
// Does not treat rich.print as an entry-script signal.
func pythonShebangIsEntryScript(src string) bool {
	if !isPythonShebangScript(src) {
		return false
	}
	if strings.Contains(src, "__name__") && strings.Contains(src, "__main__") {
		return true
	}
	if pythonHasArgparseCLI(src) {
		return true
	}
	return strings.Contains(src, "import click") ||
		strings.Contains(src, "from click ") ||
		strings.Contains(src, "import typer") ||
		strings.Contains(src, "from typer ") ||
		strings.Contains(src, "import cyclopts") ||
		strings.Contains(src, "from cyclopts ")
}

func pushAt(unit *core.ParsedUnit, meta *rules.RuleMetadata, byteOffset int, message string, out *[]rules.Finding) {
	if unit == nil || meta == nil || out == nil {
		return
	}
	line, col := unit.LineCol(byteOffset)
	rules.PushFinding(meta, fileDisplayPath(unit), line, col, message, out)
}

// stripPyComment removes a trailing # comment outside of string literals.
func stripPyComment(line string) string {
	inStr := byte(0)
	escape := false
	triple := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inStr == 0 {
			if c == '"' || c == '\'' {
				// Triple-quoted?
				if i+2 < len(line) && line[i+1] == c && line[i+2] == c {
					inStr = c
					triple = true
					i += 2
					continue
				}
				inStr = c
				continue
			}
			if c == '#' {
				return strings.TrimRight(line[:i], " \t")
			}
			continue
		}
		if escape {
			escape = false
			continue
		}
		if c == '\\' && !triple {
			escape = true
			continue
		}
		if triple {
			if c == inStr && i+2 < len(line) && line[i+1] == inStr && line[i+2] == inStr {
				inStr = 0
				triple = false
				i += 2
			}
			continue
		}
		if c == inStr {
			inStr = 0
		}
	}
	return line
}

// codeLine is one source line with 0-based index and byte offset.
type codeLine struct {
	idx  int
	text string // comment-stripped
	raw  string
	byte int
}

func buildCodeLines(source string) []codeLine {
	if source == "" {
		return nil
	}
	n := 1
	for i := 0; i < len(source); i++ {
		if source[i] == '\n' {
			n++
		}
	}
	out := make([]codeLine, 0, n)
	byteOff := 0
	lineIdx := 0
	for {
		nl := strings.IndexByte(source[byteOff:], '\n')
		var raw string
		if nl < 0 {
			raw = source[byteOff:]
			out = append(out, codeLine{idx: lineIdx, text: stripPyComment(raw), raw: raw, byte: byteOff})
			break
		}
		raw = source[byteOff : byteOff+nl]
		out = append(out, codeLine{idx: lineIdx, text: stripPyComment(raw), raw: raw, byte: byteOff})
		byteOff += nl + 1
		lineIdx++
		if byteOff >= len(source) {
			out = append(out, codeLine{idx: lineIdx, text: "", raw: "", byte: byteOff})
			break
		}
	}
	return out
}

func codeLinesFacts(facts *bpFacts, source string) []codeLine {
	if facts != nil && facts.lines != nil && (source == "" || source == facts.Source) {
		return facts.lines
	}
	return buildCodeLines(source)
}

func isIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// isSimpleIdent reports whether s is a non-empty Python identifier (no dots).
func isSimpleIdent(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isIdentByte(s[i]) {
			return false
		}
	}
	c := s[0]
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// indexOfIdent finds needle as a whole identifier/attribute fragment (not mid-ident).
func indexOfIdent(source, needle string) int {
	start := 0
	for {
		idx := strings.Index(source[start:], needle)
		if idx < 0 {
			return -1
		}
		abs := start + idx
		if abs > 0 {
			prev := source[abs-1]
			if isIdentByte(prev) {
				start = abs + len(needle)
				continue
			}
		}
		end := abs + len(needle)
		if end < len(source) && isIdentByte(source[end]) {
			start = end
			continue
		}
		return abs
	}
}

// findAllIdent returns all absolute offsets of needle as whole identifiers.
func findAllIdent(source, needle string) []int {
	var out []int
	start := 0
	for {
		idx := indexOfIdent(source[start:], needle)
		if idx < 0 {
			return out
		}
		abs := start + idx
		out = append(out, abs)
		start = abs + len(needle)
		if start >= len(source) {
			return out
		}
	}
}

// indentWidth returns leading space/tab count (tabs count as 1 for relative compare).
func indentWidth(line string) int {
	n := 0
	for _, r := range line {
		if r == ' ' || r == '\t' {
			n++
			continue
		}
		break
	}
	return n
}

// isStringLiteral reports whether s is a Python string/bytes literal (possibly with prefix).
func isStringLiteral(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// Strip common prefixes: r, u, b, f, fr, rf, br, rb (case-insensitive for b/u/r).
	for {
		if len(s) < 2 {
			break
		}
		c0 := s[0]
		if c0 == 'r' || c0 == 'R' || c0 == 'u' || c0 == 'U' || c0 == 'b' || c0 == 'B' || c0 == 'f' || c0 == 'F' {
			s = s[1:]
			continue
		}
		break
	}
	if len(s) < 2 {
		return false
	}
	q := s[0]
	if q != '"' && q != '\'' {
		return false
	}
	// Triple quotes.
	if len(s) >= 6 && s[1] == q && s[2] == q {
		return strings.HasSuffix(s, string([]byte{q, q, q}))
	}
	return s[len(s)-1] == q
}

// firstCallArg extracts the first top-level argument text of a call starting at
// openParen (index of '('). Returns arg and ok.
func firstCallArg(source string, openParen int) (string, bool) {
	if openParen < 0 || openParen >= len(source) || source[openParen] != '(' {
		return "", false
	}
	depth := 0
	inStr := byte(0)
	escape := false
	start := openParen + 1
	for i := openParen; i < len(source); i++ {
		c := source[i]
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
				// Full arg list region.
				inner := source[start:i]
				// Split first top-level comma.
				return strings.TrimSpace(splitFirstArg(inner)), true
			}
		case ',':
			if depth == 1 {
				return strings.TrimSpace(source[start:i]), true
			}
		}
	}
	return "", false
}

func splitFirstArg(inner string) string {
	depth := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(inner); i++ {
		c := inner[i]
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
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				return strings.TrimSpace(inner[:i])
			}
		}
	}
	return strings.TrimSpace(inner)
}

func findMatchingClose(source string, openParen int) int {
	depth := 0
	inStr := byte(0)
	escape := false
	for i := openParen; i < len(source); i++ {
		c := source[i]
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
				return i
			}
		}
	}
	return -1
}

// callArgsRegion returns the full argument list text inside (... ) for a call.
func callArgsRegion(source string, openParen int) (string, bool) {
	closeParen := findMatchingClose(source, openParen)
	if closeParen < 0 {
		return "", false
	}
	return source[openParen+1 : closeParen], true
}

// lineContainsWithOpen reports whether the line uses open as a with context.
func lineContainsWithOpen(line string) bool {
	t := strings.TrimSpace(line)
	if !containsWithOpenCall(t) {
		return false
	}
	return strings.Contains(t, "open(")
}

func containsWithOpenCall(t string) bool {
	// crude: "with " then later "open("
	idx := strings.Index(t, "with ")
	if idx < 0 {
		idx = strings.Index(t, "async with ")
	}
	if idx < 0 {
		return false
	}
	return strings.Contains(t[idx:], "open(")
}

func looksLikePlaceholderSecret(val string) bool {
	v := strings.ToLower(strings.TrimSpace(val))
	// strip bytes/raw/f prefixes then quotes
	for {
		if len(v) >= 2 {
			c0 := v[0]
			if c0 == 'b' || c0 == 'r' || c0 == 'u' || c0 == 'f' {
				v = v[1:]
				continue
			}
		}
		break
	}
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			v = v[1 : len(v)-1]
		}
	}
	if v == "" {
		return true
	}
	// Exact / whole-value placeholders only (conservative FP policy).
	exact := map[string]struct{}{
		"changeme": {}, "change-me": {}, "change_me": {},
		"xxx": {}, "xxxxx": {}, "todo": {}, "fix": {}, "placeholder": {},
		"secret": {}, "password": {}, "api_key": {}, "apikey": {}, "token": {},
		"your_secret": {}, "yoursecret": {}, "your_password": {}, "your_token": {},
		"example": {}, "sample": {}, "dummy": {}, "replace_me": {}, "insert_here": {},
		"notasecret": {}, "not-a-secret": {}, "not_a_secret": {},
		"***": {}, "...": {}, "none": {}, "null": {}, "n/a": {}, "na": {},
		"secret_key": {}, "supersecret": {},
	}
	if _, ok := exact[v]; ok {
		return true
	}
	// Prefix patterns for documented placeholders.
	prefixes := []string{"changeme", "your_", "replace", "todo", "example_", "dummy_", "placeholder",
		"bench-", "bench_", "test-", "test_", "testing-", "for-testing", "fake-", "fake_"}
	for _, p := range prefixes {
		if strings.HasPrefix(v, p) {
			return true
		}
	}
	if strings.Contains(v, "testing-only") || strings.Contains(v, "for-testing-only") ||
		strings.Contains(v, "not-for-production") {
		return true
	}
	return false
}

package badpractices

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/pytext"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	metaByID["BP-PY-46"] = &rules.RuleMetadata{
		ID: "BP-PY-46", Title: "print Debugging In Library Code",
		Description: "`print` is used for operational logging in non-script modules.",
		Severity:    rules.SeverityInfo, Pack: rules.PackBadPractice,
		Fix: "Use the logging module; keep print under `if __name__ == \"__main__\"` for CLIs.",
	}
	metaByID["BP-PY-47"] = &rules.RuleMetadata{
		ID: "BP-PY-47", Title: "logging With String Format Before Logger",
		Description: "Log messages are eagerly formatted with f-strings instead of lazy `%s`/`{}` args.",
		Severity:    rules.SeverityInfo, Pack: rules.PackBadPractice,
		Fix: "Prefer lazy logging: logger.info(\"x=%s\", val) so formatting is skipped when disabled.",
	}
	RegisterRule("BP-PY-46", detectBPPY46)
	RegisterRule("BP-PY-47", detectBPPY47)
}

// BP-PY-46: print( in library code (not under __main__ guard, not tests).
func detectBPPY46(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-46")
	if isPythonTestFile(unit) {
		return
	}
	if isPythonBenchmarkFile(unit) {
		return
	}
	if isPythonScriptModule(unit) {
		return
	}
	// Shebang-less standalone data scripts whose only prints sit at module
	// top level (FlashySurf data-process.py / semantic-classification.py).
	// Package modules and function-body printers keep firing.
	if pythonIsLooseStandaloneModulePrintScript(unit) {
		return
	}
	if isRequirementsPath(unit) {
		return
	}
	// Shebang alone is not an entry-script signal (vestigial #! on libraries).
	if pythonShebangIsEntryScript(unit.Source) {
		return
	}
	// Do NOT whole-module-skip on `from rich import print` — that over-suppressed
	// caniscrape upload_handler/telemetry audited TPs. Rich is only a CLI signal
	// inside pythonHasClickishCLI / presentation-name skips below.
	if !facts.has("print(") && !strings.Contains(unit.Source, "print(") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	// Masked twin: blank string/comment interiors so docstring/epilog lines do
	// not reset CLI indent and so print( tokens inside literals are ignored.
	maskedLines := buildCodeLines(pytext.Mask(unit.Source))
	mainInvoked := pythonMainGuardInvokesMain(lines)
	hasArgparse := pythonHasArgparseCLI(unit.Source)
	hasClickishCLI := pythonHasClickishCLI(unit.Source)
	registeredCmds := pythonArgparseSetDefaultsFuncs(unit.Source)
	// Track whether we are inside if __name__ == "__main__": block.
	mainIndent := -1
	cliDecorator := false
	cliIndent := -1
	for i, line := range lines {
		maskedTrim := ""
		if i < len(maskedLines) {
			maskedTrim = strings.TrimSpace(maskedLines[i].text)
		}
		// Wholly blanked lines are string/comment interiors — skip without
		// touching indent trackers (fixes argparse epilog "Examples:" resets).
		if maskedTrim == "" {
			continue
		}
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		ind := indentWidth(line.raw)
		if mainIndent >= 0 && ind <= mainIndent {
			mainIndent = -1
		}
		if cliIndent >= 0 && ind <= cliIndent {
			cliIndent = -1
		}
		if isMainGuard(t) {
			mainIndent = ind
			continue
		}
		if mainIndent >= 0 && ind > mainIndent {
			// Under main guard — skip print.
			continue
		}
		if isPythonCLIDecorator(t) {
			cliDecorator = true
			continue
		}
		if cliDecorator {
			if strings.HasPrefix(t, "@") {
				continue
			}
			if strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") {
				cliIndent = ind
				cliDecorator = false
				continue
			}
			cliDecorator = false
		}
		if name, ok := pythonDefFuncName(t); ok &&
			pythonCLIPrintSkipFunc(name, mainInvoked, hasArgparse, hasClickishCLI, registeredCmds) {
			cliIndent = ind
			continue
		}
		if cliIndent >= 0 && ind > cliIndent {
			continue
		}
		// Flag print( calls that survive masking (executable, not in a literal).
		if !strings.Contains(maskedTrim, "print(") {
			continue
		}
		if printCallOutsideString(t) {
			off := line.byte
			if i := indexOfIdent(t, "print"); i >= 0 {
				// Use offset in raw text when possible
				if j := strings.Index(line.text, "print"); j >= 0 {
					off = line.byte + j
				}
			}
			pushAt(unit, meta, off,
				"print used in library code; prefer logging (or keep print under if __name__ == \"__main__\")",
				out)
		}
	}
}

// isPythonCLIDecorator reports Click/Typer/Cyclopts/Flask CLI command decorators
// whose function bodies are user-facing presentation, not library logging.
func isPythonCLIDecorator(t string) bool {
	if !strings.HasPrefix(t, "@") {
		return false
	}
	switch {
	case strings.Contains(t, ".cli.command("),
		strings.Contains(t, "click.command("),
		strings.Contains(t, "click.group("),
		strings.Contains(t, "cli.command("),
		strings.Contains(t, "cli.group("),
		strings.Contains(t, ".command("),
		strings.Contains(t, ".default("),
		strings.Contains(t, ".callback("):
		return true
	default:
		return false
	}
}

// pythonHasArgparseCLI reports modules that build an argparse CLI entrypoint.
func pythonHasArgparseCLI(src string) bool {
	return strings.Contains(src, "ArgumentParser") ||
		strings.Contains(src, "argparse.") ||
		strings.Contains(src, "import argparse")
}

// pythonHasClickishCLI reports Click/Typer/Cyclopts/Rich-CLI presentation modules.
func pythonHasClickishCLI(src string) bool {
	return strings.Contains(src, "import click") ||
		strings.Contains(src, "from click ") ||
		strings.Contains(src, "import typer") ||
		strings.Contains(src, "from typer ") ||
		strings.Contains(src, "import cyclopts") ||
		strings.Contains(src, "from cyclopts ") ||
		pythonUsesRichPrint(src) ||
		strings.Contains(src, "from rich.console import")
}

// pythonUsesRichPrint reports modules that rebind builtin print to rich.print for UX.
func pythonUsesRichPrint(src string) bool {
	return strings.Contains(src, "from rich import print") ||
		strings.Contains(src, "from rich import print,") ||
		strings.Contains(src, "import rich.print")
}

// pythonMainGuardInvokesMain reports that if __name__ == "__main__" calls main(...).
func pythonMainGuardInvokesMain(lines []codeLine) bool {
	mainIndent := -1
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		ind := indentWidth(line.raw)
		if mainIndent >= 0 && ind <= mainIndent {
			mainIndent = -1
		}
		if isMainGuard(t) {
			mainIndent = ind
			continue
		}
		if mainIndent >= 0 && ind > mainIndent && strings.Contains(t, "main(") {
			return true
		}
	}
	return false
}

// pythonModuleDefNames collects module-level (indent-0) function names.
func pythonModuleDefNames(src string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, raw := range strings.Split(src, "\n") {
		if indentWidth(raw) != 0 {
			continue
		}
		t := strings.TrimSpace(raw)
		for _, pref := range []string{"async def ", "def "} {
			if !strings.HasPrefix(t, pref) {
				continue
			}
			if name := pythonLeadingIdent(t[len(pref):]); name != "" {
				out[name] = struct{}{}
			}
			break
		}
	}
	return out
}

// pythonCallNames extracts bare identifiers immediately followed by '(' from a
// raw line. Attribute-qualified calls (module.fn( / obj.method() are skipped;
// commented-out calls are included so demo scripts with disabled example calls
// still count as self-running.
func pythonCallNames(line string) []string {
	var out []string
	seen := make(map[string]struct{})
	for i := 0; i < len(line); i++ {
		if line[i] != '(' {
			continue
		}
		j := i - 1
		for j >= 0 && isIdentByte(line[j]) {
			j--
		}
		if j >= 0 && line[j] == '.' {
			continue // attribute call (module.fn( / obj.method()
		}
		name := line[j+1 : i]
		if name == "" || (name[0] >= '0' && name[0] <= '9') {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}

// pythonMainGuardCallsLocal reports that the module's own __main__ guard body
// calls a function defined at module level in this file — the shape of a
// self-running entry script whose prints are the program's own output.
func pythonMainGuardCallsLocal(src string) bool {
	defs := pythonModuleDefNames(src)
	if len(defs) == 0 {
		return false
	}
	mainIndent := -1
	for _, raw := range strings.Split(src, "\n") {
		t := strings.TrimSpace(raw)
		if t == "" {
			continue
		}
		ind := indentWidth(raw)
		if mainIndent >= 0 && ind <= mainIndent {
			mainIndent = -1
		}
		if isMainGuard(t) {
			mainIndent = ind
			continue
		}
		if mainIndent >= 0 && ind > mainIndent {
			for _, name := range pythonCallNames(raw) {
				if _, ok := defs[name]; ok {
					return true
				}
			}
		}
	}
	return false
}

// pythonHasMainGuard reports the presence of an if __name__ == "__main__" guard.
func pythonHasMainGuard(src string) bool {
	for _, raw := range strings.Split(src, "\n") {
		if isMainGuard(strings.TrimSpace(raw)) {
			return true
		}
	}
	return false
}

// pythonHasTopLevelPrint reports an executable print( call at module top level
// (indent 0) — the shape of a runnable script whose prints are its program
// output rather than library logging. Strings/comments are masked.
func pythonHasTopLevelPrint(src string) bool {
	for _, raw := range strings.Split(pytext.Mask(src), "\n") {
		if strings.HasPrefix(raw, "print(") {
			return true
		}
	}
	return false
}

// pythonSingleTopLevelCompletionPrint reports a dev-tool script whose only
// print call is a single top-level completion message. Any other print
// anywhere in the module disqualifies it, so operational scripts that log via
// print (WeThePeople scripts/veritas_*_patch.py) keep firing.
func pythonSingleTopLevelCompletionPrint(src string) bool {
	found := false
	for _, raw := range strings.Split(pytext.Mask(src), "\n") {
		if !strings.Contains(raw, "print(") {
			continue
		}
		if !strings.HasPrefix(raw, "print(") {
			return false
		}
		if found {
			return false
		}
		found = true
	}
	return found
}

// pythonAllPrintsModuleLevel reports that every print( call is at column 0
// (module top level). Used for shebang debug/CLI scripts under scripts/tools
// that print progress at module scope (pdf_oxide scripts/*). Function-body
// prints (WeThePeople operational scripts) return false.
func pythonAllPrintsModuleLevel(src string) bool {
	saw := false
	for _, raw := range strings.Split(pytext.Mask(src), "\n") {
		if !strings.Contains(raw, "print(") {
			continue
		}
		if !strings.HasPrefix(raw, "print(") {
			return false
		}
		saw = true
	}
	return saw
}

// pythonIsPureModuleScript reports a file with no def/class/async def — a
// linear module-level script (pdf_oxide scripts/check_span_spacing.py). Prints
// may sit under module-level if/else. WeThePeople operational scripts define
// functions and return false.
func pythonIsPureModuleScript(src string) bool {
	for _, raw := range strings.Split(pytext.Mask(src), "\n") {
		t := strings.TrimSpace(raw)
		if strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") ||
			strings.HasPrefix(t, "class ") {
			return false
		}
	}
	return strings.Contains(src, "print(")
}

// pythonArgparseSetDefaultsFuncs collects names from set_defaults(func=...).
func pythonArgparseSetDefaultsFuncs(src string) map[string]struct{} {
	out := make(map[string]struct{})
	start := 0
	for {
		idx := strings.Index(src[start:], "set_defaults(")
		if idx < 0 {
			return out
		}
		abs := start + idx + len("set_defaults(")
		end := abs
		depth := 1
		for end < len(src) && depth > 0 {
			switch src[end] {
			case '(':
				depth++
			case ')':
				depth--
			}
			if depth > 0 {
				end++
			}
		}
		if depth != 0 {
			return out
		}
		args := src[abs:end]
		if i := strings.Index(args, "func="); i >= 0 {
			name := pythonLeadingIdent(strings.TrimSpace(args[i+len("func="):]))
			if isSimpleIdent(name) {
				out[name] = struct{}{}
			}
		}
		start = end + 1
		if start >= len(src) {
			return out
		}
	}
}

func pythonDefFuncName(t string) (string, bool) {
	for _, pref := range []string{"async def ", "def "} {
		if !strings.HasPrefix(t, pref) {
			continue
		}
		name := pythonLeadingIdent(t[len(pref):])
		if name == "" {
			return "", false
		}
		return name, true
	}
	return "", false
}

func pythonLeadingIdent(s string) string {
	end := 0
	for end < len(s) && isIdentByte(s[end]) {
		end++
	}
	if end == 0 {
		return ""
	}
	return s[:end]
}

// pythonCLIPrintSkipFunc reports argparse/__main__/Click CLI presentation
// entrypoints whose print calls are intentional user output, not library debug logging.
func pythonCLIPrintSkipFunc(
	name string,
	mainInvoked, hasArgparse, hasClickishCLI bool,
	registeredCmds map[string]struct{},
) bool {
	if mainInvoked && name == "main" {
		return true
	}
	if _, ok := registeredCmds[name]; ok {
		return true
	}
	if mainInvoked {
		if strings.HasPrefix(name, "print_") || strings.HasPrefix(name, "cmd_") ||
			strings.HasPrefix(name, "show_") {
			return true
		}
	}
	if hasArgparse && mainInvoked {
		if strings.HasPrefix(name, "run_") {
			return true
		}
	}
	// Click-style helpers: init_command / push_command called from @cli.command wrappers.
	if hasClickishCLI && strings.HasSuffix(name, "_command") {
		return true
	}
	// Presentation helpers under a CLI framework (rich/click/typer/cyclopts):
	// prompt/show/enable/disable/display/request_ are user-facing UX flow, not
	// library operational logging. Worker printers (contribute_*/try_*/upload_*)
	// stay reportable (caniscrape contribute_scan / try_upload_scan TPs).
	if hasClickishCLI {
		for _, prefix := range []string{
			"prompt_", "show_", "enable_", "disable_", "display_",
			"request_", "confirm_", "ask_", "warn_if_",
		} {
			if strings.HasPrefix(name, prefix) {
				return true
			}
		}
	}
	return false
}

func isMainGuard(t string) bool {
	// if __name__ == "__main__":  and variants with single quotes / spaces
	if !strings.HasPrefix(t, "if ") {
		return false
	}
	if !strings.Contains(t, "__name__") {
		return false
	}
	if !strings.Contains(t, "__main__") {
		return false
	}
	return strings.HasSuffix(strings.TrimSpace(t), ":") || strings.Contains(t, ":")
}

func printCallOutsideString(line string) bool {
	inStr := byte(0)
	escape := false
	for i := 0; i < len(line); i++ {
		c := line[i]
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
		if c == '"' || c == '\'' {
			inStr = c
			continue
		}
		if c == 'p' && i+6 <= len(line) && line[i:i+6] == "print(" {
			if i == 0 || (line[i-1] != '.' && !isIdentByte(line[i-1])) {
				return true
			}
		}
	}
	return false
}

// BP-PY-47: logger.*/logging.* with eager f-string or .format( as first arg.
func detectBPPY47(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-47")
	if isRequirementsPath(unit) {
		return
	}
	// Eager format in tests is not a production logging-hot-path concern
	// (logxide test_* f-string message construction FPs). Production modules
	// and examples keep firing (CourtScrapper / logxide examples TPs).
	if isPythonTestFile(unit) {
		return
	}
	if !facts.hasAny("logger.", "logging.") &&
		!strings.Contains(unit.Source, "logger.") && !strings.Contains(unit.Source, "logging.") {
		return
	}
	src := unit.Source
	methods := []string{
		"debug", "info", "warning", "error", "critical", "exception", "log",
	}
	// Patterns: logger.METHOD( or logging.METHOD(
	prefixes := []string{"logger.", "logging."}
	for _, pref := range prefixes {
		for _, m := range methods {
			needle := pref + m + "("
			start := 0
			for {
				idx := strings.Index(src[start:], needle)
				if idx < 0 {
					break
				}
				abs := start + idx
				open := abs + len(needle) - 1 // '('
				if open >= len(src) || src[open] != '(' {
					start = abs + len(needle)
					continue
				}
				arg, ok := firstCallArg(src, open)
				if !ok {
					start = abs + len(needle)
					continue
				}
				if isEagerLogFormat(arg) {
					pushAt(unit, meta, abs,
						"log message eagerly formatted (f-string/.format); prefer lazy %s args so work is skipped when disabled",
						out)
				}
				start = abs + len(needle)
			}
		}
	}
}

func isEagerLogFormat(arg string) bool {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return false
	}
	// f-string / rf / fr prefixes
	lower := arg
	// Strip b/u not relevant; f/F and combinations
	a := arg
	// Common f-string prefixes
	for _, p := range []string{"f", "F", "rf", "Rf", "rF", "RF", "fr", "Fr", "fR", "FR"} {
		if strings.HasPrefix(a, p+"\"") || strings.HasPrefix(a, p+"'") {
			// All f-strings are flagged — including constant f"hello" without
			// braces (among-llms audited TPs). Lazy logging still prefers %s.
			return true
		}
	}
	// .format( on the first arg expression: "...".format( or 'x{}'.format(
	if strings.Contains(arg, ".format(") {
		return true
	}
	// % formatting done before call: ("x %s" % val) — paren wrapped
	// Conservative: only flag if top-level % between quotes and rest.
	if isPercentFormattedArg(arg) {
		return true
	}
	_ = lower
	return false
}

// fStringHasInterpolation reports an f/rf/fr-string first arg that actually
// interpolates ({...}). Constant f"hello" is not eager formatting work.
func fStringHasInterpolation(arg string) bool {
	arg = strings.TrimSpace(arg)
	// Strip prefix to the opening quote.
	i := 0
	for i < len(arg) && (arg[i] == 'f' || arg[i] == 'F' || arg[i] == 'r' || arg[i] == 'R' ||
		arg[i] == 'b' || arg[i] == 'B' || arg[i] == 'u' || arg[i] == 'U') {
		i++
	}
	if i >= len(arg) || (arg[i] != '"' && arg[i] != '\'') {
		return true // not a plain string; be conservative
	}
	quote := arg[i]
	// Triple?
	if i+2 < len(arg) && arg[i+1] == quote && arg[i+2] == quote {
		body := arg[i+3:]
		if end := strings.Index(body, string([]byte{quote, quote, quote})); end >= 0 {
			body = body[:end]
		}
		return strings.Contains(body, "{")
	}
	body := arg[i+1:]
	escape := false
	for j := 0; j < len(body); j++ {
		c := body[j]
		if escape {
			escape = false
			continue
		}
		if c == '\\' {
			escape = true
			continue
		}
		if c == quote {
			return false
		}
		if c == '{' {
			// {{ is escaped brace, not interpolation
			if j+1 < len(body) && body[j+1] == '{' {
				j++
				continue
			}
			return true
		}
	}
	return false
}

// pythonIsLooseStandaloneModulePrintScript reports a non-package file whose
// every print( is at column 0 (module top level). Covers shebang-less one-shot
// data scripts (FlashySurf data-process.py) without silencing package modules
// or function-body printers (Cronboard / WHEN-Language / WeThePeople).
func pythonIsLooseStandaloneModulePrintScript(unit *core.ParsedUnit) bool {
	if unit == nil || unit.Source == "" {
		return false
	}
	if !pythonAllPrintsModuleLevel(unit.Source) {
		return false
	}
	// Must actually print something.
	if !strings.Contains(unit.Source, "print(") {
		return false
	}
	return pythonPathOutsidePackage(unit)
}

// pythonPathOutsidePackage reports that the unit does not live inside a Python
// package (no __init__.py in the file's directory or parent directories up to a
// few hops). Repo-root scripts qualify; src/pkg/module.py does not.
func pythonPathOutsidePackage(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	path := unit.Path
	if path == "" {
		path = fileDisplayPath(unit)
	}
	if path == "" {
		return false
	}
	dir := filepath.Dir(path)
	for hops := 0; hops < 8; hops++ {
		if dir == "" || dir == "." || dir == string(filepath.Separator) {
			break
		}
		initPath := filepath.Join(dir, "__init__.py")
		if st, err := os.Stat(initPath); err == nil && !st.IsDir() {
			return false
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		// Stop at common project roots so we do not walk into unrelated trees.
		base := filepath.Base(dir)
		if base == "src" || base == "lib" || base == "site-packages" {
			dir = parent
			continue
		}
		dir = parent
	}
	return true
}

func isPercentFormattedArg(arg string) bool {
	// ("msg %s" % x) or 'msg %s' % x
	// Not a plain "msg %s" literal (lazy-style first arg without %).
	depth := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(arg); i++ {
		c := arg[i]
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
		case '%':
			if depth == 0 {
				// binary % operator outside string
				// ensure not %% and not start of f-string nonsense
				if i+1 < len(arg) && arg[i+1] == '%' {
					i++
					continue
				}
				// Must have something after that looks like RHS
				rest := strings.TrimSpace(arg[i+1:])
				if rest != "" {
					return true
				}
			}
		}
	}
	return false
}

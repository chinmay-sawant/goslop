package cwe

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/pytext"
)

const (
	confidence65 = 0.65
	confidence68 = 0.68
	confidence70 = 0.70
	confidence72 = 0.72
	confidence74 = 0.74
	confidence75 = 0.75
	confidence76 = 0.76
	confidence78 = 0.78
	confidence80 = 0.80
	confidence82 = 0.82
	confidence84 = 0.84
	confidence85 = 0.85
	confidence86 = 0.86
	confidence88 = 0.88
	confidence90 = 0.90
	confidence92 = 0.92

	maxInsecureTokenBytes       = 15
	maxSmallRandomTokenValue    = 9999
	userLookupContextWindow     = 400
	uninitializedResourceWindow = 220
	closedResourceWindow        = 180
	minimumRouteBranches        = 12
	pathDivisionContextWindow   = 80
)

// unitFile returns the display path for findings.
func unitFile(unit *core.ParsedUnit) string {
	if unit == nil {
		return ""
	}
	if unit.DisplayPath != "" {
		return unit.DisplayPath
	}
	return unit.Path
}

// isPythonTestModule identifies test modules so source-only rules can avoid
// treating deliberate test fixtures and assertions as deployed code.
func isPythonTestModule(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, path := range []string{unit.DisplayPath, unit.Path} {
		if path == "" {
			continue
		}
		base := filepath.Base(path)
		if base == "tests.py" || base == "conftest.py" ||
			(strings.HasPrefix(base, "test_") && strings.HasSuffix(base, ".py")) ||
			strings.HasSuffix(base, "_test.py") {
			return true
		}
		normalized := filepath.ToSlash(path)
		if strings.Contains(normalized, "/tests/") || strings.Contains(normalized, "/test/") ||
			strings.HasPrefix(normalized, "tests/") || strings.HasPrefix(normalized, "test/") {
			return true
		}
	}
	return false
}

// isPythonBenchmarkFile identifies harness code whose literals and console
// output are synthetic benchmark data rather than deployed application code.
// Match path components only so a project name containing "bench" is not
// enough to suppress a finding.
func isPythonBenchmarkFile(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, path := range []string{unit.DisplayPath, unit.Path} {
		for _, component := range strings.Split(strings.Trim(filepath.ToSlash(path), "/"), "/") {
			if component == "bench" || component == "benchmarks" {
				return true
			}
		}
	}
	return false
}

// callSite is a lightweight function/method call match in source text.
type callSite struct {
	Name     string // callee text matched (e.g. "pickle.loads")
	Start    int    // byte offset of Name
	End      int    // past-the-end of matching call
	ArgsText string // interior of top-level argument list
}

// findCalls finds call-like sites using facts.Masked when source is the full
// file. Prefer findCallsMasked with facts.codeMask for known fragment offsets.
func findCalls(facts *PyCweFacts, source string, names ...string) []callSite {
	var masked string
	if facts != nil {
		masked = facts.codeMask(source, fragStartHint(facts, source))
	} else {
		masked = pythonCodeMask(source)
	}
	return findCallsMasked(source, masked, names...)
}

// findCallsMasked finds call-like sites for each name using pre-masked text
// (byte-aligned with source). Does not remask.
func findCallsMasked(source, masked string, names ...string) []callSite {
	var out []callSite
	if masked == "" {
		masked = pythonCodeMask(source)
	}
	for _, name := range names {
		start := 0
		for {
			idx := strings.Index(masked[start:], name)
			if idx < 0 {
				break
			}
			abs := start + idx
			if !identBoundaryOK(masked, abs, abs+len(name)) {
				start = abs + len(name)
				continue
			}
			after := abs + len(name)
			j := skipWS(masked, after)
			if j >= len(masked) || masked[j] != '(' {
				start = after
				continue
			}
			closeAt, args := scanCallArgs(source, j)
			if closeAt < 0 {
				start = after
				continue
			}
			out = append(out, callSite{
				Name:     name,
				Start:    abs,
				End:      closeAt + 1,
				ArgsText: args,
			})
			start = closeAt + 1
		}
	}
	return out
}

// fragStartHint returns 0 when source is the full facts.Source, else -1
// (unknown offset → codeMask may remask the fragment).
func fragStartHint(facts *PyCweFacts, source string) int {
	if facts != nil && source == facts.Source {
		return 0
	}
	return -1
}

// pythonCodeMask keeps byte offsets stable while blanking comments and string
// literals. Prefer facts.Masked / facts.codeMask; this remasks and is the
// fallback when no facts or fragment offset is available.
func pythonCodeMask(source string) string {
	return pytext.Mask(source)
}

// firstCodeMatchStart returns the start of the first pattern match on masked
// code text (comments/strings already blanked). Prefer this for code-oriented
// patterns; use firstLiteralMatchStart when the pattern must see string quotes.
func firstCodeMatchStart(facts *PyCweFacts, source string, pattern *regexp.Regexp) int {
	if pattern == nil || source == "" {
		return -1
	}
	var masked string
	if facts != nil {
		masked = facts.codeMask(source, fragStartHint(facts, source))
	} else {
		masked = pythonCodeMask(source)
	}
	return firstCodeMatchStartMasked(source, masked, pattern)
}

// firstCodeMatchStartMasked runs a single FindStringIndex on masked text.
// masked must be byte-aligned with source. Hits are already "code" hits.
func firstCodeMatchStartMasked(source, masked string, pattern *regexp.Regexp) int {
	if pattern == nil {
		return -1
	}
	if masked == "" {
		if source == "" {
			return -1
		}
		masked = pythonCodeMask(source)
	}
	loc := pattern.FindStringIndex(masked)
	if loc == nil {
		return -1
	}
	return loc[0]
}

// firstLiteralMatchStart scans unmasked source with pattern and keeps the first
// hit whose masked span is non-blank. Use when the RE must see string-literal
// quotes or contents (hard-coded credentials, quoted config keys, etc.).
func firstLiteralMatchStart(facts *PyCweFacts, source string, pattern *regexp.Regexp) int {
	if pattern == nil || source == "" {
		return -1
	}
	var masked string
	if facts != nil {
		masked = facts.codeMask(source, fragStartHint(facts, source))
	} else {
		masked = pythonCodeMask(source)
	}
	return firstLiteralMatchStartMasked(source, masked, pattern)
}

// firstLiteralMatchStartMasked is the source-then-filter path for literal REs.
func firstLiteralMatchStartMasked(source, masked string, pattern *regexp.Regexp) int {
	if pattern == nil || source == "" {
		return -1
	}
	if masked == "" {
		masked = pythonCodeMask(source)
	}
	search := 0
	for search <= len(source) {
		loc := pattern.FindStringIndex(source[search:])
		if loc == nil {
			return -1
		}
		start := search + loc[0]
		end := search + loc[1]
		if end > len(masked) {
			end = len(masked)
		}
		if start < end && strings.TrimSpace(masked[start:end]) != "" {
			return start
		}
		if end <= search {
			search++
			continue
		}
		search = end
	}
	return -1
}

// containsAnyNeedle reports whether any non-empty needle appears in s.
func containsAnyNeedle(s string, needles ...string) bool {
	for _, n := range needles {
		if n != "" && strings.Contains(s, n) {
			return true
		}
	}
	return false
}

// eachLiteralMatch invokes fn for each source match whose masked span is
// non-blank. fn returns false to stop. Prefer over FindAllStringIndex.
func eachLiteralMatch(facts *PyCweFacts, source string, pattern *regexp.Regexp, fn func(start, end int) bool) {
	if pattern == nil || source == "" || fn == nil {
		return
	}
	var masked string
	if facts != nil {
		masked = facts.codeMask(source, fragStartHint(facts, source))
	} else {
		masked = pythonCodeMask(source)
	}
	if masked == "" {
		masked = pythonCodeMask(source)
	}
	search := 0
	for search <= len(source) {
		loc := pattern.FindStringIndex(source[search:])
		if loc == nil {
			return
		}
		start := search + loc[0]
		end := search + loc[1]
		maskedEnd := end
		if maskedEnd > len(masked) {
			maskedEnd = len(masked)
		}
		if start < maskedEnd && strings.TrimSpace(masked[start:maskedEnd]) != "" {
			if !fn(start, end) {
				return
			}
		}
		if end <= search {
			search++
			continue
		}
		search = end
	}
}

// eachCodeMatch invokes fn for each FindStringIndex hit on masked text.
// fn returns false to stop. Prefer over FindAllStringIndex on masked.
func eachCodeMatch(masked string, pattern *regexp.Regexp, fn func(start, end int) bool) {
	if pattern == nil || masked == "" || fn == nil {
		return
	}
	search := 0
	for search <= len(masked) {
		loc := pattern.FindStringIndex(masked[search:])
		if loc == nil {
			return
		}
		start := search + loc[0]
		end := search + loc[1]
		if !fn(start, end) {
			return
		}
		if end <= search {
			search++
			continue
		}
		search = end
	}
}

func identBoundaryOK(source string, start, end int) bool {
	// Method-form needles like ".execute" intentionally start with '.' so the
	// left side may be an identifier (cursor.execute). Only enforce the left
	// boundary for names that do not begin with '.'.
	if start > 0 && (start >= len(source) || source[start] != '.') {
		r, _ := utf8.DecodeLastRuneInString(source[:start])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' {
			return false
		}
	}
	if end < len(source) {
		r, _ := utf8.DecodeRuneInString(source[end:])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return false
		}
	}
	return true
}

func skipWS(s string, i int) int {
	for i < len(s) {
		switch s[i] {
		case ' ', '\t', '\n', '\r':
			i++
		default:
			return i
		}
	}
	return i
}

// scanCallArgs starts at '(' and returns index of matching ')' and interior args.
func scanCallArgs(source string, open int) (int, string) {
	if open >= len(source) || source[open] != '(' {
		return -1, ""
	}
	depth := 0
	inStr := byte(0)
	escape := false
	// Python triple-quote tracking: """ or '''
	triple := 0 // 0 none, 3-ish when opening triple
	for i := open; i < len(source); i++ {
		c := source[i]
		if inStr != 0 {
			i, inStr, triple, escape = advancePythonQuotedString(source, i, inStr, triple, escape)
			continue
		}
		// not in string
		if c == '"' || c == '\'' {
			// check for triple quote
			if i+2 < len(source) && source[i+1] == c && source[i+2] == c {
				inStr = c
				triple = 3
				i += 2
				continue
			}
			inStr = c
			continue
		}
		switch c {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i, source[open+1 : i]
			}
		}
	}
	return -1, ""
}

// splitTopLevelArgs splits a call argument list on top-level commas.
func splitTopLevelArgs(args string) []string {
	var out []string
	depth := 0
	inStr := byte(0)
	escape := false
	triple := 0
	start := 0
	for i := 0; i < len(args); i++ {
		c := args[i]
		if inStr != 0 {
			i, inStr, triple, escape = advancePythonQuotedString(args, i, inStr, triple, escape)
			continue
		}
		if c == '"' || c == '\'' {
			if i+2 < len(args) && args[i+1] == c && args[i+2] == c {
				inStr = c
				triple = 3
				i += 2
				continue
			}
			inStr = c
			continue
		}
		switch c {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(args[start:i]))
				start = i + 1
			}
		}
	}
	rest := strings.TrimSpace(args[start:])
	if rest != "" || len(out) > 0 {
		out = append(out, rest)
	}
	return out
}

// isPureStringLiteral reports whether expr is exactly one Python string literal
// (optional string prefix b/r/u/f/fr… then ' or " or triple quotes). Pure f-strings
// are treated as dynamic (not pure) because they embed expressions.
func isPureStringLiteral(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	// Strip string prefixes (order-insensitive common set).
	i := 0
	hasF := false
	for i < len(t) {
		c := t[i]
		if c == 'f' || c == 'F' {
			hasF = true
			i++
			continue
		}
		if c == 'b' || c == 'B' || c == 'r' || c == 'R' || c == 'u' || c == 'U' {
			i++
			continue
		}
		break
	}
	if hasF {
		// f-strings embed expressions → dynamic for injection purposes
		return false
	}
	if i >= len(t) {
		return false
	}
	quote := t[i]
	if quote != '"' && quote != '\'' {
		return false
	}
	// triple?
	if i+2 < len(t) && t[i+1] == quote && t[i+2] == quote {
		end := strings.Index(t[i+3:], string([]byte{quote, quote, quote}))
		if end < 0 {
			return false
		}
		return strings.TrimSpace(t[i+3+end+3:]) == ""
	}
	// single-quoted string with escapes
	escape := false
	for j := i + 1; j < len(t); j++ {
		if escape {
			escape = false
			continue
		}
		if t[j] == '\\' {
			escape = true
			continue
		}
		if t[j] == quote {
			return strings.TrimSpace(t[j+1:]) == ""
		}
	}
	return false
}

// isDynamicExpr reports whether expr looks non-constant (variable, format, concat, f-string).
func isDynamicExpr(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if isPureStringLiteral(t) {
		return false
	}
	// list/tuple of pure string literals is still "static argv"
	if looksStaticStringList(t) {
		return false
	}
	return true
}

// looksStaticStringList is true for ["a","b"] / ('a', 'b') of pure literals only.
func looksStaticStringList(expr string) bool {
	t := strings.TrimSpace(expr)
	if len(t) < 2 {
		return false
	}
	var closeDelimiter byte
	switch t[0] {
	case '[':
		closeDelimiter = ']'
	case '(':
		closeDelimiter = ')'
	default:
		return false
	}
	if t[len(t)-1] != closeDelimiter {
		return false
	}
	inner := strings.TrimSpace(t[1 : len(t)-1])
	if inner == "" {
		return true
	}
	parts := splitTopLevelArgs(inner)
	if len(parts) == 0 {
		return false
	}
	for _, p := range parts {
		if !isPureStringLiteral(p) {
			return false
		}
	}
	return true
}

// hasKwargTrue reports whether argsText contains an exact name=True keyword.
func hasKwargTrue(argsText, name string) bool {
	return hasBooleanKwarg(argsText, name, "True")
}

// hasKwargFalse reports an exact name=False keyword.
func hasKwargFalse(argsText, name string) bool {
	return hasBooleanKwarg(argsText, name, "False")
}

func hasBooleanKwarg(argsText, name, value string) bool {
	for _, arg := range splitTopLevelArgs(argsText) {
		key, candidate, ok := strings.Cut(arg, "=")
		if !ok || strings.TrimSpace(key) != name {
			continue
		}
		return strings.TrimSpace(candidate) == value
	}
	return false
}

func compactWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// containsIdent reports whole-identifier presence of ident in text.
func containsIdent(text, ident string) bool {
	if ident == "" {
		return false
	}
	start := 0
	for {
		idx := strings.Index(text[start:], ident)
		if idx < 0 {
			return false
		}
		abs := start + idx
		okLeft := true
		if abs > 0 {
			r, _ := utf8.DecodeLastRuneInString(text[:abs])
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				okLeft = false
			}
		}
		okRight := true
		end := abs + len(ident)
		if end < len(text) {
			r, _ := utf8.DecodeRuneInString(text[end:])
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				okRight = false
			}
		}
		if okLeft && okRight {
			return true
		}
		start = abs + len(ident)
	}
}

// looksSQLFormatted reports f-string / % / .format style string construction in expr.
func looksSQLFormatted(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	// f"SELECT ... {x}" or F'...'
	if isFStringExpr(t) {
		return true
	}
	// "... %s ..." % var  or  "...".format(...)
	if strings.Contains(t, ".format(") {
		return true
	}
	// binary % formatting: lhs % rhs (not modulo of numbers only)
	if idx := indexTopLevelPercent(t); idx > 0 {
		lhs := strings.TrimSpace(t[:idx])
		// string % something
		if isPureStringLiteral(lhs) || isFStringExpr(lhs) || strings.HasPrefix(strings.TrimSpace(lhs), "\"") || strings.HasPrefix(strings.TrimSpace(lhs), "'") {
			return true
		}
		// parenthesized concat etc.
		if strings.Contains(lhs, "SELECT") || strings.Contains(lhs, "INSERT") ||
			strings.Contains(lhs, "UPDATE") || strings.Contains(lhs, "DELETE") ||
			strings.Contains(lhs, "select") || strings.Contains(lhs, "insert") {
			return true
		}
	}
	// "SELECT " + var
	if strings.Contains(t, "+") && !isPureStringLiteral(t) {
		if strings.Contains(strings.ToUpper(t), "SELECT") ||
			strings.Contains(strings.ToUpper(t), "INSERT") ||
			strings.Contains(strings.ToUpper(t), "UPDATE") ||
			strings.Contains(strings.ToUpper(t), "DELETE") ||
			strings.Contains(t, "\"") || strings.Contains(t, "'") {
			// dynamic concat of SQL-ish pieces
			if !isPureStringLiteral(t) {
				return true
			}
		}
	}
	return false
}

func isFStringExpr(expr string) bool {
	t := strings.TrimSpace(expr)
	// leading prefixes may include f/F among b/r/u
	i := 0
	hasF := false
	for i < len(t) {
		c := t[i]
		if c == 'f' || c == 'F' {
			hasF = true
			i++
			continue
		}
		if c == 'b' || c == 'B' || c == 'r' || c == 'R' || c == 'u' || c == 'U' {
			i++
			continue
		}
		break
	}
	if !hasF || i >= len(t) {
		return false
	}
	return t[i] == '"' || t[i] == '\''
}

// indexTopLevelPercent finds first '%' not inside strings/parens at depth 0.
func indexTopLevelPercent(s string) int {
	depth := 0
	inStr := byte(0)
	escape := false
	triple := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			i, inStr, triple, escape = advancePythonQuotedString(s, i, inStr, triple, escape)
			continue
		}
		if c == '"' || c == '\'' {
			if i+2 < len(s) && s[i+1] == c && s[i+2] == c {
				inStr = c
				triple = 3
				i += 2
				continue
			}
			inStr = c
			continue
		}
		switch c {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '%':
			if depth == 0 {
				// skip %=
				if i+1 < len(s) && s[i+1] == '=' {
					continue
				}
				return i
			}
		}
	}
	return -1
}

func advancePythonQuotedString(source string, index int, quote byte, triple int, escaped bool) (int, byte, int, bool) {
	current := source[index]
	if triple > 0 {
		if current == quote && index+2 < len(source) && source[index+1] == quote && source[index+2] == quote {
			return index + 2, 0, 0, false
		}
		return index, quote, triple, false
	}
	if escaped {
		return index, quote, 0, false
	}
	if current == '\\' {
		return index, quote, 0, true
	}
	if current == quote {
		return index, 0, 0, false
	}
	return index, quote, 0, false
}

// hasSafePathConfinement reports file-level safe patterns for CWE-22.
// Documented suppressions: basename-only policy, or resolve + startswith root check.
func hasSafePathConfinement(src string) bool {
	if strings.Contains(src, "os.path.basename(") || strings.Contains(src, "pathlib.Path(...).name") {
		return true
	}
	// resolve + startswith / is_relative_to confinement
	hasResolve := strings.Contains(src, ".resolve(") || strings.Contains(src, "os.path.realpath(") || strings.Contains(src, "os.path.abspath(")
	hasPrefix := strings.Contains(src, "startswith(") || strings.Contains(src, "is_relative_to(") || strings.Contains(src, "commonpath(")
	return hasResolve && hasPrefix
}

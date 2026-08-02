package cwe

import (
	"regexp"
	"strings"
	"unicode"
)

// False-positive guard helpers shared by CWE-117 / 1341 / 367 / 88 / 93 / 215 / 502.

var (
	fStringInterpRE    = regexp.MustCompile(`\{([^{][^}]*)\}`)
	upperSnakeRE       = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	closeCallRE        = regexp.MustCompile(`(?is)(\w+)\.close\s*\(\s*\)`)
	toctouExistsCallRE = regexp.MustCompile(`(?is)os\.path\.(?:exists|lexists)\s*\(\s*(\w+)\s*\)`)
	toctouUseCallRE    = regexp.MustCompile(`(?is)(?:open|os\.remove|os\.unlink)\s*\(\s*(\w+)\b`)
)

// logMessageHasCRLFCapableValue is true when a formatted log argument can carry
// attacker-controlled CR/LF. Pure constants, len(...), and numeric casts cannot.
func logMessageHasCRLFCapableValue(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if isFStringExpr(t) {
		matches := fStringInterpRE.FindAllStringSubmatch(t, -1)
		if len(matches) == 0 {
			return false
		}
		for _, m := range matches {
			if !logInterpLooksCRLFSafe(m[1]) {
				return true
			}
		}
		return false
	}
	if strings.Contains(t, ".format(") {
		open := strings.Index(t, ".format(")
		inner, ok := callArgsRegion(t, open+len(".format"))
		if !ok {
			return true
		}
		for _, arg := range splitTopLevelArgs(inner) {
			if !logInterpLooksCRLFSafe(arg) {
				return true
			}
		}
		return false
	}
	if idx := indexTopLevelPercent(t); idx > 0 {
		rhs := strings.TrimSpace(t[idx+1:])
		if strings.HasPrefix(rhs, "(") && strings.HasSuffix(rhs, ")") {
			rhs = strings.TrimSpace(rhs[1 : len(rhs)-1])
		}
		for _, arg := range splitTopLevelArgs(rhs) {
			if !logInterpLooksCRLFSafe(arg) {
				return true
			}
		}
		return false
	}
	if strings.Contains(t, "+") {
		for _, part := range splitTopLevelConcat(t) {
			part = strings.TrimSpace(part)
			if isPureStringLiteral(part) || logInterpLooksCRLFSafe(part) {
				continue
			}
			return true
		}
		return false
	}
	return isDynamicExpr(t)
}

func logInterpLooksCRLFSafe(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return true
	}
	// Drop f-string conversion/format specs: {elapsed:.2f} / {n!r}
	if i := strings.IndexAny(t, ":!"); i >= 0 {
		t = strings.TrimSpace(t[:i])
	}
	if t == "" {
		return true
	}
	if isPureStringLiteral(t) || isNumericLiteral(t) {
		return true
	}
	compact := compactWhitespace(t)
	if strings.HasPrefix(compact, "len(") && strings.HasSuffix(compact, ")") {
		return true
	}
	if headerValueIsInternalNumeric(compact) {
		return true
	}
	if looksConstantStringRepeat(compact) {
		return true
	}
	if strings.HasPrefix(compact, "str(") && strings.HasSuffix(compact, ")") {
		inner := strings.TrimSpace(compact[4 : len(compact)-1])
		return logInterpLooksCRLFSafe(inner)
	}
	if strings.HasPrefix(compact, "int(") || strings.HasPrefix(compact, "float(") ||
		strings.HasPrefix(compact, "round(") {
		return true
	}
	return false
}

func looksConstantStringRepeat(expr string) bool {
	t := strings.TrimSpace(expr)
	if !strings.Contains(t, "*") {
		return false
	}
	parts := strings.SplitN(t, "*", 2)
	if len(parts) != 2 {
		return false
	}
	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])
	return (isPureStringLiteral(left) && isNumericLiteral(right)) ||
		(isNumericLiteral(left) && isPureStringLiteral(right))
}

func isNumericLiteral(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if strings.HasPrefix(t, "-") || strings.HasPrefix(t, "+") {
		t = strings.TrimSpace(t[1:])
	}
	if t == "" {
		return false
	}
	dot := 0
	for _, r := range t {
		if r == '.' {
			dot++
			if dot > 1 {
				return false
			}
			continue
		}
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// callArgsRegion returns the interior of a call whose '(' is at openIdx.
func callArgsRegion(src string, openIdx int) (string, bool) {
	if openIdx < 0 || openIdx >= len(src) || src[openIdx] != '(' {
		return "", false
	}
	closeAt, args := scanCallArgs(src, openIdx)
	if closeAt < 0 {
		return "", false
	}
	return args, true
}

// sameHandleDoubleCloseStart reports the first double-close of the same
// receiver within one function body. Lifecycle hooks (__exit__/__del__) and
// distinct receivers are ignored.
func sameHandleDoubleCloseStart(masked string) int {
	if masked == "" {
		return -1
	}
	matches := closeCallRE.FindAllStringSubmatchIndex(masked, -1)
	for i := 0; i < len(matches); i++ {
		recv1 := masked[matches[i][2]:matches[i][3]]
		end1 := matches[i][1]
		for j := i + 1; j < len(matches); j++ {
			start2 := matches[j][0]
			if start2-end1 > closedResourceWindow {
				break
			}
			recv2 := masked[matches[j][2]:matches[j][3]]
			if recv1 != recv2 {
				continue
			}
			between := masked[end1:start2]
			if strings.Contains(between, "\ndef ") || strings.Contains(between, "\nasync def ") ||
				strings.Contains(between, "\nclass ") {
				continue
			}
			if closeSiteInLifecycleHook(masked, matches[i][0]) || closeSiteInLifecycleHook(masked, start2) {
				continue
			}
			return matches[i][0]
		}
	}
	return -1
}

func closeSiteInLifecycleHook(src string, offset int) bool {
	if offset < 0 || offset > len(src) {
		return false
	}
	head := src[:offset]
	defAt := strings.LastIndex(head, "\ndef ")
	asyncAt := strings.LastIndex(head, "\nasync def ")
	at := defAt
	prefix := "\ndef "
	if asyncAt > defAt {
		at = asyncAt
		prefix = "\nasync def "
	}
	if at < 0 {
		// module-top def without leading newline
		if strings.HasPrefix(strings.TrimLeft(head, " \t"), "def ") {
			at = strings.Index(head, "def ")
			prefix = "def "
		} else {
			return false
		}
	} else {
		at++ // skip newline for \ndef
		prefix = strings.TrimPrefix(prefix, "\n")
	}
	rest := head[at:]
	if !strings.HasPrefix(rest, prefix) && !strings.HasPrefix(rest, "def ") && !strings.HasPrefix(rest, "async def ") {
		rest = src[at:]
	}
	lineEnd := strings.IndexByte(rest, '\n')
	if lineEnd < 0 {
		lineEnd = len(rest)
	}
	header := rest[:lineEnd]
	return strings.Contains(header, "def __exit__") || strings.Contains(header, "def __del__") ||
		strings.Contains(header, "async def __aexit__")
}

func toctouSamePathStart(facts *PyCweFacts, source string) int {
	if source == "" {
		return -1
	}
	var masked string
	if facts != nil {
		masked = facts.codeMask(source, fragStartHint(facts, source))
	} else {
		masked = pythonCodeMask(source)
	}
	exists := toctouExistsCallRE.FindAllStringSubmatchIndex(masked, -1)
	uses := toctouUseCallRE.FindAllStringSubmatchIndex(masked, -1)
	for _, ex := range exists {
		pathName := masked[ex[2]:ex[3]]
		exEnd := ex[1]
		for _, use := range uses {
			if use[0] < exEnd {
				continue
			}
			const toctouWindowBytes = 300
			if use[0]-exEnd > toctouWindowBytes {
				break
			}
			useName := masked[use[2]:use[3]]
			if useName == pathName {
				return ex[0]
			}
		}
	}
	return -1
}

// splitAssignmentEq splits a true assignment (=), rejecting ==, !=, <=, >=, :=.
func splitAssignmentEq(line string) (string, string, bool) {
	inStr := byte(0)
	esc := false
	triple := false
	depth := 0
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inStr != 0 {
			i, inStr, esc, triple = advanceInString(line, i, inStr, esc, triple)
			continue
		}
		switch c {
		case '\'', '"':
			inStr = c
			if i+2 < len(line) && line[i+1] == c && line[i+2] == c {
				triple = true
				i += 2
			}
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '=':
			if depth != 0 {
				continue
			}
			if i > 0 {
				prev := line[i-1]
				if prev == '=' || prev == '!' || prev == '<' || prev == '>' || prev == ':' {
					continue
				}
			}
			if i+1 < len(line) && line[i+1] == '=' {
				continue
			}
			return line[:i], line[i+1:], true
		}
	}
	return "", "", false
}

// advanceInString steps one byte inside a string literal scanner state.
func advanceInString(line string, i int, inStr byte, esc, triple bool) (int, byte, bool, bool) {
	c := line[i]
	if triple {
		if c == inStr && i+2 < len(line) && line[i+1] == inStr && line[i+2] == inStr {
			return i + 2, 0, false, false
		}
		return i, inStr, esc, triple
	}
	if esc {
		return i, inStr, false, triple
	}
	if c == '\\' {
		return i, inStr, true, triple
	}
	if c == inStr {
		return i, 0, false, triple
	}
	return i, inStr, esc, triple
}

// sensitiveIdentOutsideLiterals reports a sensitive identifier that is not only
// present as an English word inside a string/f-string literal.
func sensitiveIdentOutsideLiterals(args string) bool {
	if args == "" {
		return false
	}
	masked := pythonCodeMask(args)
	// Keep f-string interpolations visible: Mask blanks the whole literal,
	// including {expr}. Re-inject {...} bodies from the original.
	if isFStringExpr(strings.TrimSpace(args)) || strings.Contains(args, "f\"") || strings.Contains(args, "f'") ||
		strings.Contains(args, "F\"") || strings.Contains(args, "F'") {
		var b strings.Builder
		b.Grow(len(args))
		for _, m := range fStringInterpRE.FindAllStringSubmatchIndex(args, -1) {
			b.WriteByte(' ')
			b.WriteString(args[m[2]:m[3]])
			b.WriteByte(' ')
		}
		masked += b.String()
	}
	return sensitiveValueRE.MatchString(masked)
}

func argvSegmentLooksTrusted(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return true
	}
	if t == "sys.executable" {
		return true
	}
	if upperSnakeRE.MatchString(t) {
		return true
	}
	if strings.HasPrefix(t, "str(") && strings.HasSuffix(t, ")") {
		inner := strings.TrimSpace(t[4 : len(t)-1])
		if upperSnakeRE.MatchString(inner) || inner == "sys.executable" {
			return true
		}
	}
	if looksTrustedArgvConcat(t) {
		return true
	}
	return false
}

func looksTrustedArgvConcat(expr string) bool {
	t := strings.TrimSpace(expr)
	if !strings.Contains(t, "+") {
		return false
	}
	parts := splitTopLevelConcat(t)
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if isPureStringLiteral(p) || isPureBytesLiteral(p) || argvSegmentLooksTrusted(p) {
			continue
		}
		return false
	}
	return true
}

func sourceUsesRuamelYAML(src string) bool {
	return strings.Contains(src, "ruamel.yaml") ||
		strings.Contains(src, "from ruamel") ||
		strings.Contains(src, "import ruamel")
}

// yamlLoadLooksLikeRuamel is true for ruamel.YAML().load-style calls (no
// PyYAML Loader= kwarg). PyYAML yaml.load(..., Loader=...) stays reportable.
func yamlLoadLooksLikeRuamel(src, args string) bool {
	compact := compactWhitespace(args)
	if strings.Contains(compact, "Loader=") {
		return false
	}
	return sourceUsesRuamelYAML(src) || strings.Contains(src, "YAML()")
}

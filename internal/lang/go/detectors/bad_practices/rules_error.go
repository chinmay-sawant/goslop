package badpractices

import (
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-1", detectBP1)
	RegisterRule("BP-2", detectBP2)
	RegisterRule("BP-4", detectBP4)
	RegisterRule("BP-5", detectBP5)
	RegisterRule("BP-67", detectBP67)
	RegisterRule("BP-68", detectBP68)
	RegisterRule("BP-70", detectBP70)
	RegisterRule("BP-154", detectBP154)
}

// BP-1: discarded error-shaped return via `_`.
func detectBP1(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-1")
	if !facts.has("_") {
		return
	}
	msg := "discarded error return; handle or explicitly ignore with a comment"

	if len(facts.assignNodes) > 0 {
		for _, a := range facts.assignNodes {
			lhs, rhs, ok := splitAssign(a.text)
			if !ok || !strings.Contains(rhs, "(") || isNonErrorBuiltinRHS(rhs) {
				continue
			}
			if lhsDiscardsPossibleError(lhs) {
				pushAt(unit, meta, a.start, msg, out)
			}
		}
		return
	}

	// Text fallback when AST unavailable.
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		lhs, rhs, ok := splitAssign(t)
		if !ok || !strings.Contains(rhs, "(") || isNonErrorBuiltinRHS(rhs) {
			continue
		}
		if lhsDiscardsPossibleError(lhs) {
			pushAt(unit, meta, line.byte+strings.Index(line.text, t), msg, out)
		}
	}
}

func splitAssign(text string) (lhs, rhs string, ok bool) {
	if i := strings.Index(text, ":="); i >= 0 {
		return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+2:]), true
	}
	// Prefer single '=' not part of ==, !=, <=, >=
	for i := 0; i < len(text); i++ {
		if text[i] != '=' {
			continue
		}
		if i > 0 {
			prev := text[i-1]
			if prev == '!' || prev == '<' || prev == '>' || prev == '=' {
				continue
			}
		}
		if i+1 < len(text) && text[i+1] == '=' {
			continue
		}
		return strings.TrimSpace(text[:i]), strings.TrimSpace(text[i+1:]), true
	}
	return "", "", false
}

func lhsDiscardsPossibleError(lhs string) bool {
	names := strings.Split(lhs, ",")
	hasBlank := false
	bindsErr := false
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "_" {
			hasBlank = true
		}
		if n == "err" || n == "error" || strings.HasSuffix(n, "Err") {
			bindsErr = true
		}
	}
	return hasBlank && !bindsErr
}

func isNonErrorBuiltinRHS(rhs string) bool {
	t := strings.TrimSpace(rhs)
	name := t
	if i := strings.Index(t, "("); i >= 0 {
		name = strings.TrimSpace(t[:i])
	}
	if j := strings.LastIndex(name, "."); j >= 0 {
		name = name[j+1:]
	}
	switch name {
	case "len", "cap", "append", "make", "new", "copy", "delete", "clear",
		"min", "max", "real", "imag", "complex", "close", "panic", "recover",
		"print", "println":
		return true
	default:
		return false
	}
}

// BP-2: naked `return err` without wrapping.
func detectBP2(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-2")
	if !strings.Contains(unit.Source, "return err") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		if strings.TrimSpace(line.text) == "return err" {
			pushAt(unit, meta, line.byte, "naked return err loses operation context; wrap it before returning", out)
		}
	}
}

// BP-4: recover() without nearby logging.
func detectBP4(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-4")
	if !facts.has("recover()") && !strings.Contains(unit.Source, "recover()") {
		return
	}
	src := unit.Source
	reports := strings.Contains(src, "log.") ||
		strings.Contains(src, "Logger.") ||
		strings.Contains(src, ".Error(") ||
		strings.Contains(src, ".Warn(") ||
		strings.Contains(src, "fmt.Printf(") ||
		strings.Contains(src, "fmt.Fprintf(") ||
		strings.Contains(src, "slog.")
	if reports {
		return
	}
	if pos := strings.Index(src, "recover()"); pos >= 0 {
		pushAt(unit, meta, pos, "recover() suppresses panic information without logging or reporting it", out)
	}
}

// BP-5: Close() errors ignored through bare or deferred calls.
func detectBP5(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-5")
	if !strings.Contains(unit.Source, ".Close(") {
		return
	}
	msg := "Close() return value is ignored; check the close error where it can affect correctness"

	if len(facts.callNodes) > 0 {
		for _, c := range facts.callNodes {
			if !strings.HasSuffix(strings.TrimSpace(c.text), ".Close()") {
				continue
			}
			if closeCallIsHandled(c, facts) {
				continue
			}
			pushAt(unit, meta, c.start, msg, out)
		}
		return
	}

	// Text fallback: flag defer x.Close() and bare x.Close().
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, ".Close()") {
			if strings.HasPrefix(t, "defer ") || (!strings.Contains(t, "=") && !strings.HasPrefix(t, "return ") && !strings.Contains(t, "if ")) {
				// skip if assigned: err := x.Close() / if err := x.Close()
				if strings.Contains(t, ":=") || strings.Contains(t, "err") && strings.Contains(t, "=") {
					continue
				}
				if idx := strings.Index(line.text, ".Close()"); idx >= 0 {
					// approximate start of call
					pushAt(unit, meta, line.byte+idx, msg, out)
				}
			}
		}
	}
}

func closeCallIsHandled(c callSpan, facts *bpFacts) bool {
	// Parent kind from AST walk.
	switch c.parent {
	case "return_statement", "short_var_declaration":
		return true
	case "assignment_statement":
		// Check surrounding assignment text
		for _, a := range facts.assignNodes {
			if c.start >= a.start && c.end <= a.end {
				lhs, _, ok := splitAssign(a.text)
				return ok && strings.TrimSpace(lhs) != "_"
			}
		}
		return false
	case "defer_statement", "expression_statement":
		return false
	}
	// Heuristic: if the call is on a line with := or err = it's handled.
	lineStart := c.start
	for lineStart > 0 && facts.Source[lineStart-1] != '\n' {
		lineStart--
	}
	lineEnd := c.end
	for lineEnd < len(facts.Source) && facts.Source[lineEnd] != '\n' {
		lineEnd++
	}
	line := facts.Source[lineStart:lineEnd]
	if strings.Contains(line, ":=") {
		return true
	}
	if strings.Contains(line, "return ") {
		return true
	}
	if i := strings.Index(line, "="); i >= 0 {
		lhs := strings.TrimSpace(line[:i])
		return lhs != "_" && !strings.HasPrefix(strings.TrimSpace(line), "defer ")
	}
	return false
}

// BP-67: errors.As target not a pointer.
func detectBP67(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-67")
	if !strings.Contains(unit.Source, "errors.As(") {
		return
	}
	// Flag errors.As(err, target) where second arg lacks &
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "errors.As(")
		if idx < 0 {
			break
		}
		abs := start + idx
		// find args
		open := abs + len("errors.As(") - 1
		depth := 0
		end := -1
		for i := open; i < len(src); i++ {
			switch src[i] {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					end = i
				}
			}
			if end >= 0 {
				break
			}
		}
		if end < 0 {
			start = abs + 1
			continue
		}
		args := src[open+1 : end]
		parts := splitTopLevel(args)
		if len(parts) >= 2 {
			target := strings.TrimSpace(parts[1])
			if target != "" && !strings.HasPrefix(target, "&") && target != "nil" {
				// var target *T; errors.As(err, target) — still needs &target
				pushAt(unit, meta, abs, "errors.As target must be a pointer to the destination variable", out)
			}
		}
		start = end + 1
	}
}

func splitTopLevel(args string) []string {
	var out []string
	depth := 0
	start := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(args); i++ {
		c := args[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if (inStr == '"' || inStr == '\'') && c == '\\' {
				escape = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '`', '\'':
			inStr = c
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

// BP-68: discarded errors.Join result.
func detectBP68(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-68")
	if !strings.Contains(unit.Source, "errors.Join(") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "errors.Join(") || strings.HasPrefix(t, "_ = errors.Join(") {
			pushAt(unit, meta, line.byte, "errors.Join result is discarded; capture or return the joined error", out)
		}
	}
}

// BP-70: error branch logs with log.Print* and continues (Rust core_language_deferred).
// log.Fatal / panic / os.Exit / return are explicit exits and must not fire.
func detectBP70(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-70")
	lines := codeLines(unit.Source)
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i].text)
		// Only non-fatal log.Print family (Rust is_error_log_call).
		if !isBP70ErrorLogCall(t) {
			continue
		}
		// Must mention err-shaped identifier in the call args.
		if !strings.Contains(t, "err") && !strings.Contains(t, "Err") {
			continue
		}
		// Require error-ish message literal (Rust error_message_literal).
		if !bp70ErrorMessageLiteral(t) {
			// Still allow log.Print(err) / log.Println(err) with no string.
			if !(strings.Contains(t, "log.Print(err") || strings.Contains(t, "log.Println(err") ||
				strings.Contains(t, "log.Print( err") || strings.Contains(t, "log.Println( err")) {
				continue
			}
		}
		// Find enclosing if err != nil by scanning upward for condition.
		ifStart := -1
		for j := i; j >= 0 && j >= i-12; j-- {
			prev := strings.TrimSpace(lines[j].text)
			if strings.Contains(prev, "if ") && (strings.Contains(prev, "err != nil") ||
				strings.Contains(prev, "err!= nil") || strings.Contains(prev, "err !=nil") ||
				strings.Contains(prev, "err!=nil")) {
				ifStart = j
				break
			}
			if prev == "}" {
				break
			}
		}
		if ifStart < 0 {
			continue
		}
		// Explicit exit in the same branch (until closing brace of the if).
		if bp70BranchHasExplicitExit(lines, i) {
			continue
		}
		pushAt(unit, meta, lines[ifStart].byte,
			"error is logged and execution continues; return, panic, or otherwise handle the failure", out)
	}
}

func isBP70ErrorLogCall(line string) bool {
	// Fatal/Panic are exits, not "continue after log".
	if strings.Contains(line, "log.Fatal") || strings.Contains(line, "log.Panic") ||
		strings.Contains(line, ".Fatal(") || strings.Contains(line, ".Fatalf(") ||
		strings.Contains(line, ".Fatalln(") || strings.Contains(line, ".Panic(") {
		return false
	}
	return strings.Contains(line, "log.Print(") ||
		strings.Contains(line, "log.Printf(") ||
		strings.Contains(line, "log.Println(")
}

func bp70ErrorMessageLiteral(line string) bool {
	lower := strings.ToLower(line)
	for _, n := range []string{
		"error", "failed", "failure", "unable", "cannot", "could not", "invalid", "timeout",
	} {
		if strings.Contains(lower, n) {
			return true
		}
	}
	return false
}

func bp70BranchHasExplicitExit(lines []struct {
	idx  int
	text string
	byte int
}, logIdx int) bool {
	for j := logIdx; j < len(lines) && j < logIdx+12; j++ {
		n := strings.TrimSpace(lines[j].text)
		if n == "}" {
			return false
		}
		if strings.HasPrefix(n, "return") || strings.HasPrefix(n, "break") ||
			strings.HasPrefix(n, "continue") || strings.HasPrefix(n, "panic(") ||
			strings.Contains(n, "os.Exit(") || strings.Contains(n, "log.Fatal") ||
			strings.Contains(n, "log.Fatalf") || strings.Contains(n, "log.Fatalln") ||
			strings.Contains(n, "runtime.Goexit") {
			return true
		}
	}
	return false
}

// BP-154: discarded json.Unmarshal error (expression statement).
func detectBP154(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-154")
	if !strings.Contains(unit.Source, "json.Unmarshal(") {
		return
	}
	msg := "json.Unmarshal error is discarded; check the returned error"
	if len(facts.callNodes) > 0 {
		for _, c := range facts.callNodes {
			if c.callee != "json.Unmarshal" && !strings.HasSuffix(c.callee, ".Unmarshal") {
				if !strings.Contains(c.text, "json.Unmarshal(") {
					continue
				}
			}
			if c.parent == "expression_statement" || c.parent == "" {
				// confirm not assigned
				lineStart := c.start
				for lineStart > 0 && unit.Source[lineStart-1] != '\n' {
					lineStart--
				}
				lineEnd := c.end
				for lineEnd < len(unit.Source) && unit.Source[lineEnd] != '\n' {
					lineEnd++
				}
				line := strings.TrimSpace(unit.Source[lineStart:lineEnd])
				if strings.HasPrefix(line, "json.Unmarshal(") || strings.HasPrefix(line, "_ = json.Unmarshal(") {
					pushAt(unit, meta, c.start, msg, out)
				}
			}
		}
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "json.Unmarshal(") || strings.HasPrefix(t, "_ = json.Unmarshal(") {
			pushAt(unit, meta, line.byte, msg, out)
		}
	}
}

package perf

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
)

// pyPerfFacts is a single per-file source view shared by all PERF-PY rules.
type pyPerfFacts struct {
	Source string
	lines  []codeLine
}

type codeLine struct {
	text string
	raw  string
	byte int
}

func buildFacts(unit *core.ParsedUnit) *pyPerfFacts {
	if unit == nil {
		return &pyPerfFacts{}
	}
	return &pyPerfFacts{Source: unit.Source, lines: buildCodeLines(unit.Source)}
}

func buildCodeLines(source string) []codeLine {
	if source == "" {
		return nil
	}
	lines := strings.SplitAfter(source, "\n")
	out := make([]codeLine, 0, len(lines))
	offset := 0
	inTriple := byte(0)
	for _, rawWithNL := range lines {
		raw := strings.TrimSuffix(rawWithNL, "\n")
		text, nextTriple := stripPyLineForFacts(raw, inTriple)
		inTriple = nextTriple
		out = append(out, codeLine{text: text, raw: raw, byte: offset})
		offset += len(rawWithNL)
	}
	return out
}

// stripPyLineForFacts blanks comments and triple-quoted string bodies while
// keeping ordinary single-line string literals (so keyword heuristics like
// status="pending" still work).
func stripPyLineForFacts(line string, inTriple byte) (string, byte) {
	if inTriple != 0 {
		if at := strings.Index(line, string([]byte{inTriple, inTriple, inTriple})); at >= 0 {
			// Keep prefix blanked; resume after closing delimiter.
			rest, next := stripPyLineForFacts(line[at+3:], 0)
			return strings.Repeat(" ", at+3) + rest, next
		}
		return strings.Repeat(" ", len(line)), inTriple
	}

	buf := []byte(line)
	inQuote := byte(0)
	escaped := false
	for i := 0; i < len(buf); i++ {
		c := buf[i]
		if inQuote != 0 {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == inQuote {
				inQuote = 0
			}
			continue
		}
		if c == '#' {
			return strings.TrimRight(string(buf[:i]), " \t"), 0
		}
		if c == '\'' || c == '"' {
			if i+2 < len(buf) && buf[i+1] == c && buf[i+2] == c {
				// Enter triple-quoted string; blank the opener and the rest of the line.
				for j := i; j < len(buf); j++ {
					buf[j] = ' '
				}
				return string(buf), c
			}
			inQuote = c
		}
	}
	return string(buf), 0
}

func stripPyComment(line string) string {
	text, _ := stripPyLineForFacts(line, 0)
	return text
}

func indentWidth(line string) int { return len(line) - len(strings.TrimLeft(line, " \t")) }

func isLoopLine(line string) bool {
	t := strings.TrimSpace(line)
	return strings.HasPrefix(t, "for ") || strings.HasPrefix(t, "while ")
}

func inLoop(lines []codeLine, idx int) bool {
	if idx < 0 || idx >= len(lines) {
		return false
	}
	indent := indentWidth(lines[idx].raw)
	for i := idx - 1; i >= 0; i-- {
		t := strings.TrimSpace(lines[i].text)
		if t == "" {
			continue
		}
		if indentWidth(lines[i].raw) < indent && isLoopLine(t) {
			return true
		}
	}
	return false
}

func functionWindow(lines []codeLine, at int) (int, int) {
	start, end := 0, len(lines)
	for i := at; i >= 0; i-- {
		t := strings.TrimSpace(lines[i].text)
		if strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") {
			start = i
			base := indentWidth(lines[i].raw)
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j].text) != "" && indentWidth(lines[j].raw) <= base {
					end = j
					break
				}
			}
			break
		}
	}
	return start, end
}

func windowHas(lines []codeLine, start, end int, needles ...string) bool {
	start, end = safeLineRange(lines, start, end)
	for _, line := range lines[start:end] {
		for _, needle := range needles {
			if strings.Contains(line.text, needle) {
				return true
			}
		}
	}
	return false
}

// safeLineRange keeps heuristic window calculations total. A detector may
// infer an end before the current line when a top-level loop follows a
// function; that should mean an empty window, never a process panic.
func safeLineRange(lines []codeLine, start, end int) (int, int) {
	if start < 0 {
		start = 0
	}
	if start > len(lines) {
		start = len(lines)
	}
	if end < 0 {
		end = 0
	}
	if end > len(lines) {
		end = len(lines)
	}
	if end < start {
		start = end
	}
	return start, end
}

package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/ast"
	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/pytext"
)

// PyCweFacts is the per-file fact bag for Python CWE source-pattern detectors.
// Priority batch uses SourceIndex needle presence; call heuristics scan Source
// against a once-built Masked view (comments/strings blanked, offsets stable).
type PyCweFacts struct {
	Index  ast.SourceIndex
	Source string
	Masked string
	Funcs  []pythonFunction

	lines      []pyMaskedLine
	linesReady bool
}

var pyFunctionDefRE = regexp.MustCompile(`(?m)^([ \t]*)def\s+([A-Za-z_][A-Za-z0-9_]*)\s*\([^\n]*\)\s*:`)

// pythonFunction is a def span in Source (body is Source[bodyStart:bodyEnd]).
type pythonFunction struct {
	name      string
	start     int
	bodyStart int
	body      string
}

// BuildFacts constructs PyCweFacts for a unit.
//
// Masks the file once and extracts function spans once. Rules reuse Masked /
// Funcs instead of calling pytext.Mask per detector.
func BuildFacts(unit *core.ParsedUnit) *PyCweFacts {
	facts := &PyCweFacts{}
	if unit == nil || unit.Source == "" {
		facts.Index = ast.Build("", pyCweNeedles)
		return facts
	}
	facts.Source = unit.Source
	facts.Masked = pytext.Mask(unit.Source)
	facts.Index = ast.Build(unit.Source, pyCweNeedles)
	facts.Funcs = buildPythonFunctions(facts.Source, facts.Masked)
	return facts
}

// Functions returns memoized def spans for this file.
func (f *PyCweFacts) Functions() []pythonFunction {
	if f == nil {
		return nil
	}
	return f.Funcs
}

// codeMask returns a masked view of fragment aligned with Source when possible.
// fragStart is the byte offset of fragment in Source (0 for the full file).
// Pass fragStart < 0 only when the offset is unknown; that path remasks.
func (f *PyCweFacts) codeMask(fragment string, fragStart int) string {
	if f == nil || f.Masked == "" {
		return pythonCodeMask(fragment)
	}
	if fragment == f.Source || (fragStart == 0 && len(fragment) == len(f.Source)) {
		return f.Masked
	}
	if fragStart >= 0 && fragStart+len(fragment) <= len(f.Masked) {
		return f.Masked[fragStart : fragStart+len(fragment)]
	}
	return pythonCodeMask(fragment)
}

// MaskedLines returns masked line spans, built once lazily from Masked.
func (f *PyCweFacts) MaskedLines() []pyMaskedLine {
	if f == nil {
		return nil
	}
	if f.linesReady {
		return f.lines
	}
	f.lines = buildMaskedPythonLines(f.Masked)
	f.linesReady = true
	return f.lines
}

func buildPythonFunctions(source, masked string) []pythonFunction {
	matches := pyFunctionDefRE.FindAllStringSubmatchIndex(masked, -1)
	out := make([]pythonFunction, 0, len(matches))
	for _, match := range matches {
		indent := len(masked[match[2]:match[3]])
		bodyStart := match[1]
		bodyEnd := pythonFunctionBodyEnd(masked, bodyStart, indent)
		out = append(out, pythonFunction{
			name:      source[match[4]:match[5]],
			start:     match[0],
			bodyStart: bodyStart,
			body:      source[bodyStart:bodyEnd],
		})
	}
	return out
}

func pythonFunctionBodyEnd(code string, bodyStart, indent int) int {
	for lineStart := bodyStart; lineStart < len(code); {
		if code[lineStart] == '\n' {
			lineStart++
		}
		lineEnd := len(code)
		if next := strings.IndexByte(code[lineStart:], '\n'); next >= 0 {
			lineEnd = lineStart + next
		}
		line := code[lineStart:lineEnd]
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lineIndent := len(line) - len(strings.TrimLeft(line, " \t"))
			if lineIndent <= indent {
				return lineStart
			}
		}
		if lineEnd == len(code) {
			break
		}
		lineStart = lineEnd + 1
	}
	return len(code)
}

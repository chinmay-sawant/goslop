package badpractices

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
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

func isTestFile(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	p := fileDisplayPath(unit)
	if p == "" {
		p = unit.Path
	}
	return strings.HasSuffix(p, "_test.go")
}

func isMaterializedFixture(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, p := range []string{unit.DisplayPath, unit.Path} {
		if strings.Contains(p, "target/codehound-fixtures/") ||
			strings.Contains(p, `target\codehound-fixtures\`) ||
			strings.Contains(p, "codehound-fixtures") {
			return true
		}
	}
	return false
}

func pushAt(unit *core.ParsedUnit, meta *rules.RuleMetadata, byteOffset int, message string, out *[]rules.Finding) {
	if unit == nil || meta == nil || out == nil {
		return
	}
	line, col := unit.LineCol(byteOffset)
	rules.PushFinding(meta, fileDisplayPath(unit), line, col, message, out)
}

func lineStartByte(source string, lineIdx int) int {
	byteOff := 0
	for i, line := range strings.Split(source, "\n") {
		if i == lineIdx {
			return byteOff
		}
		byteOff += len(line) + 1
	}
	return byteOff
}

func packageName(source string) string {
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "package ") {
			return strings.TrimSpace(strings.TrimPrefix(t, "package "))
		}
	}
	return ""
}

func enclosingFuncName(source string, startByte int) (string, bool) {
	if startByte > len(source) {
		startByte = len(source)
	}
	if startByte < 0 {
		startByte = 0
	}
	head := source[:startByte]
	funcKw := strings.LastIndex(head, "func ")
	if funcKw < 0 {
		return "", false
	}
	after := strings.TrimLeft(source[funcKw+len("func "):startByte], " \t")
	if strings.HasPrefix(after, "(") {
		close := strings.Index(after, ")")
		if close < 0 {
			return "", false
		}
		after = strings.TrimLeft(after[close+1:], " \t")
	}
	end := 0
	for end < len(after) {
		r := rune(after[end])
		if end == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				break
			}
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			break
		}
		end++
	}
	name := after[:end]
	if name == "" {
		return "", false
	}
	return name, true
}

func stripLineComment(line string) string {
	inStr := byte(0)
	escape := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if inStr == '"' || inStr == '\'' {
				if c == '\\' {
					escape = true
					continue
				}
				if c == inStr {
					inStr = 0
				}
				continue
			}
			if inStr == '`' && c == '`' {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '`', '\'':
			inStr = c
		case '/':
			if i+1 < len(line) && line[i+1] == '/' {
				return strings.TrimRight(line[:i], " \t")
			}
		}
	}
	return line
}

func codeLines(source string) []struct {
	idx  int
	text string
	byte int
} {
	var out []struct {
		idx  int
		text string
		byte int
	}
	byteOff := 0
	for i, line := range strings.Split(source, "\n") {
		out = append(out, struct {
			idx  int
			text string
			byte int
		}{i, stripLineComment(line), byteOff})
		byteOff += len(line) + 1
	}
	return out
}

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
			if isIdentByte(prev) || prev == '.' {
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

func isIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func pathBase(p string) string {
	return filepath.Base(p)
}

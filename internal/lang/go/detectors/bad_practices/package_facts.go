package badpractices

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/chinmay/goslop/internal/core"
)

// packageTypeFacts holds same-package method sets and exported interfaces
// (used by BP-30 / BP-31 for cross-file package awareness).
type packageTypeFacts struct {
	// methods maps type name (no pointer) → exported method names.
	methods map[string]map[string]struct{}
	// interfaces maps exported interface name → method names.
	interfaces map[string][]string
}

var reExportedIface = regexp.MustCompile(`(?m)^type\s+([A-Z]\w*)\s+interface\s*\{`)

func packageTypeFactsForUnit(unit *core.ParsedUnit) *packageTypeFacts {
	facts := &packageTypeFacts{
		methods:    map[string]map[string]struct{}{},
		interfaces: map[string][]string{},
	}
	if unit == nil {
		return facts
	}
	dir := filepath.Dir(unit.Path)
	if dir == "" || dir == "." {
		dir = filepath.Dir(fileDisplayPath(unit))
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		// Fall back to the current unit only.
		collectTypeFactsFromSource(unit.Source, facts)
		return facts
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		p := filepath.Join(dir, name)
		data, rerr := os.ReadFile(p)
		if rerr != nil {
			continue
		}
		collectTypeFactsFromSource(string(data), facts)
	}
	return facts
}

func collectTypeFactsFromSource(src string, facts *packageTypeFacts) {
	if facts == nil || src == "" {
		return
	}
	// Method sets: func (r *Type) Method(
	for _, line := range strings.Split(src, "\n") {
		t := strings.TrimSpace(line)
		if !strings.HasPrefix(t, "func (") {
			continue
		}
		rest := strings.TrimPrefix(t, "func (")
		closeParen := strings.Index(rest, ")")
		if closeParen < 0 {
			continue
		}
		recv := strings.TrimSpace(rest[:closeParen])
		fields := strings.Fields(recv)
		if len(fields) == 0 {
			continue
		}
		typ := strings.TrimPrefix(fields[len(fields)-1], "*")
		after := strings.TrimSpace(rest[closeParen+1:])
		meth := firstIdent(after)
		if meth == "" || !unicode.IsUpper(rune(meth[0])) {
			continue
		}
		if facts.methods[typ] == nil {
			facts.methods[typ] = map[string]struct{}{}
		}
		facts.methods[typ][meth] = struct{}{}
	}

	// Exported interfaces.
	for _, m := range reExportedIface.FindAllStringSubmatchIndex(src, -1) {
		name := src[m[2]:m[3]]
		open := strings.Index(src[m[0]:], "{")
		if open < 0 {
			continue
		}
		absOpen := m[0] + open
		depth := 0
		end := absOpen
		for i := absOpen; i < len(src); i++ {
			switch src[i] {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					end = i
					i = len(src)
				}
			}
		}
		if end <= absOpen {
			continue
		}
		block := src[absOpen+1 : end]
		var methods []string
		for _, line := range strings.Split(block, "\n") {
			t := strings.TrimSpace(line)
			if t == "" || strings.HasPrefix(t, "//") {
				continue
			}
			if i := strings.Index(t, "("); i > 0 {
				meth := strings.TrimSpace(t[:i])
				// strip type params / embedded interfaces roughly
				if meth != "" && unicode.IsUpper(rune(meth[0])) && !strings.Contains(meth, ".") {
					methods = append(methods, meth)
				}
			}
		}
		if len(methods) > 0 {
			facts.interfaces[name] = methods
		}
	}
}

// typeImplements reports whether typeName's method set covers all methods.
func typeImplements(facts *packageTypeFacts, typeName string, methods []string) bool {
	if facts == nil || len(methods) == 0 {
		return false
	}
	impl := facts.methods[typeName]
	if impl == nil {
		return false
	}
	for _, m := range methods {
		if _, ok := impl[m]; !ok {
			return false
		}
	}
	return true
}

// countWordOccurrences counts non-overlapping word-boundary matches of word in source.
func countWordOccurrences(source, word string) int {
	if word == "" || source == "" {
		return 0
	}
	n := 0
	start := 0
	for {
		idx := strings.Index(source[start:], word)
		if idx < 0 {
			return n
		}
		abs := start + idx
		beforeOK := abs == 0 || !isIdentByteBP(source[abs-1])
		after := abs + len(word)
		afterOK := after >= len(source) || !isIdentByteBP(source[after])
		if beforeOK && afterOK {
			n++
		}
		start = abs + len(word)
	}
}

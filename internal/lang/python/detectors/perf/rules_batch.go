package perf

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("PERF-PY-3", detectPYPERF3)
	RegisterRule("PERF-PY-4", detectPYPERF4)
	RegisterRule("PERF-PY-10", detectPYPERF10)
	RegisterRule("PERF-PY-11", detectPYPERF11)
}

var (
	createCallRE    = regexp.MustCompile(`\b[A-Za-z_][A-Za-z0-9_]*\.objects\.create\s*\(`)
	counterAssignRE = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\.([A-Za-z_][A-Za-z0-9_]*)\s*(?:\+=|-=)`)
	fieldAssignRE   = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\.([A-Za-z_][A-Za-z0-9_]*)\s*=`)
)

func detectPYPERF3(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !inLoop(facts.lines, i) || !createCallRE.MatchString(line.text) || strings.Contains(strings.ToLower(line.raw), "signal") || strings.Contains(strings.ToLower(line.raw), "hook") {
			continue
		}
		if i > 0 && (strings.Contains(strings.ToLower(facts.lines[i-1].raw), "signal") || strings.Contains(strings.ToLower(facts.lines[i-1].raw), "hook")) {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if strings.Contains(facts.lines[start].text, "range(2)") || strings.Contains(facts.lines[start].text, "range(1)") {
			continue
		}
		if createdValueIsImmediatelyRequired(facts.lines, i, end) {
			continue
		}
		if isDependentCreate(facts.lines, i, start) {
			continue
		}
		pushLine(unit, "PERF-PY-3", line, ".objects.create", "Django create runs once per batch item; use bulk_create when hooks and generated IDs are not required", out)
	}
}

// isDependentCreate preserves the common parent/child transaction shape: the
// child create consumes an object made in the same iteration and cannot be
// flattened into one independent bulk_create without changing identity flow.
func isDependentCreate(lines []codeLine, at, start int) bool {
	for i := at - 1; i >= start && i >= at-3; i-- {
		m := regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*.*\.objects\.create\s*\(`).FindStringSubmatch(lines[i].text)
		if len(m) > 1 && strings.Contains(lines[at].text, "="+m[1]) {
			return true
		}
	}
	return false
}

func createdValueIsImmediatelyRequired(lines []codeLine, at, end int) bool {
	assign := regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*.*\.objects\.create\s*\(`).FindStringSubmatch(lines[at].text)
	if len(assign) < 2 {
		return false
	}
	name := assign[1]
	for i := at + 1; i < end && i <= at+3; i++ {
		t := lines[i].text
		if strings.Contains(t, name+".") || (strings.Contains(t, ".objects.create(") && strings.Contains(t, "="+name)) {
			return true
		}
	}
	return false
}

func detectPYPERF4(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		m := counterAssignRE.FindStringSubmatch(line.text)
		if len(m) < 3 || !numericField(m[2]) {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "models.F(", ".update(") || windowHasRaw(facts.lines, start, end, "save hook", "save-hook") {
			continue
		}
		laterStart, laterEnd := safeLineRange(facts.lines, i+1, end)
		for _, later := range facts.lines[laterStart:laterEnd] {
			if strings.Contains(later.text, m[1]+".save(") {
				pushLine(unit, "PERF-PY-4", line, m[1]+".", "numeric model field is mutated then saved; use QuerySet.update with an F expression when hooks are unnecessary", out)
				break
			}
		}
	}
}

func numericField(field string) bool {
	field = strings.ToLower(field)
	for _, needle := range []string{"quantity", "stock", "count", "balance", "available", "reserved", "total", "amount"} {
		if strings.Contains(field, needle) {
			return true
		}
	}
	return false
}

func detectPYPERF10(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !isLoopLine(line.text) {
			continue
		}
		loopIndent := indentWidth(line.raw)
		end := len(facts.lines)
		for j := i + 1; j < len(facts.lines); j++ {
			if strings.TrimSpace(facts.lines[j].text) != "" && indentWidth(facts.lines[j].raw) <= loopIndent {
				end = j
				break
			}
		}
		processed := ""
		for j := i + 1; j < end; j++ {
			if m := regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(?:await\s+)?(?:process|handle|run|claim)[A-Za-z_0-9]*\s*\(`).FindStringSubmatch(facts.lines[j].text); len(m) > 1 {
				processed = m[1]
			}
			if processed == "" || !strings.Contains(facts.lines[j].text, "if "+processed) {
				continue
			}
			if successThenSleeps(facts.lines, j, end) {
				pushLine(unit, "PERF-PY-10", facts.lines[j], processed, "worker sleeps even after successful work; continue immediately and sleep only when idle", out)
				break
			}
		}
	}
}

func successThenSleeps(lines []codeLine, ifAt, end int) bool {
	ifIndent := indentWidth(lines[ifAt].raw)
	hasExit := false
	for i := ifAt + 1; i < end; i++ {
		indent := indentWidth(lines[i].raw)
		if strings.TrimSpace(lines[i].text) == "" {
			continue
		}
		if indent > ifIndent && (strings.TrimSpace(lines[i].text) == "continue" || strings.HasPrefix(strings.TrimSpace(lines[i].text), "return")) {
			hasExit = true
		}
		if indent <= ifIndent && (strings.Contains(lines[i].text, "time.sleep(") || strings.Contains(lines[i].text, "asyncio.sleep(")) {
			return !hasExit
		}
	}
	return false
}

func detectPYPERF11(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		isSave := strings.Contains(line.text, ".save(")
		isCommit := strings.Contains(line.text, ".commit(")
		if (!isSave && !isCommit) || (isSave && !inLoop(facts.lines, i)) {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, ".update(") || windowHasRaw(facts.lines, start, end, "save hook", "external effect", "derived value") {
			continue
		}
		for j := i - 1; j >= start; j-- {
			m := fieldAssignRE.FindStringSubmatch(facts.lines[j].text)
			if len(m) < 3 || (isSave && !strings.Contains(line.text, m[1]+".save(")) {
				continue
			}
			loop, ok := enclosingLoopHeader(facts.lines, j)
			if !ok || (!strings.Contains(loop.text, ".objects.") && !strings.Contains(loop.text, ".query(") && !strings.Contains(loop.text, ".scalars().all")) {
				continue
			}
			if strings.Contains(facts.lines[j].text, m[1]+".") && !strings.Contains(facts.lines[j].text, m[1]+"."+m[2]+" ==") {
				pushLine(unit, "PERF-PY-11", line, ".save", "ORM rows are mutated and saved one at a time; use a set-based update when every row receives the same value", out)
				break
			}
		}
	}
}

func enclosingLoopHeader(lines []codeLine, at int) (codeLine, bool) {
	if at < 0 || at >= len(lines) {
		return codeLine{}, false
	}
	indent := indentWidth(lines[at].raw)
	for i := at - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i].text) == "" {
			continue
		}
		if indentWidth(lines[i].raw) < indent && isLoopLine(lines[i].text) {
			return lines[i], true
		}
	}
	return codeLine{}, false
}

func windowHasRaw(lines []codeLine, start, end int, needles ...string) bool {
	start, end = safeLineRange(lines, start, end)
	for _, line := range lines[start:end] {
		lower := strings.ToLower(line.raw)
		for _, needle := range needles {
			if strings.Contains(lower, needle) {
				return true
			}
		}
	}
	return false
}

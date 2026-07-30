package badpractices

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	metaByID["BP-PY-38"] = &rules.RuleMetadata{
		ID: "BP-PY-38", Title: "asyncio create_task Without Reference",
		Description: "`asyncio.create_task` result is discarded, risking silent task loss and missed errors.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Store the task and await/gather it with exception handling.",
	}
	metaByID["BP-PY-39"] = &rules.RuleMetadata{
		ID: "BP-PY-39", Title: "time.sleep In Async Function",
		Description: "`time.sleep` blocks the event loop inside `async def`.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Use `await asyncio.sleep(...)` instead of `time.sleep`.",
	}
	metaByID["BP-PY-40"] = &rules.RuleMetadata{
		ID: "BP-PY-40", Title: "threading Without Join Or Shutdown",
		Description: "A Thread is started without join or a clear shutdown protocol.",
		Severity:    rules.SeverityLow, Pack: rules.PackBadPractice,
		Fix: "Call `.join()` (or use a pool with shutdown); avoid fire-and-forget non-daemon threads.",
	}
	RegisterRule("BP-PY-38", detectBPPY38)
	RegisterRule("BP-PY-39", detectBPPY39)
	RegisterRule("BP-PY-40", detectBPPY40)
}

// BP-PY-38: bare create_task / ensure_future as a statement (return unused).
func detectBPPY38(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-38")
	if !facts.hasAny("create_task", "ensure_future") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" || !isBareTaskSpawn(t) {
			continue
		}
		// Consumed forms on the same statement.
		if strings.Contains(t, ".append(") || strings.Contains(t, "gather(") ||
			strings.Contains(t, "wait(") || strings.Contains(t, "as_completed(") {
			continue
		}
		off := line.byte
		if i := strings.Index(line.text, "create_task"); i >= 0 {
			off = line.byte + i
		} else if i := strings.Index(line.text, "ensure_future"); i >= 0 {
			off = line.byte + i
		}
		pushAt(unit, meta, off,
			"create_task/ensure_future result discarded; store the task and await/gather with exception handling",
			out)
	}
}

// isBareTaskSpawn reports a statement that is only a create_task/ensure_future call
// (optionally asyncio.-qualified), not assigned or otherwise consumed.
func isBareTaskSpawn(t string) bool {
	t = strings.TrimSpace(t)
	t = strings.TrimSuffix(t, ";")
	t = strings.TrimSpace(t)
	if t == "" {
		return false
	}
	// await create_task(...) — task is driven by await; miss.
	if strings.HasPrefix(t, "await ") {
		return false
	}
	// Assignment / augmented assignment at statement level.
	if hasTopLevelAssign(t) {
		return false
	}
	prefixes := []string{
		"asyncio.create_task(",
		"create_task(",
		"asyncio.ensure_future(",
		"ensure_future(",
		"loop.create_task(",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(t, p) {
			return true
		}
	}
	return false
}

func hasTopLevelAssign(t string) bool {
	depth := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(t); i++ {
		c := t[i]
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
		case '=':
			if depth != 0 {
				continue
			}
			// Skip ==, !=, <=, >=, :=
			if i > 0 {
				prev := t[i-1]
				if prev == '=' || prev == '!' || prev == '<' || prev == '>' || prev == ':' {
					continue
				}
			}
			if i+1 < len(t) && t[i+1] == '=' {
				continue
			}
			return true
		}
	}
	return false
}

// BP-PY-39: time.sleep inside async def body (blocks event loop).
func detectBPPY39(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-39")
	if !strings.Contains(unit.Source, "time.sleep") {
		return
	}
	if !strings.Contains(unit.Source, "async def ") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	// Stack of function frames only (def / async def).
	type frame struct {
		indent  int
		isAsync bool
	}
	var stack []frame
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		ind := indentWidth(line.raw)
		for len(stack) > 0 && ind <= stack[len(stack)-1].indent {
			stack = stack[:len(stack)-1]
		}
		if strings.HasPrefix(t, "async def ") {
			stack = append(stack, frame{indent: ind, isAsync: true})
			continue
		}
		if strings.HasPrefix(t, "def ") {
			stack = append(stack, frame{indent: ind, isAsync: false})
			continue
		}
		if !strings.Contains(t, "time.sleep") {
			continue
		}
		if len(stack) == 0 || !stack[len(stack)-1].isAsync {
			continue
		}
		// Prefer call form.
		if !strings.Contains(t, "time.sleep(") && !strings.Contains(t, "time.sleep (") {
			continue
		}
		off := line.byte
		if i := strings.Index(line.text, "time.sleep"); i >= 0 {
			off = line.byte + i
		}
		pushAt(unit, meta, off,
			"time.sleep blocks the event loop inside async def; use await asyncio.sleep",
			out)
	}
}

// BP-PY-40: Thread started without join in same unit (review-only file-level heuristic).
func detectBPPY40(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-40")
	src := unit.Source
	if !strings.Contains(src, "threading") && !facts.has("threading.") {
		return
	}
	if !strings.Contains(src, "Thread") || !strings.Contains(src, ".start(") {
		return
	}
	// File-level: any .join( → miss (v0 heuristic).
	if strings.Contains(src, ".join(") {
		return
	}
	// Daemon-only policy: if every Thread construction uses daemon=True, miss.
	// Simpler: skip lines that include daemon=True on the same start/construct line.
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		if strings.Contains(t, "daemon=True") || strings.Contains(t, "daemon = True") {
			continue
		}
		if !strings.Contains(t, ".start(") {
			continue
		}
		// Prefer Thread(...).start() same line, or .start after threading.Thread in file.
		if strings.Contains(t, "Thread") || looksThreadStart(t, src) {
			off := line.byte
			if i := strings.Index(line.text, ".start("); i >= 0 {
				off = line.byte + i
			}
			pushAt(unit, meta, off,
				"threading.Thread started without .join in this file (review-only heuristic); join or use a pool with shutdown",
				out)
			return
		}
	}
}

func looksThreadStart(line, src string) bool {
	if !strings.Contains(src, "threading.Thread") && !strings.Contains(src, "Thread(") {
		return false
	}
	lower := strings.ToLower(line)
	for _, bad := range []string{"process.start", "server.start", "app.start", "loop.start", "timer.start"} {
		if strings.Contains(lower, bad) {
			return false
		}
	}
	return strings.Contains(line, ".start(")
}

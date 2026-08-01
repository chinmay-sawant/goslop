package perf

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("PERF-PY-2", detectPERFPY2)
	RegisterRule("PERF-PY-6", detectPERFPY6)
	RegisterRule("PERF-PY-8", detectPERFPY8)
	RegisterRule("PERF-PY-13", detectPERFPY13)
	RegisterRule("PERF-PY-19", detectPERFPY19)
	RegisterRule("PERF-PY-21", detectPERFPY21)
}

// PERF-PY-2 finds terminal Django ORM lookups that run for each item.
func detectPERFPY2(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for i, line := range facts.lines {
		if !inLoop(facts.lines, i) || !strings.Contains(line.text, ".objects.") {
			continue
		}
		loop, ok := enclosingLoopHeader(facts.lines, i)
		if !ok || perfSmallExplicitLoop(loop.text) {
			continue
		}
		loopVar, _, bindOK := perfLoopBinding(loop.text)
		if !bindOK || !strings.Contains(line.text, loopVar) {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if containsAnyFold(linesText(facts.lines, start, end), "page_size", "paginator", "page_number", "management command", "pagination") {
			continue
		}
		if strings.Contains(line.text, ".get(") || strings.Contains(line.text, ".first()") || strings.Contains(line.text, ".all()") {
			pushLine(unit, "PERF-PY-2", line, ".objects.", "Django ORM lookup executes inside an item loop; hoist or batch the lookup", out)
		}
	}
}

func perfSmallExplicitLoop(loop string) bool {
	if strings.Contains(loop, "range(1)") || strings.Contains(loop, "range(2)") {
		return true
	}
	if !strings.Contains(loop, " in (") && !strings.Contains(loop, " in [") {
		return false
	}
	inside := loop[strings.Index(loop, " in ")+4:]
	inside = strings.TrimRight(strings.TrimSpace(inside), ":")
	return strings.Count(inside, ",") < 2
}

var perfClaimAssignRE = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*[^\n]*\.objects\.filter\([^\n]*\)\s*\.first\s*\(`)

// PERF-PY-6 finds a visible select-then-mark work claim without a locking primitive.
// Status mutation must target the same bound name as the claim lookup.
func detectPERFPY6(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for i, line := range facts.lines {
		if !strings.Contains(line.text, ".objects.filter(") || !strings.Contains(line.text, ".first()") || !containsAnyFold(line.text, "pending", "queued", "ready", "new") {
			continue
		}
		m := perfClaimAssignRE.FindStringSubmatch(line.text)
		if len(m) < 2 {
			continue
		}
		name := m[1]
		start, end := functionWindow(facts.lines, i)
		body := linesText(facts.lines, start, end)
		if containsAnyFold(body, "select_for_update", "with_for_update", "skip_locked", "atomic update", "single-process", "single process") {
			continue
		}
		statusRE := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\s*\.\s*status\s*=`)
		saveRE := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\s*\.\s*save\s*\(`)
		if !statusRE.MatchString(body) || !saveRE.MatchString(body) {
			continue
		}
		pushLine(unit, "PERF-PY-6", line, ".objects.filter(", "work claim selects and updates status without a visible row lock or atomic update", out)
	}
}

// PERF-PY-8 is deliberately narrow: it only considers familiar relation-shaped names.
func detectPERFPY8(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for i, line := range facts.lines {
		loopVar, iterable, ok := perfLoopBinding(line.text)
		if !ok || !perfSQLAlchemyIterable(facts.lines, i, iterable) {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		body := linesText(facts.lines, start, end)
		if containsAnyFold(body, "joinedload", "selectinload", "contains_eager", "subqueryload") {
			continue
		}
		if perfRelationAccess(linesText(facts.lines, i+1, end), loopVar) {
			pushLine(unit, "PERF-PY-8", line, "for ", "possible SQLAlchemy lazy relationship access in a batch loop (review-level)", out)
		}
	}
}

// PERF-PY-13 detects a fully hydrated queryset where the loop body uses one scalar field.
func detectPERFPY13(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for i, line := range facts.lines {
		loopVar, iterable, ok := perfLoopBinding(line.text)
		if !ok || !strings.Contains(iterable, ".objects.") || containsAnyFold(iterable, ".values(", ".values_list(", ".only(", ".defer(") {
			continue
		}
		_, end := functionWindow(facts.lines, i)
		body := linesText(facts.lines, i+1, end)
		if perfUsesOnlyScalar(body, loopVar) {
			pushLine(unit, "PERF-PY-13", line, "for ", "ORM rows are fully hydrated while the loop uses only scalar fields; use a projection", out)
		}
	}
}

// PERF-PY-19 detects one transaction holding locks across an unbounded iteration.
func detectPERFPY19(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for i, line := range facts.lines {
		if !isLoopLine(line.text) || (!strings.Contains(line.text, "select_for_update") && !strings.Contains(line.text, "with_for_update")) {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		decoratorStart := start
		for decoratorStart > 0 && strings.HasPrefix(strings.TrimSpace(facts.lines[decoratorStart-1].text), "@") {
			decoratorStart--
		}
		body := linesText(facts.lines, decoratorStart, end)
		if !containsAnyFold(body, "transaction.atomic", "session.begin", "with session.begin") || containsAnyFold(line.text, "[:", ".limit(", "batch", "chunk", "keyset") || containsAnyFold(body, "per-batch commit", "per batch commit") {
			continue
		}
		needle := "select_for_update"
		if strings.Contains(line.text, "with_for_update") {
			needle = "with_for_update"
		}
		pushLine(unit, "PERF-PY-19", line, needle, "locking query is iterated without a visible batch bound inside a transaction", out)
	}
}

// PERF-PY-21 detects maintenance deletes that do not visibly bound the candidate set.
func detectPERFPY21(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for i, line := range facts.lines {
		if !strings.Contains(line.text, ".delete()") || !strings.Contains(line.text, ".filter(") {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		body := linesText(facts.lines, start, end)
		header := ""
		if start >= 0 && start < len(facts.lines) {
			header = facts.lines[start].text
		}
		// Require a maintenance marker in the function name or nearby comments/body —
		// do not treat the ".delete()" trigger itself as evidence.
		if !containsAnyFold(header+"\n"+body, "purge", "cleanup", "retention", "expire", "maintenance") ||
			containsAnyFold(body, "[:", ".limit(", "chunk", "batch", "keyset", "partition") {
			continue
		}
		pushLine(unit, "PERF-PY-21", line, ".delete()", "maintenance delete has no visible batch bound; delete in bounded chunks", out)
	}
}

var perfForBinding = regexp.MustCompile(`^\s*for\s+([A-Za-z_][A-Za-z0-9_]*)\s+in\s+(.+):\s*$`)

func perfLoopBinding(line string) (string, string, bool) {
	matches := perfForBinding.FindStringSubmatch(line)
	if len(matches) != 3 {
		return "", "", false
	}
	return matches[1], matches[2], true
}

func perfRelationAccess(body, variable string) bool {
	for _, name := range []string{"partner_endpoint", "inbound_event", "author", "customer", "profile", "relation", "parent", "owner", "account", "organization"} {
		if strings.Contains(body, variable+"."+name+".") || strings.Contains(body, variable+"."+name+")") || strings.Contains(body, variable+"."+name+"]") {
			return true
		}
	}
	return false
}

func perfSQLAlchemyIterable(lines []codeLine, loopAt int, iterable string) bool {
	if strings.Contains(iterable, ".all()") && containsAnyFold(iterable, ".query(", "session.execute(", "select(") {
		return true
	}
	name := strings.TrimSpace(iterable)
	if name == "" || strings.ContainsAny(name, ".([{") {
		return false
	}
	for i := loopAt - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i].text)
		if indentWidth(lines[i].raw) < indentWidth(lines[loopAt].raw) {
			break
		}
		if !strings.HasPrefix(line, name+" =") && !strings.HasPrefix(line, name+"=") {
			continue
		}
		return strings.Contains(line, ".all()") && containsAnyFold(line, ".query(", "session.execute(", "select(")
	}
	return false
}

func perfUsesOnlyScalar(body, variable string) bool {
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(variable) + `\.([A-Za-z_][A-Za-z0-9_]*)`)
	matches := re.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 || len(matches) > 2 {
		return false
	}
	for _, match := range matches {
		if containsAnyFold(match[1], "author", "customer", "profile", "relation", "parent", "owner", "account", "organization") {
			return false
		}
	}
	return true
}

func linesText(lines []codeLine, start, end int) string {
	start, end = safeLineRange(lines, start, end)
	if start >= end {
		return ""
	}
	parts := make([]string, 0, end-start)
	for _, line := range lines[start:end] {
		parts = append(parts, line.text)
	}
	return strings.Join(parts, "\n")
}

func containsAnyFold(value string, needles ...string) bool {
	lower := strings.ToLower(value)
	for _, needle := range needles {
		if strings.Contains(lower, strings.ToLower(needle)) {
			return true
		}
	}
	return false
}

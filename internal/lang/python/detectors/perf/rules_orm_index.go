package perf

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("PERF-PY-15", detectPERFPY15)
	RegisterRule("PERF-PY-16", detectPERFPY16)
	RegisterRule("PERF-PY-20", detectPERFPY20)
}

// PERF-PY-15 reports only when the queried model and its complete local declaration are visible.
func detectPERFPY15(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if perfExternalIndexEvidence(facts.Source) {
		return
	}
	for _, line := range facts.lines {
		model, ok := perfQueryModel(line.text)
		if !ok || !strings.Contains(line.text, ".filter(") || !perfCompositeTenantTimeFilter(line.text) {
			continue
		}
		decl, visible := perfModelDeclaration(facts, model)
		if !visible || perfDeclHasIndex(decl) {
			continue
		}
		pushLine(unit, "PERF-PY-15", line, ".filter(", "tenant/time filter has no visible supporting model index (review-level)", out)
	}
}

// PERF-PY-16 has the same declaration-scope guard as PERF-PY-15.
func detectPERFPY16(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if perfExternalIndexEvidence(facts.Source) {
		return
	}
	for _, line := range facts.lines {
		model, ok := perfQueryModel(line.text)
		if !ok || !strings.Contains(line.text, ".filter(") || !perfRetentionPredicate(line.text) {
			continue
		}
		decl, visible := perfModelDeclaration(facts, model)
		if !visible || perfDeclHasIndex(decl) {
			continue
		}
		pushLine(unit, "PERF-PY-16", line, ".filter(", "retention timestamp predicate has no visible supporting model index (review-level)", out)
	}
}

// PERF-PY-20 requires an equality filter plus sort and a complete local model declaration.
func detectPERFPY20(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if perfExternalIndexEvidence(facts.Source) {
		return
	}
	for _, line := range facts.lines {
		model, ok := perfQueryModel(line.text)
		if !ok || !strings.Contains(line.text, ".filter(") || !strings.Contains(line.text, ".order_by(") || strings.Contains(line.text, "pk=") {
			continue
		}
		decl, visible := perfModelDeclaration(facts, model)
		if !visible || perfDeclHasIndex(decl) {
			continue
		}
		pushLine(unit, "PERF-PY-20", line, ".order_by(", "filtered ORM sort has no visible supporting composite index (review-level)", out)
	}
}

var perfORMQueryModel = regexp.MustCompile(`\b([A-Z][A-Za-z0-9_]*)\.objects\.`)
var perfModelClass = regexp.MustCompile(`^class\s+([A-Z][A-Za-z0-9_]*)\s*\(`)

func perfQueryModel(line string) (string, bool) {
	matches := perfORMQueryModel.FindStringSubmatch(line)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

func perfModelDeclaration(facts *pyPerfFacts, model string) (string, bool) {
	for i, line := range facts.lines {
		matches := perfModelClass.FindStringSubmatch(strings.TrimSpace(line.text))
		if len(matches) != 2 || matches[1] != model || !strings.Contains(line.text, "models.Model") {
			continue
		}
		baseIndent := indentWidth(line.raw)
		end := len(facts.lines)
		for j := i + 1; j < len(facts.lines); j++ {
			trimmed := strings.TrimSpace(facts.lines[j].text)
			if trimmed != "" && indentWidth(facts.lines[j].raw) <= baseIndent && (strings.HasPrefix(trimmed, "class ") || strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ")) {
				end = j
				break
			}
		}
		return linesText(facts.lines, i, end), true
	}
	return "", false
}

// Any local index declaration is treated as sufficient. This intentionally prefers a miss
// over claiming that a deployment lacks an index we cannot prove is absent.
func perfDeclHasIndex(declaration string) bool {
	return containsAnyFold(declaration, "models.index(", "index(", "db_index=true", "__table_args__", "partition")
}

func perfExternalIndexEvidence(source string) bool {
	return containsAnyFold(source, "external migration", "managed externally", "external index", "migration-managed")
}

func perfCompositeTenantTimeFilter(line string) bool {
	lower := strings.ToLower(line)
	hasTenant := strings.Contains(lower, "tenant") || strings.Contains(lower, "owner")
	hasTime := strings.Contains(lower, "timestamp__") || strings.Contains(lower, "created_at__") || strings.Contains(lower, "occurred_at__") || strings.Contains(lower, "expires_at__")
	return hasTenant && hasTime
}

func perfRetentionPredicate(line string) bool {
	lower := strings.ToLower(line)
	return (strings.Contains(lower, "created_at__lt") || strings.Contains(lower, "created_at__lte") || strings.Contains(lower, "expires_at__lt") || strings.Contains(lower, "timestamp__lt")) && containsAnyFold(lower, "cutoff", "retention", "expire", "purge")
}

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
		if !visible || perfDeclCoversQuery(decl, perfFilterColumns(line.text), nil) {
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
		if !visible || perfDeclCoversQuery(decl, perfFilterColumns(line.text), nil) {
			continue
		}
		pushLine(unit, "PERF-PY-16", line, ".filter(", "retention timestamp predicate has no visible supporting model index (review-level)", out)
	}
}

var perfPKConstraintRE = regexp.MustCompile(`\b(?:pk|id)\s*=`)

// PERF-PY-20 requires an equality filter plus sort and a complete local model declaration.
// Primary-key-constrained queries are suppressed (catalogue polarity).
func detectPERFPY20(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if perfExternalIndexEvidence(facts.Source) {
		return
	}
	for _, line := range facts.lines {
		model, ok := perfQueryModel(line.text)
		if !ok || !strings.Contains(line.text, ".filter(") || !strings.Contains(line.text, ".order_by(") || perfPKConstraintRE.MatchString(line.text) {
			continue
		}
		decl, visible := perfModelDeclaration(facts, model)
		filterCols := perfFilterColumns(line.text)
		orderCols := perfOrderByColumns(line.text)
		if !visible || perfDeclCoversQuery(decl, filterCols, orderCols) {
			continue
		}
		pushLine(unit, "PERF-PY-20", line, ".order_by(", "filtered ORM sort has no visible supporting composite index (review-level)", out)
	}
}

var perfORMQueryModel = regexp.MustCompile(`\b([A-Z][A-Za-z0-9_]*)\.objects\.`)
var perfModelClass = regexp.MustCompile(`^class\s+([A-Z][A-Za-z0-9_]*)\s*\(`)
var perfFilterKwRE = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=`)
var perfOrderByRE = regexp.MustCompile(`\.order_by\s*\(([^)]*)\)`)
var perfIndexFieldsRE = regexp.MustCompile(`(?i)(?:models\.)?Index\s*\(\s*fields\s*=\s*\[([^\]]*)\]`)
var perfDBIndexFieldRE = regexp.MustCompile(`(?m)^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*models\.\w+\([^;\n]*db_index\s*=\s*True`)
var perfFieldNameRE = regexp.MustCompile(`["'](-?[A-Za-z_][A-Za-z0-9_]*)["']`)

func perfQueryModel(line string) (string, bool) {
	matches := perfORMQueryModel.FindStringSubmatch(line)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

func perfModelDeclaration(facts *pyPerfFacts, model string) (string, bool) {
	for i, line := range facts.lines {
		matches := perfModelClass.FindStringSubmatch(line.trim)
		if len(matches) != 2 || matches[1] != model || !strings.Contains(line.text, "models.Model") {
			continue
		}
		baseIndent := indentWidth(line.raw)
		end := len(facts.lines)
		for j := i + 1; j < len(facts.lines); j++ {
			trimmed := facts.lines[j].trim
			if trimmed != "" && indentWidth(facts.lines[j].raw) <= baseIndent && (strings.HasPrefix(trimmed, "class ") || strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ")) {
				end = j
				break
			}
		}
		return linesText(facts.lines, i, end), true
	}
	return "", false
}

// perfDeclCoversQuery reports whether the model declaration has a visible index
// that covers filter columns (and optional order_by columns as a composite prefix).
// FN-safe: unparseable index markers still suppress findings.
func perfDeclCoversQuery(declaration string, filterCols, orderCols []string) bool {
	needed := append([]string{}, filterCols...)
	needed = append(needed, orderCols...)
	needed = normalizeIndexCols(needed)
	if len(needed) == 0 {
		// No parseable query columns → prefer miss over claiming a missing index.
		return perfDeclHasIndexMarker(declaration)
	}

	indexes, opaque := perfParseIndexes(declaration)
	if opaque {
		return true
	}
	if len(indexes) == 0 {
		return false
	}
	for _, idx := range indexes {
		if indexCovers(idx, needed) {
			return true
		}
	}
	return false
}

func perfDeclHasIndexMarker(declaration string) bool {
	return containsAnyFold(declaration, "models.index(", "index(", "db_index=true", "__table_args__", "partition")
}

func perfParseIndexes(declaration string) (indexes [][]string, opaque bool) {
	lower := strings.ToLower(declaration)
	if strings.Contains(lower, "__table_args__") || strings.Contains(lower, "partition") {
		// Cannot reliably parse SQLAlchemy/table-args shapes → FN-safe suppress.
		return nil, true
	}

	for _, m := range perfIndexFieldsRE.FindAllStringSubmatch(declaration, -1) {
		cols := extractQuotedNames(m[1])
		if len(cols) == 0 {
			opaque = true
			continue
		}
		indexes = append(indexes, normalizeIndexCols(cols))
	}
	for _, m := range perfDBIndexFieldRE.FindAllStringSubmatch(declaration, -1) {
		indexes = append(indexes, normalizeIndexCols([]string{m[1]}))
	}

	// Index( / Meta.indexes present but no fields= parsed → opaque.
	if (strings.Contains(lower, "models.index(") || strings.Contains(lower, "indexes") || strings.Contains(lower, "db_index=true")) &&
		len(indexes) == 0 {
		return nil, true
	}
	return indexes, opaque
}

func indexCovers(index, needed []string) bool {
	if len(needed) == 0 || len(index) < len(needed) {
		return false
	}
	// Covering prefix: needed columns appear as a prefix of the index (order-sensitive for composites).
	for i, col := range needed {
		if index[i] != col {
			return false
		}
	}
	return true
}

func normalizeIndexCols(cols []string) []string {
	out := make([]string, 0, len(cols))
	seen := map[string]struct{}{}
	for _, c := range cols {
		c = strings.TrimSpace(c)
		c = strings.TrimPrefix(c, "-")
		if at := strings.Index(c, "__"); at > 0 {
			c = c[:at]
		}
		c = strings.ToLower(c)
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		out = append(out, c)
	}
	return out
}

func extractQuotedNames(inner string) []string {
	var cols []string
	for _, m := range perfFieldNameRE.FindAllStringSubmatch(inner, -1) {
		cols = append(cols, m[1])
	}
	return cols
}

func perfFilterColumns(line string) []string {
	start := strings.Index(line, ".filter(")
	if start < 0 {
		return nil
	}
	open := start + len(".filter(") - 1
	inner, ok := balancedCallArgs(line, open)
	if !ok {
		inner = line[start:]
	}
	var cols []string
	for _, m := range perfFilterKwRE.FindAllStringSubmatch(inner, -1) {
		name := m[1]
		switch strings.ToLower(name) {
		case "true", "false", "none":
			continue
		}
		cols = append(cols, name)
	}
	return normalizeIndexCols(cols)
}

func perfOrderByColumns(line string) []string {
	m := perfOrderByRE.FindStringSubmatch(line)
	if len(m) != 2 {
		return nil
	}
	return normalizeIndexCols(extractQuotedNames(m[1]))
}

func balancedCallArgs(src string, openParen int) (string, bool) {
	if openParen < 0 || openParen >= len(src) || src[openParen] != '(' {
		return "", false
	}
	depth := 0
	inQuote := byte(0)
	escaped := false
	for i := openParen; i < len(src); i++ {
		c := src[i]
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
		switch c {
		case '\'', '"':
			inQuote = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return src[openParen+1 : i], true
			}
		}
	}
	return "", false
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

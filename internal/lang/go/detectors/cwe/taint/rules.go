package taint

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

// Rule detection for taint-core CWEs. Metadata is supplied by the caller
// (cwe package owns catalogue Meta vars).

// DetectCWE22 finds path-traversal flows (UserInput → FileOpen without Path sanitizer).
func DetectCWE22(unit *core.ParsedUnit, g *TaintGraph, meta *rules.RuleMetadata, out *[]rules.Finding) {
	if g == nil || meta == nil || unit == nil {
		return
	}
	file := displayPath(unit)
	paths := FindTaintPaths(g, SourceUserInput, SinkFileOpen, []SanitizerKind{SanitizerPath})
	for _, path := range paths {
		if path.Sanitized {
			continue
		}
		if !isFirstArgTainted(g, path) {
			continue
		}
		sink := g.Nodes[path.SinkID]
		line, col := unit.LineCol(sink.ByteRange.Start)
		rules.PushFindingWithConfidence(meta, file, line, col,
			"user-controlled input reaches a file-open sink without path sanitization",
			0.75, out)
	}
}

// DetectCWE78 finds command-injection flows.
func DetectCWE78(unit *core.ParsedUnit, g *TaintGraph, meta *rules.RuleMetadata, out *[]rules.Finding) {
	if g == nil || meta == nil || unit == nil {
		return
	}
	file := displayPath(unit)
	// Path sanitizers do NOT count for command injection.
	paths := FindTaintPaths(g, SourceUserInput, SinkCommandExec, []SanitizerKind{SanitizerValidation, SanitizerBounded})
	for _, path := range paths {
		if path.Sanitized {
			continue
		}
		sink := g.Nodes[path.SinkID]
		// exec.Command without shell argv is not shell injection.
		if sink.Function == "exec.Command" && !hasShellArgs(unit.Source, sink.ByteRange) {
			continue
		}
		line, col := unit.LineCol(sink.ByteRange.Start)
		rules.PushFindingWithConfidence(meta, file, line, col,
			"user-controlled input reaches a shell command execution sink",
			0.75, out)
	}
}

// DetectCWE79 finds XSS flows into Template / HTTPWrite without HTML sanitizer.
func DetectCWE79(unit *core.ParsedUnit, g *TaintGraph, meta *rules.RuleMetadata, out *[]rules.Finding) {
	if g == nil || meta == nil || unit == nil {
		return
	}
	file := displayPath(unit)
	for _, item := range []struct {
		kind SinkKind
		msg  string
	}{
		{SinkTemplate, "user-controlled input reaches a template execution sink without HTML escaping"},
		{SinkHTTPWrite, "user-controlled input reaches an HTTP write sink without HTML escaping"},
	} {
		paths := FindTaintPaths(g, SourceUserInput, item.kind, []SanitizerKind{SanitizerHTML})
		for _, path := range paths {
			if path.Sanitized {
				continue
			}
			sink := g.Nodes[path.SinkID]
			line, col := unit.LineCol(sink.ByteRange.Start)
			rules.PushFindingWithConfidence(meta, file, line, col, item.msg, 0.7, out)
		}
	}
}

// DetectCWE89 finds SQL injection flows.
func DetectCWE89(unit *core.ParsedUnit, g *TaintGraph, ann *TaintAnnotations, meta *rules.RuleMetadata, out *[]rules.Finding) {
	if g == nil || meta == nil || unit == nil {
		return
	}
	file := displayPath(unit)
	paths := FindTaintPaths(g, SourceUserInput, SinkSQLQuery, []SanitizerKind{SanitizerSQL})
	for _, path := range paths {
		if path.Sanitized {
			continue
		}
		sink := g.Nodes[path.SinkID]
		if isParameterizedQuery(unit.Source, sink.ByteRange) {
			continue
		}
		if isPreparedStmtParameterized(sink.Function, sink.ByteRange, ann) {
			continue
		}
		line, col := unit.LineCol(sink.ByteRange.Start)
		rules.PushFindingWithConfidence(meta, file, line, col,
			"user-controlled input reaches an SQL execution sink (heuristic; not full SQLi coverage)",
			0.7, out)
	}
}

func displayPath(unit *core.ParsedUnit) string {
	if unit.DisplayPath != "" {
		return unit.DisplayPath
	}
	return unit.Path
}

func isFirstArgTainted(g *TaintGraph, path TaintPath) bool {
	for _, nodeID := range path.NodeIDs {
		for _, e := range g.Edges {
			if e.From == nodeID && e.To == path.SinkID && e.Kind.IsArgument() && e.Kind.ArgumentIndex() == 0 {
				return true
			}
		}
	}
	return false
}

func hasShellArgs(source string, br ByteRange) bool {
	if br.Start < 0 || br.End > len(source) || br.Start >= br.End {
		return false
	}
	call := CompactWhitespace(source[br.Start:br.End])
	return strings.Contains(call, `"sh","-c"`) || strings.Contains(call, `"bash","-c"`)
}

func isParameterizedQuery(source string, br ByteRange) bool {
	if br.Start < 0 || br.End > len(source) || br.Start >= br.End {
		return false
	}
	call := source[br.Start:br.End]
	p := strings.Index(call, "(")
	if p < 0 {
		return false
	}
	first, ok := NthTopLevelArg(call[p+1:], 0)
	return ok && IsPureStringLiteral(first)
}

func isPreparedStmtParameterized(sinkFn string, sinkRange ByteRange, ann *TaintAnnotations) bool {
	if ann == nil {
		return false
	}
	recv, method, ok := strings.Cut(sinkFn, ".")
	if !ok || strings.Contains(recv, ".") || strings.Contains(recv, "(") || strings.Contains(recv, "[") {
		return false
	}
	switch method {
	case "Query", "Exec", "QueryRow", "QueryContext", "ExecContext", "QueryRowContext":
	default:
		return false
	}
	// Find enclosing function range.
	var fnRange *ByteRange
	for _, r := range ann.FunctionRanges {
		rr := r
		if rr.Start <= sinkRange.Start && sinkRange.End <= rr.End {
			fnRange = &rr
			break
		}
	}
	var latest *AssignmentDetail
	for i := range ann.Assignments {
		a := &ann.Assignments[i]
		if a.LHS != recv {
			continue
		}
		if a.ByteRange.Start >= sinkRange.Start {
			continue
		}
		if fnRange != nil && (a.ByteRange.Start < fnRange.Start || a.ByteRange.Start >= fnRange.End) {
			continue
		}
		if latest == nil || a.ByteRange.Start > latest.ByteRange.Start {
			latest = a
		}
	}
	if latest == nil {
		return false
	}
	return prepareRHSHasLiteralSQL(latest.RHSText)
}

func prepareRHSHasLiteralSQL(rhs string) bool {
	const prepareCtx = ".PrepareContext("
	const prepare = ".Prepare("
	var after string
	var sqlIdx int
	if i := strings.Index(rhs, prepareCtx); i >= 0 {
		after = rhs[i+len(prepareCtx):]
		sqlIdx = 1
	} else if i := strings.Index(rhs, prepare); i >= 0 {
		after = rhs[i+len(prepare):]
		sqlIdx = 0
	} else {
		return false
	}
	arg, ok := NthTopLevelArg(after, sqlIdx)
	return ok && IsPureStringLiteral(arg)
}

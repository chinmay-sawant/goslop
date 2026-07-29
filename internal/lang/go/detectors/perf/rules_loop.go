package perf

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

// detectPERF1: regexp.MustCompile / Compile inside a loop.
func detectPERF1(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		if call.Callee != "regexp.MustCompile" && call.Callee != "regexp.Compile" {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(&MetaPERF1, file, line, col, "regular expression compiled inside loop body", out)
	}
}

// detectPERF2: string concatenation via += inside a loop.
func detectPERF2(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, a := range facts.Assignments {
		if !IsAssignmentInLoop(a) {
			continue
		}
		text, expr := a.Text, a.Expr
		isConcat := strings.Contains(text, " += ") ||
			strings.Contains(text, "= s +") ||
			strings.Contains(expr, "s = s +")
		if !isConcat {
			// also catch `name = name + ...`
			if !strings.Contains(text, " = ") || !strings.Contains(expr, a.Name+" +") {
				if !strings.Contains(text, "+=") {
					continue
				}
			}
		}
		if strings.Contains(text, "strings.Builder") ||
			strings.Contains(text, "bytes.Buffer") ||
			strings.Contains(text, "strings.Join(") {
			continue
		}
		if k, ok := facts.VarKinds[a.Name]; ok && k == VarNumeric {
			continue
		}
		isKnownString := facts.VarKinds[a.Name] == VarString
		hasStringLiteral := strings.Contains(text, "\"") || strings.Contains(text, "`")
		if !isKnownString && !hasStringLiteral {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(
			&MetaPERF2, file, line, col,
			"string is built by repeated concatenation inside a loop body", out,
		)
	}
}

// detectPERF3: make([]T, ...) rebuilt inside a loop without capacity-style hint.
func detectPERF3(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, a := range facts.Assignments {
		if !IsAssignmentInLoop(a) {
			continue
		}
		expr := a.Expr
		if !strings.Contains(expr, "make([]") || !strings.Contains(expr, ",") {
			continue
		}
		if strings.Contains(expr, ", 0, ") {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(
			&MetaPERF3, file, line, col,
			"working slice is rebuilt with make inside a loop body", out,
		)
	}
}

// detectPERF4: make(map[...]) inside a loop without size hint.
func detectPERF4(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, a := range facts.Assignments {
		if !IsAssignmentInLoop(a) {
			continue
		}
		expr := a.Expr
		if !strings.Contains(expr, "make(map[") {
			continue
		}
		// make(map[K]V, hint) has a comma after the closing ]
		if close := strings.Index(expr, "]"); close >= 0 && strings.Contains(expr[close:], ",") {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(
			&MetaPERF4, file, line, col,
			"map is allocated with make inside a loop body", out,
		)
	}
}

// detectPERF5: json marshal/unmarshal in loop.
func detectPERF5(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		switch call.Callee {
		case "json.Marshal", "json.Unmarshal", "json.NewEncoder", "json.NewDecoder":
		default:
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF5, file, line, col,
			"JSON conversion is performed inside a loop body", out,
		)
	}
}

// detectPERF6: fmt.Sprintf / Fprintf in loop.
func detectPERF6(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		if call.Callee != "fmt.Sprintf" && call.Callee != "fmt.Fprintf" {
			continue
		}
		if call.Callee == "fmt.Fprintf" && len(call.Arguments) >= 1 {
			first := call.Arguments[0]
			if first == "&buf" || first == "buf" {
				continue
			}
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF6, file, line, col,
			"fmt-based formatting is performed inside a loop body", out,
		)
	}
}

// detectPERF7: defer inside a loop (not inside a per-iteration closure).
func detectPERF7(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, dr := range facts.DeferStarts {
		startByte := dr[0]
		var bestLoop *[2]int
		bestSpan := int(^uint(0) >> 1)
		for i := range facts.ForRanges {
			fr := &facts.ForRanges[i]
			if fr[0] <= startByte && startByte <= fr[1] {
				span := fr[1] - fr[0]
				if span < bestSpan {
					bestSpan = span
					bestLoop = fr
				}
			}
		}
		if bestLoop == nil {
			continue
		}
		loopStart, loopEnd := bestLoop[0], bestLoop[1]
		// skip defer inside a func literal launched by the loop
		skip := false
		for _, fl := range facts.FunctionLiteralRanges {
			if loopStart <= fl[0] && fl[0] < startByte && startByte < fl[1] && fl[1] <= loopEnd {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		line, col := unit.LineCol(startByte)
		rules.PushFinding(
			&MetaPERF7, file, line, col,
			"defer statement is placed inside a loop body", out,
		)
	}
}

// detectPERF8: time.Parse with literal layout in loop.
func detectPERF8(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		if call.Callee != "time.Parse" && call.Callee != "time.ParseInLocation" {
			continue
		}
		if len(call.Arguments) == 0 {
			continue
		}
		layout := call.Arguments[0]
		if !(strings.HasPrefix(layout, "\"") && strings.HasSuffix(layout, "\"")) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF8, file, line, col,
			"time.Parse is called inside a loop body with a literal layout", out,
		)
	}
}

// detectPERF50: regexp.Match* inside a loop.
func detectPERF50(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	for _, call := range facts.Calls {
		switch call.Callee {
		case "regexp.MatchString", "regexp.Match", "regexp.MatchReader":
		default:
			continue
		}
		if !IsInLoop(call) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF50, file, line, col,
			"regexp match is invoked inside a loop; compile the pattern once and reuse it", out,
		)
		return
	}
}

// detectPERF116: strings.Index compared to -1 (prefer Contains).
func detectPERF116(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	src := unit.Source
	if !strings.Contains(src, "strings.Index") {
		return
	}
	for _, call := range facts.Calls {
		if call.Callee != "strings.Index" {
			continue
		}
		// window: call through end of line
		windowEnd := call.StartByte + 64
		if windowEnd > len(src) {
			windowEnd = len(src)
		}
		if nl := strings.IndexByte(src[call.StartByte:], '\n'); nl >= 0 {
			if call.StartByte+nl > windowEnd {
				windowEnd = call.StartByte + nl
			}
		}
		window := src[call.StartByte:windowEnd]
		if !strings.Contains(window, "!= -1") && !strings.Contains(window, "== -1") &&
			!strings.Contains(window, "!=-1") && !strings.Contains(window, "==-1") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF116, file, line, col,
			"strings.Index(s, sub) compared to -1 should be strings.Contains(s, sub)", out,
		)
	}
}

func unitFile(unit *core.ParsedUnit) string {
	if unit.DisplayPath != "" {
		return unit.DisplayPath
	}
	return unit.Path
}

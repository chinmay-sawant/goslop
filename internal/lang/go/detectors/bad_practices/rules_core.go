package badpractices

import (
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("BP-3", detectBP3)
	RegisterRule("BP-6", detectBP6)
	RegisterRule("BP-7", detectBP7)
	RegisterRule("BP-8", detectBP8)
	RegisterRule("BP-9", detectBP9)
	RegisterRule("BP-10", detectBP10)
	RegisterRule("BP-11", detectBP11)
	RegisterRule("BP-12", detectBP12)
	RegisterRule("BP-13", detectBP13)
	RegisterRule("BP-14", detectBP14)
	RegisterRule("BP-15", detectBP15)
	RegisterRule("BP-72", detectBP72)
	RegisterRule("BP-73", detectBP73)
	RegisterRule("BP-75", detectBP75)
	RegisterRule("BP-79", detectBP79)
	RegisterRule("BP-80", detectBP80)
	RegisterRule("BP-81", detectBP81)
	RegisterRule("BP-82", detectBP82)
	RegisterRule("BP-83", detectBP83)
	RegisterRule("BP-84", detectBP84)
	RegisterRule("BP-86", detectBP86)
	RegisterRule("BP-88", detectBP88)
	RegisterRule("BP-89", detectBP89)
	RegisterRule("BP-90", detectBP90)
	RegisterRule("BP-92", detectBP92)
	RegisterRule("BP-93", detectBP93)
	RegisterRule("BP-94", detectBP94)
	RegisterRule("BP-100", detectBP100)
}

// BP-3: panic outside main/test.
func detectBP3(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-3")
	if isTestFile(unit) || !strings.Contains(unit.Source, "panic(") {
		return
	}
	msg := "panic outside main() or test files; prefer returning errors up the call stack"
	if len(facts.callNodes) > 0 {
		for _, c := range facts.callNodes {
			if c.callee != "panic" {
				continue
			}
			fn := facts.enclosingFunc(c.start)
			if fn != nil && fn.isMain {
				continue
			}
			// also skip if name is main via text
			if name, ok := enclosingFuncName(unit.Source, c.start); ok && name == "main" {
				continue
			}
			pushAt(unit, meta, c.start, msg, out)
		}
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.Contains(t, "panic(") {
			continue
		}
		if name, ok := enclosingFuncName(unit.Source, line.byte); ok && name == "main" {
			continue
		}
		pushAt(unit, meta, line.byte+strings.Index(line.text, "panic("), msg, out)
	}
}

// BP-6: WaitGroup.Add inside goroutine.
func detectBP6(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-6")
	if !strings.Contains(unit.Source, "go ") || !strings.Contains(unit.Source, ".Add(") {
		return
	}
	msg := "WaitGroup.Add inside the goroutine it tracks races with Wait; call Add before go"
	for _, g := range facts.goNodes {
		body := g.text
		if strings.Contains(body, ".Add(") && (strings.Contains(body, "WaitGroup") || strings.Contains(body, "wg.") || strings.Contains(unit.Source, "sync.WaitGroup")) {
			// more precise: .Add( inside the go body
			if idx := strings.Index(body, ".Add("); idx >= 0 {
				pushAt(unit, meta, g.start+idx, msg, out)
			}
		}
	}
	// text fallback for go func() { wg.Add
	if len(facts.goNodes) == 0 {
		src := unit.Source
		for _, needle := range []string{"go func", "go func("} {
			start := 0
			for {
				idx := strings.Index(src[start:], needle)
				if idx < 0 {
					break
				}
				abs := start + idx
				// take a window
				end := abs + 200
				if end > len(src) {
					end = len(src)
				}
				window := src[abs:end]
				if strings.Contains(window, ".Add(") {
					if a := strings.Index(window, ".Add("); a >= 0 {
						pushAt(unit, meta, abs+a, msg, out)
					}
				}
				start = abs + len(needle)
			}
		}
	}
}

// BP-7: mutex passed by value.
func detectBP7(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-7")
	if !strings.Contains(unit.Source, "sync.Mutex") {
		return
	}
	msg := "sync.Mutex is passed by value; pass *sync.Mutex to avoid copying lock state"
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "func ") && strings.Contains(t, " sync.Mutex") && !strings.Contains(t, "*sync.Mutex") {
			pushAt(unit, meta, line.byte+strings.Index(line.text, "sync.Mutex"), msg, out)
		}
	}
}

// BP-8: defer unlock on by-value mutex parameter.
func detectBP8(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-8")
	if !strings.Contains(unit.Source, "defer ") || !strings.Contains(unit.Source, ".Unlock()") {
		return
	}
	// If function has by-value mutex param and defers Unlock on it.
	for _, fn := range facts.funcDecls {
		if !strings.Contains(fn.params, "sync.Mutex") || strings.Contains(fn.params, "*sync.Mutex") {
			continue
		}
		// extract param names with sync.Mutex
		// simple: look for defer x.Unlock in body
		body := unit.Source[fn.bodyS:fn.bodyE]
		if strings.Contains(body, "defer ") && strings.Contains(body, ".Unlock()") {
			if idx := strings.Index(body, ".Unlock()"); idx >= 0 {
				// find defer start
				pushAt(unit, meta, fn.bodyS+idx, "Unlock deferred on a mutex held by value; use *sync.Mutex", out)
			}
		}
	}
}

// BP-9: select without default/timeout/context escape (Rust parity).
// Escape only on direct communication cases: default, time.After/NewTimer,
// ctx.Done()/context.Done(), or bare <-stop/<-done. Nested selects and
// arbitrary .Done() receivers (e.g. localCtx.Done) do not suppress the parent.
func detectBP9(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-9")
	if unit == nil || !strings.Contains(unit.Source, "select") {
		return
	}
	msg := "select can block indefinitely without default, timeout, or context cancellation"
	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		tree := facts.tree
		goast.Inspect(tree.File, func(n goast.Node) bool {
			sel, ok := n.(*goast.SelectStmt)
			if !ok || sel.Body == nil {
				return true
			}
			if !selectStmtHasEscape(tree, sel) {
				pushAt(unit, meta, tree.Offset(sel.Pos()), msg, out)
			}
			return true
		})
		return
	}
	// Text fallback: brace-matched select blocks; no nested-select suppression.
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "select {")
		if idx < 0 {
			break
		}
		abs := start + idx
		end := abs + len("select {")
		depth := 1
		for end < len(src) && depth > 0 {
			switch src[end] {
			case '{':
				depth++
			case '}':
				depth--
			}
			end++
		}
		block := src[abs:end]
		if !selectBlockHasEscape(block) {
			pushAt(unit, meta, abs, msg, out)
		}
		start = end
	}
}

func selectStmtHasEscape(tree interface {
	NodeText(n goast.Node) string
}, sel *goast.SelectStmt) bool {
	for _, stmt := range sel.Body.List {
		cc, ok := stmt.(*goast.CommClause)
		if !ok {
			continue
		}
		if cc.Comm == nil {
			return true // default
		}
		if selectCommHasEscape(tree.NodeText(cc.Comm)) {
			return true
		}
	}
	return false
}

func selectCommHasEscape(comm string) bool {
	// Mirror Rust select_has_escape token set on the communication expression only.
	for _, tok := range []string{
		"time.After(",
		"time.NewTimer(",
		"ctx.Done()",
		"context.Done()",
		"<-stop",
		"<-done",
		"<-ctx.Done()",
	} {
		if strings.Contains(comm, tok) {
			return true
		}
	}
	return false
}

func selectBlockHasEscape(block string) bool {
	if strings.Contains(block, "default:") {
		return true
	}
	// Approximate case-communication lines only (text after "case " until ':').
	for _, line := range strings.Split(block, "\n") {
		t := strings.TrimSpace(line)
		if !strings.HasPrefix(t, "case ") {
			continue
		}
		comm := strings.TrimPrefix(t, "case ")
		if i := strings.Index(comm, ":"); i >= 0 {
			comm = comm[:i]
		}
		if selectCommHasEscape(comm) {
			return true
		}
	}
	return false
}

// BP-10: time.After inside loop.
func detectBP10(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-10")
	if !strings.Contains(unit.Source, "time.After") {
		return
	}
	msg := "time.After inside a loop allocates a new timer per iteration"
	if len(facts.callNodes) > 0 {
		for _, c := range facts.callNodes {
			if c.callee == "time.After" && facts.insideLoop(c.start) {
				pushAt(unit, meta, c.start, msg, out)
			}
		}
		return
	}
	for _, line := range codeLines(unit.Source) {
		if strings.Contains(line.text, "time.After") {
			// crude: if any for earlier in file without closing — weak; use facts.forRanges if empty skip unless for in source
			if strings.Contains(unit.Source, "for ") {
				pushAt(unit, meta, line.byte+strings.Index(line.text, "time.After"), msg, out)
			}
		}
	}
}

// BP-11: defer inside loop.
func detectBP11(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-11")
	if !strings.Contains(unit.Source, "defer ") {
		return
	}
	msg := "defer inside a loop defers cleanup until the surrounding function returns"
	if len(facts.deferNodes) > 0 {
		for _, d := range facts.deferNodes {
			if facts.insideLoop(d.start) {
				pushAt(unit, meta, d.start, msg, out)
			}
		}
		return
	}
	// text: look for for { ... defer
	inFor := 0
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "for ") || t == "for {" || strings.HasPrefix(t, "for{") {
			inFor++
		}
		if strings.HasPrefix(t, "defer ") && inFor > 0 {
			pushAt(unit, meta, line.byte, msg, out)
		}
		// rough brace tracking
		inFor += strings.Count(t, "{") - strings.Count(t, "}")
		if inFor < 0 {
			inFor = 0
		}
	}
}

// BP-12: unbuffered channel send from multiple goroutines — light heuristic.
func detectBP12(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-12")
	src := unit.Source
	if !strings.Contains(src, "make(chan ") || !strings.Contains(src, "go ") {
		return
	}
	// unbuffered: make(chan T) without size
	if !strings.Contains(src, "make(chan ") {
		return
	}
	// if multiple go and channel send without buffer
	if strings.Count(src, "go ") >= 2 && strings.Contains(src, "<-") {
		// only fire if make(chan X) without comma capacity
		for _, line := range codeLines(src) {
			t := strings.TrimSpace(line.text)
			if strings.Contains(t, "make(chan ") && !strings.Contains(t, ",") {
				pushAt(unit, meta, line.byte, "unbuffered channel receives sends from multiple goroutines; consider buffering or synchronizing", out)
				return
			}
		}
	}
}

// BP-13: context.Background in library.
func detectBP13(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-13")
	if isTestFile(unit) || !strings.Contains(unit.Source, "context.Background") {
		return
	}
	msg := "context.Background used in library code; accept and propagate a caller context"
	if len(facts.callNodes) > 0 {
		for _, c := range facts.callNodes {
			if c.callee != "context.Background" {
				continue
			}
			fn := facts.enclosingFunc(c.start)
			if fn != nil && (fn.isMain || fn.name == "init") {
				continue
			}
			if name, ok := enclosingFuncName(unit.Source, c.start); ok && (name == "main" || name == "init") {
				continue
			}
			pushAt(unit, meta, c.start, msg, out)
		}
		return
	}
	for _, line := range codeLines(unit.Source) {
		if strings.Contains(line.text, "context.Background") {
			if name, ok := enclosingFuncName(unit.Source, line.byte); ok && (name == "main" || name == "init") {
				continue
			}
			pushAt(unit, meta, line.byte+strings.Index(line.text, "context.Background"), msg, out)
		}
	}
}

// BP-14: goroutine without ctx cancellation.
func detectBP14(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-14")
	if !strings.Contains(unit.Source, "go ") && !strings.Contains(unit.Source, "go func") {
		return
	}
	msg := "goroutine launched without observing context cancellation"
	for _, g := range facts.goNodes {
		body := g.text
		if strings.Contains(body, "ctx.Done()") || strings.Contains(body, "context.WithCancel") ||
			strings.Contains(body, "<-ctx.Done") || strings.Contains(body, ".Done()") {
			continue
		}
		// long-running shape: for { without done
		if strings.Contains(body, "for ") || strings.Contains(body, "for{") {
			pushAt(unit, meta, g.start, msg, out)
		}
	}
	if len(facts.goNodes) == 0 && strings.Contains(unit.Source, "go func") {
		// crude: go func with for and no Done
		if strings.Contains(unit.Source, "go func") &&
			(strings.Contains(unit.Source, "for {") || strings.Contains(unit.Source, "for{")) &&
			!strings.Contains(unit.Source, "ctx.Done()") &&
			!strings.Contains(unit.Source, ".Done()") {
			if pos := strings.Index(unit.Source, "go func"); pos >= 0 {
				pushAt(unit, meta, pos, msg, out)
			}
		}
	}
}

// BP-15: recursive Once.Do — light heuristic.
func detectBP15(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-15")
	if !strings.Contains(unit.Source, ".Do(") || !strings.Contains(unit.Source, "sync.Once") {
		return
	}
	// Flag Once.Do that calls itself by name inside
	// hard; skip unless obvious once.Do(func(){ once.Do
	if strings.Contains(unit.Source, ".Do(func") {
		// look for nested .Do(
		src := unit.Source
		if strings.Count(src, ".Do(") >= 2 {
			if pos := strings.Index(src, ".Do("); pos >= 0 {
				pushAt(unit, meta, pos, "sync.Once.Do appears nested/recursive; avoid re-entering the same Once", out)
			}
		}
	}
}

func detectBP72(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-72")
	// return nil for interface when local *T is nil — common pattern: var p *T; return p
	src := unit.Source
	if !strings.Contains(src, "return ") {
		return
	}
	// function returns interface/error interface and returns typed nil
	// heuristic: var x *Something\n return x
	lines := codeLines(src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "var ") && strings.Contains(t, "*") {
			// next return same name
			parts := strings.Fields(t)
			if len(parts) < 3 {
				continue
			}
			name := parts[1]
			for j := i + 1; j < len(lines) && j < i+6; j++ {
				rt := strings.TrimSpace(lines[j].text)
				if rt == "return "+name {
					pushAt(unit, meta, lines[j].byte, "typed nil returned as interface; return an untyped nil instead", out)
				}
			}
		}
	}
}

func detectBP73(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-73")
	// m[k] = v without make — hard; flag var m map[...] then assignment
	if !strings.Contains(unit.Source, "map[") {
		return
	}
	varMaps := map[string]int{}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "var ") && strings.Contains(t, "map[") && !strings.Contains(t, "make(") {
			parts := strings.Fields(t)
			if len(parts) >= 2 {
				varMaps[parts[1]] = line.byte
			}
		}
		// m[k] =
		for name, pos := range varMaps {
			if strings.Contains(t, name+"[") && strings.Contains(t, "=") && !strings.Contains(t, "==") {
				pushAt(unit, meta, pos, "write to a nil map without initialization; make the map first", out)
				delete(varMaps, name)
			}
		}
	}
}

func detectBP75(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-75")
	if !strings.Contains(unit.Source, "copy(") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		// copy(dst, src) where dst is make([]T, 0) or nil
		if strings.Contains(t, "copy(") && (strings.Contains(unit.Source, "make([]") && strings.Contains(unit.Source, ", 0)")) {
			// weak
			if strings.Contains(t, "copy(") {
				// look for zero-length destination nearby
				pushAt(unit, meta, line.byte, "copy into a zero-length slice does nothing; allocate with length or use append", out)
				return
			}
		}
	}
}

func detectBP79(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-79")
	if !strings.Contains(unit.Source, "context.WithCancel") && !strings.Contains(unit.Source, "context.WithTimeout") && !strings.Contains(unit.Source, "context.WithDeadline") {
		return
	}
	// cancel not deferred
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "context.WithCancel") || strings.Contains(t, "context.WithTimeout") || strings.Contains(t, "context.WithDeadline") {
			// if cancel assigned and no defer cancel nearby in function
			if strings.Contains(t, "cancel") || strings.Contains(t, ",") {
				// check whole source for defer cancel
				if !strings.Contains(unit.Source, "defer cancel") && !strings.Contains(unit.Source, "defer cancel(") {
					pushAt(unit, meta, line.byte, "context cancel function is not deferred; resource may leak", out)
					return
				}
			}
		}
	}
}

func detectBP80(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-80")
	if isTestFile(unit) || !strings.Contains(unit.Source, "context.TODO") {
		return
	}
	msg := "context.TODO used in production code; pass a real context or context.Background at the edge"
	if len(facts.callNodes) > 0 {
		for _, c := range facts.callNodes {
			if c.callee == "context.TODO" {
				fn := facts.enclosingFunc(c.start)
				if fn != nil && fn.isMain {
					continue
				}
				pushAt(unit, meta, c.start, msg, out)
			}
		}
		return
	}
	if pos := strings.Index(unit.Source, "context.TODO"); pos >= 0 {
		pushAt(unit, meta, pos, msg, out)
	}
}

func detectBP81(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-81")
	if strings.Count(unit.Source, "time.Now()") >= 2 {
		// in condition
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if strings.Contains(t, "time.Now()") && (strings.Contains(t, "if ") || strings.Contains(t, "for ")) {
				pushAt(unit, meta, line.byte, "time.Now() evaluated repeatedly in a condition; capture once", out)
				return
			}
		}
	}
}

func detectBP82(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-82")
	if !strings.Contains(unit.Source, "time.Parse(") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "time.Parse(") && !strings.Contains(t, "ParseInLocation") && !strings.Contains(unit.Source, "time.ParseInLocation") {
			// time.Parse is UTC-only for no zone; flag if no location elsewhere
			if !strings.Contains(unit.Source, "time.Local") && !strings.Contains(unit.Source, "time.UTC") {
				pushAt(unit, meta, line.byte, "time.Parse without location; prefer ParseInLocation for local times", out)
				return
			}
		}
	}
}

// BP-83: time.Sleep used as synchronization without a coordination primitive.
// Rust only flags sync-shaped function names (wait/ready/sync/…) or go-launched
// closures, excludes backoff/retry-shaped bodies, and requires no visible
// channel/lock/atomic boundary in the same function.
func detectBP83(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-83")
	if isTestFile(unit) {
		return
	}
	if !strings.Contains(unit.Source, "time.Sleep(") {
		return
	}
	msg := "time.Sleep is being used as synchronization without a visible coordination primitive"
	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		tree := facts.tree
		goast.Inspect(tree.File, func(n goast.Node) bool {
			call, ok := n.(*goast.CallExpr)
			if !ok {
				return true
			}
			if tree.NodeText(call.Fun) != "time.Sleep" {
				return true
			}
			off := tree.Offset(call.Pos())
			name, text, goLaunched, ok := enclosingFuncTextBP83(facts, off)
			if !ok {
				return true
			}
			if !funcNameIsSyncShape(name) && !goLaunched {
				return true
			}
			if containsVisibleSyncBP83(text) || isBackoffOrRetryBP83(text) {
				return true
			}
			pushAt(unit, meta, off, msg, out)
			return true
		})
		return
	}
	// Text fallback: only sleep inside functions whose names look sync-shaped.
	for _, line := range codeLines(unit.Source) {
		if !strings.Contains(line.text, "time.Sleep(") {
			continue
		}
		name, ok := enclosingFuncName(unit.Source, line.byte)
		if !ok || !funcNameIsSyncShape(name) {
			continue
		}
		pushAt(unit, meta, line.byte, msg, out)
	}
}

// enclosingFuncTextBP83 finds the function/method/literal enclosing offset and returns
// (func name or "", full func source, go-launched?, ok).
func enclosingFuncTextBP83(facts *bpFacts, offset int) (name, text string, goLaunched, ok bool) {
	if facts == nil || facts.tree == nil || facts.tree.File == nil {
		return "", "", false, false
	}
	tree := facts.tree
	var best goast.Node
	var bestName string
	bestSpan := int(^uint(0) >> 1)
	goast.Inspect(tree.File, func(n goast.Node) bool {
		if n == nil {
			return true
		}
		start := tree.Offset(n.Pos())
		end := tree.Offset(n.End())
		if offset < start || offset >= end {
			return true
		}
		span := end - start
		switch x := n.(type) {
		case *goast.FuncDecl:
			if span <= bestSpan {
				best = x
				bestSpan = span
				if x.Name != nil {
					bestName = x.Name.Name
				} else {
					bestName = ""
				}
			}
		case *goast.FuncLit:
			if span < bestSpan {
				best = x
				bestSpan = span
				bestName = ""
			}
		}
		return true
	})
	if best == nil {
		return "", "", false, false
	}
	text = tree.NodeText(best)
	if _, isLit := best.(*goast.FuncLit); isLit {
		pos := tree.Offset(best.Pos())
		for _, g := range facts.goNodes {
			if pos >= g.start && pos < g.end {
				goLaunched = true
				break
			}
		}
	}
	return bestName, text, goLaunched, true
}

func funcNameIsSyncShape(name string) bool {
	lower := strings.ToLower(name)
	for _, needle := range []string{
		"wait", "ready", "sync", "until", "drain", "flush", "shutdown", "startup",
		"start", "stop", "signal", "notify", "done",
	} {
		if strings.Contains(lower, needle) {
			return true
		}
	}
	return false
}

func containsVisibleSyncBP83(text string) bool {
	for _, needle := range []string{
		"<-", "select", ".Wait(", ".Lock(", ".Unlock(", "sync.", "atomic.", "Cond", "Once", "close(",
	} {
		if strings.Contains(text, needle) {
			return true
		}
	}
	return false
}

func isBackoffOrRetryBP83(text string) bool {
	lower := strings.ToLower(text)
	for _, needle := range []string{
		"backoff", "retry", "throttle", "rate_limit", "ratelimit", "jitter",
		"cooldown", "debounce", "poll", "periodic", "interval", "timeout", "delay",
	} {
		if strings.Contains(lower, needle) {
			return true
		}
	}
	return false
}

func detectBP84(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-84")
	// x * 100 / y integer percentage
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "* 100 /") || strings.Contains(t, "*100/") || strings.Contains(t, "* 100/") {
			pushAt(unit, meta, line.byte, "integer percentage truncation; multiply carefully or use floating math", out)
		}
	}
}

func detectBP86(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-86")
	if !strings.Contains(unit.Source, ".Lock()") {
		return
	}
	// Lock without Unlock/defer Unlock in same function — crude
	if strings.Contains(unit.Source, ".Lock()") && !strings.Contains(unit.Source, ".Unlock()") {
		if pos := strings.Index(unit.Source, ".Lock()"); pos >= 0 {
			pushAt(unit, meta, pos, "mutex Lock without a corresponding Unlock", out)
		}
	}
}

func detectBP88(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-88")
	// var ch chan T without make then send/recv
	varChans := map[string]int{}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "var ") && strings.Contains(t, "chan ") && !strings.Contains(t, "make(") {
			parts := strings.Fields(t)
			if len(parts) >= 2 {
				varChans[parts[1]] = line.byte
			}
		}
		for name, pos := range varChans {
			if strings.Contains(t, "<-"+name) || strings.Contains(t, name+"<-") || strings.Contains(t, "<- "+name) || strings.Contains(t, name+" <-") {
				pushAt(unit, meta, pos, "send/receive on a nil channel blocks forever", out)
				delete(varChans, name)
			}
		}
	}
}

func detectBP89(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-89")
	// close(ch) more than once on same channel name
	closes := map[string]int{}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "close(") {
			inner := strings.TrimSuffix(strings.TrimPrefix(t, "close("), ")")
			inner = strings.TrimSpace(strings.Split(inner, ")")[0])
			if prev, ok := closes[inner]; ok {
				pushAt(unit, meta, line.byte, "channel closed more than once", out)
				_ = prev
			} else {
				closes[inner] = line.byte
			}
		}
	}
}

// detectBP90: infinite `for {` with a channel receive and no local select/break/return.
// Mirrors Rust detect_bp_90_channel_receive_loop_without_exit (package-level funcs only;
// nested func literals are not walked, matching Rust inspect_functions).
func detectBP90(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-90")
	msg := "infinite loop receives from a channel without a visible local exit"

	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		for _, decl := range facts.tree.File.Decls {
			fd, ok := decl.(*goast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}
			goast.Inspect(fd.Body, func(n goast.Node) bool {
				if n == nil {
					return true
				}
				// Rust walk_scope stops at nested functions/func literals.
				if _, isLit := n.(*goast.FuncLit); isLit {
					return false
				}
				fs, ok := n.(*goast.ForStmt)
				if !ok {
					return true
				}
				// Infinite loop: for { ... } — no init/cond/post.
				if fs.Init != nil || fs.Cond != nil || fs.Post != nil {
					return true
				}
				text := facts.tree.NodeText(fs)
				trimmed := strings.TrimLeft(text, " \t\r\n")
				if !strings.HasPrefix(trimmed, "for {") && !strings.HasPrefix(trimmed, "for{") {
					return true
				}
				// Visible local escape (substring checks, Rust parity).
				if strings.Contains(text, "select") ||
					strings.Contains(text, "break") ||
					strings.Contains(text, "return") {
					return true
				}
				if !forContainsChannelReceive(fs) {
					return true
				}
				pushAt(unit, meta, facts.tree.Offset(fs.Pos()), msg, out)
				return true
			})
		}
		return
	}

	// Text fallback: per-`for {` block with brace matching.
	detectBP90Text(unit, meta, msg, out)
}

func forContainsChannelReceive(n goast.Node) bool {
	found := false
	goast.Inspect(n, func(c goast.Node) bool {
		if found || c == nil {
			return false
		}
		// Nested funcs are outside local scope for the outer for (Rust walk_scope).
		if _, isLit := c.(*goast.FuncLit); isLit {
			return false
		}
		if ue, ok := c.(*goast.UnaryExpr); ok && ue.Op == token.ARROW {
			found = true
			return false
		}
		return true
	})
	return found
}

func detectBP90Text(unit *core.ParsedUnit, meta *rules.RuleMetadata, msg string, out *[]rules.Finding) {
	src := unit.Source
	// Find `for {` / `for{` occurrences and take the statement via brace match.
	for i := 0; i < len(src); {
		// cheap scan for "for"
		idx := strings.Index(src[i:], "for")
		if idx < 0 {
			return
		}
		pos := i + idx
		// word boundary
		if pos > 0 {
			prev := src[pos-1]
			if (prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z') || (prev >= '0' && prev <= '9') || prev == '_' {
				i = pos + 3
				continue
			}
		}
		rest := strings.TrimLeft(src[pos+3:], " \t")
		if !strings.HasPrefix(rest, "{") {
			i = pos + 3
			continue
		}
		// locate opening brace after "for"
		brace := pos + 3
		for brace < len(src) && src[brace] != '{' {
			brace++
		}
		if brace >= len(src) {
			return
		}
		end := matchBrace(src, brace)
		if end < 0 {
			i = brace + 1
			continue
		}
		text := src[pos : end+1]
		if strings.Contains(text, "select") ||
			strings.Contains(text, "break") ||
			strings.Contains(text, "return") {
			i = end + 1
			continue
		}
		// bare receive: <- somewhere that looks like a receive (not send ch<-)
		if textContainsBareReceive(text) {
			pushAt(unit, meta, pos, msg, out)
		}
		i = end + 1
	}
}

func matchBrace(src string, open int) int {
	if open >= len(src) || src[open] != '{' {
		return -1
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for i := open; i < len(src); i++ {
		c := src[i]
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
		case '"', '\'', '`':
			inStr = c
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func textContainsBareReceive(text string) bool {
	// Look for "<-" not preceded by an identifier character (send is name<-).
	for i := 0; i+1 < len(text); i++ {
		if text[i] != '<' || text[i+1] != '-' {
			continue
		}
		// skip if this is a send: ident immediately before
		j := i - 1
		for j >= 0 && (text[j] == ' ' || text[j] == '\t') {
			j--
		}
		if j >= 0 {
			c := text[j]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == ')' {
				// likely send or cast; still could be receive in paren — treat as send-ish skip
				continue
			}
		}
		return true
	}
	return false
}

func detectBP92(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-92")
	if strings.Contains(unit.Source, "errgroup.Group") || strings.Contains(unit.Source, "new(errgroup.Group)") || strings.Contains(unit.Source, "errgroup.Group{}") {
		if !strings.Contains(unit.Source, "errgroup.WithContext") {
			if pos := strings.Index(unit.Source, "errgroup"); pos >= 0 {
				pushAt(unit, meta, pos, "errgroup without context; prefer errgroup.WithContext", out)
			}
		}
	}
}

func detectBP93(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-93")
	if !strings.Contains(unit.Source, ".Go(func") {
		return
	}
	// g.Go(func() error { ... }) that doesn't return err
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, ".Go(func") && strings.Contains(t, "error") {
			// check body discards — hard; flag if `_ =` inside next lines
		}
	}
	if strings.Contains(unit.Source, ".Go(func") && strings.Contains(unit.Source, "return nil") {
		// if also has err ignored
		if strings.Contains(unit.Source, "_ =") && strings.Contains(unit.Source, ".Go(") {
			if pos := strings.Index(unit.Source, ".Go("); pos >= 0 {
				pushAt(unit, meta, pos, "errgroup closure may discard an error; return it from the Go func", out)
			}
		}
	}
}

func detectBP94(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-94")
	if !strings.Contains(unit.Source, "go ") || !strings.Contains(unit.Source, "map[") {
		return
	}
	// go func writing map without mutex
	if strings.Contains(unit.Source, "go ") && strings.Contains(unit.Source, "] =") &&
		!strings.Contains(unit.Source, "sync.Mutex") && !strings.Contains(unit.Source, "sync.RWMutex") &&
		!strings.Contains(unit.Source, "sync.Map") {
		if pos := strings.Index(unit.Source, "go "); pos >= 0 {
			pushAt(unit, meta, pos, "map write from goroutine without synchronization", out)
		}
	}
}

func detectBP100(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-100")
	if !strings.Contains(unit.Source, "go ") {
		return
	}
	msg := "unbounded goroutine fan-out; bound concurrency with a worker pool or semaphore"
	// for range + go inside loop
	if len(facts.goNodes) > 0 && len(facts.forRanges) > 0 {
		for _, g := range facts.goNodes {
			if facts.insideLoop(g.start) {
				// skip if semaphore/errgroup/WaitGroup.Add with limited workers
				if strings.Contains(unit.Source, "make(chan") && strings.Contains(unit.Source, "cap(") {
					continue
				}
				pushAt(unit, meta, g.start, msg, out)
			}
		}
		return
	}
	// text: for _, x := range ... { go
	lines := codeLines(unit.Source)
	inFor := false
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "for ") && strings.Contains(t, " range ") {
			inFor = true
		}
		if inFor && strings.HasPrefix(t, "go ") {
			pushAt(unit, meta, line.byte, msg, out)
			inFor = false
		}
		if t == "}" {
			inFor = false
		}
	}
}

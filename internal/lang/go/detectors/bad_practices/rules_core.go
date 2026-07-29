package badpractices

import (
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

// BP-9: select without default/timeout escape.
func detectBP9(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-9")
	if !strings.Contains(unit.Source, "select {") {
		return
	}
	// Flag select with single case and no default/timeout — heuristic.
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "select {")
		if idx < 0 {
			break
		}
		abs := start + idx
		// extract block
		end := abs + len("select {")
		depth := 1
		for end < len(src) && depth > 0 {
			if src[end] == '{' {
				depth++
			} else if src[end] == '}' {
				depth--
			}
			end++
		}
		block := src[abs:end]
		if !strings.Contains(block, "default:") &&
			!strings.Contains(block, "time.After") &&
			!strings.Contains(block, "ctx.Done()") &&
			!strings.Contains(block, ".Done()") {
			// only flag if looks blocking with one case
			caseCount := strings.Count(block, "case ")
			if caseCount == 1 {
				pushAt(unit, meta, abs, "blocking select lacks a timeout or default escape path", out)
			}
		}
		start = end
	}
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

func detectBP83(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-83")
	if isTestFile(unit) {
		return
	}
	if !strings.Contains(unit.Source, "time.Sleep(") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		if strings.Contains(line.text, "time.Sleep(") {
			pushAt(unit, meta, line.byte, "time.Sleep used for synchronization; prefer channels, waitgroups, or condition variables", out)
		}
	}
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

func detectBP90(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-90")
	// for range ch without ok or default exit
	if strings.Contains(unit.Source, "for range ") || strings.Contains(unit.Source, "for {") {
		if strings.Contains(unit.Source, "<-") && !strings.Contains(unit.Source, ", ok :=") && !strings.Contains(unit.Source, ", ok:=") {
			// weak signal
			if strings.Contains(unit.Source, "for {") && strings.Contains(unit.Source, "<-ch") {
				if pos := strings.Index(unit.Source, "for {"); pos >= 0 {
					pushAt(unit, meta, pos, "channel receive loop without a clear exit condition", out)
				}
			}
		}
	}
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

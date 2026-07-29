package badpractices

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("BP-46", detectBP46)
	RegisterRule("BP-47", detectBP47)
	RegisterRule("BP-48", detectBP48)
	RegisterRule("BP-49", detectBP49)
	RegisterRule("BP-50", detectBP50)
	RegisterRule("BP-51", detectBP51)
	RegisterRule("BP-52", detectBP52)
	RegisterRule("BP-53", detectBP53)
	RegisterRule("BP-54", detectBP54)
	RegisterRule("BP-55", detectBP55)
	RegisterRule("BP-56", detectBP56)
}

func detectBP46(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-46")
	if isTestFile(unit) || !strings.Contains(unit.Source, "http.Server") {
		return
	}
	msg := "http.Server should set both ReadTimeout and WriteTimeout"
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "http.Server{")
		if idx < 0 {
			// also &http.Server{
			idx = strings.Index(src[start:], "http.Server{")
			if idx < 0 {
				// try without type name after composite
				break
			}
		}
		abs := start + idx
		end := abs
		depth := 0
		for i := abs; i < len(src); i++ {
			if src[i] == '{' {
				depth++
			} else if src[i] == '}' {
				depth--
				if depth == 0 {
					end = i + 1
					break
				}
			}
		}
		if end > abs {
			lit := src[abs:end]
			if !strings.Contains(lit, "ReadTimeout:") || !strings.Contains(lit, "WriteTimeout:") {
				pushAt(unit, meta, abs, msg, out)
			}
		}
		start = abs + len("http.Server{")
	}
	// also handle &http.Server{ pattern already covered
	// and multi-line with http.Server{\n
	if !strings.Contains(src, "http.Server{") && strings.Contains(src, "http.Server") {
		// var s http.Server; s.Addr= — skip
	}
}

func detectBP47(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-47")
	if isMaterializedFixture(unit) {
		return
	}
	snap := projectSnapshot(unit)
	if !snap.HasServerStart || snap.HasShutdown {
		return
	}
	pushAt(unit, meta, 0, "server startup should include a graceful shutdown path", out)
}

func detectBP48(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-48")
	if isTestFile(unit) {
		return
	}
	msg := "library code should return errors instead of exiting the process"
	callees := []string{"log.Fatal", "log.Fatalf", "log.Fatalln", "os.Exit"}
	if len(facts.callNodes) > 0 {
		for _, c := range facts.callNodes {
			match := false
			for _, name := range callees {
				if c.callee == name {
					match = true
					break
				}
			}
			if !match {
				continue
			}
			fn := facts.enclosingFunc(c.start)
			if fn != nil && (fn.isMain || fn.name == "TestMain") {
				continue
			}
			if name, ok := enclosingFuncName(unit.Source, c.start); ok && (name == "main" || name == "TestMain") {
				continue
			}
			pushAt(unit, meta, c.start, msg, out)
		}
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		for _, name := range callees {
			if strings.Contains(t, name+"(") {
				if n, ok := enclosingFuncName(unit.Source, line.byte); ok && (n == "main" || n == "TestMain") {
					continue
				}
				pushAt(unit, meta, line.byte, msg, out)
			}
		}
	}
}

func detectBP49(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-49")
	if isTestFile(unit) {
		return
	}
	msg := "deferred cleanup drops an error; wrap it in a deferred function and check the result"
	for _, d := range facts.deferNodes {
		text := d.text
		if strings.Contains(text, "func()") {
			continue
		}
		if strings.Contains(text, ".Close()") || strings.Contains(text, ".Flush()") || strings.Contains(text, ".Sync()") {
			pushAt(unit, meta, d.start, msg, out)
		}
	}
	if len(facts.deferNodes) == 0 {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if strings.HasPrefix(t, "defer ") && !strings.Contains(t, "func()") {
				if strings.Contains(t, ".Close()") || strings.Contains(t, ".Flush()") || strings.Contains(t, ".Sync()") {
					pushAt(unit, meta, line.byte, msg, out)
				}
			}
		}
	}
}

func detectBP50(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-50")
	if isMaterializedFixture(unit) {
		return
	}
	snap := projectSnapshot(unit)
	if !snap.HasServerStart || snap.HasSignalHandling {
		return
	}
	pushAt(unit, meta, 0, "long-running server should handle SIGTERM or SIGINT", out)
}

func detectBP51(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-51")
	if isTestFile(unit) || packageName(unit.Source) == "main" {
		return
	}
	if !strings.Contains(unit.Source, "recover()") {
		return
	}
	if strings.Contains(unit.Source, "panic(") || strings.Contains(unit.Source, "log.") || strings.Contains(unit.Source, "fmt.") {
		return
	}
	if pos := strings.Index(unit.Source, "recover()"); pos >= 0 {
		pushAt(unit, meta, pos, "library recover paths should re-panic or convert the panic into an explicit error contract", out)
	}
}

func detectBP52(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-52")
	if isTestFile(unit) || !strings.Contains(unit.Source, "make(") || !strings.Contains(unit.Source, "*") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "make(") && strings.Contains(t, "*") {
			if !strings.Contains(unit.Source, "MaxInt") && !strings.Contains(unit.Source, "overflow") && !strings.Contains(unit.Source, "bits.Mul") {
				pushAt(unit, meta, line.byte, "multiplication used in an allocation path without an obvious overflow guard", out)
				return
			}
		}
	}
}

func detectBP53(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-53")
	if isTestFile(unit) || !strings.Contains(unit.Source, "gob.Register(") {
		return
	}
	if !strings.Contains(unit.Source, ".Encode(") && !strings.Contains(unit.Source, ".Decode(") {
		if pos := strings.Index(unit.Source, "gob.Register("); pos >= 0 {
			pushAt(unit, meta, pos, "gob.Register without nearby Encode/Decode; registration may not match payloads", out)
		}
	}
}

func detectBP54(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-54")
	if isMaterializedFixture(unit) {
		return
	}
	snap := projectSnapshot(unit)
	if !snap.HasServerStart || !snap.HasPublicRoute || snap.HasRateLimiting {
		return
	}
	pushAt(unit, meta, 0, "public HTTP handlers should enforce a rate-limiting guard", out)
}

func detectBP55(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-55")
	if isMaterializedFixture(unit) {
		return
	}
	snap := projectSnapshot(unit)
	if !snap.HasServerStart || !snap.HasPublicRoute || !snap.HasLogging || snap.HasRequestID {
		return
	}
	pushAt(unit, meta, 0, "request-handling code logs traffic without a visible request-id propagation path", out)
}

func detectBP56(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-56")
	deprecated := []string{
		`"io/ioutil"`, `"golang.org/x/net/context"`, `"github.com/golang/protobuf"`,
	}
	for _, d := range deprecated {
		if strings.Contains(unit.Source, d) {
			if pos := strings.Index(unit.Source, d); pos >= 0 {
				pushAt(unit, meta, pos, "deprecated package import; migrate to the supported replacement", out)
			}
		}
	}
}

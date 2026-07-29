package perf

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

// Batch 5: PERF-215..242 excluding PERF-230 (seeded in seed_register.go).
func init() {
	RegisterRule("PERF-215", detectPERF215, &MetaPERF215)
	RegisterRule("PERF-216", detectPERF216, &MetaPERF216)
	RegisterRule("PERF-217", detectPERF217, &MetaPERF217)
	RegisterRule("PERF-218", detectPERF218, &MetaPERF218)
	RegisterRule("PERF-219", detectPERF219, &MetaPERF219)
	RegisterRule("PERF-220", detectPERF220, &MetaPERF220)
	RegisterRule("PERF-221", detectPERF221, &MetaPERF221)
	RegisterRule("PERF-222", detectPERF222, &MetaPERF222)
	RegisterRule("PERF-223", detectPERF223, &MetaPERF223)
	RegisterRule("PERF-224", detectPERF224, &MetaPERF224)
	RegisterRule("PERF-225", detectPERF225, &MetaPERF225)
	RegisterRule("PERF-226", detectPERF226, &MetaPERF226)
	RegisterRule("PERF-227", detectPERF227, &MetaPERF227)
	RegisterRule("PERF-228", detectPERF228, &MetaPERF228)
	RegisterRule("PERF-229", detectPERF229, &MetaPERF229)
	RegisterRule("PERF-231", detectPERF231, &MetaPERF231)
	RegisterRule("PERF-232", detectPERF232, &MetaPERF232)
	RegisterRule("PERF-233", detectPERF233, &MetaPERF233)
	RegisterRule("PERF-234", detectPERF234, &MetaPERF234)
	RegisterRule("PERF-235", detectPERF235, &MetaPERF235)
	RegisterRule("PERF-236", detectPERF236, &MetaPERF236)
	RegisterRule("PERF-237", detectPERF237, &MetaPERF237)
	RegisterRule("PERF-238", detectPERF238, &MetaPERF238)
	RegisterRule("PERF-239", detectPERF239, &MetaPERF239)
	RegisterRule("PERF-240", detectPERF240, &MetaPERF240)
	RegisterRule("PERF-241", detectPERF241, &MetaPERF241)
	RegisterRule("PERF-242", detectPERF242, &MetaPERF242)
}

// --- PERF-215: bytes.Buffer / strings.Builder without Grow ---

func detectPERF215(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source

	for _, name := range b5_bufferOrBuilderNames(source) {
		if writeByte, ok := b5_firstWriteWithKnownLen(source, name); ok {
			body := b5_enclosingFunctionBody(source, writeByte)
			if body == "" {
				body = source
			}
			if !strings.Contains(body, name+".Grow(") {
				line, col := unit.LineCol(writeByte)
				rules.PushFinding(
					&MetaPERF215, file, line, col,
					"bytes.Buffer or strings.Builder writes without a preceding Grow(expectedSize)", out,
				)
				return
			}
		}
		writeNeedle := name + ".Write"
		searchFrom := 0
		for {
			rel := strings.Index(source[searchFrom:], writeNeedle)
			if rel < 0 {
				break
			}
			bytePos := searchFrom + rel
			body := b5_enclosingFunctionBody(source, bytePos)
			if body == "" {
				searchFrom = bytePos + len(writeNeedle)
				continue
			}
			if strings.Contains(body, name+".Grow(") {
				searchFrom = bytePos + len(writeNeedle)
				continue
			}
			writeCount := 0
			for _, m := range []string{"WriteString(", "Write(", "WriteByte(", "WriteRune("} {
				writeCount += strings.Count(body, name+"."+m)
			}
			sizeHint := strings.Contains(body, "len(") || strings.Contains(body, ".Len()") || strings.Contains(body, "cap(")
			if writeCount >= 3 && (sizeHint || writeCount >= 6) && b5_enclosingFunctionIsHot(source, bytePos) {
				line, col := unit.LineCol(bytePos)
				rules.PushFinding(
					&MetaPERF215, file, line, col,
					"bytes.Buffer or strings.Builder does many writes without Grow; pre-size when output size is estimable", out,
				)
				return
			}
			searchFrom = bytePos + len(writeNeedle)
		}
	}
}

// --- PERF-216: TreeNode{} alloc in loop ---

func detectPERF216(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if strings.Contains(source, "sync.Pool") {
		return
	}
	for _, a := range facts.Assignments {
		if a.EnclosingLoop == nil {
			continue
		}
		expr := a.Expr
		if !(strings.Contains(expr, "TreeNode{") || strings.Contains(expr, "&TreeNode{")) {
			continue
		}
		line, col := unit.LineCol(a.StartByte)
		rules.PushFinding(
			&MetaPERF216, file, line, col,
			"struct literal allocation inside a loop on the hot path; reuse pooled objects instead", out,
		)
		return
	}
}

// --- PERF-217: static builder rebuilt on hot path ---

func detectPERF217(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, call := range facts.Calls {
		callee := call.Callee
		if strings.HasPrefix(callee, "func") || strings.Contains(callee, "{") || strings.Contains(callee, "\n") {
			continue
		}
		if b5_isPoolOrResetAccessor(callee) {
			continue
		}
		var nonEmpty []string
		for _, a := range call.Arguments {
			a = strings.TrimSpace(a)
			if a != "" {
				nonEmpty = append(nonEmpty, a)
			}
		}
		argsOK := len(nonEmpty) == 0
		if !argsOK {
			argsOK = true
			for _, a := range nonEmpty {
				if !b5_isConstLikeArg(a) {
					argsOK = false
					break
				}
			}
		}
		if !argsOK || !b5_looksLikeStaticBuilder(callee) {
			continue
		}
		if !IsHotPath(source, call.StartByte, IsInLoop(call)) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF217, file, line, col,
			"deterministic static computation is rebuilt on a hot path instead of cached at init", out,
		)
		return
	}
}

// --- PERF-218: unsharded zero-value sync.Pool under concurrency ---

func detectPERF218(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "sync.Pool") {
		return
	}
	if strings.Contains(source, "[runtime.NumCPU()]sync.Pool") ||
		strings.Contains(source, "runtime_procPin") ||
		strings.Contains(source, "shard") ||
		strings.Contains(source, "[]sync.Pool") {
		return
	}
	concurrent := b5_fileHasHandler(source) || len(facts.GoStarts) > 0 || b5_fileHasConcurrency(source)
	if !concurrent {
		return
	}
	for _, call := range facts.Calls {
		method, _ := splitCallee(call.Callee)
		if method != "Get" && method != "Put" {
			continue
		}
		recv := b5_receiverName(call.Callee)
		if recv == "" || !strings.Contains(source, "var "+recv+" sync.Pool") {
			continue
		}
		if !b5_hasZeroValuePoolDecl(source, recv) {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF218, file, line, col,
			"single global sync.Pool is used on a hot concurrent path without sharding", out,
		)
		return
	}
}

// --- PERF-219: Put without capacity guard ---

func detectPERF219(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "sync.Pool") {
		return
	}
	for _, call := range facts.Calls {
		bare, _ := splitCallee(call.Callee)
		if bare != "Put" {
			continue
		}
		if len(call.Arguments) == 0 {
			continue
		}
		arg := strings.TrimSpace(call.Arguments[0])
		if !b5_looksLikeBufferArg(arg) {
			continue
		}
		if !b5_enclosingFuncHasSliceBufParam(source, call.StartByte, arg) {
			continue
		}
		windowStart := b5_charBoundary(source, call.StartByte-200)
		window := source[windowStart:call.StartByte]
		if strings.Contains(window, "cap("+arg+") >") ||
			strings.Contains(window, "cap("+arg+") >=") ||
			strings.Contains(window, "cap("+arg+")>") {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF219, file, line, col,
			"object is returned to sync.Pool without a capacity guard, so oversized buffers stay retained", out,
		)
		return
	}
}

// --- PERF-220: consecutive loops over same "row" ---

func detectPERF220(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if len(facts.ForRanges) < 2 {
		return
	}
	loops := append([][2]int(nil), facts.ForRanges...)
	// sort by start
	for i := 0; i < len(loops); i++ {
		for j := i + 1; j < len(loops); j++ {
			if loops[j][0] < loops[i][0] {
				loops[i], loops[j] = loops[j], loops[i]
			}
		}
	}
	for i := 0; i+1 < len(loops); i++ {
		first, second := loops[i], loops[i+1]
		a := b5_loopTarget(source, first[0], first[1])
		b := b5_loopTarget(source, second[0], second[1])
		if a == "" || a != b || a != "row" {
			continue
		}
		if second[0]-first[1] > 16 {
			continue
		}
		line, col := unit.LineCol(second[0])
		rules.PushFinding(
			&MetaPERF220, file, line, col,
			"two consecutive loops scan the same collection; merge them into one pass", out,
		)
		return
	}
}

// --- PERF-221: map[int] sequential keys ---
//
// Index writes (`m[i+1] = v`) may be missing from AssignmentFact because
// extractIdents only keeps simple names. Scan source text as well.

func detectPERF221(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "map[int]") && !strings.Contains(source, "map[int64]") {
		return
	}
	for _, a := range facts.Assignments {
		if b5_sequentialMapWrite(a.Text) {
			line, col := unit.LineCol(a.StartByte)
			rules.PushFinding(
				&MetaPERF221, file, line, col,
				"map[int] is being filled with dense sequential keys; use a slice instead", out,
			)
			return
		}
	}
	// Source scan for name[seqKey] = form.
	search := 0
	for search < len(source) {
		bracket := strings.Index(source[search:], "[")
		if bracket < 0 {
			return
		}
		abs := search + bracket
		if abs == 0 || !b5_isIdentChar(rune(source[abs-1])) {
			search = abs + 1
			continue
		}
		end := abs + 48
		if end > len(source) {
			end = len(source)
		}
		// Include preceding ident so sequentialMapWrite sees name[
		start := abs - 1
		for start > 0 && b5_isIdentChar(rune(source[start-1])) {
			start--
		}
		window := source[start:end]
		if b5_sequentialMapWrite(window) {
			line, col := unit.LineCol(abs)
			rules.PushFinding(
				&MetaPERF221, file, line, col,
				"map[int] is being filled with dense sequential keys; use a slice instead", out,
			)
			return
		}
		search = abs + 1
	}
}

func b5_sequentialMapWrite(text string) bool {
	bracket := strings.Index(text, "[")
	if bracket < 0 {
		return false
	}
	if bracket > 0 {
		prev := text[bracket-1]
		if !b5_isIdentChar(rune(prev)) {
			return false
		}
	}
	key := strings.ToLower(text[bracket:])
	// Require assignment into the index (not a bare read / compare-only window).
	eq := strings.Index(key, "=")
	if eq < 0 {
		return false
	}
	// Skip == / !=
	if eq+1 < len(key) && key[eq+1] == '=' {
		return false
	}
	if eq > 0 && (key[eq-1] == '!' || key[eq-1] == '<' || key[eq-1] == '>') {
		return false
	}
	return strings.Contains(key, "[i]") ||
		strings.Contains(key, "[i+") ||
		strings.Contains(key, "[i +") ||
		strings.Contains(key, "[idx") ||
		strings.Contains(key, "[index") ||
		strings.Contains(key, "[len(") ||
		strings.Contains(key, "[j]") ||
		strings.Contains(key, "[n]")
}

// --- PERF-222: generic instantiation in loop ---

func detectPERF222(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	genericNames := b5_genericFuncNames(source)
	if len(genericNames) == 0 {
		return
	}
	for _, fr := range facts.ForRanges {
		end := fr[1] + 64
		if end > len(source) {
			end = len(source)
		}
		if fr[0] > end {
			continue
		}
		loopText := source[fr[0]:b5_charBoundary(source, end)]
		for _, name := range genericNames {
			needle := name + "["
			search := 0
			for {
				rel := strings.Index(loopText[search:], needle)
				if rel < 0 {
					break
				}
				at := search + rel
				after := loopText[at+len(needle):]
				if b5_isTypeInstantiationCall(after) {
					line, col := unit.LineCol(fr[0] + at)
					rules.PushFinding(
						&MetaPERF222, file, line, col,
						"generic function call appears on a measured hot path; prefer a concrete specialization", out,
					)
					return
				}
				search = at + len(needle)
			}
		}
	}
}

// --- PERF-223: x = nil before pool.Put(x) ---

func detectPERF223(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "sync.Pool") {
		return
	}
	for _, call := range facts.Calls {
		bare, _ := splitCallee(call.Callee)
		if bare != "Put" {
			continue
		}
		if len(call.Arguments) == 0 {
			continue
		}
		arg := strings.TrimSpace(call.Arguments[0])
		if arg == "" || !isSimpleIdent(arg) {
			continue
		}
		windowStart := b5_charBoundary(source, call.StartByte-160)
		window := source[windowStart:call.StartByte]
		if strings.Contains(window, arg+" = nil") ||
			strings.Contains(window, arg+"=nil") ||
			strings.Contains(window, arg+"= nil") {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(
				&MetaPERF223, file, line, col,
				"value is set to nil before pool return, so the backing array is discarded instead of reused", out,
			)
			return
		}
	}
}

// --- PERF-224: recursive walk on request path with flat representation ---

func detectPERF224(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !b5_fileHasHandler(source) && !IsRequestPath(facts.SourceIndex) {
		return
	}
	for _, fn := range b5_recursiveFunctions(source) {
		if !strings.Contains(source, fn+"(") {
			continue
		}
		if !strings.Contains(source, "flat") && !strings.Contains(source, "[]*Node") && !strings.Contains(source, "[]Node") {
			continue
		}
		bytePos, ok := b5_handlerCallSite(source, fn)
		if !ok {
			continue
		}
		line, col := unit.LineCol(bytePos)
		rules.PushFinding(
			&MetaPERF224, file, line, col,
			"recursive tree walk is invoked on the request path even though a flat representation already exists", out,
		)
		return
	}
}

// --- PERF-225: redundant large slice clone ---

func detectPERF225(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source

	type site struct {
		byte int
		arg  string
	}
	var sites []site
	for _, call := range facts.Calls {
		if call.Callee != "slices.Clone" {
			continue
		}
		arg := ""
		if len(call.Arguments) > 0 {
			arg = strings.TrimSpace(call.Arguments[0])
		}
		sites = append(sites, site{call.StartByte, arg})
	}
	for _, s := range b5_appendNilClones(source) {
		sites = append(sites, s)
	}
	if len(sites) < 2 {
		return
	}
	for i := 0; i < len(sites); i++ {
		for j := i + 1; j < len(sites); j++ {
			if sites[i].arg != "" && sites[i].arg == sites[j].arg {
				line, col := unit.LineCol(sites[j].byte)
				rules.PushFinding(
					&MetaPERF225, file, line, col,
					"large slice is fully cloned more than once; keep a single owned buffer", out,
				)
				return
			}
		}
	}
	// sort by byte
	for i := 0; i < len(sites); i++ {
		for j := i + 1; j < len(sites); j++ {
			if sites[j].byte < sites[i].byte {
				sites[i], sites[j] = sites[j], sites[i]
			}
		}
	}
	first, second := sites[0].byte, sites[1].byte
	n1, ok1 := EnclosingFunctionName(source, first)
	n2, ok2 := EnclosingFunctionName(source, second)
	if ok1 && ok2 && n1 == n2 {
		line, col := unit.LineCol(second)
		rules.PushFinding(
			&MetaPERF225, file, line, col,
			"large slice is fully cloned more than once; keep a single owned buffer", out,
		)
	}
}

// --- PERF-226: post-producer make+copy ---

func detectPERF226(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	emitted := 0

	search := 0
	for {
		rel := strings.Index(source[search:], "make([]byte")
		if rel < 0 {
			break
		}
		start := search + rel
		windowEnd := b5_charBoundary(source, min(start+200, len(source)))
		window := source[start:windowEnd]
		makeLen := strings.Contains(window, ".Len()") || strings.Contains(window, "len(")
		hasCopy := strings.Contains(window, "copy(")
		fromBytes := strings.Contains(window, ".Bytes()") || strings.Contains(window, "copy(")
		if makeLen && hasCopy && fromBytes {
			line, col := unit.LineCol(start)
			rules.PushFinding(
				&MetaPERF226, file, line, col,
				"buffer is re-copied after production (make+copy); take ownership of the producer buffer", out,
			)
			emitted++
			if emitted >= 8 {
				return
			}
		}
		search = start + 4
	}

	for _, needle := range []string{".Bytes()", ".Close()"} {
		search = 0
		for {
			rel := strings.Index(source[search:], needle)
			if rel < 0 {
				break
			}
			prod := search + rel
			windowEnd := b5_charBoundary(source, min(prod+480, len(source)))
			window := source[prod:windowEnd]
			if b5_windowHasRecopy(window) {
				recopyRel := strings.Index(window, "make([]byte")
				if recopyRel < 0 {
					recopyRel = strings.Index(window, "slices.Clone(")
				}
				if recopyRel < 0 {
					recopyRel = 0
				}
				abs := prod + recopyRel
				already := false
				absLine, _ := unit.LineCol(abs)
				for _, f := range *out {
					if f.RuleID == "PERF-226" && f.Line == absLine {
						already = true
						break
					}
				}
				if !already {
					line, col := unit.LineCol(abs)
					rules.PushFinding(
						&MetaPERF226, file, line, col,
						"buffer is re-copied immediately after production; take ownership instead of make+copy", out,
					)
					emitted++
					if emitted >= 8 {
						return
					}
				}
			}
			search = prod + len(needle)
		}
	}
}

// --- PERF-227: compress NewWriter without pool ---

func detectPERF227(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	triggers := map[string]bool{
		"flate.NewWriter": true, "flate.NewWriterDict": true,
		"zlib.NewWriter": true, "zlib.NewWriterLevel": true,
		"gzip.NewWriter": true, "gzip.NewWriterLevel": true,
	}
	for _, call := range facts.Calls {
		if !triggers[call.Callee] {
			continue
		}
		fname, _ := EnclosingFunctionName(source, call.StartByte)
		fnameL := strings.ToLower(fname)
		if strings.Contains(fnameL, "getzlib") || strings.Contains(fnameL, "getflate") ||
			strings.Contains(fnameL, "getgzip") || strings.Contains(fnameL, "newpool") ||
			(strings.HasPrefix(fnameL, "get") && strings.Contains(fnameL, "writer")) {
			continue
		}
		if b5_functionBodyHasWriterReset(source, call.StartByte) {
			continue
		}
		if !IsHotPath(source, call.StartByte, IsInLoop(call)) {
			if fname == "" || fname == "init" || fname == "main" {
				continue
			}
			if !b5_compressShapedFname(fnameL) {
				continue
			}
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF227, file, line, col,
			"compress writer is allocated without local pool/Reset reuse; pool writers on hot paths", out,
		)
	}
}

// --- PERF-228: parallel fan-out over tiny workset ---

func detectPERF228(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	tinyNames := b5_tinyCompositeSliceNames(source)
	for _, fr := range facts.ForRanges {
		end := fr[1]
		if end > len(source) {
			end = len(source)
		}
		if fr[0] >= end {
			continue
		}
		loopText := source[fr[0]:end]
		if !b5_loopHasParallelFanout(loopText) {
			continue
		}
		if n, ok := b5_compositeElemCountAfterRange(loopText); ok && n >= 1 && n <= 2 {
			line, col := unit.LineCol(fr[0])
			rules.PushFinding(
				&MetaPERF228, file, line, col,
				"parallel fan-out over a 1–2 element workset; prefer a serial path for tiny N", out,
			)
			return
		}
		if target, ok := b5_rangeTargetName(loopText); ok {
			for _, n := range tinyNames {
				if n == target {
					line, col := unit.LineCol(fr[0])
					rules.PushFinding(
						&MetaPERF228, file, line, col,
						"parallel fan-out over a 1–2 element workset; prefer a serial path for tiny N", out,
					)
					return
				}
			}
		}
	}
}

// --- PERF-229: Itoa/Sprintf then append string ---

func detectPERF229(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, a := range facts.Assignments {
		expr := a.Expr
		isFmt := strings.Contains(expr, "strconv.Itoa(") ||
			strings.Contains(expr, "strconv.FormatInt(") ||
			strings.Contains(expr, "strconv.FormatUint(") ||
			strings.Contains(expr, "fmt.Sprintf(")
		if !isFmt {
			continue
		}
		name := a.Name
		if name == "" || !isSimpleIdent(name) {
			continue
		}
		after := source[a.StartByte:]
		window := after
		if len(window) > 200 {
			window = window[:200]
		}
		if strings.Contains(window, name+"...") ||
			strings.Contains(window, "WriteString("+name+")") ||
			strings.Contains(window, "[]byte("+name+")") {
			line, col := unit.LineCol(a.StartByte)
			rules.PushFinding(
				&MetaPERF229, file, line, col,
				"temporary string is built then appended to bytes; use AppendInt/append-style APIs", out,
			)
			return
		}
	}
	if bytePos := strings.Index(source, "strconv.Itoa("); bytePos >= 0 {
		windowEnd := b5_charBoundary(source, min(bytePos+160, len(source)))
		window := source[bytePos:windowEnd]
		if strings.Contains(window, "append(") && strings.Contains(window, "...") {
			line, col := unit.LineCol(bytePos)
			rules.PushFinding(
				&MetaPERF229, file, line, col,
				"temporary string is built then appended to bytes; use AppendInt/append-style APIs", out,
			)
		}
	}
}

// --- PERF-231: PEM/key parse on hot path ---

func detectPERF231(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	triggers := map[string]bool{
		"pem.Decode":            true,
		"x509.ParseCertificate": true, "x509.ParsePKCS1PrivateKey": true,
		"x509.ParsePKCS8PrivateKey": true, "x509.ParseECPrivateKey": true,
		"tls.X509KeyPair": true, "tls.LoadX509KeyPair": true,
	}
	for _, call := range facts.Calls {
		if !triggers[call.Callee] {
			continue
		}
		hot := IsHotPath(source, call.StartByte, IsInLoop(call))
		if !hot && !IsRequestPath(facts.SourceIndex) {
			continue
		}
		if strings.Contains(source, "sync.Once") && strings.Contains(source, "Do(") {
			if !hot {
				continue
			}
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF231, file, line, col,
			"PEM/key material is parsed on a hot path; parse once at startup and reuse", out,
		)
		return
	}
}

// --- PERF-232: unbounded errgroup fan-out ---

func detectPERF232(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, fr := range facts.ForRanges {
		end := fr[1]
		if end > len(source) {
			end = len(source)
		}
		if fr[0] >= end {
			continue
		}
		loopText := source[fr[0]:end]
		if !strings.Contains(loopText, ".Go(") {
			continue
		}
		body := b5_enclosingFunctionBody(source, fr[0])
		if body == "" {
			body = source
		}
		if !b5_usesErrgroup(body) || b5_hasConcurrencyBound(body) {
			continue
		}
		line, col := unit.LineCol(fr[0])
		rules.PushFinding(
			&MetaPERF232, file, line, col,
			"parallel work fan-out has no SetLimit or semaphore bound; cap concurrency before spawning per-item work", out,
		)
		return
	}
}

// --- PERF-233: slow compress level on hot path ---

func detectPERF233(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	emitted := 0
	for _, call := range facts.Calls {
		argsJoined := strings.Join(call.Arguments, ",")
		if strings.Contains(argsJoined, "BestSpeed") ||
			strings.Contains(argsJoined, "HuffmanOnly") ||
			strings.Contains(argsJoined, "NoCompression") {
			continue
		}
		usesSlow := false
		switch call.Callee {
		case "zlib.NewWriter", "gzip.NewWriter":
			usesSlow = true
		case "zlib.NewWriterLevel", "flate.NewWriter", "flate.NewWriterLevel", "gzip.NewWriterLevel":
			usesSlow = strings.Contains(argsJoined, "DefaultCompression") ||
				strings.Contains(argsJoined, "BestCompression")
		}
		if !usesSlow {
			continue
		}
		if !IsHotPath(source, call.StartByte, IsInLoop(call)) {
			fname, _ := EnclosingFunctionName(source, call.StartByte)
			fnameL := strings.ToLower(fname)
			if fname == "" || fname == "init" || fname == "main" {
				continue
			}
			if !b5_compressShapedFname(fnameL) {
				continue
			}
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF233, file, line, col,
			"compress uses Default/BestCompression on a hot path; consider BestSpeed (or level 1) when size budget allows", out,
		)
		emitted++
		if emitted >= 8 {
			return
		}
	}
}

// --- PERF-234: bulk fixed Grow / pool Reset without Grow ---

const b5_bulkGrowMin uint64 = 4096

func detectPERF234(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	search := 0
	for {
		rel := strings.Index(source[search:], ".Grow(")
		if rel < 0 {
			break
		}
		start := search + rel
		after := strings.TrimLeft(source[start+len(".Grow("):], " \t")
		digits := ""
		for _, r := range after {
			if r >= '0' && r <= '9' {
				digits += string(r)
			} else {
				break
			}
		}
		if digits != "" {
			if n, err := strconv.ParseUint(digits, 10, 64); err == nil && n >= b5_bulkGrowMin {
				line, col := unit.LineCol(start)
				rules.PushFinding(
					&MetaPERF234, file, line, col,
					"bulk buffer uses a fixed Grow size; derive capacity from the input workload when it is known", out,
				)
				return
			}
		}
		search = start + 4
	}
	if strings.Contains(source, "Get().(*bytes.Buffer)") &&
		strings.Contains(source, ".Reset()") &&
		(strings.Contains(source, ".Write(") || strings.Contains(source, ".WriteString(") || strings.Contains(source, ".WriteByte(")) &&
		!strings.Contains(source, ".Grow(") {
		bytePos := strings.Index(source, "Get().(*bytes.Buffer)")
		if bytePos < 0 {
			bytePos = strings.Index(source, ".Reset()")
		}
		if bytePos < 0 {
			bytePos = 0
		}
		line, col := unit.LineCol(bytePos)
		rules.PushFinding(
			&MetaPERF234, file, line, col,
			"reused bulk buffer is reset without a workload-based Grow before assembly writes", out,
		)
	}
}

// --- PERF-235: strings.Builder .String() bridge ---

func detectPERF235(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "strings.Builder") {
		return
	}
	search := 0
	for {
		rel := strings.Index(source[search:], ".String()")
		if rel < 0 {
			break
		}
		start := search + rel
		body := b5_enclosingFunctionBody(source, start)
		if body == "" {
			body = source
		}
		if !strings.Contains(body, "strings.Builder") {
			search = start + 8
			continue
		}
		if b5_isStringBridgedIntoSink(source, start) {
			line, col := unit.LineCol(start)
			rules.PushFinding(
				&MetaPERF235, file, line, col,
				"temporary strings.Builder is flushed through .String() into a sink; write into the destination buffer directly", out,
			)
			return
		}
		search = start + 8
	}
}

// --- PERF-236: bytes.Clone on signing path ---

func detectPERF236(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	clone := strings.Index(source, "bytes.Clone(")
	if clone < 0 {
		return
	}
	function, _ := EnclosingFunctionName(source, clone)
	function = strings.ToLower(function)
	if !strings.Contains(function, "sign") && !strings.Contains(function, "signature") {
		return
	}
	line, col := unit.LineCol(clone)
	rules.PushFinding(
		&MetaPERF236, file, line, col,
		"signing path clones the complete buffer; prefer an owned writable buffer or in-place patching of reserved holes", out,
	)
}

// --- PERF-237: errgroup without tiny-N serial path ---

func detectPERF237(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, fr := range facts.ForRanges {
		end := fr[1]
		if end > len(source) {
			end = len(source)
		}
		if fr[0] >= end {
			continue
		}
		loopText := source[fr[0]:end]
		if !strings.Contains(loopText, ".Go(") {
			continue
		}
		body := b5_enclosingFunctionBody(source, fr[0])
		if body == "" {
			body = source
		}
		if !b5_usesErrgroup(body) {
			continue
		}
		if !strings.Contains(loopText, " range ") {
			continue
		}
		if b5_hasSerialShortCircuitNear(source, fr[0]) {
			continue
		}
		line, col := unit.LineCol(fr[0])
		rules.PushFinding(
			&MetaPERF237, file, line, col,
			"errgroup fan-out has no serial short-circuit for tiny worksets; run len(items) <= 2 serially before spawning", out,
		)
		return
	}
}

// --- PERF-238: map[rune]bool in loop ---

func detectPERF238(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "map[rune]bool") {
		return
	}
	for _, fr := range facts.ForRanges {
		end := fr[1]
		if end > len(source) {
			end = len(source)
		}
		if fr[0] >= end {
			continue
		}
		loopText := source[fr[0]:end]
		if !(strings.Contains(loopText, "] = true") ||
			strings.Contains(loopText, "]= true") ||
			strings.Contains(loopText, "]=true") ||
			strings.Contains(loopText, "] =true")) {
			continue
		}
		body := b5_enclosingFunctionBody(source, fr[0])
		if body == "" {
			body = source
		}
		if !strings.Contains(body, "map[rune]bool") && !strings.Contains(source, "map[rune]bool") {
			continue
		}
		line, col := unit.LineCol(fr[0])
		rules.PushFinding(
			&MetaPERF238, file, line, col,
			"rune membership is updated via map[rune]bool in a loop; prefer a bitset or denser set when the domain is bounded (e.g. BMP)", out,
		)
		return
	}
}

// --- PERF-239: dense map[int] write churn ---

func detectPERF239(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !strings.Contains(source, "map[int]") && !strings.Contains(source, "map[int64]") {
		return
	}
	search := 0
	for {
		rel := strings.Index(source[search:], "make(map[int")
		if rel < 0 {
			break
		}
		makeAt := search + rel
		head := source[:makeAt]
		name, ok := b5_mapNameBeforeMake(head)
		if !ok {
			search = makeAt + 4
			continue
		}
		body := b5_enclosingFunctionBody(source, makeAt)
		if body == "" {
			body = source
		}
		needle := name + "["
		assigns := strings.Count(body, needle)
		if assigns >= 6 {
			abs := makeAt
			if r := strings.Index(source[makeAt:], needle); r >= 0 {
				abs = makeAt + r
			}
			line, col := unit.LineCol(abs)
			rules.PushFinding(
				&MetaPERF239, file, line, col,
				"dense integer-keyed map is written many times in one function; prefer a slice or append-only records with one final index", out,
			)
			return
		}
		search = makeAt + 4
	}
}

// --- PERF-240: unpooled len-sized scratch ---

func detectPERF240(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	search := 0
	for {
		rel := strings.Index(source[search:], "make([]byte")
		if rel < 0 {
			break
		}
		start := search + rel
		afterEnd := b5_charBoundary(source, min(start+80, len(source)))
		after := source[start:afterEnd]
		if !strings.Contains(after, "len(") {
			search = start + 4
			continue
		}
		body := b5_enclosingFunctionBody(source, start)
		if body == "" {
			body = source
		}
		if strings.Contains(body, "sync.Pool") || strings.Contains(body, "Pool.Get") ||
			strings.Contains(body, "scratchPool") || strings.Contains(body, "bufPool") {
			search = start + 4
			continue
		}
		inLoop := false
		for _, fr := range facts.ForRanges {
			if start >= fr[0] && start < fr[1] {
				inLoop = true
				break
			}
		}
		fname, _ := EnclosingFunctionName(source, start)
		fnameL := strings.ToLower(fname)
		hotName := strings.Contains(fnameL, "encode") || strings.Contains(fnameL, "build") ||
			strings.Contains(fnameL, "subset") || strings.Contains(fnameL, "process") ||
			strings.Contains(fnameL, "render") || strings.Contains(fnameL, "compress") ||
			strings.Contains(fnameL, "write") || strings.Contains(fnameL, "generate")
		if !inLoop && !hotName {
			search = start + 4
			continue
		}
		line, col := unit.LineCol(start)
		rules.PushFinding(
			&MetaPERF240, file, line, col,
			"large []byte scratch is allocated from len(source) without pool reuse; pool and reset a scratch buffer on hot paths", out,
		)
		return
	}
}

// --- PERF-241: asn1.Marshal + time.Now on sign path ---

func detectPERF241(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	_ = facts
	file := unitFile(unit)
	source := unit.Source
	if !(strings.Contains(source, "asn1.Marshal") && strings.Contains(source, "time.Now")) {
		return
	}
	search := 0
	for {
		rel := strings.Index(source[search:], "asn1.Marshal")
		if rel < 0 {
			break
		}
		start := search + rel
		fname, _ := EnclosingFunctionName(source, start)
		fnameL := strings.ToLower(fname)
		body := b5_enclosingFunctionBody(source, start)
		signish := strings.Contains(fnameL, "sign") || strings.Contains(fnameL, "signature") ||
			strings.Contains(fnameL, "pkcs") || strings.Contains(fnameL, "cms")
		if signish && strings.Contains(body, "time.Now") {
			line, col := unit.LineCol(start)
			rules.PushFinding(
				&MetaPERF241, file, line, col,
				"signing path re-marshals ASN.1 with a fresh time.Now; cache immutable DER and only re-marshal time-varying attributes", out,
			)
			return
		}
		search = start + 8
	}
}

// --- PERF-242: per-iteration make([]byte, len*N) ---

func detectPERF242(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	source := unit.Source
	for _, fr := range facts.ForRanges {
		end := fr[1]
		if end > len(source) {
			end = len(source)
		}
		if fr[0] >= end {
			continue
		}
		loopText := source[fr[0]:end]
		if !strings.Contains(loopText, "make([]byte") {
			continue
		}
		if !(strings.Contains(loopText, "len(") && strings.Contains(loopText, "*")) {
			continue
		}
		rel := strings.Index(loopText, "make([]byte")
		if rel < 0 {
			continue
		}
		abs := fr[0] + rel
		line, col := unit.LineCol(abs)
		rules.PushFinding(
			&MetaPERF242, file, line, col,
			"loop allocates make([]byte, … len(x)*N …) each iteration; reuse a scratch buffer with [:0] growth", out,
		)
		return
	}
}

// ---------------------------------------------------------------------------
// Batch-5 helpers (ported from Rust stdlib_misuse helpers; scoped to this file).
// ---------------------------------------------------------------------------

func b5_bufferOrBuilderNames(source string) []string {
	var names []string
	for _, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(trimmed, "var "); ok {
			rest = strings.TrimLeft(rest, " \t")
			nameEnd := strings.IndexFunc(rest, unicode.IsSpace)
			if nameEnd < 0 {
				continue
			}
			name := strings.TrimSpace(rest[:nameEnd])
			ty := strings.TrimLeft(rest[nameEnd:], " \t")
			if (strings.HasPrefix(ty, "bytes.Buffer") || strings.HasPrefix(ty, "strings.Builder")) && isSimpleIdent(name) {
				names = append(names, name)
			}
			continue
		}
		if eq := strings.Index(trimmed, ":="); eq >= 0 {
			name := strings.TrimSpace(trimmed[:eq])
			rhs := strings.TrimSpace(trimmed[eq+2:])
			if (strings.HasPrefix(rhs, "bytes.Buffer{") || strings.HasPrefix(rhs, "strings.Builder{")) && isSimpleIdent(name) {
				names = append(names, name)
			}
		}
	}
	return names
}

func b5_firstWriteWithKnownLen(source, name string) (int, bool) {
	for _, method := range []string{"WriteString(", "Write("} {
		needle := name + "." + method
		searchFrom := 0
		for {
			rel := strings.Index(source[searchFrom:], needle)
			if rel < 0 {
				break
			}
			start := searchFrom + rel
			argStart := start + len(needle)
			rest := source[argStart:]
			argEnd := strings.IndexAny(rest, "),")
			if argEnd < 0 {
				searchFrom = argStart
				continue
			}
			arg := strings.TrimSpace(rest[:argEnd])
			if isSimpleIdent(arg) {
				body := b5_enclosingFunctionBody(source, start)
				if body == "" {
					body = source
				}
				if strings.Contains(body, "len("+arg+")") {
					return start, true
				}
			}
			searchFrom = argStart
		}
	}
	return 0, false
}

func b5_looksLikeBufferArg(arg string) bool {
	if !isSimpleIdent(arg) {
		return false
	}
	lower := strings.ToLower(arg)
	return strings.Contains(lower, "buf") ||
		strings.Contains(lower, "scratch") ||
		strings.Contains(lower, "tmp") ||
		strings.HasSuffix(lower, "buffer")
}

func b5_hasZeroValuePoolDecl(source, name string) bool {
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "var "+name+" sync.Pool") && !strings.Contains(t, "=") {
			return true
		}
	}
	return false
}

func b5_enclosingFuncHasSliceBufParam(source string, startByte int, name string) bool {
	if startByte > len(source) {
		startByte = len(source)
	}
	head := source[:startByte]
	funcKw := strings.LastIndex(head, "func ")
	if funcKw < 0 {
		return false
	}
	sig := source[funcKw:startByte]
	brace := strings.Index(sig, "{")
	if brace < 0 {
		return false
	}
	sig = sig[:brace]
	return strings.Contains(sig, name+" []byte") ||
		strings.Contains(sig, name+" []") ||
		strings.Contains(sig, name+" *bytes.Buffer")
}

func b5_receiverName(callee string) string {
	if i := strings.Index(callee, "."); i >= 0 {
		return callee[:i]
	}
	return ""
}

func b5_looksLikeStaticBuilder(callee string) bool {
	bare, _ := splitCallee(callee)
	lower := strings.ToLower(bare)
	return strings.Contains(lower, "build") ||
		strings.Contains(lower, "profile") ||
		strings.Contains(lower, "template") ||
		strings.Contains(lower, "serialize") ||
		strings.Contains(lower, "generate") ||
		strings.Contains(lower, "metadata") ||
		strings.Contains(lower, "defaultconfig") ||
		strings.Contains(lower, "loadconfig") ||
		(strings.Contains(lower, "compress") &&
			(strings.Contains(lower, "static") || strings.Contains(lower, "once") || strings.Contains(lower, "profile")))
}

func b5_isPoolOrResetAccessor(callee string) bool {
	bare, _ := splitCallee(callee)
	bare = strings.ToLower(bare)
	if strings.HasPrefix(bare, "reset") || strings.HasPrefix(bare, "clear") || strings.HasPrefix(bare, "recycle") {
		return true
	}
	switch bare {
	case "get", "put", "load", "store", "delete", "pop", "push", "take", "borrow":
		return true
	}
	getOrPut := strings.HasPrefix(bare, "get") || strings.HasPrefix(bare, "put")
	if getOrPut && (strings.Contains(bare, "buffer") || strings.Contains(bare, "writer") ||
		strings.Contains(bare, "pool") || strings.HasSuffix(bare, "buf") || strings.Contains(bare, "scratch")) {
		return true
	}
	recv := strings.ToLower(b5_receiverName(callee))
	if strings.Contains(recv, "pool") {
		switch bare {
		case "get", "put", "new":
			return true
		}
	}
	return false
}

func b5_isConstLikeArg(a string) bool {
	a = strings.TrimSpace(a)
	if a == "true" || a == "false" || a == "nil" {
		return true
	}
	if len(a) >= 2 && a[0] == '"' && a[len(a)-1] == '"' {
		return true
	}
	for _, r := range a {
		if r < '0' || r > '9' {
			return false
		}
	}
	return a != ""
}

func b5_genericFuncNames(source string) []string {
	var names []string
	search := 0
	for {
		rel := strings.Index(source[search:], "func ")
		if rel < 0 {
			break
		}
		at := search + rel
		after := strings.TrimLeft(source[at+len("func "):], " \t")
		if strings.HasPrefix(after, "(") {
			close := strings.Index(after, ")")
			if close < 0 {
				search = at + 4
				continue
			}
			after = strings.TrimLeft(after[close+1:], " \t")
		}
		nameEnd := 0
		for nameEnd < len(after) {
			r := rune(after[nameEnd])
			if !b5_isIdentChar(r) {
				break
			}
			nameEnd++
		}
		if nameEnd == 0 {
			search = at + 4
			continue
		}
		name := after[:nameEnd]
		rest := strings.TrimLeft(after[nameEnd:], " \t")
		if strings.HasPrefix(rest, "[") {
			names = append(names, name)
		}
		search = at + 4
	}
	return names
}

func b5_isTypeInstantiationCall(afterBracket string) bool {
	s := strings.TrimLeft(afterBracket, " \t")
	if s == "" {
		return false
	}
	close := strings.Index(s, "]")
	if close < 0 {
		return false
	}
	inside := strings.TrimSpace(s[:close])
	if inside == "" {
		return false
	}
	first := rune(inside[0])
	if !(unicode.IsLetter(first) || first == '_') {
		return false
	}
	return strings.HasPrefix(strings.TrimLeft(s[close+1:], " \t"), "(")
}

func b5_loopTarget(source string, start, end int) string {
	end = b5_charBoundary(source, min(end+64, len(source)))
	if start >= end || start < 0 {
		return ""
	}
	text := source[start:end]
	rangeIdx := strings.Index(text, "range")
	if rangeIdx < 0 {
		return ""
	}
	rest := strings.TrimLeft(text[rangeIdx+len("range"):], " \t")
	// stop at { or newline
	stop := len(rest)
	for i, r := range rest {
		if r == '{' || r == '\n' {
			stop = i
			break
		}
	}
	return strings.TrimSpace(rest[:stop])
}

func b5_recursiveFunctions(source string) []string {
	var out []string
	searchFrom := 0
	for {
		rel := strings.Index(source[searchFrom:], "func ")
		if rel < 0 {
			break
		}
		start := searchFrom + rel + len("func ")
		after := source[start:]
		nameEnd := 0
		for nameEnd < len(after) {
			r := rune(after[nameEnd])
			if !b5_isIdentChar(r) {
				break
			}
			nameEnd++
		}
		if nameEnd == 0 {
			searchFrom = start
			continue
		}
		name := after[:nameEnd]
		bodyStartRel := strings.Index(source[start:], "{")
		if bodyStartRel < 0 {
			break
		}
		bodyStart := start + bodyStartRel
		bodyEnd, ok := b5_matchBrace(source, bodyStart)
		if !ok {
			break
		}
		body := source[bodyStart:bodyEnd]
		if strings.Contains(body, name+"(") {
			out = append(out, name)
		}
		searchFrom = bodyEnd
	}
	return out
}

func b5_handlerCallSite(source, name string) (int, bool) {
	searchFrom := 0
	needle := name + "("
	for {
		rel := strings.Index(source[searchFrom:], needle)
		if rel < 0 {
			return 0, false
		}
		bytePos := searchFrom + rel
		if IsHandlerShaped(source, bytePos) {
			return bytePos, true
		}
		searchFrom = bytePos + len(name)
	}
}

func b5_matchBrace(source string, open int) (int, bool) {
	depth := 0
	for i := open; i < len(source); i++ {
		switch source[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1, true
			}
		}
	}
	return 0, false
}

func b5_appendNilClones(source string) []struct {
	byte int
	arg  string
} {
	var out []struct {
		byte int
		arg  string
	}
	for _, pattern := range []string{"append([]byte(nil), ", "append([]byte{}, ", "append(nil, "} {
		search := 0
		for {
			rel := strings.Index(source[search:], pattern)
			if rel < 0 {
				break
			}
			start := search + rel
			argStart := start + len(pattern)
			rest := source[argStart:]
			if end := strings.Index(rest, "..."); end >= 0 {
				arg := strings.TrimSpace(rest[:end])
				if isSimpleIdent(arg) {
					out = append(out, struct {
						byte int
						arg  string
					}{start, arg})
				}
			}
			search = argStart
		}
	}
	return out
}

func b5_windowHasRecopy(window string) bool {
	hasMake := strings.Contains(window, "make([]byte")
	hasCopy := strings.Contains(window, "copy(")
	hasClone := strings.Contains(window, "slices.Clone(")
	makeFromLen := strings.Contains(window, ".Len()") || strings.Contains(window, "len(")
	return (hasMake && hasCopy && makeFromLen) || hasClone
}

func b5_loopHasParallelFanout(loopText string) bool {
	return strings.Contains(loopText, ".Go(") ||
		strings.Contains(loopText, "g.Go(") ||
		strings.Contains(loopText, "group.Go(") ||
		strings.Contains(loopText, "go func") ||
		(strings.Contains(loopText, "wg.Add(") && strings.Contains(loopText, "go ")) ||
		(strings.Contains(loopText, "WaitGroup") && strings.Contains(loopText, "go "))
}

func b5_tinyCompositeSliceNames(source string) []string {
	var names []string
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		eq := strings.Index(t, ":=")
		if eq < 0 {
			eq = strings.Index(t, "=")
		}
		if eq < 0 {
			continue
		}
		lhs := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(t[:eq]), ":"))
		rhs := strings.TrimSpace(t[eq:])
		rhs = strings.TrimPrefix(rhs, ":")
		rhs = strings.TrimPrefix(rhs, "=")
		rhs = strings.TrimSpace(rhs)
		if !isSimpleIdent(lhs) {
			continue
		}
		if n, ok := b5_compositeLiteralElemCount(rhs); ok && n >= 1 && n <= 2 {
			names = append(names, lhs)
		}
	}
	return names
}

func b5_rangeTargetName(loopText string) (string, bool) {
	idx := strings.Index(loopText, "range")
	if idx < 0 {
		return "", false
	}
	rest := strings.TrimLeft(loopText[idx+len("range"):], " \t")
	// first token
	end := 0
	for end < len(rest) {
		r := rune(rest[end])
		if r == '{' || unicode.IsSpace(r) {
			break
		}
		end++
	}
	name := strings.TrimSpace(rest[:end])
	if isSimpleIdent(name) {
		return name, true
	}
	return "", false
}

func b5_compositeElemCountAfterRange(loopText string) (int, bool) {
	idx := strings.Index(loopText, "range")
	if idx < 0 {
		return 0, false
	}
	rest := strings.TrimLeft(loopText[idx+len("range"):], " \t")
	return b5_compositeLiteralElemCount(rest)
}

func b5_compositeLiteralElemCount(s string) (int, bool) {
	s = strings.TrimLeft(s, " \t")
	if !strings.HasPrefix(s, "[") {
		return 0, false
	}
	brace := strings.Index(s, "{")
	if brace < 0 {
		return 0, false
	}
	after := s[brace+1:]
	close := strings.Index(after, "}")
	if close < 0 {
		return 0, false
	}
	inner := strings.TrimSpace(after[:close])
	if inner == "" {
		return 0, true
	}
	n := 0
	for _, p := range strings.Split(inner, ",") {
		if strings.TrimSpace(p) != "" {
			n++
		}
	}
	return n, true
}

func b5_compressShapedFname(fnameL string) bool {
	return strings.Contains(fnameL, "compress") ||
		strings.Contains(fnameL, "encode") ||
		strings.Contains(fnameL, "write") ||
		strings.Contains(fnameL, "generate") ||
		strings.Contains(fnameL, "render") ||
		strings.Contains(fnameL, "export") ||
		strings.Contains(fnameL, "build") ||
		strings.Contains(fnameL, "stream") ||
		strings.Contains(fnameL, "serialize") ||
		strings.Contains(fnameL, "marshal")
}

func b5_functionBodyHasWriterReset(source string, startByte int) bool {
	if startByte > len(source) {
		startByte = len(source)
	}
	head := source[:startByte]
	funcKw := strings.LastIndex(head, "func ")
	if funcKw < 0 {
		return false
	}
	braceRel := strings.Index(source[funcKw:], "{")
	if braceRel < 0 {
		return false
	}
	bodyOpen := funcKw + braceRel
	end := bodyOpen
	depth := 0
	for i := bodyOpen; i < len(source); i++ {
		switch source[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i
				i = len(source)
			}
		}
	}
	body := source[bodyOpen:min(end, len(source))]
	return strings.Contains(body, ".Reset(") || strings.Contains(body, "Reset(")
}

func b5_usesErrgroup(body string) bool {
	return strings.Contains(body, "errgroup.Group") ||
		strings.Contains(body, "errgroup.WithContext") ||
		strings.Contains(body, "errgroup.")
}

func b5_hasConcurrencyBound(body string) bool {
	return strings.Contains(body, "SetLimit(") ||
		strings.Contains(body, "semaphore") ||
		strings.Contains(body, "sem.Acquire(") ||
		strings.Contains(body, "Acquire(ctx")
}

func b5_hasSerialShortCircuitNear(source string, loopStart int) bool {
	windowStart := loopStart - 600
	if windowStart < 0 {
		windowStart = 0
	}
	window := source[windowStart:loopStart]
	if !strings.Contains(window, "len(") {
		return false
	}
	for _, n := range []string{"<= 2", "<=2", "< 2", "<2", "< 3", "<3", "== 1", "==1", "<= 1", "<=1"} {
		if strings.Contains(window, n) {
			return true
		}
	}
	return false
}

func b5_isStringBridgedIntoSink(source string, stringDot int) bool {
	from := stringDot - 96
	if from < 0 {
		from = 0
	}
	before := source[from:stringDot]
	for _, needle := range []string{"WriteString(", "WriteString (", "[]byte(", "append(", "append ("} {
		pos := strings.LastIndex(before, needle)
		if pos < 0 {
			continue
		}
		between := strings.TrimSpace(before[pos+len(needle):])
		if isSimpleIdent(between) {
			return true
		}
		if strings.HasPrefix(needle, "append") {
			parts := strings.Split(between, ",")
			if len(parts) > 0 {
				last := strings.TrimSpace(parts[len(parts)-1])
				if isSimpleIdent(last) {
					return true
				}
			}
		}
	}
	return false
}

func b5_mapNameBeforeMake(head string) (string, bool) {
	trimmed := strings.TrimRight(head, " \t\n\r")
	assign := strings.LastIndex(trimmed, ":=")
	if assign < 0 {
		assign = strings.LastIndex(trimmed, "=")
	}
	if assign < 0 {
		return "", false
	}
	before := strings.TrimRight(trimmed[:assign], " \t")
	// last identifier
	i := len(before)
	for i > 0 {
		r := rune(before[i-1])
		if b5_isIdentChar(r) {
			i--
			continue
		}
		break
	}
	name := before[i:]
	if name == "" || !isSimpleIdent(name) {
		return "", false
	}
	return name, true
}

func b5_enclosingFunctionBody(source string, startByte int) string {
	if startByte > len(source) {
		startByte = len(source)
	}
	if startByte < 0 {
		return ""
	}
	head := source[:startByte]
	funcKw := strings.LastIndex(head, "func ")
	if funcKw < 0 {
		return ""
	}
	// body open after this func
	relHead := head[funcKw:]
	braceRel := strings.Index(relHead, "{")
	if braceRel < 0 {
		return ""
	}
	bodyOpen := funcKw + braceRel
	// ensure still inside
	depth := 0
	for i := bodyOpen; i < startByte; i++ {
		switch source[i] {
		case '{':
			depth++
		case '}':
			depth--
		}
	}
	if depth <= 0 {
		return ""
	}
	end := bodyOpen
	d := 0
	for i := bodyOpen; i < len(source); i++ {
		switch source[i] {
		case '{':
			d++
		case '}':
			d--
			if d == 0 {
				end = i + 1
				return source[bodyOpen:end]
			}
		}
	}
	return ""
}

func b5_enclosingFunctionIsHot(source string, startByte int) bool {
	name, ok := EnclosingFunctionName(source, startByte)
	return ok && FunctionNameIsHot(name)
}

func b5_fileHasHandler(source string) bool {
	return IsHandlerShaped(source, len(source))
}

func b5_fileHasConcurrency(source string) bool {
	return strings.Contains(source, "go ") ||
		strings.Contains(source, "go\t") ||
		strings.Contains(source, "go\n") ||
		strings.Contains(source, "errgroup") ||
		strings.Contains(source, "WaitGroup")
}

func b5_charBoundary(s string, index int) int {
	if index > len(s) {
		index = len(s)
	}
	if index < 0 {
		return 0
	}
	for index > 0 && index < len(s) && !utf8Start(s, index) {
		index--
	}
	return index
}

func b5_isIdentChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

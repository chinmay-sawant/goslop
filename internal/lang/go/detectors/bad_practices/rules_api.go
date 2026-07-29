package badpractices

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("BP-26", detectBP26)
	RegisterRule("BP-27", detectBP27)
	RegisterRule("BP-28", detectBP28)
	RegisterRule("BP-29", detectBP29)
	RegisterRule("BP-32", detectBP32)
	RegisterRule("BP-34", detectBP34)
	RegisterRule("BP-36", detectBP36)
	RegisterRule("BP-37", detectBP37)
	RegisterRule("BP-38", detectBP38)
	RegisterRule("BP-39", detectBP39)
	RegisterRule("BP-41", detectBP41)
	RegisterRule("BP-43", detectBP43)
	RegisterRule("BP-44", detectBP44)
	RegisterRule("BP-66", detectBP66)
	RegisterRule("BP-76", detectBP76)
	RegisterRule("BP-85", detectBP85)
	RegisterRule("BP-91", detectBP91)
	RegisterRule("BP-138", detectBP138)
	RegisterRule("BP-141", detectBP141)
	RegisterRule("BP-164", detectBP164)
}

func detectBP26(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-26")
	if !strings.Contains(unit.Source, "context.Context") {
		return
	}
	msg := "context.Context should be the first parameter"
	for _, fn := range facts.funcDecls {
		p := fn.params
		if !strings.Contains(p, "context.Context") {
			continue
		}
		// strip outer parens
		inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(p, "("), ")"))
		if inner == "" {
			continue
		}
		// first param should contain context.Context
		first := strings.Split(inner, ",")[0]
		if !strings.Contains(first, "context.Context") {
			pushAt(unit, meta, fn.start, msg, out)
		}
	}
	if len(facts.funcDecls) == 0 {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if !strings.HasPrefix(t, "func ") || !strings.Contains(t, "context.Context") {
				continue
			}
			// find param list
			open := strings.Index(t, "(")
			close := strings.LastIndex(t, ")")
			if open < 0 || close <= open {
				continue
			}
			// may be method: func (r *T) Name(
			rest := t[open+1:]
			if strings.HasPrefix(strings.TrimSpace(rest), "*") || strings.Contains(rest[:min(20, len(rest))], " ") && !strings.Contains(rest[:min(30, len(rest))], "context") {
				// method receiver — find second (
				if second := strings.Index(t[open+1:], "("); second >= 0 {
					open = open + 1 + second
					rest = t[open+1:]
					close = strings.Index(rest, ")")
					if close < 0 {
						continue
					}
					inner := rest[:close]
					first := strings.Split(inner, ",")[0]
					if strings.Contains(inner, "context.Context") && !strings.Contains(first, "context.Context") {
						pushAt(unit, meta, line.byte, msg, out)
					}
				}
				continue
			}
			close = strings.Index(rest, ")")
			if close < 0 {
				continue
			}
			inner := rest[:close]
			if !strings.Contains(inner, "context.Context") {
				continue
			}
			first := strings.Split(inner, ",")[0]
			if !strings.Contains(first, "context.Context") {
				pushAt(unit, meta, line.byte, msg, out)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func detectBP27(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-27")
	// exported func returns unexported type
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "func ") {
			continue
		}
		// func Name(...) unexported
		nameStart := len("func ")
		rest := t[nameStart:]
		if strings.HasPrefix(rest, "(") {
			// method
			continue
		}
		end := 0
		for end < len(rest) && (unicode.IsLetter(rune(rest[end])) || unicode.IsDigit(rune(rest[end])) || rest[end] == '_') {
			end++
		}
		if end == 0 {
			continue
		}
		name := rest[:end]
		if name == "" || !unicode.IsUpper(rune(name[0])) {
			continue
		}
		// look at return types after )
		if i := strings.LastIndex(t, ")"); i >= 0 {
			rets := strings.TrimSpace(t[i+1:])
			rets = strings.TrimPrefix(rets, "(")
			rets = strings.TrimSuffix(rets, "{")
			rets = strings.TrimSpace(strings.TrimSuffix(rets, ")"))
			for _, r := range strings.Split(rets, ",") {
				r = strings.TrimSpace(r)
				r = strings.TrimPrefix(r, "*")
				if r == "" || r == "error" || r == "string" || r == "int" || r == "bool" || r == "byte" {
					continue
				}
				if unicode.IsLetter(rune(r[0])) && unicode.IsLower(rune(r[0])) && !strings.Contains(r, ".") {
					pushAt(unit, meta, line.byte, "exported function returns an unexported type", out)
					break
				}
			}
		}
	}
}

func detectBP28(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-28")
	// single-method interface
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "interface {")
		if idx < 0 {
			idx = strings.Index(src[start:], "interface{")
			if idx < 0 {
				break
			}
		}
		abs := start + idx
		// count methods roughly
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
		block := src[abs:end]
		// method lines: name(
		methods := 0
		for _, line := range strings.Split(block, "\n") {
			t := strings.TrimSpace(line)
			if t == "" || strings.HasPrefix(t, "//") || t == "interface {" || t == "interface{" || t == "}" {
				continue
			}
			if strings.Contains(t, "(") {
				methods++
			}
		}
		if methods == 1 {
			pushAt(unit, meta, abs, "single-method interface; consider accepting a function type instead", out)
		}
		start = end
	}
}

func detectBP29(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-29")
	src := unit.Source
	start := 0
	for {
		idx := strings.Index(src[start:], "interface {")
		if idx < 0 {
			break
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
		block := src[abs:end]
		methods := 0
		for _, line := range strings.Split(block, "\n") {
			t := strings.TrimSpace(line)
			if strings.Contains(t, "(") && !strings.HasPrefix(t, "interface") {
				methods++
			}
		}
		if methods >= 8 {
			pushAt(unit, meta, abs, "interface has many methods; consider splitting responsibilities", out)
		}
		start = end
	}
}

func detectBP32(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-32")
	// type X string with Error() method — string alias error
	if strings.Contains(unit.Source, "type ") && strings.Contains(unit.Source, " string") && strings.Contains(unit.Source, "Error() string") {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if strings.HasPrefix(t, "type ") && strings.HasSuffix(t, " string") {
				pushAt(unit, meta, line.byte, "string alias used as error type; prefer a struct implementing error", out)
				return
			}
		}
	}
}

func detectBP34(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-34")
	// fmt.Errorf without %w when wrapping err
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "fmt.Errorf(") && strings.Contains(t, "err") && !strings.Contains(t, "%w") {
			pushAt(unit, meta, line.byte, "error wrapping without %w loses the error chain", out)
		}
	}
}

func detectBP36(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-36")
	if !strings.Contains(unit.Source, "func init(") {
		return
	}
	// init with network/db/side effects
	side := []string{"http.", "sql.Open", "Listen", "Dial", "os.Exit", "log.Fatal", "time.Sleep"}
	// extract init body roughly
	if pos := strings.Index(unit.Source, "func init("); pos >= 0 {
		body := unit.Source[pos:]
		if end := strings.Index(body, "\nfunc "); end > 0 {
			body = body[:end]
		}
		for _, s := range side {
			if strings.Contains(body, s) {
				pushAt(unit, meta, pos, "init has side effects; prefer explicit setup from main", out)
				return
			}
		}
	}
}

func detectBP37(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust parity: package-level var that is later written (not just declared).
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-37")
	globals := packageLevelVarNames(unit.Source)
	if len(globals) == 0 {
		return
	}
	written := packageLevelVarsWritten(unit.Source, globals)
	if len(written) == 0 {
		return
	}
	// Emit at first package-level var that is written.
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "var ") {
			continue
		}
		for name := range written {
			if strings.Contains(t, name) {
				// Skip pure error sentinels
				if strings.HasPrefix(name, "Err") {
					continue
				}
				pushAt(unit, meta, line.byte, "package-level mutable global state makes behavior harder to reason about", out)
				return
			}
		}
	}
}

func packageLevelVarNames(source string) map[string]struct{} {
	out := map[string]struct{}{}
	inFunc := 0
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		// crude brace depth for top-level only
		open := strings.Count(t, "{")
		closeN := strings.Count(t, "}")
		if inFunc == 0 && strings.HasPrefix(t, "var ") {
			// var name T or var ( ... )
			rest := strings.TrimSpace(strings.TrimPrefix(t, "var "))
			if strings.HasPrefix(rest, "(") {
				// multi-line block handled poorly — skip paren form for now
			} else {
				name := firstIdent(rest)
				if name != "" && !strings.HasPrefix(name, "Err") {
					out[name] = struct{}{}
				}
			}
		}
		inFunc += open - closeN
		if inFunc < 0 {
			inFunc = 0
		}
	}
	return out
}

func packageLevelVarsWritten(source string, globals map[string]struct{}) map[string]struct{} {
	written := map[string]struct{}{}
	// After first function starts, look for name = or name[ or name.Field =
	inFunc := false
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "func ") {
			inFunc = true
		}
		if !inFunc {
			continue
		}
		for name := range globals {
			// assignment to global (not short decl of local with same name)
			if strings.Contains(t, name+" =") || strings.Contains(t, name+"=") ||
				strings.Contains(t, name+"[") || strings.Contains(t, name+".") {
				// exclude := which is short decl (often shadows)
				if strings.Contains(t, name+" :=") || strings.Contains(t, name+":=") {
					continue
				}
				written[name] = struct{}{}
			}
		}
	}
	return written
}

func firstIdent(s string) string {
	s = strings.TrimSpace(s)
	i := 0
	for i < len(s) {
		r := rune(s[i])
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return ""
			}
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			break
		}
		i++
	}
	return s[:i]
}

func detectBP38(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Unexported helper with no same-file callers.
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-38")
	src := unit.Source
	helpers := unexportedHelpers(src)
	if len(helpers) == 0 {
		return
	}
	calls := localCallNames(src)
	for name, byteOff := range helpers {
		if _, used := calls[name]; used {
			continue
		}
		// skip common generated / interface method noise
		if name == "init" || name == "main" {
			continue
		}
		pushAt(unit, meta, byteOff, "unexported helper has no same-file callers", out)
	}
}

func unexportedHelpers(source string) map[string]int {
	// Package-scope funcs and methods (Rust collects both function_declaration
	// and method_declaration). Name filter: helper/must/build/parse*.
	out := map[string]int{}
	byteOff := 0
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "func ") {
			rest := strings.TrimPrefix(t, "func ")
			// methods: func (r *T) name(
			if strings.HasPrefix(rest, "(") {
				if close := strings.Index(rest, ")"); close >= 0 {
					rest = strings.TrimSpace(rest[close+1:])
				} else {
					byteOff += len(line) + 1
					continue
				}
			}
			name := firstIdent(rest)
			if name != "" && unicode.IsLower(rune(name[0])) && looksLikeHelperName(name) {
				if name != "init" && name != "main" &&
					!strings.HasPrefix(name, "Test") &&
					!strings.HasPrefix(name, "Benchmark") &&
					!strings.HasPrefix(name, "Example") {
					// first declaration wins for multi-line
					if _, ok := out[name]; !ok {
						out[name] = byteOff
					}
				}
			}
		}
		byteOff += len(line) + 1
	}
	return out
}

func looksLikeHelperName(name string) bool {
	// Rust parity: only helper/must/build/parse prefixes.
	lower := strings.ToLower(name)
	return lower == "helper" ||
		strings.HasPrefix(lower, "helper") ||
		strings.HasPrefix(lower, "must") ||
		strings.HasPrefix(lower, "build") ||
		strings.HasPrefix(lower, "parse")
}

func localCallNames(source string) map[string]struct{} {
	out := map[string]struct{}{}
	// crude: bare name( not preceded by ., excluding func/method declarations.
	for i := 0; i < len(source); i++ {
		if source[i] != '(' {
			continue
		}
		// walk back over ident
		j := i - 1
		for j >= 0 && isIdentByte(source[j]) {
			j--
		}
		name := source[j+1 : i]
		if name == "" {
			continue
		}
		// skip method/selector call: x.Name(
		if j >= 0 && source[j] == '.' {
			continue
		}
		// skip function/method declarations:
		//   func name(
		//   func (recv T) name(
		k := j
		for k >= 0 && (source[k] == ' ' || source[k] == '\t') {
			k--
		}
		if k >= 0 && source[k] == ')' {
			// method declaration receiver — skip back over "(...)" then look for func
			depth := 1
			k--
			for k >= 0 && depth > 0 {
				switch source[k] {
				case ')':
					depth++
				case '(':
					depth--
				}
				k--
			}
			for k >= 0 && (source[k] == ' ' || source[k] == '\t') {
				k--
			}
		}
		if k >= 3 && source[k-3:k+1] == "func" {
			if k-3 == 0 || !isIdentByte(source[k-4]) {
				continue // declaration, not a call
			}
		}
		if unicode.IsLower(rune(name[0])) {
			out[name] = struct{}{}
		}
	}
	return out
}

func detectBP39(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-39")
	lines := strings.Split(unit.Source, "\n")
	byteOff := 0
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "func ") {
			// get name
			rest := strings.TrimPrefix(t, "func ")
			if strings.HasPrefix(rest, "(") {
				byteOff += len(line) + 1
				continue // method; still check
			}
			nameEnd := 0
			for nameEnd < len(rest) && (unicode.IsLetter(rune(rest[nameEnd])) || unicode.IsDigit(rune(rest[nameEnd])) || rest[nameEnd] == '_') {
				nameEnd++
			}
			name := rest[:nameEnd]
			if name != "" && unicode.IsUpper(rune(name[0])) {
				// previous non-empty line should be doc comment
				hasDoc := false
				for j := i - 1; j >= 0; j-- {
					pt := strings.TrimSpace(lines[j])
					if pt == "" {
						continue
					}
					if strings.HasPrefix(pt, "//") || strings.HasPrefix(pt, "/*") {
						hasDoc = true
					}
					break
				}
				if !hasDoc {
					pushAt(unit, meta, byteOff, "exported function lacks a doc comment", out)
				}
			}
		}
		byteOff += len(line) + 1
	}
}

func detectBP41(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust parity: only the package anchor file reports missing package doc.
	if isTestFile(unit) || isMaterializedFixture(unit) {
		return
	}
	meta := MetadataForID("BP-41")
	pkg := packageName(unit.Source)
	if pkg == "" {
		return
	}
	snap := packageDocSnapshotForUnit(unit)
	anchor, ok := snap.Anchors[pkg]
	if !ok {
		return
	}
	// Compare absolute paths when possible.
	unitPath := unit.Path
	if abs, err := filepath.Abs(unitPath); err == nil {
		unitPath = abs
	}
	if abs, err := filepath.Abs(anchor); err == nil {
		anchor = abs
	}
	if unitPath != anchor {
		return
	}
	if _, documented := snap.DocumentedPackages[pkg]; documented {
		return
	}
	pushAt(unit, meta, 0, "package is missing a package-level doc comment", out)
}

func detectBP43(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-43")
	if isTestFile(unit) {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, ". \"") || strings.Contains(t, "\t. \"") || strings.Contains(t, " . \"") {
			pushAt(unit, meta, line.byte, "dot import outside tests pollutes the namespace", out)
		}
		if t == `. "` || strings.HasPrefix(t, `. "`) {
			pushAt(unit, meta, line.byte, "dot import outside tests pollutes the namespace", out)
		}
	}
	// import ( . "fmt" )
	if strings.Contains(unit.Source, `. "`) {
		if pos := strings.Index(unit.Source, `. "`); pos >= 0 {
			// avoid false positive in strings
			pushAt(unit, meta, pos, "dot import outside tests pollutes the namespace", out)
		}
	}
}

func detectBP44(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-44")
	// blank import without comment
	lines := strings.Split(unit.Source, "\n")
	byteOff := 0
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "_ \"") || strings.HasPrefix(t, "_ \"") {
			// comment on same line or previous?
			if !strings.Contains(line, "//") {
				prevOK := false
				if i > 0 && strings.Contains(lines[i-1], "//") {
					prevOK = true
				}
				if !prevOK {
					pushAt(unit, meta, byteOff, "blank import without a justifying comment", out)
				}
			}
		}
		byteOff += len(line) + 1
	}
}

func detectBP66(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-66")
	// err == SomeSentinel when wrapping used
	if strings.Contains(unit.Source, "fmt.Errorf") && strings.Contains(unit.Source, "%w") {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if (strings.Contains(t, "err ==") || strings.Contains(t, "err !=")) && strings.Contains(t, "Err") {
				if !strings.Contains(t, "errors.Is") {
					pushAt(unit, meta, line.byte, "wrapped error compared directly; use errors.Is", out)
				}
			}
		}
	}
}

func detectBP76(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-76")
	// for k, v := range map used to build ordered output
	if strings.Contains(unit.Source, "range ") && strings.Contains(unit.Source, "map[") {
		if strings.Contains(unit.Source, "sort.") || strings.Contains(unit.Source, "json.Marshal") || strings.Contains(unit.Source, "fmt.Print") {
			// weak
			if !strings.Contains(unit.Source, "sort.") {
				for _, line := range codeLines(unit.Source) {
					t := strings.TrimSpace(line.text)
					if strings.HasPrefix(t, "for ") && strings.Contains(t, " range ") {
						// check if ranging over map var
						pushAt(unit, meta, line.byte, "map range used where ordered output may be expected", out)
						return
					}
				}
			}
		}
	}
}

func detectBP85(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-85")
	// r.Context().Value without ok check - type assert
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, ".Value(") && strings.Contains(t, ".(") && !strings.Contains(t, "ok") {
			pushAt(unit, meta, line.byte, "unchecked type assertion on request context value", out)
		}
	}
}

func detectBP91(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-91")
	// chan struct{} for data-bearing? opposite: make(chan T) for notification
	// flag make(chan bool) used as signal
	if strings.Contains(unit.Source, "make(chan bool") || strings.Contains(unit.Source, "make(chan int") {
		if strings.Contains(unit.Source, "<-") && !strings.Contains(unit.Source, "make(chan struct{}") {
			if pos := strings.Index(unit.Source, "make(chan "); pos >= 0 {
				pushAt(unit, meta, pos, "notification channel carries data; prefer chan struct{}", out)
			}
		}
	}
}

func detectBP138(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-138")
	// GORM hook with http/external
	if strings.Contains(unit.Source, "BeforeCreate") || strings.Contains(unit.Source, "AfterCreate") ||
		strings.Contains(unit.Source, "BeforeSave") || strings.Contains(unit.Source, "AfterSave") {
		if strings.Contains(unit.Source, "http.") || strings.Contains(unit.Source, "smtp") || strings.Contains(unit.Source, "Dial") {
			pos := strings.Index(unit.Source, "Before")
			if pos < 0 {
				pos = strings.Index(unit.Source, "After")
			}
			if pos >= 0 {
				pushAt(unit, meta, pos, "external I/O inside a GORM hook; keep hooks local and fast", out)
			}
		}
	}
}

func detectBP141(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-141")
	if !strings.Contains(unit.Source, "Named") && !strings.Contains(unit.Source, "sqlx") {
		return
	}
	// :name placeholders without matching db tags — hard; skip unless obvious
	if strings.Contains(unit.Source, "NamedExec") || strings.Contains(unit.Source, "NamedQuery") {
		if !strings.Contains(unit.Source, "`db:") {
			if pos := strings.Index(unit.Source, "Named"); pos >= 0 {
				pushAt(unit, meta, pos, "sqlx named query may lack matching struct db tags", out)
			}
		}
	}
}

func detectBP164(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-164")
	// functional option mutates package default
	if strings.Contains(unit.Source, "func ") && strings.Contains(unit.Source, "Option") {
		// defaultOptions or Default
		if strings.Contains(unit.Source, "default") && strings.Contains(unit.Source, "func With") {
			for _, line := range codeLines(unit.Source) {
				t := strings.TrimSpace(line.text)
				if strings.HasPrefix(t, "func With") && strings.Contains(unit.Source, "default") {
					// if assigns to package var
					if strings.Contains(unit.Source, "default") && strings.Contains(t, "func With") {
						// look for assignment to package-level in option body — weak
						if strings.Contains(unit.Source, "defaultConfig") || strings.Contains(unit.Source, "defaultOptions") {
							pushAt(unit, meta, line.byte, "functional option mutates a package-level default", out)
							return
						}
					}
				}
			}
		}
	}
}

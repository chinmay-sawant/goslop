package badpractices

import (
	goast "go/ast"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-26", detectBP26)
	RegisterRule("BP-27", detectBP27)
	RegisterRule("BP-28", detectBP28)
	RegisterRule("BP-29", detectBP29)
	RegisterRule("BP-30", detectBP30)
	RegisterRule("BP-31", detectBP31)
	RegisterRule("BP-32", detectBP32)
	RegisterRule("BP-33", detectBP33)
	RegisterRule("BP-34", detectBP34)
	RegisterRule("BP-36", detectBP36)
	RegisterRule("BP-37", detectBP37)
	RegisterRule("BP-38", detectBP38)
	RegisterRule("BP-39", detectBP39)
	RegisterRule("BP-40", detectBP40)
	RegisterRule("BP-41", detectBP41)
	RegisterRule("BP-42", detectBP42)
	RegisterRule("BP-43", detectBP43)
	RegisterRule("BP-44", detectBP44)
	RegisterRule("BP-45", detectBP45)
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
		needle := "interface {"
		idx := strings.Index(src[start:], needle)
		if idx < 0 {
			needle = "interface{"
			idx = strings.Index(src[start:], needle)
			if idx < 0 {
				break
			}
		}
		abs := start + idx
		// count methods roughly
		end, ok := matchBraceBlock(src, abs)
		if !ok {
			// Needle often appears inside string literals (including this
			// detector's own source). Always advance past the match.
			start = abs + len(needle)
			continue
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
	const needle = "interface {"
	start := 0
	for {
		idx := strings.Index(src[start:], needle)
		if idx < 0 {
			break
		}
		abs := start + idx
		end, ok := matchBraceBlock(src, abs)
		if !ok {
			// Unbalanced braces from a false match (e.g. needle in a string)
			// would otherwise set start=abs and loop forever.
			start = abs + len(needle)
			continue
		}
		block := src[abs:end]
		methods := 0
		for _, line := range strings.Split(block, "\n") {
			t := strings.TrimSpace(line)
			if strings.Contains(t, "(") && !strings.HasPrefix(t, "interface") {
				methods++
			}
		}
		// Rust: more than five methods.
		if methods > 5 {
			pushAt(unit, meta, abs, "interface declares more than five methods and is likely too broad", out)
		}
		start = end
	}
}

// matchBraceBlock finds the end offset (exclusive) of the brace-balanced region
// starting at abs. abs must point at or before the opening '{'. Returns false
// when braces never balance (string/comment false positives).
func matchBraceBlock(src string, abs int) (int, bool) {
	if abs < 0 || abs >= len(src) {
		return abs, false
	}
	depth := 0
	for i := abs; i < len(src); i++ {
		switch src[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1, true
			}
		}
	}
	return abs, false
}

// detectBP30: exported interface with no evident same-package implementation.
// Uses the whole package directory (Rust package_method_sets parity).
func detectBP30(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	if !strings.Contains(unit.Source, "interface") {
		return
	}
	meta := MetadataForID("BP-30")
	src := unit.Source
	pkgFacts := packageTypeFactsForUnit(unit)

	// Only report interfaces declared in this file (per-file findings).
	for _, m := range reExportedIface.FindAllStringSubmatchIndex(src, -1) {
		ifaceName := src[m[2]:m[3]]
		methods, ok := pkgFacts.interfaces[ifaceName]
		if !ok || len(methods) == 0 {
			// re-parse local body if package scan missed embedding noise
			continue
		}
		hasImpl := false
		for typ := range pkgFacts.methods {
			if typ == ifaceName {
				continue
			}
			if typeImplements(pkgFacts, typ, methods) {
				hasImpl = true
				break
			}
		}
		if !hasImpl {
			pushAt(unit, meta, m[0], "exported interface has no evident same-package implementation", out)
		}
	}
}

// detectBP31: New* constructor returns a concrete type even though a fitting
// package interface already exists (Rust api_design::detect_bp_31).
func detectBP31(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	if !strings.Contains(unit.Source, "func New") {
		return
	}
	meta := MetadataForID("BP-31")
	pkgFacts := packageTypeFactsForUnit(unit)
	if len(pkgFacts.interfaces) == 0 {
		return
	}
	lines := strings.Split(unit.Source, "\n")
	byteOff := 0
	for _, line := range lines {
		t := strings.TrimSpace(line)
		// package-level constructors only (not methods)
		if strings.HasPrefix(t, "func New") || strings.HasPrefix(t, "func Must") {
			// func NewX(...) *Type / Type
			if strings.HasPrefix(t, "func (") {
				byteOff += len(line) + 1
				continue
			}
			name := firstIdent(strings.TrimPrefix(t, "func "))
			if name == "" || !unicode.IsUpper(rune(name[0])) {
				byteOff += len(line) + 1
				continue
			}
			if !strings.HasPrefix(name, "New") && !strings.HasPrefix(name, "Must") {
				byteOff += len(line) + 1
				continue
			}
			returned := firstResultTypeFromFuncLine(t)
			if returned == "" || !unicode.IsUpper(rune(returned[0])) {
				byteOff += len(line) + 1
				continue
			}
			methods := pkgFacts.methods[returned]
			if len(methods) == 0 {
				byteOff += len(line) + 1
				continue
			}
			exposes := false
			for ifaceName, ifaceMethods := range pkgFacts.interfaces {
				if ifaceName == returned || len(ifaceMethods) == 0 {
					continue
				}
				if typeImplements(pkgFacts, returned, ifaceMethods) {
					exposes = true
					break
				}
			}
			if exposes {
				pushAt(unit, meta, byteOff,
					"constructor returns a concrete type even though the package already exposes a fitting interface",
					out)
			}
		}
		byteOff += len(line) + 1
	}
}

// firstResultTypeFromFuncLine extracts the first named result type from a single-line
// func signature (best-effort; multi-line signatures are skipped).
func firstResultTypeFromFuncLine(line string) string {
	// func Name(...) results {
	closeParams := strings.Index(line, ")")
	if closeParams < 0 {
		return ""
	}
	rest := strings.TrimSpace(line[closeParams+1:])
	if rest == "" || strings.HasPrefix(rest, "{") {
		return ""
	}
	// strip trailing {
	if i := strings.Index(rest, "{"); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	// (T, error) or *T or (a T, err error)
	rest = strings.TrimSpace(rest)
	if strings.HasPrefix(rest, "(") {
		inner := strings.TrimSuffix(strings.TrimPrefix(rest, "("), ")")
		parts := strings.Split(inner, ",")
		if len(parts) == 0 {
			return ""
		}
		return normalizeResultType(parts[0])
	}
	// single result
	if i := strings.IndexAny(rest, " \t"); i > 0 {
		rest = rest[:i]
	}
	return normalizeResultType(rest)
}

func normalizeResultType(s string) string {
	s = strings.TrimSpace(s)
	// named result: "err error" / "v *Type"
	fields := strings.Fields(s)
	if len(fields) >= 2 {
		s = fields[len(fields)-1]
	}
	s = strings.TrimPrefix(s, "*")
	// drop package qualifier
	if i := strings.LastIndex(s, "."); i >= 0 {
		s = s[i+1:]
	}
	// only plain identifiers
	if s == "" || !unicode.IsLetter(rune(s[0])) {
		return ""
	}
	return s
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
	// fmt.Errorf without %w when wrapping err — require bare `err` arg, not errCount.
	emitted := 0
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.Contains(t, "fmt.Errorf(") || strings.Contains(t, "%w") {
			continue
		}
		if !bp34WrapsErr(t) {
			continue
		}
		pushAt(unit, meta, line.byte, "error wrapping without %w loses the error chain", out)
		emitted++
		if emitted >= 2 {
			return
		}
	}
}

// bp34WrapsErr reports whether fmt.Errorf call arguments include the bare
// identifier `err` (word boundary), not errCount / errno / etc.
func bp34WrapsErr(line string) bool {
	// After fmt.Errorf(, scan for , err or (err as standalone ident.
	idx := strings.Index(line, "fmt.Errorf(")
	if idx < 0 {
		return false
	}
	args := line[idx+len("fmt.Errorf("):]
	// Strip trailing comments.
	if i := strings.Index(args, "//"); i >= 0 {
		args = args[:i]
	}
	for i := 0; i < len(args); {
		// find "err"
		j := strings.Index(args[i:], "err")
		if j < 0 {
			return false
		}
		j += i
		// word start
		if j > 0 {
			prev := args[j-1]
			if (prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z') ||
				(prev >= '0' && prev <= '9') || prev == '_' {
				i = j + 3
				continue
			}
		}
		// word end
		end := j + 3
		if end < len(args) {
			next := args[end]
			if (next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') ||
				(next >= '0' && next <= '9') || next == '_' {
				i = j + 3
				continue
			}
		}
		return true
	}
	return false
}

func detectBP36(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust: init body contains call_expression / go / defer.
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-36")
	pos := strings.Index(unit.Source, "func init(")
	if pos < 0 {
		return
	}
	body := unit.Source[pos:]
	if end := strings.Index(body, "\nfunc "); end > 0 {
		body = body[:end]
	}
	// Extract brace body
	open := strings.Index(body, "{")
	if open < 0 {
		return
	}
	depth := 0
	end := open
	for i := open; i < len(body); i++ {
		if body[i] == '{' {
			depth++
		} else if body[i] == '}' {
			depth--
			if depth == 0 {
				end = i
				break
			}
		}
	}
	inner := body[open+1 : end]
	// Side effect: any call/conversion (contains '(' not in string — crude) or go/defer.
	if strings.Contains(inner, "go ") || strings.Contains(inner, "defer ") ||
		strings.Contains(inner, "(") {
		pushAt(unit, meta, pos, "init() performs side effects beyond simple package setup", out)
	}
}

// BP-37: package-level var that receives a post-init write (Rust parity).
// Flags the package var_declaration when any of its names is written via
// assignment, index write, field write, send, or ++/--. Reads and method calls
// are not writes. Shadowed locals (params, short decls) are excluded.
func detectBP37(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-37")
	msg := "package-level mutable global state makes behavior harder to reason about"
	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		detectBP37AST(unit, facts, meta, msg, out)
		return
	}
	// Text fallback (no AST): package-level names later clearly assigned.
	globals := packageLevelVarNames(unit.Source)
	if len(globals) == 0 {
		return
	}
	written := packageLevelVarsWritten(unit.Source, globals)
	if len(written) == 0 {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "var ") {
			continue
		}
		for name := range written {
			if strings.HasPrefix(name, "Err") {
				continue
			}
			rest := strings.TrimSpace(strings.TrimPrefix(t, "var "))
			if !strings.HasPrefix(rest, name) {
				continue
			}
			afterName := rest[len(name):]
			if afterName != "" {
				c := afterName[0]
				if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
					continue
				}
			}
			pushAt(unit, meta, line.byte, msg, out)
			break
		}
	}
}

func detectBP37AST(unit *core.ParsedUnit, facts *bpFacts, meta *rules.RuleMetadata, msg string, out *[]rules.Finding) {
	tree := facts.tree
	type varDecl struct {
		offset int
		names  []string
	}
	var decls []varDecl
	globals := map[string]struct{}{}
	for _, decl := range tree.File.Decls {
		gd, ok := decl.(*goast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		var names []string
		for _, sp := range gd.Specs {
			vs, ok := sp.(*goast.ValueSpec)
			if !ok {
				continue
			}
			for _, n := range vs.Names {
				if n == nil || n.Name == "" || strings.HasPrefix(n.Name, "Err") {
					continue
				}
				names = append(names, n.Name)
				globals[n.Name] = struct{}{}
			}
		}
		if len(names) > 0 {
			decls = append(decls, varDecl{offset: tree.Offset(gd.Pos()), names: names})
		}
	}
	if len(globals) == 0 {
		return
	}
	written := collectWrittenGlobalsAST(tree.File, globals)
	for _, d := range decls {
		for _, name := range d.names {
			if _, ok := written[name]; ok {
				pushAt(unit, meta, d.offset, msg, out)
				break
			}
		}
	}
}

func collectWrittenGlobalsAST(file *goast.File, globals map[string]struct{}) map[string]struct{} {
	written := map[string]struct{}{}
	var walk func(n goast.Node, shadowed map[string]struct{})
	walk = func(n goast.Node, shadowed map[string]struct{}) {
		if n == nil {
			return
		}
		switch x := n.(type) {
		case *goast.FuncDecl:
			sh := cloneShadowSet(shadowed)
			if x.Recv != nil {
				addFieldListNames(x.Recv, globals, sh)
			}
			if x.Type != nil {
				if x.Type.Params != nil {
					addFieldListNames(x.Type.Params, globals, sh)
				}
				if x.Type.Results != nil {
					addFieldListNames(x.Type.Results, globals, sh)
				}
			}
			if x.Body != nil {
				walkBlock(x.Body, globals, sh, written)
			}
			return
		case *goast.FuncLit:
			sh := cloneShadowSet(shadowed)
			if x.Type != nil {
				if x.Type.Params != nil {
					addFieldListNames(x.Type.Params, globals, sh)
				}
				if x.Type.Results != nil {
					addFieldListNames(x.Type.Results, globals, sh)
				}
			}
			if x.Body != nil {
				walkBlock(x.Body, globals, sh, written)
			}
			return
		}
		// Top-level: only walk into func decls (handled above via File.Decls).
		goast.Inspect(n, func(child goast.Node) bool {
			if child == nil || child == n {
				return true
			}
			switch child.(type) {
			case *goast.FuncDecl, *goast.FuncLit:
				walk(child, shadowed)
				return false
			}
			return true
		})
	}
	for _, decl := range file.Decls {
		if fd, ok := decl.(*goast.FuncDecl); ok {
			walk(fd, map[string]struct{}{})
		}
	}
	return written
}

func walkBlock(block *goast.BlockStmt, globals, shadowed map[string]struct{}, written map[string]struct{}) {
	if block == nil {
		return
	}
	sh := cloneShadowSet(shadowed)
	for _, stmt := range block.List {
		walkStmt(stmt, globals, sh, written)
		// Bindings from this statement shadow later siblings (Rust-like).
		addStmtBindings(stmt, globals, sh)
	}
}

func walkStmt(stmt goast.Stmt, globals, shadowed map[string]struct{}, written map[string]struct{}) {
	if stmt == nil {
		return
	}
	switch s := stmt.(type) {
	case *goast.AssignStmt:
		// Short declarations bind locals — do not treat LHS idents as package writes.
		// Walk RHS first so nested func lits can still write globals.
		for _, rhs := range s.Rhs {
			walkExpr(rhs, globals, shadowed, written)
		}
		if s.Tok == token.DEFINE {
			return
		}
		for _, lhs := range s.Lhs {
			collectWriteTargets(lhs, globals, shadowed, written)
		}
	case *goast.IncDecStmt:
		collectWriteTargets(s.X, globals, shadowed, written)
	case *goast.SendStmt:
		collectWriteTargets(s.Chan, globals, shadowed, written)
		walkExpr(s.Value, globals, shadowed, written)
	case *goast.DeclStmt:
		// local var decl — not a package write; bindings added by addStmtBindings
		if gd, ok := s.Decl.(*goast.GenDecl); ok {
			for _, sp := range gd.Specs {
				if vs, ok := sp.(*goast.ValueSpec); ok {
					for _, v := range vs.Values {
						walkExpr(v, globals, shadowed, written)
					}
				}
			}
		}
	case *goast.BlockStmt:
		walkBlock(s, globals, shadowed, written)
	case *goast.IfStmt:
		sh := cloneShadowSet(shadowed)
		if s.Init != nil {
			walkStmt(s.Init, globals, sh, written)
			addStmtBindings(s.Init, globals, sh)
		}
		walkExpr(s.Cond, globals, sh, written)
		walkBlock(s.Body, globals, sh, written)
		if s.Else != nil {
			walkStmt(s.Else, globals, sh, written)
		}
	case *goast.ForStmt:
		sh := cloneShadowSet(shadowed)
		if s.Init != nil {
			walkStmt(s.Init, globals, sh, written)
			addStmtBindings(s.Init, globals, sh)
		}
		walkExpr(s.Cond, globals, sh, written)
		if s.Post != nil {
			walkStmt(s.Post, globals, sh, written)
		}
		walkBlock(s.Body, globals, sh, written)
	case *goast.RangeStmt:
		sh := cloneShadowSet(shadowed)
		// range binds key/value before body
		if s.Tok == token.DEFINE {
			collectBindingIdents(s.Key, globals, sh)
			collectBindingIdents(s.Value, globals, sh)
		} else {
			collectWriteTargets(s.Key, globals, sh, written)
			collectWriteTargets(s.Value, globals, sh, written)
		}
		walkExpr(s.X, globals, sh, written)
		walkBlock(s.Body, globals, sh, written)
	case *goast.SwitchStmt:
		sh := cloneShadowSet(shadowed)
		if s.Init != nil {
			walkStmt(s.Init, globals, sh, written)
			addStmtBindings(s.Init, globals, sh)
		}
		walkExpr(s.Tag, globals, sh, written)
		if s.Body != nil {
			for _, c := range s.Body.List {
				cc, ok := c.(*goast.CaseClause)
				if !ok {
					continue
				}
				csh := cloneShadowSet(sh)
				for _, e := range cc.List {
					walkExpr(e, globals, csh, written)
				}
				for _, st := range cc.Body {
					walkStmt(st, globals, csh, written)
					addStmtBindings(st, globals, csh)
				}
			}
		}
	case *goast.TypeSwitchStmt:
		sh := cloneShadowSet(shadowed)
		if s.Init != nil {
			walkStmt(s.Init, globals, sh, written)
			addStmtBindings(s.Init, globals, sh)
		}
		// assign: v := x.(type) — walk RHS (and nested writes) before binding the alias,
		// matching Rust type_switch_statement handling for guard RHS writes.
		if as, ok := s.Assign.(*goast.AssignStmt); ok {
			for _, rhs := range as.Rhs {
				walkExpr(rhs, globals, sh, written)
			}
			if as.Tok == token.DEFINE {
				for _, lhs := range as.Lhs {
					collectBindingIdents(lhs, globals, sh)
				}
			} else {
				for _, lhs := range as.Lhs {
					collectWriteTargets(lhs, globals, sh, written)
				}
			}
		}
		if s.Body != nil {
			for _, c := range s.Body.List {
				cc, ok := c.(*goast.CaseClause)
				if !ok {
					continue
				}
				csh := cloneShadowSet(sh)
				for _, st := range cc.Body {
					walkStmt(st, globals, csh, written)
					addStmtBindings(st, globals, csh)
				}
			}
		}
	case *goast.SelectStmt:
		if s.Body == nil {
			return
		}
		for _, c := range s.Body.List {
			cc, ok := c.(*goast.CommClause)
			if !ok {
				continue
			}
			csh := cloneShadowSet(shadowed)
			if cc.Comm != nil {
				if as, ok := cc.Comm.(*goast.AssignStmt); ok {
					if as.Tok == token.DEFINE {
						for _, lhs := range as.Lhs {
							collectBindingIdents(lhs, globals, csh)
						}
					} else {
						for _, lhs := range as.Lhs {
							collectWriteTargets(lhs, globals, csh, written)
						}
					}
					for _, rhs := range as.Rhs {
						walkExpr(rhs, globals, csh, written)
					}
				} else if send, ok := cc.Comm.(*goast.SendStmt); ok {
					collectWriteTargets(send.Chan, globals, csh, written)
					walkExpr(send.Value, globals, csh, written)
				} else if es, ok := cc.Comm.(*goast.ExprStmt); ok {
					walkExpr(es.X, globals, csh, written)
				}
			}
			for _, st := range cc.Body {
				walkStmt(st, globals, csh, written)
				addStmtBindings(st, globals, csh)
			}
		}
	case *goast.GoStmt:
		walkExpr(s.Call, globals, shadowed, written)
	case *goast.DeferStmt:
		walkExpr(s.Call, globals, shadowed, written)
	case *goast.ExprStmt:
		walkExpr(s.X, globals, shadowed, written)
	case *goast.ReturnStmt:
		for _, r := range s.Results {
			walkExpr(r, globals, shadowed, written)
		}
	default:
		// Best-effort: discover nested function literals under uncommon stmts.
		goast.Inspect(stmt, func(n goast.Node) bool {
			if fl, ok := n.(*goast.FuncLit); ok {
				walkExpr(fl, globals, shadowed, written)
				return false
			}
			return true
		})
	}
}

func walkExpr(e goast.Expr, globals, shadowed map[string]struct{}, written map[string]struct{}) {
	if e == nil {
		return
	}
	goast.Inspect(e, func(n goast.Node) bool {
		fl, ok := n.(*goast.FuncLit)
		if !ok {
			return true
		}
		sh := cloneShadowSet(shadowed)
		if fl.Type != nil {
			if fl.Type.Params != nil {
				addFieldListNames(fl.Type.Params, globals, sh)
			}
			if fl.Type.Results != nil {
				addFieldListNames(fl.Type.Results, globals, sh)
			}
		}
		if fl.Body != nil {
			walkBlock(fl.Body, globals, sh, written)
		}
		return false
	})
}

func collectWriteTargets(e goast.Expr, globals, shadowed map[string]struct{}, written map[string]struct{}) {
	if e == nil {
		return
	}
	switch x := e.(type) {
	case *goast.Ident:
		if _, ok := globals[x.Name]; ok {
			if _, sh := shadowed[x.Name]; !sh {
				written[x.Name] = struct{}{}
			}
		}
	case *goast.SelectorExpr:
		// opts.Field or imgCache.cache — base ident is the package global
		if id, ok := x.X.(*goast.Ident); ok {
			if _, g := globals[id.Name]; g {
				if _, sh := shadowed[id.Name]; !sh {
					written[id.Name] = struct{}{}
				}
			}
		} else {
			collectWriteTargets(x.X, globals, shadowed, written)
		}
	case *goast.IndexExpr:
		// name[i] or name.field[i]
		collectWriteTargets(x.X, globals, shadowed, written)
	case *goast.IndexListExpr:
		collectWriteTargets(x.X, globals, shadowed, written)
	case *goast.StarExpr:
		collectWriteTargets(x.X, globals, shadowed, written)
	case *goast.ParenExpr:
		collectWriteTargets(x.X, globals, shadowed, written)
	}
}

func collectBindingIdents(e goast.Expr, globals, shadowed map[string]struct{}) {
	if e == nil {
		return
	}
	if id, ok := e.(*goast.Ident); ok && id.Name != "_" {
		if _, g := globals[id.Name]; g {
			shadowed[id.Name] = struct{}{}
		}
	}
}

func addFieldListNames(fl *goast.FieldList, globals, shadowed map[string]struct{}) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, n := range f.Names {
			if n != nil && n.Name != "_" {
				if _, g := globals[n.Name]; g {
					shadowed[n.Name] = struct{}{}
				}
			}
		}
	}
}

func addStmtBindings(stmt goast.Stmt, globals, shadowed map[string]struct{}) {
	if stmt == nil {
		return
	}
	switch s := stmt.(type) {
	case *goast.AssignStmt:
		if s.Tok != token.DEFINE {
			return
		}
		for _, lhs := range s.Lhs {
			collectBindingIdents(lhs, globals, shadowed)
		}
	case *goast.DeclStmt:
		gd, ok := s.Decl.(*goast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			return
		}
		for _, sp := range gd.Specs {
			vs, ok := sp.(*goast.ValueSpec)
			if !ok {
				continue
			}
			for _, n := range vs.Names {
				if n != nil && n.Name != "_" {
					if _, g := globals[n.Name]; g {
						shadowed[n.Name] = struct{}{}
					}
				}
			}
		}
	}
}

func cloneShadowSet(in map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{}, len(in))
	for k := range in {
		out[k] = struct{}{}
	}
	return out
}

func packageLevelVarNames(source string) map[string]struct{} {
	out := map[string]struct{}{}
	inFunc := 0
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		open := strings.Count(t, "{")
		closeN := strings.Count(t, "}")
		if inFunc == 0 && strings.HasPrefix(t, "var ") {
			rest := strings.TrimSpace(strings.TrimPrefix(t, "var "))
			if !strings.HasPrefix(rest, "(") {
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
	inFunc := false
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "func ") {
			inFunc = true
		}
		if !inFunc {
			continue
		}
		// Package-level vars declared after functions are not in-function writes.
		if strings.HasPrefix(t, "var ") || strings.HasPrefix(t, "const ") || strings.HasPrefix(t, "type ") {
			continue
		}
		for name := range globals {
			if strings.Contains(t, name+" :=") || strings.Contains(t, name+":=") {
				continue
			}
			// Direct assign / compound assign (not bare index read).
			if strings.Contains(t, name+" =") || strings.Contains(t, name+"=") ||
				strings.Contains(t, name+" +=") || strings.Contains(t, name+"+=") ||
				strings.Contains(t, name+"++") || strings.Contains(t, name+"--") {
				if isMethodCallOnly(t, name) {
					continue
				}
				written[name] = struct{}{}
				continue
			}
			// Index write: name[i] = ...
			if strings.Contains(t, name+"[") && strings.Contains(t, "=") && !strings.Contains(t, ":=") {
				// likely write if '=' after the index
				idx := strings.Index(t, name+"[")
				rest := t[idx:]
				if eq := strings.Index(rest, "="); eq >= 0 && !strings.Contains(rest[:eq], ":=") {
					written[name] = struct{}{}
					continue
				}
			}
			// Field write: name.field =
			if strings.Contains(t, name+".") {
				dot := strings.Index(t, name+".")
				after := t[dot+len(name)+1:]
				identEnd := 0
				for identEnd < len(after) && isIdentByteBP(after[identEnd]) {
					identEnd++
				}
				if identEnd < len(after) && after[identEnd] == '(' {
					continue // method call
				}
				if strings.Contains(after, "=") {
					written[name] = struct{}{}
				}
			}
		}
	}
	return written
}

func isMethodCallOnly(line, name string) bool {
	if strings.Contains(line, name+" =") || strings.Contains(line, name+"=") ||
		strings.Contains(line, name+"[") {
		return false
	}
	return strings.Contains(line, name+".")
}

func isIdentByteBP(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
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
	// Rust: exported API (funcs + methods on exported receivers) need a doc
	// comment that starts with the function/method name. Emit all hits (no
	// early return) for real-repos parity.
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-39")
	lines := strings.Split(unit.Source, "\n")
	byteOff := 0
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if !strings.HasPrefix(t, "func ") {
			byteOff += len(line) + 1
			continue
		}
		rest := strings.TrimPrefix(t, "func ")
		name := ""
		if strings.HasPrefix(rest, "(") {
			// method: func (r *Type) Name
			closeParen := strings.Index(rest, ")")
			if closeParen < 0 {
				byteOff += len(line) + 1
				continue
			}
			recv := strings.TrimSpace(rest[:closeParen+1]) // includes leading (
			recvInner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(recv, "("), ")"))
			fields := strings.Fields(recvInner)
			if len(fields) == 0 {
				byteOff += len(line) + 1
				continue
			}
			recvType := strings.TrimPrefix(fields[len(fields)-1], "*")
			// Rust is_exported_api: only methods on exported receivers.
			if recvType == "" || !unicode.IsUpper(rune(recvType[0])) {
				byteOff += len(line) + 1
				continue
			}
			after := strings.TrimSpace(rest[closeParen+1:])
			name = firstIdent(after)
		} else {
			name = firstIdent(rest)
		}
		if name == "" || !unicode.IsUpper(rune(name[0])) {
			byteOff += len(line) + 1
			continue
		}
		// Collect contiguous // comments immediately above (blank lines stop after docs start).
		var docs []string
		for j := i - 1; j >= 0; j-- {
			pt := strings.TrimSpace(lines[j])
			if pt == "" {
				break // Rust breaks on empty (after reading upward)
			}
			if strings.HasPrefix(pt, "//") {
				docs = append([]string{strings.TrimSpace(strings.TrimPrefix(pt, "//"))}, docs...)
				continue
			}
			break
		}
		ok := len(docs) > 0 && strings.HasPrefix(docs[0], name)
		if !ok {
			pushAt(unit, meta, byteOff, "exported API should have a doc comment that starts with its name", out)
		}
		byteOff += len(line) + 1
	}
}

// detectBP40: package-level const block groups unrelated names (Rust
// detect_bp_40_unrelated_constants_in_one_block). Requires ≥3 names whose
// constant_prefix set has size > 2.
func detectBP40(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-40")
	msg := "const block groups unrelated constants together"

	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		for _, decl := range facts.tree.File.Decls {
			gd, ok := decl.(*goast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}
			names := constDeclNames(gd)
			if len(names) < 3 {
				continue
			}
			prefixes := make(map[string]struct{}, len(names))
			for _, name := range names {
				prefixes[constantPrefix(name)] = struct{}{}
			}
			if len(prefixes) > 2 {
				pushAt(unit, meta, facts.tree.Offset(gd.Pos()), msg, out)
			}
		}
		return
	}

	// Text fallback when AST is unavailable.
	detectBP40Text(unit, meta, msg, out)
}

func constDeclNames(gd *goast.GenDecl) []string {
	var names []string
	for _, spec := range gd.Specs {
		vs, ok := spec.(*goast.ValueSpec)
		if !ok {
			continue
		}
		for _, id := range vs.Names {
			if id != nil && id.Name != "" {
				names = append(names, id.Name)
			}
		}
	}
	return names
}

// constantPrefix mirrors Rust code_organization::constant_prefix:
// underscore segment before first '_', else leading camelCase run until the
// next ASCII uppercase letter.
func constantPrefix(name string) string {
	if i := strings.IndexByte(name, '_'); i >= 0 {
		return name[:i]
	}
	var b strings.Builder
	for _, ch := range name {
		if ch >= 'A' && ch <= 'Z' && b.Len() > 0 {
			break
		}
		b.WriteRune(ch)
	}
	if b.Len() == 0 {
		return name
	}
	return b.String()
}

func detectBP40Text(unit *core.ParsedUnit, meta *rules.RuleMetadata, msg string, out *[]rules.Finding) {
	src := unit.Source
	lines := strings.Split(src, "\n")
	byteOff := 0
	i := 0
	for i < len(lines) {
		line := lines[i]
		t := strings.TrimSpace(stripLineComment(line))
		// const ( ... ) block only (Rust const_declaration with multiple specs)
		if t == "const (" || strings.HasPrefix(t, "const(") {
			blockStart := byteOff
			var names []string
			i++
			byteOff += len(line) + 1
			for i < len(lines) {
				inner := lines[i]
				it := strings.TrimSpace(stripLineComment(inner))
				if it == ")" {
					break
				}
				if it != "" && !strings.HasPrefix(it, "//") {
					// left of '=' or whole line; first identifier per line (may be multi-name)
					left := it
					if eq := strings.Index(it, "="); eq >= 0 {
						left = strings.TrimSpace(it[:eq])
					}
					for _, part := range strings.Split(left, ",") {
						part = strings.TrimSpace(part)
						if part == "" {
							continue
						}
						// take first token (name [type])
						name := strings.Fields(part)
						if len(name) == 0 {
							continue
						}
						ident := name[0]
						if isSimpleIdent(ident) {
							names = append(names, ident)
						}
					}
				}
				byteOff += len(inner) + 1
				i++
			}
			if len(names) >= 3 {
				prefixes := make(map[string]struct{}, len(names))
				for _, name := range names {
					prefixes[constantPrefix(name)] = struct{}{}
				}
				if len(prefixes) > 2 {
					pushAt(unit, meta, blockStart, msg, out)
				}
			}
			if i < len(lines) {
				byteOff += len(lines[i]) + 1
				i++
			}
			continue
		}
		byteOff += len(line) + 1
		i++
	}
}

func isSimpleIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
			continue
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

// detectBP42: import alias used only once (Rust count_word_occurrences <= 2).
func detectBP42(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	if isTestFile(unit) {
		return
	}
	meta := MetadataForID("BP-42")
	src := unit.Source
	// Single-line: import alias "path"  (not import "path")
	reSingle := regexp.MustCompile(`(?m)^import\s+([A-Za-z_]\w*)\s+"[^"]+"`)
	// Inside import ( ): alias "path"  (not bare "path")
	reBlock := regexp.MustCompile(`(?m)^\s*([A-Za-z_]\w*)\s+"[^"]+"`)
	inImport := false
	byteOff := 0
	for _, line := range strings.Split(src, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "import ") || t == "import (" || strings.HasPrefix(t, "import(") {
			if strings.Contains(t, "(") {
				inImport = true
				// also handle: import ( alias "path" ) rare one-liners
			} else if m := reSingle.FindStringSubmatchIndex(line); m != nil {
				alias := line[m[2]:m[3]]
				if isUsefulImportAlias(alias) && countWordOccurrences(src, alias) <= 2 {
					pushAt(unit, meta, byteOff+m[2],
						"import alias is only used once and likely adds indirection without value", out)
				}
			}
			byteOff += len(line) + 1
			continue
		}
		if inImport {
			if t == ")" {
				inImport = false
				byteOff += len(line) + 1
				continue
			}
			// skip bare "path" imports (no alias)
			if strings.HasPrefix(t, "\"") || strings.HasPrefix(t, ". \"") || strings.HasPrefix(t, "_ \"") {
				byteOff += len(line) + 1
				continue
			}
			if m := reBlock.FindStringSubmatchIndex(line); m != nil {
				alias := line[m[2]:m[3]]
				if isUsefulImportAlias(alias) && countWordOccurrences(src, alias) <= 2 {
					pushAt(unit, meta, byteOff+m[2],
						"import alias is only used once and likely adds indirection without value", out)
				}
			}
		}
		byteOff += len(line) + 1
	}
}

func isUsefulImportAlias(alias string) bool {
	if alias == "" || alias == "_" || alias == "." {
		return false
	}
	// never treat the import keyword / common false parses as an alias
	switch alias {
	case "import", "package", "var", "const", "type", "func", "go", "defer", "map", "chan", "interface", "struct":
		return false
	}
	return true
}

func detectBP41(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust parity: only the package anchor file reports missing package doc.
	// Skip flat materializations (…/go/file.go) to avoid noise, but still run
	// on nested package fixtures such as go/bp41/*.go (BP-41-vulnerable).
	if isTestFile(unit) || isFlatMaterializedFixture(unit) {
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
	if isTestFile(unit) {
		return
	}
	msg := "blank import should carry a justification or match a standard registration pattern"
	lines := strings.Split(unit.Source, "\n")
	byteOff := 0
	for i, line := range lines {
		t := strings.TrimSpace(line)
		path, isBlank, pathOff := parseBlankImportLine(t)
		if isBlank {
			if !isAllowedBlankImport(path) && !hasBlankImportJustification(lines, i) {
				pushAt(unit, meta, byteOff+pathOff, msg, out)
			}
		}
		byteOff += len(line) + 1
	}
}

// parseBlankImportLine recognizes `import _ "path"`, `_ "path"`, and backtick forms.
// Returns (importPath, ok, byteOffsetOfBlankInLine).
func parseBlankImportLine(t string) (path string, ok bool, off int) {
	// import _ "path"
	for _, prefix := range []string{`import _ "`, "import _ `", `import _"`, "import _`"} {
		if strings.HasPrefix(t, prefix) {
			rest := t[len(prefix):]
			quote := byte('"')
			if strings.Contains(prefix, "`") {
				quote = '`'
			}
			end := strings.IndexByte(rest, quote)
			if end < 0 {
				return "", false, 0
			}
			return rest[:end], true, strings.Index(t, "_")
		}
	}
	// _ "path" inside import block
	for _, prefix := range []string{`_ "`, "_ `", `_"`, "_`"} {
		if strings.HasPrefix(t, prefix) {
			rest := t[len(prefix):]
			quote := byte('"')
			if strings.Contains(prefix, "`") {
				quote = '`'
			}
			end := strings.IndexByte(rest, quote)
			if end < 0 {
				return "", false, 0
			}
			return rest[:end], true, 0
		}
	}
	return "", false, 0
}

func isAllowedBlankImport(path string) bool {
	return strings.HasPrefix(path, "image/") ||
		strings.Contains(path, "driver") ||
		strings.HasSuffix(path, "/pprof") ||
		strings.Contains(path, "plugin")
}

func hasBlankImportJustification(lines []string, lineNo int) bool {
	current := ""
	if lineNo < len(lines) {
		current = lines[lineNo]
	}
	previous := ""
	if lineNo > 0 {
		previous = lines[lineNo-1]
	}
	context := strings.ToLower(previous + "\n" + current)
	return strings.Contains(context, "register") ||
		strings.Contains(context, "side effect") ||
		strings.Contains(context, "side-effect") ||
		strings.Contains(context, "plugin") ||
		strings.Contains(current, "//")
}

// BP-45: methods on the same type should use a consistent receiver name.
func detectBP45(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-45")
	if isTestFile(unit) {
		return
	}
	msg := "methods on the same receiver type should use a consistent receiver name"
	byType := map[string]string{} // type → first receiver name
	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		tree := facts.tree
		for _, decl := range tree.File.Decls {
			fd, ok := decl.(*goast.FuncDecl)
			if !ok || fd.Recv == nil || len(fd.Recv.List) == 0 {
				continue
			}
			field := fd.Recv.List[0]
			typeName := receiverTypeName(field.Type)
			if typeName == "" {
				continue
			}
			recvName := ""
			if len(field.Names) > 0 && field.Names[0] != nil {
				recvName = field.Names[0].Name
			}
			if recvName == "" || recvName == "_" {
				continue
			}
			if prev, ok := byType[typeName]; ok {
				if prev != recvName {
					pushAt(unit, meta, tree.Offset(fd.Pos()), msg, out)
				}
			} else {
				byType[typeName] = recvName
			}
		}
		return
	}
	// Text fallback: func (name *Type) / func (name Type)
	re := regexp.MustCompile(`func\s+\((\w+)\s+\*?(\w+)\)`)
	for _, line := range codeLines(unit.Source) {
		m := re.FindStringSubmatch(line.text)
		if m == nil {
			continue
		}
		recvName, typeName := m[1], m[2]
		if prev, ok := byType[typeName]; ok {
			if prev != recvName {
				pushAt(unit, meta, line.byte, msg, out)
			}
		} else {
			byType[typeName] = recvName
		}
	}
}

func receiverTypeName(expr goast.Expr) string {
	switch e := expr.(type) {
	case *goast.Ident:
		return e.Name
	case *goast.StarExpr:
		return receiverTypeName(e.X)
	case *goast.IndexExpr:
		return receiverTypeName(e.X)
	default:
		return ""
	}
}

// BP-33: sentinel-style error type (var ErrX = T{}) missing Is method.
func detectBP33(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-33")
	if isTestFile(unit) {
		return
	}
	msg := "sentinel-style error type is missing an Is(error) bool method"
	// Types that implement Error() string
	errorTypes := map[string]struct{}{}
	isTypes := map[string]struct{}{}
	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		tree := facts.tree
		for _, decl := range tree.File.Decls {
			fd, ok := decl.(*goast.FuncDecl)
			if !ok || fd.Recv == nil || fd.Name == nil || len(fd.Recv.List) == 0 {
				continue
			}
			typeName := receiverTypeName(fd.Recv.List[0].Type)
			if typeName == "" {
				continue
			}
			switch fd.Name.Name {
			case "Error":
				errorTypes[typeName] = struct{}{}
			case "Is":
				isTypes[typeName] = struct{}{}
			}
		}
	} else {
		// Text: func (e T) Error() / Is(
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if !strings.HasPrefix(t, "func (") {
				continue
			}
			// func (e notFoundError) Error()
			if strings.Contains(t, ") Error(") {
				if name := methodReceiverTypeText(t); name != "" {
					errorTypes[name] = struct{}{}
				}
			}
			if strings.Contains(t, ") Is(") {
				if name := methodReceiverTypeText(t); name != "" {
					isTypes[name] = struct{}{}
				}
			}
		}
	}
	if len(errorTypes) == 0 {
		return
	}
	// Sentinel: var Err… = Type{} or Type()
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "var Err") && !strings.HasPrefix(t, "const Err") {
			continue
		}
		for typeName := range errorTypes {
			if _, hasIs := isTypes[typeName]; hasIs {
				continue
			}
			if strings.Contains(t, typeName+"{") || strings.Contains(t, typeName+"(") ||
				strings.Contains(t, " "+typeName) || strings.HasSuffix(t, typeName) {
				pushAt(unit, meta, line.byte, msg, out)
				return
			}
		}
	}
}

func methodReceiverTypeText(funcLine string) string {
	// func (e *notFoundError) Error() string {
	start := strings.Index(funcLine, "(")
	end := strings.Index(funcLine, ")")
	if start < 0 || end <= start {
		return ""
	}
	inner := strings.TrimSpace(funcLine[start+1 : end])
	parts := strings.Fields(inner)
	if len(parts) == 0 {
		return ""
	}
	ty := parts[len(parts)-1]
	return strings.TrimPrefix(ty, "*")
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
	// Map range feeds strings.Join (or similar ordered output) without sort.
	src := unit.Source
	if !strings.Contains(src, "range ") || !strings.Contains(src, "map[") {
		return
	}
	if strings.Contains(src, "sort.") {
		return
	}
	// Ordered-output sink: strings.Join is the fixture expectation; also fmt prints.
	if !strings.Contains(src, "strings.Join") &&
		!strings.Contains(src, "strings.Join(") &&
		!strings.Contains(src, "fmt.Print") &&
		!strings.Contains(src, "json.Marshal") {
		return
	}
	msg := "map iteration feeds ordered output without sorting; collect keys or values and sort before strings.Join"
	for _, line := range codeLines(src) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "for ") && strings.Contains(t, " range ") {
			pushAt(unit, meta, line.byte, msg, out)
			return
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
	src := unit.Source
	// Notification-shaped bool/int channels: constant send + discarded receive.
	// Prefer chan struct{} for pure signals.
	channels := notificationChannelNamesBP91(src)
	if len(channels) == 0 {
		return
	}
	for _, ch := range channels {
		if !hasConstantNotificationSendBP91(src, ch) {
			continue
		}
		if !hasDiscardedReceiveBP91(src, ch) {
			continue
		}
		pos := strings.Index(src, ch)
		if pos < 0 {
			pos = 0
		}
		pushAt(unit, meta, pos, "notification channel carries a boolean/integer payload; use chan struct{}", out)
		return
	}
}

func notificationChannelNamesBP91(src string) []string {
	var names []string
	seen := map[string]bool{}
	for _, line := range codeLines(src) {
		t := strings.TrimSpace(line.text)
		// params or locals: name chan bool / name chan int / make(chan bool
		if !strings.Contains(t, "chan bool") && !strings.Contains(t, "chan int") &&
			!strings.Contains(t, "make(chan bool") && !strings.Contains(t, "make(chan int") {
			continue
		}
		// extract identifiers before "chan"
		for _, part := range strings.FieldsFunc(t, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
		}) {
			if part == "chan" || part == "bool" || part == "int" || part == "func" ||
				part == "make" || part == "var" || part == "type" || part == "struct" {
				continue
			}
			if len(part) == 0 {
				continue
			}
			// rough: name appears near chan bool/int on this line
			if strings.Contains(t, part+" chan") || strings.Contains(t, part+" := make(chan") ||
				strings.Contains(t, part+"= make(chan") || strings.Contains(t, part+" = make(chan") {
				if !seen[part] {
					seen[part] = true
					names = append(names, part)
				}
			}
		}
	}
	return names
}

func hasConstantNotificationSendBP91(src, ch string) bool {
	for _, line := range codeLines(src) {
		t := strings.TrimSpace(line.text)
		// ch <- true / ch <- 1
		if !strings.Contains(t, "<-") {
			continue
		}
		left, right, ok := strings.Cut(t, "<-")
		if !ok {
			continue
		}
		if strings.TrimSpace(left) != ch {
			continue
		}
		val := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(right), ";"))
		if val == "true" || val == "1" {
			return true
		}
	}
	return false
}

func hasDiscardedReceiveBP91(src, ch string) bool {
	for _, line := range codeLines(src) {
		t := strings.TrimSpace(line.text)
		// bare receive statement or select case with discarded value
		if t == "<-"+ch || t == "<- "+ch {
			return true
		}
		if t == "case <-"+ch+":" || t == "case <- "+ch+":" {
			return true
		}
		// select case with empty body line after is still discarded
		if strings.HasPrefix(t, "case <-") && strings.Contains(t, ch) && strings.HasSuffix(t, ":") {
			rest := strings.TrimPrefix(t, "case")
			rest = strings.TrimSpace(rest)
			if rest == "<-"+ch+":" || rest == "<- "+ch+":" {
				return true
			}
		}
	}
	return false
}

func detectBP138(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-138")
	// GORM lifecycle hooks with direct external I/O (http/smtp).
	hooks := []string{
		"BeforeCreate", "AfterCreate", "BeforeSave", "AfterSave",
		"BeforeUpdate", "AfterUpdate", "BeforeDelete", "AfterDelete",
		"BeforeFind", "AfterFind",
	}
	hasHook := false
	hookPos := -1
	for _, h := range hooks {
		if i := strings.Index(unit.Source, h); i >= 0 {
			hasHook = true
			if hookPos < 0 || i < hookPos {
				hookPos = i
			}
		}
	}
	if !hasHook {
		return
	}
	// Require gorm.DB-shaped param somewhere and direct external calls.
	if !strings.Contains(unit.Source, "gorm.DB") && !strings.Contains(unit.Source, "*gorm.DB") {
		return
	}
	if strings.Contains(unit.Source, "http.Get") || strings.Contains(unit.Source, "http.Post") ||
		strings.Contains(unit.Source, "http.Head") || strings.Contains(unit.Source, "http.PostForm") ||
		strings.Contains(unit.Source, "smtp.SendMail") {
		pushAt(unit, meta, hookPos, "external I/O inside a GORM hook; keep hooks local and fast", out)
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
	if isTestFile(unit) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "func With") || !strings.Contains(src, "Option") {
		return
	}
	// Collect package-level var names.
	globals := packageLevelVarNames(src)
	if len(globals) == 0 {
		return
	}
	// Walk With* functions that return Option; flag assignment to package globals.
	lines := codeLines(src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "func With") || !strings.Contains(t, "Option") {
			continue
		}
		// Scan body until next top-level func or end.
		depth := 0
		started := false
		for j := i; j < len(lines); j++ {
			lt := lines[j].text
			for _, ch := range lt {
				if ch == '{' {
					depth++
					started = true
				} else if ch == '}' {
					depth--
				}
			}
			if !started {
				continue
			}
			bodyLine := strings.TrimSpace(lt)
			// assignment: globalX = or globalX[key] =
			if strings.Contains(bodyLine, "=") && !strings.Contains(bodyLine, ":=") &&
				!strings.Contains(bodyLine, "==") && !strings.Contains(bodyLine, "!=") &&
				!strings.Contains(bodyLine, "<=") && !strings.Contains(bodyLine, ">=") {
				for g := range globals {
					if strings.HasPrefix(bodyLine, g+" ") || strings.HasPrefix(bodyLine, g+"=") ||
						strings.HasPrefix(bodyLine, g+".") || strings.HasPrefix(bodyLine, g+"[") {
						pushAt(unit, meta, lines[j].byte, "functional option mutates a package-level default", out)
						return
					}
				}
			}
			if started && depth <= 0 {
				break
			}
		}
	}
}

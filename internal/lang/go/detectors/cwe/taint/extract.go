package taint

import (
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/lang/go/goparse"
)

// ExtractTaintFacts walks the unit AST and collects sources/sinks/sanitizers/assignments.
func ExtractTaintFacts(unit *core.ParsedUnit) TaintAnnotations {
	if unit == nil {
		return TaintAnnotations{}
	}
	state := newExtractionState(unit.Source)
	state.pushScope(ScopePackage, ByteRange{0, len(unit.Source)})

	tree := unitTree(unit)
	if tree != nil && tree.File != nil {
		walkExtract(tree.File, tree, state, nil)
	}
	state.popScope()

	transfers, unsupported := pairChannelTransfers(state.channelSends, state.channelRecvs, state.assignments)
	unsupported = append(unsupported, state.unsupported...)

	return TaintAnnotations{
		Sources:          state.sources,
		Sinks:            state.sinks,
		Sanitizers:       state.sanitizers,
		Assignments:      state.assignments,
		Scopes:           state.scopes,
		FunctionParams:   state.functionParams,
		FunctionRanges:   state.functionRanges,
		UnsupportedFlows: unsupported,
		ChannelTransfers: transfers,
	}
}

// unitTree returns the *goparse.Tree for the unit (shared via unit.Tree).
func unitTree(unit *core.ParsedUnit) *goparse.Tree {
	return goparse.TreeForUnit(unit)
}

type extractionState struct {
	srcBytes        []byte
	source          string
	scopes          []ScopeInfo
	scopeStack      []int
	nextScopeID     int
	currentFunction string
	functionScopes  []int
	sources         []TaintSourceAnnotation
	sinks           []TaintSinkAnnotation
	sanitizers      []TaintSanitizerAnnotation
	assignments     []AssignmentDetail
	functionParams  map[string][]string
	functionRanges  map[string]ByteRange
	unsupported     []UnsupportedFlow
	channelSends    []ChannelSendSite
	channelRecvs    []ChannelRecvSite
	hasHTMLTemplate bool
}

func newExtractionState(source string) *extractionState {
	return &extractionState{
		srcBytes:        []byte(source),
		source:          source,
		functionParams:  map[string][]string{},
		functionRanges:  map[string]ByteRange{},
		hasHTMLTemplate: HasHTMLTemplateImport(source),
	}
}

func (s *extractionState) pushScope(kind ScopeKind, br ByteRange) int {
	id := s.nextScopeID
	s.nextScopeID++
	var parent *int
	if len(s.scopeStack) > 0 {
		p := s.scopeStack[len(s.scopeStack)-1]
		parent = &p
	}
	s.scopes = append(s.scopes, ScopeInfo{
		ID: id, Parent: parent, Kind: kind, ByteRange: br, Function: s.currentFunction,
	})
	s.scopeStack = append(s.scopeStack, id)
	if kind == ScopeFunction {
		s.functionScopes = append(s.functionScopes, id)
	}
	return id
}

func (s *extractionState) popScope() {
	if len(s.scopeStack) == 0 {
		return
	}
	id := s.scopeStack[len(s.scopeStack)-1]
	s.scopeStack = s.scopeStack[:len(s.scopeStack)-1]
	if id < len(s.scopes) && s.scopes[id].Kind == ScopeFunction && len(s.functionScopes) > 0 {
		s.functionScopes = s.functionScopes[:len(s.functionScopes)-1]
	}
}

func (s *extractionState) currentScope() int {
	if len(s.scopeStack) == 0 {
		return 0
	}
	return s.scopeStack[len(s.scopeStack)-1]
}

func (s *extractionState) currentFunctionScope() int {
	if len(s.functionScopes) == 0 {
		return s.currentScope()
	}
	return s.functionScopes[len(s.functionScopes)-1]
}

func walkExtract(n goast.Node, tree *goparse.Tree, state *extractionState, parents []goast.Node) {
	if n == nil {
		return
	}
	var entered *scopeEntry
	var restoreFunc *string

	switch x := n.(type) {
	case *goast.FuncDecl:
		name := functionIdentityDecl(x, tree)
		prev := state.currentFunction
		restoreFunc = &prev
		state.currentFunction = name
		state.functionParams[name] = extractParamNamesFieldList(x.Type)
		state.functionRanges[name] = nodeRange(tree, x)
		entered = &scopeEntry{ScopeFunction, nodeRange(tree, x)}
	case *goast.FuncLit:
		name := "<anonymous>"
		prev := state.currentFunction
		restoreFunc = &prev
		state.currentFunction = name
		state.functionParams[name] = extractParamNamesFieldList(x.Type)
		state.functionRanges[name] = nodeRange(tree, x)
		entered = &scopeEntry{ScopeFunction, nodeRange(tree, x)}
	case *goast.BlockStmt:
		entered = &scopeEntry{ScopeBlock, nodeRange(tree, x)}
	case *goast.IfStmt:
		entered = &scopeEntry{ScopeIf, nodeRange(tree, x)}
	case *goast.ForStmt, *goast.RangeStmt:
		entered = &scopeEntry{ScopeFor, nodeRange(tree, x)}
	case *goast.SwitchStmt, *goast.TypeSwitchStmt:
		entered = &scopeEntry{ScopeSwitch, nodeRange(tree, x)}
	case *goast.CaseClause:
		entered = &scopeEntry{ScopeCase, nodeRange(tree, x)}
	case *goast.CallExpr:
		recordCall(x, tree, state, parents)
	case *goast.SendStmt:
		recordSend(x, tree, state, parents)
	case *goast.GoStmt:
		state.unsupported = append(state.unsupported, UnsupportedFlow{
			Kind:      UnsupportedGoroutine,
			ByteRange: nodeRange(tree, x),
			Note:      "goroutine spawn is not tracked by taint (explicit FN)",
		})
	case *goast.AssignStmt:
		recordAssignment(x, tree, state, parents)
	case *goast.CommClause:
		// Select receive forms that are not AssignStmt/SendStmt.
		if x.Comm != nil {
			if ue, ok := receiveUnary(x.Comm); ok {
				recordSelectReceiveUnary(ue, tree, state, parents)
			}
		}
	}

	if entered != nil {
		state.pushScope(entered.kind, entered.br)
	}

	// Full-slice append so child walks cannot clobber sibling parent stacks.
	nextParents := append(parents[:len(parents):len(parents)], n)
	forEachChild(n, func(child goast.Node) {
		walkExtract(child, tree, state, nextParents)
	})

	if entered != nil {
		state.popScope()
	}
	if restoreFunc != nil {
		state.currentFunction = *restoreFunc
	}
}

type scopeEntry struct {
	kind ScopeKind
	br   ByteRange
}

func nodeRange(tree *goparse.Tree, n goast.Node) ByteRange {
	if n == nil || tree == nil {
		return ByteRange{}
	}
	return ByteRange{tree.Offset(n.Pos()), tree.Offset(n.End())}
}

// forEachChild visits immediate children of n (not n itself, not grandchildren).
func forEachChild(n goast.Node, fn func(goast.Node)) {
	if n == nil {
		return
	}
	goast.Inspect(n, func(c goast.Node) bool {
		if c == n {
			return true
		}
		if c != nil {
			fn(c)
		}
		return false
	})
}

func functionIdentityDecl(fn *goast.FuncDecl, tree *goparse.Tree) string {
	name := "<anonymous>"
	if fn.Name != nil {
		name = fn.Name.Name
	}
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return name
	}
	recv := NormalizeReceiverType(tree.NodeText(fn.Recv))
	if recv != "" {
		return recv + "." + name
	}
	return name
}

func extractParamNamesFieldList(ft *goast.FuncType) []string {
	if ft == nil || ft.Params == nil {
		return nil
	}
	var out []string
	for _, field := range ft.Params.List {
		if len(field.Names) == 0 {
			continue
		}
		for _, id := range field.Names {
			if id == nil {
				continue
			}
			name := strings.TrimSpace(id.Name)
			if name != "" && name != "_" {
				out = append(out, name)
			}
		}
	}
	return out
}

func recordCall(node *goast.CallExpr, tree *goparse.Tree, state *extractionState, parents []goast.Node) {
	if node == nil || node.Fun == nil {
		return
	}
	if isChainedCall(node.Fun) {
		return
	}
	funcText := strings.TrimSpace(tree.NodeText(node.Fun))
	if funcText == "" {
		return
	}
	br := nodeRange(tree, node)
	args := argumentTexts(node, tree)

	if kind, ok := ClassifySource(funcText); ok {
		rv := resultVariableOfCall(parents, tree)
		state.sources = append(state.sources, TaintSourceAnnotation{
			Function:       funcText,
			Kind:           kind,
			ByteRange:      br,
			ResultVariable: rv,
			Arguments:      args,
		})
		return
	}

	receiver := ""
	if sel, ok := node.Fun.(*goast.SelectorExpr); ok && sel.X != nil {
		receiver = strings.TrimSpace(tree.NodeText(sel.X))
	}
	firstArg := ""
	if len(args) > 0 {
		firstArg = args[0]
	}

	// HTTP write sinks need ResponseWriter check.
	if sk, idx, ok := ClassifySinkHTTPWrite(funcText, looksLikeResponseWriter(parents, tree, state, receiver, firstArg, funcText), strings.HasPrefix(strings.TrimSpace(firstArg), "[]string")); ok {
		argText := ""
		if idx < len(args) {
			argText = args[idx]
		}
		state.sinks = append(state.sinks, TaintSinkAnnotation{
			Function: funcText, Kind: sk, ArgumentIndex: idx,
			ArgumentText: argText, AllArguments: args, ByteRange: br,
		})
		return
	}

	if sk, idx, ok := ClassifySink(funcText, state.source, receiver, firstArg, state.hasHTMLTemplate); ok {
		argText := ""
		if idx < len(args) {
			argText = args[idx]
		}
		state.sinks = append(state.sinks, TaintSinkAnnotation{
			Function: funcText, Kind: sk, ArgumentIndex: idx,
			ArgumentText: argText, AllArguments: args, ByteRange: br,
		})
		return
	}

	if kind, ok := ClassifySanitizer(funcText); ok {
		rv := resultVariableOfCall(parents, tree)
		state.sanitizers = append(state.sanitizers, TaintSanitizerAnnotation{
			Function: funcText, Kind: kind, ByteRange: br,
			ResultVariable: rv, Arguments: args,
		})
	}
}

func looksLikeResponseWriter(parents []goast.Node, tree *goparse.Tree, state *extractionState, receiver, firstArg, funcText string) bool {
	// fmt.Fprintf: check arg0; method Write: check receiver
	name := receiver
	if funcText == "fmt.Fprintf" {
		name = firstArg
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	// Reject csv.Writer.Write([]string)
	if strings.HasPrefix(strings.TrimSpace(firstArg), "[]string") {
		return false
	}
	params := enclosingFunctionParams(parents, tree)
	if params == "" {
		return false
	}
	return strings.Contains(params, name+" http.ResponseWriter") ||
		strings.Contains(params, name+" *http.ResponseWriter")
}

func enclosingFunctionParams(parents []goast.Node, tree *goparse.Tree) string {
	for i := len(parents) - 1; i >= 0; i-- {
		switch p := parents[i].(type) {
		case *goast.FuncDecl:
			if p.Type != nil && p.Type.Params != nil {
				return tree.NodeText(p.Type.Params)
			}
			return ""
		case *goast.FuncLit:
			if p.Type != nil && p.Type.Params != nil {
				return tree.NodeText(p.Type.Params)
			}
			return ""
		}
	}
	return ""
}

func recordAssignment(node *goast.AssignStmt, tree *goparse.Tree, state *extractionState, parents []goast.Node) {
	text := tree.NodeText(node)
	lhs, rhs, ok := SplitAssignment(text)
	if !ok {
		// Fallback from AST parts when text split fails.
		if len(node.Lhs) == 0 {
			return
		}
		var lhsParts []string
		for _, l := range node.Lhs {
			lhsParts = append(lhsParts, strings.TrimSpace(tree.NodeText(l)))
		}
		lhs = strings.Join(lhsParts, ", ")
		var rhsParts []string
		for _, r := range node.Rhs {
			rhsParts = append(rhsParts, strings.TrimSpace(tree.NodeText(r)))
		}
		rhs = strings.Join(rhsParts, ", ")
		if lhs == "" {
			return
		}
	}
	names := ExtractLHSNames(lhs)
	if len(names) == 0 {
		return
	}
	scope := state.currentScope()
	br := nodeRange(tree, node)
	fromCall := IsSourceOrSanitizerCall(rhs)
	recvCh := ChannelFromReceiveRHS(rhs)
	for _, name := range names {
		state.assignments = append(state.assignments, AssignmentDetail{
			LHS:                   name,
			RHSText:               rhs,
			Scope:                 scope,
			ByteRange:             br,
			FromSourceOrSanitizer: fromCall,
			IsChannelSend:         false,
		})
	}
	if recvCh != "" {
		lhsName := ""
		if len(names) > 0 && names[0] != "_" {
			lhsName = names[0]
		}
		state.channelRecvs = append(state.channelRecvs, ChannelRecvSite{
			Channel: recvCh, LHS: lhsName,
			FunctionScope: state.currentFunctionScope(),
			RecvScope:     scope,
			ByteRange:     br,
			InSelect:      isInsideSelect(parents),
		})
	}
}

func recordSend(node *goast.SendStmt, tree *goparse.Tree, state *extractionState, parents []goast.Node) {
	ch := ""
	if node.Chan != nil {
		ch = strings.TrimSpace(tree.NodeText(node.Chan))
	}
	val := ""
	if node.Value != nil {
		val = strings.TrimSpace(tree.NodeText(node.Value))
	}
	if ch == "" {
		return
	}
	state.channelSends = append(state.channelSends, ChannelSendSite{
		Channel: ch, ValueText: val,
		FunctionScope: state.currentFunctionScope(),
		ByteRange:     nodeRange(tree, node),
		InSelect:      isInsideSelect(parents),
	})
}

func receiveUnary(n goast.Node) (*goast.UnaryExpr, bool) {
	switch x := n.(type) {
	case *goast.UnaryExpr:
		if x.Op == token.ARROW {
			return x, true
		}
	case *goast.ExprStmt:
		if ue, ok := x.X.(*goast.UnaryExpr); ok && ue.Op == token.ARROW {
			return ue, true
		}
	}
	return nil, false
}

func recordSelectReceiveUnary(ue *goast.UnaryExpr, tree *goparse.Tree, state *extractionState, parents []goast.Node) {
	if ue == nil || ue.X == nil {
		return
	}
	ch := ChannelFromReceiveRHS("<-" + strings.TrimSpace(tree.NodeText(ue.X)))
	if ch == "" {
		// ChannelFromReceiveRHS expects full "<-ch"; also try operand alone when simple ident.
		ch = strings.TrimSpace(tree.NodeText(ue.X))
		if !IsIdent(ch) {
			return
		}
	}
	state.channelRecvs = append(state.channelRecvs, ChannelRecvSite{
		Channel: ch, LHS: "",
		FunctionScope: state.currentFunctionScope(),
		RecvScope:     state.currentScope(),
		ByteRange:     nodeRange(tree, ue),
		InSelect:      true,
	})
}

func isInsideSelect(parents []goast.Node) bool {
	for i := len(parents) - 1; i >= 0; i-- {
		switch parents[i].(type) {
		case *goast.SelectStmt:
			return true
		case *goast.FuncDecl, *goast.FuncLit:
			return false
		}
	}
	return false
}

func isChainedCall(fun goast.Expr) bool {
	sel, ok := fun.(*goast.SelectorExpr)
	if !ok {
		return false
	}
	_, ok = sel.X.(*goast.CallExpr)
	return ok
}

func resultVariableOfCall(parents []goast.Node, tree *goparse.Tree) string {
	for i := len(parents) - 1; i >= 0; i-- {
		switch p := parents[i].(type) {
		case *goast.AssignStmt:
			if len(p.Lhs) == 0 {
				return ""
			}
			var parts []string
			for _, l := range p.Lhs {
				parts = append(parts, strings.TrimSpace(tree.NodeText(l)))
			}
			return strings.Join(parts, ", ")
		case *goast.SendStmt:
			// Send value must not attribute channel as result.
			return ""
		}
	}
	return ""
}

func argumentTexts(call *goast.CallExpr, tree *goparse.Tree) []string {
	if call == nil {
		return nil
	}
	out := make([]string, 0, len(call.Args))
	for _, a := range call.Args {
		out = append(out, strings.TrimSpace(tree.NodeText(a)))
	}
	return out
}

func pairChannelTransfers(sends []ChannelSendSite, recvs []ChannelRecvSite, assigns []AssignmentDetail) ([]ChannelTransfer, []UnsupportedFlow) {
	type key struct {
		scope int
		ch    string
	}
	sendG := map[key][]ChannelSendSite{}
	recvG := map[key][]ChannelRecvSite{}
	for _, s := range sends {
		k := key{s.FunctionScope, s.Channel}
		sendG[k] = append(sendG[k], s)
	}
	for _, r := range recvs {
		k := key{r.FunctionScope, r.Channel}
		recvG[k] = append(recvG[k], r)
	}
	seen := map[key]struct{}{}
	var keys []key
	for k := range sendG {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			keys = append(keys, k)
		}
	}
	for k := range recvG {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			keys = append(keys, k)
		}
	}

	var transfers []ChannelTransfer
	var unsupported []UnsupportedFlow
	for _, k := range keys {
		ss := sendG[k]
		rs := recvG[k]
		decline := len(ss) != 1 || len(rs) != 1 || ss[0].InSelect || rs[0].InSelect || rs[0].LHS == "" || channelLooksBuffered(assigns, k.ch)
		if decline {
			for _, s := range ss {
				unsupported = append(unsupported, UnsupportedFlow{
					Kind: UnsupportedChannel, ByteRange: s.ByteRange,
					Note: "channel send/receive declined by G5 v0 pairing (explicit FN)",
				})
			}
			for _, r := range rs {
				unsupported = append(unsupported, UnsupportedFlow{
					Kind: UnsupportedChannel, ByteRange: r.ByteRange,
					Note: "channel send/receive declined by G5 v0 pairing (explicit FN)",
				})
			}
			continue
		}
		transfers = append(transfers, ChannelTransfer{
			Channel: k.ch, SendValueText: ss[0].ValueText, RecvLHS: rs[0].LHS,
			RecvScope: rs[0].RecvScope, SendByteRange: ss[0].ByteRange, RecvByteRange: rs[0].ByteRange,
		})
	}
	return transfers, unsupported
}

func channelLooksBuffered(assigns []AssignmentDetail, channel string) bool {
	for _, a := range assigns {
		if a.LHS != channel {
			continue
		}
		rhs := strings.TrimSpace(a.RHSText)
		if !strings.HasPrefix(rhs, "make(chan") {
			continue
		}
		if open := strings.Index(rhs, "("); open >= 0 {
			inner := rhs[open+1:]
			if close := strings.LastIndex(inner, ")"); close >= 0 {
				if strings.Contains(inner[:close], ",") {
					return true
				}
			}
		}
	}
	return false
}

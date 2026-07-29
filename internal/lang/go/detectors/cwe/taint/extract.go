package taint

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// ExtractTaintFacts walks the unit AST and collects sources/sinks/sanitizers/assignments.
func ExtractTaintFacts(unit *core.ParsedUnit) TaintAnnotations {
	if unit == nil {
		return TaintAnnotations{}
	}
	src := []byte(unit.Source)
	state := newExtractionState(unit.Source)
	state.pushScope(ScopePackage, ByteRange{0, len(unit.Source)})

	root := unitRoot(unit)
	if root != nil {
		walkExtract(root, src, state)
	} else {
		// Source-only fallback: no tree → empty annotations.
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

func unitRoot(unit *core.ParsedUnit) *sitter.Node {
	if unit == nil || unit.Tree == nil {
		return nil
	}
	type rooted interface{ RootNode() *sitter.Node }
	if t, ok := unit.Tree.(rooted); ok {
		return t.RootNode()
	}
	return nil
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

func walkExtract(node *sitter.Node, src []byte, state *extractionState) {
	if node == nil {
		return
	}
	var entered *scopeEntry
	var restoreFunc *string

	kind := node.Kind()
	switch kind {
	case "function_declaration", "func_literal", "method_declaration":
		name := functionIdentity(node, src)
		prev := state.currentFunction
		restoreFunc = &prev
		state.currentFunction = name
		state.functionParams[name] = extractParamNames(node, src)
		state.functionRanges[name] = ByteRange{int(node.StartByte()), int(node.EndByte())}
		entered = &scopeEntry{ScopeFunction, ByteRange{int(node.StartByte()), int(node.EndByte())}}
	case "block":
		entered = &scopeEntry{ScopeBlock, ByteRange{int(node.StartByte()), int(node.EndByte())}}
	case "if_statement":
		entered = &scopeEntry{ScopeIf, ByteRange{int(node.StartByte()), int(node.EndByte())}}
	case "for_statement", "range_clause":
		entered = &scopeEntry{ScopeFor, ByteRange{int(node.StartByte()), int(node.EndByte())}}
	case "switch_statement", "expression_switch_statement", "type_switch_statement":
		entered = &scopeEntry{ScopeSwitch, ByteRange{int(node.StartByte()), int(node.EndByte())}}
	case "case_clause", "default_case":
		entered = &scopeEntry{ScopeCase, ByteRange{int(node.StartByte()), int(node.EndByte())}}
	case "call_expression":
		recordCall(node, state)
	case "send_statement":
		recordSend(node, state)
	case "receive_statement":
		recordSelectReceive(node, state)
	case "go_statement":
		state.unsupported = append(state.unsupported, UnsupportedFlow{
			Kind:      UnsupportedGoroutine,
			ByteRange: ByteRange{int(node.StartByte()), int(node.EndByte())},
			Note:      "goroutine spawn is not tracked by taint (explicit FN)",
		})
	case "assignment_statement", "short_var_declaration":
		recordAssignment(node, state)
	}

	if entered != nil {
		state.pushScope(entered.kind, entered.br)
	}

	// Manual child walk (pre-order with scopes).
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child == nil || child.IsNamed() == false && !isInterestingUnnamed(child) {
			// Still walk all children; tree-sitter Go fields are on named nodes.
			// Walk named children primarily.
		}
		if child != nil {
			walkExtract(child, src, state)
		}
	}

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

func isInterestingUnnamed(n *sitter.Node) bool { return false }

func functionIdentity(node *sitter.Node, src []byte) string {
	name := "<anonymous>"
	if n := node.ChildByFieldName("name"); n != nil {
		if t := strings.TrimSpace(n.Utf8Text(src)); t != "" {
			name = t
		}
	}
	if node.Kind() != "method_declaration" {
		return name
	}
	recv := ""
	if r := node.ChildByFieldName("receiver"); r != nil {
		recv = NormalizeReceiverType(r.Utf8Text(src))
	}
	if recv != "" {
		return recv + "." + name
	}
	return name
}

func extractParamNames(node *sitter.Node, src []byte) []string {
	params := node.ChildByFieldName("parameters")
	if params == nil {
		return nil
	}
	var out []string
	cursor := params.Walk()
	defer cursor.Close()
	for _, p := range params.NamedChildren(cursor) {
		// parameter_declaration may have multiple name identifiers.
		if n := p.ChildByFieldName("name"); n != nil {
			name := strings.TrimSpace(n.Utf8Text(src))
			if name != "" && name != "_" {
				out = append(out, name)
			}
			continue
		}
		// Walk identifiers in the declaration.
		pc := p.Walk()
		for _, ch := range p.NamedChildren(pc) {
			if ch.Kind() == "identifier" {
				name := strings.TrimSpace(ch.Utf8Text(src))
				if name != "" && name != "_" {
					out = append(out, name)
				}
			}
		}
		pc.Close()
	}
	return out
}

func recordCall(node *sitter.Node, state *extractionState) {
	fn := node.ChildByFieldName("function")
	if fn == nil {
		return
	}
	if isChainedCall(fn) {
		return
	}
	funcText := strings.TrimSpace(fn.Utf8Text(state.srcBytes))
	if funcText == "" {
		return
	}
	br := ByteRange{int(node.StartByte()), int(node.EndByte())}
	args := argumentTexts(node, state.srcBytes)

	if kind, ok := ClassifySource(funcText); ok {
		rv := resultVariableOfCall(node, state.srcBytes)
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
	if fn.Kind() == "selector_expression" {
		if op := fn.ChildByFieldName("operand"); op != nil {
			receiver = strings.TrimSpace(op.Utf8Text(state.srcBytes))
		}
	}
	firstArg := ""
	if len(args) > 0 {
		firstArg = args[0]
	}

	// HTTP write sinks need ResponseWriter check.
	if sk, idx, ok := ClassifySinkHTTPWrite(funcText, looksLikeResponseWriter(node, state, receiver, firstArg, funcText), strings.HasPrefix(strings.TrimSpace(firstArg), "[]string")); ok {
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
		rv := resultVariableOfCall(node, state.srcBytes)
		state.sanitizers = append(state.sanitizers, TaintSanitizerAnnotation{
			Function: funcText, Kind: kind, ByteRange: br,
			ResultVariable: rv, Arguments: args,
		})
	}
}

func looksLikeResponseWriter(call *sitter.Node, state *extractionState, receiver, firstArg, funcText string) bool {
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
	// Enclosing function parameters must declare name as http.ResponseWriter
	fn := enclosingFunction(call)
	if fn == nil {
		return false
	}
	params := fn.ChildByFieldName("parameters")
	if params == nil {
		return false
	}
	ptext := params.Utf8Text(state.srcBytes)
	return strings.Contains(ptext, name+" http.ResponseWriter") ||
		strings.Contains(ptext, name+" *http.ResponseWriter")
}

func enclosingFunction(node *sitter.Node) *sitter.Node {
	for p := node.Parent(); p != nil; p = p.Parent() {
		switch p.Kind() {
		case "function_declaration", "method_declaration", "func_literal":
			return p
		}
	}
	return nil
}

func recordAssignment(node *sitter.Node, state *extractionState) {
	text := node.Utf8Text(state.srcBytes)
	lhs, rhs, ok := SplitAssignment(text)
	if !ok {
		return
	}
	names := ExtractLHSNames(lhs)
	if len(names) == 0 {
		return
	}
	scope := state.currentScope()
	br := ByteRange{int(node.StartByte()), int(node.EndByte())}
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
			InSelect:      isInsideSelect(node),
		})
	}
}

func recordSend(node *sitter.Node, state *extractionState) {
	ch := ""
	if n := node.ChildByFieldName("channel"); n != nil {
		ch = strings.TrimSpace(n.Utf8Text(state.srcBytes))
	}
	val := ""
	if n := node.ChildByFieldName("value"); n != nil {
		val = strings.TrimSpace(n.Utf8Text(state.srcBytes))
	}
	if ch == "" {
		return
	}
	state.channelSends = append(state.channelSends, ChannelSendSite{
		Channel: ch, ValueText: val,
		FunctionScope: state.currentFunctionScope(),
		ByteRange:     ByteRange{int(node.StartByte()), int(node.EndByte())},
		InSelect:      isInsideSelect(node),
	})
}

func recordSelectReceive(node *sitter.Node, state *extractionState) {
	right := node.ChildByFieldName("right")
	if right == nil {
		return
	}
	ch := ChannelFromReceiveRHS(strings.TrimSpace(right.Utf8Text(state.srcBytes)))
	if ch == "" {
		return
	}
	lhs := ""
	if left := node.ChildByFieldName("left"); left != nil {
		names := ExtractLHSNames(left.Utf8Text(state.srcBytes))
		if len(names) > 0 {
			lhs = names[0]
		}
	}
	state.channelRecvs = append(state.channelRecvs, ChannelRecvSite{
		Channel: ch, LHS: lhs,
		FunctionScope: state.currentFunctionScope(),
		RecvScope:     state.currentScope(),
		ByteRange:     ByteRange{int(node.StartByte()), int(node.EndByte())},
		InSelect:      true,
	})
}

func isInsideSelect(node *sitter.Node) bool {
	for p := node.Parent(); p != nil; p = p.Parent() {
		if p.Kind() == "select_statement" {
			return true
		}
		switch p.Kind() {
		case "function_declaration", "method_declaration", "func_literal":
			return false
		}
	}
	return false
}

func isChainedCall(funcNode *sitter.Node) bool {
	if funcNode.Kind() != "selector_expression" {
		return false
	}
	op := funcNode.ChildByFieldName("operand")
	return op != nil && op.Kind() == "call_expression"
}

func resultVariableOfCall(call *sitter.Node, src []byte) string {
	parent := call.Parent()
	for parent != nil {
		k := parent.Kind()
		if k == "assignment_statement" || k == "short_var_declaration" || k == "send_statement" {
			break
		}
		parent = parent.Parent()
	}
	if parent == nil {
		return ""
	}
	// Send value must not attribute channel as result.
	if parent.Kind() == "send_statement" {
		return ""
	}
	left := parent.ChildByFieldName("left")
	if left == nil {
		return ""
	}
	return strings.TrimSpace(left.Utf8Text(src))
}

func argumentTexts(call *sitter.Node, src []byte) []string {
	args := call.ChildByFieldName("arguments")
	if args == nil {
		return nil
	}
	cursor := args.Walk()
	defer cursor.Close()
	var out []string
	for _, n := range args.NamedChildren(cursor) {
		out = append(out, strings.TrimSpace(n.Utf8Text(src)))
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

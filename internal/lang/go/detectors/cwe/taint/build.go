package taint

import (
	"sort"
	"strings"
)

// BuildTaintGraph constructs an intra-procedural taint graph from annotations.
func BuildTaintGraph(ann *TaintAnnotations) *TaintGraph {
	g := NewTaintGraph()
	if ann == nil {
		return g
	}

	// decl_nodes[scope][name] = sorted (decl_byte, node_id)
	type version struct {
		byte int
		id   int
	}
	declNodes := map[int]map[string][]version{}

	scopeByID := map[int]*ScopeInfo{}
	for i := range ann.Scopes {
		s := &ann.Scopes[i]
		scopeByID[s.ID] = s
	}
	scopeOrder := make([]int, 0, len(ann.Scopes))
	for _, s := range ann.Scopes {
		scopeOrder = append(scopeOrder, s.ID)
	}
	sort.Slice(scopeOrder, func(i, j int) bool {
		a, b := scopeByID[scopeOrder[i]], scopeByID[scopeOrder[j]]
		da := a.ByteRange.End - a.ByteRange.Start
		db := b.ByteRange.End - b.ByteRange.Start
		if da != db {
			return da < db
		}
		if a.ByteRange.Start != b.ByteRange.Start {
			return a.ByteRange.Start < b.ByteRange.Start
		}
		return a.ID < b.ID
	})

	addDecl := func(scope int, name string, byteOff, id int) {
		if declNodes[scope] == nil {
			declNodes[scope] = map[string][]version{}
		}
		declNodes[scope][name] = append(declNodes[scope][name], version{byteOff, id})
	}

	// Parameters as variables.
	for funcName, params := range ann.FunctionParams {
		var funcScope *ScopeInfo
		for i := range ann.Scopes {
			s := &ann.Scopes[i]
			if s.Kind == ScopeFunction && s.Function == funcName {
				funcScope = s
				break
			}
		}
		if funcScope == nil {
			continue
		}
		for _, param := range params {
			id := g.AddNode(TaintNode{
				Kind: "variable", Name: param, Scope: funcScope.ID,
				DeclByte: funcScope.ByteRange.Start,
			})
			addDecl(funcScope.ID, param, funcScope.ByteRange.Start, id)
		}
	}

	// Assignment variables (versioned).
	for _, a := range ann.Assignments {
		if a.IsChannelSend {
			continue
		}
		id := g.AddNode(TaintNode{
			Kind: "variable", Name: a.LHS, Scope: a.Scope, DeclByte: a.ByteRange.Start,
		})
		addDecl(a.Scope, a.LHS, a.ByteRange.Start, id)
	}

	// Sort versions.
	for _, names := range declNodes {
		for n := range names {
			vs := names[n]
			sort.Slice(vs, func(i, j int) bool { return vs[i].byte < vs[j].byte })
			names[n] = vs
		}
	}

	resolveDeclAt := func(scope int, name string, useByte int) (int, bool) {
		names := declNodes[scope]
		if names == nil {
			return 0, false
		}
		vs := names[name]
		for i := len(vs) - 1; i >= 0; i-- {
			if vs[i].byte <= useByte {
				return vs[i].id, true
			}
		}
		return 0, false
	}

	resolveVariable := func(byteOff int, name string) (int, bool) {
		// Innermost scope containing byteOff (scopeOrder sorted by size ascending).
		var current *ScopeInfo
		for _, id := range scopeOrder {
			s := scopeByID[id]
			if s.ByteRange.Start <= byteOff && byteOff < s.ByteRange.End {
				current = s
				break
			}
		}
		if current == nil {
			return 0, false
		}
		for {
			if id, ok := resolveDeclAt(current.ID, name, byteOff); ok {
				return id, true
			}
			if base, _, ok := strings.Cut(name, "."); ok {
				if id, ok := resolveDeclAt(current.ID, base, byteOff); ok {
					return id, true
				}
			}
			if current.Parent == nil {
				return 0, false
			}
			current = scopeByID[*current.Parent]
			if current == nil {
				return 0, false
			}
		}
	}

	wireArgs := func(nodeID, byteOff int, args []string) {
		for idx, arg := range args {
			for _, name := range ReferencedNames(arg) {
				if srcID, ok := resolveVariable(byteOff, name); ok {
					g.AddEdge(srcID, nodeID, Argument(idx))
				}
			}
		}
	}

	// Sources.
	for _, src := range ann.Sources {
		id := g.AddNode(TaintNode{
			Kind: "source", Function: src.Function, SourceKind: src.Kind, ByteRange: src.ByteRange,
		})
		if src.ResultVariable != "" {
			if target, ok := resolveVariable(src.ByteRange.Start, src.ResultVariable); ok {
				g.AddEdge(id, target, EdgeAssignment)
			}
		}
		wireArgs(id, src.ByteRange.Start, src.Arguments)
	}

	// Sanitizers.
	for _, san := range ann.Sanitizers {
		id := g.AddNode(TaintNode{
			Kind: "sanitizer", Function: san.Function, SanitizerKind: san.Kind, ByteRange: san.ByteRange,
		})
		if san.ResultVariable != "" {
			if target, ok := resolveVariable(san.ByteRange.Start, san.ResultVariable); ok {
				g.AddEdge(id, target, EdgeAssignment)
			}
		}
		wireArgs(id, san.ByteRange.Start, san.Arguments)
	}

	// Sinks.
	for _, sink := range ann.Sinks {
		id := g.AddNode(TaintNode{
			Kind: "sink", Function: sink.Function, SinkKind: sink.Kind,
			ArgumentIndex: sink.ArgumentIndex, ByteRange: sink.ByteRange,
		})
		for idx, arg := range sink.AllArguments {
			for _, name := range ReferencedNames(arg) {
				if srcID, ok := resolveVariable(sink.ByteRange.Start, name); ok {
					g.AddEdge(srcID, id, Argument(idx))
				}
			}
		}
		// Deserialization pointer bridge.
		for _, outIdx := range taintedOutputArgs(sink.Function) {
			if outIdx >= len(sink.AllArguments) {
				continue
			}
			outNames := ReferencedNames(sink.AllArguments[outIdx])
			for inIdx, inArg := range sink.AllArguments {
				if inIdx == outIdx {
					continue
				}
				for _, inName := range ReferencedNames(inArg) {
					for _, outName := range outNames {
						srcID, sok := resolveVariable(sink.ByteRange.Start, inName)
						dstID, dok := resolveVariable(sink.ByteRange.Start, outName)
						if sok && dok {
							g.AddEdge(srcID, dstID, EdgeAssignment)
						}
					}
				}
			}
		}
	}

	// Assignment edges.
	for _, a := range ann.Assignments {
		if a.IsChannelSend {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(a.RHSText), "<-") {
			continue
		}
		target, ok := resolveDeclAt(a.Scope, a.LHS, a.ByteRange.Start)
		if !ok {
			continue
		}
		if a.FromSourceOrSanitizer {
			continue
		}
		callName := strings.TrimSpace(strings.SplitN(a.RHSText, "(", 2)[0])
		isOpaque := strings.Contains(a.RHSText, "(") &&
			!isSourceOrSanitizerAssignment(a.RHSText) &&
			!isKnownPropagator(callName)
		if isOpaque {
			continue
		}
		for _, name := range ReferencedNames(a.RHSText) {
			if srcID, ok := resolveVariable(a.ByteRange.Start, name); ok {
				g.AddEdge(srcID, target, EdgeAssignment)
			}
		}
	}

	// Channel transfers.
	for _, t := range ann.ChannelTransfers {
		target, ok := resolveDeclAt(t.RecvScope, t.RecvLHS, t.RecvByteRange.Start)
		if !ok {
			continue
		}
		if isSourceOrSanitizerAssignment(t.SendValueText) {
			if srcID, ok := findSourceOrSanitizerAt(g, t.SendByteRange.Start, t.SendByteRange.End); ok {
				g.AddEdge(srcID, target, EdgeChannelTransfer)
			}
		} else {
			for _, name := range ReferencedNames(t.SendValueText) {
				if srcID, ok := resolveVariable(t.SendByteRange.Start, name); ok {
					g.AddEdge(srcID, target, EdgeChannelTransfer)
				}
			}
		}
	}

	// Map/slice index base taint.
	for _, a := range ann.Assignments {
		if a.IsChannelSend || strings.HasPrefix(strings.TrimSpace(a.RHSText), "<-") {
			continue
		}
		bracket := strings.Index(a.LHS, "[")
		if bracket < 0 {
			continue
		}
		base := strings.TrimSpace(a.LHS[:bracket])
		if baseID, ok := resolveDeclAt(a.Scope, base, a.ByteRange.Start); ok && !a.FromSourceOrSanitizer {
			for _, name := range ReferencedNames(a.RHSText) {
				if srcID, ok := resolveVariable(a.ByteRange.Start, name); ok {
					g.AddEdge(srcID, baseID, EdgeAssignment)
				}
			}
		}
	}

	return g
}

func findSourceOrSanitizerAt(g *TaintGraph, start, end int) (int, bool) {
	for id, n := range g.Nodes {
		if n.Kind == "source" || n.Kind == "sanitizer" {
			if n.ByteRange.Start >= start && n.ByteRange.End <= end {
				return id, true
			}
		}
	}
	return 0, false
}

func isSourceOrSanitizerAssignment(rhs string) bool {
	callName := strings.TrimSpace(strings.SplitN(rhs, "(", 2)[0])
	if callName == "" {
		return false
	}
	if strings.Contains(callName, ".URL.Query") ||
		strings.Contains(callName, ".FormValue") ||
		strings.Contains(callName, ".PostForm") ||
		strings.Contains(callName, ".Header.Get") ||
		strings.Contains(callName, ".GetRawData") ||
		strings.HasSuffix(callName, ".PathValue") ||
		strings.HasSuffix(callName, ".Param") ||
		callName == "c.Query" || callName == "c.DefaultQuery" || callName == "c.QueryArray" ||
		callName == "os.Args" || callName == "flag.Args" || callName == "flag.String" ||
		callName == "os.Getenv" || callName == "os.LookupEnv" || callName == "io.ReadAll" {
		return true
	}
	if callName == "filepath.Base" || callName == "html.EscapeString" ||
		callName == "url.QueryEscape" || callName == "url.PathEscape" ||
		callName == "ldap.EscapeFilter" || callName == "xml.EscapeText" || callName == "xml.Marshal" {
		return true
	}
	name := callName
	if i := strings.LastIndex(callName, "."); i >= 0 {
		name = callName[i+1:]
	}
	lower := strings.ToLower(name)
	return strings.HasPrefix(lower, "sanitize") || strings.HasPrefix(lower, "escape") ||
		strings.HasPrefix(lower, "validate") || strings.HasPrefix(lower, "purify")
}

func isKnownPropagator(funcName string) bool {
	switch funcName {
	case "filepath.Join", "filepath.Clean", "strings.Join", "strings.Replace", "strings.Repeat",
		"strings.Trim", "strings.TrimSpace", "fmt.Sprintf", "fmt.Errorf", "path.Join", "path.Clean",
		"append", "json.Marshal", "strconv.Itoa", "strconv.FormatInt", "html.UnescapeString":
		return true
	}
	return false
}

func taintedOutputArgs(funcName string) []int {
	if funcName == "json.Unmarshal" || funcName == "xml.Unmarshal" {
		return []int{0}
	}
	return nil
}

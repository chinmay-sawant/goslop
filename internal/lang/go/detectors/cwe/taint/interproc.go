package taint

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

// ProjectUnit is accumulated per-file state for project finalize.
type ProjectUnit struct {
	Path       string
	Source     string
	LineStarts []int
	Package    PackageIdentity
	CallGraph  *CallGraph
	Annot      TaintAnnotations
	Graph      *TaintGraph
	ImportMap  map[string]string
}

// MetaSet holds rule metadata for inter-procedural emission.
type MetaSet struct {
	CWE22 *rules.RuleMetadata
	CWE78 *rules.RuleMetadata
	CWE79 *rules.RuleMetadata
	CWE89 *rules.RuleMetadata
}

type fileData struct {
	path      string
	graph     *TaintGraph
	index     *TaintGraphIndex
	summaries map[string]*TaintSummary
}

// FinalizeInterProcedural runs bounded same-package inter-procedural analysis.
func FinalizeInterProcedural(units []ProjectUnit, ctx *core.ScanContext, meta MetaSet, out *[]rules.Finding) {
	if len(units) == 0 || ctx == nil || !ctx.TaintEnabled {
		return
	}
	sort.Slice(units, func(i, j int) bool { return units[i].Path < units[j].Path })

	maxDepth := ctx.EffectiveTaintMaxDepth()
	perFile := make([]fileData, 0, len(units))
	declIndex := map[TaintSymbolKey]int{}
	packageNameIndex := map[pkgNameKey][]TaintSymbolKey{}

	for idx := range units {
		u := &units[idx]
		g := u.Graph
		if g == nil {
			g = BuildTaintGraph(&u.Annot)
		}
		index := BuildIndex(g)
		summaries := ComputeAllSummariesWithIndex(g, index, &u.Annot, u.Source)
		if maxDepth > 1 && u.CallGraph != nil {
			RefineSummariesMultihopWithContext(u.CallGraph, &u.Annot, summaries, maxDepth)
		}
		perFile = append(perFile, fileData{path: u.Path, graph: g, index: index, summaries: summaries})
		if u.CallGraph == nil {
			continue
		}
		for _, decl := range u.CallGraph.Declarations {
			var key TaintSymbolKey
			if decl.IsMethod {
				key = MethodKey(u.Package, decl.ReceiverType, decl.Name)
			} else {
				key = FunctionKey(u.Package, decl.Name)
			}
			if _, ok := declIndex[key]; !ok {
				declIndex[key] = idx
			}
			pk := pkgNameKey{u.Package, decl.Name}
			packageNameIndex[pk] = append(packageNameIndex[pk], key)
		}
	}

	varIndexes := make([]map[string][]int, len(perFile))
	for i, f := range perFile {
		varIndexes[i] = buildVariableIndex(f.graph)
	}

	summaryIndex := map[TaintSymbolKey]int{}
	for fileIdx, u := range units {
		for name := range perFile[fileIdx].summaries {
			if u.CallGraph != nil {
				if decl, ok := u.CallGraph.Declarations[name]; ok {
					var key TaintSymbolKey
					if decl.IsMethod {
						key = MethodKey(u.Package, decl.ReceiverType, decl.Name)
					} else {
						key = FunctionKey(u.Package, decl.Name)
					}
					if _, seen := summaryIndex[key]; !seen {
						summaryIndex[key] = fileIdx
					}
					continue
				}
			}
			// Fall back: name is free function or *Type.Method identity.
			if strings.Contains(name, ".") {
				parts := strings.SplitN(name, ".", 2)
				key := MethodKey(u.Package, parts[0], parts[1])
				if _, seen := summaryIndex[key]; !seen {
					summaryIndex[key] = fileIdx
				}
			} else {
				key := FunctionKey(u.Package, name)
				if _, seen := summaryIndex[key]; !seen {
					summaryIndex[key] = fileIdx
				}
			}
		}
	}

	for callerIdx, u := range units {
		if u.CallGraph == nil {
			continue
		}
		callerGraph := perFile[callerIdx].graph
		callerIndex := perFile[callerIdx].index
		callerVars := varIndexes[callerIdx]
		callerImports := map[string]struct{}{}
		for k := range u.ImportMap {
			callerImports[k] = struct{}{}
		}

		for si := range u.CallGraph.Sites {
			site := &u.CallGraph.Sites[si]
			rawCallee := site.Callee
			if site.IsMethodCall {
				if dot := strings.LastIndex(rawCallee, "."); dot >= 0 {
					prefix := rawCallee[:dot]
					if _, ok := callerImports[prefix]; ok {
						continue
					}
				}
			}
			calleeName := resolveCalleeName(rawCallee, site.IsMethodCall)
			var callerDecl *FunctionDecl
			if d, ok := u.CallGraph.Declarations[site.Caller]; ok {
				callerDecl = &d
			}
			calleeSum := findSamePackageSummary(
				perFile, summaryIndex, packageNameIndex, declIndex,
				u.Package, rawCallee, calleeName, site.IsMethodCall, callerDecl,
			)
			if calleeSum == nil {
				continue
			}

			for i, isSrc := range calleeSum.ParamSources {
				if !isSrc || i >= len(site.Arguments) {
					continue
				}
				arg := site.Arguments[i]
				if !isIdentifierTainted(callerGraph, callerIndex, callerVars, arg) {
					continue
				}
				emitInterProc(u.Path, u.LineStarts, site, calleeSum.SinkKinds, arg, ctx, meta, out)
			}

			for retIdx, isRet := range calleeSum.ReturnSources {
				if !isRet {
					continue
				}
				resultVar := ResultVariableAtReturnIndex(site.AssignmentLHS, retIdx)
				if resultVar == "" {
					continue
				}
				reached := sinkKindsReachedByVar(callerGraph, callerIndex, callerVars, resultVar)
				if len(reached) == 0 {
					continue
				}
				emitInterProc(u.Path, u.LineStarts, site, reached, resultVar, ctx, meta, out)
			}

			for _, outIdx := range calleeSum.OutputPointerParams {
				if outIdx >= len(site.Arguments) {
					continue
				}
				arg := strings.TrimSpace(site.Arguments[outIdx])
				varName := strings.TrimSpace(strings.TrimPrefix(arg, "&"))
				reached := sinkKindsReachedByVar(callerGraph, callerIndex, callerVars, varName)
				if len(reached) == 0 {
					continue
				}
				emitInterProc(u.Path, u.LineStarts, site, reached, varName, ctx, meta, out)
			}
		}
	}
}

func buildVariableIndex(g *TaintGraph) map[string][]int {
	out := map[string][]int{}
	if g == nil {
		return out
	}
	for id, n := range g.Nodes {
		if n.Kind == "variable" {
			out[n.Name] = append(out[n.Name], id)
		}
	}
	return out
}

func resolveCalleeName(callee string, isMethod bool) string {
	if isMethod {
		if dot := strings.LastIndex(callee, "."); dot >= 0 {
			return callee[dot+1:]
		}
	}
	return callee
}

type pkgNameKey struct {
	pkg  PackageIdentity
	name string
}

func findSamePackageSummary(
	perFile []fileData,
	summaryIndex map[TaintSymbolKey]int,
	packageNameIndex map[pkgNameKey][]TaintSymbolKey,
	declIndex map[TaintSymbolKey]int,
	callerPackage PackageIdentity,
	rawCallee, bareName string,
	isMethodCall bool,
	callerDecl *FunctionDecl,
) *TaintSummary {
	// Adapt packageNameIndex type from Finalize — we used pkgName which equals pkgNameKey fields.
	// Rebuild lookup using bare name in same package.
	bare := bareName

	if !isMethodCall {
		fnKey := FunctionKey(callerPackage, bare)
		if s := summaryForKey(perFile, summaryIndex, fnKey); s != nil {
			return s
		}
		if rawCallee != bare {
			rawKey := FunctionKey(callerPackage, rawCallee)
			if s := summaryForKey(perFile, summaryIndex, rawKey); s != nil {
				return s
			}
		}
		return nil
	}

	// Collect candidates with this bare method name in the package.
	var candidates []TaintSymbolKey
	for key := range summaryIndex {
		if key.Package == callerPackage && key.Name == bare {
			candidates = append(candidates, key)
		}
	}
	for key := range declIndex {
		if key.Package == callerPackage && key.Name == bare {
			found := false
			for _, c := range candidates {
				if c == key {
					found = true
					break
				}
			}
			if !found {
				candidates = append(candidates, key)
			}
		}
	}
	// Also use packageNameIndex if populated with matching type.
	if packageNameIndex != nil {
		if ks, ok := packageNameIndex[pkgNameKey{callerPackage, bare}]; ok {
			candidates = append(candidates, ks...)
		}
	}
	// Dedup
	seen := map[TaintSymbolKey]struct{}{}
	var uniq []TaintSymbolKey
	for _, k := range candidates {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		uniq = append(uniq, k)
	}
	candidates = uniq

	if inferred := inferCallSiteReceiverType(rawCallee, callerDecl); inferred != "" {
		exact := MethodKey(callerPackage, inferred, bare)
		if s := summaryForKey(perFile, summaryIndex, exact); s != nil {
			return s
		}
		if s := summaryFromDecl(perFile, declIndex, exact); s != nil {
			return s
		}
		return nil
	}

	// Unique receiver only.
	var withSummary []TaintSymbolKey
	for _, key := range candidates {
		if summaryForKey(perFile, summaryIndex, key) != nil || summaryFromDecl(perFile, declIndex, key) != nil {
			withSummary = append(withSummary, key)
		}
	}
	if len(withSummary) == 0 {
		return nil
	}
	firstRecv := withSummary[0].Receiver
	for _, k := range withSummary {
		if k.Receiver != firstRecv {
			return nil // ambiguous
		}
	}
	key := withSummary[0]
	if s := summaryForKey(perFile, summaryIndex, key); s != nil {
		return s
	}
	return summaryFromDecl(perFile, declIndex, key)
}

func inferCallSiteReceiverType(rawCallee string, callerDecl *FunctionDecl) string {
	dot := strings.LastIndex(rawCallee, ".")
	if dot < 0 {
		return ""
	}
	recvExpr := strings.TrimSpace(rawCallee[:dot])
	if recvExpr == "" || !IsIdent(recvExpr) || callerDecl == nil || callerDecl.ReceiverType == "" {
		return ""
	}
	normalized := NormalizeReceiverType(callerDecl.ReceiverType)
	if normalized == "" {
		return ""
	}
	s := strings.TrimSpace(callerDecl.ReceiverType)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	s = strings.TrimSpace(s)
	fields := strings.Fields(s)
	// "h *Handler" → first token is param name
	if len(fields) >= 2 && fields[0] == recvExpr {
		return normalized
	}
	return ""
}

func summaryForKey(perFile []fileData, summaryIndex map[TaintSymbolKey]int, key TaintSymbolKey) *TaintSummary {
	fileIdx, ok := summaryIndex[key]
	if !ok || fileIdx < 0 || fileIdx >= len(perFile) {
		return nil
	}
	name := SummaryNameForKey(key)
	return perFile[fileIdx].summaries[name]
}

func summaryFromDecl(perFile []fileData, declIndex map[TaintSymbolKey]int, key TaintSymbolKey) *TaintSummary {
	fileIdx, ok := declIndex[key]
	if !ok || fileIdx < 0 || fileIdx >= len(perFile) {
		return nil
	}
	name := SummaryNameForKey(key)
	return perFile[fileIdx].summaries[name]
}

func isIdentifierTainted(g *TaintGraph, idx *TaintGraphIndex, varIndex map[string][]int, name string) bool {
	varIDs := varIndex[name]
	if len(varIDs) > 0 {
		for _, sourceIDs := range g.BySource {
			for _, sid := range sourceIDs {
				if UnsanitizedReachesAnyWithIndex(g, idx, sid, varIDs) {
					return true
				}
			}
		}
	}
	// Direct source call expression as argument.
	callFunc := strings.TrimSpace(strings.SplitN(name, "(", 2)[0])
	if callFunc != "" {
		for _, sourceIDs := range g.BySource {
			for _, sid := range sourceIDs {
				if sid < len(g.Nodes) && g.Nodes[sid].Kind == "source" && g.Nodes[sid].Function == callFunc {
					return true
				}
			}
		}
	}
	return false
}

func sinkKindsReachedByVar(g *TaintGraph, idx *TaintGraphIndex, varIndex map[string][]int, varName string) []SinkKind {
	varIDs := varIndex[varName]
	if len(varIDs) == 0 {
		return nil
	}
	var reached []SinkKind
	for _, sk := range []SinkKind{
		SinkFileOpen, SinkCommandExec, SinkSQLQuery, SinkTemplate,
		SinkHTTPWrite, SinkLDAPQuery, SinkXMLQuery, SinkDeserialization,
	} {
		if sinks := g.BySink[sk]; len(sinks) > 0 {
			if ForwardReachesAnyWithIndex(g, idx, varIDs, sinks) {
				reached = append(reached, sk)
			}
		}
	}
	return reached
}

func emitInterProc(
	file string,
	lineStarts []int,
	site *CallSite,
	sinkKinds []SinkKind,
	argText string,
	ctx *core.ScanContext,
	meta MetaSet,
	out *[]rules.Finding,
) {
	line, col := lineCol(lineStarts, site.ByteRange.Start)
	for _, sk := range sinkKinds {
		m := metaForSink(sk, meta)
		if m == nil {
			continue
		}
		if ctx != nil && !ctx.Allows(m.ID) {
			continue
		}
		msg := fmt.Sprintf("tainted data reaches %s via call crossing function boundary", m.Title)
		rules.PushFindingWithConfidence(m, file, line, col, msg, 0.7, out)
	}
}

func metaForSink(sk SinkKind, meta MetaSet) *rules.RuleMetadata {
	switch sk {
	case SinkFileOpen:
		return meta.CWE22
	case SinkCommandExec:
		return meta.CWE78
	case SinkTemplate, SinkHTTPWrite:
		return meta.CWE79
	case SinkSQLQuery:
		return meta.CWE89
	default:
		return nil
	}
}

func lineCol(lineStarts []int, byteOffset int) (line, col int) {
	if len(lineStarts) == 0 {
		return 1, 1
	}
	if byteOffset < 0 {
		byteOffset = 0
	}
	lo, hi := 0, len(lineStarts)
	for lo < hi {
		mid := (lo + hi) / 2
		if lineStarts[mid] <= byteOffset {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	idx := lo - 1
	if idx < 0 {
		return 1, 1
	}
	return idx + 1, byteOffset - lineStarts[idx] + 1
}

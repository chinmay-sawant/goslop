package taint

import "strings"

// ComputeAllSummaries builds per-function taint summaries.
func ComputeAllSummaries(g *TaintGraph, ann *TaintAnnotations, source string) map[string]*TaintSummary {
	idx := BuildIndex(g)
	return ComputeAllSummariesWithIndex(g, idx, ann, source)
}

// ComputeAllSummariesWithIndex uses a prebuilt index.
func ComputeAllSummariesWithIndex(g *TaintGraph, idx *TaintGraphIndex, ann *TaintAnnotations, source string) map[string]*TaintSummary {
	out := map[string]*TaintSummary{}
	if ann == nil {
		return out
	}
	for funcName, params := range ann.FunctionParams {
		out[funcName] = computeSummaryFor(g, idx, ann, source, funcName, params)
	}
	return out
}

// RefineSummariesMultihopWithContext propagates param/return taint through the call graph.
func RefineSummariesMultihopWithContext(cg *CallGraph, ann *TaintAnnotations, summaries map[string]*TaintSummary, maxDepth int) {
	if cg == nil || summaries == nil {
		return
	}
	depth := maxDepth
	if depth < 1 {
		depth = 1
	}
	if depth > 4 {
		depth = 4
	}
	for hop := 1; hop < depth; hop++ {
		changed := false
		for _, site := range cg.Sites {
			calleeKey := site.Callee
			if i := strings.LastIndex(calleeKey, "."); i >= 0 {
				calleeKey = calleeKey[i+1:]
			}
			// Try full callee and bare name for method identities.
			calleeSum := summaries[site.Callee]
			if calleeSum == nil {
				calleeSum = summaries[calleeKey]
			}
			// Also try matching method identities that end with .bare
			if calleeSum == nil {
				for name, s := range summaries {
					if strings.HasSuffix(name, "."+calleeKey) || name == calleeKey {
						calleeSum = s
						break
					}
				}
			}
			if calleeSum == nil {
				continue
			}
			callerSum := summaries[site.Caller]
			if callerSum == nil {
				continue
			}
			if ann != nil {
				callerParams := ann.FunctionParams[site.Caller]
				for calleeIdx, isSrc := range calleeSum.ParamSources {
					if !isSrc || calleeIdx >= len(site.Arguments) {
						continue
					}
					arg := strings.TrimSpace(site.Arguments[calleeIdx])
					arg = strings.TrimPrefix(arg, "&")
					arg = strings.TrimSpace(arg)
					callerIdx := -1
					for i, p := range callerParams {
						if p == arg {
							callerIdx = i
							break
						}
					}
					if callerIdx < 0 {
						continue
					}
					for len(callerSum.ParamSources) <= callerIdx {
						callerSum.ParamSources = append(callerSum.ParamSources, false)
					}
					if !callerSum.ParamSources[callerIdx] {
						callerSum.ParamSources[callerIdx] = true
						changed = true
					}
					for _, sk := range calleeSum.SinkKinds {
						if !containsSink(callerSum.SinkKinds, sk) {
							callerSum.SinkKinds = append(callerSum.SinkKinds, sk)
						}
					}
				}
			}
			if !site.ReturnsResult {
				continue
			}
			returnsTaint := false
			for _, b := range calleeSum.ReturnSources {
				if b {
					returnsTaint = true
					break
				}
			}
			if !returnsTaint {
				continue
			}
			allFalse := true
			for _, b := range callerSum.ReturnSources {
				if b {
					allFalse = false
					break
				}
			}
			if allFalse {
				if len(callerSum.ReturnSources) == 0 {
					callerSum.ReturnSources = []bool{true}
				} else {
					for i := range callerSum.ReturnSources {
						callerSum.ReturnSources[i] = true
					}
				}
				for _, sk := range calleeSum.SinkKinds {
					if !containsSink(callerSum.SinkKinds, sk) {
						callerSum.SinkKinds = append(callerSum.SinkKinds, sk)
					}
				}
				changed = true
			}
		}
		if !changed {
			break
		}
	}
}

func computeSummaryFor(g *TaintGraph, idx *TaintGraphIndex, ann *TaintAnnotations, source, funcName string, params []string) *TaintSummary {
	sum := &TaintSummary{}
	if g == nil {
		return sum
	}

	paramNodeIDs := make([][]int, len(params))
	for i, param := range params {
		for id, n := range g.Nodes {
			if n.Kind == "variable" && n.Name == param && nodeInFunction(g, id, ann, funcName) {
				paramNodeIDs[i] = append(paramNodeIDs[i], id)
			}
		}
	}

	allSinks := []SinkKind{
		SinkCommandExec, SinkSQLQuery, SinkFileOpen, SinkTemplate,
		SinkHTTPWrite, SinkDeserialization, SinkLDAPQuery, SinkXMLQuery,
	}

	for _, ids := range paramNodeIDs {
		reaches := false
		for _, sk := range allSinks {
			var ok bool
			if sk == SinkSQLQuery {
				ok = ReachesSinkArgumentFromNodes(g, idx, ids, sk, 0)
			} else {
				ok = ReachesSinkFromNodes(g, idx, ids, sk)
			}
			if ok {
				reaches = true
				if !containsSink(sum.SinkKinds, sk) {
					sum.SinkKinds = append(sum.SinkKinds, sk)
				}
			}
		}
		sum.ParamSources = append(sum.ParamSources, reaches)
	}

	// Direct sources in function.
	var sourceIDs []int
	for _, ids := range g.BySource {
		for _, id := range ids {
			if nodeInFunction(g, id, ann, funcName) {
				sourceIDs = append(sourceIDs, id)
			}
		}
	}
	if len(sourceIDs) > 0 {
		for _, sk := range allSinks {
			if len(g.BySink[sk]) == 0 {
				continue
			}
			if ReachesSinkFromNodes(g, idx, sourceIDs, sk) {
				sum.HasDirectSink = true
				if !containsSink(sum.SinkKinds, sk) {
					sum.SinkKinds = append(sum.SinkKinds, sk)
				}
			}
		}
	}

	sum.ReturnSources = computeReturnSources(ann, source, funcName)
	sum.OutputPointerParams = computeOutputPointerParams(ann, source, funcName, params)
	return sum
}

func nodeInFunction(g *TaintGraph, nodeID int, ann *TaintAnnotations, funcName string) bool {
	if ann == nil || nodeID < 0 || nodeID >= len(g.Nodes) {
		return false
	}
	fr, ok := ann.FunctionRanges[funcName]
	if !ok {
		return false
	}
	n := g.Nodes[nodeID]
	switch n.Kind {
	case "variable":
		for _, s := range ann.Scopes {
			if s.ID == n.Scope {
				return s.Function == funcName
			}
		}
		return false
	case "source", "sink", "sanitizer":
		return n.ByteRange.Start >= fr.Start && n.ByteRange.End <= fr.End
	case "return":
		return n.Function == funcName
	}
	return false
}

func computeReturnSources(ann *TaintAnnotations, source, funcName string) []bool {
	if ann == nil {
		return []bool{false}
	}
	fr, ok := ann.FunctionRanges[funcName]
	if !ok {
		return []bool{false}
	}
	// Source returned from body.
	for _, src := range ann.Sources {
		if src.ByteRange.Start >= fr.Start && src.ByteRange.End <= fr.End {
			if sourceIsReturned(src, source, fr) {
				return []bool{true}
			}
		}
	}
	// Parameter returned.
	params := ann.FunctionParams[funcName]
	end := fr.End
	if end > len(source) {
		end = len(source)
	}
	start := fr.Start
	if start > end {
		start = end
	}
	body := source[start:end]
	for _, param := range params {
		for _, line := range strings.Split(body, "\n") {
			rest, ok := strings.CutPrefix(strings.TrimSpace(line), "return")
			if !ok {
				continue
			}
			rest = strings.TrimSpace(rest)
			endI := 0
			for endI < len(rest) {
				c := rest[endI]
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
					break
				}
				endI++
			}
			if rest[:endI] == param {
				return []bool{true}
			}
		}
	}
	return []bool{false}
}

func sourceIsReturned(src TaintSourceAnnotation, source string, fr ByteRange) bool {
	end := fr.End
	if end > len(source) {
		end = len(source)
	}
	start := fr.Start
	if start > end {
		start = end
	}
	body := source[start:end]
	ss := src.ByteRange.Start - start
	se := src.ByteRange.End - start
	if ss < 0 {
		ss = 0
	}
	if se > len(body) {
		se = len(body)
	}
	if se < ss {
		se = ss
	}
	// Same line has return?
	lineStart := strings.LastIndex(body[:ss], "\n") + 1
	lineEnd := len(body)
	if i := strings.Index(body[se:], "\n"); i >= 0 {
		lineEnd = se + i
	}
	if strings.Contains(body[lineStart:lineEnd], "return") {
		return true
	}
	if src.ResultVariable == "" {
		return false
	}
	for _, line := range strings.Split(body, "\n") {
		rest, ok := strings.CutPrefix(strings.TrimSpace(line), "return")
		if !ok {
			continue
		}
		rest = strings.TrimSpace(rest)
		endI := 0
		for endI < len(rest) {
			c := rest[endI]
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
				break
			}
			endI++
		}
		if rest[:endI] == src.ResultVariable {
			return true
		}
	}
	return false
}

func computeOutputPointerParams(ann *TaintAnnotations, source, funcName string, params []string) []int {
	if ann == nil {
		return nil
	}
	fr, ok := ann.FunctionRanges[funcName]
	if !ok {
		return nil
	}
	// Need a source in the body.
	hasSource := false
	for _, s := range ann.Sources {
		if s.ByteRange.Start >= fr.Start && s.ByteRange.End <= fr.End {
			hasSource = true
			break
		}
	}
	if !hasSource {
		return nil
	}
	end := fr.End
	if end > len(source) {
		end = len(source)
	}
	start := fr.Start
	if start > end {
		start = end
	}
	body := source[start:end]
	var out []int
	for i, param := range params {
		if strings.Contains(body, "*"+param+" =") || strings.Contains(body, "*"+param+"=") {
			out = append(out, i)
		}
	}
	return out
}

func containsSink(ss []SinkKind, k SinkKind) bool {
	for _, s := range ss {
		if s == k {
			return true
		}
	}
	return false
}

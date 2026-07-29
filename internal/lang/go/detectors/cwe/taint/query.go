package taint

// TaintGraphIndex is a forward adjacency map.
type TaintGraphIndex struct {
	Adj map[int][]adjEdge
}

type adjEdge struct {
	to   int
	kind EdgeKind
}

// BuildIndex builds adjacency from graph edges.
func BuildIndex(g *TaintGraph) *TaintGraphIndex {
	adj := map[int][]adjEdge{}
	if g == nil {
		return &TaintGraphIndex{Adj: adj}
	}
	for _, e := range g.Edges {
		adj[e.From] = append(adj[e.From], adjEdge{e.To, e.Kind})
	}
	return &TaintGraphIndex{Adj: adj}
}

type searchState struct {
	node      int
	sanitized bool
}

// FindTaintPaths finds paths from any source of sourceKind to sinks of sinkKind.
func FindTaintPaths(g *TaintGraph, sourceKind SourceKind, sinkKind SinkKind, allowed []SanitizerKind) []TaintPath {
	if g == nil {
		return nil
	}
	sourceIDs := g.BySource[sourceKind]
	sinkIDs := g.BySink[sinkKind]
	idx := BuildIndex(g)
	var paths []TaintPath
	for _, sinkID := range sinkIDs {
		if p, ok := bfsPath(g, idx.Adj, sourceIDs, sinkID, allowed); ok {
			paths = append(paths, p)
		}
	}
	return paths
}

// FindTaintPathsFromNodes finds paths from arbitrary start nodes to sinks.
func FindTaintPathsFromNodes(g *TaintGraph, startIDs []int, sinkKind SinkKind, allowed []SanitizerKind) []TaintPath {
	if g == nil {
		return nil
	}
	sinkIDs := g.BySink[sinkKind]
	idx := BuildIndex(g)
	var paths []TaintPath
	for _, sinkID := range sinkIDs {
		if p, ok := bfsPath(g, idx.Adj, startIDs, sinkID, allowed); ok {
			paths = append(paths, p)
		}
	}
	return paths
}

// ForwardReachesAny reports whether any start reaches any target.
func ForwardReachesAny(g *TaintGraph, starts, targets []int) bool {
	if g == nil || len(starts) == 0 || len(targets) == 0 {
		return false
	}
	idx := BuildIndex(g)
	return ForwardReachesAnyWithIndex(g, idx, starts, targets)
}

// ForwardReachesAnyWithIndex is the indexed form of ForwardReachesAny.
func ForwardReachesAnyWithIndex(g *TaintGraph, idx *TaintGraphIndex, starts, targets []int) bool {
	if g == nil || idx == nil || len(starts) == 0 || len(targets) == 0 {
		return false
	}
	targetSet := map[int]struct{}{}
	for _, t := range targets {
		targetSet[t] = struct{}{}
	}
	visited := make([]bool, len(g.Nodes))
	queue := make([]int, 0, len(starts))
	for _, s := range starts {
		if s >= 0 && s < len(visited) && !visited[s] {
			visited[s] = true
			queue = append(queue, s)
		}
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if _, ok := targetSet[cur]; ok {
			return true
		}
		for _, e := range idx.Adj[cur] {
			if e.to < len(visited) && !visited[e.to] {
				visited[e.to] = true
				queue = append(queue, e.to)
			}
		}
	}
	return false
}

// UnsanitizedReachesAny reports whether an unsanitized path exists from start to any target.
func UnsanitizedReachesAny(g *TaintGraph, start int, targets []int) bool {
	if g == nil {
		return false
	}
	idx := BuildIndex(g)
	return UnsanitizedReachesAnyWithIndex(g, idx, start, targets)
}

// UnsanitizedReachesAnyWithIndex is the indexed form.
func UnsanitizedReachesAnyWithIndex(g *TaintGraph, idx *TaintGraphIndex, start int, targets []int) bool {
	if g == nil || idx == nil {
		return false
	}
	targetSet := map[int]struct{}{}
	for _, t := range targets {
		targetSet[t] = struct{}{}
	}
	type st struct {
		node int
		san  bool
	}
	visited := map[st]struct{}{}
	queue := []st{{start, false}}
	visited[st{start, false}] = struct{}{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		san := cur.san
		if cur.node < len(g.Nodes) && g.Nodes[cur.node].Kind == "sanitizer" {
			san = true
		}
		if _, ok := targetSet[cur.node]; ok && !san {
			return true
		}
		for _, e := range idx.Adj[cur.node] {
			ns := st{e.to, san}
			if e.to < len(g.Nodes) {
				if _, seen := visited[ns]; !seen {
					visited[ns] = struct{}{}
					queue = append(queue, ns)
				}
			}
		}
	}
	return false
}

// ReachesSinkFromNodes reports whether any start reaches a sink of the given kind.
func ReachesSinkFromNodes(g *TaintGraph, idx *TaintGraphIndex, starts []int, sinkKind SinkKind) bool {
	return reachesSink(g, idx, starts, sinkKind, -1)
}

// ReachesSinkArgumentFromNodes requires the edge into the sink to be Argument(argIdx).
func ReachesSinkArgumentFromNodes(g *TaintGraph, idx *TaintGraphIndex, starts []int, sinkKind SinkKind, argIdx int) bool {
	return reachesSink(g, idx, starts, sinkKind, argIdx)
}

func reachesSink(g *TaintGraph, idx *TaintGraphIndex, starts []int, sinkKind SinkKind, argIdx int) bool {
	if g == nil || idx == nil {
		return false
	}
	targets := g.BySink[sinkKind]
	if len(starts) == 0 || len(targets) == 0 {
		return false
	}
	targetSet := map[int]struct{}{}
	for _, t := range targets {
		targetSet[t] = struct{}{}
	}
	visited := make([]bool, len(g.Nodes))
	queue := make([]int, 0, len(starts))
	for _, s := range starts {
		if s >= 0 && s < len(visited) && !visited[s] {
			visited[s] = true
			queue = append(queue, s)
		}
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, e := range idx.Adj[cur] {
			if _, ok := targetSet[e.to]; ok {
				if argIdx < 0 || (e.kind.IsArgument() && e.kind.ArgumentIndex() == argIdx) {
					return true
				}
			}
			if e.to < len(visited) && !visited[e.to] {
				visited[e.to] = true
				queue = append(queue, e.to)
			}
		}
	}
	return false
}

func bfsPath(g *TaintGraph, adj map[int][]adjEdge, sourceIDs []int, sinkID int, allowed []SanitizerKind) (TaintPath, bool) {
	type key = searchState
	predecessors := map[key]*key{}
	queue := make([]key, 0, len(sourceIDs))
	for _, sid := range sourceIDs {
		san := isSanitizerNode(g, sid, allowed)
		st := key{sid, san}
		queue = append(queue, st)
		predecessors[st] = nil
	}
	var bestSan *key
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.node == sinkID {
			if !cur.sanitized {
				path := reconstruct(predecessors, cur)
				return TaintPath{SourceID: path[0], SinkID: sinkID, NodeIDs: path, Sanitized: false}, true
			}
			if bestSan == nil {
				c := cur
				bestSan = &c
			}
			continue
		}
		for _, e := range adj[cur.node] {
			next := key{
				node:      e.to,
				sanitized: cur.sanitized || isSanitizerNode(g, e.to, allowed),
			}
			if e.to < len(g.Nodes) {
				if _, seen := predecessors[next]; !seen {
					c := cur
					predecessors[next] = &c
					queue = append(queue, next)
				}
			}
		}
	}
	if bestSan != nil {
		path := reconstruct(predecessors, *bestSan)
		return TaintPath{SourceID: path[0], SinkID: sinkID, NodeIDs: path, Sanitized: true}, true
	}
	return TaintPath{}, false
}

func reconstruct(pred map[searchState]*searchState, terminal searchState) []int {
	var path []int
	cur := terminal
	path = append(path, cur.node)
	for {
		p, ok := pred[cur]
		if !ok || p == nil {
			break
		}
		cur = *p
		path = append(path, cur.node)
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func isSanitizerNode(g *TaintGraph, nodeID int, allowed []SanitizerKind) bool {
	if nodeID < 0 || nodeID >= len(g.Nodes) {
		return false
	}
	n := g.Nodes[nodeID]
	if n.Kind != "sanitizer" {
		return false
	}
	for _, a := range allowed {
		if a == n.SanitizerKind {
			return true
		}
	}
	return false
}

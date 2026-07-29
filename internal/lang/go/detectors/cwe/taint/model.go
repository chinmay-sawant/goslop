// Package taint implements experimental intra/inter-procedural taint tracking
// for Go CWE-22/78/79/89. Name-string model — triage aid, not a security gate.
package taint

// ByteRange is a half-open [Start, End) source span.
type ByteRange struct {
	Start int
	End   int
}

// SourceKind classifies a taint source.
type SourceKind int

const (
	SourceUserInput SourceKind = iota
	SourceArgs
	SourceEnvVar
	SourceFile
	SourceNetwork
)

func (k SourceKind) String() string {
	switch k {
	case SourceUserInput:
		return "UserInput"
	case SourceArgs:
		return "Args"
	case SourceEnvVar:
		return "EnvVar"
	case SourceFile:
		return "File"
	case SourceNetwork:
		return "Network"
	default:
		return "SourceKind(?)"
	}
}

// SinkKind classifies a taint sink.
type SinkKind int

const (
	SinkCommandExec SinkKind = iota
	SinkSQLQuery
	SinkFileOpen
	SinkTemplate
	SinkHTTPWrite
	SinkDeserialization
	SinkLDAPQuery
	SinkXMLQuery
)

func (k SinkKind) String() string {
	switch k {
	case SinkCommandExec:
		return "CommandExec"
	case SinkSQLQuery:
		return "SQLQuery"
	case SinkFileOpen:
		return "FileOpen"
	case SinkTemplate:
		return "Template"
	case SinkHTTPWrite:
		return "HTTPWrite"
	case SinkDeserialization:
		return "Deserialization"
	case SinkLDAPQuery:
		return "LDAPQuery"
	case SinkXMLQuery:
		return "XMLQuery"
	default:
		return "SinkKind(?)"
	}
}

// SanitizerKind classifies a sanitizer.
type SanitizerKind int

const (
	SanitizerPath SanitizerKind = iota
	SanitizerHTML
	SanitizerURL
	SanitizerSQL
	SanitizerValidation
	SanitizerBounded
	SanitizerLDAP
	SanitizerXML
)

func (k SanitizerKind) String() string {
	switch k {
	case SanitizerPath:
		return "Path"
	case SanitizerHTML:
		return "HTML"
	case SanitizerURL:
		return "URL"
	case SanitizerSQL:
		return "SQL"
	case SanitizerValidation:
		return "Validation"
	case SanitizerBounded:
		return "Bounded"
	case SanitizerLDAP:
		return "LDAP"
	case SanitizerXML:
		return "XML"
	default:
		return "SanitizerKind(?)"
	}
}

// EdgeKind is edge semantics in the taint graph.
type EdgeKind int

const (
	EdgeAssignment EdgeKind = iota
	EdgeChannelTransfer
	EdgeReturn
	EdgePassThrough
	// EdgeArgument encodes argument index as EdgeArgument+idx (idx 0..N).
	// Use Argument(i) / ArgumentIndex().
	edgeArgumentBase = 100
)

// Argument returns an EdgeKind for positional argument i.
func Argument(i int) EdgeKind { return EdgeKind(edgeArgumentBase + i) }

// IsArgument reports whether k is an Argument edge.
func (k EdgeKind) IsArgument() bool { return int(k) >= edgeArgumentBase }

// ArgumentIndex returns the argument index for Argument edges.
func (k EdgeKind) ArgumentIndex() int {
	if !k.IsArgument() {
		return -1
	}
	return int(k) - edgeArgumentBase
}

// ScopeKind is a lexical scope kind.
type ScopeKind int

const (
	ScopePackage ScopeKind = iota
	ScopeFunction
	ScopeBlock
	ScopeIf
	ScopeFor
	ScopeSwitch
	ScopeCase
)

// ScopeInfo is scope metadata for variable resolution.
type ScopeInfo struct {
	ID        int
	Parent    *int
	Kind      ScopeKind
	ByteRange ByteRange
	Function  string // empty when package-level
}

// TaintNode is a node in the taint graph.
type TaintNode struct {
	// Kind discriminates: "source" | "variable" | "sink" | "sanitizer" | "return"
	Kind string

	// Shared fields
	Function  string
	ByteRange ByteRange

	// Source
	SourceKind SourceKind

	// Variable
	Name     string
	Scope    int
	DeclByte int

	// Sink
	SinkKind      SinkKind
	ArgumentIndex int

	// Sanitizer
	SanitizerKind SanitizerKind

	// Return
	ReturnIndex int
}

// TaintEdge is a directed edge.
type TaintEdge struct {
	From int
	To   int
	Kind EdgeKind
}

// TaintGraph is the intra-procedural data-flow graph for one unit.
type TaintGraph struct {
	Nodes      []TaintNode
	Edges      []TaintEdge
	ByVariable map[varKey][]int
	BySink     map[SinkKind][]int
	BySource   map[SourceKind][]int
}

type varKey struct {
	scope int
	name  string
}

// NewTaintGraph returns an empty graph with indexes.
func NewTaintGraph() *TaintGraph {
	return &TaintGraph{
		ByVariable: map[varKey][]int{},
		BySink:     map[SinkKind][]int{},
		BySource:   map[SourceKind][]int{},
	}
}

// AddNode appends a node and updates indexes.
func (g *TaintGraph) AddNode(n TaintNode) int {
	id := len(g.Nodes)
	if n.Kind == "variable" {
		k := varKey{scope: n.Scope, name: n.Name}
		g.ByVariable[k] = append(g.ByVariable[k], id)
	}
	if n.Kind == "source" {
		g.BySource[n.SourceKind] = append(g.BySource[n.SourceKind], id)
	}
	if n.Kind == "sink" {
		g.BySink[n.SinkKind] = append(g.BySink[n.SinkKind], id)
	}
	g.Nodes = append(g.Nodes, n)
	return id
}

// AddEdge appends an edge.
func (g *TaintGraph) AddEdge(from, to int, kind EdgeKind) {
	g.Edges = append(g.Edges, TaintEdge{From: from, To: to, Kind: kind})
}

// TaintSourceAnnotation is a raw source site.
type TaintSourceAnnotation struct {
	Function       string
	Kind           SourceKind
	ByteRange      ByteRange
	ResultVariable string // empty if none
	Arguments      []string
}

// TaintSinkAnnotation is a raw sink site.
type TaintSinkAnnotation struct {
	Function      string
	Kind          SinkKind
	ArgumentIndex int
	ArgumentText  string
	AllArguments  []string
	ByteRange     ByteRange
}

// TaintSanitizerAnnotation is a raw sanitizer site.
type TaintSanitizerAnnotation struct {
	Function       string
	Kind           SanitizerKind
	ByteRange      ByteRange
	ResultVariable string
	Arguments      []string
}

// AssignmentDetail is an assignment / short-decl site.
type AssignmentDetail struct {
	LHS                   string
	RHSText               string
	Scope                 int
	ByteRange             ByteRange
	FromSourceOrSanitizer bool
	IsChannelSend         bool
}

// ChannelTransfer is a paired same-function channel handoff.
type ChannelTransfer struct {
	Channel       string
	SendValueText string
	RecvLHS       string
	RecvScope     int
	SendByteRange ByteRange
	RecvByteRange ByteRange
}

// ChannelSendSite stages a send before pairing.
type ChannelSendSite struct {
	Channel       string
	ValueText     string
	FunctionScope int
	ByteRange     ByteRange
	InSelect      bool
}

// ChannelRecvSite stages a receive before pairing.
type ChannelRecvSite struct {
	Channel       string
	LHS           string // empty if no LHS
	FunctionScope int
	RecvScope     int
	ByteRange     ByteRange
	InSelect      bool
}

// UnsupportedFlowKind is an intentionally unmodeled flow.
type UnsupportedFlowKind int

const (
	UnsupportedChannel UnsupportedFlowKind = iota
	UnsupportedGoroutine
)

// UnsupportedFlow records an honest FN site.
type UnsupportedFlow struct {
	Kind      UnsupportedFlowKind
	ByteRange ByteRange
	Note      string
}

// TaintAnnotations are raw facts extracted from one unit.
type TaintAnnotations struct {
	Sources          []TaintSourceAnnotation
	Sinks            []TaintSinkAnnotation
	Sanitizers       []TaintSanitizerAnnotation
	Assignments      []AssignmentDetail
	Scopes           []ScopeInfo
	FunctionParams   map[string][]string
	FunctionRanges   map[string]ByteRange
	UnsupportedFlows []UnsupportedFlow
	ChannelTransfers []ChannelTransfer
}

// FunctionDecl is a discovered function/method declaration.
type FunctionDecl struct {
	Name         string
	ParamCount   int
	IsMethod     bool
	ReceiverType string // raw receiver text; empty for free functions
}

// CallSite is one call expression.
type CallSite struct {
	Caller        string
	Callee        string
	ByteRange     ByteRange
	Arguments     []string
	AssignmentLHS string
	ReturnsResult bool
	IsMethodCall  bool
	IsClosure     bool
}

// CallGraph is a per-file call graph.
type CallGraph struct {
	Sites        []CallSite
	ByCaller     map[string][]int
	ByCallee     map[string][]int
	Declarations map[string]FunctionDecl // identity → decl
}

// NewCallGraph returns an empty call graph.
func NewCallGraph() *CallGraph {
	return &CallGraph{
		ByCaller:     map[string][]int{},
		ByCallee:     map[string][]int{},
		Declarations: map[string]FunctionDecl{},
	}
}

// AddDeclaration records a function/method.
func (cg *CallGraph) AddDeclaration(identity string, d FunctionDecl) {
	if cg.Declarations == nil {
		cg.Declarations = map[string]FunctionDecl{}
	}
	cg.Declarations[identity] = d
}

// AddSite records a call site.
func (cg *CallGraph) AddSite(site CallSite) {
	idx := len(cg.Sites)
	if cg.ByCaller == nil {
		cg.ByCaller = map[string][]int{}
	}
	if cg.ByCallee == nil {
		cg.ByCallee = map[string][]int{}
	}
	cg.ByCaller[site.Caller] = append(cg.ByCaller[site.Caller], idx)
	cg.ByCallee[site.Callee] = append(cg.ByCallee[site.Callee], idx)
	cg.Sites = append(cg.Sites, site)
}

// TaintSummary is a function-level taint summary for inter-proc.
type TaintSummary struct {
	ParamSources        []bool // true = param reaches a sink
	ReturnSources       []bool
	HasDirectSink       bool
	SinkKinds           []SinkKind
	OutputPointerParams []int
}

// PackageIdentity keys same-package resolution.
type PackageIdentity struct {
	Dir  string
	Name string
}

// PackageFromUnit builds package identity from path + source.
func PackageFromUnit(path, source string) PackageIdentity {
	normalized := path
	for i := 0; i < len(normalized); i++ {
		if normalized[i] == '\\' {
			// normalize path separators
			b := []byte(normalized)
			for j := range b {
				if b[j] == '\\' {
					b[j] = '/'
				}
			}
			normalized = string(b)
			break
		}
	}
	dir := "."
	if i := lastSlash(normalized); i >= 0 {
		if i > 0 {
			dir = normalized[:i]
		}
	}
	name := packageClauseName(source)
	if name == "" {
		name = "_"
	}
	return PackageIdentity{Dir: dir, Name: name}
}

func lastSlash(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return i
		}
	}
	return -1
}

// PackageClauseName extracts the package clause identifier.
func packageClauseName(source string) string {
	for _, line := range splitLines(source) {
		trimmed := trimSpace(line)
		if len(trimmed) >= 8 && trimmed[:8] == "package " {
			rest := trimmed[8:]
			// take first token
			end := 0
			for end < len(rest) {
				c := rest[end]
				if c == ' ' || c == '\t' || c == '/' {
					break
				}
				end++
			}
			name := rest[:end]
			if name != "" {
				return name
			}
		}
	}
	return ""
}

// TaintSymbolKey is a project-level function/method key.
type TaintSymbolKey struct {
	Package  PackageIdentity
	Receiver string // normalized; empty for free functions
	Name     string
}

// FunctionKey builds a free-function key.
func FunctionKey(pkg PackageIdentity, name string) TaintSymbolKey {
	return TaintSymbolKey{Package: pkg, Name: name}
}

// MethodKey builds a method key with optional raw receiver text.
func MethodKey(pkg PackageIdentity, receiverRaw, name string) TaintSymbolKey {
	recv := ""
	if receiverRaw != "" {
		recv = NormalizeReceiverType(receiverRaw)
	}
	return TaintSymbolKey{Package: pkg, Receiver: recv, Name: name}
}

// NormalizeReceiverType normalizes receiver parameter text to a type key.
func NormalizeReceiverType(raw string) string {
	s := trimSpace(raw)
	if len(s) > 0 && s[0] == '(' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == ')' {
		s = s[:len(s)-1]
	}
	s = trimSpace(s)
	parts := fields(s)
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	default:
		return parts[1]
	}
}

// SummaryNameForKey matches per-file extraction identity.
func SummaryNameForKey(key TaintSymbolKey) string {
	if key.Receiver != "" {
		return key.Receiver + "." + key.Name
	}
	return key.Name
}

// TaintPath is a source→sink path through the graph.
type TaintPath struct {
	SourceID  int
	SinkID    int
	NodeIDs   []int
	Sanitized bool
}

// helpers without importing strings to keep model self-contained for types

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			out = append(out, line)
			start = i + 1
		}
	}
	if start <= len(s) {
		out = append(out, s[start:])
	}
	return out
}

func fields(s string) []string {
	var out []string
	start := -1
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			if start >= 0 {
				out = append(out, s[start:i])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		out = append(out, s[start:])
	}
	return out
}

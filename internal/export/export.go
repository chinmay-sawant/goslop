// Package export writes opt-in finding context and chunk files (Rust export parity).
package export

import (
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/chinmay/goslop/internal/cwe"
	"github.com/chinmay/goslop/internal/rules"
)

// Default dirs match the Rust product.
const (
	DefaultContextDir = "scripts/findings/functions"
	DefaultChunksDir  = "scripts/chunks"
	DefaultChunkSize  = 25
	// DefaultWholeFunction is the default Context expansion mode for exports.
	DefaultWholeFunction = true
)

// Options controls on-disk finding context/chunk export.
type Options struct {
	ExportContext    bool
	ExportChunks     bool
	ChunkSize        int
	ContextOutputDir string
	ChunksOutputDir  string
	// WholeFunction, when nil, defaults to DefaultWholeFunction (true): expand
	// Context to the enclosing FuncDecl/FuncLit. When false, uses a nearby
	// ~4-line window. Applies to both per-finding refs and chunks.
	WholeFunction *bool
}

// Summary counts files written by a successful export pass.
type Summary struct {
	ContextFilesWritten int
	ChunkFilesWritten   int
}

// Enabled reports whether any export is requested.
func (o Options) Enabled() bool {
	return o.ExportContext || o.ExportChunks
}

// Normalize fills defaults for dirs and chunk size.
func (o Options) Normalize() Options {
	if o.ContextOutputDir == "" {
		o.ContextOutputDir = DefaultContextDir
	}
	if o.ChunksOutputDir == "" {
		o.ChunksOutputDir = DefaultChunksDir
	}
	if o.ChunkSize <= 0 {
		o.ChunkSize = DefaultChunkSize
	}
	return o
}

// wholeFunctionEnabled reports whether whole-function Context expansion is on.
func (o Options) wholeFunctionEnabled() bool {
	if o.WholeFunction == nil {
		return DefaultWholeFunction
	}
	return *o.WholeFunction
}

// BoolPtr returns a *bool for Options.WholeFunction and tests.
func BoolPtr(v bool) *bool { return &v }

// ExportFindings writes context and/or chunk files for findings.
// sourceCache maps display path → source text (from AnalysisResult.SourceCache).
func ExportFindings(findings []rules.Finding, opts Options, sourceCache map[string]string) (Summary, error) {
	opts = opts.Normalize()
	var sum Summary
	if !opts.Enabled() {
		return sum, nil
	}
	if opts.ExportContext && opts.ExportChunks {
		if same, err := dirsEqual(opts.ContextOutputDir, opts.ChunksOutputDir); err != nil {
			return sum, err
		} else if same {
			return sum, fmt.Errorf("context and chunk exports must use different output directories")
		}
	}

	fileCache := map[string]string{}
	astCache := map[string]*parsedFile{}
	total := len(findings)
	wholeFn := opts.wholeFunctionEnabled()

	if opts.ExportContext {
		if err := os.MkdirAll(opts.ContextOutputDir, 0o755); err != nil {
			return sum, fmt.Errorf("create context dir: %w", err)
		}
		if err := cleanMatchingTxt(opts.ContextOutputDir, func(name string) bool {
			// numbered N.txt files only
			for _, r := range name {
				if r < '0' || r > '9' {
					return strings.HasSuffix(name, ".txt") && name != "" && !strings.Contains(name, "Chunk")
				}
			}
			return strings.HasSuffix(name, ".txt")
		}); err != nil {
			return sum, err
		}
		// Clean all *.txt in context dir (simple parity with fresh stage).
		if err := cleanMatchingTxt(opts.ContextOutputDir, func(name string) bool {
			return strings.HasSuffix(name, ".txt")
		}); err != nil {
			return sum, err
		}
		for i, f := range findings {
			text := formatFindingBlock(f, i+1, total, fileCache, sourceCache, astCache, wholeFn)
			path := filepath.Join(opts.ContextOutputDir, fmt.Sprintf("%d.txt", i+1))
			if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
				return sum, fmt.Errorf("write context %s: %w", path, err)
			}
			sum.ContextFilesWritten++
		}
	}

	if opts.ExportChunks {
		n, err := writeChunkFiles(findings, opts.ChunksOutputDir, opts.ChunkSize, fileCache, sourceCache, astCache, wholeFn)
		if err != nil {
			return sum, err
		}
		sum.ChunkFilesWritten = n
	}
	return sum, nil
}

func dirsEqual(a, b string) (bool, error) {
	if err := os.MkdirAll(a, 0o755); err != nil {
		return false, fmt.Errorf("mkdir context dir: %w", err)
	}
	if err := os.MkdirAll(b, 0o755); err != nil {
		return false, fmt.Errorf("mkdir chunks dir: %w", err)
	}
	ca, err := filepath.Abs(a)
	if err != nil {
		return false, fmt.Errorf("abs context dir: %w", err)
	}
	cb, err := filepath.Abs(b)
	if err != nil {
		return false, fmt.Errorf("abs chunks dir: %w", err)
	}
	return ca == cb, nil
}

func cleanMatchingTxt(dir string, keep func(name string) bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if keep(name) {
			_ = os.Remove(filepath.Join(dir, name))
		}
	}
	return nil
}

func writeChunkFiles(
	findings []rules.Finding,
	outputDir string,
	chunkSize int,
	fileCache map[string]string,
	sourceCache map[string]string,
	astCache map[string]*parsedFile,
	wholeFunction bool,
) (int, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return 0, fmt.Errorf("create chunks dir: %w", err)
	}
	if err := cleanMatchingTxt(outputDir, func(name string) bool {
		return strings.HasPrefix(name, "Chunk_") && strings.HasSuffix(name, ".txt")
	}); err != nil {
		return 0, err
	}
	if len(findings) == 0 {
		return 0, nil
	}
	if chunkSize < 1 {
		chunkSize = DefaultChunkSize
	}
	sep := strings.Repeat("=", 100)
	total := len(findings)
	chunks := 0
	for start := 0; start < total; start += chunkSize {
		end := start + chunkSize
		if end > total {
			end = total
		}
		startIdx := start + 1
		endIdx := end
		path := filepath.Join(outputDir, fmt.Sprintf("Chunk_%d_%d.txt", startIdx, endIdx))
		var b strings.Builder
		_, _ = fmt.Fprintf(&b, "Findings %d-%d of %d\n\n", startIdx, endIdx, total)
		for i := start; i < end; i++ {
			if i > start {
				b.WriteByte('\n')
				b.WriteString(sep)
				b.WriteString("\n\n")
			}
			block := formatFindingBlock(findings[i], i+1, total, fileCache, sourceCache, astCache, wholeFunction)
			b.WriteString(strings.TrimRight(block, "\n"))
		}
		b.WriteByte('\n')
		if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
			return chunks, fmt.Errorf("write chunk %s: %w", path, err)
		}
		chunks++
	}
	return chunks, nil
}

func formatFindingBlock(
	f rules.Finding,
	index, total int,
	fileCache map[string]string,
	sourceCache map[string]string,
	astCache map[string]*parsedFile,
	wholeFunction bool,
) string {
	lines := make([]string, 0, 16)
	lines = append(lines,
		fmt.Sprintf("Finding %d/%d", index, total),
		fmt.Sprintf("Source: %s:%d:%d", f.File, f.Line, f.Column),
		fmt.Sprintf("Rule: %s", f.RuleID),
		fmt.Sprintf("Fingerprint: %s", f.Fingerprint),
		fmt.Sprintf("Rule title: %s", f.RuleTitle),
		fmt.Sprintf("Severity: %s", f.Severity.String()),
		fmt.Sprintf("Message: %s", f.Message),
	)
	if len(f.CWE) > 0 {
		lines = append(lines, fmt.Sprintf("CWEs: %s", cwe.FormatList(f.CWE)))
	}
	if f.Fix != "" {
		lines = append(lines, fmt.Sprintf("Fix: %s", f.Fix))
	}
	if f.Confidence > 0 {
		lines = append(lines, fmt.Sprintf("Confidence: %g", f.Confidence))
	}
	if len(f.Tags) > 0 {
		lines = append(lines, fmt.Sprintf("Tags: %s", strings.Join(f.Tags, ", ")))
	}
	if f.Remediation != "" {
		lines = append(lines, fmt.Sprintf("Remediation: %s", f.Remediation))
	}
	lines = append(lines, "Context:")
	for _, cl := range findingContextLines(f, fileCache, sourceCache, astCache, wholeFunction) {
		lines = append(lines, "    "+cl)
	}
	return strings.Join(lines, "\n") + "\n"
}

func findingContextLines(
	f rules.Finding,
	fileCache map[string]string,
	sourceCache map[string]string,
	astCache map[string]*parsedFile,
	wholeFunction bool,
) []string {
	// Nearby-window mode: prefer an explicit detector snippet when present.
	if !wholeFunction && f.Snippet != "" {
		if out := snippetLines(f.Snippet); len(out) > 0 {
			return out
		}
	}
	content := loadSource(f.File, fileCache, sourceCache)
	if content == "" {
		if f.Snippet != "" {
			if out := snippetLines(f.Snippet); len(out) > 0 {
				return out
			}
		}
		return []string{"<context unavailable>"}
	}
	if wholeFunction {
		if lines := functionWindow(f, content, astCache); len(lines) > 0 {
			return lines
		}
	}
	return lineWindow(f, content)
}

func snippetLines(snippet string) []string {
	var out []string
	for _, line := range strings.Split(snippet, "\n") {
		line = strings.TrimRight(line, "\r")
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func loadSource(path string, fileCache map[string]string, sourceCache map[string]string) string {
	if content := sourceCache[path]; content != "" {
		return content
	}
	if cached, ok := fileCache[path]; ok {
		return cached
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fileCache[path] = ""
		return ""
	}
	content := string(data)
	fileCache[path] = content
	return content
}

// parsedFile caches a go/ast parse of one source file for whole-function spans.
type parsedFile struct {
	fset *token.FileSet
	file *goast.File
	// failed marks a permanent parse failure (no usable File).
	failed bool
}

func getParsed(path, content string, astCache map[string]*parsedFile) *parsedFile {
	if astCache != nil {
		if p, ok := astCache[path]; ok {
			return p
		}
	}
	fset := token.NewFileSet()
	// SkipObjectResolution keeps export parsing cheap; partial File is ok on error.
	file, err := parser.ParseFile(fset, path, content, parser.SkipObjectResolution)
	p := &parsedFile{fset: fset, file: file}
	if file == nil {
		p.failed = true
		_ = err
	}
	if astCache != nil {
		astCache[path] = p
	}
	return p
}

// functionWindow returns numbered lines for the innermost FuncDecl/FuncLit that
// contains the finding line. Empty when no enclosing function is found.
func functionWindow(f rules.Finding, content string, astCache map[string]*parsedFile) []string {
	if f.Line < 1 {
		return nil
	}
	p := getParsed(f.File, content, astCache)
	if p == nil || p.failed || p.file == nil {
		return nil
	}
	startLine, endLine := enclosingFunctionLines(p.fset, p.file, f.Line)
	if startLine == 0 {
		return nil
	}
	return numberedLines(content, startLine, endLine, f.Line)
}

// enclosingFunctionLines returns the 1-based [start, end] line span of the
// function-like node containing hitLine.
//
// Preference (for agent/human export context):
//  1. Outermost *FuncDecl* that contains the hit (named method/function).
//  2. Else outermost *FuncLit* (closures / defer func() when not inside a decl).
//
// Named decls win over tiny inner closures so a hit inside `defer func() { ... }`
// still exports the full surrounding method.
func enclosingFunctionLines(fset *token.FileSet, file *goast.File, hitLine int) (startLine, endLine int) {
	var declStart, declEnd, litStart, litEnd int
	var declSpan, litSpan int // track largest (outermost)
	goast.Inspect(file, func(n goast.Node) bool {
		if n == nil {
			return false
		}
		var pos, end token.Pos
		isDecl := false
		switch n := n.(type) {
		case *goast.FuncDecl:
			pos, end = n.Pos(), n.End()
			isDecl = true
		case *goast.FuncLit:
			pos, end = n.Pos(), n.End()
		default:
			return true
		}
		s := fset.Position(pos).Line
		e := fset.Position(end).Line
		if hitLine < s || hitLine > e {
			return true
		}
		span := e - s
		if isDecl {
			if declStart == 0 || span > declSpan {
				declStart, declEnd, declSpan = s, e, span
			}
			return true
		}
		if litStart == 0 || span > litSpan {
			litStart, litEnd, litSpan = s, e, span
		}
		return true
	})
	if declStart != 0 {
		return declStart, declEnd
	}
	return litStart, litEnd
}

func lineWindow(f rules.Finding, content string) []string {
	start := f.Line - 2
	if start < 1 {
		start = 1
	}
	end := f.Line + 1
	return numberedLines(content, start, end, f.Line)
}

func numberedLines(content string, start, end, hitLine int) []string {
	if start < 1 {
		start = 1
	}
	all := strings.Split(content, "\n")
	// strings.Split yields a trailing empty element when content ends with \n;
	// that empty last element is still a valid "line" only if end reaches it.
	lines := make([]string, 0, end-start+1)
	for i, line := range all {
		lineNo := i + 1
		if lineNo < start {
			continue
		}
		if lineNo > end {
			break
		}
		marker := " "
		if lineNo == hitLine {
			marker = ">"
		}
		lines = append(lines, fmt.Sprintf("%s %5d: %s", marker, lineNo, line))
	}
	if len(lines) == 0 {
		return []string{"<context unavailable>"}
	}
	return lines
}

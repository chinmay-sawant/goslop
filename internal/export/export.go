// Package export writes opt-in finding context and chunk files (Rust export parity).
package export

import (
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
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

	// When both surfaces export, format each finding once and reuse (P1.3).
	dual := opts.ExportContext && opts.ExportChunks
	var blocks []string
	if dual {
		blocks = make([]string, total)
		for i, f := range findings {
			blocks[i] = formatFindingBlock(f, i+1, total, fileCache, sourceCache, astCache, wholeFn)
		}
	}

	if opts.ExportContext {
		if err := os.MkdirAll(opts.ContextOutputDir, 0o755); err != nil {
			return sum, fmt.Errorf("create context dir: %w", err)
		}
		if err := cleanOwnedFiles(opts.ContextOutputDir, isContextOutputFile); err != nil {
			return sum, fmt.Errorf("clean context output: %w", err)
		}
		for i, f := range findings {
			var text string
			if dual {
				text = blocks[i]
			} else {
				text = formatFindingBlock(f, i+1, total, fileCache, sourceCache, astCache, wholeFn)
			}
			path := filepath.Join(opts.ContextOutputDir, fmt.Sprintf("%d.txt", i+1))
			if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
				return sum, fmt.Errorf("write context %s: %w", path, err)
			}
			sum.ContextFilesWritten++
		}
	}

	if opts.ExportChunks {
		n, err := writeChunkFiles(findings, opts.ChunksOutputDir, opts.ChunkSize, fileCache, sourceCache, astCache, wholeFn, blocks)
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

// cleanOwnedFiles removes only files owned by this export surface. Callers may
// choose a directory that also contains their own files, so broad extension
// based cleanup is never safe here.
func cleanOwnedFiles(dir string, owns func(name string) bool) error {
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
		if owns(name) {
			if err := removeOwnedFile(filepath.Join(dir, name)); err != nil {
				return fmt.Errorf("remove %s: %w", filepath.Join(dir, name), err)
			}
		}
	}
	return nil
}

// removeOwnedFile is kept as a narrow filesystem seam so cleanup failures can
// be tested deterministically. Production behavior is os.Remove.
var removeOwnedFile = os.Remove

func isContextOutputFile(name string) bool {
	if !strings.HasSuffix(name, ".txt") {
		return false
	}
	base := strings.TrimSuffix(name, ".txt")
	if base == "" {
		return false
	}
	for _, r := range base {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func writeChunkFiles(
	findings []rules.Finding,
	outputDir string,
	chunkSize int,
	fileCache map[string]string,
	sourceCache map[string]string,
	astCache map[string]*parsedFile,
	wholeFunction bool,
	preformatted []string, // optional: reuse dual-export blocks (same index as findings)
) (int, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return 0, fmt.Errorf("create chunks dir: %w", err)
	}
	if err := cleanOwnedFiles(outputDir, func(name string) bool {
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
			var block string
			if i < len(preformatted) && preformatted[i] != "" {
				block = preformatted[i]
			} else {
				block = formatFindingBlock(findings[i], i+1, total, fileCache, sourceCache, astCache, wholeFunction)
			}
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

// funcSpan is a 1-based line range of a FuncDecl or FuncLit.
type funcSpan struct {
	start, end int
	isDecl     bool
}

// parsedFile caches a go/ast parse of one source file for whole-function spans.
type parsedFile struct {
	fset *token.FileSet
	file *goast.File
	// spans is built once at parse time (P0.2): all FuncDecl/FuncLit line ranges.
	spans []funcSpan
	// lineStarts[i] = byte offset of 1-based line i+1 (for numberedLines).
	lineStarts []int
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
	p := &parsedFile{fset: fset, file: file, lineStarts: computeLineStarts(content)}
	if file == nil {
		p.failed = true
		_ = err
	} else {
		p.spans = collectFuncSpans(fset, file)
	}
	if astCache != nil {
		astCache[path] = p
	}
	return p
}

func computeLineStarts(content string) []int {
	// lineStarts[0] is byte offset of line 1.
	starts := make([]int, 0, 64)
	starts = append(starts, 0)
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			starts = append(starts, i+1)
		}
	}
	return starts
}

func collectFuncSpans(fset *token.FileSet, file *goast.File) []funcSpan {
	var spans []funcSpan
	goast.Inspect(file, func(n goast.Node) bool {
		if n == nil {
			return false
		}
		var pos, end token.Pos
		isDecl := false
		switch x := n.(type) {
		case *goast.FuncDecl:
			pos, end = x.Pos(), x.End()
			isDecl = true
		case *goast.FuncLit:
			pos, end = x.Pos(), x.End()
		default:
			return true
		}
		spans = append(spans, funcSpan{
			start:  fset.Position(pos).Line,
			end:    fset.Position(end).Line,
			isDecl: isDecl,
		})
		return true
	})
	return spans
}

// functionWindow returns numbered lines for the outermost FuncDecl (preferred)
// or FuncLit that contains the finding line. Empty when none found.
func functionWindow(f rules.Finding, content string, astCache map[string]*parsedFile) []string {
	if f.Line < 1 {
		return nil
	}
	// P2.3: skip expensive parse when whole-function span is impossible.
	if !strings.HasSuffix(f.File, ".go") && !strings.HasSuffix(f.File, ".go.txt") {
		// Display paths are almost always .go; non-Go sources use line window.
		return nil
	}
	if content == "" || (!strings.Contains(content, "func ") && !strings.Contains(content, "func(")) {
		return nil
	}
	// Package-level / go.mod-style findings: no enclosing function to expand.
	if isPackageLevelExport(f) {
		return nil
	}
	p := getParsed(f.File, content, astCache)
	if p == nil || p.failed || p.file == nil {
		return nil
	}
	startLine, endLine := enclosingFromSpans(p.spans, f.Line)
	if startLine == 0 {
		return nil
	}
	return numberedLinesCached(content, p.lineStarts, startLine, endLine, f.Line)
}

// isPackageLevelExport reports findings that never sit inside a function body
// (module/package hygiene and similar). Avoids parse for whole-function Context.
func isPackageLevelExport(f rules.Finding) bool {
	switch f.RuleID {
	case "BP-41", // missing package doc
		"BP-57", "BP-58", "BP-59", "BP-60", "BP-61", "BP-62", "BP-63", "BP-64", "BP-65":
		return true
	default:
		return false
	}
}

// enclosingFromSpans picks outermost FuncDecl containing hitLine, else outermost FuncLit.
func enclosingFromSpans(spans []funcSpan, hitLine int) (int, int) {
	var declStart, declEnd, litStart, litEnd int
	var declSpan, litSpan int
	for _, s := range spans {
		if hitLine < s.start || hitLine > s.end {
			continue
		}
		span := s.end - s.start
		if s.isDecl {
			if declStart == 0 || span > declSpan {
				declStart, declEnd, declSpan = s.start, s.end, span
			}
			continue
		}
		if litStart == 0 || span > litSpan {
			litStart, litEnd, litSpan = s.start, s.end, span
		}
	}
	if declStart != 0 {
		return declStart, declEnd
	}
	return litStart, litEnd
}

// enclosingFunctionLines returns the 1-based [start, end] line span of the
// function-like node containing hitLine. Prefer outermost FuncDecl over FuncLit.
// Used by tests and as a thin wrapper over collectFuncSpans.
func enclosingFunctionLines(fset *token.FileSet, file *goast.File, hitLine int) (int, int) {
	return enclosingFromSpans(collectFuncSpans(fset, file), hitLine)
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
	return numberedLinesCached(content, computeLineStarts(content), start, end, hitLine)
}

// numberedLinesCached formats lines without fmt.Sprintf or re-splitting the file (P1.2).
func numberedLinesCached(content string, lineStarts []int, start, end, hitLine int) []string {
	if start < 1 {
		start = 1
	}
	nLines := len(lineStarts)
	if nLines == 0 {
		return []string{"<context unavailable>"}
	}
	if end > nLines {
		end = nLines
	}
	if start > end {
		return []string{"<context unavailable>"}
	}
	out := make([]string, 0, end-start+1)
	var b strings.Builder
	for lineNo := start; lineNo <= end; lineNo++ {
		off := lineStarts[lineNo-1]
		lineEnd := len(content)
		if lineNo < nLines {
			lineEnd = lineStarts[lineNo] - 1 // exclude '\n'
			if lineEnd < off {
				lineEnd = off
			}
		}
		// Drop trailing \r if present (Windows sources).
		for lineEnd > off && content[lineEnd-1] == '\r' {
			lineEnd--
		}
		line := content[off:lineEnd]
		marker := " "
		if lineNo == hitLine {
			marker = ">"
		}
		b.Reset()
		b.Grow(8 + len(line))
		b.WriteByte(marker[0])
		b.WriteByte(' ')
		// right-align line number in 5 columns (product-compatible)
		num := strconv.Itoa(lineNo)
		for pad := 5 - len(num); pad > 0; pad-- {
			b.WriteByte(' ')
		}
		b.WriteString(num)
		b.WriteString(": ")
		b.WriteString(line)
		out = append(out, b.String())
	}
	if len(out) == 0 {
		return []string{"<context unavailable>"}
	}
	return out
}

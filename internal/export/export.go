// Package export writes opt-in finding context and chunk files (Rust export parity).
package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chinmay/codehound/internal/cwe"
	"github.com/chinmay/codehound/internal/rules"
)

// Default dirs match the Rust product.
const (
	DefaultContextDir = "scripts/findings/functions"
	DefaultChunksDir  = "scripts/chunks"
	DefaultChunkSize  = 25
)

// Options controls on-disk finding context/chunk export.
type Options struct {
	ExportContext    bool
	ExportChunks     bool
	ChunkSize        int
	ContextOutputDir string
	ChunksOutputDir  string
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
	total := len(findings)

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
			text := formatFindingBlock(f, i+1, total, fileCache, sourceCache)
			path := filepath.Join(opts.ContextOutputDir, fmt.Sprintf("%d.txt", i+1))
			if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
				return sum, fmt.Errorf("write context %s: %w", path, err)
			}
			sum.ContextFilesWritten++
		}
	}

	if opts.ExportChunks {
		n, err := writeChunkFiles(findings, opts.ChunksOutputDir, opts.ChunkSize, fileCache, sourceCache)
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
			block := formatFindingBlock(findings[i], i+1, total, fileCache, sourceCache)
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
	for _, cl := range findingContextLines(f, fileCache, sourceCache) {
		lines = append(lines, "    "+cl)
	}
	return strings.Join(lines, "\n") + "\n"
}

func findingContextLines(
	f rules.Finding,
	fileCache map[string]string,
	sourceCache map[string]string,
) []string {
	if f.Snippet != "" {
		var out []string
		for _, line := range strings.Split(f.Snippet, "\n") {
			line = strings.TrimRight(line, "\r")
			if line != "" {
				out = append(out, line)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	content := sourceCache[f.File]
	if content == "" {
		if cached, ok := fileCache[f.File]; ok {
			content = cached
		} else {
			data, err := os.ReadFile(f.File)
			if err != nil {
				fileCache[f.File] = ""
			} else {
				content = string(data)
				fileCache[f.File] = content
			}
		}
	}
	if content == "" {
		return []string{"<context unavailable>"}
	}
	return lineWindow(f, content)
}

func lineWindow(f rules.Finding, content string) []string {
	start := f.Line - 2
	if start < 1 {
		start = 1
	}
	end := f.Line + 1
	all := strings.Split(content, "\n")
	lines := make([]string, 0, 4)
	for i, line := range all {
		lineNo := i + 1
		if lineNo < start {
			continue
		}
		if lineNo > end {
			break
		}
		marker := " "
		if lineNo == f.Line {
			marker = ">"
		}
		lines = append(lines, fmt.Sprintf("%s %5d: %s", marker, lineNo, line))
	}
	if len(lines) == 0 {
		return []string{"<context unavailable>"}
	}
	return lines
}

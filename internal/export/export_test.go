package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay/goslop/internal/export"
	"github.com/chinmay/goslop/internal/rules"
)

func TestExportContextAndChunks(t *testing.T) {
	dir := t.TempDir()
	ctxDir := filepath.Join(dir, "ctx")
	chunkDir := filepath.Join(dir, "chunks")
	srcPath := filepath.Join(dir, "sample.go")
	src := "package sample\n\nfunc Bad() {\n\t_ = 1\n}\n"
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	findings := []rules.Finding{
		{
			RuleID:      "PERF-6",
			RuleTitle:   "Fmt In Loop",
			File:        srcPath,
			Line:        4,
			Column:      2,
			Message:     "fmt in loop",
			Severity:    rules.SeverityMedium,
			Fingerprint: "goslop:2:PERF-6:sample.go:deadbeef",
		},
		{
			RuleID:      "CWE-78",
			RuleTitle:   "Command Injection",
			File:        srcPath,
			Line:        3,
			Column:      1,
			Message:     "shell",
			Severity:    rules.SeverityHigh,
			Fingerprint: "goslop:2:CWE-78:sample.go:cafebabe",
		},
	}
	// Two findings → chunk size 1 → 2 chunk files; context → 2 files.
	sum, err := export.ExportFindings(findings, export.Options{
		ExportContext:    true,
		ExportChunks:     true,
		ChunkSize:        1,
		ContextOutputDir: ctxDir,
		ChunksOutputDir:  chunkDir,
	}, map[string]string{srcPath: src})
	if err != nil {
		t.Fatal(err)
	}
	if sum.ContextFilesWritten != 2 {
		t.Fatalf("context files=%d", sum.ContextFilesWritten)
	}
	if sum.ChunkFilesWritten != 2 {
		t.Fatalf("chunk files=%d", sum.ChunkFilesWritten)
	}
	data, err := os.ReadFile(filepath.Join(ctxDir, "1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "Finding 1/2") || !strings.Contains(text, "Rule: PERF-6") {
		t.Fatalf("context body:\n%s", text)
	}
	if !strings.Contains(text, ">") {
		t.Fatal("expected line marker in context")
	}
}

func TestExportDisabledNoOp(t *testing.T) {
	sum, err := export.ExportFindings(nil, export.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sum.ContextFilesWritten != 0 || sum.ChunkFilesWritten != 0 {
		t.Fatalf("%+v", sum)
	}
}

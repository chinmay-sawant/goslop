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
	// Default whole_function: Context should include the full Bad() body.
	if !strings.Contains(text, "func Bad()") {
		t.Fatalf("default whole-function context missing func:\n%s", text)
	}
}

func TestExportWholeFunctionVsWindow(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "sample.go")
	// Long enough that line-2..line+1 around the hit is not the whole body.
	src := `package sample

func Process(x int) int {
	a := x + 1
	b := a * 2
	c := b + 3
	d := c * 4
	e := d - 5
	return e
}
`
	// Line numbers: 3=func, 4=a, 5=b, 6=c, 7=d, 8=e, 9=return, 10=}
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	finding := rules.Finding{
		RuleID:      "PERF-1",
		RuleTitle:   "Test",
		File:        srcPath,
		Line:        7, // d := c * 4  → window 5..8, no signature / return
		Column:      2,
		Message:     "hit",
		Severity:    rules.SeverityLow,
		Fingerprint: "goslop:2:PERF-1:sample.go:abcd",
	}
	cache := map[string]string{srcPath: src}

	// whole_function = true (default / explicit)
	wholeDir := filepath.Join(dir, "whole")
	_, err := export.ExportFindings([]rules.Finding{finding}, export.Options{
		ExportContext:    true,
		ContextOutputDir: wholeDir,
		WholeFunction:    export.BoolPtr(true),
	}, cache)
	if err != nil {
		t.Fatal(err)
	}
	wholeBody := readFile(t, filepath.Join(wholeDir, "1.txt"))
	if !strings.Contains(wholeBody, "func Process") {
		t.Fatalf("whole function missing signature:\n%s", wholeBody)
	}
	if !strings.Contains(wholeBody, "return e") {
		t.Fatalf("whole function missing end of body:\n%s", wholeBody)
	}
	if !strings.Contains(wholeBody, ">     7:") {
		t.Fatalf("expected hit marker on line 7:\n%s", wholeBody)
	}

	// whole_function = false → ~4-line window (line-2..line+1)
	windowDir := filepath.Join(dir, "window")
	_, err = export.ExportFindings([]rules.Finding{finding}, export.Options{
		ExportContext:    true,
		ContextOutputDir: windowDir,
		WholeFunction:    export.BoolPtr(false),
	}, cache)
	if err != nil {
		t.Fatal(err)
	}
	windowBody := readFile(t, filepath.Join(windowDir, "1.txt"))
	if strings.Contains(windowBody, "func Process") {
		t.Fatalf("window mode should not include signature:\n%s", windowBody)
	}
	if strings.Contains(windowBody, "return e") {
		t.Fatalf("window mode should not include return:\n%s", windowBody)
	}
	if !strings.Contains(windowBody, ">     7:") {
		t.Fatalf("window mode missing hit marker:\n%s", windowBody)
	}
	// Window is at most ~4 content lines under Context:
	ctxPart := windowBody[strings.Index(windowBody, "Context:"):]
	lineCount := strings.Count(ctxPart, "\n")
	if lineCount > 6 { // "Context:" + up to 4 lines + trailing
		t.Fatalf("window context too large (%d lines):\n%s", lineCount, ctxPart)
	}
}

func TestExportPrefersOuterFuncDeclOverClosure(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "sample.go")
	src := `package sample

func Handler() error {
	f, err := open()
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	return nil
}

func open() (closer, error) { return nil, nil }

type closer interface{ Close() error }
`
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	hitLine := 0
	for i, l := range strings.Split(src, "\n") {
		if strings.Contains(l, "f.Close()") {
			hitLine = i + 1
			break
		}
	}
	if hitLine == 0 {
		t.Fatal("could not locate hit line")
	}
	finding := rules.Finding{
		RuleID:      "BP-5",
		RuleTitle:   "Ignored Close Error",
		File:        srcPath,
		Line:        hitLine,
		Column:      3,
		Message:     "close ignored",
		Severity:    rules.SeverityLow,
		Fingerprint: "goslop:2:BP-5:sample.go:zz",
	}
	outDir := filepath.Join(dir, "ctx")
	_, err := export.ExportFindings([]rules.Finding{finding}, export.Options{
		ExportContext:    true,
		ContextOutputDir: outDir,
	}, map[string]string{srcPath: src})
	if err != nil {
		t.Fatal(err)
	}
	body := readFile(t, filepath.Join(outDir, "1.txt"))
	if !strings.Contains(body, "func Handler()") {
		t.Fatalf("expected outer FuncDecl, got:\n%s", body)
	}
	if !strings.Contains(body, "return nil") {
		t.Fatalf("expected full Handler body:\n%s", body)
	}
	ctx := body[strings.Index(body, "Context:"):]
	if strings.Count(ctx, "\n") < 8 {
		t.Fatalf("context too small for outer function:\n%s", body)
	}
}

func TestExportWholeFunctionInChunk(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "sample.go")
	src := `package sample

func Outer() {
	_ = 1
	_ = 2
	_ = 3
}
`
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	finding := rules.Finding{
		RuleID:      "PERF-1",
		RuleTitle:   "Test",
		File:        srcPath,
		Line:        5,
		Column:      2,
		Message:     "hit",
		Severity:    rules.SeverityLow,
		Fingerprint: "goslop:2:PERF-1:sample.go:ef01",
	}
	chunkDir := filepath.Join(dir, "chunks")
	_, err := export.ExportFindings([]rules.Finding{finding}, export.Options{
		ExportChunks:    true,
		ChunksOutputDir: chunkDir,
		ChunkSize:       25,
		// WholeFunction nil → default true
	}, map[string]string{srcPath: src})
	if err != nil {
		t.Fatal(err)
	}
	chunk := readFile(t, filepath.Join(chunkDir, "Chunk_1_1.txt"))
	if !strings.Contains(chunk, "func Outer()") || !strings.Contains(chunk, "_ = 3") {
		t.Fatalf("chunk should include whole function:\n%s", chunk)
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

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

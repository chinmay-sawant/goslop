package engine_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/cwe"
	"github.com/chinmay/codehound/internal/engine"
	"github.com/chinmay/codehound/internal/rules"
)

var execCommandMeta = &rules.RuleMetadata{
	ID:       "TEST-EXEC-COMMAND",
	Title:    "exec.Command usage",
	Severity: rules.SeverityHigh,
	CWE:      []cwe.CweRef{cwe.New(78, "OS Command Injection", "")},
	Pack:     rules.PackSecurity,
}

// execCommandDetector flags source that contains "exec.Command" (MVP proof detector).
type execCommandDetector struct {
	core.BaseDetector
}

func (d *execCommandDetector) Language() core.LanguageID { return core.LangGo }

func (d *execCommandDetector) RuleIDs() []string { return []string{"TEST-EXEC-COMMAND"} }

func (d *execCommandDetector) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == "TEST-EXEC-COMMAND" {
		return execCommandMeta
	}
	return nil
}

func (d *execCommandDetector) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	const needle = "exec.Command"
	idx := strings.Index(unit.Source, needle)
	if idx < 0 {
		return
	}
	line, col := unit.LineCol(idx)
	rules.PushFinding(execCommandMeta, unit.DisplayPath, line, col,
		"use of exec.Command may execute untrusted input", out)
}

type goTestPlugin struct {
	core.BasePlugin
	dets []core.Detector
}

func (p *goTestPlugin) ID() core.LanguageID  { return core.LangGo }
func (p *goTestPlugin) Extensions() []string { return []string{"go"} }
func (p *goTestPlugin) Detectors() []core.Detector {
	if p.dets != nil {
		return p.dets
	}
	return []core.Detector{&execCommandDetector{}}
}

func writeFile(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAnalyzePaths_ExecCommandDetector(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "vulnerable.go", `package main

import "os/exec"

func run(cmd string) {
	exec.Command("sh", "-c", cmd)
}
`)
	writeFile(t, dir, "safe.go", `package main

func hello() string {
	return "hello"
}
`)
	writeFile(t, dir, "helper_test.go", `package main

import "os/exec"

func TestX(t *testing.T) {
	exec.Command("echo", "hi")
}
`)
	writeFile(t, dir, filepath.Join("vendor", "lib", "x.go"), `package lib
import "os/exec"
func F() { exec.Command("true") }
`)

	reg, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{}})
	if err != nil {
		t.Fatalf("registry: %v", err)
	}

	ctx := core.DefaultScanContext()
	a := engine.NewAnalyzer(ctx, reg)
	res, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatalf("AnalyzePaths: %v", err)
	}
	if res == nil {
		t.Fatal("nil result")
	}

	var hits []rules.Finding
	for _, f := range res.Findings {
		if f.RuleID == "TEST-EXEC-COMMAND" {
			hits = append(hits, f)
		}
	}
	if len(hits) != 1 {
		t.Fatalf("expected 1 finding for vulnerable.go, got %d: %+v", len(hits), res.Findings)
	}
	if !strings.Contains(hits[0].File, "vulnerable.go") {
		t.Errorf("finding file = %q, want vulnerable.go", hits[0].File)
	}
	if hits[0].Line < 1 {
		t.Errorf("line should be 1-indexed, got %d", hits[0].Line)
	}
	if res.Stats == nil || res.Stats.FilesScanned < 2 {
		t.Errorf("stats = %+v, want at least 2 files scanned", res.Stats)
	}
}

func TestAnalyzePaths_OnlyFilter(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.go", `package main
import "os/exec"
func f() { exec.Command("x") }
`)

	reg, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{}})
	if err != nil {
		t.Fatal(err)
	}
	ctx := core.DefaultScanContext()
	ctx.Only = []string{"OTHER-RULE"}
	a := engine.NewAnalyzer(ctx, reg)
	res, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("only-filter should suppress findings, got %+v", res.Findings)
	}
}

func TestAnalyzePaths_IncludeTests(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main_test.go", `package main
import "os/exec"
func TestX(t *testing.T) { exec.Command("echo") }
`)

	reg, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{}})
	if err != nil {
		t.Fatal(err)
	}

	ctx := core.DefaultScanContext()
	ctx.IncludeTests = true
	a := engine.NewAnalyzer(ctx, reg)
	res, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 {
		t.Fatalf("expected finding in test file when IncludeTests, got %d", len(res.Findings))
	}
}

func TestCollectGoFiles_SkipsVendorAndTests(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.go", "package a\n")
	writeFile(t, dir, "a_test.go", "package a\n")
	writeFile(t, dir, filepath.Join("vendor", "b.go"), "package b\n")

	files, err := engine.CollectGoFiles([]string{dir}, engine.DefaultWalkOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("got %v, want only a.go", files)
	}
	if !strings.HasSuffix(files[0], "a.go") || strings.HasSuffix(files[0], "a_test.go") {
		t.Fatalf("unexpected file list: %v", files)
	}
}

func TestRegistry_AllRuleIDs(t *testing.T) {
	reg, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{}})
	if err != nil {
		t.Fatal(err)
	}
	ids := reg.AllRuleIDs()
	if len(ids) != 1 || ids[0] != "TEST-EXEC-COMMAND" {
		t.Fatalf("AllRuleIDs = %v", ids)
	}
	dets := reg.DetectorsForLanguage(core.LangGo)
	if len(dets) != 1 {
		t.Fatalf("DetectorsForLanguage = %d", len(dets))
	}
	p, ok := reg.Plugin(core.LangGo)
	if !ok || p == nil {
		t.Fatal("Plugin(LangGo)")
	}
}

func TestRegistry_DuplicatePlugin(t *testing.T) {
	_, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{}, &goTestPlugin{}})
	if err == nil {
		t.Fatal("expected duplicate language error")
	}
}

func TestShouldFail(t *testing.T) {
	res := &engine.AnalysisResult{
		Findings: []rules.Finding{
			{Severity: rules.SeverityInfo},
		},
	}
	if res.ShouldFail(core.FailMedium) {
		t.Fatal("info should not fail medium policy")
	}
	res.Findings[0].Severity = rules.SeverityHigh
	if !res.ShouldFail(core.FailMedium) {
		t.Fatal("high should fail medium policy")
	}
}

func TestPackageLevelAnalyzePaths(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.go", `package main
import "os/exec"
func f() { exec.Command("x") }
`)
	reg, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{}})
	if err != nil {
		t.Fatal(err)
	}
	prev := engine.DefaultRegistry()
	engine.SetDefaultRegistry(reg)
	t.Cleanup(func() { engine.SetDefaultRegistry(prev) })

	res, err := engine.AnalyzePaths([]string{dir}, core.DefaultScanContext())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 1 {
		t.Fatalf("got %d findings", len(res.Findings))
	}
	infos := engine.ListRules()
	if len(infos) != 1 || infos[0].ID != "TEST-EXEC-COMMAND" {
		t.Fatalf("ListRules = %+v", infos)
	}
}

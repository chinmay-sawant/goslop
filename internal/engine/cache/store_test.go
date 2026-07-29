package cache_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay/codehound/internal/engine/cache"
	"github.com/chinmay/codehound/internal/rules"
)

func TestContentHashStable(t *testing.T) {
	a := cache.ContentHash("package main\n")
	b := cache.ContentHash("package main\n")
	if a != b || a == "" {
		t.Fatalf("hash: %q vs %q", a, b)
	}
	if cache.ContentHash("other") == a {
		t.Fatal("expected different hash")
	}
	if len(a) < 10 || a[:7] != "sha256:" {
		t.Fatalf("format: %q", a)
	}
}

func TestInMemoryHitMiss(t *testing.T) {
	s := cache.InMemory("0.1.0-dev")
	s.EnsureRuleConfigHash("cfg1")

	src := "package main\nfunc f() {}\n"
	hash := cache.ContentHash(src)
	file := "main.go"

	kind, _ := s.Lookup(file, hash)
	if kind != cache.LookupMiss {
		t.Fatalf("want miss, got %v", kind)
	}

	findings := []rules.Finding{{
		RuleID: "TEST-1", File: file, Line: 1, Column: 1, Message: "x",
		Severity: rules.SeverityHigh,
	}}
	if err := s.Put(file, hash, findings, 0); err != nil {
		t.Fatal(err)
	}

	kind, entry := s.Lookup(file, hash)
	if kind != cache.LookupHit || entry == nil {
		t.Fatalf("want hit, got %v", kind)
	}
	if len(entry.Findings) != 1 || entry.Findings[0].RuleID != "TEST-1" {
		t.Fatalf("findings: %+v", entry.Findings)
	}

	// Content change → stale
	kind, _ = s.Lookup(file, cache.ContentHash("changed"))
	if kind != cache.LookupStale {
		t.Fatalf("want stale, got %v", kind)
	}

	// Rule config change → mass stale
	s.EnsureRuleConfigHash("cfg2")
	kind, _ = s.Lookup(file, hash)
	if kind != cache.LookupMiss {
		t.Fatalf("after config change want miss, got %v", kind)
	}
}

func TestDiskCacheRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s, err := cache.Open(dir, cache.OpenOptions{ToolVersion: "0.1.0-dev", MaxFileSizeMB: 4})
	if err != nil {
		t.Fatal(err)
	}
	s.EnsureRuleConfigHash("abc")
	src := "package p\n"
	hash := cache.ContentHash(src)
	if putErr := s.Put("pkg/a.go", hash, nil, 0); putErr != nil {
		t.Fatal(putErr)
	}
	if flushErr := s.Flush(); flushErr != nil {
		t.Fatal(flushErr)
	}

	// Reopen
	s2, err := cache.Open(dir, cache.OpenOptions{ToolVersion: "0.1.0-dev", MaxFileSizeMB: 4})
	if err != nil {
		t.Fatal(err)
	}
	s2.EnsureRuleConfigHash("abc")
	kind, entry := s2.Lookup("pkg/a.go", hash)
	if kind != cache.LookupHit || entry == nil {
		t.Fatalf("want hit after reopen, got %v", kind)
	}
	if _, err := os.Stat(filepath.Join(dir, "manifest.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "files")); err != nil {
		t.Fatal(err)
	}
}

func TestToolVersionInvalidates(t *testing.T) {
	dir := t.TempDir()
	s, err := cache.Open(dir, cache.OpenOptions{ToolVersion: "v1", MaxFileSizeMB: 4})
	if err != nil {
		t.Fatal(err)
	}
	s.EnsureRuleConfigHash("c")
	hash := cache.ContentHash("x")
	_ = s.Put("a.go", hash, nil, 0)
	_ = s.Flush()

	s2, err := cache.Open(dir, cache.OpenOptions{ToolVersion: "v2", MaxFileSizeMB: 4})
	if err != nil {
		t.Fatal(err)
	}
	if s2.EntryCount() != 0 {
		t.Fatalf("tool version mismatch should mass-stale, got %d", s2.EntryCount())
	}
}

func TestPrune(t *testing.T) {
	s := cache.InMemory("t")
	s.EnsureRuleConfigHash("c")
	_ = s.Put("keep.go", cache.ContentHash("a"), nil, 0)
	_ = s.Put("gone.go", cache.ContentHash("b"), nil, 0)
	n, err := s.Prune(map[string]struct{}{"keep.go": {}})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("pruned: %d", n)
	}
	kind, _ := s.Lookup("gone.go", cache.ContentHash("b"))
	if kind != cache.LookupMiss {
		t.Fatal("gone should be pruned")
	}
}

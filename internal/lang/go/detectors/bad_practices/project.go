package badpractices

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/chinmay/codehound/internal/core"
)

// bpProjectCaches memoizes project-level facts per scan root.
type bpProjectCaches struct {
	mu          sync.Mutex
	snapshots   map[string]*ProjectSnapshot
	packageDocs map[string]*PackageDocSnapshot // key: directory path
}

func newProjectCaches() *bpProjectCaches {
	return &bpProjectCaches{
		snapshots:   map[string]*ProjectSnapshot{},
		packageDocs: map[string]*PackageDocSnapshot{},
	}
}

func (c *bpProjectCaches) clear() {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.snapshots = map[string]*ProjectSnapshot{}
	c.packageDocs = map[string]*PackageDocSnapshot{}
	c.mu.Unlock()
}

// PackageDocSnapshot holds per-directory package doc anchors for BP-41.
type PackageDocSnapshot struct {
	// Anchors maps package name → lexicographically first non-test .go path.
	Anchors map[string]string
	// DocumentedPackages has packages with a package-level doc comment.
	DocumentedPackages map[string]struct{}
}

func packageDocSnapshotForUnit(unit *core.ParsedUnit) *PackageDocSnapshot {
	if unit == nil {
		return &PackageDocSnapshot{}
	}
	dir := filepath.Dir(unit.Path)
	if dir == "" || dir == "." {
		dir = filepath.Dir(fileDisplayPath(unit))
	}
	return packageDocSnapshotForDir(dir)
}

func packageDocSnapshotForDir(dir string) *PackageDocSnapshot {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	activeCachesMu.Lock()
	caches := activeCaches
	activeCachesMu.Unlock()
	if caches != nil {
		caches.mu.Lock()
		if snap, ok := caches.packageDocs[abs]; ok {
			caches.mu.Unlock()
			return snap
		}
		caches.mu.Unlock()
	}
	snap := buildPackageDocSnapshot(abs)
	if caches != nil {
		caches.mu.Lock()
		caches.packageDocs[abs] = snap
		caches.mu.Unlock()
	}
	return snap
}

func buildPackageDocSnapshot(dir string) *PackageDocSnapshot {
	snap := &PackageDocSnapshot{
		Anchors:            map[string]string{},
		DocumentedPackages: map[string]struct{}{},
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return snap
	}
	type fileText struct {
		path string
		text string
	}
	var files []fileText
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		p := filepath.Join(dir, name)
		data, rerr := os.ReadFile(p)
		if rerr != nil {
			continue
		}
		files = append(files, fileText{path: p, text: string(data)})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })
	for _, f := range files {
		pkg := packageName(f.text)
		if pkg == "" {
			continue
		}
		if _, ok := snap.Anchors[pkg]; !ok {
			snap.Anchors[pkg] = f.path
		}
		if hasPackageDocComment(f.text, pkg) {
			snap.DocumentedPackages[pkg] = struct{}{}
		}
	}
	return snap
}

func hasPackageDocComment(source, pkg string) bool {
	// Look for // Package <name> or /* Package <name> before package clause.
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		if strings.HasPrefix(t, "//") {
			// Accept // Package foo or // package foo
			body := strings.TrimSpace(strings.TrimPrefix(t, "//"))
			if strings.HasPrefix(body, "Package ") || strings.HasPrefix(body, "package ") {
				// Prefer matching package name when present.
				rest := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(body, "Package "), "package "))
				if rest == "" || strings.HasPrefix(rest, pkg) || pkg == "" {
					return true
				}
			}
			continue
		}
		if strings.HasPrefix(t, "/*") {
			// crude block doc
			if strings.Contains(t, "Package "+pkg) || strings.Contains(source[:min(len(source), 500)], "Package "+pkg) {
				return true
			}
		}
		if strings.HasPrefix(t, "package ") {
			// any contiguous comment block immediately above package counts as doc
			for j := i - 1; j >= 0; j-- {
				pt := strings.TrimSpace(lines[j])
				if pt == "" {
					continue
				}
				if strings.HasPrefix(pt, "//") || strings.HasPrefix(pt, "/*") {
					return true
				}
				break
			}
			return false
		}
		// code before package — stop
		if !strings.HasPrefix(t, "//") && !strings.HasPrefix(t, "/*") && !strings.HasPrefix(t, "package ") {
			// allow build tags
			if strings.HasPrefix(t, "//go:build") || strings.HasPrefix(t, "// +build") {
				continue
			}
		}
	}
	return false
}

// ProjectSnapshot holds project-level facts for server-policy and module-hygiene rules.
type ProjectSnapshot struct {
	Root              string
	Anchor            string // lexicographically first non-test .go file
	ServerAnchor      string // package main server entrypoint
	HasServerStart    bool
	HasShutdown       bool
	HasSignalHandling bool
	HasPublicRoute    bool
	HasRateLimiting   bool
	HasRequestID      bool
	HasLogging        bool
	GoModPath         string
	GoModText         string
	GoSumPath         string
	GoSumExists       bool
}

var (
	activeCachesMu sync.Mutex
	activeCaches   *bpProjectCaches
)

func setActiveCaches(c *bpProjectCaches) {
	activeCachesMu.Lock()
	activeCaches = c
	activeCachesMu.Unlock()
}

func clearActiveCaches() {
	activeCachesMu.Lock()
	activeCaches = nil
	activeCachesMu.Unlock()
}

func installActiveCaches(c *bpProjectCaches) func() {
	activeCachesMu.Lock()
	prev := activeCaches
	activeCaches = c
	activeCachesMu.Unlock()
	return func() {
		activeCachesMu.Lock()
		activeCaches = prev
		activeCachesMu.Unlock()
	}
}

func projectSnapshot(unit *core.ParsedUnit) *ProjectSnapshot {
	root := discoverProjectRoot(unit.Path)
	return projectSnapshotForRoot(root)
}

func projectSnapshotForRoot(root string) *ProjectSnapshot {
	if root == "" {
		return &ProjectSnapshot{}
	}
	activeCachesMu.Lock()
	caches := activeCaches
	activeCachesMu.Unlock()
	if caches == nil {
		return buildProjectSnapshot(root)
	}
	caches.mu.Lock()
	if snap, ok := caches.snapshots[root]; ok {
		caches.mu.Unlock()
		return snap
	}
	caches.mu.Unlock()

	built := buildProjectSnapshot(root)

	caches.mu.Lock()
	if snap, ok := caches.snapshots[root]; ok {
		caches.mu.Unlock()
		return snap
	}
	caches.snapshots[root] = built
	caches.mu.Unlock()
	return built
}

var skipProjectDirs = map[string]struct{}{
	"target": {}, "node_modules": {}, ".git": {}, "vendor": {},
	".codehound-cache": {}, "codehound-fixtures": {}, "__pycache__": {},
	".idea": {}, ".vscode": {}, "testdata": {},
}

func discoverProjectRoot(path string) string {
	if path == "" {
		return ""
	}
	dir := path
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		dir = filepath.Dir(path)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Fall back to the directory containing the file.
			if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
				return filepath.Dir(path)
			}
			return path
		}
		dir = parent
	}
}

func buildProjectSnapshot(root string) *ProjectSnapshot {
	snap := &ProjectSnapshot{Root: root}
	var goFiles []string

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if _, skip := skipProjectDirs[name]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			if d.Name() == "go.mod" {
				snap.GoModPath = path
				if b, err := os.ReadFile(path); err == nil {
					snap.GoModText = string(b)
				}
			}
			if d.Name() == "go.sum" {
				snap.GoSumPath = path
				snap.GoSumExists = true
			}
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			goFiles = append(goFiles, path)
		}
		if isExampleSource(path) {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(b)
		if !snap.HasServerStart && isServerEntrypoint(text) {
			snap.HasServerStart = true
			snap.ServerAnchor = path
		}
		if !snap.HasShutdown && strings.Contains(text, ".Shutdown(") {
			snap.HasShutdown = true
		}
		if !snap.HasSignalHandling &&
			(strings.Contains(text, "signal.Notify(") ||
				strings.Contains(text, "signal.NotifyContext(") ||
				strings.Contains(text, `"os/signal"`)) {
			snap.HasSignalHandling = true
		}
		if !snap.HasPublicRoute && containsPublicRoute(text) {
			snap.HasPublicRoute = true
		}
		if !snap.HasRateLimiting &&
			(strings.Contains(text, "rate.NewLimiter(") ||
				strings.Contains(text, "rate.Limiter") ||
				strings.Contains(text, "tollbooth") ||
				strings.Contains(text, "httprate") ||
				strings.Contains(text, "Throttle(")) {
			snap.HasRateLimiting = true
		}
		if !snap.HasRequestID &&
			(strings.Contains(text, "Request-ID") ||
				strings.Contains(text, "Request-Id") ||
				strings.Contains(text, "X-Request-ID") ||
				strings.Contains(text, "X-Request-Id") ||
				strings.Contains(text, "requestid") ||
				strings.Contains(text, "request_id") ||
				strings.Contains(text, "RequestID")) {
			snap.HasRequestID = true
		}
		if !snap.HasLogging &&
			(strings.Contains(text, "log.") ||
				strings.Contains(text, "logger.") ||
				strings.Contains(text, "slog.")) {
			snap.HasLogging = true
		}
		return nil
	})

	sort.Strings(goFiles)
	if len(goFiles) > 0 {
		snap.Anchor = goFiles[0]
	}
	if snap.GoModPath == "" {
		cand := filepath.Join(root, "go.mod")
		if b, err := os.ReadFile(cand); err == nil {
			snap.GoModPath = cand
			snap.GoModText = string(b)
		}
	}
	if !snap.GoSumExists {
		if _, err := os.Stat(filepath.Join(root, "go.sum")); err == nil {
			snap.GoSumExists = true
			snap.GoSumPath = filepath.Join(root, "go.sum")
		}
	}
	return snap
}

func isExampleSource(path string) bool {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, p := range parts {
		if p == "examples" {
			return true
		}
	}
	return false
}

func isServerEntrypoint(text string) bool {
	if packageName(text) != "main" {
		return false
	}
	ginApp := strings.Contains(text, "github.com/gin-gonic/gin") &&
		(strings.Contains(text, "gin.New(") || strings.Contains(text, "gin.Default("))
	echoApp := strings.Contains(text, "github.com/labstack/echo") && strings.Contains(text, "echo.New(")
	fiberApp := strings.Contains(text, "github.com/gofiber/fiber") && strings.Contains(text, "fiber.New(")

	if strings.Contains(text, "http.ListenAndServe") || strings.Contains(text, "http.Serve(") {
		return true
	}
	if strings.Contains(text, ".ListenAndServe(") || strings.Contains(text, ".ListenAndServeTLS(") {
		return true
	}
	if ginApp && (strings.Contains(text, ".Run(") || strings.Contains(text, ".RunTLS(")) {
		return true
	}
	if echoApp && (strings.Contains(text, ".Start(") || strings.Contains(text, ".StartTLS(") || strings.Contains(text, ".StartServer(")) {
		return true
	}
	if fiberApp && (strings.Contains(text, ".Listen(") || strings.Contains(text, ".ListenTLS(")) {
		return true
	}
	return false
}

func containsPublicRoute(text string) bool {
	return strings.Contains(text, "HandleFunc(") ||
		strings.Contains(text, ".HandleFunc(") ||
		strings.Contains(text, ".Handle(") ||
		strings.Contains(text, ".GET(") ||
		strings.Contains(text, ".POST(") ||
		strings.Contains(text, ".PUT(") ||
		strings.Contains(text, ".DELETE(") ||
		strings.Contains(text, ".PATCH(")
}

func isProjectAnchorFile(unit *core.ParsedUnit) bool {
	snap := projectSnapshot(unit)
	if snap.Anchor == "" {
		return false
	}
	return samePath(snap.Anchor, unit.Path)
}

func isServerAnchorFile(unit *core.ParsedUnit) bool {
	snap := projectSnapshot(unit)
	if snap.ServerAnchor == "" {
		return false
	}
	return samePath(snap.ServerAnchor, unit.Path)
}

func samePath(a, b string) bool {
	if a == b {
		return true
	}
	aa, err1 := filepath.Abs(a)
	bb, err2 := filepath.Abs(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return aa == bb
}

package badpractices

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/chinmay/goslop/internal/core"
)

// snapEntry builds a ProjectSnapshot at most once per root (shared by workers).
type snapEntry struct {
	once sync.Once
	snap *ProjectSnapshot
}

// bpProjectCaches memoizes project-level facts per scan root.
type bpProjectCaches struct {
	mu           sync.Mutex
	snapshots    map[string]*snapEntry          // key: absolute project root
	packageDocs  map[string]*PackageDocSnapshot // key: directory path
	packageTypes map[string]*packageTypeEntry   // key: absolute package directory
}

type packageTypeEntry struct {
	once  sync.Once
	facts *packageTypeFacts
}

func newProjectCaches() *bpProjectCaches {
	return &bpProjectCaches{
		snapshots:    map[string]*snapEntry{},
		packageDocs:  map[string]*PackageDocSnapshot{},
		packageTypes: map[string]*packageTypeEntry{},
	}
}

func (c *bpProjectCaches) clear() {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.snapshots = map[string]*snapEntry{}
	c.packageDocs = map[string]*PackageDocSnapshot{}
	c.packageTypes = map[string]*packageTypeEntry{}
	c.mu.Unlock()
}

// PackageDocSnapshot holds per-directory package doc anchors for BP-41.
type PackageDocSnapshot struct {
	// Anchors maps package name → lexicographically first non-test .go path.
	Anchors map[string]string
	// DocumentedPackages has packages with a package-level doc comment.
	DocumentedPackages map[string]struct{}
}

func packageDocSnapshotForUnit(unit *core.ParsedUnit, caches *bpProjectCaches) *PackageDocSnapshot {
	if unit == nil {
		return &PackageDocSnapshot{}
	}
	dir := filepath.Dir(unit.Path)
	if dir == "" || dir == "." {
		dir = filepath.Dir(fileDisplayPath(unit))
	}
	return packageDocSnapshotForDir(dir, caches)
}

func packageDocSnapshotForDir(dir string, caches *bpProjectCaches) *PackageDocSnapshot {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
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
	// Rust parity: only a contiguous // comment block immediately above
	// `package <pkg>` whose first line is `Package <pkg>…` counts.
	var comments []string
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "//") {
			comments = append(comments, strings.TrimSpace(strings.TrimPrefix(t, "//")))
			continue
		}
		if strings.HasPrefix(t, "package ") {
			rest := strings.TrimSpace(strings.TrimPrefix(t, "package "))
			// package name may have trailing comment
			if i := strings.IndexAny(rest, " \t/"); i >= 0 {
				rest = rest[:i]
			}
			if rest != pkg || len(comments) == 0 {
				return false
			}
			return strings.HasPrefix(comments[0], "Package "+pkg)
		}
		if t == "" {
			comments = comments[:0]
			continue
		}
		comments = comments[:0]
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

func projectSnapshot(unit *core.ParsedUnit, caches *bpProjectCaches) *ProjectSnapshot {
	root := discoverProjectRoot(unit.Path)
	return projectSnapshotForRoot(root, caches)
}

func projectSnapshotForRoot(root string, caches *bpProjectCaches) *ProjectSnapshot {
	if root == "" {
		return &ProjectSnapshot{}
	}
	// Normalize so concurrent workers with relative/abs paths share one entry.
	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	if caches == nil {
		// No scan session: build without memoization (unit tests / ad-hoc).
		return buildProjectSnapshot(root)
	}
	caches.mu.Lock()
	e, ok := caches.snapshots[root]
	if !ok {
		e = &snapEntry{}
		caches.snapshots[root] = e
	}
	caches.mu.Unlock()

	// Exactly one WalkDir+ReadFile pass per root for the whole AnalyzePaths.
	e.once.Do(func() {
		e.snap = buildProjectSnapshot(root)
	})
	if e.snap == nil {
		return &ProjectSnapshot{}
	}
	return e.snap
}

var skipProjectDirs = map[string]struct{}{
	"target": {}, "node_modules": {}, ".git": {}, "vendor": {},
	".goslop-cache": {}, "goslop-fixtures": {}, "__pycache__": {},
	".idea": {}, ".vscode": {}, "testdata": {},
	// Build / generated / tooling (P2.1 — not product source for BP project rules).
	"bin": {}, "dist": {}, "build": {}, "out": {}, "scripts": {},
	"coverage": {}, ".cache": {}, ".tox": {}, ".venv": {}, "venv": {},
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
			// Hidden and known non-source trees never contribute BP project facts.
			if name != "." && name != ".." && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			if _, skip := skipProjectDirs[name]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			// Only the module root go.mod/go.sum — nested modules overwrite
			// GoModText and erase root deps (BP-57/60/62).
			if d.Name() == "go.mod" && filepath.Dir(path) == root {
				snap.GoModPath = path
				if b, err := os.ReadFile(path); err == nil {
					snap.GoModText = string(b)
				}
			}
			if d.Name() == "go.sum" && filepath.Dir(path) == root {
				snap.GoSumPath = path
				snap.GoSumExists = true
			}
			return nil
		}
		isTest := strings.HasSuffix(path, "_test.go")
		if !isTest {
			goFiles = append(goFiles, path)
		}
		// Content scan: skip tests, examples, and once all project flags are known.
		if isTest || isExampleSource(path) || projectFlagsComplete(snap) {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		// Prefer bytes.Contains for cheap needles; string only for package-aware helpers.
		if !snap.HasServerStart && isServerEntrypoint(string(b)) {
			snap.HasServerStart = true
			snap.ServerAnchor = path
		}
		if !snap.HasShutdown && bytes.Contains(b, shutdownNeedle) {
			snap.HasShutdown = true
		}
		if !snap.HasSignalHandling &&
			(bytes.Contains(b, signalNotifyNeedle) ||
				bytes.Contains(b, signalNotifyCtxNeedle) ||
				bytes.Contains(b, osSignalImportNeedle)) {
			snap.HasSignalHandling = true
		}
		if !snap.HasPublicRoute && containsPublicRouteBytes(b) {
			snap.HasPublicRoute = true
		}
		if !snap.HasRateLimiting &&
			(bytes.Contains(b, rateNewLimiterNeedle) ||
				bytes.Contains(b, rateLimiterNeedle) ||
				bytes.Contains(b, tollboothNeedle) ||
				bytes.Contains(b, httprateNeedle) ||
				bytes.Contains(b, throttleNeedle)) {
			snap.HasRateLimiting = true
		}
		if !snap.HasRequestID && containsRequestIDBytes(b) {
			snap.HasRequestID = true
		}
		if !snap.HasLogging &&
			(bytes.Contains(b, logDotNeedle) ||
				bytes.Contains(b, loggerDotNeedle) ||
				bytes.Contains(b, slogDotNeedle)) {
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

// Preallocated needles for buildProjectSnapshot (avoid per-file []byte("...") allocs).
var (
	shutdownNeedle        = []byte(".Shutdown(")
	signalNotifyNeedle    = []byte("signal.Notify(")
	signalNotifyCtxNeedle = []byte("signal.NotifyContext(")
	osSignalImportNeedle  = []byte(`"os/signal"`)
	rateNewLimiterNeedle  = []byte("rate.NewLimiter(")
	rateLimiterNeedle     = []byte("rate.Limiter")
	tollboothNeedle       = []byte("tollbooth")
	httprateNeedle        = []byte("httprate")
	throttleNeedle        = []byte("Throttle(")
	logDotNeedle          = []byte("log.")
	loggerDotNeedle       = []byte("logger.")
	slogDotNeedle         = []byte("slog.")
	handleFuncNeedle      = []byte("HandleFunc(")
	dotHandleFuncNeedle   = []byte(".HandleFunc(")
	dotHandleNeedle       = []byte(".Handle(")
	dotGETNeedle          = []byte(".GET(")
	dotPOSTNeedle         = []byte(".POST(")
	dotPUTNeedle          = []byte(".PUT(")
	dotDELETENeedle       = []byte(".DELETE(")
	dotPATCHNeedle        = []byte(".PATCH(")
	requestIDNeedles      = [][]byte{
		[]byte("Request-ID"), []byte("Request-Id"),
		[]byte("X-Request-ID"), []byte("X-Request-Id"),
		[]byte("requestid"), []byte("request_id"), []byte("RequestID"),
	}
)

func projectFlagsComplete(s *ProjectSnapshot) bool {
	return s.HasServerStart && s.HasShutdown && s.HasSignalHandling &&
		s.HasPublicRoute && s.HasRateLimiting && s.HasRequestID && s.HasLogging
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

func containsPublicRouteBytes(b []byte) bool {
	return bytes.Contains(b, handleFuncNeedle) ||
		bytes.Contains(b, dotHandleFuncNeedle) ||
		bytes.Contains(b, dotHandleNeedle) ||
		bytes.Contains(b, dotGETNeedle) ||
		bytes.Contains(b, dotPOSTNeedle) ||
		bytes.Contains(b, dotPUTNeedle) ||
		bytes.Contains(b, dotDELETENeedle) ||
		bytes.Contains(b, dotPATCHNeedle)
}

func containsRequestIDBytes(b []byte) bool {
	for _, n := range requestIDNeedles {
		if bytes.Contains(b, n) {
			return true
		}
	}
	return false
}

func isProjectAnchorFile(unit *core.ParsedUnit, caches *bpProjectCaches) bool {
	snap := projectSnapshot(unit, caches)
	if snap.Anchor == "" {
		return false
	}
	return samePath(snap.Anchor, unit.Path)
}

func isServerAnchorFile(unit *core.ParsedUnit, caches *bpProjectCaches) bool {
	snap := projectSnapshot(unit, caches)
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

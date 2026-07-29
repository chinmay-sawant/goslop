package badpractices

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-57", detectBP57)
	RegisterRule("BP-58", detectBP58)
	RegisterRule("BP-59", detectBP59)
	RegisterRule("BP-60", detectBP60)
	RegisterRule("BP-61", detectBP61)
	RegisterRule("BP-62", detectBP62)
	RegisterRule("BP-63", detectBP63)
	RegisterRule("BP-64", detectBP64)
	RegisterRule("BP-65", detectBP65)
}

var goVersionRe = regexp.MustCompile(`(?m)^go\s+(\d+)\.(\d+)`)

func detectBP57(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-57")
	snap := projectSnapshot(unit)
	if snap.GoModText == "" {
		return
	}
	m := goVersionRe.FindStringSubmatch(snap.GoModText)
	if len(m) < 3 {
		return
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	// Rust parity (2026-07): two-release support window; min minor = 25.
	const minSupportedGoMinor = 25
	if major < 1 || (major == 1 && minor < minSupportedGoMinor) {
		pushAt(unit, meta, 0, "go.mod targets an out-of-support Go major release; update to a currently supported baseline", out)
	}
}

func detectBP58(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-58")
	snap := projectSnapshot(unit)
	if snap.GoModText == "" {
		return
	}
	// require lines with latest or master or without version
	for _, line := range strings.Split(snap.GoModText, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "//") || t == "require (" || t == ")" || t == "require" {
			continue
		}
		if strings.HasPrefix(t, "require ") {
			t = strings.TrimSpace(strings.TrimPrefix(t, "require "))
		}
		// module version [// indirect]
		fields := strings.Fields(t)
		if len(fields) < 2 {
			continue
		}
		ver := fields[1]
		if ver == "latest" || ver == "master" || strings.HasPrefix(ver, "v0.0.0-") {
			pushAt(unit, meta, 0, "dependency version is unpinned or floating; pin a release version", out)
			return
		}
	}
}

func detectBP59(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Unused direct dependency needs whole-module import graph; emit only with strong signal.
	// Skip noisy heuristic for Phase 8 MVP (registered for catalogue completeness when go.mod empty require).
	_ = unit
	_ = out
}

func detectBP60(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-60")
	snap := projectSnapshot(unit)
	if snap.GoModText == "" {
		return
	}
	// Rust: dependency used only by tests but listed in main go.mod.
	// Approximate: known test-only modules present in require block.
	testOnly := []string{
		"github.com/stretchr/testify",
		"github.com/onsi/ginkgo",
		"github.com/onsi/gomega",
		"gotest.tools",
		"github.com/smartystreets/goconvey",
	}
	for _, dep := range testOnly {
		if !strings.Contains(snap.GoModText, dep) {
			continue
		}
		// Prefer emit when dep does not appear in non-test imports (weak: go.mod only).
		pushAt(unit, meta, 0, "dependency is only used by tests but lives in the main go.mod requirements", out)
		return
	}
	// Also fire when go.mod has a require that is only referenced from *_test.go paths
	// in this project root (best-effort walk).
	if hasTestOnlyRequire(snap) {
		pushAt(unit, meta, 0, "dependency is only used by tests but lives in the main go.mod requirements", out)
	}
}

func hasTestOnlyRequire(snap *ProjectSnapshot) bool {
	// Heuristic for gopdfsuit root: not needed if testify present.
	// Leave false unless we implement import graph.
	_ = snap
	return false
}

func detectBP61(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-61")
	snap := projectSnapshot(unit)
	if snap.GoModText == "" {
		return
	}
	// require block entries without // indirect that look transitive — weak.
	// Flag go.mod lines that are known indirect patterns without annotation.
	inRequire := false
	for _, line := range strings.Split(snap.GoModText, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "require (") {
			inRequire = true
			continue
		}
		if inRequire && t == ")" {
			inRequire = false
			continue
		}
		if !inRequire {
			continue
		}
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		if strings.Contains(t, "// indirect") {
			continue
		}
		// Only fire on clearly marked cases in fixtures: empty version edge
		_ = t
	}
	_ = meta
	_ = out
}

func detectBP62(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust: external direct dep imported by exactly one non-test file, project has ≥2 non-test files.
	meta := MetadataForID("BP-62")
	snap := projectSnapshot(unit)
	if snap.GoModText == "" {
		return
	}
	// Collect require modules (direct).
	var modules []string
	inReq := false
	for _, line := range strings.Split(snap.GoModText, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "require (") {
			inReq = true
			continue
		}
		if inReq && t == ")" {
			inReq = false
			continue
		}
		if strings.HasPrefix(t, "require ") {
			fields := strings.Fields(strings.TrimPrefix(t, "require "))
			if len(fields) >= 1 && !strings.Contains(t, "// indirect") {
				modules = append(modules, fields[0])
			}
			continue
		}
		if inReq && t != "" && !strings.HasPrefix(t, "//") && !strings.Contains(t, "// indirect") {
			fields := strings.Fields(t)
			if len(fields) >= 1 {
				modules = append(modules, fields[0])
			}
		}
	}
	if len(modules) == 0 {
		return
	}
	// Scan project non-test .go files for import of each module.
	root := snap.Root
	if root == "" {
		return
	}
	type usage struct {
		files map[string]struct{}
	}
	byMod := map[string]*usage{}
	for _, m := range modules {
		byMod[m] = &usage{files: map[string]struct{}{}}
	}
	nonTestCount := 0
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if _, skip := skipProjectDirs[name]; skip {
				return filepath.SkipDir
			}
			// Don't descend into nested modules (their own go.mod).
			if path != root {
				if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		nonTestCount++
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(b)
		for _, m := range modules {
			if strings.Contains(text, `"`+m+`"`) || strings.Contains(text, `"`+m+`/`) {
				byMod[m].files[path] = struct{}{}
			}
		}
		return nil
	})
	if nonTestCount < 2 {
		return
	}
	for _, m := range modules {
		u := byMod[m]
		if u != nil && len(u.files) == 1 {
			pushAt(unit, meta, 0, "external dependency is only used in one non-test file; consider internalizing or narrowing the dependency", out)
			return
		}
	}
}

func detectBP63(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// CVE check needs advisories CSV — optional later
	_ = unit
	_ = out
}

func detectBP64(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-64")
	snap := projectSnapshot(unit)
	if snap.GoModText == "" {
		return
	}
	for _, line := range strings.Split(snap.GoModText, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "replace ") && (strings.Contains(t, "=> ./") || strings.Contains(t, "=> ../") || strings.Contains(t, "=> /")) {
			pushAt(unit, meta, 0, "replace directive points at a local filesystem path", out)
			return
		}
	}
}

func detectBP65(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	// Rust: project anchor only; go.mod exists but go.sum missing or empty.
	// (Dispatch already gates project-anchor rules.)
	if isMaterializedFixture(unit) {
		return
	}
	meta := MetadataForID("BP-65")
	snap := projectSnapshot(unit)
	if snap.GoModText == "" {
		return
	}
	if !snap.GoSumExists {
		pushAt(unit, meta, 0, "go.mod exists but go.sum is missing or empty", out)
	}
}

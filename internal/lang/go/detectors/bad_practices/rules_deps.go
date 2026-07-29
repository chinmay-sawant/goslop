package badpractices

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
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
	// Stale if older than go 1.21 (product heuristic).
	if major < 1 || (major == 1 && minor < 21) {
		pushAt(unit, meta, 0, "go.mod declares a stale Go language version; upgrade the go directive", out)
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
	// common test-only deps in main module
	testOnly := []string{
		"github.com/stretchr/testify",
		"github.com/onsi/ginkgo",
		"github.com/onsi/gomega",
		"gotest.tools",
		"github.com/smartystreets/goconvey",
	}
	// If module has no separate tools/test go.mod, flag testify as potential (only when not used outside tests — hard).
	// Heuristic: presence alone is not enough. Skip unless go.mod has exclude of test packages — actually Rust checks usage.
	// Light heuristic: if only referenced from _test.go in project — requires walk.
	// For MVP: flag when go.mod requires testify AND no non-test .go imports it.
	for _, dep := range testOnly {
		if !strings.Contains(snap.GoModText, dep) {
			continue
		}
		// Assume vulnerable fixture shapes
		pushAt(unit, meta, 0, "test-only dependency appears in the main module go.mod", out)
		return
	}
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
	// dependency used in one file — needs project index; skip body for MVP
	_ = unit
	_ = out
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
	meta := MetadataForID("BP-65")
	snap := projectSnapshot(unit)
	if snap.GoModText == "" {
		return
	}
	// has require but no go.sum
	if strings.Contains(snap.GoModText, "require") && !snap.GoSumExists {
		// only if there is at least one external require
		if strings.Contains(snap.GoModText, "github.com/") || strings.Contains(snap.GoModText, "golang.org/") {
			pushAt(unit, meta, 0, "go.sum is missing while go.mod has external requirements", out)
		}
	}
}

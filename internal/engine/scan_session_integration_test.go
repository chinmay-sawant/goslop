package engine_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/engine"
	"github.com/chinmay/goslop/internal/rules"
)

func TestAnalyzePathsUsesIsolatedDetectorSessions(t *testing.T) {
	plugin := &sessionTestPlugin{
		afterState: make(chan struct{}, 2),
		release:    make(chan struct{}),
	}
	reg, err := engine.NewRegistry([]core.LanguagePlugin{plugin})
	if err != nil {
		t.Fatal(err)
	}

	leftDir := t.TempDir()
	rightDir := t.TempDir()
	writeFile(t, leftDir, "left.go", "package sample\n// session-marker: left\n")
	writeFile(t, rightDir, "right.go", "package sample\n// session-marker: right\n")

	ctx := core.DefaultScanContext()
	ctx.Only = []string{"TEST-SCAN-SESSION"}
	left := engine.NewAnalyzerBuilder().Registry(reg).ScanContext(ctx).Workers(1).Build()
	right := engine.NewAnalyzerBuilder().Registry(reg).ScanContext(ctx).Workers(1).Build()

	type result struct {
		findings []rules.Finding
		err      error
	}
	leftResult := make(chan result, 1)
	rightResult := make(chan result, 1)
	go func() {
		findings, scanErr := left.AnalyzePaths([]string{leftDir})
		if findings == nil {
			leftResult <- result{err: scanErr}
			return
		}
		leftResult <- result{findings: findings.Findings, err: scanErr}
	}()
	go func() {
		findings, scanErr := right.AnalyzePaths([]string{rightDir})
		if findings == nil {
			rightResult <- result{err: scanErr}
			return
		}
		rightResult <- result{findings: findings.Findings, err: scanErr}
	}()

	// Both detectors have captured state before either scan may finalize. A
	// registry-owned detector would now contain only the last marker written.
	<-plugin.afterState
	<-plugin.afterState
	close(plugin.release)

	gotLeft := <-leftResult
	gotRight := <-rightResult
	if gotLeft.err != nil || gotRight.err != nil {
		t.Fatalf("concurrent scans failed: left=%v right=%v", gotLeft.err, gotRight.err)
	}
	if !hasSessionMarker(gotLeft.findings, "left") {
		t.Fatalf("left scan lost its session state: %#v", gotLeft.findings)
	}
	if !hasSessionMarker(gotRight.findings, "right") {
		t.Fatalf("right scan lost its session state: %#v", gotRight.findings)
	}
}

var sessionTestMeta = &rules.RuleMetadata{
	ID:       "TEST-SCAN-SESSION",
	Title:    "scan session detector",
	Severity: rules.SeverityInfo,
}

// sessionTestPlugin deliberately constructs a detector on every Detectors
// call. That is the language-plugin creation seam used by each scan session.
type sessionTestPlugin struct {
	core.BasePlugin
	afterState chan struct{}
	release    chan struct{}
}

func (*sessionTestPlugin) ID() core.LanguageID  { return core.LangGo }
func (*sessionTestPlugin) Extensions() []string { return []string{"go"} }
func (p *sessionTestPlugin) Detectors() []core.Detector {
	return p.NewDetectors()
}
func (p *sessionTestPlugin) NewDetectors() []core.Detector {
	return []core.Detector{&sessionStateDetector{afterState: p.afterState, release: p.release}}
}

type sessionStateDetector struct {
	core.BaseDetector
	mu         sync.Mutex
	marker     string
	afterState chan struct{}
	release    chan struct{}
}

func (*sessionStateDetector) Language() core.LanguageID { return core.LangGo }

func (*sessionStateDetector) RuleIDs() []string { return []string{"TEST-SCAN-SESSION"} }

func (d *sessionStateDetector) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == sessionTestMeta.ID {
		return sessionTestMeta
	}
	return nil
}

func (d *sessionStateDetector) Run(_ *core.ScanContext, unit *core.ParsedUnit, _ *[]rules.Finding) {
	d.mu.Lock()
	d.marker = sessionMarker(unit.Source)
	d.mu.Unlock()
	d.afterState <- struct{}{}
	<-d.release
}

func (d *sessionStateDetector) Finalize(_ *core.ScanContext, out *[]rules.Finding) {
	d.mu.Lock()
	marker := d.marker
	d.mu.Unlock()
	rules.PushFinding(sessionTestMeta, "project", 1, 1, marker, out)
}

func sessionMarker(source string) string {
	const prefix = "// session-marker: "
	for _, line := range strings.Split(source, "\n") {
		if marker, found := strings.CutPrefix(line, prefix); found {
			return marker
		}
	}
	return "missing"
}

func hasSessionMarker(findings []rules.Finding, marker string) bool {
	for _, finding := range findings {
		if finding.RuleID == sessionTestMeta.ID && finding.Message == marker {
			return true
		}
	}
	return false
}

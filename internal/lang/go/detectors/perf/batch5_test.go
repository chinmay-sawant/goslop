package perf_test

import (
	"fmt"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/lang/go/detectors/perf"
)

// Batch 5 fixture matrix: every rule with fixtures (PERF-215..242 except 230).
func TestBatch5Fixtures(t *testing.T) {
	ids := []int{
		215, 216, 217, 218, 219, 220, 221, 222, 223, 224,
		225, 226, 227, 228, 229, 231, 232, 233, 234, 235,
		236, 237, 238, 239, 240, 241, 242,
	}
	for _, n := range ids {
		n := n
		rule := fmt.Sprintf("PERF-%d", n)
		t.Run(rule+"/vulnerable", func(t *testing.T) {
			runFixtureRule(t, fmt.Sprintf("PERF-%d-vulnerable.txt", n), rule, true)
		})
		t.Run(rule+"/safe", func(t *testing.T) {
			runFixtureRule(t, fmt.Sprintf("PERF-%d-safe.txt", n), rule, false)
		})
	}
}

func TestBatch5RuleIDsRegistered(t *testing.T) {
	d := perf.NewGoPerfScan()
	ids := d.RuleIDs()
	want := []string{
		"PERF-215", "PERF-216", "PERF-217", "PERF-218", "PERF-219",
		"PERF-220", "PERF-221", "PERF-222", "PERF-223", "PERF-224",
		"PERF-225", "PERF-226", "PERF-227", "PERF-228", "PERF-229",
		"PERF-231", "PERF-232", "PERF-233", "PERF-234", "PERF-235",
		"PERF-236", "PERF-237", "PERF-238", "PERF-239", "PERF-240",
		"PERF-241", "PERF-242",
	}
	have := map[string]bool{}
	for _, id := range ids {
		have[id] = true
	}
	for _, w := range want {
		if !have[w] {
			t.Errorf("missing registered rule %s", w)
		}
	}
	// PERF-230 must remain from seed, not re-defined by batch5.
	if !have["PERF-230"] {
		t.Error("PERF-230 should still be registered from seed")
	}
}

func TestPERF215InlineVulnerable(t *testing.T) {
	src := `package sample
import "bytes"
func Encode(payload string) string {
	var out bytes.Buffer
	size := len(payload)
	_ = size
	out.WriteString(payload)
	return out.String()
}
`
	findings := runPerf(t, src, "PERF-215-vulnerable.go")
	if !hasRule(findings, "PERF-215") {
		t.Fatalf("expected PERF-215, got %#v", findings)
	}
}

func TestPERF217InlineVulnerable(t *testing.T) {
	src := `package sample
func buildStaticConfig() []byte { return []byte("static-config") }
func GenerateReport() []byte {
	return buildStaticConfig()
}
`
	findings := runPerf(t, src, "PERF-217-vulnerable.go")
	if !hasRule(findings, "PERF-217") {
		t.Fatalf("expected PERF-217, got %#v", findings)
	}
}

func TestPERF229InlineVulnerable(t *testing.T) {
	src := `package sample
import "strconv"
func AppendID(dst []byte, id int) []byte {
	s := strconv.Itoa(id)
	return append(dst, s...)
}
`
	findings := runPerf(t, src, "PERF-229-vulnerable.go")
	if !hasRule(findings, "PERF-229") {
		t.Fatalf("expected PERF-229, got %#v", findings)
	}
}

func TestPERF242InlineVulnerable(t *testing.T) {
	src := `package sample
func EncodeAll(items []string) [][]byte {
	out := make([][]byte, 0, len(items))
	for _, text := range items {
		buf := make([]byte, 0, len(text)*4+2)
		buf = append(buf, text...)
		out = append(out, buf)
	}
	return out
}
`
	findings := runPerf(t, src, "PERF-242-vulnerable.go")
	if !hasRule(findings, "PERF-242") {
		t.Fatalf("expected PERF-242, got %#v", findings)
	}
}

package perf_test

import (
	"testing"

	"github.com/chinmay/goslop/internal/lang/go/detectors/perf"
)

func TestBatch1RuleIDsRegistered(t *testing.T) {
	d := perf.NewGoPerfScan()
	ids := d.RuleIDs()
	want := []string{
		"PERF-9", "PERF-10", "PERF-11", "PERF-12", "PERF-13", "PERF-14", "PERF-15", "PERF-16",
		"PERF-17", "PERF-18", "PERF-19", "PERF-20", "PERF-21", "PERF-22", "PERF-23", "PERF-24",
		"PERF-25", "PERF-26", "PERF-27", "PERF-28", "PERF-29", "PERF-30", "PERF-31", "PERF-33",
		"PERF-34", "PERF-35", "PERF-36", "PERF-37", "PERF-38", "PERF-39", "PERF-40", "PERF-41",
		"PERF-42", "PERF-43", "PERF-44", "PERF-45", "PERF-46", "PERF-47", "PERF-48", "PERF-49",
		"PERF-51", "PERF-52", "PERF-53", "PERF-54", "PERF-55", "PERF-56", "PERF-57", "PERF-58",
		"PERF-59", "PERF-60",
	}
	set := map[string]struct{}{}
	for _, id := range ids {
		set[id] = struct{}{}
	}
	for _, id := range want {
		if _, ok := set[id]; !ok {
			t.Errorf("missing registered rule %s", id)
		}
	}
}

func TestPERF9URLParseInLoop(t *testing.T) {
	src := `package sample
import (
	"net/http"
	"net/url"
)
func FetchAll(urls []string) ([]int, error) {
	statuses := make([]int, 0, len(urls))
	for _, raw := range urls {
		u, err := url.Parse(raw)
		if err != nil {
			return nil, err
		}
		resp, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, resp.StatusCode)
		resp.Body.Close()
	}
	return statuses, nil
}
`
	findings := runPerf(t, src, "PERF-009-vulnerable.go")
	if !hasRule(findings, "PERF-9") {
		t.Fatalf("expected PERF-9, got %#v", findings)
	}
}

func TestPERF9Safe(t *testing.T) {
	src := `package sample
import "net/url"
func ParseOnce(raw string) (*url.URL, error) {
	return url.Parse(raw)
}
`
	findings := runPerf(t, src, "PERF-009-safe.go")
	if hasRule(findings, "PERF-9") {
		t.Fatalf("safe should not emit PERF-9, got %#v", findings)
	}
}

func TestPERF10TemplateOnRequestPath(t *testing.T) {
	runFixtureRule(t, "PERF-010-vulnerable.txt", "PERF-10", true)
	runFixtureRule(t, "PERF-010-safe.txt", "PERF-10", false)
}

func TestPERF16BufferInLoop(t *testing.T) {
	runFixtureRule(t, "PERF-016-vulnerable.txt", "PERF-16", true)
	runFixtureRule(t, "PERF-016-safe.txt", "PERF-16", false)
}

func TestPERF24HasherInLoop(t *testing.T) {
	runFixtureRule(t, "PERF-024-vulnerable.txt", "PERF-24", true)
	runFixtureRule(t, "PERF-024-safe.txt", "PERF-24", false)
}

func TestPERF27LargeMakeInLoop(t *testing.T) {
	runFixtureRule(t, "PERF-027-vulnerable.txt", "PERF-27", true)
}

func TestPERF29UnboundedGo(t *testing.T) {
	runFixtureRule(t, "PERF-029-vulnerable.txt", "PERF-29", true)
}

func TestPERF36LoopCapture(t *testing.T) {
	runFixtureRule(t, "PERF-036-vulnerable.txt", "PERF-36", true)
	runFixtureRule(t, "PERF-036-safe.txt", "PERF-36", false)
}

func TestPERF38UnbufferedChan(t *testing.T) {
	runFixtureRule(t, "PERF-038-vulnerable.txt", "PERF-38", true)
}

func TestPERF40TimeNowRepeated(t *testing.T) {
	runFixtureRule(t, "PERF-040-vulnerable.txt", "PERF-40", true)
}

func TestPERF47SplitInLoop(t *testing.T) {
	runFixtureRule(t, "PERF-047-vulnerable.txt", "PERF-47", true)
}

func TestPERF56JSONInLoop(t *testing.T) {
	runFixtureRule(t, "PERF-056-vulnerable.txt", "PERF-56", true)
}

func TestPERF57MiddlewareHeavy(t *testing.T) {
	runFixtureRule(t, "PERF-057-vulnerable.txt", "PERF-57", true)
}

func TestPERF52RuntimeGC(t *testing.T) {
	src := `package sample
import "runtime"
func Cleanup() {
	runtime.GC()
}
`
	findings := runPerf(t, src, "PERF-052-vulnerable.go")
	if !hasRule(findings, "PERF-52") {
		t.Fatalf("expected PERF-52, got %#v", findings)
	}
}

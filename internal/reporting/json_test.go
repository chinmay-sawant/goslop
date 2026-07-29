package reporting

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/chinmay/codehound/internal/rules"
)

func TestJSONReporterRoundTrip(t *testing.T) {
	findings := []rules.Finding{
		{
			RuleID:    "CWE-89",
			RuleTitle: "SQL injection",
			File:      "main.go",
			Line:      10,
			Column:    3,
			Message:   "unsanitized query",
			Severity:  rules.SeverityHigh,
		},
	}
	var buf bytes.Buffer
	r := JSONReporter{Version: "0.1.0-dev"}
	if err := r.Write(findings, &buf); err != nil {
		t.Fatal(err)
	}

	var env struct {
		Findings []rules.Finding `json:"findings"`
		Version  string          `json:"version"`
	}
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v\nraw: %s", err, buf.String())
	}
	if env.Version != "0.1.0-dev" {
		t.Fatalf("version: %q", env.Version)
	}
	if len(env.Findings) != 1 {
		t.Fatalf("findings: %#v", env.Findings)
	}
	got := env.Findings[0]
	if got.RuleID != "CWE-89" || got.File != "main.go" || got.Line != 10 {
		t.Fatalf("finding: %+v", got)
	}
	if got.Fingerprint == "" {
		t.Fatal("expected fingerprint")
	}
}

func TestJSONReporterEmptyFindings(t *testing.T) {
	var buf bytes.Buffer
	if err := (JSONReporter{Version: "v"}).Write(nil, &buf); err != nil {
		t.Fatal(err)
	}
	var env map[string]any
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	raw, ok := env["findings"]
	if !ok {
		t.Fatal("missing findings")
	}
	arr, ok := raw.([]any)
	if !ok {
		t.Fatalf("findings type %T", raw)
	}
	if len(arr) != 0 {
		t.Fatalf("want empty, got %v", arr)
	}
}

func TestTextReporterLine(t *testing.T) {
	var buf bytes.Buffer
	err := (TextReporter{}).Write([]rules.Finding{{
		RuleID:  "PERF-1",
		File:    "a.go",
		Line:    2,
		Column:  4,
		Message: "hot path alloc",
	}}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	want := "PERF-1 a.go:2:4 hot path alloc\n"
	if buf.String() != want {
		t.Fatalf("got %q want %q", buf.String(), want)
	}
}

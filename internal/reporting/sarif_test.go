package reporting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestSARIFReporterMinimalShape(t *testing.T) {
	findings := []rules.Finding{
		{
			RuleID:    "CWE-78",
			RuleTitle: "OS Command Injection",
			File:      "cmd.go",
			Line:      12,
			Column:    5,
			Message:   "unsanitized command",
			Severity:  rules.SeverityHigh,
		},
		{
			RuleID:   "PERF-6",
			File:     "loop.go",
			Line:     3,
			Column:   1,
			Message:  "fmt in loop",
			Severity: rules.SeverityMedium,
		},
	}
	var buf bytes.Buffer
	if err := (SARIFReporter{Version: "0.1.0-dev"}).Write(findings, &buf); err != nil {
		t.Fatal(err)
	}

	var log map[string]any
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Fatalf("unmarshal: %v\nraw: %s", err, buf.String())
	}

	// SARIF 2.1.0 required top-level keys.
	if log["version"] != "2.1.0" {
		t.Fatalf("version: %#v", log["version"])
	}
	schema, _ := log["$schema"].(string)
	if schema == "" {
		t.Fatal("missing $schema")
	}
	runs, ok := log["runs"].([]any)
	if !ok || len(runs) != 1 {
		t.Fatalf("runs: %#v", log["runs"])
	}
	run, ok := runs[0].(map[string]any)
	if !ok {
		t.Fatalf("run type %T", runs[0])
	}
	tool, ok := run["tool"].(map[string]any)
	if !ok {
		t.Fatal("missing tool")
	}
	driver, ok := tool["driver"].(map[string]any)
	if !ok {
		t.Fatal("missing driver")
	}
	if driver["name"] != "goslop" {
		t.Fatalf("driver name: %#v", driver["name"])
	}
	if driver["version"] != "0.1.0-dev" {
		t.Fatalf("driver version: %#v", driver["version"])
	}
	rulesArr, ok := driver["rules"].([]any)
	if !ok || len(rulesArr) != 2 {
		t.Fatalf("rules: %#v", driver["rules"])
	}
	results, ok := run["results"].([]any)
	if !ok || len(results) != 2 {
		t.Fatalf("results: %#v", run["results"])
	}

	// First result: high → error level, location present.
	r0, ok := results[0].(map[string]any)
	if !ok {
		t.Fatalf("result0 type %T", results[0])
	}
	if r0["ruleId"] != "CWE-78" {
		t.Fatalf("ruleId: %#v", r0["ruleId"])
	}
	if r0["level"] != "error" {
		t.Fatalf("level for high: %#v", r0["level"])
	}
	locs, ok := r0["locations"].([]any)
	if !ok || len(locs) < 1 {
		t.Fatalf("locations: %#v", r0["locations"])
	}
	fp, ok := r0["partialFingerprints"].(map[string]any)
	if !ok || fp["primaryLocationLineHash"] == "" {
		t.Fatalf("partialFingerprints: %#v", r0["partialFingerprints"])
	}
}

func TestSARIFReporterEmptyFindings(t *testing.T) {
	var buf bytes.Buffer
	if err := (SARIFReporter{Version: "v"}).Write(nil, &buf); err != nil {
		t.Fatal(err)
	}
	var log map[string]any
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Fatal(err)
	}
	runs := log["runs"].([]any)
	run := runs[0].(map[string]any)
	results, ok := run["results"].([]any)
	if !ok {
		// empty array may decode as nil depending on omitempty; both ok if no panic
		if run["results"] != nil {
			t.Fatalf("results type %T", run["results"])
		}
		return
	}
	if len(results) != 0 {
		t.Fatalf("want empty results, got %v", results)
	}
}

func TestReportersDoNotMutateFindings(t *testing.T) {
	findings := []rules.Finding{{
		RuleID: "PERF-1", File: "main.go", Line: 1, Column: 1,
		Message: "test", Severity: rules.SeverityLow,
	}}
	for _, reporter := range []Reporter{
		JSONReporter{Version: "v"},
		SARIFReporter{Version: "v"},
	} {
		var buf bytes.Buffer
		if err := reporter.Write(findings, &buf); err != nil {
			t.Fatal(err)
		}
		if findings[0].Fingerprint != "" {
			t.Fatalf("%T mutated caller finding", reporter)
		}
	}
}

func TestMachineReportersPreserveFingerprintAndLocationParity(t *testing.T) {
	findings := []rules.Finding{{
		RuleID: "CWE-78", RuleTitle: "Command Injection", File: "cmd.go",
		Line: 17, Column: 4, Message: "unsanitized command", Severity: rules.SeverityHigh,
	}}
	var jsonBuf, sarifBuf bytes.Buffer
	if err := (JSONReporter{Version: "v"}).Write(findings, &jsonBuf); err != nil {
		t.Fatal(err)
	}
	if err := (SARIFReporter{Version: "v"}).Write(findings, &sarifBuf); err != nil {
		t.Fatal(err)
	}

	var jsonOutput struct {
		Findings []rules.Finding `json:"findings"`
	}
	if err := json.Unmarshal(jsonBuf.Bytes(), &jsonOutput); err != nil {
		t.Fatal(err)
	}
	var sarifOutput sarifLog
	if err := json.Unmarshal(sarifBuf.Bytes(), &sarifOutput); err != nil {
		t.Fatal(err)
	}
	if len(jsonOutput.Findings) != 1 || len(sarifOutput.Runs) != 1 || len(sarifOutput.Runs[0].Results) != 1 {
		t.Fatalf("unexpected report shapes: json=%+v sarif=%+v", jsonOutput, sarifOutput)
	}
	jsonFinding := jsonOutput.Findings[0]
	sarifFinding := sarifOutput.Runs[0].Results[0]
	location := sarifFinding.Locations[0].PhysicalLocation
	if jsonFinding.Fingerprint == "" || jsonFinding.Fingerprint != sarifFinding.PartialFingerprints["primaryLocationLineHash"] {
		t.Fatalf("fingerprints differ: json=%q sarif=%q", jsonFinding.Fingerprint, sarifFinding.PartialFingerprints)
	}
	if jsonFinding.File != location.ArtifactLocation.URI || jsonFinding.Line != location.Region.StartLine || jsonFinding.Column != location.Region.StartColumn {
		t.Fatalf("locations differ: json=%+v sarif=%+v", jsonFinding, location)
	}
	if findings[0].Fingerprint != "" {
		t.Fatalf("reporters mutated caller finding: %q", findings[0].Fingerprint)
	}
}

func TestMachineReportersSupportConcurrentWrites(t *testing.T) {
	findings := []rules.Finding{{
		RuleID: "PERF-1", File: "main.go", Line: 1, Column: 1,
		Message: "allocation", Severity: rules.SeverityLow,
	}}
	reporters := []Reporter{JSONReporter{Version: "v"}, SARIFReporter{Version: "v"}}
	errs := make(chan error, len(reporters)*16)
	var wg sync.WaitGroup
	for _, reporter := range reporters {
		for range 16 {
			wg.Add(1)
			go func(reporter Reporter) {
				defer wg.Done()
				var out bytes.Buffer
				if err := reporter.Write(findings, &out); err != nil {
					errs <- err
					return
				}
				if !json.Valid(out.Bytes()) {
					errs <- fmt.Errorf("%T wrote invalid JSON: %q", reporter, out.String())
				}
			}(reporter)
		}
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	if findings[0].Fingerprint != "" {
		t.Fatalf("concurrent reporters mutated caller finding: %q", findings[0].Fingerprint)
	}
}

func TestSARIFSeverityMapping(t *testing.T) {
	cases := []struct {
		sev   rules.Severity
		level string
	}{
		{rules.SeverityInfo, "note"},
		{rules.SeverityLow, "warning"},
		{rules.SeverityMedium, "warning"},
		{rules.SeverityHigh, "error"},
		{rules.SeverityCritical, "error"},
	}
	for _, tc := range cases {
		if got := severityToSARIF(tc.sev); got != tc.level {
			t.Errorf("severity %v → %q want %q", tc.sev, got, tc.level)
		}
	}
}

func TestJSONReporterRequiredFields(t *testing.T) {
	// Schema-ish smoke: every finding in the envelope has the wire fields
	// consumers expect (rule_id, file, line, message, severity, fingerprint).
	f := rules.NewFinding(rules.FindingInputs{
		RuleID:    "PERF-1",
		RuleTitle: "alloc",
		File:      "a.go",
		Location:  rules.LineCol{Line: 1, Column: 2},
		Message:   "hot path",
		Severity:  rules.SeverityLow,
	})
	var buf bytes.Buffer
	if err := (JSONReporter{Version: "0.1.0-dev"}).Write([]rules.Finding{f}, &buf); err != nil {
		t.Fatal(err)
	}
	var env struct {
		Findings []map[string]any `json:"findings"`
		Version  string           `json:"version"`
	}
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env.Version != "0.1.0-dev" || len(env.Findings) != 1 {
		t.Fatalf("envelope: %+v", env)
	}
	got := env.Findings[0]
	for _, key := range []string{"rule_id", "file", "line", "column", "message", "severity", "fingerprint"} {
		if _, ok := got[key]; !ok {
			t.Errorf("missing key %q in %v", key, got)
		}
	}
	if got["rule_id"] != "PERF-1" {
		t.Fatalf("rule_id: %#v", got["rule_id"])
	}
}

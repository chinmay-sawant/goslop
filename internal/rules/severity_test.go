package rules

import (
	"encoding/json"
	"testing"
)

func TestParseSeverity(t *testing.T) {
	cases := []struct {
		in   string
		want Severity
	}{
		{"info", SeverityInfo},
		{"INFO", SeverityInfo},
		{"low", SeverityLow},
		{"medium", SeverityMedium},
		{"high", SeverityHigh},
		{"critical", SeverityCritical},
		{" Critical ", SeverityCritical},
	}
	for _, tc := range cases {
		got, err := ParseSeverity(tc.in)
		if err != nil {
			t.Fatalf("ParseSeverity(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseSeverity(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
	if _, err := ParseSeverity("nope"); err == nil {
		t.Fatal("expected error for unknown severity")
	}
}

func TestSeverityStringAndJSON(t *testing.T) {
	for _, s := range []Severity{
		SeverityInfo, SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical,
	} {
		b, err := json.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		var back Severity
		if err := json.Unmarshal(b, &back); err != nil {
			t.Fatal(err)
		}
		if back != s {
			t.Fatalf("json round-trip: got %v want %v (raw %s)", back, s, b)
		}
		var raw string
		if err := json.Unmarshal(b, &raw); err != nil {
			t.Fatal(err)
		}
		if raw != s.String() {
			t.Fatalf("json string = %q, want %q", raw, s.String())
		}
	}
}

func TestSeverityIsFailure(t *testing.T) {
	if SeverityInfo.IsFailure() || SeverityLow.IsFailure() {
		t.Fatal("info/low should not be failure severities")
	}
	if !SeverityMedium.IsFailure() || !SeverityHigh.IsFailure() || !SeverityCritical.IsFailure() {
		t.Fatal("medium+ should be failure severities")
	}
}

package rules

import "testing"

func TestFingerprintV2Stable(t *testing.T) {
	a := FingerprintV2("CWE-78", "main.go", "msg")
	b := FingerprintV2("CWE-78", "main.go", "msg")
	if a != b {
		t.Fatalf("unstable fingerprint: %q vs %q", a, b)
	}
	if len(a) < len("goslop:2:") || a[:len("goslop:2:")] != "goslop:2:" {
		t.Fatalf("prefix: %q", a)
	}
	c := FingerprintV2("CWE-78", "main.go", "other")
	if a == c {
		t.Fatal("different messages should differ")
	}
	f := NewFinding(FindingInputs{
		RuleID: "CWE-78", RuleTitle: "t", File: "main.go",
		Location: LineCol{Line: 10, Column: 2},
		Message:  "msg", Severity: SeverityHigh,
	})
	if f.Fingerprint != a {
		t.Fatalf("NewFinding fingerprint %q != %q", f.Fingerprint, a)
	}
	f.Fingerprint = ""
	f.EnsureFingerprint()
	if f.Fingerprint != a {
		t.Fatalf("EnsureFingerprint %q != %q", f.Fingerprint, a)
	}
}

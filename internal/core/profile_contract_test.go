package core

import "testing"

// Profile contract tests (Phase 12.1) — keep pack defaults stable for CI docs.

func TestProfileFailPolicies(t *testing.T) {
	cases := []struct {
		p    ScanProfile
		want FailPolicy
	}{
		{ProfileRecommended, FailHigh},
		{ProfileSecurity, FailHigh},
		{ProfilePerf, FailHigh},
		{ProfileStyle, FailNone},
		{ProfileAll, FailMedium},
	}
	for _, tc := range cases {
		if got := tc.p.DefaultFailPolicy(); got != tc.want {
			t.Errorf("%s fail=%v want %v", tc.p, got, tc.want)
		}
	}
}

func TestProfileOnlyPatternsContract(t *testing.T) {
	if ProfileAll.OnlyPatterns() != nil {
		t.Fatal("all pack should not set only patterns")
	}
	if got := ProfilePerf.OnlyPatterns(); len(got) != 1 || got[0] != "PERF-*" {
		t.Fatalf("perf: %#v", got)
	}
	if got := ProfileSecurity.OnlyPatterns(); len(got) != 1 || got[0] != "CWE-*" {
		t.Fatalf("security: %#v", got)
	}
	if got := ProfileStyle.OnlyPatterns(); len(got) != 1 || got[0] != "BP-*" {
		t.Fatalf("style: %#v", got)
	}
}

func TestBuildScanContextCLIUnion(t *testing.T) {
	ctx := BuildScanContext(ProfilePerf, []string{"CWE-78"}, []string{"PERF-1"})
	if !ctx.Allows("PERF-6") {
		t.Fatal("perf pack should allow PERF-6")
	}
	if !ctx.Allows("CWE-78") {
		t.Fatal("CLI only should union CWE-78 onto perf pack")
	}
	if ctx.Allows("PERF-1") {
		t.Fatal("skip should deny PERF-1")
	}
	if ctx.Allows("BP-1") {
		t.Fatal("perf pack should not enable BP by default")
	}
}

func TestParseScanProfileAliases(t *testing.T) {
	aliases := map[string]ScanProfile{
		"recommended": ProfileRecommended,
		"ci":          ProfileRecommended,
		"default":     ProfileRecommended,
		"perf":        ProfilePerf,
		"performance": ProfilePerf,
		"security":    ProfileSecurity,
		"sec":         ProfileSecurity,
		"style":       ProfileStyle,
		"bp":          ProfileStyle,
		"all":         ProfileAll,
		"full":        ProfileAll,
	}
	for name, want := range aliases {
		got, err := ParseScanProfile(name)
		if err != nil || got != want {
			t.Errorf("ParseScanProfile(%q)=%v,%v want %v", name, got, err, want)
		}
	}
	if _, err := ParseScanProfile("nope"); err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

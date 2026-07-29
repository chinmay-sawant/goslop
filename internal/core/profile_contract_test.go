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
	// Perf = S-tier + A-tier explicit IDs (not PERF-*).
	if got := ProfilePerf.OnlyPatterns(); len(got) < 2 {
		t.Fatalf("perf pack too small: %#v", got)
	}
	// Security = exact SECURITY_PACK_RULES list.
	if got := ProfileSecurity.OnlyPatterns(); len(got) < 5 {
		t.Fatalf("security pack too small: %#v", got)
	}
	if got := ProfileStyle.OnlyPatterns(); len(got) != 1 || got[0] != "BP-*" {
		t.Fatalf("style: %#v", got)
	}
	// Recommended: S-tier PERF + taint-core CWE.
	rec := ProfileRecommended.OnlyPatterns()
	hasPERF1, hasCWE22 := false, false
	for _, id := range rec {
		if id == "PERF-1" {
			hasPERF1 = true
		}
		if id == "CWE-22" {
			hasCWE22 = true
		}
	}
	if !hasPERF1 || !hasCWE22 {
		t.Fatalf("recommended missing core ids: %#v", rec)
	}
}

func TestBuildScanContextCLIUnion(t *testing.T) {
	// Perf S-tier includes PERF-1; skip should deny it.
	ctx := BuildScanContext(ProfilePerf, []string{"CWE-78"}, []string{"PERF-1"})
	if !ctx.Allows("PERF-7") {
		t.Fatal("perf pack should allow PERF-7 (S-tier)")
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
	// Security enables taint.
	sec := BuildScanContext(ProfileSecurity, nil, nil)
	if !sec.TaintEnabled {
		t.Fatal("security profile should enable taint")
	}
	// Style skips noisy BPs by default.
	style := BuildScanContext(ProfileStyle, nil, nil)
	if style.Allows("BP-21") {
		t.Fatal("style should skip BP-21 by default")
	}
	if !style.Allows("BP-1") {
		t.Fatal("style should allow BP-1")
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

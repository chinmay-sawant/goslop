package core

import (
	"fmt"
	"strings"
)

// ScanProfile is a named product pack. CLI default is ProfileRecommended.
type ScanProfile int

const (
	// ProfileRecommended is the curated CI pack. BP off. Fail high.
	ProfileRecommended ScanProfile = iota
	// ProfilePerf is framework + hot-path PERF. BP off.
	ProfilePerf
	// ProfileSecurity is taint CWE core + structural CWEs. Taint on. BP off.
	ProfileSecurity
	// ProfileStyle is bad-practice / style pack only (advisory).
	ProfileStyle
	// ProfileAll is the full catalog.
	ProfileAll
)

// String returns the stable lowercase profile name.
func (p ScanProfile) String() string {
	switch p {
	case ProfileRecommended:
		return "recommended"
	case ProfilePerf:
		return "perf"
	case ProfileSecurity:
		return "security"
	case ProfileStyle:
		return "style"
	case ProfileAll:
		return "all"
	default:
		return fmt.Sprintf("ScanProfile(%d)", int(p))
	}
}

// ParseScanProfile parses a profile name (case-insensitive).
func ParseScanProfile(s string) (ScanProfile, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "recommended", "default", "ci":
		return ProfileRecommended, nil
	case "perf", "performance":
		return ProfilePerf, nil
	case "security", "sec":
		return ProfileSecurity, nil
	case "style", "bp", "bad-practices", "bad_practices":
		return ProfileStyle, nil
	case "all", "full":
		return ProfileAll, nil
	default:
		return 0, fmt.Errorf("unknown scan profile %q", s)
	}
}

// ParseProfile is like ParseScanProfile but returns ok bool.
func ParseProfile(s string) (ScanProfile, bool) {
	p, err := ParseScanProfile(s)
	return p, err == nil
}

// DefaultFailPolicy is the fail policy for this profile when CLI did not set fail flags.
func (p ScanProfile) DefaultFailPolicy() FailPolicy {
	switch p {
	case ProfileRecommended, ProfileSecurity, ProfilePerf:
		return FailHigh
	case ProfileStyle:
		return FailNone
	case ProfileAll:
		return FailMedium
	default:
		return FailMedium
	}
}

// EnablesTaint reports whether taint should be enabled for this profile.
func (p ScanProfile) EnablesTaint() bool {
	return p == ProfileSecurity
}

// EnablesBadPractices reports whether BP rules run for this profile.
func (p ScanProfile) EnablesBadPractices() bool {
	return p == ProfileStyle || p == ProfileAll
}

// OnlyPatterns returns a coarse allow-list for the profile.
// nil means no only filter (full catalog subject to skip/BP).
func (p ScanProfile) OnlyPatterns() []string {
	switch p {
	case ProfileRecommended:
		// Curated CI pack placeholders until full rule tables land.
		return []string{
			"PERF-101",
			"CWE-22", "CWE-78", "CWE-79", "CWE-89",
		}
	case ProfilePerf:
		return []string{"PERF-*"}
	case ProfileSecurity:
		return []string{"CWE-*"}
	case ProfileStyle:
		return []string{"BP-*"}
	case ProfileAll:
		return nil
	default:
		return nil
	}
}

// BuildScanContext creates a ScanContext from a profile and CLI only/skip lists.
// CLI only is unioned with pack patterns; CLI skip is appended.
func BuildScanContext(profile ScanProfile, only, skip []string) *ScanContext {
	ctx := DefaultScanContext()
	ctx.Profile = profile
	ctx.FailPolicy = profile.DefaultFailPolicy()
	ctx.TaintEnabled = profile.EnablesTaint()
	ctx.BadPracticesEnabled = profile.EnablesBadPractices()

	pack := profile.OnlyPatterns()
	if pack != nil || len(only) > 0 {
		merged := make([]string, 0, len(pack)+len(only))
		merged = append(merged, pack...)
		merged = append(merged, only...)
		ctx.Only = merged
	}
	if len(skip) > 0 {
		ctx.Skip = append([]string(nil), skip...)
	}
	return ctx
}

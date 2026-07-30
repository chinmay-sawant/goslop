package integration

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Fixture quarantine: known heuristic gaps vs Rust fixture expectations.
// Entries are paths relative to tests/fixtures/ (slash form), one per line.
// Blank lines and # comments are ignored.
//
// Goal: keep the full Rust-style matrix in CI while not blocking on detectors
// that are still incomplete. Unexpected failures still fail the suite;
// empty quarantine means full hard gate.

var (
	quarantineOnce sync.Once
	quarantineSet  map[string]struct{}
	quarantineErr  error
)

// QuarantineFileName is the default allowlist file under tests/integration/.
const QuarantineFileName = "fixture_quarantine.txt"

// LoadQuarantine loads the allowlist once from tests/integration/fixture_quarantine.txt.
func LoadQuarantine() (map[string]struct{}, error) {
	quarantineOnce.Do(func() {
		quarantineSet = map[string]struct{}{}
		root, err := RepoRoot()
		if err != nil {
			quarantineErr = err
			return
		}
		path := filepath.Join(root, "tests", "integration", QuarantineFileName)
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				return // empty allowlist
			}
			quarantineErr = err
			return
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			// allow optional trailing comment
			if i := strings.Index(line, " #"); i >= 0 {
				line = strings.TrimSpace(line[:i])
			}
			quarantineSet[filepath.ToSlash(line)] = struct{}{}
		}
		if err := sc.Err(); err != nil {
			quarantineErr = fmt.Errorf("scan quarantine file: %w", err)
		}
	})
	return quarantineSet, quarantineErr
}

// IsQuarantined reports whether a fixture rel path is allowlisted.
func IsQuarantined(relPath string) bool {
	set, err := LoadQuarantine()
	if err != nil || set == nil {
		return false
	}
	_, ok := set[filepath.ToSlash(relPath)]
	return ok
}

// PartitionFailures splits failures into unexpected vs quarantined.
// Each failure string should start with the fixture rel path (as produced by matrix tests).
func PartitionFailures(failures []string) ([]string, []string) {
	unexpected := make([]string, 0, len(failures))
	known := make([]string, 0, len(failures)/4)
	for _, f := range failures {
		rel := failureRelPath(f)
		if rel != "" && IsQuarantined(rel) {
			known = append(known, f)
			continue
		}
		unexpected = append(unexpected, f)
	}
	return unexpected, known
}

func failureRelPath(failure string) string {
	// formats:
	//   go/bad_practices/BP-1-safe.txt: expected ...
	//   go/frameworks/CWE-15-vulnerable.txt: expected ...
	//   BP-18 vuln: materialize error...
	s := failure
	if i := strings.Index(s, ": "); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "go/") {
		return filepath.ToSlash(s)
	}
	return ""
}

// FormatMatrixReport builds a multi-line failure report.
func FormatMatrixReport(unexpected, known []string, limit int) string {
	var b strings.Builder
	if len(unexpected) > 0 {
		fmt.Fprintf(&b, "unexpected failures (%d):\n  %s\n", len(unexpected), joinN(unexpected, limit))
	}
	if len(known) > 0 {
		fmt.Fprintf(&b, "quarantined failures (%d, known detector gaps):\n  %s\n", len(known), joinN(known, limit))
	}
	return b.String()
}

func joinN(ss []string, n int) string {
	total := len(ss)
	if total > n {
		ss = ss[:n]
	}
	return strings.Join(ss, "\n  ") + func() string {
		if total > n {
			return fmt.Sprintf("\n  ... and %d more", total-n)
		}
		return ""
	}()
}

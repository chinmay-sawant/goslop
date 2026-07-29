// Package engine implements the analysis pipeline: walk, registry, and parallel scan.
package engine

import (
	"fmt"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

// ScanErrorKind is a coarse per-file error category.
type ScanErrorKind int

const (
	// ScanErrorIO is a failure reading a file or directory.
	ScanErrorIO ScanErrorKind = iota
	// ScanErrorEncoding is invalid UTF-8 source.
	ScanErrorEncoding
	// ScanErrorParse is a language parse failure.
	ScanErrorParse
	// ScanErrorEngine is a detector/engine internal error.
	ScanErrorEngine
)

// String returns a short category name.
func (k ScanErrorKind) String() string {
	switch k {
	case ScanErrorIO:
		return "io"
	case ScanErrorEncoding:
		return "encoding"
	case ScanErrorParse:
		return "parse"
	case ScanErrorEngine:
		return "engine"
	default:
		return "unknown"
	}
}

// ScanError is a non-fatal error for one file; the scan continues.
type ScanError struct {
	Path    string
	Kind    ScanErrorKind
	Message string
}

// Error implements the error interface.
func (e ScanError) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ScanStats holds optional operational counters for a scan.
type ScanStats struct {
	FilesScanned       int
	FilesSkipped       int
	FilesErrored       int
	BytesScanned       int64
	LinesScanned       int
	FindingsTotal      int
	FindingsSuppressed int
	FindingsBaselined  int
	CacheHits          int
	CacheMisses        int
	DetectorsLoaded    int
}

// AnalysisResult holds findings and per-file errors from a scan run.
type AnalysisResult struct {
	Findings []rules.Finding
	// Errors are non-fatal per-file errors (scan continues).
	Errors []ScanError
	// Stats is populated with basic counts for every scan.
	Stats *ScanStats
	// SourceCache maps display path → source when ScanContext.RetainSources is set.
	SourceCache map[string]string
	// SuppressedCount is findings removed/tagged by codehound-ignore directives.
	SuppressedCount int
	// BaselinedCount is findings filtered by the baseline store.
	BaselinedCount int
}

// ShouldFail reports whether any finding matches the fail policy.
func (r *AnalysisResult) ShouldFail(policy core.FailPolicy) bool {
	if r == nil {
		return false
	}
	for i := range r.Findings {
		if policy.ShouldFail(r.Findings[i].Severity) {
			return true
		}
	}
	return false
}

package rules

import "github.com/chinmay/codehound/internal/cwe"

// LineCol is a 1-indexed line and column in a source file.
type LineCol struct {
	Line   int
	Column int
}

// FindingInputs holds the core fields required to construct a Finding.
type FindingInputs struct {
	RuleID    string
	RuleTitle string
	File      string
	Location  LineCol
	Message   string
	Severity  Severity
	// CWE accepts string ids ("CWE-89") for convenience.
	CWE []string
	// CWERefs is used when full refs are already available (takes precedence when non-empty).
	CWERefs []cwe.Ref
	Fix     string
}

// Finding is a single static-analysis finding.
type Finding struct {
	RuleID      string         `json:"rule_id"`
	RuleTitle   string         `json:"rule_title"`
	File        string         `json:"file"`
	Line        int            `json:"line"`
	Column      int            `json:"column"`
	Message     string         `json:"message"`
	Severity    Severity       `json:"severity"`
	CWE         []cwe.Ref      `json:"cwe"`
	Fingerprint string         `json:"fingerprint,omitempty"`
	Fix         string         `json:"fix,omitempty"`
	Snippet     string         `json:"snippet,omitempty"`
	Evidence    map[string]any `json:"evidence,omitempty"`
	Confidence  float32        `json:"confidence,omitempty"`
	Suppressed  bool           `json:"suppressed,omitempty"`
	Remediation string         `json:"remediation,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
}

// NewFinding constructs a Finding from inputs and computes the v2 fingerprint.
func NewFinding(in FindingInputs) Finding {
	cwes := in.CWERefs
	if len(cwes) == 0 && len(in.CWE) > 0 {
		cwes = cwe.RefsFromIDs(in.CWE)
	}
	if cwes == nil {
		cwes = []cwe.Ref{}
	}
	return Finding{
		RuleID:      in.RuleID,
		RuleTitle:   in.RuleTitle,
		File:        in.File,
		Line:        in.Location.Line,
		Column:      in.Location.Column,
		Message:     in.Message,
		Severity:    in.Severity,
		CWE:         cwes,
		Fix:         in.Fix,
		Fingerprint: Fingerprint(in.RuleID, in.File, in.Message),
	}
}

// NewFindingFromMeta constructs a Finding from rule metadata and a location.
func NewFindingFromMeta(meta *RuleMetadata, file string, line, col int, message string) Finding {
	if meta == nil {
		return Finding{CWE: []cwe.Ref{}, File: file, Line: line, Column: col, Message: message}
	}
	cwes := append([]cwe.Ref(nil), meta.CWE...)
	if cwes == nil {
		cwes = []cwe.Ref{}
	}
	return Finding{
		RuleID:      meta.ID,
		RuleTitle:   meta.Title,
		File:        file,
		Line:        line,
		Column:      col,
		Message:     message,
		Severity:    meta.Severity,
		CWE:         cwes,
		Fix:         meta.Fix,
		Fingerprint: Fingerprint(meta.ID, file, message),
	}
}

// EnsureFingerprint fills Fingerprint when empty (stable for reporters).
func (f *Finding) EnsureFingerprint() {
	if f == nil {
		return
	}
	if f.Fingerprint == "" {
		f.Fingerprint = Fingerprint(f.RuleID, f.File, f.Message)
	}
}

// FingerprintString returns the content-stable fingerprint.
func (f Finding) FingerprintString() string {
	if f.Fingerprint != "" {
		return f.Fingerprint
	}
	return Fingerprint(f.RuleID, f.File, f.Message)
}

// WithFix attaches a short suggested fix.
func (f Finding) WithFix(fix string) Finding {
	f.Fix = fix
	return f
}

// WithSnippet attaches a source snippet.
func (f Finding) WithSnippet(snippet string) Finding {
	f.Snippet = snippet
	return f
}

// WithEvidence attaches structured evidence.
func (f Finding) WithEvidence(evidence map[string]any) Finding {
	f.Evidence = evidence
	return f
}

// Category returns the coarse pack category for the finding's rule id.
func (f Finding) Category() string {
	return PackFromRuleID(f.RuleID).CategoryStr()
}

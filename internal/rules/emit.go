package rules

import "github.com/chinmay-sawant/goslop/internal/cwe"

// PushFinding appends a finding built from metadata.
func PushFinding(meta *RuleMetadata, file string, line, col int, message string, out *[]Finding) {
	PushFindingWithConfidence(meta, file, line, col, message, 0, out)
}

// PushFindingWithConfidence appends a finding with optional confidence (0 = omit).
func PushFindingWithConfidence(meta *RuleMetadata, file string, line, col int, message string, confidence float32, out *[]Finding) {
	PushFindingWithEvidence(meta, file, line, col, message, confidence, nil, out)
}

// PushFindingWithEvidence appends a finding with optional confidence and evidence.
func PushFindingWithEvidence(meta *RuleMetadata, file string, line, col int, message string, confidence float32, evidence map[string]any, out *[]Finding) {
	if meta == nil || out == nil {
		return
	}
	if line <= 0 {
		line = 1
	}
	if col <= 0 {
		col = 1
	}
	var cwes []cwe.CweRef
	if len(meta.CWE) > 0 {
		cwes = append([]cwe.CweRef(nil), meta.CWE...)
	} else {
		cwes = []cwe.CweRef{}
	}
	f := Finding{
		RuleID:     meta.ID,
		RuleTitle:  meta.Title,
		File:       file,
		Line:       line,
		Column:     col,
		Message:    message,
		Severity:   meta.Severity,
		CWE:        cwes,
		Fix:        meta.Fix,
		Confidence: confidence,
		Evidence:   evidence,
	}
	f.EnsureFingerprint()
	*out = append(*out, f)
}

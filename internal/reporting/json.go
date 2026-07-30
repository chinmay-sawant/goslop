package reporting

import (
	"encoding/json"
	"io"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

// JSONReporter writes a simple envelope:
//
//	{"findings":[...],"version":"..."}
type JSONReporter struct {
	Version string
}

type jsonEnvelope struct {
	Findings []rules.Finding `json:"findings"`
	Version  string          `json:"version"`
}

// Write implements Reporter.
func (r JSONReporter) Write(findings []rules.Finding, w io.Writer) error {
	version := r.Version
	if version == "" {
		version = DefaultVersion
	}
	// Ensure fingerprints for stable machine output without mutating callers.
	out := normalizedFindings(findings)
	env := jsonEnvelope{
		Findings: out,
		Version:  version,
	}
	if env.Findings == nil {
		env.Findings = []rules.Finding{}
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(env)
}

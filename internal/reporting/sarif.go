package reporting

import (
	"encoding/json"
	"io"

	"github.com/chinmay/goslop/internal/rules"
)

// SARIFReporter emits a minimal valid SARIF 2.1.0 log.
type SARIFReporter struct {
	Version string // tool version
}

// Write implements Reporter.
func (r SARIFReporter) Write(findings []rules.Finding, w io.Writer) error {
	version := r.Version
	if version == "" {
		version = DefaultVersion
	}

	// Collect unique rules.
	type ruleMeta struct {
		id, name string
	}
	seen := make(map[string]ruleMeta)
	order := make([]string, 0)
	for i := range findings {
		f := &findings[i]
		if _, ok := seen[f.RuleID]; ok {
			continue
		}
		name := f.RuleTitle
		if name == "" {
			name = f.RuleID
		}
		seen[f.RuleID] = ruleMeta{id: f.RuleID, name: name}
		order = append(order, f.RuleID)
	}

	rulesArr := make([]sarifRule, 0, len(order))
	for _, id := range order {
		m := seen[id]
		rulesArr = append(rulesArr, sarifRule{
			ID:               m.id,
			Name:             m.name,
			ShortDescription: sarifText{Text: m.name},
		})
	}

	results := make([]sarifResult, 0, len(findings))
	for i := range findings {
		f := &findings[i]
		f.EnsureFingerprint()
		line := f.Line
		col := f.Column
		if line <= 0 {
			line = 1
		}
		if col <= 0 {
			col = 1
		}
		level := severityToSARIF(f.Severity)
		results = append(results, sarifResult{
			RuleID: f.RuleID,
			Level:  level,
			Message: sarifText{
				Text: f.Message,
			},
			Locations: []sarifLocation{{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifactLocation{URI: f.File},
					Region: sarifRegion{
						StartLine:   line,
						StartColumn: col,
					},
				},
			}},
			PartialFingerprints: map[string]string{
				"primaryLocationLineHash": f.Fingerprint,
			},
		})
	}

	log := sarifLog{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{
				Driver: sarifDriver{
					Name:           "goslop",
					InformationURI: "https://github.com/chinmay/goslop",
					Version:        version,
					Rules:          rulesArr,
				},
			},
			Results: results,
		}},
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(log)
}

func severityToSARIF(s rules.Severity) string {
	switch s {
	case rules.SeverityInfo:
		return "note"
	case rules.SeverityLow, rules.SeverityMedium:
		return "warning"
	case rules.SeverityHigh, rules.SeverityCritical:
		return "error"
	default:
		return "warning"
	}
}

// Minimal SARIF 2.1.0 shapes (subset).

type sarifLog struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri,omitempty"`
	Version        string      `json:"version,omitempty"`
	Rules          []sarifRule `json:"rules,omitempty"`
}

type sarifRule struct {
	ID               string    `json:"id"`
	Name             string    `json:"name,omitempty"`
	ShortDescription sarifText `json:"shortDescription,omitempty"`
}

type sarifText struct {
	Text string `json:"text"`
}

type sarifResult struct {
	RuleID              string            `json:"ruleId"`
	Level               string            `json:"level,omitempty"`
	Message             sarifText         `json:"message"`
	Locations           []sarifLocation   `json:"locations"`
	PartialFingerprints map[string]string `json:"partialFingerprints,omitempty"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn,omitempty"`
}

package core

import "github.com/chinmay/codehound/internal/rules"

// Detector walks parsed units and appends findings.
type Detector interface {
	Language() LanguageID
	RuleIDs() []string
	BeginScan(ctx *ScanContext)
	EndScan()
	Run(ctx *ScanContext, unit *ParsedUnit, out *[]rules.Finding)
	AccumulateState(ctx *ScanContext, unit *ParsedUnit)
	RequiresCacheState(ctx *ScanContext) bool
	Finalize(ctx *ScanContext, out *[]rules.Finding)
	ResetState()
	// MetadataFor returns catalogue metadata for a rule id (optional; may return nil).
	MetadataFor(ruleID string) *rules.RuleMetadata
}

// BaseDetector provides no-op lifecycle defaults for stateless detectors.
type BaseDetector struct{}

func (BaseDetector) BeginScan(*ScanContext)                    {}
func (BaseDetector) EndScan()                                  {}
func (BaseDetector) AccumulateState(*ScanContext, *ParsedUnit) {}
func (BaseDetector) RequiresCacheState(*ScanContext) bool      { return false }
func (BaseDetector) Finalize(*ScanContext, *[]rules.Finding)   {}
func (BaseDetector) ResetState()                               {}
func (BaseDetector) MetadataFor(string) *rules.RuleMetadata    { return nil }

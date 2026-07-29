// Package cwe provides CWE catalogue references attached to findings.
package cwe

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// CweRef is a CWE identifier with display metadata.
// ID is the display form "CWE-N".
type CweRef struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Ref is an alias for CweRef (used by detectors).
type Ref = CweRef

// New builds a CweRef from a numeric id.
func New(id uint, name, url string) CweRef {
	if url == "" && id > 0 {
		url = fmt.Sprintf("https://cwe.mitre.org/data/definitions/%d.html", id)
	}
	return CweRef{
		ID:   fmt.Sprintf("CWE-%d", id),
		Name: name,
		URL:  url,
	}
}

// NewFromID parses "CWE-N", "cwe-N", or bare "N".
func NewFromID(s string) CweRef {
	s = strings.TrimSpace(s)
	if s == "" {
		return CweRef{}
	}
	n := parseNumeric(s)
	if n == 0 {
		return CweRef{ID: s}
	}
	return New(n, "", "")
}

// NumericID returns the numeric portion of the CWE id.
func (r CweRef) NumericID() uint {
	return parseNumeric(r.ID)
}

// String returns the display id.
func (r CweRef) String() string {
	if r.ID != "" {
		return r.ID
	}
	return ""
}

// DisplayID returns the display id.
func (r CweRef) DisplayID() string { return r.String() }

// RefsFromIDs parses a list of CWE id strings.
func RefsFromIDs(ids []string) []CweRef {
	out := make([]CweRef, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		out = append(out, NewFromID(id))
	}
	return out
}

// FormatList formats refs as "CWE-N (Name), ..." sorted by id.
func FormatList(refs []CweRef) string {
	if len(refs) == 0 {
		return ""
	}
	cp := append([]CweRef(nil), refs...)
	sort.Slice(cp, func(i, j int) bool {
		return cp[i].NumericID() < cp[j].NumericID()
	})
	parts := make([]string, 0, len(cp))
	for _, r := range cp {
		if r.Name != "" {
			parts = append(parts, fmt.Sprintf("%s (%s)", r.ID, r.Name))
		} else {
			parts = append(parts, r.ID)
		}
	}
	return strings.Join(parts, ", ")
}

func parseNumeric(s string) uint {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "CWE-")
	s = strings.TrimPrefix(s, "cwe-")
	s = strings.TrimPrefix(s, "Cwe-")
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(n)
}

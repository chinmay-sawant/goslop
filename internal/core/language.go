package core

import (
	"fmt"
	"strings"
)

// LanguageID identifies a supported analysis language.
type LanguageID int

const (
	// LanguageGo is the production default.
	LanguageGo LanguageID = iota
	// LanguagePython is reserved for multi-language WIP (plugin stub / fixtures).
	// Not registered in engine.DefaultRegistry until config language filtering lands.
	LanguagePython
)

// Aliases matching alternate naming in early scaffolds.
const (
	LangGo     = LanguageGo
	LangPython = LanguagePython
)

// String returns the canonical lowercase language id.
func (l LanguageID) String() string {
	switch l {
	case LanguageGo:
		return "go"
	case LanguagePython:
		return "python"
	default:
		return "unknown"
	}
}

// DefaultEnabledLanguages returns the production default enabled set: Go only.
// Phase 3 config merge uses this when the languages field is unset.
// Callers must treat the slice as read-only (do not append/mutate in place).
func DefaultEnabledLanguages() []LanguageID {
	return []LanguageID{LanguageGo}
}

// LanguageFromExtension maps a file extension (without the dot) to a language.
// Known extensions: "go" → Go; "py" → Python. Matching is case-insensitive.
func LanguageFromExtension(ext string) (LanguageID, bool) {
	switch strings.ToLower(strings.TrimPrefix(strings.TrimSpace(ext), ".")) {
	case "go":
		return LanguageGo, true
	case "py":
		return LanguagePython, true
	default:
		return 0, false
	}
}

// ParseLanguage parses a config/cache language name.
// Aliases: "go"; "python" or "py" for Python. Matching is case-insensitive.
func ParseLanguage(s string) (LanguageID, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "go":
		return LanguageGo, true
	case "python", "py":
		return LanguagePython, true
	default:
		return 0, false
	}
}

// ParseLanguages parses language name tokens into LanguageIDs for config merge.
//
// Behavior:
//   - Stable order: first-seen order of the input (after trim).
//   - Empty / whitespace-only tokens are skipped.
//   - Duplicates collapse (including aliases such as "python" then "py").
//   - Unknown tokens return a clear error; no partial result is returned.
//   - An empty list after skipping blanks returns (nil, nil) so callers can
//     decide whether to reject or fall back to DefaultEnabledLanguages.
func ParseLanguages(names []string) ([]LanguageID, error) {
	if len(names) == 0 {
		return nil, nil
	}
	out := make([]LanguageID, 0, len(names))
	seen := make(map[LanguageID]struct{}, len(names))
	for _, raw := range names {
		token := strings.TrimSpace(raw)
		if token == "" {
			continue
		}
		id, ok := ParseLanguage(token)
		if !ok {
			return nil, fmt.Errorf("unknown language %q (want go, python, or py)", token)
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

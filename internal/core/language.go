package core

import "strings"

// LanguageID identifies a supported analysis language.
type LanguageID int

const (
	// LanguageGo is the production default.
	LanguageGo LanguageID = iota
	// LanguagePython is reserved / deferred in the Go port.
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

// LanguageFromExtension maps a file extension (without the dot) to a language.
func LanguageFromExtension(ext string) (LanguageID, bool) {
	switch strings.ToLower(ext) {
	case "go":
		return LanguageGo, true
	case "py":
		return LanguagePython, true
	default:
		return 0, false
	}
}

// ParseLanguage parses a config/cache language name.
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

// DefaultEnabledLanguages is the product default when languages is unset in config.
// Production remains Go-only until additional plugins are registered and opted in.
func DefaultEnabledLanguages() []LanguageID {
	return []LanguageID{LanguageGo}
}

// ParseLanguages parses a list of language names into LanguageIDs.
// Order is stable (first occurrence wins); duplicates are dropped.
// Empty tokens are ignored. An empty or all-blank list returns an error —
// callers that want the product default should use DefaultEnabledLanguages
// when the config field is unset, not when it is explicitly empty.
func ParseLanguages(names []string) ([]LanguageID, error) {
	if len(names) == 0 {
		return nil, &LanguageError{Msg: "languages list must not be empty"}
	}
	seen := make(map[LanguageID]struct{}, len(names))
	out := make([]LanguageID, 0, len(names))
	for _, raw := range names {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		id, ok := ParseLanguage(raw)
		if !ok {
			return nil, &LanguageError{Msg: "unknown language " + quoteLang(raw)}
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	if len(out) == 0 {
		return nil, &LanguageError{Msg: "languages list must not be empty"}
	}
	return out, nil
}

// LanguageError is a configuration/parse error for language tokens.
type LanguageError struct {
	Msg string
}

func (e *LanguageError) Error() string {
	if e == nil {
		return "language error"
	}
	return e.Msg
}

func quoteLang(s string) string {
	return `"` + s + `"`
}

// LanguageEnabled reports whether id is in the enabled set.
func LanguageEnabled(enabled []LanguageID, id LanguageID) bool {
	for _, e := range enabled {
		if e == id {
			return true
		}
	}
	return false
}

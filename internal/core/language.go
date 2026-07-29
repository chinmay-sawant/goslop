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

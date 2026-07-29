// Package fixture parses and materializes goslop .txt text fixtures.
// Format matches the Rust goslop fixture module (header + "---" + body).
package fixture

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// FixtureExtension is the file extension for text fixtures (without the leading dot).
const FixtureExtension = "txt"

// FIXTURE_EXTENSION is a compatibility alias for FixtureExtension.
const FIXTURE_EXTENSION = FixtureExtension

// Separator between header metadata and source body.
const Separator = "---"

// FixtureLanguage is the language declared in a fixture header.
type FixtureLanguage string

const (
	LangGo     FixtureLanguage = "go"
	LangPython FixtureLanguage = "python"
)

// Extension returns the source-file extension for this language.
func (l FixtureLanguage) Extension() string {
	switch l {
	case LangGo:
		return "go"
	case LangPython:
		return "py"
	default:
		return "txt"
	}
}

// String returns the canonical header spelling.
func (l FixtureLanguage) String() string {
	return string(l)
}

// ParseLanguage parses a fixture language token (case-insensitive).
func ParseLanguage(s string) (FixtureLanguage, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "go":
		return LangGo, nil
	case "python", "py":
		return LangPython, nil
	default:
		return "", fmt.Errorf("unknown fixture language: %s", strings.TrimSpace(s))
	}
}

// TextFixture is a parsed .txt fixture (header + source body).
type TextFixture struct {
	Language FixtureLanguage
	Filename string
	Source   string
}

// DefaultFilename builds the output filename from the .txt path stem + language extension.
func DefaultFilename(txtPath string, language FixtureLanguage) string {
	stem := strings.TrimSuffix(filepath.Base(txtPath), filepath.Ext(txtPath))
	if stem == "" || stem == "." {
		stem = "fixture"
	}
	return stem + "." + language.Extension()
}

// ParseFixture parses raw .txt fixture contents.
//
// Header runs until the first occurrence of the "---" separator (Rust parity).
// Known keys: lang/language, file/filename. Unknown keys and "#" comments are ignored.
func ParseFixture(text string, txtPath string) (TextFixture, error) {
	header, body, err := splitHeaderBody(text)
	if err != nil {
		return TextFixture{}, err
	}

	var (
		language *FixtureLanguage
		filename string
	)

	for _, raw := range strings.Split(header, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "lang", "language":
			lang, err := ParseLanguage(value)
			if err != nil {
				return TextFixture{}, err
			}
			language = &lang
		case "file", "filename":
			filename = strings.TrimSpace(value)
		default:
			// ignore unknown keys (variant, expect, …)
		}
	}

	if language == nil {
		return TextFixture{}, errors.New("fixture header missing `lang:` (go | python)")
	}
	if filename == "" {
		filename = DefaultFilename(txtPath, *language)
	}
	// Match Rust: strip only leading newlines from the body (not other whitespace).
	source := strings.TrimPrefix(body, "\n")
	// Also handle body that starts with \r\n after separator on Windows-ish fixtures.
	source = strings.TrimPrefix(source, "\r\n")
	source = strings.TrimPrefix(source, "\r")

	return TextFixture{
		Language: *language,
		Filename: filename,
		Source:   source,
	}, nil
}

func splitHeaderBody(text string) (string, string, error) {
	idx := strings.Index(text, Separator)
	if idx < 0 {
		return "", "", fmt.Errorf("fixture must contain a `%s` separator between header and source", Separator)
	}
	header := strings.TrimSpace(text[:idx])
	body := text[idx+len(Separator):]
	return header, body, nil
}

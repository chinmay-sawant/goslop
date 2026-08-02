package perf

import (
	"path/filepath"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func fileDisplayPath(unit *core.ParsedUnit) string {
	if unit == nil {
		return ""
	}
	if unit.DisplayPath != "" {
		return unit.DisplayPath
	}
	return unit.Path
}

func isPythonTestFile(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, path := range []string{unit.Path, unit.DisplayPath} {
		base := filepath.Base(path)
		norm := filepath.ToSlash(path)
		if base == "conftest.py" || strings.HasPrefix(base, "test_") || strings.HasSuffix(base, "_test.py") || strings.Contains(norm, "/tests/") {
			return true
		}
	}
	return false
}

// isPythonOfflineToolsPath reports offline ETL/CLI tool scripts under tools/
// or backend/tools/ (or scripts/tools/). These intentionally load and build
// records once per run; hot-path allocation rules are not meaningful there.
func isPythonOfflineToolsPath(unit *core.ParsedUnit) bool {
	if unit == nil {
		return false
	}
	for _, path := range []string{unit.Path, unit.DisplayPath} {
		if path == "" {
			continue
		}
		norm := filepath.ToSlash(path)
		if strings.Contains(norm, "/tools/") || strings.HasPrefix(norm, "tools/") ||
			strings.Contains(norm, "/backend/tools/") || strings.HasPrefix(norm, "backend/tools/") {
			return true
		}
	}
	return false
}

func pushAt(unit *core.ParsedUnit, ruleID string, offset int, message string, out *[]rules.Finding) {
	meta := MetadataForID(ruleID)
	if unit == nil || meta == nil || out == nil {
		return
	}
	line, col := unit.LineCol(offset)
	rules.PushFinding(meta, fileDisplayPath(unit), line, col, message, out)
}

func pushLine(unit *core.ParsedUnit, ruleID string, line codeLine, needle, message string, out *[]rules.Finding) {
	offset := line.byte
	if at := strings.Index(line.text, needle); at >= 0 {
		offset += at
	}
	pushAt(unit, ruleID, offset, message, out)
}

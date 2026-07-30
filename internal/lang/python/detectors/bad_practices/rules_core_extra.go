package badpractices

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	metaByID["BP-PY-3"] = &rules.RuleMetadata{
		ID: "BP-PY-3", Title: "Raise Generic Exception",
		Description: "Code raises bare `Exception` or `BaseException` instead of a domain-specific error type.",
		Severity:    rules.SeverityLow, Pack: rules.PackBadPractice,
		Fix: "Prefer project-defined exceptions or specific stdlib types (ValueError, RuntimeError, OSError).",
	}
	metaByID["BP-PY-5"] = &rules.RuleMetadata{
		ID: "BP-PY-5", Title: "Wildcard Import",
		Description: "`from module import *` pollutes the namespace and breaks static analysis.",
		Severity:    rules.SeverityLow, Pack: rules.PackBadPractice,
		Fix: "Import explicit names; re-export via __all__ and named imports in package __init__ if needed.",
	}
	RegisterRule("BP-PY-3", detectBPPY3)
	RegisterRule("BP-PY-5", detectBPPY5)
}

// raiseGenericRe matches raise Exception / raise BaseException (call optional).
// Does not match raise ValueError, raise MyException, etc. (\b after the type name).
var raiseGenericRe = regexp.MustCompile(`\braise\s+(Exception|BaseException)\b`)

// BP-PY-3: raise Exception(...) / raise BaseException(...) outside tests.
func detectBPPY3(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-3")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.has("raise") {
		return
	}
	if !facts.has("Exception") && !facts.has("BaseException") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "raise ") {
			continue
		}
		if !raiseGenericRe.MatchString(t) {
			continue
		}
		// Double-check we did not match something like raise ExceptionalThing (word boundary handles).
		pushAt(unit, meta, line.byte, "raise Exception/BaseException loses type information; use a specific exception type", out)
	}
}

// fromImportStarRe matches "from <module> import *" (optional trailing comment already stripped).
var fromImportStarRe = regexp.MustCompile(`^from\s+\S+\s+import\s+\*\s*$`)

// BP-PY-5: from module import *
// Policy: flag all wildcard imports except package re-exports in __init__.py.
func detectBPPY5(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-5")
	if !facts.has("import *") && !strings.Contains(unit.Source, "import *") {
		return
	}
	// Allow __init__.py re-export patterns (common package public API).
	base := filepath.Base(fileDisplayPath(unit))
	if base == "__init__.py" {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.Contains(t, "import *") {
			continue
		}
		if fromImportStarRe.MatchString(t) {
			pushAt(unit, meta, line.byte, "wildcard import pollutes the namespace; import explicit names", out)
		}
	}
}

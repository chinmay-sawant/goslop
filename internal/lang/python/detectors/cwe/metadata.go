package cwe

import (
	"strconv"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Hand-authored metadata for the priority Python CWE batch (#52).
// Titles/names align with ruleset/python/chunks for these IDs.
// Pack is PackSecurity via PackFromRuleID("CWE-*").

func cweNumber(ruleID string) uint {
	number, err := strconv.ParseUint(strings.TrimPrefix(ruleID, "CWE-"), 10, 0)
	if err != nil {
		panic("invalid CWE rule ID: " + ruleID)
	}
	return uint(number)
}

var (
	// MetaCWE22 — path traversal (cwe-001-050.json).
	MetaCWE22 = rules.Meta(
		"CWE-22",
		"Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')",
		"The product uses external input to construct a pathname that is intended to identify a file or directory under a restricted parent, but does not properly neutralize special elements (e.g. ..) that can escape the restricted directory. In Python this commonly appears via open/pathlib/os.path.join on request-derived segments without resolve+prefix confinement or basename-only policy.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-22"), "Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')", "https://cwe.mitre.org/data/definitions/22.html")},
		"Constrain paths: os.path.basename for filename-only, or Path.resolve() then require startswith(root.resolve()); never join untrusted segments into open/remove sinks.",
	)

	// MetaCWE78 — OS command injection (cwe-051-100.json).
	MetaCWE78 = rules.Meta(
		"CWE-78",
		"Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection')",
		"The product constructs all or part of an OS command using externally-influenced input without neutralizing special elements. In Python this typically involves os.system, os.popen, or subprocess with shell=True and a dynamic command string.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-78"), "Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection')", "https://cwe.mitre.org/data/definitions/78.html")},
		"Prefer subprocess with a list argv and shell=False; never interpolate user input into shell strings.",
	)

	// MetaCWE79 — XSS (cwe-051-100.json).
	MetaCWE79 = rules.Meta(
		"CWE-79",
		"Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')",
		"The product does not neutralize user-controllable input before it is placed in web page output. In Python web apps this often appears as mark_safe/Markup/render_template_string with dynamic HTML rather than autoescaped render_template.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-79"), "Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')", "https://cwe.mitre.org/data/definitions/79.html")},
		"Use autoescaping templates (render_template) and avoid mark_safe/Markup/|safe on untrusted data; escape HTML when building responses manually.",
	)

	// MetaCWE89 — SQL injection (cwe-051-100.json).
	MetaCWE89 = rules.Meta(
		"CWE-89",
		"Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')",
		"The product constructs all or part of an SQL command using externally-influenced input without neutralizing special elements. In Python DB-API this appears when execute/executemany build SQL via f-strings, % formatting, or .format instead of bound parameters.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-89"), "Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')", "https://cwe.mitre.org/data/definitions/89.html")},
		"Use parameterized queries: cursor.execute(\"… ? …\", (value,)) or driver-native placeholders with a bound args sequence; never interpolate untrusted data into SQL strings.",
	)

	// MetaCWE502 — unsafe deserialization (cwe-501-550.json, python_relevance High).
	MetaCWE502 = rules.Meta(
		"CWE-502",
		"Deserialization of Untrusted Data",
		"The product deserializes untrusted data without sufficiently ensuring that the resulting data will be valid. In Python this commonly involves pickle.loads/load or yaml.load without a SafeLoader on attacker-controlled bytes.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-502"), "Deserialization of Untrusted Data", "https://cwe.mitre.org/data/definitions/502.html")},
		"Never unpickle untrusted data; use yaml.safe_load or Loader=yaml.SafeLoader; prefer json for untrusted interchange formats.",
	)
)

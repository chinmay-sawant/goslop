package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for narrow Python source heuristics that detect unsafe validation,
// export, and parsing boundaries. Each rule requires an executable, local
// pattern; none attempts to infer application-wide validation policy.
var (
	MetaCWE1173 = rules.Meta(
		"CWE-1173", "Improper Use of Validation Framework",
		"A route persists JSON request data without an observable schema or serializer validation step.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(cweNumber("CWE-1173"), "Improper Use of Validation Framework", "https://cwe.mitre.org/data/definitions/1173.html")},
		"Validate request data with an explicit schema or serializer before constructing and saving application models.",
	)
	MetaCWE1230 = rules.Meta(
		"CWE-1230", "Exposure of Sensitive Information Through Metadata",
		"A response Content-Disposition header includes a request-controlled original filename.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-1230"), "Exposure of Sensitive Information Through Metadata", "https://cwe.mitre.org/data/definitions/1230.html")},
		"Expose an application-generated download name or authorize and sanitize metadata before returning it.",
	)
	MetaCWE1236 = rules.Meta(
		"CWE-1236", "Improper Neutralization of Formula Elements in a CSV File",
		"A CSV row writes request-controlled fields without an observable spreadsheet-formula neutralization step.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(cweNumber("CWE-1236"), "Improper Neutralization of Formula Elements in a CSV File", "https://cwe.mitre.org/data/definitions/1236.html")},
		"Prefix or reject fields beginning with =, +, -, or @ before writing user data to CSV.",
	)
	MetaCWE1286 = rules.Meta(
		"CWE-1286", "Improper Validation of Syntactic Correctness of Input",
		"A request-controlled JSON value is used directly as an outbound URL without syntactic URL validation.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(cweNumber("CWE-1286"), "Improper Validation of Syntactic Correctness of Input", "https://cwe.mitre.org/data/definitions/1286.html")},
		"Parse and allowlist URL scheme, host, and syntax before making an outbound request.",
	)
	MetaCWE1289 = rules.Meta(
		"CWE-1289", "Improper Validation of Unsafe Equivalence in Input",
		"A route protects a request path with exact string equality before using it as a filesystem resource.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(cweNumber("CWE-1289"), "Improper Validation of Unsafe Equivalence in Input", "https://cwe.mitre.org/data/definitions/1289.html")},
		"Canonicalize the path and enforce that it remains beneath an approved base directory instead of using a deny-list equality check.",
	)
	MetaCWE1333 = rules.Meta(
		"CWE-1333", "Inefficient Regular Expression Complexity",
		"A compiled regular expression contains a directly visible nested quantifier with potentially catastrophic backtracking.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(cweNumber("CWE-1333"), "Inefficient Regular Expression Complexity", "https://cwe.mitre.org/data/definitions/1333.html")},
		"Replace nested unbounded quantifiers with a linear-time pattern and bound input length before matching.",
	)
	MetaCWE1389 = rules.Meta(
		"CWE-1389", "Incorrect Parsing of Numbers with Different Radices",
		"A request-controlled numeric value is parsed with Python base 0, accepting alternate radix prefixes.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-1389"), "Incorrect Parsing of Numbers with Different Radices", "https://cwe.mitre.org/data/definitions/1389.html")},
		"Parse externally supplied decimal values with base 10 and validate their allowed textual form first.",
	)
	MetaCWE140 = rules.Meta(
		"CWE-140", "Improper Neutralization of Delimiters",
		"A response manually joins request-controlled fields with a delimiter instead of using a structured encoder.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-140"), "Improper Neutralization of Delimiters", "https://cwe.mitre.org/data/definitions/140.html")},
		"Use a structured CSV or response encoder that quotes delimiter-bearing fields instead of concatenating user values.",
	)
)

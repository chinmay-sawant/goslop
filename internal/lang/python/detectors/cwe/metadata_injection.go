package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for the v0.0.2 injection expansion batch. Titles match the Python
// CWE catalogue chunks; the source-pattern detectors intentionally report only
// dynamic values at their documented sinks.
var (
	MetaCWE88 = rules.Meta(
		"CWE-88",
		"Improper Neutralization of Argument Delimiters in a Command ('Argument Injection')",
		"The product constructs a command argument vector using externally influenced input without delimiting intended arguments, options, or switches.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-88"), "Improper Neutralization of Argument Delimiters in a Command ('Argument Injection')", "https://cwe.mitre.org/data/definitions/88.html")},
		"Keep command arguments in a fixed allowlisted form; validate untrusted values and use -- before user-controlled operands where the invoked tool supports it.",
	)

	MetaCWE90 = rules.Meta(
		"CWE-90",
		"Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection')",
		"The product constructs an LDAP search filter with externally influenced input without escaping LDAP filter special characters.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-90"), "Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection')", "https://cwe.mitre.org/data/definitions/90.html")},
		"Use the LDAP library's filter escaping helper, such as ldap3.utils.conv.escape_filter_chars, before interpolating an untrusted value into a filter.",
	)

	MetaCWE91 = rules.Meta(
		"CWE-91",
		"XML Injection (aka Blind XPath Injection)",
		"The product constructs XML or XPath expressions from externally influenced input without neutralizing XML or XPath special characters.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-91"), "XML Injection (aka Blind XPath Injection)", "https://cwe.mitre.org/data/definitions/91.html")},
		"Use XML element-builder APIs and bind XPath variables instead of formatting untrusted values into XML or XPath expressions.",
	)

	MetaCWE93 = rules.Meta(
		"CWE-93",
		"Improper Neutralization of CRLF Sequences ('CRLF Injection')",
		"The product writes an externally influenced value into an HTTP response header without removing carriage-return and line-feed characters.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-93"), "Improper Neutralization of CRLF Sequences ('CRLF Injection')", "https://cwe.mitre.org/data/definitions/93.html")},
		"Reject CR and LF in header values, or remove them before setting a response header; prefer framework redirect and cookie APIs.",
	)

	MetaCWE94 = rules.Meta(
		"CWE-94",
		"Improper Control of Generation of Code ('Code Injection')",
		"The product passes externally influenced text to a Python code-generation or dynamic-import sink.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-94"), "Improper Control of Generation of Code ('Code Injection')", "https://cwe.mitre.org/data/definitions/94.html")},
		"Avoid eval, exec, compile, and dynamic imports for untrusted data; use an allowlist or a data-only parser instead.",
	)

	MetaCWE117 = rules.Meta(
		"CWE-117",
		"Improper Output Neutralization for Logs",
		"The product constructs a log message from externally influenced input without neutralizing line-breaking control characters.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-117"), "Improper Output Neutralization for Logs", "https://cwe.mitre.org/data/definitions/117.html")},
		"Use structured logging and reject or normalize CR and LF in externally controlled values before writing them to logs.",
	)
)

package cwe

import (
	"github.com/chinmay/codehound/internal/cwe"
	"github.com/chinmay/codehound/internal/rules"
)

// MetaCWE22 is catalogue metadata for path traversal.
var MetaCWE22 = rules.Meta(
	"CWE-22",
	"Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')",
	"The product uses external input to construct a pathname that is intended to identify a file or directory that is located underneath a restricted parent directory, but the product does not properly neutralize special elements within the pathname that can cause the pathname to resolve to a location that is outside of the restricted directory.",
	rules.SeverityHigh,
	[]cwe.CweRef{cwe.New(22, "Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')", "https://cwe.mitre.org/data/definitions/22.html")},
	"Confine paths with filepath.Base or Abs+HasPrefix checks; never trust Clean alone.",
)

// MetaCWE78 is catalogue metadata for OS command injection.
var MetaCWE78 = rules.Meta(
	"CWE-78",
	"Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection')",
	"The product constructs all or part of an OS command using externally-influenced input from an upstream component, but it does not neutralize or incorrectly neutralizes special elements that could modify the intended OS command when it is sent to a downstream component.",
	rules.SeverityHigh,
	[]cwe.CweRef{cwe.New(78, "Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection')", "https://cwe.mitre.org/data/definitions/78.html")},
	"Validate and allowlist command arguments; avoid shell (-c) with user input. Prefer exec.Command with fixed argv and no shell.",
)

// MetaCWE79 is catalogue metadata for XSS.
var MetaCWE79 = rules.Meta(
	"CWE-79",
	"Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')",
	"The product does not neutralize or incorrectly neutralizes user-controllable input before it is placed in output that is used as a web page that is served to other users.",
	rules.SeverityHigh,
	[]cwe.CweRef{cwe.New(79, "Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')", "https://cwe.mitre.org/data/definitions/79.html")},
	"Escape untrusted data with html.EscapeString or use html/template auto-escaping; never pass template.HTML from user input.",
)

// MetaCWE89 is catalogue metadata for SQL injection.
var MetaCWE89 = rules.Meta(
	"CWE-89",
	"Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')",
	"Improper neutralization of special elements used in an SQL command. In GORM and sqlx, this occurs when developers use Raw(), string formatting, or unsanitized input in Where/Order clauses instead of proper parameterized queries or the built-in escaping mechanisms.",
	rules.SeverityHigh,
	[]cwe.CweRef{cwe.New(89, "Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')", "https://cwe.mitre.org/data/definitions/89.html")},
	"Use parameterized queries or prepared statements instead of string formatting for SQL.",
)

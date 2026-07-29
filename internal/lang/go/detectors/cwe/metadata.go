package cwe

import (
	"github.com/chinmay/codehound/internal/cwe"
	"github.com/chinmay/codehound/internal/rules"
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

// MetaCWE89 is catalogue metadata for SQL injection.
var MetaCWE89 = rules.Meta(
	"CWE-89",
	"Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')",
	"Improper neutralization of special elements used in an SQL command. In GORM and sqlx, this occurs when developers use Raw(), string formatting, or unsanitized input in Where/Order clauses instead of proper parameterized queries or the built-in escaping mechanisms.",
	rules.SeverityHigh,
	[]cwe.CweRef{cwe.New(89, "Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')", "https://cwe.mitre.org/data/definitions/89.html")},
	"Use parameterized queries or prepared statements instead of string formatting for SQL.",
)

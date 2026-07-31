package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for the v0.0.2 information-exposure batch. These source-only
// heuristics deliberately use narrow, executable patterns rather than trying
// to infer whole-application authorization policy.
var (
	MetaCWE201 = rules.Meta(
		"CWE-201",
		"Insertion of Sensitive Information Into Sent Data",
		"A Python HTTP response directly includes a password, token, financial identifier, or other sensitive field.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-201"), "Insertion of Sensitive Information Into Sent Data", "https://cwe.mitre.org/data/definitions/201.html")},
		"Return an allowlisted response DTO that omits secrets and sensitive personal data.",
	)

	MetaCWE204 = rules.Meta(
		"CWE-204",
		"Observable Response Discrepancy",
		"A route returns distinct account-existence and password-failure messages, revealing authentication state to an unauthenticated caller.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-204"), "Observable Response Discrepancy", "https://cwe.mitre.org/data/definitions/204.html")},
		"Use one generic authentication-failure response for unknown accounts and invalid credentials.",
	)

	MetaCWE208 = rules.Meta(
		"CWE-208",
		"Observable Timing Discrepancy",
		"Security-sensitive values are compared with Python equality instead of a constant-time comparison API.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-208"), "Observable Timing Discrepancy", "https://cwe.mitre.org/data/definitions/208.html")},
		"Use hmac.compare_digest or secrets.compare_digest for secret and authentication-value comparisons.",
	)

	MetaCWE209 = rules.Meta(
		"CWE-209",
		"Generation of Error Message Containing Sensitive Information",
		"A Python HTTP response serializes an exception or traceback instead of a safe public error message.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-209"), "Generation of Error Message Containing Sensitive Information", "https://cwe.mitre.org/data/definitions/209.html")},
		"Log exception details privately and return a stable, generic error response to callers.",
	)

	MetaCWE212 = rules.Meta(
		"CWE-212",
		"Improper Removal of Sensitive Information Before Storage or Transfer",
		"A serialization sink stores a sensitive field without an allowlisted or redacted export representation.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-212"), "Improper Removal of Sensitive Information Before Storage or Transfer", "https://cwe.mitre.org/data/definitions/212.html")},
		"Create a redacted export DTO and serialize only the fields approved for transfer or storage.",
	)

	MetaCWE213 = rules.Meta(
		"CWE-213",
		"Exposure of Sensitive Information Due to Incompatible Policies",
		"A request to a guest, public, or anonymous endpoint includes a sensitive data field.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-213"), "Exposure of Sensitive Information Due to Incompatible Policies", "https://cwe.mitre.org/data/definitions/213.html")},
		"Apply the strictest stakeholder policy and exclude sensitive fields from public or guest integrations.",
	)

	MetaCWE488 = rules.Meta(
		"CWE-488",
		"Exposure of Data Element to Wrong Session",
		"A route writes request identity into a module-global variable, which can leak state across concurrent sessions.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-488"), "Exposure of Data Element to Wrong Session", "https://cwe.mitre.org/data/definitions/488.html")},
		"Keep request identity in framework request-local storage or an explicit per-request context, never a module global.",
	)

	MetaCWE497 = rules.Meta(
		"CWE-497",
		"Exposure of Sensitive System Information to an Unauthorized Control Sphere",
		"A Python HTTP response exposes environment, filesystem, hostname, or traceback details.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-497"), "Exposure of Sensitive System Information to an Unauthorized Control Sphere", "https://cwe.mitre.org/data/definitions/497.html")},
		"Keep system diagnostics in protected logs and return a generic public error response.",
	)
)

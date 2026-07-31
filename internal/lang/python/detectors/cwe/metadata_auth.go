package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for deliberately narrow Python authorization and session rules.
// The rules report direct framework-shaped source evidence only; they do not
// attempt to infer policy configured in a separate application module.
var (
	MetaCWE306 = rules.Meta(
		"CWE-306", "Missing Authentication for Critical Function",
		"A critical administration route is declared without a same-route authentication or authorization decorator.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(306, "Missing Authentication for Critical Function", "https://cwe.mitre.org/data/definitions/306.html")},
		"Require authentication and authorization for every critical route, using the framework's reviewed access-control mechanism.",
	)
	MetaCWE307 = rules.Meta(
		"CWE-307", "Improper Restriction of Excessive Authentication Attempts",
		"A password-authentication route has no same-route rate-limit or throttle decorator.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(307, "Improper Restriction of Excessive Authentication Attempts", "https://cwe.mitre.org/data/definitions/307.html")},
		"Apply a reviewed rate limit, throttle, or account-lockout policy to password-authentication endpoints.",
	)
	MetaCWE346 = rules.Meta(
		"CWE-346", "Origin Validation Error",
		"A CORS policy allows every origin while also allowing credentialed browser requests.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(346, "Origin Validation Error", "https://cwe.mitre.org/data/definitions/346.html")},
		"Allow credentials only for an explicit allowlist of trusted origins; never combine credentials with a wildcard origin.",
	)
	MetaCWE359 = rules.Meta(
		"CWE-359", "Exposure of Private Personal Information to an Unauthorized Actor",
		"The code sends a user personal-data field directly to a log or console output sink.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(359, "Exposure of Private Personal Information to an Unauthorized Actor", "https://cwe.mitre.org/data/definitions/359.html")},
		"Do not log passwords, government identifiers, contact details, or other personal data; log a non-sensitive identifier instead.",
	)
	MetaCWE613 = rules.Meta(
		"CWE-613", "Insufficient Session Expiration",
		"A session lifetime setting explicitly disables server-side session expiration.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(613, "Insufficient Session Expiration", "https://cwe.mitre.org/data/definitions/613.html")},
		"Set a short, server-enforced session lifetime and rotate or invalidate session credentials on logout and privilege changes.",
	)
	MetaCWE565 = rules.Meta(
		"CWE-565", "Reliance on Cookies without Validation and Integrity Checking",
		"A security-sensitive identity or privilege value is read from a request cookie without same-function validation evidence.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(565, "Reliance on Cookies without Validation and Integrity Checking", "https://cwe.mitre.org/data/definitions/565.html")},
		"Use a signed server-side session or validate cookie values against trusted server-side state before using them for authorization.",
	)
	MetaCWE807 = rules.Meta(
		"CWE-807", "Reliance on Untrusted Inputs in a Security Decision",
		"A request-controlled header, query parameter, or cookie directly controls an authorization decision.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(807, "Reliance on Untrusted Inputs in a Security Decision", "https://cwe.mitre.org/data/definitions/807.html")},
		"Derive authorization from a verified identity and server-side policy, not a client-controlled request value.",
	)
	MetaCWE698 = rules.Meta(
		"CWE-698", "Execution After Redirect (EAR)",
		"A redirect response is constructed without returning it and executable code follows in the same function.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(698, "Execution After Redirect (EAR)", "https://cwe.mitre.org/data/definitions/698.html")},
		"Return or raise the redirect response immediately so no protected action executes after redirecting the client.",
	)
)

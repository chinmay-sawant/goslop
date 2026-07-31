package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for the Python SSRF and communication-channel heuristics. These
// rules require direct, same-file evidence of a request-controlled value at a
// sensitive sink; they intentionally do not attempt interprocedural taint.
var (
	MetaCWE918 = rules.Meta(
		"CWE-918",
		"Server-Side Request Forgery (SSRF)",
		"The server passes a request-controlled URL to an outbound HTTP client without establishing an allowlisted destination.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(918, "Server-Side Request Forgery (SSRF)", "https://cwe.mitre.org/data/definitions/918.html")},
		"Allowlist outbound hosts and schemes, resolve and reject private or link-local addresses, and do not fetch client-supplied URLs directly.",
	)

	MetaCWE601 = rules.Meta(
		"CWE-601",
		"URL Redirection to Untrusted Site ('Open Redirect')",
		"The application passes a request-controlled URL directly to a redirect response without validating the destination.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(601, "URL Redirection to Untrusted Site ('Open Redirect')", "https://cwe.mitre.org/data/definitions/601.html")},
		"Redirect only to an allowlisted local path or host; reject absolute URLs and protocol-relative destinations supplied by clients.",
	)

	MetaCWE605 = rules.Meta(
		"CWE-605",
		"Multiple Binds to the Same Port",
		"The socket enables address or port reuse and then binds a wildcard interface, which can permit another process to bind the same port.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(605, "Multiple Binds to the Same Port", "https://cwe.mitre.org/data/definitions/605.html")},
		"Do not enable SO_REUSEADDR or SO_REUSEPORT for wildcard-bound services unless the sharing model is deliberate and access-controlled.",
	)

	MetaCWE924 = rules.Meta(
		"CWE-924",
		"Improper Enforcement of Message Integrity During Transmission in a Communication Channel",
		"A webhook-style handler consumes a received request body without a same-handler message-signature verification.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(924, "Improper Enforcement of Message Integrity During Transmission in a Communication Channel", "https://cwe.mitre.org/data/definitions/924.html")},
		"Verify a provider signature or HMAC over the raw request body with a constant-time comparison before processing a webhook event.",
	)

	MetaCWE940 = rules.Meta(
		"CWE-940",
		"Improper Verification of Source of a Communication Channel",
		"An authentication callback uses a request-controlled identity value to log in a user without a same-handler source or state verification.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(940, "Improper Verification of Source of a Communication Channel", "https://cwe.mitre.org/data/definitions/940.html")},
		"Bind callback identities to a server-side session and verify OAuth state or nonce values before creating an authenticated session.",
	)

	MetaCWE941 = rules.Meta(
		"CWE-941",
		"Incorrectly Specified Destination in a Communication Channel",
		"An outbound mail API receives a request-controlled recipient address without establishing an intended destination.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(941, "Incorrectly Specified Destination in a Communication Channel", "https://cwe.mitre.org/data/definitions/941.html")},
		"Use server-owned recipient identities or an explicit allowlist before sending sensitive data through an outbound communication channel.",
	)
)

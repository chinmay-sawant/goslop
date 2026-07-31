package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for source-only Python crypto and random-number heuristics. Each
// rule remains intentionally narrow because cryptographic intent often spans
// more than one function or file.
var (
	MetaCWE295 = rules.Meta(
		"CWE-295", "Improper Certificate Validation",
		"The code explicitly disables TLS certificate or hostname validation, allowing a network peer to present an untrusted certificate.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(295, "Improper Certificate Validation", "https://cwe.mitre.org/data/definitions/295.html")},
		"Keep certificate and hostname verification enabled; use a trusted CA bundle and fix the server certificate instead of bypassing validation.",
	)
	MetaCWE328 = rules.Meta(
		"CWE-328", "Use of Weak Hash",
		"The code constructs a legacy MD5 or SHA-1 digest that does not provide modern collision resistance for security-sensitive use.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(328, "Use of Weak Hash", "https://cwe.mitre.org/data/definitions/328.html")},
		"Use SHA-256 or SHA-3 for general hashing, or a purpose-built modern construction for the security protocol.",
	)
	MetaCWE335 = rules.Meta(
		"CWE-335", "Incorrect Usage of Seeds in Pseudo-Random Number Generator (PRNG)",
		"The code initializes a pseudo-random generator with a fixed source literal, making its generated sequence reproducible.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(335, "Incorrect Usage of Seeds in Pseudo-Random Number Generator (PRNG)", "https://cwe.mitre.org/data/definitions/335.html")},
		"Do not use a fixed seed for security-relevant random values; use the operating system CSPRNG through the secrets module.",
	)
	MetaCWE338 = rules.Meta(
		"CWE-338", "Use of Cryptographically Weak Pseudo-Random Number Generator (PRNG)",
		"The code uses Python's non-cryptographic random module to generate a security-sensitive value.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(338, "Use of Cryptographically Weak Pseudo-Random Number Generator (PRNG)", "https://cwe.mitre.org/data/definitions/338.html")},
		"Generate tokens, nonces, session identifiers, and secrets with secrets.token_urlsafe, secrets.token_bytes, or another CSPRNG.",
	)
	MetaCWE347 = rules.Meta(
		"CWE-347", "Improper Verification of Cryptographic Signature",
		"The code explicitly disables JWT signature verification, allowing unverified claims to be accepted.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(347, "Improper Verification of Cryptographic Signature", "https://cwe.mitre.org/data/definitions/347.html")},
		"Verify signatures with an expected algorithm and trusted key; never disable signature verification for untrusted tokens.",
	)
	MetaCWE1204 = rules.Meta(
		"CWE-1204", "Generation of Weak Initialization Vector (IV)",
		"A cipher invocation receives a fixed literal initialization vector, which can break confidentiality or integrity guarantees for IV-based modes.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(1204, "Generation of Weak Initialization Vector (IV)", "https://cwe.mitre.org/data/definitions/1204.html")},
		"Generate a fresh unpredictable IV with secrets.token_bytes for every encryption operation and store it with the ciphertext when required.",
	)
	MetaCWE1240 = rules.Meta(
		"CWE-1240", "Use of a Cryptographic Primitive with a Risky Implementation",
		"The code implements an XOR cipher directly instead of using a reviewed, standard cryptographic implementation.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(1240, "Use of a Cryptographic Primitive with a Risky Implementation", "https://cwe.mitre.org/data/definitions/1240.html")},
		"Use a well-reviewed library and a standard authenticated-encryption mode instead of a home-grown cipher.",
	)
	MetaCWE1241 = rules.Meta(
		"CWE-1241", "Use of Predictable Algorithm in Random Number Generator",
		"The code uses a predictable random algorithm to create a security-sensitive value.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(1241, "Use of Predictable Algorithm in Random Number Generator", "https://cwe.mitre.org/data/definitions/1241.html")},
		"Use the secrets module or another operating-system-backed CSPRNG for values that authenticate, authorize, or protect data.",
	)
	MetaCWE1392 = rules.Meta(
		"CWE-1392", "Use of Default Credentials",
		"The code assigns a common default password that attackers routinely try against deployed services.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(1392, "Use of Default Credentials", "https://cwe.mitre.org/data/definitions/1392.html")},
		"Require a unique deployment-specific credential from a protected configuration source and rotate any default value already in use.",
	)
)

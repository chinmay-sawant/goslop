package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for the Python secrets and credentials batch. The source-pattern
// rules intentionally require direct, reviewable credential or transport
// evidence rather than inferring values across functions or files.
var (
	MetaCWE798 = rules.Meta(
		"CWE-798",
		"Use of Hard-coded Credentials",
		"The code assigns a credential or secret directly to a Python source literal, making the credential available to anyone who can read the source.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(798, "Use of Hard-coded Credentials", "https://cwe.mitre.org/data/definitions/798.html")},
		"Load credentials from a managed secret store or protected environment variable; rotate any credential already committed to source.",
	)

	MetaCWE256 = rules.Meta(
		"CWE-256",
		"Plaintext Storage of a Password",
		"The code stores a password as a plaintext Python string literal instead of using a password-hash or protected credential mechanism.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(256, "Plaintext Storage of a Password", "https://cwe.mitre.org/data/definitions/256.html")},
		"Do not store passwords in source or plaintext resources; use a modern password hash for user passwords and a secret manager for service credentials.",
	)

	MetaCWE260 = rules.Meta(
		"CWE-260",
		"Password in Configuration File",
		"A Python configuration mapping contains a literal password that may be exposed to readers of the configuration source.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(260, "Password in Configuration File", "https://cwe.mitre.org/data/definitions/260.html")},
		"Read passwords from a protected secret provider or environment at deployment time rather than committing them to application configuration.",
	)

	MetaCWE261 = rules.Meta(
		"CWE-261",
		"Weak Encoding for Password",
		"The code applies a reversible encoding such as Base64, hexadecimal, or ROT13 to a password, which does not protect the password.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(261, "Weak Encoding for Password", "https://cwe.mitre.org/data/definitions/261.html")},
		"Use a modern password-hashing function such as Argon2id, bcrypt, scrypt, or PBKDF2; encoding is not password protection.",
	)

	MetaCWE312 = rules.Meta(
		"CWE-312",
		"Cleartext Storage of Sensitive Information",
		"The code stores a security-sensitive key or token as a cleartext Python source literal.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(312, "Cleartext Storage of Sensitive Information", "https://cwe.mitre.org/data/definitions/312.html")},
		"Keep keys and tokens in a protected secret manager or environment variable, and rotate values that were stored in source.",
	)

	MetaCWE319 = rules.Meta(
		"CWE-319",
		"Cleartext Transmission of Sensitive Information",
		"The code sends credentials through an HTTP or other cleartext protocol that can be observed by network attackers.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(319, "Cleartext Transmission of Sensitive Information", "https://cwe.mitre.org/data/definitions/319.html")},
		"Use TLS-protected protocols and HTTPS URLs whenever credentials or other sensitive values are transmitted.",
	)

	MetaCWE547 = rules.Meta(
		"CWE-547",
		"Use of Hard-coded, Security-relevant Constants",
		"A Django-style security setting is hard-coded to disable an important transport, cookie, host, or HSTS protection.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(547, "Use of Hard-coded, Security-relevant Constants", "https://cwe.mitre.org/data/definitions/547.html")},
		"Configure security-sensitive settings from a reviewed policy and keep HTTPS redirects, secure cookies, host validation, and HSTS enabled in production.",
	)

	MetaCWE523 = rules.Meta(
		"CWE-523",
		"Unprotected Transport of Credentials",
		"The code disables certificate validation or hostname verification while communicating over a credential-bearing transport.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(523, "Unprotected Transport of Credentials", "https://cwe.mitre.org/data/definitions/523.html")},
		"Keep TLS certificate and hostname verification enabled, and use HTTPS for every login or credential-bearing request.",
	)
)

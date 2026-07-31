package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for the Python mass-assignment and password-hash CWE batch.
// Titles align with ruleset/python/chunks/cwe-901-950.json.
var (
	MetaCWE914 = rules.Meta(
		"CWE-914",
		"Improper Control of Dynamically-Identified Variables",
		"The product does not properly restrict reading from or writing to dynamically-identified variables. In Python this can occur when request-controlled names are used to select globals, locals, object variables, or attributes.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-914"), "Improper Control of Dynamically-Identified Variables", "https://cwe.mitre.org/data/definitions/914.html")},
		"Do not use request-controlled names to select variables or attributes; map accepted names to a fixed allowlist.",
	)

	MetaCWE915 = rules.Meta(
		"CWE-915",
		"Improperly Controlled Modification of Dynamically-Determined Object Attributes",
		"The product receives input that specifies multiple attributes to initialize or update, but does not properly control which attributes can be modified. In Python this commonly appears as request-data mass assignment or a setattr loop over request keys.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-915"), "Improperly Controlled Modification of Dynamically-Determined Object Attributes", "https://cwe.mitre.org/data/definitions/915.html")},
		"Allowlist writable fields and explicitly construct the update payload; never pass request data directly as object attributes.",
	)

	MetaCWE916 = rules.Meta(
		"CWE-916",
		"Use of Password Hash With Insufficient Computational Effort",
		"The product hashes a password with a scheme that does not make password cracking sufficiently expensive. In Python this includes fast MD5 or SHA-1 password hashes and md5_crypt.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-916"), "Use of Password Hash With Insufficient Computational Effort", "https://cwe.mitre.org/data/definitions/916.html")},
		"Use a password-hashing function such as Argon2id, bcrypt, scrypt, or PBKDF2 with current parameters; do not use fast general-purpose hashes for passwords.",
	)
)

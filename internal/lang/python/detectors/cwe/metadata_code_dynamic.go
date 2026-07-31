package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for the v0.0.2 code-dynamic batch. Titles match the Python CWE
// catalogue chunks; the rules intentionally remain fixture-only until broader
// corpus calibration promotes them.
var (
	MetaCWE749 = rules.Meta(
		"CWE-749",
		"Exposed Dangerous Method or Function",
		"An externally reachable Python handler exposes dynamic code execution or another dangerous runtime capability without restricting access.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(749, "Exposed Dangerous Method or Function", "https://cwe.mitre.org/data/definitions/749.html")},
		"Do not expose eval, exec, dynamic imports, or shell execution through request handlers; use an explicit allowlist and authorization boundary.",
	)

	MetaCWE829 = rules.Meta(
		"CWE-829",
		"Inclusion of Functionality from Untrusted Control Sphere",
		"The code loads executable Python functionality from a dynamically selected module or file path, allowing an untrusted control sphere to select what executes.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(829, "Inclusion of Functionality from Untrusted Control Sphere", "https://cwe.mitre.org/data/definitions/829.html")},
		"Load only allowlisted, package-controlled modules and files; never pass request-derived module names or paths to import or runpy APIs.",
	)

	MetaCWE695 = rules.Meta(
		"CWE-695",
		"Use of Low-Level Functionality",
		"The code invokes native-memory or low-level runtime functionality that may bypass the framework or platform safety controls.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(695, "Use of Low-Level Functionality", "https://cwe.mitre.org/data/definitions/695.html")},
		"Avoid ctypes, cffi, and raw mmap interfaces unless the platform explicitly requires them and their boundary is reviewed.",
	)

	MetaCWE214 = rules.Meta(
		"CWE-214",
		"Invocation of Process Using Visible Sensitive Information",
		"A subprocess invocation passes credentials or another secret through command-line arguments or environment values visible to other processes.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(214, "Invocation of Process Using Visible Sensitive Information", "https://cwe.mitre.org/data/definitions/214.html")},
		"Pass secrets through protected files, standard input, or the platform's credential mechanism rather than process arguments or inherited environment values.",
	)

	MetaCWE215 = rules.Meta(
		"CWE-215",
		"Insertion of Sensitive Information Into Debugging Code",
		"Debug output logs a password, token, secret, API key, or credential that could be exposed when diagnostics are enabled.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(215, "Insertion of Sensitive Information Into Debugging Code", "https://cwe.mitre.org/data/definitions/215.html")},
		"Remove sensitive values from debug output and log only redacted identifiers or deliberately masked values.",
	)
)

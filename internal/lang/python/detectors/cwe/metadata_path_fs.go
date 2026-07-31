package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for narrow, source-only filesystem and execution-boundary checks.
// The corresponding rules report only explicit patterns that can be assessed
// within one Python source file; interprocedural provenance remains out of scope.
var (
	MetaCWE73 = rules.Meta(
		"CWE-73", "External Control of File Name or Path",
		"The code passes a directly request-controlled file name or path to a filesystem API without an evident same-expression restriction.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(73, "External Control of File Name or Path", "https://cwe.mitre.org/data/definitions/73.html")},
		"Constrain untrusted path components with an allowlist and resolve them beneath a trusted base directory before file access.",
	)
	MetaCWE59 = rules.Meta(
		"CWE-59", "Improper Link Resolution Before File Access ('Link Following')",
		"The code checks a path for a symbolic link and subsequently accesses it by name, leaving a check-then-use link-following race.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(59, "Improper Link Resolution Before File Access ('Link Following')", "https://cwe.mitre.org/data/definitions/59.html")},
		"Open the target using a descriptor API with O_NOFOLLOW where available and operate on that descriptor instead of checking a path before use.",
	)
	MetaCWE41 = rules.Meta(
		"CWE-41", "Improper Resolution of Path Equivalence",
		"The code normalizes a directly request-controlled path before file access without resolving it to a trusted canonical location.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(41, "Improper Resolution of Path Equivalence", "https://cwe.mitre.org/data/definitions/41.html")},
		"Resolve the canonical path and verify it remains within the intended trusted base directory rather than relying on normpath alone.",
	)
	MetaCWE276 = rules.Meta(
		"CWE-276", "Incorrect Default Permissions",
		"The code explicitly creates or changes filesystem permissions to a world-writable mode or disables the process umask.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(276, "Incorrect Default Permissions", "https://cwe.mitre.org/data/definitions/276.html")},
		"Use least-privilege file modes such as 0o600 or 0o640 and retain a restrictive process umask.",
	)
	MetaCWE378 = rules.Meta(
		"CWE-378", "Creation of Temporary File With Insecure Permissions",
		"The code creates a temporary pathname with tempfile.mktemp, which does not safely reserve the file before later use.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(378, "Creation of Temporary File With Insecure Permissions", "https://cwe.mitre.org/data/definitions/378.html")},
		"Use tempfile.NamedTemporaryFile, TemporaryFile, or mkstemp so the file is created securely and atomically.",
	)
	MetaCWE426 = rules.Meta(
		"CWE-426", "Untrusted Search Path",
		"The code adds the current directory or a directly request-controlled value to Python's import search path.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(426, "Untrusted Search Path", "https://cwe.mitre.org/data/definitions/426.html")},
		"Use a fixed, trusted absolute import location and do not prepend the current directory or externally supplied paths to sys.path.",
	)
	MetaCWE250 = rules.Meta(
		"CWE-250", "Execution with Unnecessary Privileges",
		"The code explicitly changes its effective user or group identity to root before continuing execution.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(250, "Execution with Unnecessary Privileges", "https://cwe.mitre.org/data/definitions/250.html")},
		"Run the service with the minimum required account privileges and drop elevated privileges before processing untrusted input.",
	)
	MetaCWE494 = rules.Meta(
		"CWE-494", "Download of Code Without Integrity Check",
		"The code immediately executes text downloaded through an HTTP client without an integrity-verification step.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(494, "Download of Code Without Integrity Check", "https://cwe.mitre.org/data/definitions/494.html")},
		"Do not execute downloaded code; if a download is unavoidable, authenticate it and verify a pinned cryptographic digest or signature before use.",
	)
)

package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for source-only resource, upload, and lifecycle heuristics. The
// patterns are intentionally narrow and do not infer cross-module policy.
var (
	MetaCWE434 = rules.Meta(
		"CWE-434", "Unrestricted Upload of File with Dangerous Type",
		"An uploaded Flask file is saved without a same-function dangerous-type allowlist decision.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(cweNumber("CWE-434"), "Unrestricted Upload of File with Dangerous Type", "https://cwe.mitre.org/data/definitions/434.html")},
		"Restrict accepted file types, generate server-side names, and normalize any retained filename with secure_filename before storing the upload.",
	)
	MetaCWE427 = rules.Meta(
		"CWE-427", "Uncontrolled Search Path Element",
		"A dynamic value is assigned to a process library or Python module search-path environment variable.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-427"), "Uncontrolled Search Path Element", "https://cwe.mitre.org/data/definitions/427.html")},
		"Use fixed, trusted absolute search paths and do not derive LD_LIBRARY_PATH or PYTHONPATH from external input.",
	)
	MetaCWE379 = rules.Meta(
		"CWE-379", "Creation of Temporary File in Directory with Insecure Permissions",
		"A predictable temporary pathname in a shared temporary directory is opened directly.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-379"), "Creation of Temporary File in Directory with Insecure Permissions", "https://cwe.mitre.org/data/definitions/379.html")},
		"Use tempfile.NamedTemporaryFile, TemporaryFile, or mkstemp with restrictive permissions instead of a predictable shared-directory pathname.",
	)
	MetaCWE459 = rules.Meta(
		"CWE-459", "Incomplete Cleanup",
		"A persistent temporary file is created without same-function unlink cleanup evidence.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-459"), "Incomplete Cleanup", "https://cwe.mitre.org/data/definitions/459.html")},
		"Remove persistent temporary files in a finally block or use a context-managed temporary-file API that owns cleanup.",
	)
	MetaCWE772 = rules.Meta(
		"CWE-772", "Missing Release of Resource after Effective Lifetime",
		"A file, socket, or URL response is assigned without same-function close or context-manager release evidence.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-772"), "Missing Release of Resource after Effective Lifetime", "https://cwe.mitre.org/data/definitions/772.html")},
		"Use a with statement or close the resource in a finally block after its effective lifetime ends.",
	)
	MetaCWE770 = rules.Meta(
		"CWE-770", "Allocation of Resources Without Limits or Throttling",
		"A Flask request body is read directly without module-level maximum content-length configuration evidence.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-770"), "Allocation of Resources Without Limits or Throttling", "https://cwe.mitre.org/data/definitions/770.html")},
		"Set MAX_CONTENT_LENGTH or an equivalent request-size limit before reading request bodies.",
	)
	MetaCWE708 = rules.Meta(
		"CWE-708", "Incorrect Ownership Assignment",
		"A resource is explicitly assigned to the root user or group.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-708"), "Incorrect Ownership Assignment", "https://cwe.mitre.org/data/definitions/708.html")},
		"Assign files to the least-privileged service account and avoid root ownership unless the boundary is explicitly reviewed.",
	)
	MetaCWE477 = rules.Meta(
		"CWE-477", "Use of Obsolete Function",
		"The code calls an obsolete or deprecated Python standard-library function.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-477"), "Use of Obsolete Function", "https://cwe.mitre.org/data/definitions/477.html")},
		"Replace obsolete standard-library APIs with their maintained alternatives and remove compatibility calls where no longer needed.",
	)
)

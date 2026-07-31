package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for narrow Python web-debug and configuration heuristics. These
// rules require explicit same-file configuration or framework evidence and do
// not infer deployment environment, permissions, or data flow.
var (
	MetaCWE756 = rules.Meta(
		"CWE-756", "Missing Custom Error Page",
		"The application explicitly enables debug error output, which can expose framework details instead of a custom error page.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(756, "Missing Custom Error Page", "https://cwe.mitre.org/data/definitions/756.html")},
		"Disable framework debug mode in deployed environments and configure generic error pages that do not disclose exception details.",
	)
	MetaCWE489 = rules.Meta(
		"CWE-489", "Active Debug Code",
		"The application enables debug mode or leaves an executable debugger breakpoint in source that can expose internals or interrupt service.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(489, "Active Debug Code", "https://cwe.mitre.org/data/definitions/489.html")},
		"Disable debug mode before deployment and remove pdb.set_trace or breakpoint calls from application code.",
	)
	MetaCWE15 = rules.Meta(
		"CWE-15", "External Control of System or Configuration Setting",
		"A request-controlled value directly changes an application or process configuration setting.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(15, "External Control of System or Configuration Setting", "https://cwe.mitre.org/data/definitions/15.html")},
		"Keep configuration under server control; validate and map client input to a narrow allowlist instead of assigning it to settings or environment values.",
	)
	MetaCWE1051 = rules.Meta(
		"CWE-1051", "Initialization with Hard-Coded Network Resource Configuration Data",
		"An outbound HTTP client uses a hard-coded private-network or localhost resource URL with an explicit port.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(1051, "Initialization with Hard-Coded Network Resource Configuration Data", "https://cwe.mitre.org/data/definitions/1051.html")},
		"Inject network endpoints through reviewed deployment configuration, and avoid embedding private addresses and ports in application source.",
	)
	MetaCWE1052 = rules.Meta(
		"CWE-1052", "Excessive Use of Hard-Coded Literals in Initialization",
		"Security or database initialization uses a hard-coded literal, making deployment-specific configuration difficult to review and rotate.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(1052, "Excessive Use of Hard-Coded Literals in Initialization", "https://cwe.mitre.org/data/definitions/1052.html")},
		"Load security keys and database connection data from a protected deployment configuration or secret manager.",
	)
	MetaCWE1125 = rules.Meta(
		"CWE-1125", "Excessive Attack Surface",
		"The application exposes an explicit debug-only HTTP route that increases its externally reachable attack surface.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(1125, "Excessive Attack Surface", "https://cwe.mitre.org/data/definitions/1125.html")},
		"Remove debug routes from deployed applications or protect them with strong server-side access controls and an operational allowlist.",
	)
	MetaCWE1188 = rules.Meta(
		"CWE-1188", "Initialization of a Resource with an Insecure Default",
		"A host-validation policy or network client is initialized with an insecure default that weakens access or TLS protection.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(1188, "Initialization of a Resource with an Insecure Default", "https://cwe.mitre.org/data/definitions/1188.html")},
		"Use explicit trusted hosts and keep TLS certificate verification enabled for all network clients.",
	)
	MetaCWE921 = rules.Meta(
		"CWE-921", "Storage of Sensitive Data in a Mechanism without Access Control",
		"A secret-shaped file is opened for writing in a common shared temporary directory without an access-control guarantee.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(921, "Storage of Sensitive Data in a Mechanism without Access Control", "https://cwe.mitre.org/data/definitions/921.html")},
		"Store sensitive data in an application-owned protected directory with restrictive permissions, not in shared temporary locations.",
	)
)

package cwe

import (
	"github.com/chinmay/codehound/internal/cwe"
	"github.com/chinmay/codehound/internal/rules"
)

// Generated catalogue metadata for registry CWE rules.

var (
	MetaCWE15 = rules.Meta(
		"CWE-15",
		"External Control of System or Configuration Setting",
		"One or more system settings or configuration elements can be externally controlled by a user. In Go + Gin + GORM applications, this commonly occurs when database connection strings, API keys, or runtime flags are loaded from environment variables or config files without proper validation or access control.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(15, "External Control of System or Configuration Setting", "https://cwe.mitre.org/data/definitions/15.html")},
		"Review and harden code against External Control of System or Configuration Setting (CWE-15); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE22 = rules.Meta(
		"CWE-22",
		"Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')",
		"Improper limitation of a pathname to a restricted directory (Path Traversal). Frequently seen in Gin applications when serving static files, handling uploads, or constructing file paths from user input without proper canonicalization and validation.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(22, "Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')", "https://cwe.mitre.org/data/definitions/22.html")},
		"Review and harden code against Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal') (CWE-22); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE41 = rules.Meta(
		"CWE-41",
		"Improper Resolution of Path Equivalence",
		"The product is vulnerable to file system contents disclosure through path equivalence. Path equivalence involves the use of special characters in file and directory names. The associated manipulations are intended to generate multiple names for the same object.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(41, "Improper Resolution of Path Equivalence", "https://cwe.mitre.org/data/definitions/41.html")},
		"Review and harden code against Improper Resolution of Path Equivalence (CWE-41); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE59 = rules.Meta(
		"CWE-59",
		"Improper Link Resolution Before File Access ('Link Following')",
		"The product attempts to access a file based on the filename, but it does not properly prevent that filename from identifying a link or shortcut that resolves to an unintended resource.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(59, "Improper Link Resolution Before File Access ('Link Following')", "https://cwe.mitre.org/data/definitions/59.html")},
		"Review and harden code against Improper Link Resolution Before File Access ('Link Following') (CWE-59); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE76 = rules.Meta(
		"CWE-76",
		"Improper Neutralization of Equivalent Special Elements",
		"The product correctly neutralizes certain special elements, but it improperly neutralizes equivalent special elements.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(76, "Improper Neutralization of Equivalent Special Elements", "https://cwe.mitre.org/data/definitions/76.html")},
		"Review and harden code against Improper Neutralization of Equivalent Special Elements (CWE-76); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE78 = rules.Meta(
		"CWE-78",
		"Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection')",
		"The product constructs all or part of an OS command using externally-influenced input from an upstream component, but it does not neutralize or incorrectly neutralizes special elements that could modify the intended OS command when it is sent to a downstream component.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(78, "Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection')", "https://cwe.mitre.org/data/definitions/78.html")},
		"Review and harden code against Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection') (CWE-78); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE79 = rules.Meta(
		"CWE-79",
		"Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')",
		"Improper neutralization of input during web page generation (XSS). Common in Gin when rendering templates with user-controlled data without proper escaping, returning JSON with unescaped HTML, or displaying error messages containing user input.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(79, "Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')", "https://cwe.mitre.org/data/definitions/79.html")},
		"Review and harden code against Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting') (CWE-79); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE89 = rules.Meta(
		"CWE-89",
		"Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')",
		"Improper neutralization of special elements used in an SQL command. In GORM and sqlx, this occurs when developers use Raw(), string formatting, or unsanitized input in Where/Order clauses instead of proper parameterized queries or the built-in escaping mechanisms.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(89, "Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')", "https://cwe.mitre.org/data/definitions/89.html")},
		"Review and harden code against Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection') (CWE-89); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE90 = rules.Meta(
		"CWE-90",
		"Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection')",
		"The product constructs all or part of an LDAP query using externally-influenced input from an upstream component, but it does not neutralize or incorrectly neutralizes special elements that could modify the intended LDAP query when it is sent to a downstream component.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(90, "Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection')", "https://cwe.mitre.org/data/definitions/90.html")},
		"Review and harden code against Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection') (CWE-90); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE91 = rules.Meta(
		"CWE-91",
		"XML Injection (aka Blind XPath Injection)",
		"The product does not properly neutralize special elements that are used in XML, allowing attackers to modify the syntax, content, or commands of the XML before it is processed by an end system.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(91, "XML Injection (aka Blind XPath Injection)", "https://cwe.mitre.org/data/definitions/91.html")},
		"Review and harden code against XML Injection (aka Blind XPath Injection) (CWE-91); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE93 = rules.Meta(
		"CWE-93",
		"Improper Neutralization of CRLF Sequences ('CRLF Injection')",
		"The product uses CRLF (carriage return line feeds) as a special element, e.g. to separate lines or records, but it does not neutralize or incorrectly neutralizes CRLF sequences from inputs.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(93, "Improper Neutralization of CRLF Sequences ('CRLF Injection')", "https://cwe.mitre.org/data/definitions/93.html")},
		"Review and harden code against Improper Neutralization of CRLF Sequences ('CRLF Injection') (CWE-93); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE112 = rules.Meta(
		"CWE-112",
		"Missing XML Validation",
		"The product accepts XML from an untrusted source but does not validate the XML against the proper schema.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(112, "Missing XML Validation", "https://cwe.mitre.org/data/definitions/112.html")},
		"Review and harden code against Missing XML Validation (CWE-112); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE140 = rules.Meta(
		"CWE-140",
		"Improper Neutralization of Delimiters",
		"The product does not neutralize or incorrectly neutralizes delimiters.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(140, "Improper Neutralization of Delimiters", "https://cwe.mitre.org/data/definitions/140.html")},
		"Review and harden code against Improper Neutralization of Delimiters (CWE-140); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE178 = rules.Meta(
		"CWE-178",
		"Improper Handling of Case Sensitivity",
		"The product does not properly account for differences in case sensitivity when accessing or determining the properties of a resource, leading to inconsistent results.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(178, "Improper Handling of Case Sensitivity", "https://cwe.mitre.org/data/definitions/178.html")},
		"Review and harden code against Improper Handling of Case Sensitivity (CWE-178); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE179 = rules.Meta(
		"CWE-179",
		"Incorrect Behavior Order: Early Validation",
		"The product validates input before applying protection mechanisms that modify the input, which could allow an attacker to bypass the validation via dangerous inputs that only arise after the modification.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(179, "Incorrect Behavior Order: Early Validation", "https://cwe.mitre.org/data/definitions/179.html")},
		"Review and harden code against Incorrect Behavior Order: Early Validation (CWE-179); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE182 = rules.Meta(
		"CWE-182",
		"Collapse of Data into Unsafe Value",
		"The product filters data in a way that causes it to be reduced or collapsed into an unsafe value that violates an expected security property.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(182, "Collapse of Data into Unsafe Value", "https://cwe.mitre.org/data/definitions/182.html")},
		"Review and harden code against Collapse of Data into Unsafe Value (CWE-182); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE184 = rules.Meta(
		"CWE-184",
		"Incomplete List of Disallowed Inputs",
		"The product implements a protection mechanism that relies on a list of inputs (or properties of inputs) that are not allowed by policy or otherwise require other action to neutralize before additional processing takes place, but the list is incomplete.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(184, "Incomplete List of Disallowed Inputs", "https://cwe.mitre.org/data/definitions/184.html")},
		"Review and harden code against Incomplete List of Disallowed Inputs (CWE-184); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE186 = rules.Meta(
		"CWE-186",
		"Overly Restrictive Regular Expression",
		"A regular expression is overly restrictive, which prevents dangerous values from being detected.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(186, "Overly Restrictive Regular Expression", "https://cwe.mitre.org/data/definitions/186.html")},
		"Review and harden code against Overly Restrictive Regular Expression (CWE-186); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE201 = rules.Meta(
		"CWE-201",
		"Insertion of Sensitive Information Into Sent Data",
		"The code transmits data to another actor, but a portion of the data includes sensitive information that should not be accessible to that actor.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(201, "Insertion of Sensitive Information Into Sent Data", "https://cwe.mitre.org/data/definitions/201.html")},
		"Review and harden code against Insertion of Sensitive Information Into Sent Data (CWE-201); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE204 = rules.Meta(
		"CWE-204",
		"Observable Response Discrepancy",
		"The product provides different responses to incoming requests in a way that reveals internal state information to an unauthorized actor outside of the intended control sphere.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(204, "Observable Response Discrepancy", "https://cwe.mitre.org/data/definitions/204.html")},
		"Review and harden code against Observable Response Discrepancy (CWE-204); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE208 = rules.Meta(
		"CWE-208",
		"Observable Timing Discrepancy",
		"Two separate operations in a product require different amounts of time to complete, in a way that is observable to an actor and reveals security-relevant information about the state of the product, such as whether a particular operation was successful or not.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(208, "Observable Timing Discrepancy", "https://cwe.mitre.org/data/definitions/208.html")},
		"Review and harden code against Observable Timing Discrepancy (CWE-208); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE209 = rules.Meta(
		"CWE-209",
		"Generation of Error Message Containing Sensitive Information",
		"The product generates an error message that includes sensitive information about its environment, users, or associated data.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(209, "Generation of Error Message Containing Sensitive Information", "https://cwe.mitre.org/data/definitions/209.html")},
		"Review and harden code against Generation of Error Message Containing Sensitive Information (CWE-209); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE212 = rules.Meta(
		"CWE-212",
		"Improper Removal of Sensitive Information Before Storage or Transfer",
		"The product stores, transfers, or shares a resource that contains sensitive information, but it does not properly remove that information before the product makes the resource available to unauthorized actors.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(212, "Improper Removal of Sensitive Information Before Storage or Transfer", "https://cwe.mitre.org/data/definitions/212.html")},
		"Review and harden code against Improper Removal of Sensitive Information Before Storage or Transfer (CWE-212); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE213 = rules.Meta(
		"CWE-213",
		"Exposure of Sensitive Information Due to Incompatible Policies",
		"The product's intended functionality exposes information to certain actors in accordance with the developer's security policy, but this information is regarded as sensitive according to the intended security policies of other stakeholders such as the product's administrator, users, or others whose i...",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(213, "Exposure of Sensitive Information Due to Incompatible Policies", "https://cwe.mitre.org/data/definitions/213.html")},
		"Review and harden code against Exposure of Sensitive Information Due to Incompatible Policies (CWE-213); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE214 = rules.Meta(
		"CWE-214",
		"Invocation of Process Using Visible Sensitive Information",
		"A process is invoked with sensitive command-line arguments, environment variables, or other elements that can be seen by other processes on the operating system.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(214, "Invocation of Process Using Visible Sensitive Information", "https://cwe.mitre.org/data/definitions/214.html")},
		"Review and harden code against Invocation of Process Using Visible Sensitive Information (CWE-214); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE215 = rules.Meta(
		"CWE-215",
		"Insertion of Sensitive Information Into Debugging Code",
		"The product inserts sensitive information into debugging code, which could expose this information if the debugging code is not disabled in production.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(215, "Insertion of Sensitive Information Into Debugging Code", "https://cwe.mitre.org/data/definitions/215.html")},
		"Review and harden code against Insertion of Sensitive Information Into Debugging Code (CWE-215); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE250 = rules.Meta(
		"CWE-250",
		"Execution with Unnecessary Privileges",
		"The product performs an operation at a privilege level that is higher than the minimum level required, which creates new weaknesses or amplifies the consequences of other weaknesses.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(250, "Execution with Unnecessary Privileges", "https://cwe.mitre.org/data/definitions/250.html")},
		"Review and harden code against Execution with Unnecessary Privileges (CWE-250); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE252 = rules.Meta(
		"CWE-252",
		"Unchecked Return Value",
		"The product does not check the return value from a method or function, which can prevent it from detecting unexpected states and conditions.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(252, "Unchecked Return Value", "https://cwe.mitre.org/data/definitions/252.html")},
		"Review and harden code against Unchecked Return Value (CWE-252); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE256 = rules.Meta(
		"CWE-256",
		"Plaintext Storage of a Password",
		"The product stores a password in plaintext within resources such as memory or files.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(256, "Plaintext Storage of a Password", "https://cwe.mitre.org/data/definitions/256.html")},
		"Review and harden code against Plaintext Storage of a Password (CWE-256); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE257 = rules.Meta(
		"CWE-257",
		"Storing Passwords in a Recoverable Format",
		"The storage of passwords in a recoverable format makes them subject to password reuse attacks by malicious users. In fact, it should be noted that recoverable encrypted passwords provide no significant benefit over plaintext passwords since they are subject not only to reuse by malicious attackers b...",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(257, "Storing Passwords in a Recoverable Format", "https://cwe.mitre.org/data/definitions/257.html")},
		"Review and harden code against Storing Passwords in a Recoverable Format (CWE-257); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE260 = rules.Meta(
		"CWE-260",
		"Password in Configuration File",
		"The product stores a password in a configuration file that might be accessible to actors who do not know the password.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(260, "Password in Configuration File", "https://cwe.mitre.org/data/definitions/260.html")},
		"Review and harden code against Password in Configuration File (CWE-260); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE261 = rules.Meta(
		"CWE-261",
		"Weak Encoding for Password",
		"Obscuring a password with a trivial encoding does not protect the password.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(261, "Weak Encoding for Password", "https://cwe.mitre.org/data/definitions/261.html")},
		"Review and harden code against Weak Encoding for Password (CWE-261); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE262 = rules.Meta(
		"CWE-262",
		"Not Using Password Aging",
		"The product does not have a mechanism in place for managing password aging.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(262, "Not Using Password Aging", "https://cwe.mitre.org/data/definitions/262.html")},
		"Review and harden code against Not Using Password Aging (CWE-262); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE263 = rules.Meta(
		"CWE-263",
		"Password Aging with Long Expiration",
		"The product supports password aging, but the expiration period is too long.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(263, "Password Aging with Long Expiration", "https://cwe.mitre.org/data/definitions/263.html")},
		"Review and harden code against Password Aging with Long Expiration (CWE-263); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE266 = rules.Meta(
		"CWE-266",
		"Incorrect Privilege Assignment",
		"A product incorrectly assigns a privilege to a particular actor, creating an unintended sphere of control for that actor.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(266, "Incorrect Privilege Assignment", "https://cwe.mitre.org/data/definitions/266.html")},
		"Review and harden code against Incorrect Privilege Assignment (CWE-266); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE267 = rules.Meta(
		"CWE-267",
		"Privilege Defined With Unsafe Actions",
		"A particular privilege, role, capability, or right can be used to perform unsafe actions that were not intended, even when it is assigned to the correct entity.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(267, "Privilege Defined With Unsafe Actions", "https://cwe.mitre.org/data/definitions/267.html")},
		"Review and harden code against Privilege Defined With Unsafe Actions (CWE-267); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE268 = rules.Meta(
		"CWE-268",
		"Privilege Chaining",
		"Two distinct privileges, roles, capabilities, or rights can be combined in a way that allows an entity to perform unsafe actions that would not be allowed without that combination.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(268, "Privilege Chaining", "https://cwe.mitre.org/data/definitions/268.html")},
		"Review and harden code against Privilege Chaining (CWE-268); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE270 = rules.Meta(
		"CWE-270",
		"Privilege Context Switching Error",
		"The product does not properly manage privileges while it is switching between different contexts that have different privileges or spheres of control.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(270, "Privilege Context Switching Error", "https://cwe.mitre.org/data/definitions/270.html")},
		"Review and harden code against Privilege Context Switching Error (CWE-270); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE272 = rules.Meta(
		"CWE-272",
		"Least Privilege Violation",
		"The elevated privilege level required to perform operations such as chroot() should be dropped immediately after the operation is performed.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(272, "Least Privilege Violation", "https://cwe.mitre.org/data/definitions/272.html")},
		"Review and harden code against Least Privilege Violation (CWE-272); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE273 = rules.Meta(
		"CWE-273",
		"Improper Check for Dropped Privileges",
		"The product attempts to drop privileges but does not check or incorrectly checks to see if the drop succeeded.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(273, "Improper Check for Dropped Privileges", "https://cwe.mitre.org/data/definitions/273.html")},
		"Review and harden code against Improper Check for Dropped Privileges (CWE-273); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE274 = rules.Meta(
		"CWE-274",
		"Improper Handling of Insufficient Privileges",
		"The product does not handle or incorrectly handles when it has insufficient privileges to perform an operation, leading to resultant weaknesses.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(274, "Improper Handling of Insufficient Privileges", "https://cwe.mitre.org/data/definitions/274.html")},
		"Review and harden code against Improper Handling of Insufficient Privileges (CWE-274); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE276 = rules.Meta(
		"CWE-276",
		"Incorrect Default Permissions",
		"During installation, installed file permissions are set to allow anyone to modify those files.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(276, "Incorrect Default Permissions", "https://cwe.mitre.org/data/definitions/276.html")},
		"Review and harden code against Incorrect Default Permissions (CWE-276); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE277 = rules.Meta(
		"CWE-277",
		"Insecure Inherited Permissions",
		"A product defines a set of insecure permissions that are inherited by objects that are created by the program.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(277, "Insecure Inherited Permissions", "https://cwe.mitre.org/data/definitions/277.html")},
		"Review and harden code against Insecure Inherited Permissions (CWE-277); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE278 = rules.Meta(
		"CWE-278",
		"Insecure Preserved Inherited Permissions",
		"A product inherits a set of insecure permissions for an object, e.g. when copying from an archive file, without user awareness or involvement.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(278, "Insecure Preserved Inherited Permissions", "https://cwe.mitre.org/data/definitions/278.html")},
		"Review and harden code against Insecure Preserved Inherited Permissions (CWE-278); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE279 = rules.Meta(
		"CWE-279",
		"Incorrect Execution-Assigned Permissions",
		"While it is executing, the product sets the permissions of an object in a way that violates the intended permissions that have been specified by the user.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(279, "Incorrect Execution-Assigned Permissions", "https://cwe.mitre.org/data/definitions/279.html")},
		"Review and harden code against Incorrect Execution-Assigned Permissions (CWE-279); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE280 = rules.Meta(
		"CWE-280",
		"Improper Handling of Insufficient Permissions or Privileges",
		"The product does not handle or incorrectly handles when it has insufficient privileges to access resources or functionality as specified by their permissions. This may cause it to follow unexpected code paths that may leave the product in an invalid state.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(280, "Improper Handling of Insufficient Permissions or Privileges", "https://cwe.mitre.org/data/definitions/280.html")},
		"Review and harden code against Improper Handling of Insufficient Permissions or Privileges (CWE-280); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE281 = rules.Meta(
		"CWE-281",
		"Improper Preservation of Permissions",
		"The product does not preserve permissions or incorrectly preserves permissions when copying, restoring, or sharing objects, which can cause them to have less restrictive permissions than intended.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(281, "Improper Preservation of Permissions", "https://cwe.mitre.org/data/definitions/281.html")},
		"Review and harden code against Improper Preservation of Permissions (CWE-281); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE283 = rules.Meta(
		"CWE-283",
		"Unverified Ownership",
		"The product does not properly verify that a critical resource is owned by the proper entity.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(283, "Unverified Ownership", "https://cwe.mitre.org/data/definitions/283.html")},
		"Review and harden code against Unverified Ownership (CWE-283); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE289 = rules.Meta(
		"CWE-289",
		"Authentication Bypass by Alternate Name",
		"The product performs authentication based on the name of a resource being accessed, or the name of the actor performing the access, but it does not properly check all possible names for that resource or actor.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(289, "Authentication Bypass by Alternate Name", "https://cwe.mitre.org/data/definitions/289.html")},
		"Review and harden code against Authentication Bypass by Alternate Name (CWE-289); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE290 = rules.Meta(
		"CWE-290",
		"Authentication Bypass by Spoofing",
		"This attack-focused weakness is caused by incorrectly implemented authentication schemes that are subject to spoofing attacks.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(290, "Authentication Bypass by Spoofing", "https://cwe.mitre.org/data/definitions/290.html")},
		"Review and harden code against Authentication Bypass by Spoofing (CWE-290); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE294 = rules.Meta(
		"CWE-294",
		"Authentication Bypass by Capture-replay",
		"A capture-replay flaw exists when the design of the product makes it possible for a malicious user to sniff network traffic and bypass authentication by replaying it to the server in question to the same effect as the original message (or with minor changes).",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(294, "Authentication Bypass by Capture-replay", "https://cwe.mitre.org/data/definitions/294.html")},
		"Review and harden code against Authentication Bypass by Capture-replay (CWE-294); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE301 = rules.Meta(
		"CWE-301",
		"Reflection Attack in an Authentication Protocol",
		"Simple authentication protocols are subject to reflection attacks if a malicious user can use the target machine to impersonate a trusted user.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(301, "Reflection Attack in an Authentication Protocol", "https://cwe.mitre.org/data/definitions/301.html")},
		"Review and harden code against Reflection Attack in an Authentication Protocol (CWE-301); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE303 = rules.Meta(
		"CWE-303",
		"Incorrect Implementation of Authentication Algorithm",
		"The requirements for the product dictate the use of an established authentication algorithm, but the implementation of the algorithm is incorrect.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(303, "Incorrect Implementation of Authentication Algorithm", "https://cwe.mitre.org/data/definitions/303.html")},
		"Review and harden code against Incorrect Implementation of Authentication Algorithm (CWE-303); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE305 = rules.Meta(
		"CWE-305",
		"Authentication Bypass by Primary Weakness",
		"The authentication algorithm is sound, but the implemented mechanism can be bypassed as the result of a separate weakness that is primary to the authentication error.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(305, "Authentication Bypass by Primary Weakness", "https://cwe.mitre.org/data/definitions/305.html")},
		"Review and harden code against Authentication Bypass by Primary Weakness (CWE-305); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE306 = rules.Meta(
		"CWE-306",
		"Missing Authentication for Critical Function",
		"The product does not perform any authentication for functionality that requires a provable user identity or consumes a significant amount of resources.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(306, "Missing Authentication for Critical Function", "https://cwe.mitre.org/data/definitions/306.html")},
		"Review and harden code against Missing Authentication for Critical Function (CWE-306); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE307 = rules.Meta(
		"CWE-307",
		"Improper Restriction of Excessive Authentication Attempts",
		"The product does not implement sufficient measures to prevent multiple failed authentication attempts within a short time frame.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(307, "Improper Restriction of Excessive Authentication Attempts", "https://cwe.mitre.org/data/definitions/307.html")},
		"Review and harden code against Improper Restriction of Excessive Authentication Attempts (CWE-307); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE308 = rules.Meta(
		"CWE-308",
		"Use of Single-factor Authentication",
		"The product uses an authentication algorithm that uses a single factor (e.g., a password) in a security context that should require more than one factor.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(308, "Use of Single-factor Authentication", "https://cwe.mitre.org/data/definitions/308.html")},
		"Review and harden code against Use of Single-factor Authentication (CWE-308); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE309 = rules.Meta(
		"CWE-309",
		"Use of Password System for Primary Authentication",
		"The use of password systems as the primary means of authentication may be subject to several flaws or shortcomings, each reducing the effectiveness of the mechanism.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(309, "Use of Password System for Primary Authentication", "https://cwe.mitre.org/data/definitions/309.html")},
		"Review and harden code against Use of Password System for Primary Authentication (CWE-309); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE312 = rules.Meta(
		"CWE-312",
		"Cleartext Storage of Sensitive Information",
		"The product stores sensitive information in cleartext within a resource that might be accessible to another control sphere.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(312, "Cleartext Storage of Sensitive Information", "https://cwe.mitre.org/data/definitions/312.html")},
		"Review and harden code against Cleartext Storage of Sensitive Information (CWE-312); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE319 = rules.Meta(
		"CWE-319",
		"Cleartext Transmission of Sensitive Information",
		"The product transmits sensitive or security-critical data in cleartext in a communication channel that can be sniffed by unauthorized actors.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(319, "Cleartext Transmission of Sensitive Information", "https://cwe.mitre.org/data/definitions/319.html")},
		"Review and harden code against Cleartext Transmission of Sensitive Information (CWE-319); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE322 = rules.Meta(
		"CWE-322",
		"Key Exchange without Entity Authentication",
		"The product performs a key exchange with an actor without verifying the identity of that actor.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(322, "Key Exchange without Entity Authentication", "https://cwe.mitre.org/data/definitions/322.html")},
		"Review and harden code against Key Exchange without Entity Authentication (CWE-322); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE323 = rules.Meta(
		"CWE-323",
		"Reusing a Nonce, Key Pair in Encryption",
		"Nonces should be used for the present occasion and only once.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(323, "Reusing a Nonce, Key Pair in Encryption", "https://cwe.mitre.org/data/definitions/323.html")},
		"Review and harden code against Reusing a Nonce, Key Pair in Encryption (CWE-323); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE324 = rules.Meta(
		"CWE-324",
		"Use of a Key Past its Expiration Date",
		"The product uses a cryptographic key or password past its expiration date, which diminishes its safety significantly by increasing the timing window for cracking attacks against that key.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(324, "Use of a Key Past its Expiration Date", "https://cwe.mitre.org/data/definitions/324.html")},
		"Review and harden code against Use of a Key Past its Expiration Date (CWE-324); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE325 = rules.Meta(
		"CWE-325",
		"Missing Cryptographic Step",
		"The product does not implement a required step in a cryptographic algorithm, resulting in weaker encryption than advertised by the algorithm.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(325, "Missing Cryptographic Step", "https://cwe.mitre.org/data/definitions/325.html")},
		"Review and harden code against Missing Cryptographic Step (CWE-325); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE328 = rules.Meta(
		"CWE-328",
		"Use of Weak Hash",
		"The product uses an algorithm that produces a digest (output value) that does not meet security expectations for a hash function that allows an adversary to reasonably determine the original input (preimage attack), find another input that can produce the same hash (2nd preimage attack), or find mul...",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(328, "Use of Weak Hash", "https://cwe.mitre.org/data/definitions/328.html")},
		"Review and harden code against Use of Weak Hash (CWE-328); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE331 = rules.Meta(
		"CWE-331",
		"Insufficient Entropy",
		"The product uses an algorithm or scheme that produces insufficient entropy, leaving patterns or clusters of values that are more likely to occur than others.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(331, "Insufficient Entropy", "https://cwe.mitre.org/data/definitions/331.html")},
		"Review and harden code against Insufficient Entropy (CWE-331); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE334 = rules.Meta(
		"CWE-334",
		"Small Space of Random Values",
		"The number of possible random values is smaller than needed by the product, making it more susceptible to brute force attacks.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(334, "Small Space of Random Values", "https://cwe.mitre.org/data/definitions/334.html")},
		"Review and harden code against Small Space of Random Values (CWE-334); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE335 = rules.Meta(
		"CWE-335",
		"Incorrect Usage of Seeds in Pseudo-Random Number Generator (PRNG)",
		"The product uses a Pseudo-Random Number Generator (PRNG) but does not correctly manage seeds.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(335, "Incorrect Usage of Seeds in Pseudo-Random Number Generator (PRNG)", "https://cwe.mitre.org/data/definitions/335.html")},
		"Review and harden code against Incorrect Usage of Seeds in Pseudo-Random Number Generator (PRNG) (CWE-335); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE338 = rules.Meta(
		"CWE-338",
		"Use of Cryptographically Weak Pseudo-Random Number Generator (PRNG)",
		"The product uses a Pseudo-Random Number Generator (PRNG) in a security context, but the PRNG's algorithm is not cryptographically strong.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(338, "Use of Cryptographically Weak Pseudo-Random Number Generator (PRNG)", "https://cwe.mitre.org/data/definitions/338.html")},
		"Review and harden code against Use of Cryptographically Weak Pseudo-Random Number Generator (PRNG) (CWE-338); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE341 = rules.Meta(
		"CWE-341",
		"Predictable from Observable State",
		"A number or object is predictable based on observations that the attacker can make about the state of the system or network, such as time, process ID, etc.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(341, "Predictable from Observable State", "https://cwe.mitre.org/data/definitions/341.html")},
		"Review and harden code against Predictable from Observable State (CWE-341); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE342 = rules.Meta(
		"CWE-342",
		"Predictable Exact Value from Previous Values",
		"An exact value or random number can be precisely predicted by observing previous values.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(342, "Predictable Exact Value from Previous Values", "https://cwe.mitre.org/data/definitions/342.html")},
		"Review and harden code against Predictable Exact Value from Previous Values (CWE-342); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE343 = rules.Meta(
		"CWE-343",
		"Predictable Value Range from Previous Values",
		"The product's random number generator produces a series of values which, when observed, can be used to infer a relatively small range of possibilities for the next value that could be generated.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(343, "Predictable Value Range from Previous Values", "https://cwe.mitre.org/data/definitions/343.html")},
		"Review and harden code against Predictable Value Range from Previous Values (CWE-343); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE344 = rules.Meta(
		"CWE-344",
		"Use of Invariant Value in Dynamically Changing Context",
		"The product uses a constant value, name, or reference, but this value can (or should) vary across different environments.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(344, "Use of Invariant Value in Dynamically Changing Context", "https://cwe.mitre.org/data/definitions/344.html")},
		"Review and harden code against Use of Invariant Value in Dynamically Changing Context (CWE-344); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE346 = rules.Meta(
		"CWE-346",
		"Origin Validation Error",
		"The product does not properly verify that the source of data or communication is valid.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(346, "Origin Validation Error", "https://cwe.mitre.org/data/definitions/346.html")},
		"Review and harden code against Origin Validation Error (CWE-346); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE347 = rules.Meta(
		"CWE-347",
		"Improper Verification of Cryptographic Signature",
		"The product does not verify, or incorrectly verifies, the cryptographic signature for data.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(347, "Improper Verification of Cryptographic Signature", "https://cwe.mitre.org/data/definitions/347.html")},
		"Review and harden code against Improper Verification of Cryptographic Signature (CWE-347); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE349 = rules.Meta(
		"CWE-349",
		"Acceptance of Extraneous Untrusted Data With Trusted Data",
		"The product, when processing trusted data, accepts any untrusted data that is also included with the trusted data, treating the untrusted data as if it were trusted.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(349, "Acceptance of Extraneous Untrusted Data With Trusted Data", "https://cwe.mitre.org/data/definitions/349.html")},
		"Review and harden code against Acceptance of Extraneous Untrusted Data With Trusted Data (CWE-349); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE353 = rules.Meta(
		"CWE-353",
		"Missing Support for Integrity Check",
		"The product uses a transmission protocol that does not include a mechanism for verifying the integrity of the data during transmission, such as a checksum.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(353, "Missing Support for Integrity Check", "https://cwe.mitre.org/data/definitions/353.html")},
		"Review and harden code against Missing Support for Integrity Check (CWE-353); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE356 = rules.Meta(
		"CWE-356",
		"Product UI does not Warn User of Unsafe Actions",
		"The product's user interface does not warn the user before undertaking an unsafe action on behalf of that user. This makes it easier for attackers to trick users into inflicting damage to their system.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(356, "Product UI does not Warn User of Unsafe Actions", "https://cwe.mitre.org/data/definitions/356.html")},
		"Review and harden code against Product UI does not Warn User of Unsafe Actions (CWE-356); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE358 = rules.Meta(
		"CWE-358",
		"Improperly Implemented Security Check for Standard",
		"The product does not implement or incorrectly implements one or more security-relevant checks as specified by the design of a standardized algorithm, protocol, or technique.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(358, "Improperly Implemented Security Check for Standard", "https://cwe.mitre.org/data/definitions/358.html")},
		"Review and harden code against Improperly Implemented Security Check for Standard (CWE-358); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE359 = rules.Meta(
		"CWE-359",
		"Exposure of Private Personal Information to an Unauthorized Actor",
		"The product does not properly prevent a person's private, personal information from being accessed by actors who either (1) are not explicitly authorized to access the information or (2) do not have the implicit consent of the person about whom the information is collected.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(359, "Exposure of Private Personal Information to an Unauthorized Actor", "https://cwe.mitre.org/data/definitions/359.html")},
		"Review and harden code against Exposure of Private Personal Information to an Unauthorized Actor (CWE-359); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE360 = rules.Meta(
		"CWE-360",
		"Trust of System Event Data",
		"Security based on event locations are insecure and can be spoofed.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(360, "Trust of System Event Data", "https://cwe.mitre.org/data/definitions/360.html")},
		"Review and harden code against Trust of System Event Data (CWE-360); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE366 = rules.Meta(
		"CWE-366",
		"Race Condition within a Thread",
		"If two threads of execution use a resource simultaneously, there exists the possibility that resources may be used while invalid, in turn making the state of execution undefined.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(366, "Race Condition within a Thread", "https://cwe.mitre.org/data/definitions/366.html")},
		"Review and harden code against Race Condition within a Thread (CWE-366); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE367 = rules.Meta(
		"CWE-367",
		"Time-of-check Time-of-use (TOCTOU) Race Condition",
		"The product checks the state of a resource before using that resource, but the resource's state can change between the check and the use in a way that invalidates the results of the check.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(367, "Time-of-check Time-of-use (TOCTOU) Race Condition", "https://cwe.mitre.org/data/definitions/367.html")},
		"Review and harden code against Time-of-check Time-of-use (TOCTOU) Race Condition (CWE-367); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE368 = rules.Meta(
		"CWE-368",
		"Context Switching Race Condition",
		"A product performs a series of non-atomic actions to switch between contexts that cross privilege or other security boundaries, but a race condition allows an attacker to modify or misrepresent the product's behavior during the switch.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(368, "Context Switching Race Condition", "https://cwe.mitre.org/data/definitions/368.html")},
		"Review and harden code against Context Switching Race Condition (CWE-368); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE378 = rules.Meta(
		"CWE-378",
		"Creation of Temporary File With Insecure Permissions",
		"Opening temporary files without appropriate measures or controls can leave the file, its contents and any function that it impacts vulnerable to attack.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(378, "Creation of Temporary File With Insecure Permissions", "https://cwe.mitre.org/data/definitions/378.html")},
		"Review and harden code against Creation of Temporary File With Insecure Permissions (CWE-378); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE379 = rules.Meta(
		"CWE-379",
		"Creation of Temporary File in Directory with Insecure Permissions",
		"The product creates a temporary file in a directory whose permissions allow unintended actors to determine the file's existence or otherwise access that file.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(379, "Creation of Temporary File in Directory with Insecure Permissions", "https://cwe.mitre.org/data/definitions/379.html")},
		"Review and harden code against Creation of Temporary File in Directory with Insecure Permissions (CWE-379); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE385 = rules.Meta(
		"CWE-385",
		"Covert Timing Channel",
		"Covert timing channels convey information by modulating some aspect of system behavior over time, so that the program receiving the information can observe system behavior and infer protected information.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(385, "Covert Timing Channel", "https://cwe.mitre.org/data/definitions/385.html")},
		"Review and harden code against Covert Timing Channel (CWE-385); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE393 = rules.Meta(
		"CWE-393",
		"Return of Wrong Status Code",
		"A function or operation returns an incorrect return value or status code that does not indicate the true result of execution, causing the product to modify its behavior based on the incorrect result.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(393, "Return of Wrong Status Code", "https://cwe.mitre.org/data/definitions/393.html")},
		"Review and harden code against Return of Wrong Status Code (CWE-393); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE403 = rules.Meta(
		"CWE-403",
		"Exposure of File Descriptor to Unintended Control Sphere ('File Descriptor Leak')",
		"A process does not close sensitive file descriptors before invoking a child process, which allows the child to perform unauthorized I/O operations using those descriptors.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(403, "Exposure of File Descriptor to Unintended Control Sphere ('File Descriptor Leak')", "https://cwe.mitre.org/data/definitions/403.html")},
		"Review and harden code against Exposure of File Descriptor to Unintended Control Sphere ('File Descriptor Leak') (CWE-403); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE408 = rules.Meta(
		"CWE-408",
		"Incorrect Behavior Order: Early Amplification",
		"The product allows an entity to perform a legitimate but expensive operation before authentication or authorization has taken place.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(408, "Incorrect Behavior Order: Early Amplification", "https://cwe.mitre.org/data/definitions/408.html")},
		"Review and harden code against Incorrect Behavior Order: Early Amplification (CWE-408); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE412 = rules.Meta(
		"CWE-412",
		"Unrestricted Externally Accessible Lock",
		"The product properly checks for the existence of a lock, but the lock can be externally controlled or influenced by an actor that is outside of the intended sphere of control.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(412, "Unrestricted Externally Accessible Lock", "https://cwe.mitre.org/data/definitions/412.html")},
		"Review and harden code against Unrestricted Externally Accessible Lock (CWE-412); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE420 = rules.Meta(
		"CWE-420",
		"Unprotected Alternate Channel",
		"The product protects a primary channel, but it does not use the same level of protection for an alternate channel.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(420, "Unprotected Alternate Channel", "https://cwe.mitre.org/data/definitions/420.html")},
		"Review and harden code against Unprotected Alternate Channel (CWE-420); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE421 = rules.Meta(
		"CWE-421",
		"Race Condition During Access to Alternate Channel",
		"The product opens an alternate channel to communicate with an authorized user, but the channel is accessible to other actors.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(421, "Race Condition During Access to Alternate Channel", "https://cwe.mitre.org/data/definitions/421.html")},
		"Review and harden code against Race Condition During Access to Alternate Channel (CWE-421); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE425 = rules.Meta(
		"CWE-425",
		"Direct Request ('Forced Browsing')",
		"The web application does not adequately enforce appropriate authorization on all restricted URLs, scripts, or files.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(425, "Direct Request ('Forced Browsing')", "https://cwe.mitre.org/data/definitions/425.html")},
		"Review and harden code against Direct Request ('Forced Browsing') (CWE-425); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE426 = rules.Meta(
		"CWE-426",
		"Untrusted Search Path",
		"The product searches for critical resources using an externally-supplied search path that can point to resources that are not under the product's direct control.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(426, "Untrusted Search Path", "https://cwe.mitre.org/data/definitions/426.html")},
		"Review and harden code against Untrusted Search Path (CWE-426); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE427 = rules.Meta(
		"CWE-427",
		"Uncontrolled Search Path Element",
		"The product uses a fixed or controlled search path to find resources, but one or more locations in that path can be under the control of unintended actors.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(427, "Uncontrolled Search Path Element", "https://cwe.mitre.org/data/definitions/427.html")},
		"Review and harden code against Uncontrolled Search Path Element (CWE-427); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE434 = rules.Meta(
		"CWE-434",
		"Unrestricted Upload of File with Dangerous Type",
		"The product allows the upload or transfer of dangerous file types that are automatically processed within its environment.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(434, "Unrestricted Upload of File with Dangerous Type", "https://cwe.mitre.org/data/definitions/434.html")},
		"Review and harden code against Unrestricted Upload of File with Dangerous Type (CWE-434); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE454 = rules.Meta(
		"CWE-454",
		"External Initialization of Trusted Variables or Data Stores",
		"The product initializes critical internal variables or data stores using inputs that can be modified by untrusted actors.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(454, "External Initialization of Trusted Variables or Data Stores", "https://cwe.mitre.org/data/definitions/454.html")},
		"Review and harden code against External Initialization of Trusted Variables or Data Stores (CWE-454); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE455 = rules.Meta(
		"CWE-455",
		"Non-exit on Failed Initialization",
		"The product does not exit or otherwise modify its operation when security-relevant errors occur during initialization, such as when a configuration file has a format error or a hardware security module (HSM) cannot be activated, which can cause the product to execute in a less secure fashion than in...",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(455, "Non-exit on Failed Initialization", "https://cwe.mitre.org/data/definitions/455.html")},
		"Review and harden code against Non-exit on Failed Initialization (CWE-455); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE459 = rules.Meta(
		"CWE-459",
		"Incomplete Cleanup",
		"The product does not properly clean up and remove temporary or supporting resources after they have been used.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(459, "Incomplete Cleanup", "https://cwe.mitre.org/data/definitions/459.html")},
		"Review and harden code against Incomplete Cleanup (CWE-459); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE472 = rules.Meta(
		"CWE-472",
		"External Control of Assumed-Immutable Web Parameter",
		"The web application does not sufficiently verify inputs that are assumed to be immutable but are actually externally controllable, such as hidden form fields.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(472, "External Control of Assumed-Immutable Web Parameter", "https://cwe.mitre.org/data/definitions/472.html")},
		"Review and harden code against External Control of Assumed-Immutable Web Parameter (CWE-472); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE488 = rules.Meta(
		"CWE-488",
		"Exposure of Data Element to Wrong Session",
		"The product does not sufficiently enforce boundaries between the states of different sessions, causing data to be provided to, or used by, the wrong session.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(488, "Exposure of Data Element to Wrong Session", "https://cwe.mitre.org/data/definitions/488.html")},
		"Review and harden code against Exposure of Data Element to Wrong Session (CWE-488); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE494 = rules.Meta(
		"CWE-494",
		"Download of Code Without Integrity Check",
		"The product downloads source code or an executable from a remote location and executes the code without sufficiently verifying the origin and integrity of the code.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(494, "Download of Code Without Integrity Check", "https://cwe.mitre.org/data/definitions/494.html")},
		"Review and harden code against Download of Code Without Integrity Check (CWE-494); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE497 = rules.Meta(
		"CWE-497",
		"Exposure of Sensitive System Information to an Unauthorized Control Sphere",
		"The product does not properly prevent sensitive system-level information from being accessed by unauthorized actors who do not have the same level of access to the underlying system as the product does.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(497, "Exposure of Sensitive System Information to an Unauthorized Control Sphere", "https://cwe.mitre.org/data/definitions/497.html")},
		"Review and harden code against Exposure of Sensitive System Information to an Unauthorized Control Sphere (CWE-497); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE501 = rules.Meta(
		"CWE-501",
		"Trust Boundary Violation",
		"The product mixes trusted and untrusted data in the same data structure or structured message.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(501, "Trust Boundary Violation", "https://cwe.mitre.org/data/definitions/501.html")},
		"Review and harden code against Trust Boundary Violation (CWE-501); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE502 = rules.Meta(
		"CWE-502",
		"Deserialization of Untrusted Data",
		"Deserialization of untrusted data. In Go + Gin, this occurs when unmarshaling JSON, gob, or XML from user requests into structs (including GORM models) without validation, potentially leading to object injection or logic bypass.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(502, "Deserialization of Untrusted Data", "https://cwe.mitre.org/data/definitions/502.html")},
		"Review and harden code against Deserialization of Untrusted Data (CWE-502); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE515 = rules.Meta(
		"CWE-515",
		"Covert Storage Channel",
		"A covert storage channel transfers information through the setting of bits by one program and the reading of those bits by another. What distinguishes this case from that of ordinary operation is that the bits are used to convey encoded information.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(515, "Covert Storage Channel", "https://cwe.mitre.org/data/definitions/515.html")},
		"Review and harden code against Covert Storage Channel (CWE-515); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE521 = rules.Meta(
		"CWE-521",
		"Weak Password Requirements",
		"The product does not require that users should have strong passwords.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(521, "Weak Password Requirements", "https://cwe.mitre.org/data/definitions/521.html")},
		"Review and harden code against Weak Password Requirements (CWE-521); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE523 = rules.Meta(
		"CWE-523",
		"Unprotected Transport of Credentials",
		"Login pages do not use adequate measures to protect the user name and password while they are in transit from the client to the server.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(523, "Unprotected Transport of Credentials", "https://cwe.mitre.org/data/definitions/523.html")},
		"Review and harden code against Unprotected Transport of Credentials (CWE-523); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE524 = rules.Meta(
		"CWE-524",
		"Use of Cache Containing Sensitive Information",
		"The code uses a cache that contains sensitive information, but the cache can be read by an actor outside of the intended control sphere.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(524, "Use of Cache Containing Sensitive Information", "https://cwe.mitre.org/data/definitions/524.html")},
		"Review and harden code against Use of Cache Containing Sensitive Information (CWE-524); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE538 = rules.Meta(
		"CWE-538",
		"Insertion of Sensitive Information into Externally-Accessible File or Directory",
		"The product places sensitive information into files or directories that are accessible to actors who are allowed to have access to the files, but not to the sensitive information.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(538, "Insertion of Sensitive Information into Externally-Accessible File or Directory", "https://cwe.mitre.org/data/definitions/538.html")},
		"Review and harden code against Insertion of Sensitive Information into Externally-Accessible File or Directory (CWE-538); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE544 = rules.Meta(
		"CWE-544",
		"Missing Standardized Error Handling Mechanism",
		"The product does not use a standardized method for handling errors throughout the code, which might introduce inconsistent error handling and resultant weaknesses.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(544, "Missing Standardized Error Handling Mechanism", "https://cwe.mitre.org/data/definitions/544.html")},
		"Review and harden code against Missing Standardized Error Handling Mechanism (CWE-544); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE547 = rules.Meta(
		"CWE-547",
		"Use of Hard-coded, Security-relevant Constants",
		"The product uses hard-coded constants instead of symbolic names for security-critical values, which increases the likelihood of mistakes during code maintenance or security policy change.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(547, "Use of Hard-coded, Security-relevant Constants", "https://cwe.mitre.org/data/definitions/547.html")},
		"Review and harden code against Use of Hard-coded, Security-relevant Constants (CWE-547); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE549 = rules.Meta(
		"CWE-549",
		"Missing Password Field Masking",
		"The product does not mask passwords during entry, increasing the potential for attackers to observe and capture passwords.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(549, "Missing Password Field Masking", "https://cwe.mitre.org/data/definitions/549.html")},
		"Review and harden code against Missing Password Field Masking (CWE-549); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE551 = rules.Meta(
		"CWE-551",
		"Incorrect Behavior Order: Authorization Before Parsing and Canonicalization",
		"If a web server does not fully parse requested URLs before it examines them for authorization, it may be possible for an attacker to bypass authorization protection.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(551, "Incorrect Behavior Order: Authorization Before Parsing and Canonicalization", "https://cwe.mitre.org/data/definitions/551.html")},
		"Review and harden code against Incorrect Behavior Order: Authorization Before Parsing and Canonicalization (CWE-551); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE552 = rules.Meta(
		"CWE-552",
		"Files or Directories Accessible to External Parties",
		"The product makes files or directories accessible to unauthorized actors, even though they should not be.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(552, "Files or Directories Accessible to External Parties", "https://cwe.mitre.org/data/definitions/552.html")},
		"Review and harden code against Files or Directories Accessible to External Parties (CWE-552); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE565 = rules.Meta(
		"CWE-565",
		"Reliance on Cookies without Validation and Integrity Checking",
		"The product relies on the existence or values of cookies when performing security-critical operations, but it does not properly ensure that the setting is valid for the associated user.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(565, "Reliance on Cookies without Validation and Integrity Checking", "https://cwe.mitre.org/data/definitions/565.html")},
		"Review and harden code against Reliance on Cookies without Validation and Integrity Checking (CWE-565); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE601 = rules.Meta(
		"CWE-601",
		"URL Redirection to Untrusted Site ('Open Redirect')",
		"The web application accepts a user-controlled input that specifies a link to an external site, and uses that link in a redirect.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(601, "URL Redirection to Untrusted Site ('Open Redirect')", "https://cwe.mitre.org/data/definitions/601.html")},
		"Review and harden code against URL Redirection to Untrusted Site ('Open Redirect') (CWE-601); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE603 = rules.Meta(
		"CWE-603",
		"Use of Client-Side Authentication",
		"A client/server product performs authentication within client code but not in server code, allowing server-side authentication to be bypassed via a modified client that omits the authentication check.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(603, "Use of Client-Side Authentication", "https://cwe.mitre.org/data/definitions/603.html")},
		"Review and harden code against Use of Client-Side Authentication (CWE-603); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE605 = rules.Meta(
		"CWE-605",
		"Multiple Binds to the Same Port",
		"When multiple sockets are allowed to bind to the same port, other services on that port may be stolen or spoofed.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(605, "Multiple Binds to the Same Port", "https://cwe.mitre.org/data/definitions/605.html")},
		"Review and harden code against Multiple Binds to the Same Port (CWE-605); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE611 = rules.Meta(
		"CWE-611",
		"Improper Restriction of XML External Entity Reference",
		"The product processes an XML document that can contain XML entities with URIs that resolve to documents outside of the intended sphere of control, causing the product to embed incorrect documents into its output.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(611, "Improper Restriction of XML External Entity Reference", "https://cwe.mitre.org/data/definitions/611.html")},
		"Review and harden code against Improper Restriction of XML External Entity Reference (CWE-611); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE613 = rules.Meta(
		"CWE-613",
		"Insufficient Session Expiration",
		"According to WASC, Insufficient Session Expiration is when a web site permits an attacker to reuse old session credentials or session IDs for authorization.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(613, "Insufficient Session Expiration", "https://cwe.mitre.org/data/definitions/613.html")},
		"Review and harden code against Insufficient Session Expiration (CWE-613); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE618 = rules.Meta(
		"CWE-618",
		"Exposed Unsafe ActiveX Method",
		"An ActiveX control is intended for use in a web browser, but it exposes dangerous methods that perform actions that are outside of the browser's security model (e.g. the zone or domain).",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(618, "Exposed Unsafe ActiveX Method", "https://cwe.mitre.org/data/definitions/618.html")},
		"Review and harden code against Exposed Unsafe ActiveX Method (CWE-618); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE619 = rules.Meta(
		"CWE-619",
		"Dangling Database Cursor ('Cursor Injection')",
		"If a database cursor is not closed properly, then it could become accessible to other users while retaining the same privileges that were originally assigned, leaving the cursor dangling.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(619, "Dangling Database Cursor ('Cursor Injection')", "https://cwe.mitre.org/data/definitions/619.html")},
		"Review and harden code against Dangling Database Cursor ('Cursor Injection') (CWE-619); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE620 = rules.Meta(
		"CWE-620",
		"Unverified Password Change",
		"When setting a new password for a user, the product does not require knowledge of the original password, or using another form of authentication.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(620, "Unverified Password Change", "https://cwe.mitre.org/data/definitions/620.html")},
		"Review and harden code against Unverified Password Change (CWE-620); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE639 = rules.Meta(
		"CWE-639",
		"Authorization Bypass Through User-Controlled Key",
		"The system's authorization functionality does not prevent one user from gaining access to another user's data or record by modifying the key value identifying the data.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(639, "Authorization Bypass Through User-Controlled Key", "https://cwe.mitre.org/data/definitions/639.html")},
		"Review and harden code against Authorization Bypass Through User-Controlled Key (CWE-639); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE640 = rules.Meta(
		"CWE-640",
		"Weak Password Recovery Mechanism for Forgotten Password",
		"The product contains a mechanism for users to recover or change their passwords without knowing the original password, but the mechanism is weak.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(640, "Weak Password Recovery Mechanism for Forgotten Password", "https://cwe.mitre.org/data/definitions/640.html")},
		"Review and harden code against Weak Password Recovery Mechanism for Forgotten Password (CWE-640); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE645 = rules.Meta(
		"CWE-645",
		"Overly Restrictive Account Lockout Mechanism",
		"The product contains an account lockout protection mechanism, but the mechanism is too restrictive and can be triggered too easily, which allows attackers to deny service to legitimate users by causing their accounts to be locked out.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(645, "Overly Restrictive Account Lockout Mechanism", "https://cwe.mitre.org/data/definitions/645.html")},
		"Review and harden code against Overly Restrictive Account Lockout Mechanism (CWE-645); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE648 = rules.Meta(
		"CWE-648",
		"Incorrect Use of Privileged APIs",
		"The product does not conform to the API requirements for a function call that requires extra privileges. This could allow attackers to gain privileges by causing the function to be called incorrectly.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(648, "Incorrect Use of Privileged APIs", "https://cwe.mitre.org/data/definitions/648.html")},
		"Review and harden code against Incorrect Use of Privileged APIs (CWE-648); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE649 = rules.Meta(
		"CWE-649",
		"Reliance on Obfuscation or Encryption of Security-Relevant Inputs without Integrity Checking",
		"The product uses obfuscation or encryption of inputs that should not be mutable by an external actor, but the product does not use integrity checks to detect if those inputs have been modified.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(649, "Reliance on Obfuscation or Encryption of Security-Relevant Inputs without Integrity Checking", "https://cwe.mitre.org/data/definitions/649.html")},
		"Review and harden code against Reliance on Obfuscation or Encryption of Security-Relevant Inputs without Integrity Checking (CWE-649); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE653 = rules.Meta(
		"CWE-653",
		"Improper Isolation or Compartmentalization",
		"The product does not properly compartmentalize or isolate functionality, processes, or resources that require different privilege levels, rights, or permissions.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(653, "Improper Isolation or Compartmentalization", "https://cwe.mitre.org/data/definitions/653.html")},
		"Review and harden code against Improper Isolation or Compartmentalization (CWE-653); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE654 = rules.Meta(
		"CWE-654",
		"Reliance on a Single Factor in a Security Decision",
		"A protection mechanism relies exclusively, or to a large extent, on the evaluation of a single condition or the integrity of a single object or entity in order to make a decision about granting access to restricted resources or functionality.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(654, "Reliance on a Single Factor in a Security Decision", "https://cwe.mitre.org/data/definitions/654.html")},
		"Review and harden code against Reliance on a Single Factor in a Security Decision (CWE-654); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE656 = rules.Meta(
		"CWE-656",
		"Reliance on Security Through Obscurity",
		"The product uses a protection mechanism whose strength depends heavily on its obscurity, such that knowledge of its algorithms or key data is sufficient to defeat the mechanism.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(656, "Reliance on Security Through Obscurity", "https://cwe.mitre.org/data/definitions/656.html")},
		"Review and harden code against Reliance on Security Through Obscurity (CWE-656); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE708 = rules.Meta(
		"CWE-708",
		"Incorrect Ownership Assignment",
		"The product assigns an owner to a resource, but the owner is outside of the intended control sphere.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(708, "Incorrect Ownership Assignment", "https://cwe.mitre.org/data/definitions/708.html")},
		"Review and harden code against Incorrect Ownership Assignment (CWE-708); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE756 = rules.Meta(
		"CWE-756",
		"Missing Custom Error Page",
		"The product does not return custom error pages to the user, possibly exposing sensitive information.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(756, "Missing Custom Error Page", "https://cwe.mitre.org/data/definitions/756.html")},
		"Review and harden code against Missing Custom Error Page (CWE-756); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE765 = rules.Meta(
		"CWE-765",
		"Multiple Unlocks of a Critical Resource",
		"The product unlocks a critical resource more times than intended, leading to an unexpected state in the system.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(765, "Multiple Unlocks of a Critical Resource", "https://cwe.mitre.org/data/definitions/765.html")},
		"Review and harden code against Multiple Unlocks of a Critical Resource (CWE-765); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE778 = rules.Meta(
		"CWE-778",
		"Insufficient Logging",
		"When a security-critical event occurs, the product either does not record the event or omits important details about the event when logging it.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(778, "Insufficient Logging", "https://cwe.mitre.org/data/definitions/778.html")},
		"Review and harden code against Insufficient Logging (CWE-778); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE783 = rules.Meta(
		"CWE-783",
		"Operator Precedence Logic Error",
		"The product uses an expression in which operator precedence causes incorrect logic to be used.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(783, "Operator Precedence Logic Error", "https://cwe.mitre.org/data/definitions/783.html")},
		"Review and harden code against Operator Precedence Logic Error (CWE-783); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE798 = rules.Meta(
		"CWE-798",
		"Use of Hard-coded Credentials",
		"The product contains hard-coded credentials, such as a password or cryptographic key.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(798, "Use of Hard-coded Credentials", "https://cwe.mitre.org/data/definitions/798.html")},
		"Review and harden code against Use of Hard-coded Credentials (CWE-798); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE807 = rules.Meta(
		"CWE-807",
		"Reliance on Untrusted Inputs in a Security Decision",
		"The product uses a protection mechanism that relies on the existence or values of an input, but the input can be modified by an untrusted actor in a way that bypasses the protection mechanism.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(807, "Reliance on Untrusted Inputs in a Security Decision", "https://cwe.mitre.org/data/definitions/807.html")},
		"Review and harden code against Reliance on Untrusted Inputs in a Security Decision (CWE-807); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE820 = rules.Meta(
		"CWE-820",
		"Missing Synchronization",
		"The product utilizes a shared resource in a concurrent manner but does not attempt to synchronize access to the resource.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(820, "Missing Synchronization", "https://cwe.mitre.org/data/definitions/820.html")},
		"Review and harden code against Missing Synchronization (CWE-820); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE821 = rules.Meta(
		"CWE-821",
		"Incorrect Synchronization",
		"The product utilizes a shared resource in a concurrent manner, but it does not correctly synchronize access to the resource.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(821, "Incorrect Synchronization", "https://cwe.mitre.org/data/definitions/821.html")},
		"Review and harden code against Incorrect Synchronization (CWE-821); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE826 = rules.Meta(
		"CWE-826",
		"Premature Release of Resource During Expected Lifetime",
		"The product releases a resource that is still intended to be used by itself or another actor.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(826, "Premature Release of Resource During Expected Lifetime", "https://cwe.mitre.org/data/definitions/826.html")},
		"Review and harden code against Premature Release of Resource During Expected Lifetime (CWE-826); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE829 = rules.Meta(
		"CWE-829",
		"Inclusion of Functionality from Untrusted Control Sphere",
		"The product imports, requires, or includes executable functionality (such as a library) from a source that is outside of the intended control sphere.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(829, "Inclusion of Functionality from Untrusted Control Sphere", "https://cwe.mitre.org/data/definitions/829.html")},
		"Review and harden code against Inclusion of Functionality from Untrusted Control Sphere (CWE-829); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE836 = rules.Meta(
		"CWE-836",
		"Use of Password Hash Instead of Password for Authentication",
		"The product records password hashes in a data store, receives a hash of a password from a client, and compares the supplied hash to the hash obtained from the data store.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(836, "Use of Password Hash Instead of Password for Authentication", "https://cwe.mitre.org/data/definitions/836.html")},
		"Review and harden code against Use of Password Hash Instead of Password for Authentication (CWE-836); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE838 = rules.Meta(
		"CWE-838",
		"Inappropriate Encoding for Output Context",
		"The product uses or specifies an encoding when generating output to a downstream component, but the specified encoding is not the same as the encoding that is expected by the downstream component.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(838, "Inappropriate Encoding for Output Context", "https://cwe.mitre.org/data/definitions/838.html")},
		"Review and harden code against Inappropriate Encoding for Output Context (CWE-838); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE841 = rules.Meta(
		"CWE-841",
		"Improper Enforcement of Behavioral Workflow",
		"The product supports a session in which more than one behavior must be performed by an actor, but it does not properly ensure that the actor performs the behaviors in the required sequence.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(841, "Improper Enforcement of Behavioral Workflow", "https://cwe.mitre.org/data/definitions/841.html")},
		"Review and harden code against Improper Enforcement of Behavioral Workflow (CWE-841); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE842 = rules.Meta(
		"CWE-842",
		"Placement of User into Incorrect Group",
		"The product or the administrator places a user into an incorrect group.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(842, "Placement of User into Incorrect Group", "https://cwe.mitre.org/data/definitions/842.html")},
		"Review and harden code against Placement of User into Incorrect Group (CWE-842); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE909 = rules.Meta(
		"CWE-909",
		"Missing Initialization of Resource",
		"The product does not initialize a critical resource.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(909, "Missing Initialization of Resource", "https://cwe.mitre.org/data/definitions/909.html")},
		"Review and harden code against Missing Initialization of Resource (CWE-909); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE915 = rules.Meta(
		"CWE-915",
		"Improperly Controlled Modification of Dynamically-Determined Object Attributes",
		"The product receives input from an upstream component that specifies multiple attributes, properties, or fields that are to be initialized or updated in an object, but it does not properly control which attributes can be modified.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(915, "Improperly Controlled Modification of Dynamically-Determined Object Attributes", "https://cwe.mitre.org/data/definitions/915.html")},
		"Review and harden code against Improperly Controlled Modification of Dynamically-Determined Object Attributes (CWE-915); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE916 = rules.Meta(
		"CWE-916",
		"Use of Password Hash With Insufficient Computational Effort",
		"The product generates a hash for a password, but it uses a scheme that does not provide a sufficient level of computational effort that would make password cracking attacks infeasible or expensive.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(916, "Use of Password Hash With Insufficient Computational Effort", "https://cwe.mitre.org/data/definitions/916.html")},
		"Review and harden code against Use of Password Hash With Insufficient Computational Effort (CWE-916); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE917 = rules.Meta(
		"CWE-917",
		"Improper Neutralization of Special Elements used in an Expression Language Statement ('Expression Language Injection')",
		"The product constructs all or part of an expression language (EL) statement in a framework such as a Java Server Page (JSP) using externally-influenced input from an upstream component, but it does not neutralize or incorrectly neutralizes special elements that could modify the intended EL statement...",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(917, "Improper Neutralization of Special Elements used in an Expression Language Statement ('Expression Language Injection')", "https://cwe.mitre.org/data/definitions/917.html")},
		"Review and harden code against Improper Neutralization of Special Elements used in an Expression Language Statement ('Expression Language Injection') (CWE-917); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE918 = rules.Meta(
		"CWE-918",
		"Server-Side Request Forgery (SSRF)",
		"Server-Side Request Forgery (SSRF). Common when Gin handlers make outgoing HTTP requests using URLs or hosts supplied by the client without proper allowlisting or validation.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(918, "Server-Side Request Forgery (SSRF)", "https://cwe.mitre.org/data/definitions/918.html")},
		"Review and harden code against Server-Side Request Forgery (SSRF) (CWE-918); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE921 = rules.Meta(
		"CWE-921",
		"Storage of Sensitive Data in a Mechanism without Access Control",
		"The product stores sensitive information in a file system or device that does not have built-in access control.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(921, "Storage of Sensitive Data in a Mechanism without Access Control", "https://cwe.mitre.org/data/definitions/921.html")},
		"Review and harden code against Storage of Sensitive Data in a Mechanism without Access Control (CWE-921); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE924 = rules.Meta(
		"CWE-924",
		"Improper Enforcement of Message Integrity During Transmission in a Communication Channel",
		"The product establishes a communication channel with an endpoint and receives a message from that endpoint, but it does not sufficiently ensure that the message was not modified during transmission.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(924, "Improper Enforcement of Message Integrity During Transmission in a Communication Channel", "https://cwe.mitre.org/data/definitions/924.html")},
		"Review and harden code against Improper Enforcement of Message Integrity During Transmission in a Communication Channel (CWE-924); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE940 = rules.Meta(
		"CWE-940",
		"Improper Verification of Source of a Communication Channel",
		"The product establishes a communication channel to handle an incoming request that has been initiated by an actor, but it does not properly verify that the request is coming from the expected origin.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(940, "Improper Verification of Source of a Communication Channel", "https://cwe.mitre.org/data/definitions/940.html")},
		"Review and harden code against Improper Verification of Source of a Communication Channel (CWE-940); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE941 = rules.Meta(
		"CWE-941",
		"Incorrectly Specified Destination in a Communication Channel",
		"The product creates a communication channel to initiate an outgoing request to an actor, but it does not correctly specify the intended destination for that actor.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(941, "Incorrectly Specified Destination in a Communication Channel", "https://cwe.mitre.org/data/definitions/941.html")},
		"Review and harden code against Incorrectly Specified Destination in a Communication Channel (CWE-941); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1051 = rules.Meta(
		"CWE-1051",
		"Initialization with Hard-Coded Network Resource Configuration Data",
		"The product initializes data using hard-coded values that act as network resource identifiers.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(1051, "Initialization with Hard-Coded Network Resource Configuration Data", "https://cwe.mitre.org/data/definitions/1051.html")},
		"Review and harden code against Initialization with Hard-Coded Network Resource Configuration Data (CWE-1051); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1052 = rules.Meta(
		"CWE-1052",
		"Excessive Use of Hard-Coded Literals in Initialization",
		"The product initializes a data element using a hard-coded literal that is not a simple integer or static constant element.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1052, "Excessive Use of Hard-Coded Literals in Initialization", "https://cwe.mitre.org/data/definitions/1052.html")},
		"Review and harden code against Excessive Use of Hard-Coded Literals in Initialization (CWE-1052); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1067 = rules.Meta(
		"CWE-1067",
		"Excessive Execution of Sequential Searches of Data Resource",
		"The product contains a data query against an SQL table or view that is configured in a way that does not utilize an index and may cause sequential searches to be performed.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(1067, "Excessive Execution of Sequential Searches of Data Resource", "https://cwe.mitre.org/data/definitions/1067.html")},
		"Review and harden code against Excessive Execution of Sequential Searches of Data Resource (CWE-1067); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1125 = rules.Meta(
		"CWE-1125",
		"Excessive Attack Surface",
		"The product has an attack surface whose quantitative measurement exceeds a desirable maximum.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1125, "Excessive Attack Surface", "https://cwe.mitre.org/data/definitions/1125.html")},
		"Review and harden code against Excessive Attack Surface (CWE-1125); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1173 = rules.Meta(
		"CWE-1173",
		"Improper Use of Validation Framework",
		"The product does not use, or incorrectly uses, an input validation framework that is provided by the source language or an independent library.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1173, "Improper Use of Validation Framework", "https://cwe.mitre.org/data/definitions/1173.html")},
		"Review and harden code against Improper Use of Validation Framework (CWE-1173); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1204 = rules.Meta(
		"CWE-1204",
		"Generation of Weak Initialization Vector (IV)",
		"The product uses a cryptographic primitive that uses an Initialization Vector (IV), but the product does not generate IVs that are sufficiently unpredictable or unique according to the expected cryptographic requirements for that primitive.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1204, "Generation of Weak Initialization Vector (IV)", "https://cwe.mitre.org/data/definitions/1204.html")},
		"Review and harden code against Generation of Weak Initialization Vector (IV) (CWE-1204); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1220 = rules.Meta(
		"CWE-1220",
		"Insufficient Granularity of Access Control",
		"The product implements access controls via a policy or other feature with the intention to disable or restrict accesses (reads and/or writes) to assets in a system from untrusted agents. However, implemented access controls lack required granularity, which renders the control policy too broad becaus...",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1220, "Insufficient Granularity of Access Control", "https://cwe.mitre.org/data/definitions/1220.html")},
		"Review and harden code against Insufficient Granularity of Access Control (CWE-1220); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1230 = rules.Meta(
		"CWE-1230",
		"Exposure of Sensitive Information Through Metadata",
		"The product prevents direct access to a resource containing sensitive information, but it does not sufficiently limit access to metadata that is derived from the original, sensitive information.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(1230, "Exposure of Sensitive Information Through Metadata", "https://cwe.mitre.org/data/definitions/1230.html")},
		"Review and harden code against Exposure of Sensitive Information Through Metadata (CWE-1230); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1236 = rules.Meta(
		"CWE-1236",
		"Improper Neutralization of Formula Elements in a CSV File",
		"The product saves user-provided information into a Comma-Separated Value (CSV) file, but it does not neutralize or incorrectly neutralizes special elements that could be interpreted as a command when the file is opened by a spreadsheet product.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1236, "Improper Neutralization of Formula Elements in a CSV File", "https://cwe.mitre.org/data/definitions/1236.html")},
		"Review and harden code against Improper Neutralization of Formula Elements in a CSV File (CWE-1236); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1240 = rules.Meta(
		"CWE-1240",
		"Use of a Cryptographic Primitive with a Risky Implementation",
		"To fulfill the need for a cryptographic primitive, the product implements a cryptographic algorithm using a non-standard, unproven, or disallowed/non-compliant cryptographic implementation.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1240, "Use of a Cryptographic Primitive with a Risky Implementation", "https://cwe.mitre.org/data/definitions/1240.html")},
		"Review and harden code against Use of a Cryptographic Primitive with a Risky Implementation (CWE-1240); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1265 = rules.Meta(
		"CWE-1265",
		"Unintended Reentrant Invocation of Non-reentrant Code Via Nested Calls",
		"The product invokes code that is believed to be reentrant, but the code performs a call that unintentionally produces a nested invocation of the non-reentrant code.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1265, "Unintended Reentrant Invocation of Non-reentrant Code Via Nested Calls", "https://cwe.mitre.org/data/definitions/1265.html")},
		"Review and harden code against Unintended Reentrant Invocation of Non-reentrant Code Via Nested Calls (CWE-1265); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1286 = rules.Meta(
		"CWE-1286",
		"Improper Validation of Syntactic Correctness of Input",
		"The product receives input that is expected to be well-formed - i.e., to comply with a certain syntax - but it does not validate or incorrectly validates that the input complies with the syntax.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(1286, "Improper Validation of Syntactic Correctness of Input", "https://cwe.mitre.org/data/definitions/1286.html")},
		"Review and harden code against Improper Validation of Syntactic Correctness of Input (CWE-1286); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1289 = rules.Meta(
		"CWE-1289",
		"Improper Validation of Unsafe Equivalence in Input",
		"The product receives an input value that is used as a resource identifier or other type of reference, but it does not validate or incorrectly validates that the input is equivalent to a potentially-unsafe value.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1289, "Improper Validation of Unsafe Equivalence in Input", "https://cwe.mitre.org/data/definitions/1289.html")},
		"Review and harden code against Improper Validation of Unsafe Equivalence in Input (CWE-1289); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1322 = rules.Meta(
		"CWE-1322",
		"Use of Blocking Code in Single-threaded, Non-blocking Context",
		"The product uses a non-blocking model that relies on a single threaded process for features such as scalability, but it contains code that can block when it is invoked.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(1322, "Use of Blocking Code in Single-threaded, Non-blocking Context", "https://cwe.mitre.org/data/definitions/1322.html")},
		"Review and harden code against Use of Blocking Code in Single-threaded, Non-blocking Context (CWE-1322); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1327 = rules.Meta(
		"CWE-1327",
		"Binding to an Unrestricted IP Address",
		"Binding to an unrestricted IP address. Very common Go mistake: calling http.ListenAndServe(\":8080\", router) or similar, which binds to 0.0.0.0 and exposes the service (including database ports) to the entire network.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1327, "Binding to an Unrestricted IP Address", "https://cwe.mitre.org/data/definitions/1327.html")},
		"Review and harden code against Binding to an Unrestricted IP Address (CWE-1327); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1333 = rules.Meta(
		"CWE-1333",
		"Inefficient Regular Expression Complexity",
		"Inefficient regular expression complexity (ReDoS). Go's regexp package can suffer from catastrophic backtracking when evil regex patterns are used to validate or filter user input in Gin handlers or GORM query builders.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(1333, "Inefficient Regular Expression Complexity", "https://cwe.mitre.org/data/definitions/1333.html")},
		"Review and harden code against Inefficient Regular Expression Complexity (CWE-1333); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1389 = rules.Meta(
		"CWE-1389",
		"Incorrect Parsing of Numbers with Different Radices",
		"The product parses numeric input assuming base 10 (decimal) values, but it does not account for inputs that use a different base number (radix).",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(1389, "Incorrect Parsing of Numbers with Different Radices", "https://cwe.mitre.org/data/definitions/1389.html")},
		"Review and harden code against Incorrect Parsing of Numbers with Different Radices (CWE-1389); validate inputs, use safe APIs, and follow secure defaults.",
	)
	MetaCWE1392 = rules.Meta(
		"CWE-1392",
		"Use of Default Credentials",
		"The product uses default credentials (such as passwords or cryptographic keys) for potentially critical functionality.",
		rules.SeverityCritical,
		[]cwe.CweRef{cwe.New(1392, "Use of Default Credentials", "https://cwe.mitre.org/data/definitions/1392.html")},
		"Review and harden code against Use of Default Credentials (CWE-1392); validate inputs, use safe APIs, and follow secure defaults.",
	)
)

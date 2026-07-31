package cwe

import (
	"strconv"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Tier-B metadata is kept in one domain-local file. Every rule is a narrow,
// source-only warning and deliberately has medium severity until a stronger
// semantic model can establish application-specific impact.
var (
	MetaCWE66   = tierBMeta("CWE-66", "Improper Handling of File Names that Identify Virtual Resources")
	MetaCWE76   = tierBMeta("CWE-76", "Improper Neutralization of Equivalent Special Elements")
	MetaCWE178  = tierBMeta("CWE-178", "Improper Handling of Case Sensitivity")
	MetaCWE179  = tierBMeta("CWE-179", "Incorrect Behavior Order: Early Validation")
	MetaCWE182  = tierBMeta("CWE-182", "Collapse of Data into Unsafe Value")
	MetaCWE184  = tierBMeta("CWE-184", "Incomplete List of Disallowed Inputs")
	MetaCWE186  = tierBMeta("CWE-186", "Overly Restrictive Regular Expression")
	MetaCWE257  = tierBMeta("CWE-257", "Storing Passwords in a Recoverable Format")
	MetaCWE272  = tierBMeta("CWE-272", "Least Privilege Violation")
	MetaCWE279  = tierBMeta("CWE-279", "Incorrect Execution-Assigned Permissions")
	MetaCWE289  = tierBMeta("CWE-289", "Authentication Bypass by Alternate Name")
	MetaCWE290  = tierBMeta("CWE-290", "Authentication Bypass by Spoofing")
	MetaCWE323  = tierBMeta("CWE-323", "Reusing a Nonce, Key Pair in Encryption")
	MetaCWE331  = tierBMeta("CWE-331", "Insufficient Entropy")
	MetaCWE334  = tierBMeta("CWE-334", "Small Space of Random Values")
	MetaCWE367  = tierBMeta("CWE-367", "Time-of-check Time-of-use (TOCTOU) Race Condition")
	MetaCWE403  = tierBMeta("CWE-403", "Exposure of File Descriptor to Unintended Control Sphere")
	MetaCWE409  = tierBMeta("CWE-409", "Improper Handling of Highly Compressed Data")
	MetaCWE454  = tierBMeta("CWE-454", "External Initialization of Trusted Variables or Data Stores")
	MetaCWE472  = tierBMeta("CWE-472", "External Control of Assumed-Immutable Web Parameter")
	MetaCWE521  = tierBMeta("CWE-521", "Weak Password Requirements")
	MetaCWE524  = tierBMeta("CWE-524", "Use of Cache Containing Sensitive Information")
	MetaCWE538  = tierBMeta("CWE-538", "Insertion of Sensitive Information into Externally-Accessible File or Directory")
	MetaCWE552  = tierBMeta("CWE-552", "Files or Directories Accessible to External Parties")
	MetaCWE617  = tierBMeta("CWE-617", "Reachable Assertion")
	MetaCWE641  = tierBMeta("CWE-641", "Improper Restriction of Names for Files and Other Resources")
	MetaCWE648  = tierBMeta("CWE-648", "Incorrect Use of Privileged APIs")
	MetaCWE779  = tierBMeta("CWE-779", "Logging of Excessive Data")
	MetaCWE836  = tierBMeta("CWE-836", "Use of Password Hash Instead of Password for Authentication")
	MetaCWE838  = tierBMeta("CWE-838", "Inappropriate Encoding for Output Context")
	MetaCWE908  = tierBMeta("CWE-908", "Use of Uninitialized Resource")
	MetaCWE909  = tierBMeta("CWE-909", "Missing Initialization of Resource")
	MetaCWE910  = tierBMeta("CWE-910", "Use of Expired File Descriptor")
	MetaCWE911  = tierBMeta("CWE-911", "Improper Update of Reference Count")
	MetaCWE920  = tierBMeta("CWE-920", "Improper Restriction of Power Consumption")
	MetaCWE939  = tierBMeta("CWE-939", "Improper Authorization in Handler for Custom URL Scheme")
	MetaCWE1007 = tierBMeta("CWE-1007", "Insufficient Visual Distinction of Homoglyphs Presented to User")
	MetaCWE1021 = tierBMeta("CWE-1021", "Improper Restriction of Rendered UI Layers or Frames")
	MetaCWE1046 = tierBMeta("CWE-1046", "Creation of Immutable Text Using String Concatenation")
	MetaCWE1050 = tierBMeta("CWE-1050", "Excessive Platform Resource Consumption within a Loop")
	MetaCWE1060 = tierBMeta("CWE-1060", "Excessive Number of Inefficient Server-Side Data Accesses")
	MetaCWE1067 = tierBMeta("CWE-1067", "Excessive Execution of Sequential Searches of Data Resource")
	MetaCWE1071 = tierBMeta("CWE-1071", "Empty Code Block")
	MetaCWE1072 = tierBMeta("CWE-1072", "Data Resource Access without Use of Connection Pooling")
	MetaCWE1084 = tierBMeta("CWE-1084", "Invokable Control Element with Excessive File or Data Access Operations")
	MetaCWE1104 = tierBMeta("CWE-1104", "Use of Unmaintained Third Party Components")
	MetaCWE1106 = tierBMeta("CWE-1106", "Insufficient Use of Symbolic Constants")
	MetaCWE1108 = tierBMeta("CWE-1108", "Excessive Reliance on Global Variables")
	MetaCWE1121 = tierBMeta("CWE-1121", "Excessive McCabe Cyclomatic Complexity")
	MetaCWE1123 = tierBMeta("CWE-1123", "Excessive Use of Self-Modifying Code")
	MetaCWE1124 = tierBMeta("CWE-1124", "Excessively Deep Nesting")
	MetaCWE1220 = tierBMeta("CWE-1220", "Insufficient Granularity of Access Control")
	MetaCWE1265 = tierBMeta("CWE-1265", "Unintended Reentrant Invocation of Non-reentrant Code Via Nested Calls")
	MetaCWE1284 = tierBMeta("CWE-1284", "Improper Validation of Specified Quantity in Input")
	MetaCWE1285 = tierBMeta("CWE-1285", "Improper Validation of Specified Index, Position, or Offset in Input")
	MetaCWE1287 = tierBMeta("CWE-1287", "Improper Validation of Specified Type of Input")
	MetaCWE1288 = tierBMeta("CWE-1288", "Improper Validation of Consistency within Input")
	MetaCWE1322 = tierBMeta("CWE-1322", "Use of Blocking Code in Single-threaded, Non-blocking Context")
	MetaCWE1339 = tierBMeta("CWE-1339", "Insufficient Precision or Accuracy of a Real Number")
	MetaCWE1341 = tierBMeta("CWE-1341", "Multiple Releases of Same Resource or Handle")
)

func tierBMeta(id, title string) rules.RuleMetadata {
	number, _ := strconv.Atoi(strings.TrimPrefix(id, "CWE-"))
	return rules.Meta(
		id, title,
		"A narrow Python source pattern indicates "+strings.ToLower(title)+".",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(uint(number), title, "https://cwe.mitre.org/data/definitions/"+strconv.Itoa(number)+".html")},
		"Review the local source pattern, validate the relevant trust boundary, and use the platform's documented safe API or control.",
	)
}

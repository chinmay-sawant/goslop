package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for explicit Python error-handling and control-flow patterns. The
// paired source heuristics intentionally require same-file syntax that can be
// assessed without guessing interprocedural error-handling behaviour.
var (
	MetaCWE396 = rules.Meta(
		"CWE-396", "Declaration of Catch for Generic Exception",
		"The code catches Exception or BaseException, which can hide failures that require distinct handling.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-396"), "Declaration of Catch for Generic Exception", "https://cwe.mitre.org/data/definitions/396.html")},
		"Catch the specific exception types the operation is expected to raise, and handle or propagate unexpected failures deliberately.",
	)
	MetaCWE397 = rules.Meta(
		"CWE-397", "Declaration of Throws for Generic Exception",
		"The code raises Exception or BaseException directly, making it difficult for callers to distinguish actionable failure conditions.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-397"), "Declaration of Throws for Generic Exception", "https://cwe.mitre.org/data/definitions/397.html")},
		"Raise a specific exception type that communicates the failed operation and lets callers handle it safely.",
	)
	MetaCWE478 = rules.Meta(
		"CWE-478", "Missing Default Case in Multiple Condition Expression",
		"A match statement has multiple cases but no wildcard default case, leaving unexpected values unhandled.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-478"), "Missing Default Case in Multiple Condition Expression", "https://cwe.mitre.org/data/definitions/478.html")},
		"Add a case _ branch that safely handles values outside the explicitly supported cases.",
	)
	MetaCWE252 = rules.Meta(
		"CWE-252", "Unchecked Return Value",
		"A process-execution call is used as a standalone statement, discarding its return status without checking success.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-252"), "Unchecked Return Value", "https://cwe.mitre.org/data/definitions/252.html")},
		"Inspect the returned CompletedProcess or use check=True so a failed command cannot silently continue execution.",
	)
	MetaCWE390 = rules.Meta(
		"CWE-390", "Detection of Error Condition Without Action",
		"An except clause handles an error condition only with pass, allowing the program to continue without any response.",
		rules.SeverityMedium, []cwe.CweRef{cwe.New(cweNumber("CWE-390"), "Detection of Error Condition Without Action", "https://cwe.mitre.org/data/definitions/390.html")},
		"Log, recover from, or re-raise the exception so failed operations are not silently ignored.",
	)
	MetaCWE584 = rules.Meta(
		"CWE-584", "Return Inside Finally Block",
		"A return statement appears directly in a finally block and can suppress an exception raised by the protected code.",
		rules.SeverityHigh, []cwe.CweRef{cwe.New(cweNumber("CWE-584"), "Return Inside Finally Block", "https://cwe.mitre.org/data/definitions/584.html")},
		"Perform cleanup in finally, then return after the try/finally block so exceptions can propagate normally.",
	)
)

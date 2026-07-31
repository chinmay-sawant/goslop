package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/cwe"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Metadata for the v0.0.2 XML batch. The rules report explicit parser
// configurations or request-to-parser flows only; broad XML parsing is not a
// finding by itself.
var (
	MetaCWE611 = rules.Meta(
		"CWE-611",
		"Improper Restriction of XML External Entity Reference",
		"An lxml parser explicitly enables XML entity resolution, allowing an XML document to resolve entity references outside the intended control sphere.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-611"), "Improper Restriction of XML External Entity Reference", "https://cwe.mitre.org/data/definitions/611.html")},
		"Use defusedxml where possible, or configure lxml XMLParser with resolve_entities=False and prohibit untrusted external DTD resolution.",
	)

	MetaCWE776 = rules.Meta(
		"CWE-776",
		"Improper Restriction of Recursive Entity References in DTDs ('XML Entity Expansion')",
		"An lxml parser enables entity expansion together with DTD loading or huge-tree processing, allowing recursive XML entity definitions to exhaust resources.",
		rules.SeverityHigh,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-776"), "Improper Restriction of Recursive Entity References in DTDs ('XML Entity Expansion')", "https://cwe.mitre.org/data/definitions/776.html")},
		"Disable entity resolution and DTD loading for untrusted XML; do not enable huge_tree unless strict input limits and a reviewed parser boundary are in place.",
	)

	MetaCWE112 = rules.Meta(
		"CWE-112",
		"Missing XML Validation",
		"Request-controlled XML is passed directly to a parser without a schema-aware parser configuration, allowing unexpected XML structure to bypass application assumptions.",
		rules.SeverityMedium,
		[]cwe.CweRef{cwe.New(cweNumber("CWE-112"), "Missing XML Validation", "https://cwe.mitre.org/data/definitions/112.html")},
		"Validate untrusted XML against the expected XML schema before using it, or configure the parser with the expected schema where the library supports it.",
	)
)

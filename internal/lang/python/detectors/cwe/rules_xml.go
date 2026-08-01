package cwe

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// These rules intentionally have no SourceIndex gates. Their aliases and
	// bare XMLParser forms are valid only after source-level import checks, and
	// an incomplete gate could skip a real finding.
	RegisterRule("CWE-611", detectCWE611, &MetaCWE611)
	RegisterRule("CWE-776", detectCWE776, &MetaCWE776)
	RegisterRule("CWE-112", detectCWE112, &MetaCWE112)
}

// CWE-611 reports an explicit lxml entity-resolution opt-in. lxml's parser
// default and ordinary stdlib XML parsing are intentionally not reported: this
// source-only rule requires a configuration that directly enables expansion.
func detectCWE611(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range lxmlXMLParserCalls(facts, unit.Source) {
		if !hasKwargTrue(call.ArgsText, "resolve_entities") {
			continue
		}
		pushXMLFinding(unit, &MetaCWE611, call.Start,
			"lxml XMLParser explicitly enables external entity resolution", confidence86, out)
		return
	}
}

// CWE-776 stays narrower than CWE-611: recursive expansion needs an entity
// resolving parser plus a DTD-loading or huge-tree configuration. Explicitly
// disabled entity resolution is therefore a safe suppression.
func detectCWE776(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range lxmlXMLParserCalls(facts, unit.Source) {
		if !hasKwargTrue(call.ArgsText, "resolve_entities") ||
			(!hasKwargTrue(call.ArgsText, "load_dtd") && !hasKwargTrue(call.ArgsText, "huge_tree")) {
			continue
		}
		pushXMLFinding(unit, &MetaCWE776, call.Start,
			"lxml XMLParser enables entity expansion with DTD or huge-tree processing", confidence84, out)
		return
	}
}

// CWE-112 detects only a direct request-controlled value passed to a known XML
// parsing API. Schema-aware parser arguments and a same-file named parser
// configuration are intentionally suppressed; following separately validated
// variables would require data flow beyond this source heuristic.
func detectCWE112(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source,
		"etree.parse", "etree.fromstring", "etree.XML",
		"lxml.etree.parse", "lxml.etree.fromstring", "lxml.etree.XML",
		"ElementTree.parse", "ElementTree.fromstring", "ET.parse", "ET.fromstring",
		"minidom.parse", "minidom.parseString", "xml.dom.minidom.parse", "xml.dom.minidom.parseString",
		"xml.sax.parse", "xml.sax.parseString") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !looksRequestControlledXML(args[0]) || xmlParserHasSchema(facts, unit.Source, args) {
			continue
		}
		pushXMLFinding(unit, &MetaCWE112, call.Start,
			"request-controlled XML reaches a parser without a schema-aware validation configuration", confidence74, out)
		return
	}
}

func lxmlXMLParserCalls(facts *PyCweFacts, source string) []callSite {
	var code string
	if facts != nil {
		code = facts.codeMask(source, fragStartHint(facts, source))
	} else {
		code = pythonCodeMask(source)
	}
	// A bare XMLParser is only treated as lxml when its import is visible in the
	// same file. This avoids confusing a project helper or the stdlib parser
	// with lxml's entity-resolution configuration.
	if !strings.Contains(code, "lxml") {
		return nil
	}
	return findCalls(facts, source, "lxml.etree.XMLParser", "etree.XMLParser", "XMLParser")
}

func looksRequestControlledXML(expr string) bool {
	compact := compactWhitespace(expr)
	return strings.Contains(compact, "request.data") ||
		strings.Contains(compact, "request.body") ||
		strings.Contains(compact, "request.files") ||
		strings.Contains(compact, "request.get_data(") ||
		strings.Contains(compact, "request.stream")
}

func xmlParserHasSchema(facts *PyCweFacts, source string, args []string) bool {
	for _, arg := range args {
		compact := strings.ToLower(compactWhitespace(arg))
		if strings.Contains(compact, "schema=") || strings.Contains(compact, "xmlschema") || strings.Contains(compact, "relaxng") {
			return true
		}
	}
	// A parser variable is a common lxml form. This is deliberately a narrow
	// same-file suppression: only the conventional `parser` name and explicit
	// XMLParser(schema=...) setup are recognized.
	for _, arg := range args[1:] {
		if strings.Contains(compactWhitespace(arg), "parser=parser") {
			var masked string
			if facts != nil {
				masked = facts.codeMask(source, fragStartHint(facts, source))
			} else {
				masked = pythonCodeMask(source)
			}
			return strings.Contains(strings.ToLower(compactWhitespace(masked)), "xmlparser(schema=")
		}
	}
	return false
}

func pushXMLFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}

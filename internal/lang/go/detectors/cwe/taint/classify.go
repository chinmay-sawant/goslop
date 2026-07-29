package taint

import (
	"strings"
	"unicode"
)

// ClassifySource maps a call function text to a SourceKind.
func ClassifySource(funcText string) (SourceKind, bool) {
	call := funcText
	if strings.Contains(call, ".URL.Query") ||
		strings.Contains(call, ".FormValue") ||
		strings.Contains(call, ".PostForm") ||
		strings.Contains(call, ".Header.Get") ||
		strings.Contains(call, ".GetHeader") ||
		strings.Contains(call, ".GetRawData") ||
		call == "io.ReadAll(r.Body)" ||
		strings.Contains(call, ".PathValue") ||
		strings.Contains(call, ".Param") ||
		call == "c.Query" || call == "c.DefaultQuery" || call == "c.QueryArray" ||
		strings.Contains(call, "chi.URLParam") ||
		call == "c.Params" ||
		(strings.HasSuffix(call, ".Params") && strings.HasPrefix(call, "c.")) {
		return SourceUserInput, true
	}
	if call == "os.Args" || call == "flag.Args" || call == "flag.String" || call == "flag.Int" {
		return SourceArgs, true
	}
	if call == "os.Getenv" || call == "os.LookupEnv" {
		return SourceEnvVar, true
	}
	if call == "io.ReadAll" || strings.Contains(call, ".Scanner.Text") || strings.Contains(call, ".Reader.ReadString") {
		return SourceFile, true
	}
	if strings.Contains(call, ".Conn.Read") || strings.Contains(call, "http.Request.Body") {
		return SourceNetwork, true
	}
	return 0, false
}

// ClassifySink maps a call name to a sink kind and primary argument index.
// source is the full unit source (for html/template heuristics).
// receiver is the method receiver text when available.
// firstArg is the first argument text.
// isResponseWriter is a precomputed check for HTTP write sinks.
func ClassifySink(funcText, source, receiver, firstArg string, hasHTMLTemplate bool) (SinkKind, int, bool) {
	callName := funcText

	if callName == "exec.Command" || callName == "exec.CommandContext" {
		return SinkCommandExec, 0, true
	}

	if strings.HasSuffix(callName, ".Query") ||
		strings.HasSuffix(callName, ".Exec") ||
		strings.HasSuffix(callName, ".QueryRow") ||
		strings.HasSuffix(callName, ".QueryContext") ||
		strings.HasSuffix(callName, ".ExecContext") ||
		strings.HasSuffix(callName, ".QueryRowContext") ||
		strings.HasSuffix(callName, ".Raw") ||
		callName == "sqlx.Get" || callName == "sqlx.Select" || callName == "sqlx.NamedExec" {
		return SinkSQLQuery, 0, true
	}

	if callName == "os.Create" || callName == "os.Open" || callName == "os.OpenFile" ||
		callName == "os.ReadFile" || callName == "os.WriteFile" ||
		callName == "ioutil.ReadFile" || callName == "ioutil.WriteFile" {
		return SinkFileOpen, 0, true
	}

	if (strings.HasSuffix(callName, ".Execute") || strings.HasSuffix(callName, ".ExecuteTemplate")) &&
		!isPlainHTMLTemplateExecute(source, receiver, firstArg, callName) {
		return SinkTemplate, 1, true
	}

	// Method forms of command execution.
	if (strings.HasSuffix(callName, ".Run") || strings.HasSuffix(callName, ".Start") || strings.HasSuffix(callName, ".Output")) &&
		(strings.Contains(receiver, "exec.Command") || strings.HasPrefix(receiver, "exec.")) {
		return SinkCommandExec, 0, true
	}

	if callName == "xml.Unmarshal" || strings.HasSuffix(callName, ".DecodeElement") {
		return SinkXMLQuery, 0, true
	}
	if callName == "json.Unmarshal" || strings.HasSuffix(callName, ".Decode") || strings.Contains(callName, "gob.NewDecoder") {
		return SinkDeserialization, 0, true
	}
	if callName == "ldap.Dial" || callName == "ldap.Search" || callName == "ldap.SearchByAttribute" || callName == "ldap.NewSearchRequest" {
		return SinkLDAPQuery, 0, true
	}

	// template.HTML(...) trusted-content cast is itself a Template sink.
	if isTemplateHTMLCall(callName, hasHTMLTemplate) {
		return SinkTemplate, 0, true
	}

	return 0, 0, false
}

// ClassifySinkHTTPWrite classifies HTTP write sinks when the writer is a ResponseWriter.
func ClassifySinkHTTPWrite(funcText string, isResponseWriter bool, firstArgIsStringSlice bool) (SinkKind, int, bool) {
	if !isResponseWriter {
		return 0, 0, false
	}
	callName := funcText
	if callName == "fmt.Fprintf" {
		return SinkHTTPWrite, 0, true
	}
	if strings.HasSuffix(callName, ".WriteString") || strings.HasSuffix(callName, ".WriteHeader") {
		return SinkHTTPWrite, 0, true
	}
	if strings.HasSuffix(callName, ".Write") && !firstArgIsStringSlice {
		return SinkHTTPWrite, 0, true
	}
	return 0, 0, false
}

// ClassifySanitizer maps a call name to a SanitizerKind.
func ClassifySanitizer(funcText string) (SanitizerKind, bool) {
	call := funcText
	// filepath.Clean / path.Clean alone are NOT path-safe.
	if call == "filepath.Base" {
		return SanitizerPath, true
	}
	if call == "html.EscapeString" ||
		strings.Contains(call, "template.HTMLEscaper") ||
		strings.Contains(call, "template.JSEscaper") {
		return SanitizerHTML, true
	}
	if call == "url.QueryEscape" || call == "url.PathEscape" {
		return SanitizerURL, true
	}
	if strings.HasPrefix(call, "regexp.") && strings.Contains(call, ".MatchString") {
		return SanitizerValidation, true
	}
	if call == "strconv.Atoi" || call == "strconv.ParseInt" || call == "strconv.ParseFloat" || call == "strconv.ParseUint" {
		return SanitizerValidation, true
	}
	if call == "len" {
		return SanitizerBounded, true
	}
	if call == "ldap.EscapeFilter" {
		return SanitizerLDAP, true
	}
	if call == "xml.EscapeText" || call == "xml.Marshal" {
		return SanitizerXML, true
	}
	// Name-based heuristic prefixes.
	name := call
	if i := strings.LastIndex(call, "."); i >= 0 {
		name = call[i+1:]
	}
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "sanitize") ||
		strings.HasPrefix(lower, "escape") ||
		strings.HasPrefix(lower, "validate") ||
		strings.HasPrefix(lower, "purify") {
		return SanitizerValidation, true
	}
	return 0, false
}

// IsSourceOrSanitizerCall reports whether RHS is (or contains) a known source/sanitizer/sink call.
func IsSourceOrSanitizerCall(rhs string) bool {
	call := strings.TrimSpace(strings.SplitN(rhs, "(", 2)[0])
	if call != "" {
		if _, ok := ClassifySource(call); ok {
			return true
		}
		if _, ok := ClassifySanitizer(call); ok {
			return true
		}
		if isSinkCallByName(call) {
			return true
		}
	}
	for _, name := range knownSanitizerCallPrefixes {
		if strings.Contains(rhs, name) {
			return true
		}
	}
	return false
}

func isSinkCallByName(n string) bool {
	return n == "exec.Command" || n == "exec.CommandContext" ||
		strings.HasSuffix(n, ".Query") || strings.HasSuffix(n, ".Exec") || strings.HasSuffix(n, ".QueryRow") ||
		n == "os.Create" || n == "os.Open" || n == "os.OpenFile" || n == "os.ReadFile" || n == "os.WriteFile" ||
		n == "ioutil.ReadFile" || n == "ioutil.WriteFile" ||
		strings.HasSuffix(n, ".Write") || strings.HasSuffix(n, ".Execute") || strings.HasSuffix(n, ".ExecuteTemplate") ||
		n == "fmt.Fprintf" || n == "xml.Unmarshal" || strings.HasSuffix(n, ".DecodeElement") ||
		n == "json.Unmarshal" || strings.HasSuffix(n, ".Decode") || strings.Contains(n, "gob.NewDecoder") ||
		n == "ldap.Dial" || n == "ldap.Search" || n == "ldap.SearchByAttribute" || n == "ldap.NewSearchRequest"
}

var knownSanitizerCallPrefixes = []string{
	"filepath.Base(",
	"html.EscapeString(",
	"ldap.EscapeFilter(",
	"xml.EscapeText(",
	"xml.Marshal(",
}

func isTemplateHTMLCall(callName string, hasHTMLTemplate bool) bool {
	if !hasHTMLTemplate {
		return false
	}
	// template.HTML / tmpl.HTML etc. — last segment is trusted-type ctor
	base := callName
	if i := strings.LastIndex(callName, "."); i >= 0 {
		base = callName[i+1:]
	}
	switch base {
	case "HTML", "HTMLAttr", "JS", "JSStr", "URL", "Srcset", "CSS":
		return true
	}
	return false
}

// isPlainHTMLTemplateExecute: html/template auto-escapes ordinary strings.
func isPlainHTMLTemplateExecute(source, receiver, _ /*dataArg*/, callName string) bool {
	alias := importAlias(source, "html/template", "template")
	if alias == "" {
		return false
	}
	// trusted content in Execute data arg is handled as separate Template sink
	// via template.HTML cast classification; plain Execute on html/template is safe.
	if receiver == "" {
		return false
	}
	// Local decl of receiver from html/template constructors.
	qualified := alias + "."
	for _, line := range strings.Split(source, "\n") {
		line = strings.TrimSpace(line)
		if (strings.HasPrefix(line, receiver+" :=") ||
			strings.HasPrefix(line, receiver+" =") ||
			strings.HasPrefix(line, "var "+receiver)) &&
			strings.Contains(line, qualified) {
			return true
		}
	}
	// Also: method on alias itself is rare; if call is alias.Must(...).Execute, skip.
	_ = callName
	return false
}

func importAlias(source, importPath, defaultAlias string) string {
	for _, line := range strings.Split(source, "\n") {
		line = strings.TrimSpace(line)
		for _, quote := range []byte{'"', '`'} {
			marker := string(quote) + importPath + string(quote)
			if !strings.HasSuffix(line, marker) {
				continue
			}
			prefix := strings.TrimSpace(strings.TrimSuffix(line, marker))
			if prefix == "" || prefix == "import" {
				return defaultAlias
			}
			fields := strings.Fields(prefix)
			if len(fields) > 0 {
				return fields[len(fields)-1]
			}
		}
	}
	return ""
}

// HasHTMLTemplateImport reports whether html/template is imported.
func HasHTMLTemplateImport(source string) bool {
	return importAlias(source, "html/template", "template") != ""
}

// IsIdent reports a simple Go identifier.
func IsIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
			continue
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

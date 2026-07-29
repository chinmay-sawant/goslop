package perf

import (
	"strings"
	"unicode"
)

// IsInLoop reports whether a call sits inside a for_statement.
func IsInLoop(call CallFact) bool {
	return call.EnclosingLoop != nil
}

// IsAssignmentInLoop reports whether an assignment sits inside a for_statement.
func IsAssignmentInLoop(a AssignmentFact) bool {
	return a.EnclosingLoop != nil
}

// IsHandlerShaped reports whether the 1 KiB window before startByte looks like
// a request handler (net/http, Gin, Echo, Fiber, Chi).
func IsHandlerShaped(source string, startByte int) bool {
	if startByte < 0 {
		startByte = 0
	}
	if startByte > len(source) {
		startByte = len(source)
	}
	windowStart := startByte - 1024
	if windowStart < 0 {
		windowStart = 0
	}
	// clamp to rune boundary
	for windowStart > 0 && !utf8Start(source, windowStart) {
		windowStart--
	}
	window := source[windowStart:startByte]
	return strings.Contains(window, "http.ResponseWriter") ||
		strings.Contains(window, "*gin.Context") ||
		strings.Contains(window, "gin.Context") ||
		strings.Contains(window, "echo.Context") ||
		strings.Contains(window, "*fiber.Ctx") ||
		strings.Contains(window, "c *fiber.Ctx") ||
		strings.Contains(window, "chi.URLParam") ||
		strings.Contains(window, "func ServeHTTP") ||
		strings.Contains(window, "c.JSON(") ||
		strings.Contains(window, "c.String(") ||
		strings.Contains(window, "c.HTML(")
}

func utf8Start(s string, i int) bool {
	if i <= 0 || i >= len(s) {
		return true
	}
	return s[i]&0xC0 != 0x80
}

var handlerHotTokens = []string{
	"handler", "middleware", "servehttp", "handlerequest", "handlemessage",
}
var libraryHotTokens = []string{
	"encode", "decode", "marshal", "unmarshal", "serialize", "deserialize",
	"render", "compress", "generate", "sign",
}

// FunctionNameIsHot reports handler/codec-style function names.
func FunctionNameIsHot(name string) bool {
	if name == "" {
		return false
	}
	lower := strings.ToLower(name)
	for _, tok := range handlerHotTokens {
		if strings.Contains(lower, tok) {
			return true
		}
	}
	for _, tok := range libraryHotTokens {
		if strings.Contains(lower, tok) {
			return true
		}
	}
	return strings.HasSuffix(lower, "handler") || strings.HasSuffix(lower, "middleware")
}

// EnclosingFunctionName walks backward to the nearest func declaration name.
func EnclosingFunctionName(source string, startByte int) (string, bool) {
	if startByte > len(source) {
		startByte = len(source)
	}
	if startByte < 0 {
		startByte = 0
	}
	head := source[:startByte]
	funcKw := strings.LastIndex(head, "func ")
	if funcKw < 0 {
		return "", false
	}
	after := strings.TrimLeft(source[funcKw+len("func "):startByte], " \t")
	if strings.HasPrefix(after, "(") {
		close := strings.Index(after, ")")
		if close < 0 {
			return "", false
		}
		after = strings.TrimLeft(after[close+1:], " \t")
	}
	end := 0
	for end < len(after) {
		r := rune(after[end])
		if end == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				break
			}
		} else if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			break
		}
		// after is string of bytes; step carefully for ASCII-only Go idents
		end++
	}
	name := after[:end]
	if name == "" {
		return "", false
	}
	// ensure still inside function body
	rel := strings.Index(source[funcKw:startByte], "{")
	if rel < 0 {
		return "", false
	}
	bodyOpen := funcKw + rel
	depth := 0
	for _, ch := range source[bodyOpen:startByte] {
		switch ch {
		case '{':
			depth++
		case '}':
			depth--
		}
	}
	if depth <= 0 {
		return "", false
	}
	return name, true
}

// IsHotPath is the unified hot-path predicate (loop | handler window | hot name).
func IsHotPath(source string, startByte int, inLoop bool) bool {
	if inLoop {
		return true
	}
	if name, ok := EnclosingFunctionName(source, startByte); ok {
		if name == "init" || name == "main" {
			return false
		}
	} else {
		return false
	}
	if IsHandlerShaped(source, startByte) {
		return true
	}
	name, ok := EnclosingFunctionName(source, startByte)
	return ok && FunctionNameIsHot(name)
}

// IsRequestPath is a whole-file request-handler guard.
func IsRequestPath(index interface{ Has(string) bool }) bool {
	return index.Has("gin.HandlerFunc") ||
		index.Has("echo.HandlerFunc") ||
		index.Has("http.HandlerFunc") ||
		index.Has("http.ResponseWriter") ||
		index.Has("*gin.Context") ||
		index.Has("echo.Context") ||
		index.Has("*fiber.Ctx")
}

// FileDisplayPath returns unit display path or path.
func FileDisplayPath(path, display string) string {
	if display != "" {
		return display
	}
	return path
}

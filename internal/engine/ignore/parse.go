package ignore

import (
	"strings"
)

const fileIgnoreScanLines = 20

// ParseInlineIgnores returns 1-based line → directive for next-line, EOL, and block ranges.
//
// Directives are recognized only inside real line comments (`//` / `#`), never
// inside string literals or block comments (`/* … */`).
func ParseInlineIgnores(source string) map[int]Directive {
	ex := extractLineComments(source)
	ignores := make(map[int]Directive)
	var openBlocks []struct {
		start int
		dir   Directive
	}

	for _, comment := range ex.comments {
		if dir, ok := parseBlockStartBody(comment.body); ok {
			openBlocks = append(openBlocks, struct {
				start int
				dir   Directive
			}{start: comment.line + 1, dir: dir})
			continue
		}
		if isBlockEndBody(comment.body) {
			if n := len(openBlocks); n > 0 {
				blk := openBlocks[n-1]
				openBlocks = openBlocks[:n-1]
				for ln := blk.start; ln < comment.line; ln++ {
					mergeDirective(ignores, ln, blk.dir)
				}
			}
			continue
		}
		if dir, ok := parseIgnoreBody(comment.body); ok {
			if comment.hasCodeBefore {
				mergeDirective(ignores, comment.line, dir)
			} else if target, ok := nextCodeLine(ex.codeLines, comment.line); ok {
				mergeDirective(ignores, target, dir)
			}
		}
	}

	last := strings.Count(source, "\n") + 2
	for _, blk := range openBlocks {
		for ln := blk.start; ln < last; ln++ {
			mergeDirective(ignores, ln, blk.dir)
		}
	}
	return ignores
}

// ParseFileIgnore returns a file-level ignore directive from the top of source, if present.
func ParseFileIgnore(source string) (Directive, bool) {
	ex := extractLineComments(source)
	for _, comment := range ex.comments {
		if comment.line > fileIgnoreScanLines {
			break
		}
		if dir, ok := parseFileIgnoreBody(comment.body); ok {
			return dir, true
		}
	}
	return Directive{}, false
}

func mergeDirective(m map[int]Directive, line int, d Directive) {
	if existing, ok := m[line]; ok {
		existing.Merge(d)
		m[line] = existing
		return
	}
	m[line] = d
}

type lineComment struct {
	line          int
	body          string
	hasCodeBefore bool
}

type extracted struct {
	comments  []lineComment
	codeLines []int
}

func extractLineComments(source string) extracted {
	bytes := []byte(source)
	var comments []lineComment
	var codeLines []int
	line := 1
	lineHasCode := false
	i := 0

	for i < len(bytes) {
		b := bytes[i]

		if b == '\n' {
			if lineHasCode {
				codeLines = append(codeLines, line)
			}
			line++
			lineHasCode = false
			i++
			continue
		}

		// // line comment
		if b == '/' && i+1 < len(bytes) && bytes[i+1] == '/' {
			bodyStart := i + 2
			bodyEnd := findLineEnd(bytes, bodyStart)
			body := strings.TrimSpace(string(bytes[bodyStart:bodyEnd]))
			comments = append(comments, lineComment{
				line:          line,
				body:          body,
				hasCodeBefore: lineHasCode,
			})
			i = bodyEnd
			continue
		}

		// /* block comment */
		if b == '/' && i+1 < len(bytes) && bytes[i+1] == '*' {
			i = skipBlockComment(bytes, i+2, &line, &lineHasCode, &codeLines)
			continue
		}

		// # comment (Python / shell style)
		if b == '#' {
			bodyStart := i + 1
			bodyEnd := findLineEnd(bytes, bodyStart)
			body := strings.TrimSpace(string(bytes[bodyStart:bodyEnd]))
			comments = append(comments, lineComment{
				line:          line,
				body:          body,
				hasCodeBefore: lineHasCode,
			})
			i = bodyEnd
			continue
		}

		// Go raw string
		if b == '`' {
			i = skipUntil(bytes, i+1, '`', &line, &lineHasCode, &codeLines)
			continue
		}

		// Quotes / triple quotes
		if b == '"' || b == '\'' {
			quote := b
			if i+2 < len(bytes) && bytes[i+1] == quote && bytes[i+2] == quote {
				i = skipTriple(bytes, i+3, quote, &line, &lineHasCode, &codeLines)
			} else {
				i = skipQuoted(bytes, i+1, quote, &line, &lineHasCode, &codeLines)
			}
			continue
		}

		if b != ' ' && b != '\t' && b != '\r' && b != 0x0c {
			lineHasCode = true
		}
		i++
	}
	if lineHasCode {
		codeLines = append(codeLines, line)
	}
	return extracted{comments: comments, codeLines: codeLines}
}

func findLineEnd(bytes []byte, i int) int {
	for i < len(bytes) && bytes[i] != '\n' {
		i++
	}
	return i
}

func advanceLine(line *int, lineHasCode *bool, codeLines *[]int) {
	if *lineHasCode {
		*codeLines = append(*codeLines, *line)
		*lineHasCode = false
	}
	*line++
}

func skipBlockComment(bytes []byte, i int, line *int, lineHasCode *bool, codeLines *[]int) int {
	for i < len(bytes) {
		if bytes[i] == '\n' {
			advanceLine(line, lineHasCode, codeLines)
			i++
			continue
		}
		if bytes[i] == '*' && i+1 < len(bytes) && bytes[i+1] == '/' {
			return i + 2
		}
		i++
	}
	return i
}

func skipUntil(bytes []byte, i int, end byte, line *int, lineHasCode *bool, codeLines *[]int) int {
	for i < len(bytes) {
		if bytes[i] == '\n' {
			advanceLine(line, lineHasCode, codeLines)
			i++
			continue
		}
		if bytes[i] == end {
			return i + 1
		}
		i++
	}
	return i
}

func skipQuoted(bytes []byte, i int, quote byte, line *int, lineHasCode *bool, codeLines *[]int) int {
	for i < len(bytes) {
		b := bytes[i]
		if b == '\\' {
			i += 2
			if i > len(bytes) {
				return len(bytes)
			}
			continue
		}
		if b == '\n' {
			advanceLine(line, lineHasCode, codeLines)
			i++
			continue
		}
		if b == quote {
			return i + 1
		}
		i++
	}
	return i
}

func skipTriple(bytes []byte, i int, quote byte, line *int, lineHasCode *bool, codeLines *[]int) int {
	for i < len(bytes) {
		if bytes[i] == '\n' {
			advanceLine(line, lineHasCode, codeLines)
			i++
			continue
		}
		if bytes[i] == quote && i+2 < len(bytes) && bytes[i+1] == quote && bytes[i+2] == quote {
			return i + 3
		}
		i++
	}
	return i
}

func parseIgnoreBody(body string) (Directive, bool) {
	const prefix = "codehound-ignore:"
	if !strings.HasPrefix(body, prefix) {
		return Directive{}, false
	}
	raw := strings.TrimSpace(body[len(prefix):])
	return parseRuleList(raw)
}

func parseBlockStartBody(body string) (Directive, bool) {
	const prefix = "codehound-ignore-start"
	if !strings.HasPrefix(body, prefix) {
		return Directive{}, false
	}
	raw := strings.TrimSpace(body[len(prefix):])
	if raw == "" {
		return All(), true
	}
	if !strings.HasPrefix(raw, ":") {
		return Directive{}, false
	}
	raw = strings.TrimSpace(raw[1:])
	if raw == "" {
		return All(), true
	}
	return parseRuleList(raw)
}

func isBlockEndBody(body string) bool {
	return strings.HasPrefix(body, "codehound-ignore-end")
}

func parseFileIgnoreBody(body string) (Directive, bool) {
	const prefix = "codehound-ignore-file"
	if !strings.HasPrefix(body, prefix) {
		return Directive{}, false
	}
	raw := strings.TrimSpace(body[len(prefix):])
	if raw == "" {
		return All(), true
	}
	if !strings.HasPrefix(raw, ":") {
		return Directive{}, false
	}
	raw = strings.TrimSpace(raw[1:])
	if raw == "" {
		return All(), true
	}
	return parseRuleList(raw)
}

func parseRuleList(raw string) (Directive, bool) {
	if strings.EqualFold(raw, "all") {
		return All(), true
	}
	parts := strings.Split(raw, ",")
	var ids []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			ids = append(ids, p)
		}
	}
	if len(ids) == 0 {
		return Directive{}, false
	}
	return Rules(ids...), true
}

func nextCodeLine(codeLines []int, after int) (int, bool) {
	for _, ln := range codeLines {
		if ln > after {
			return ln, true
		}
	}
	return 0, false
}

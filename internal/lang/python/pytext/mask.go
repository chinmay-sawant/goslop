// Package pytext holds shared pure-Go Python source helpers used by detectors.
package pytext

// Mask keeps byte offsets stable while blanking comments and string literals.
// Source-pattern rules can use it to avoid interpreting examples in
// docstrings, comments, and quoted data as executable Python.
func Mask(source string) string {
	masked := []byte(source)
	inString := byte(0)
	triple := false
	escaped := false
	inComment := false
	for i := 0; i < len(masked); i++ {
		c := masked[i]
		if inComment {
			if c == '\n' {
				inComment = false
			} else {
				masked[i] = ' '
			}
			continue
		}
		if inString != 0 {
			i, inString, triple, escaped = maskPythonString(masked, i, inString, triple, escaped)
			continue
		}
		switch c {
		case '#':
			masked[i] = ' '
			inComment = true
		case '\'', '"':
			inString = c
			if i+2 < len(masked) && masked[i+1] == c && masked[i+2] == c {
				masked[i], masked[i+1], masked[i+2] = ' ', ' ', ' '
				i += 2
				triple = true
			} else {
				masked[i] = ' '
			}
		}
	}
	return string(masked)
}

func maskPythonString(masked []byte, index int, quote byte, triple, escaped bool) (int, byte, bool, bool) {
	current := masked[index]
	// Keep line separators intact so source-mask consumers can safely correlate
	// masked and original lines, including multiline string literals.
	if current != '\n' && current != '\r' {
		masked[index] = ' '
	}
	if triple {
		if current == quote && index+2 < len(masked) && masked[index+1] == quote && masked[index+2] == quote {
			masked[index+1], masked[index+2] = ' ', ' '
			return index + 2, 0, false, false
		}
		return index, quote, true, false
	}
	if escaped {
		return index, quote, false, false
	}
	if current == '\\' {
		return index, quote, false, true
	}
	if current == quote {
		return index, 0, false, false
	}
	return index, quote, false, false
}

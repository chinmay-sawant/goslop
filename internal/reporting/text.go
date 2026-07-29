package reporting

import (
	"fmt"
	"io"

	"github.com/chinmay/codehound/internal/rules"
)

// TextReporter emits one finding per line:
//
//	rule file:line:col message
//
// NoColor is reserved for future ANSI styling; text currently has no color.
type TextReporter struct {
	NoColor bool
}

// Write implements Reporter.
func (r TextReporter) Write(findings []rules.Finding, w io.Writer) error {
	// Honor NO_COLOR even though the skeleton emits no ANSI sequences yet.
	_ = r.NoColor

	for i := range findings {
		f := &findings[i]
		line := f.Line
		col := f.Column
		if line <= 0 {
			line = 1
		}
		if col <= 0 {
			col = 1
		}
		if _, err := fmt.Fprintf(w, "%s %s:%d:%d %s\n", f.RuleID, f.File, line, col, f.Message); err != nil {
			return err
		}
	}
	return nil
}

// Package reporting formats analysis findings for text, JSON, and SARIF.
package reporting

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chinmay/codehound/internal/rules"
)

// Reporter writes findings to an io.Writer.
type Reporter interface {
	Write(findings []rules.Finding, w io.Writer) error
}

// New returns a Reporter for format (text|json|sarif).
func New(format string) (Reporter, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "text", "":
		return TextReporter{NoColor: noColorFromEnv()}, nil
	case "json":
		return JSONReporter{Version: DefaultVersion}, nil
	case "sarif":
		return SARIFReporter{Version: DefaultVersion}, nil
	default:
		return nil, fmt.Errorf("unknown output format %q", format)
	}
}

// DefaultVersion is embedded in JSON/SARIF envelopes (overridden by app when needed).
var DefaultVersion = "0.1.0-dev"

func noColorFromEnv() bool {
	// NO_COLOR: any non-empty value disables ANSI (https://no-color.org/).
	return os.Getenv("NO_COLOR") != ""
}

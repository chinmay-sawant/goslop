// Command goslop is the goslop static analysis tool (SAT) CLI.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/chinmay-sawant/goslop/internal/app"
)

func main() {
	// This is the last-resort process boundary. Expected failures should be
	// returned as errors; this guard only prevents an unexpected main-goroutine
	// panic from leaking a Go runtime stack to CLI users.
	defer func() {
		if recovered := recover(); recovered != nil {
			_, _ = fmt.Fprintf(os.Stderr, "goslop: internal error: %v\n", recovered)
			os.Exit(app.ExitInternal)
		}
	}()

	if err := app.Run(os.Args[1:]); err != nil {
		var ece *app.ExitCodeError
		if errors.As(err, &ece) {
			if ece.Err != nil {
				fmt.Fprintln(os.Stderr, ece.Err)
			}
			// ExitFailing with nil Err: findings already on stdout; silent exit 1.
			code := ece.Code
			if code == 0 {
				code = app.ExitConfig
			}
			os.Exit(code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(app.ExitConfig)
	}
}

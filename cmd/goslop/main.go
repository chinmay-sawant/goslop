// Command goslop is the goslop static analysis tool (SAT) CLI.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/chinmay-sawant/goslop/internal/app"
)

func main() {
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

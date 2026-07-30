package app

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//go:embed init_template.toml
var initTemplate []byte

// runInit writes templates/goslop.toml content to ./goslop.toml when absent.
func runInit(stdout io.Writer) error {
	path := "goslop.toml"
	if st, err := os.Stat(path); err == nil && !st.IsDir() {
		return &ExitCodeError{
			Code: ExitConfig,
			Err:  fmt.Errorf("%s already exists in this directory", path),
		}
	} else if err != nil && !os.IsNotExist(err) {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}

	const fileMode = 0o644
	if err := os.WriteFile(path, initTemplate, fileMode); err != nil {
		return &ExitCodeError{
			Code: ExitInternal,
			Err:  fmt.Errorf("failed to write %s: %w", path, err),
		}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	_, _ = fmt.Fprintf(stdout, "wrote starter goslop.toml to %s\n", abs)
	return nil
}

package app

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed init_template.toml
var initTemplate []byte

// runInit writes templates/codehound.toml content to ./codehound.toml when absent.
func runInit() error {
	path := "codehound.toml"
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
	fmt.Printf("wrote starter codehound.toml to %s\n", abs)
	return nil
}

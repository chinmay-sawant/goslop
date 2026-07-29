// Package main demonstrates a dependency reused across multiple files.
package main

import errors "github.com/pkg/errors"

func main() {
	if err := errors.New("boom"); err != nil {
		return
	}
	if err := BuildErr(); err != nil {
		return
	}
}

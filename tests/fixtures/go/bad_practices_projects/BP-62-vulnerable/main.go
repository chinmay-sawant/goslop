// Package main demonstrates a dependency used in only one file.
package main

import errors "github.com/pkg/errors"

func main() {
	if err := errors.New("boom"); err != nil {
		return
	}
}

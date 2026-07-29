// Package main imports a fully pinned module version.
package main

import errors "github.com/pkg/errors"

func main() {
	if err := errors.New("boom"); err != nil {
		panic(err)
	}
}

// Package main demonstrates a second import site for the same dependency.
package main

import errors "github.com/pkg/errors"

// BuildErr returns a second package-scoped error value.
func BuildErr() error {
	return errors.New("again")
}

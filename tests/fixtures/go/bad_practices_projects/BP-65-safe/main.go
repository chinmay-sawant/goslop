// Package main includes a go.sum alongside go.mod.
package main

import errors "github.com/pkg/errors"

func main() {
	if err := errors.New("boom"); err != nil {
		panic(err)
	}
}

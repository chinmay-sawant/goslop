// Package main imports the direct dependency so it is not stale.
package main

import errors "github.com/pkg/errors"

func main() {
	if err := errors.New("boom"); err != nil {
		panic(err)
	}
}

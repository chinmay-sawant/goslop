// Package main avoids local replace directives.
package main

import errors "github.com/pkg/errors"

func main() {
	if err := errors.New("boom"); err != nil {
		panic(err)
	}
}

// Package main demonstrates a dependency outside the curated advisory snapshot.
package main

import errors "github.com/pkg/errors"

func main() {
	if err := errors.New("ok"); err != nil {
		return
	}
}

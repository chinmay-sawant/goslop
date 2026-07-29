// Package main uses the dependency outside tests as well.
package main

import "github.com/stretchr/testify/assert"

func main() {
	if assert.ObjectsAreEqual(1, 1) {
		return
	}
}

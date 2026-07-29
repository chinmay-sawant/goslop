// Package main starts a Gin server without signal handling.
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	if err := r.Run(":8080"); err != nil {
		return
	}
}

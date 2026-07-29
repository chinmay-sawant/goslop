// Package main exposes a public Echo endpoint without rate limiting.
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/status", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	if err := e.Start(":8080"); err != nil {
		return
	}
}

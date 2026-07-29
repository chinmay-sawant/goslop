// Package main starts a Fiber server without a graceful shutdown path.
package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	if err := app.Listen(":8080"); err != nil {
		return
	}
}

package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	handler(app)

	if err := app.Listen(":8000"); err != nil {
		return
	}
}

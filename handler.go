package main

import "github.com/gofiber/fiber/v2"

func handler(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {

		return c.SendString("Hello, World!")
	})

	app.Get("/:url", resolveURL)
	app.Post("/api/v1/shorten", shortenURL)
}

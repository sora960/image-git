package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Image-Git Engine v1.0",
	})

	app.Use(cors.New())

	// Basic Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Engine is online 🟢")
	})

	api := app.Group("/api/v1")

	// Grouped Route Registries
	registerRepoRoutes(api)
	registerRenderRoutes(api)
	registerVersionRoutes(api)

	log.Fatal(app.Listen(":3000"))
}
package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sora960/image-git/internal/gitlogic"
)

func main() {
	// Initialize Fiber with high-performance settings
	app := fiber.New(fiber.Config{
		AppName: "Image-Git Engine v1.0",
	})

	// 1. Middleware: Allow Frontend access
	app.Use(cors.New())

	// 2. Health Check (Basic heartbeat)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Engine is online 🟢")
	})

	// 3. API Versioning Group
	api := app.Group("/api/v1")

	// [GET] Retrieve Repository Status (Manifest)
	api.Get("/repo/:name", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		manifest, err := gitlogic.LoadManifest(repo)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Repository not found"})
		}
		return c.JSON(manifest)
	})

	// [DELETE] Remove a layer by name
	api.Delete("/repo/:name/layers/:layerName", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		layer := c.Params("layerName")
		if err := gitlogic.RemoveLayer(repo, layer); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	})

	// [POST] Remote Render Trigger
	api.Post("/repo/:name/render/:frame", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		frame, _ := c.ParamsInt("frame")
		if err := gitlogic.CompositeFrame(repo, frame); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Render failed"})
		}
		return c.JSON(fiber.Map{"status": "success", "frame": frame})
	})

	// Start server on Port 3000
	log.Fatal(app.Listen(":3000"))
}
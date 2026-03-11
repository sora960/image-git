package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sora960/image-git/internal/gitlogic"
)

func registerRenderRoutes(api fiber.Router) {
	// Execute Composite
	api.Post("/repo/:name/render", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		err := gitlogic.RenderRepository(repo) 
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Rendering failed: " + err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Render successful"})
	})

	// Serve Preview Image
	api.Get("/preview.png", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "image/png")
		c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		return c.SendFile("./data/repositories/art-project/preview.png")
	})
}
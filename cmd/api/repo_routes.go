package main

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/sora960/image-git/internal/gitlogic"
)

func registerRepoRoutes(api fiber.Router) {
	// Get Manifest
	api.Get("/repo/:name", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		manifest, err := gitlogic.LoadManifest(repo)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Repository not found"})
		}
		return c.JSON(manifest)
	})

	// Add Layer
	api.Post("/repo/:name/layers", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		layerName := c.FormValue("name", "web-layer")
		zIndex, _ := strconv.Atoi(c.FormValue("z", "0"))

		fileHeader, err := c.FormFile("image")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "No image provided"})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to open stream"})
		}
		defer file.Close()

		hash, err := gitlogic.StoreStream(repo, file)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		m, _ := gitlogic.LoadManifest(repo)
		m.Layers = append(m.Layers, gitlogic.Layer{
			Name:    layerName,
			Hash:    hash,
			Opacity: 1.0,
			ZIndex:  zIndex,
		})
		
		gitlogic.SaveManifest(repo, m)
		return c.Status(201).JSON(fiber.Map{"status": "success", "hash": hash})
	})

	// Update Layer (Opacity/Z-Index)
	api.Patch("/repo/:name/layers/:layerName", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		layerName := c.Params("layerName")
		type UpdateRequest struct {
			Opacity *float64 `json:"opacity"`
			ZIndex  *int     `json:"z_index"`
		}
		var req UpdateRequest
		c.BodyParser(&req)

		manifest, _ := gitlogic.LoadManifest(repo)
		for i := range manifest.Layers {
			if manifest.Layers[i].Name == layerName {
				if req.Opacity != nil { manifest.Layers[i].Opacity = *req.Opacity }
				if req.ZIndex != nil { manifest.Layers[i].ZIndex = *req.ZIndex }
				break
			}
		}
		gitlogic.SaveManifest(repo, manifest)
		return c.SendStatus(204)
	})

	// Delete Layer
	api.Delete("/repo/:name/layers/:layerName", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		layer := c.Params("layerName")
		if err := gitlogic.RemoveLayer(repo, layer); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	})

	api.Static("/objects", "./data/repositories/art-project/objects")
}
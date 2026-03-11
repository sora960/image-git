package main

import (
	"log"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sora960/image-git/internal/gitlogic"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Image-Git Engine v1.0",
	})

	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Engine is online 🟢")
	})

	api := app.Group("/api/v1")

	// [GET] Retrieve Repository Status
	api.Get("/repo/:name", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		manifest, err := gitlogic.LoadManifest(repo)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Repository not found"})
		}
		return c.JSON(manifest)
	})

	// [PATCH] Update Layer Opacity - Issue #12
api.Patch("/repo/:name/layers/:layerName", func(c *fiber.Ctx) error {
    repo := c.Params("name")
    layerName := c.Params("layerName")

    // Define a minimalist struct for the incoming request
    type UpdateRequest struct {
        Opacity float64 `json:"opacity"`
    }

    var req UpdateRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // 1. Load the current manifest
    manifest, err := gitlogic.LoadManifest(repo)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Repository not found"})
    }

    // 2. Find and update the specific layer
    updated := false
    for i := range manifest.Layers {
        if manifest.Layers[i].Name == layerName {
            manifest.Layers[i].Opacity = req.Opacity
            updated = true
            break
        }
    }

    if !updated {
        return c.Status(404).JSON(fiber.Map{"error": "Layer not found"})
    }

    // 3. Save the modified manifest back to disk
    if err := gitlogic.SaveManifest(repo, manifest); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to save manifest"})
    }

    return c.SendStatus(204) // Success (No Content)
})

// Inside cmd/api/main.go, update the POST handler:
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
            // CHANGE THIS LINE to see the real error:
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        m, _ := gitlogic.LoadManifest(repo)
        m.Layers = append(m.Layers, gitlogic.Layer{
            Name:       layerName,
            Hash:       hash,
            Opacity:    1.0,
            StartFrame: 0,
            EndFrame:   999,
            ZIndex:     zIndex,
        })
        
        if err := gitlogic.SaveManifest(repo, m); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Manifest update failed: " + err.Error()})
        }

        return c.Status(201).JSON(fiber.Map{
            "status": "success",
            "hash":   hash,
        })
    })

	// [DELETE] Remove a layer
	api.Delete("/repo/:name/layers/:layerName", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		layer := c.Params("layerName")
		if err := gitlogic.RemoveLayer(repo, layer); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	})

	// [POST] Remote Render
	api.Post("/repo/:name/render/:frame", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		frame, _ := c.ParamsInt("frame")
		if err := gitlogic.CompositeFrame(repo, frame); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Render failed"})
		}
		return c.JSON(fiber.Map{"status": "success", "frame": frame})
	})

	log.Fatal(app.Listen(":3000"))
}
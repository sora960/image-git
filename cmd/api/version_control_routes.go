package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sora960/image-git/internal/gitlogic"
)

func registerVersionRoutes(api fiber.Router) {
	// [POST] Create Snapshot (Commit) - Issue #21
	api.Post("/repo/:name/commit", func(c *fiber.Ctx) error {
		repo := c.Params("name")
		type CommitRequest struct {
			Message string `json:"message"`
		}
		var req CommitRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}
		if req.Message == "" {
			req.Message = "Manual Snapshot"
		}

		hash, err := gitlogic.SaveCommit(repo, req.Message)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(fiber.Map{
			"status": "success",
			"hash":   hash,
		})
	})
}
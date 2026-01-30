package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadMedia(c *fiber.Ctx) error {
	file, err := c.FormFile("media")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}
	uniqueName := uuid.New().String() + "_" + file.Filename
	path := "./uploads/" + uniqueName

	if err := c.SaveFile(file, path); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to upload"})
	}
	return c.JSON(fiber.Map{
		"url": "/uploads/" + uniqueName,
	})
}

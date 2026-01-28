package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	role := c.Locals("role")
	name := c.Locals("name")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "could not retrieve user data from token",
		})
	}
	return c.JSON(fiber.Map{
		"user_id": userID,
		"role":    role,
		"name":    name,
		"status":  "authenticated",
	})
}

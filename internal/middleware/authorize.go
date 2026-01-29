package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
)

func AuthorizeRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": fiber.ErrForbidden.Message})
		}

		if slices.Contains(allowedRoles, userRole) {
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You do not have permission to perform this action"})
	}
}

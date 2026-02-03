package middleware

import (
	"go_blog/internal/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Validate(v *validator.Validate, body interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := body
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		validator := helper.Validator{Validate: v}
		errs := validator.ValidateStruct(body)
		if len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errs,
			})
		}
		return c.Next()
	}
}

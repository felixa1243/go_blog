package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupLogger(app *fiber.App) (fiber.Handler, *os.File) {
	file, _ := os.OpenFile("app_audit.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil || c.Response().StatusCode() >= 400 {
			// Write to file here
			fmt.Fprintf(file, "[%s] %d %s %s | Error: %v\n",
				time.Now().Format(time.RFC3339),
				c.Response().StatusCode(),
				c.Method(),
				c.Path(),
				err,
			)
		}
		return err
	}, file
}

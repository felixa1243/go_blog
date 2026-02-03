package server

import (
	"go_blog/internal/database"
	"go_blog/internal/middleware"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type FiberServer struct {
	*fiber.App
	validator *validator.Validate
	db        database.Service
	LogFile   *os.File
}

func New() *FiberServer {
	validator := validator.New()
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go_blog",
			AppName:      "go_blog",
		}),
		validator: validator,
		db:        database.New(),
	}
	fileLogger, file := middleware.SetupLogger(server.App)
	server.LogFile = file
	server.App.Use(fileLogger)

	server.App.Static("/uploads", "./uploads")
	server.App.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	return server
}

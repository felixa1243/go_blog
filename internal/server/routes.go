package server

import (
	"go_blog/internal/controllers"
	"go_blog/internal/helper"
	"go_blog/internal/midleware"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)
	//Blogs
	pubKey, _ := helper.LoadPublicKey(os.Getenv("RSA_PUBLIC_KEY_PATH"))
	auth := midleware.NewAuthMiddleware(pubKey)
	s.App.Use("/posts", auth)
	s.App.Get("/posts", controllers.GetPosts)
	s.App.Post("/posts", controllers.CreatePosts)
	s.App.Get("/health", s.healthHandler)
	s.App.Get("/me", auth, controllers.GetMe)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

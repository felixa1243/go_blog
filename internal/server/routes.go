package server

import (
	"go_blog/internal/controllers"
	"go_blog/internal/midleware"

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
	//Auth
	s.App.Post("/auth/register", controllers.Register)
	s.App.Post("/auth/login", controllers.Login)
	//Blogs
	s.App.Get("/blogs", midleware.JWTProtected, controllers.GetBlogs)
	s.App.Post("/blogs", midleware.JWTProtected, controllers.CreateBlog)
	s.App.Get("/health", s.healthHandler)

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

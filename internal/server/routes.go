package server

import (
	"go_blog/internal/controllers"
	"go_blog/internal/helper"
	"go_blog/internal/middleware"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/samber/do/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,Cookie",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	// DI Setup
	publicKey, _ := helper.LoadPublicKey(os.Getenv("RSA_PUBLIC_KEY_PATH"))
	injector := SetupDI(s.db.GetDB(), publicKey, s.validator)

	// Invoke via Interfaces
	auth := do.MustInvoke[fiber.Handler](injector)
	categoriesController := do.MustInvoke[controllers.ICategoryController](injector)
	postController := do.MustInvoke[controllers.IPostController](injector)

	// --- Posts ---
	s.App.Get("/posts", postController.GetPosts)
	s.App.Post("/posts", auth, middleware.AuthorizeRole("Administrator", "Blog:Editor"), postController.CreatePosts)
	s.App.Put("/posts/:id", auth, middleware.AuthorizeRole("Administrator", "Blog:Editor"), postController.EditPost)
	s.App.Delete("/posts/:id", auth, middleware.AuthorizeRole("Administrator"), postController.DeletePost)

	// --- Categories ---
	s.App.Get("/categories", categoriesController.GetCategories)
	s.App.Post("/categories", auth, middleware.AuthorizeRole("Administrator"), categoriesController.CreateCategory)

	// --- System ---
	s.App.Post("/media", auth, middleware.AuthorizeRole("Administrator", "Blog:Editor"), controllers.UploadMedia)
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

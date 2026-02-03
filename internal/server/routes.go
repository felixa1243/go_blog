package server

import (
	"go_blog/internal/controllers"
	"go_blog/internal/database"
	"go_blog/internal/helper"
	"go_blog/internal/middleware"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,Cookie",
		AllowCredentials: true, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)
	pubKey, _ := helper.LoadPublicKey(os.Getenv("RSA_PUBLIC_KEY_PATH"))
	auth := middleware.NewAuthMiddleware(pubKey)
	db := database.New().GetDB()
	categoriesController := controllers.NewCategoryController(db, s.validator)
	//posts
	s.App.Get("/posts", controllers.GetPosts)
	s.App.Put("/posts/:id", auth, middleware.AuthorizeRole("Administrator", "Blog:Editor"), controllers.EditPost)
	s.App.Post("/posts", auth, middleware.AuthorizeRole("Administrator", "Blog:Editor"), controllers.CreatePosts)
	//media
	s.App.Post("/media", auth, middleware.AuthorizeRole("Administrator", "Blog:Editor"), controllers.UploadMedia)
	//categories
	s.App.Post("/categories", auth, middleware.AuthorizeRole("Administrator"), categoriesController.CreateCategory)
	s.App.Get("/categories", categoriesController.GetCategories)
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

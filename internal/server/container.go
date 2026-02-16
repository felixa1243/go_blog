package server

import (
	"crypto/rsa"
	"go_blog/internal/controllers"
	"go_blog/internal/middleware"
	"go_blog/internal/repositories"
	"go_blog/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func SetupDI(db *gorm.DB, publicKey *rsa.PublicKey, v *validator.Validate) do.Injector {
	injector := do.New()
	do.ProvideValue(injector, db)
	do.ProvideValue(injector, publicKey)
	do.ProvideValue(injector, v)
	do.Provide(injector, func(i do.Injector) (fiber.Handler, error) {
		key := do.MustInvoke[*rsa.PublicKey](i)
		return middleware.NewAuthMiddleware(key), nil
	})

	// Repositories
	do.Provide(injector, func(i do.Injector) (repositories.CategoryRepository, error) {
		return repositories.NewCategoryRepository(do.MustInvoke[*gorm.DB](i)), nil
	})
	do.Provide(injector, func(i do.Injector) (repositories.PostRepository, error) {
		return repositories.NewPostRepository(do.MustInvoke[*gorm.DB](i)), nil
	})

	// Services
	do.Provide(injector, func(i do.Injector) (services.CategoryService, error) {
		return services.NewCategoryService(
			do.MustInvoke[repositories.CategoryRepository](i),
			do.MustInvoke[*validator.Validate](i),
		), nil
	})
	do.Provide(injector, func(i do.Injector) (services.PostService, error) {
		return services.NewPostService(
			do.MustInvoke[repositories.PostRepository](i),
			do.MustInvoke[repositories.CategoryRepository](i),
		), nil
	})

	// Controllers
	do.Provide(injector, func(i do.Injector) (controllers.ICategoryController, error) {
		return controllers.NewCategoryController(
			do.MustInvoke[services.CategoryService](i),
			do.MustInvoke[*validator.Validate](i),
		), nil
	})
	do.Provide(injector, func(i do.Injector) (controllers.IPostController, error) {
		return controllers.NewPostController(
			do.MustInvoke[services.PostService](i),
			do.MustInvoke[*validator.Validate](i),
		), nil
	})

	return injector
}

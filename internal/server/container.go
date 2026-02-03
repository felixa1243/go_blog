package server

import (
	"crypto/rsa"
	"go_blog/internal/controllers"
	"go_blog/internal/middleware"

	"github.com/go-playground/validator/v10" // Ensure v10 is used here
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
	do.Provide(injector, func(i do.Injector) (controllers.ICategoryController, error) {
		return &controllers.CategoryController{
			DB:        do.MustInvoke[*gorm.DB](i),
			Validator: do.MustInvoke[*validator.Validate](i),
		}, nil
	})
	do.Provide(injector, func(i do.Injector) (controllers.IPostController, error) {
		return &controllers.PostController{
			Db:       do.MustInvoke[*gorm.DB](i),
			Validate: do.MustInvoke[*validator.Validate](i),
		}, nil
	})

	return injector
}

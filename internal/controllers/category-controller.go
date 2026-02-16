package controllers

import (
	"go_blog/internal/dto"
	"go_blog/internal/helper"
	"go_blog/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ICategoryController interface {
	GetCategories(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	EditCategory(id int) error
	DeleteCategory(id int) error
}

type CategoryController struct {
	Service   services.CategoryService
	Validator *validator.Validate
}

func NewCategoryController(service services.CategoryService, validator *validator.Validate) ICategoryController {
	return &CategoryController{Service: service, Validator: validator}
}

func (c *CategoryController) GetCategories(ctx *fiber.Ctx) error {
	datas, err := c.Service.GetCategories()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch categories",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": datas,
	})
}
func (c *CategoryController) CreateCategory(ctx *fiber.Ctx) error {
	var req dto.CategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	requestValidator := helper.Validator{Validate: c.Validator}
	if errs := requestValidator.ValidateStruct(&req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": errs})
	}

	category, err := c.Service.CreateCategory(&req)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Could not create category"})
	}
	return ctx.JSON(fiber.Map{"message": "Category created successfully", "category": category})
}

func (c *CategoryController) EditCategory(id int) error {
	return nil
}
func (c *CategoryController) DeleteCategory(id int) error {
	return nil
}

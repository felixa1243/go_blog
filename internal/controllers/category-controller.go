package controllers

import (
	"go_blog/internal/dto"
	"go_blog/internal/helper"
	"go_blog/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ICategoryController interface {
	GetCategories(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	EditCategory(id int) error
	DeleteCategory(id int) error
}

type CategoryController struct {
	db        *gorm.DB
	validator *validator.Validate
}

func (c *CategoryController) GetCategories(ctx *fiber.Ctx) error {
	var categories []models.Category
	if err := c.db.Find(&categories).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch categories",
		})
	}

	datas := make([]dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		datas = append(datas, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
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
	requestValidator := helper.Validator{Validate: c.validator}
	if errs := requestValidator.ValidateStruct(&req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": errs})
	}
	category := models.Category{
		Name: req.Name,
	}
	if err := c.db.Create(&category).Error; err != nil {
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
func NewCategoryController(db *gorm.DB, validator *validator.Validate) ICategoryController {
	return &CategoryController{db: db, validator: validator}
}

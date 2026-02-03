package controllers

import (
	"encoding/json"
	"fmt"
	"go_blog/internal/database"
	"go_blog/internal/dto"
	"go_blog/internal/models"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPostController interface {
	GetPosts(c *fiber.Ctx) error
	CreatePosts(c *fiber.Ctx) error
	EditPost(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
}
type PostController struct {
	Db       *gorm.DB
	Validate *validator.Validate
}

func NewPostController() IPostController {
	return &PostController{
		Db:       database.New().GetDB(),
		Validate: validator.New(),
	}
}

func (pc *PostController) GetPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	fmt.Println(limit)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	var modelPosts []models.Post
	var total int64

	//  Count total records
	pc.Db.Model(&models.Post{}).Count(&total)

	// Fetch paginated data
	if err := pc.Db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&modelPosts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch posts"})
	}

	//  Transform to DTOs
	var p []dto.PostResponse
	for _, post := range modelPosts {
		thumbnailpath := post.ThumbnailPath
		if strings.Contains(post.ThumbnailPath, "./") {
			parts := strings.Split(post.ThumbnailPath, "./")
			if len(parts) > 1 {
				thumbnailpath = parts[1]
			}
		}

		p = append(p, dto.PostResponse{
			ID:            post.ID,
			Title:         post.Title,
			Content:       post.Content,
			CreatedAt:     post.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     post.UpdatedAt.Format(time.RFC3339),
			ThumbnailPath: thumbnailpath,
			Slug:          post.Slug,
		})
	}
	return c.JSON(fiber.Map{
		"data": p,
		"meta": fiber.Map{
			"total_records": total,
			"current_page":  page,
			"total_pages":   math.Ceil(float64(total) / float64(limit)),
			"limit":         limit,
		},
	})
}

func (pc *PostController) CreatePosts(c *fiber.Ctx) error {
	authorIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	authorUUID, err := uuid.Parse(authorIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Author ID format"})
	}

	file, err := c.FormFile("thumbnail")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Thumbnail is required"})
	}

	filePath := fmt.Sprintf("./uploads/%d_%s", time.Now().Unix(), file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image"})
	}

	categoriesRaw := c.FormValue("categories")
	var categoryIDs []int
	if err := json.Unmarshal([]byte(categoriesRaw), &categoryIDs); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid categories format"})
	}

	var categoriesFound []models.Category
	if len(categoryIDs) > 0 {
		if err := pc.Db.Where("id IN ?", categoryIDs).Find(&categoriesFound).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Database error fetching categories"})
		}
	}

	post := models.Post{
		Title:         c.FormValue("title"),
		Content:       c.FormValue("content"),
		Slug:          c.FormValue("slug"),
		AuthorID:      authorUUID,
		ThumbnailPath: filePath,
		Categories:    categoriesFound,
	}

	if err := pc.Db.Create(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create post"})
	}

	return c.Status(201).JSON(post)
}

func (pc *PostController) EditPost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post
	if err := pc.Db.First(&post, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	userID := c.Locals("user_id").(uuid.UUID)
	role := c.Locals("role").(string)

	if post.AuthorID != userID && role != "Administrator" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	var req dto.PostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := pc.Db.Model(&post).Updates(models.Post{
		Title:   req.Title,
		Content: req.Content,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post"})
	}

	return c.JSON(fiber.Map{"message": "Post updated successfully", "post": post})
}

func (pc *PostController) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := pc.Db.Delete(&models.Post{}, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete post"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

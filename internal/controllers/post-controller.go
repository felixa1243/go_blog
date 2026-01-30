package controllers

import (
	"fmt"
	"go_blog/internal/database"
	"go_blog/internal/dto"
	"go_blog/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetPosts(c *fiber.Ctx) error {
	db := database.New().GetDB()

	var modelPosts []models.Post
	if err := db.Find(&modelPosts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch posts"})
	}

	// 3. Map Models to DTOs
	var p []dto.PostResponse
	for _, post := range modelPosts {
		p = append(p, dto.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt.String(),
			UpdatedAt: post.UpdatedAt.String(),
		})
	}
	if len(p) == 0 {
		return c.JSON(fiber.Map{"posts": []string{}})
	}
	return c.JSON(fiber.Map{"posts": p})
}

func CreatePosts(c *fiber.Ctx) error {
	db := database.New().GetDB()

	title := c.FormValue("title")
	content := c.FormValue("content")
	authorIDStr := c.Locals("user_id").(string)
	authorUUID, _ := uuid.Parse(authorIDStr)
	file, err := c.FormFile("thumbnail")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Thumbnail is required"})
	}
	filePath := fmt.Sprintf("./uploads/%d_%s", time.Now().Unix(), file.Filename)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save image"})
	}
	post := models.Post{
		Title:         title,
		Content:       content,
		AuthorID:      authorUUID,
		ThumbnailPath: filePath,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := db.Create(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create post"})
	}

	return c.Status(201).JSON(post)
}
func EditPost(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.New().GetDB()

	var post models.Post
	if err := db.First(&post, "id = ?", id).Error; err != nil {
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
	if err := db.Model(&post).Updates(models.Post{
		Title:   req.Title,
		Content: req.Content,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post"})
	}

	return c.JSON(fiber.Map{"message": "Post updated successfully", "post": post})
}

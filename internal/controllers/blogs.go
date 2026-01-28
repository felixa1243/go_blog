package controllers

import (
	"go_blog/internal/database"
	"go_blog/internal/dto"
	"go_blog/internal/models"
	"go_blog/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetBlogs(c *fiber.Ctx) error {
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

func CreateBlog(c *fiber.Ctx) error {
	db := database.New().GetDB()
	var req dto.PostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	author, ok := utils.GetUserID(c)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "User unauthorized or ID missing"})
	}
	authorUUID, errParse := uuid.Parse(author)
	if errParse != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create post"})
	}
	post := models.Post{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  authorUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create post"})
	}
	return c.JSON(fiber.Map{"post": post})
}

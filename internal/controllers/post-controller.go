package controllers

import (
	"encoding/json"
	"fmt"
	"go_blog/internal/dto"
	"go_blog/internal/services"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IPostController interface {
	GetPosts(c *fiber.Ctx) error
	CreatePosts(c *fiber.Ctx) error
	EditPost(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
}
type PostController struct {
	Service  services.PostService
	Validate *validator.Validate
}

func NewPostController(service services.PostService, validate *validator.Validate) IPostController {
	return &PostController{
		Service:  service,
		Validate: validate,
	}
}

func (pc *PostController) GetPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := pc.Service.GetPosts(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch posts"})
	}

	return c.JSON(result)
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

	req := &dto.PostRequest{
		Title:   c.FormValue("title"),
		Content: c.FormValue("content"),
	}
	slug := c.FormValue("slug")

	post, err := pc.Service.CreatePost(req, authorUUID, filePath, slug, categoryIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(post)
}

func (pc *PostController) EditPost(c *fiber.Ctx) error {
	id := c.Params("id")

	authorIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, err := uuid.Parse(authorIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID format"})
	}

	role := c.Locals("role").(string)

	var req dto.PostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	post, err := pc.Service.EditPost(id, &req, userID, role)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "post not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
		}
		if errMsg == "access denied" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post"})
	}

	return c.JSON(fiber.Map{"message": "Post updated successfully", "post": post})
}

func (pc *PostController) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := pc.Service.DeletePost(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete post"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

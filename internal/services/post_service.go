package services

import (
	"errors"
	"fmt"
	"go_blog/internal/dto"
	"go_blog/internal/models"
	"go_blog/internal/repositories"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PostService interface {
	GetPosts(page, limit int) (map[string]interface{}, error)
	CreatePost(req *dto.PostRequest, authorID uuid.UUID, thumbnailPath, slug string, categoryIDs []int) (*models.Post, error)
	EditPost(id string, req *dto.PostRequest, userID uuid.UUID, role string) (*models.Post, error)
	DeletePost(id string) error
}

type postService struct {
	postRepo     repositories.PostRepository
	categoryRepo repositories.CategoryRepository
}

func NewPostService(postRepo repositories.PostRepository, categoryRepo repositories.CategoryRepository) PostService {
	return &postService{
		postRepo:     postRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *postService) GetPosts(page, limit int) (map[string]interface{}, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	total, err := s.postRepo.Count()
	if err != nil {
		return nil, err
	}

	posts, err := s.postRepo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var p []dto.PostResponse
	for _, post := range posts {
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

	return map[string]interface{}{
		"data": p,
		"meta": map[string]interface{}{
			"total_records": total,
			"current_page":  page,
			"total_pages":   math.Ceil(float64(total) / float64(limit)),
			"limit":         limit,
		},
	}, nil
}

func (s *postService) CreatePost(req *dto.PostRequest, authorID uuid.UUID, thumbnailPath, slug string, categoryIDs []int) (*models.Post, error) {
	var categoriesFound []models.Category
	if len(categoryIDs) > 0 {
		var err error
		categoriesFound, err = s.categoryRepo.FindByIds(categoryIDs)
		if err != nil {
			return nil, fmt.Errorf("database error fetching categories: %v", err)
		}
	}

	post := models.Post{
		Title:         req.Title,
		Content:       req.Content,
		Slug:          slug,
		AuthorID:      authorID,
		ThumbnailPath: thumbnailPath,
		Categories:    categoriesFound,
	}

	if err := s.postRepo.Create(&post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *postService) EditPost(id string, req *dto.PostRequest, userID uuid.UUID, role string) (*models.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.AuthorID != userID && role != "Administrator" {
		return nil, errors.New("access denied")
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) DeletePost(id string) error {
	return s.postRepo.Delete(id)
}

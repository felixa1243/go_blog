package services

import (
	"go_blog/internal/dto"
	"go_blog/internal/models"
	"go_blog/internal/repositories"

	"github.com/go-playground/validator/v10"
)

type CategoryService interface {
	GetCategories() ([]dto.CategoryResponse, error)
	CreateCategory(req *dto.CategoryRequest) (*models.Category, error)
	GetCategoriesByIds(ids []int) ([]models.Category, error)
}

type categoryService struct {
	repo      repositories.CategoryRepository
	validator *validator.Validate
}

func NewCategoryService(repo repositories.CategoryRepository, validator *validator.Validate) CategoryService {
	return &categoryService{
		repo:      repo,
		validator: validator,
	}
}

func (s *categoryService) GetCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	datas := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		datas = append(datas, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}
	return datas, nil
}

func (s *categoryService) CreateCategory(req *dto.CategoryRequest) (*models.Category, error) {
	category := models.Category{
		Name: req.Name,
	}
	if err := s.repo.Create(&category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *categoryService) GetCategoriesByIds(ids []int) ([]models.Category, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.repo.FindByIds(ids)
}

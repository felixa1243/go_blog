package repositories

import (
	"go_blog/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	Create(category *models.Category) error
	FindByIds(ids []int) ([]models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByIds(ids []int) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Where("id IN ?", ids).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

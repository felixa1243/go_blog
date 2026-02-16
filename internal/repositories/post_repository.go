package repositories

import (
	"go_blog/internal/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	Count() (int64, error)
	FindAll(limit, offset int) ([]models.Post, error)
	Create(post *models.Post) error
	FindByID(id string) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id string) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Count() (int64, error) {
	var total int64
	if err := r.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *postRepository) FindAll(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id string) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id string) error {
	return r.db.Delete(&models.Post{}, "id = ?", id).Error
}

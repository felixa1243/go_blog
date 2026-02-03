package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID            int        `json:"id" gorm:"primaryKey"`
	Title         string     `json:"title" gorm:"type:varchar(255);not null;unique"`
	Content       string     `json:"content" gorm:"type:text;not null"`
	Slug          string     `json:"slug" gorm:"type:varchar(255);not null;unique;index"`
	ThumbnailPath string     `json:"thumbnail_path"`
	AuthorID      uuid.UUID  `json:"author_id" gorm:"type:uuid;not null;index"`
	Categories    []Category `json:"categories" gorm:"many2many:post_categories;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

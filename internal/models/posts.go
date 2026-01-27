package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	AuthorID  uuid.UUID `json:"author_id" gorm:"not null"`
}

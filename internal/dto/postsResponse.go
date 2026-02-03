package dto

import (
	"github.com/google/uuid"
)

type PostResponse struct {
	ID            int       `json:"id"`
	Title         string    `json:"title" `
	Content       string    `json:"content" `
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	AuthorID      uuid.UUID `json:"author_id"`
	ThumbnailPath string    `json:"thumbnail_path"`
	Categories    []int     `json:"categories"`
	Slug          string    `json:"slug"`
}

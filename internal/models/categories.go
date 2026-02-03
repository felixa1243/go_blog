package models

type Category struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" gorm:"not null"`
	Posts []Post `json:"posts" gorm:"many2many:post_categories"`
}

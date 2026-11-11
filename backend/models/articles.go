package models

import "time"

type Article struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	AuthorID  uint           `gorm:"not null;index" json:"author_id"`
	Images    []ArticleImage `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE;" json:"images"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
type ArticleImage struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ArticleID uint   `gorm:"not null;index" json:"article_id"`
	URL       string `gorm:"not null" json:"url"`
}

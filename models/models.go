package models

import (
	"time"

	"github.com/lib/pq"
)

type Project struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Technologies  pq.StringArray `gorm:"type:text[]" json:"technologies"`
	PictureCover  string         `gorm:"not null" json:"picture_cover"`
	Pictures      pq.StringArray `gorm:"type:text[]" json:"pictures"`
	RepositoryURL string         `gorm:"type:text" json:"repository_url"`
	DemoURL       string         `gorm:"type:text" json:"demo_url"`
	Status        string         `gorm:"default:'ongoing'" json:"status"`
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

type Article struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"not null" json:"title"`
	Slug           string         `gorm:"unique;not null" json:"slug"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	Excerpt        string         `gorm:"type:text" json:"excerpt"`
	SEOTitle       string         `json:"seo_title"`
	SEODescription string         `gorm:"type:text" json:"seo_description"`
	CoverImage     string         `json:"cover_image"`
	Tags           pq.StringArray `gorm:"type:text[]" json:"tags"`
	Status         string         `gorm:"default:'draft'" json:"status"`
	PublishedAt    *time.Time     `json:"published_at"`
	AuthorID       uint           `json:"author_id"`
	ViewCount      uint           `json:"view_count"`
	Likes          uint           `json:"likes"`
	CommentsCount  uint           `json:"comments_count"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

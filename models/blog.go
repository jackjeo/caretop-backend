package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlogCategory string

const (
	BlogCategoryTech     BlogCategory = "tech"
	BlogCategoryProduct  BlogCategory = "product"
	BlogCategoryIndustry BlogCategory = "industry"
)

type BlogPost struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Title       string         `gorm:"size:300;not null" json:"title"`
	Slug        string         `gorm:"uniqueIndex;size:300;not null" json:"slug"`
	Summary     string         `gorm:"type:text" json:"summary"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	CoverImageURL string       `gorm:"size:500" json:"cover_image_url"`
	Category    BlogCategory   `gorm:"size:20;not null;default:'tech'" json:"category"`
	AuthorID    uuid.UUID      `gorm:"type:uuid;not null" json:"author_id"`
	Author      User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	LikeCount   int            `gorm:"default:0" json:"like_count"`
	IsPublished bool           `gorm:"default:true" json:"is_published"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *BlogPost) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.PublishedAt == nil && b.IsPublished {
		now := time.Now()
		b.PublishedAt = &now
	}
	return nil
}

type BlogComment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PostID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"post_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *BlogComment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

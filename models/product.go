package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Download struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Slug        string         `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Tagline     string         `gorm:"size:500" json:"tagline"`
	Description string         `gorm:"type:text" json:"description"`
	Features    string         `gorm:"type:jsonb" json:"features"`
	Screenshots string         `gorm:"type:jsonb" json:"screenshots"`
	Downloads   string         `gorm:"type:jsonb" json:"downloads"`
	IsPublished bool           `gorm:"default:true" json:"is_published"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

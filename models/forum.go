package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumBoard struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string         `gorm:"size:500" json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (f *ForumBoard) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type ForumThread struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	BoardID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"board_id"`
	Board       ForumBoard     `gorm:"foreignKey:BoardID" json:"board,omitempty"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title       string         `gorm:"size:300;not null" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	IsPinned    bool           `gorm:"default:false" json:"is_pinned"`
	IsEssential bool           `gorm:"default:false" json:"is_essential"`
	IsLocked    bool           `gorm:"default:false" json:"is_locked"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *ForumThread) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type ForumPost struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ThreadID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"thread_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentID  *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	LikeCount int            `gorm:"default:0" json:"like_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *ForumPost) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type ForumLike struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	ThreadID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"thread_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ForumCollection struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	ThreadID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"thread_id"`
	CreatedAt time.Time `json:"created_at"`
}

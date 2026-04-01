package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketType string

const (
	TicketTypeBug       TicketType = "bug"
	TicketTypeFeature   TicketType = "feature"
	TicketTypeConsult   TicketType = "consult"
	TicketTypeBusiness  TicketType = "business"
)

type TicketStatus string

const (
	TicketStatusPending     TicketStatus = "pending"
	TicketStatusProcessing  TicketStatus = "processing"
	TicketStatusResolved    TicketStatus = "resolved"
	TicketStatusClosed      TicketStatus = "closed"
)

type Ticket struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type      TicketType     `gorm:"size:20;not null;default:'consult'" json:"type"`
	Title     string         `gorm:"size:300;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Status    TicketStatus   `gorm:"size:20;not null;default:'pending'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TicketReply struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	TicketID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"ticket_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *TicketReply) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

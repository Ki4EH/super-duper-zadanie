package models

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name" validate:"required,min=3,max=100"`
	Description    string    `gorm:"type:text" json:"description"`
	Status         string    `gorm:"size:20;not null" json:"status"`
	Version        int       `gorm:"default:1" json:"version"`
	TenderID       uuid.UUID `gorm:"type:uuid;not null;size:100" json:"tenderId" validate:"required,max=100"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;size:100" json:"organizationId"`
	AuthorID       uuid.UUID `gorm:"type:uuid;not null;size:100" json:"authorId" validate:"required,max=100"`
	AuthorType     string    `gorm:"not null" json:"authorType" validate:"required,oneof=User Organization"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

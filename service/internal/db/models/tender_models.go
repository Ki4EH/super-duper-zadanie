package models

import (
	"github.com/google/uuid"
	"time"
)

type Tender struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;size:100" json:"id"`
	Name            string    `gorm:"size:100;not null" json:"name" validate:"required,min=3,max=100"`
	Description     string    `gorm:"type:text;size:500" json:"description"`
	Status          string    `gorm:"size:20;not null;default:'CREATED'" json:"status"`
	Version         int       `gorm:"default:1" json:"version"`
	ServiceType     string    `gorm:"size:20;not null" json:"serviceType" validate:"required,oneof=Construction Delivery Manufacture"`
	OrganizationID  uuid.UUID `gorm:"type:uuid;not null;size:100" json:"organizationId" validate:"required,max=100"`
	CreatorID       uuid.UUID `gorm:"type:uuid;not null" json:"creatorId"`
	CreatorUsername string    `gorm:"not null" json:"creatorUsername" validate:"required"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

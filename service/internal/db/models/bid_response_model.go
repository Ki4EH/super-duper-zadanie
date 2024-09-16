package models

import (
	"github.com/google/uuid"
	"time"
)

type BidResponse struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required,min=3"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:20;not null;default:'CREATED'" json:"status"`
	TenderID    uuid.UUID `gorm:"type:uuid;not null" json:"tenderId" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	AuthorType  string    `gorm:"not null" json:"authorType" validate:"required"`
	AuthorID    uuid.UUID `gorm:"type:uuid;not null" json:"authorId"`
	Version     int       `gorm:"default:1" json:"version"`
}

func ToBidResponse(bid Bid) BidResponse {
	return BidResponse{
		ID:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      bid.Status,
		Version:     bid.Version,
		TenderID:    bid.TenderID,
		AuthorType:  bid.AuthorType,
		AuthorID:    bid.AuthorID,
		CreatedAt:   bid.CreatedAt,
	}
}

package models

import (
	"github.com/google/uuid"
	"time"
)

type BidReview struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	BidID       uuid.UUID `gorm:"type:uuid;not null" json:"bid_id"`
	Description string    `gorm:"type:text;not null" json:"description" validate:"required,max=1000"`
	AuthorID    uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

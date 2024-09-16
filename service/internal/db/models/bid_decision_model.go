package models

import (
	"github.com/google/uuid"
	"time"
)

type BidDecision struct {
	ID        uuid.UUID `gorm:"primaryKey;" json:"id"`
	BidID     uuid.UUID `gorm:"type:uuid;not null" json:"bid_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Decision  string    `gorm:"size:20;not null" json:"decision"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

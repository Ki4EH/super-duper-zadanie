package models

import (
	"github.com/google/uuid"
	"time"
)

type BidReviewResponse struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Description string    `gorm:"type:text;not null" json:"description" validate:"required,max=1000"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func ToBidReviewResponse(review BidReview) BidReviewResponse {
	return BidReviewResponse{
		ID:          review.ID,
		Description: review.Description,
		CreatedAt:   review.CreatedAt,
	}
}

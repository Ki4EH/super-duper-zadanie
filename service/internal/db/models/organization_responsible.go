package models

import "github.com/google/uuid"

type OrganizationResponsible struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null" json:"organization_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
}

func (OrganizationResponsible) TableName() string {
	return "organization_responsible"
}

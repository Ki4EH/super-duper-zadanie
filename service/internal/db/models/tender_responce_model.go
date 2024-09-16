package models

import (
	"github.com/google/uuid"
	"time"
)

type TenderResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceType    string    `json:"serviceType"`
	Status         string    `json:"status"`
	OrganizationId uuid.UUID `json:"organizationId"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
}

func ToTenderResponse(tender Tender) TenderResponse {
	return TenderResponse{
		ID:             tender.ID,
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    tender.ServiceType,
		Status:         tender.Status,
		OrganizationId: tender.OrganizationID,
		Version:        tender.Version,
		CreatedAt:      tender.CreatedAt,
	}
}

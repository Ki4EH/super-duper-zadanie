package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *TenderHandler) ListTendersHandler(c echo.Context) error {
	ctx := c.Request().Context()

	serviceTypes := c.QueryParams()["service_type"]

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	filter := repository.TenderFilter{
		ServiceTypes: serviceTypes,
		Limit:        limit,
		Offset:       offset,
	}

	tenders, err := h.repo.ListTenders(ctx, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить список тендеров"})
	}

	var response []models.TenderResponse
	for _, tender := range tenders {
		if tender.Status == "Published" {
			response = append(response, models.ToTenderResponse(*tender))
		}
	}

	return c.JSON(http.StatusOK, response)
}

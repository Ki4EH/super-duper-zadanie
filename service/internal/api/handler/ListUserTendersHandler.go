package handler

import (
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/models"
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *TenderHandler) ListUserTendersHandler(c echo.Context) error {
	ctx := c.Request().Context()

	userName := c.QueryParam("username")

	id, err := h.orgRepo.GetUserUUID(c.Request().Context(), userName)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "пользователь не существует или некорректен"})
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	filter := repository.TenderFilter{
		CreatorID: id,
		Limit:     limit,
		Offset:    offset,
	}

	tenders, err := h.repo.ListTenders(ctx, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"response": "не удалось получить список тендеров"})
	}

	var response []models.TenderResponse
	for _, tender := range tenders {
		switch tender.Status {
		case "Created", "Closed":
			isResponsible, err := h.orgRepo.CheckOrganizationResponsible(ctx, tender.OrganizationID, id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"response": "не удалось проверить права доступа"})
			}
			if isResponsible {
				response = append(response, models.ToTenderResponse(*tender))
			}
		case "Published":
			response = append(response, models.ToTenderResponse(*tender))
		}
	}

	return c.JSON(http.StatusOK, response)
}

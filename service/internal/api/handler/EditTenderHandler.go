package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func (h *TenderHandler) EditTenderHandler(c echo.Context) error {
	tender := new(models.Tender)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	ctx := c.Request().Context()

	if err := c.Bind(tender); err != nil {
		log.Printf("ошибка привязки данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат данных"})
	}

	tender.ID = id

	updTender, err := h.repo.UpdateTender(ctx, tender)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "серверная ошибка при обновлении тендера"})
	}

	response := models.ToTenderResponse(*updTender)

	return c.JSON(http.StatusOK, response)
}

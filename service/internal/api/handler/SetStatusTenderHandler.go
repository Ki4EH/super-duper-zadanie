package handler

import (
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *TenderHandler) SetStatusTenderHandlerHandler(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат ID"})
	}

	status := c.QueryParam("status")
	if status != "Created" && status != "Published" && status != "Closed" {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "некорректный статус тендера"})
	}

	tender, err := h.repo.GetTenderByUUID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"reason": "тендер не найден"})
	}

	tender.Status = status
	tender.UpdatedAt = time.Now()

	updTender, err := h.repo.UpdateTender(ctx, tender)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "серверная ошибка при обновлении тендера"})
	}

	response := models.ToTenderResponse(*updTender)

	return c.JSON(http.StatusOK, response)
}

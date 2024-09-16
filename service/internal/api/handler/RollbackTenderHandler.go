package handler

import (
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *TenderHandler) RollbackTenderHandler(c echo.Context) error {
	tenderIDParam := c.Param("id")
	tenderID, err := uuid.Parse(tenderIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат tenderID"})
	}

	versionParam := c.Param("version")
	version, err := strconv.Atoi(versionParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат version"})
	}

	ctx := c.Request().Context()

	oldTender, err := h.repo.GetTenderVersion(ctx, tenderID, version)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "ошибка отката версии"})
	}

	oldTender.Version = version

	tender, err := h.repo.UpdateTender(ctx, oldTender)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "ошибка сохранения новой версии"})
	}

	response := models.ToTenderResponse(*tender)

	return c.JSON(http.StatusOK, response)
}

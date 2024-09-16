package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *BidHandler) SetStatusBidHandler(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат ID"})
	}

	status := c.QueryParam("status")
	if status != "Created" && status != "Published" && status != "Canceled" {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "некорректный статус предложения"})
	}

	bid, err := h.repo.GetBidByUUID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"reason": "предложение не найдено"})
	}

	bid.Status = status
	bid.UpdatedAt = time.Now()

	// Сохранение изменений в базе данных
	updBid, err := h.repo.UpdateBid(ctx, bid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "серверная ошибка при обновлении тендера"})
	}

	response := models.ToBidResponse(*updBid)

	return c.JSON(http.StatusOK, response)
}

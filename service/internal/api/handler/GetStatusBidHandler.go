package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *BidHandler) GetStatusBidHandler(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат ID"})
	}

	bid, err := h.repo.GetBidByUUID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"reason": "предложение не найдено"})
	}

	return c.JSON(http.StatusOK, bid.Status)
}

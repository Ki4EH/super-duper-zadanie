package handler

import (
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *BidHandler) RollbackBidHandler(c echo.Context) error {
	bidIDParam := c.Param("id")
	bidID, err := uuid.Parse(bidIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат bidID"})
	}

	versionParam := c.Param("version")
	version, err := strconv.Atoi(versionParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат version"})
	}

	ctx := c.Request().Context()

	oldBid, err := h.repo.GetBidVersion(ctx, bidID, version)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "ошибка отката версии"})
	}

	oldBid.Version = version

	bid, err := h.repo.UpdateBid(ctx, oldBid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "ошибка сохранения новой версии"})
	}

	response := models.ToBidResponse(*bid)

	return c.JSON(http.StatusOK, response)
}

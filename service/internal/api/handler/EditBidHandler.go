package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func (h *BidHandler) EditBidHandler(c echo.Context) error {

	bid := new(models.Bid)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		log.Printf("ошибка привязки данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат данных id"})
	}

	ctx := c.Request().Context()

	if err := c.Bind(bid); err != nil {
		log.Printf("ошибка привязки данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат данных"})
	}

	bid.ID = id

	updBid, err := h.repo.UpdateBid(ctx, bid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "серверная ошибка при обновлении предложения"})
	}

	response := models.ToBidResponse(*updBid)

	return c.JSON(http.StatusOK, response)
}

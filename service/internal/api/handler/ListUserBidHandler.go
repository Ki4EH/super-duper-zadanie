package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *BidHandler) ListUserBidHandler(c echo.Context) error {

	ctx := c.Request().Context()
	userName := c.QueryParam("username")

	// Получаем ID пользователя по username
	id, err := h.orgRepo.GetUserUUID(ctx, userName)
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

	// Фильтр для получения предложений пользователя
	filter := repository.BidFilter{
		CreatorID: id,
		Limit:     limit,
		Offset:    offset,
	}

	// Получаем список предложений пользователя
	bids, err := h.repo.ListUserBids(ctx, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"response": "не удалось получить список предложений"})
	}

	var response []models.BidResponse
	for _, bid := range bids {
		response = append(response, models.ToBidResponse(*bid))
	}

	if len(response) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"response": "предложение не найдено"})
	}

	return c.JSON(http.StatusOK, response)
}

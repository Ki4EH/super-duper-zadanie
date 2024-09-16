package handler

import (
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *BidHandler) GetAuthorBidReviewsHandler(c echo.Context) error {
	tenderID := c.Param("id")
	authorUsername := c.QueryParam("authorUsername")
	requesterUsername := c.QueryParam("requesterUsername")

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	ctx := c.Request().Context()

	requesterID, err := h.orgRepo.GetUserUUID(ctx, requesterUsername)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "Запрашивающий пользователь не найден"})
	}

	tenderUUID, err := uuid.Parse(tenderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "Некорректный идентификатор тендера"})
	}

	isResponsible, err := h.repo.CheckTenderResponsibility(ctx, tenderUUID, requesterID)
	if err != nil || !isResponsible {
		return c.JSON(http.StatusForbidden, map[string]string{"reason": "Недостаточно прав для просмотра отзывов"})
	}

	authorID, err := h.orgRepo.GetUserUUID(ctx, authorUsername)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "Автор предложений не найден"})
	}

	reviews, err := h.repo.GetBidReviewsByAuthor(ctx, tenderUUID, authorID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "Ошибка при получении отзывов"})
	}

	var response []models.BidReviewResponse
	for _, r := range reviews {
		response = append(response, models.ToBidReviewResponse(r))
	}
	return c.JSON(http.StatusOK, response)

}

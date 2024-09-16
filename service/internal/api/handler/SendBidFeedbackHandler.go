package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func (h *BidHandler) SendBidFeedbackHandler(c echo.Context) error {
	bidID := c.Param("id")
	bidFeedback := c.QueryParam("bidFeedback")
	username := c.QueryParam("username")
	ctx := c.Request().Context()
	if len(bidFeedback) == 0 || len(bidFeedback) > 1000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "Отзыв не может быть пустым и должен быть не более 1000 символов"})
	}

	authorID, err := h.orgRepo.GetUserUUID(ctx, username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "Пользователь не найден"})
	}

	bidUUID, err := uuid.Parse(bidID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "Некорректный идентификатор предложения"})
	}

	review := models.BidReview{
		BidID:       bidUUID,
		Description: bidFeedback,
		AuthorID:    authorID,
	}

	err = h.repo.AddBidReview(ctx, &review)
	if err != nil {
		log.Printf("Ошибка сохранения отзыва: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "Ошибка при сохранении отзыва"})
	}

	bid, err := h.repo.GetBidByUUID(ctx, bidUUID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "серверная ошибка при обновлении предложения"})
	}

	response := models.ToBidResponse(*bid)

	return c.JSON(http.StatusOK, response)
}

package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
)

func (h *BidHandler) SendBidDecisionHandler(c echo.Context) error {
	ctx := c.Request().Context()

	bidID := c.Param("id")
	decision := c.QueryParam("decision")
	username := c.QueryParam("username")

	if decision != "Approved" && decision != "Rejected" {
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "Некорректное решение. Доступные значения: Approved, Rejected"})
	}

	// Проверка существования пользователя и получение его ID
	userID, err := h.orgRepo.GetUserUUID(ctx, username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "Пользователь не найден или некорректен"})
	}

	// Получение предложения по ID
	bid, err := h.repo.GetBidByUUID(ctx, uuid.MustParse(bidID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"reason": "Предложение не найдено"})
	}

	isResponsible, err := h.orgRepo.CheckOrganizationResponsible(ctx, bid.OrganizationID, userID)
	if err != nil || !isResponsible {
		return c.JSON(http.StatusForbidden, map[string]string{"reason": "Недостаточно прав для принятия решения"})
	}

	decisionErr := h.repo.AddBidDecision(ctx, bid.ID, userID, decision)
	if decisionErr != nil {
		return c.JSON(http.StatusConflict, map[string]string{"reason": "Решение уже принято или другая ошибка"})
	}

	// Проверка текущего кворума и обновление статуса предложения и тендера
	err = h.processBidDecision(ctx, bid, decision)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "Ошибка обработки решения"})
	}

	return c.JSON(http.StatusOK, models.ToBidResponse(*bid))
}

func (h *BidHandler) processBidDecision(ctx context.Context, bid *models.Bid, decision string) error {
	decisions, err := h.repo.GetBidDecisions(ctx, bid.ID)
	if err != nil {
		return err
	}

	// Если есть хотя бы одно решение "Rejected", предложение отклоняется
	for _, d := range decisions {
		if d.Decision == "Rejected" {
			bid.Status = "Canceled"
			return h.repo.UpdateBidStatus(ctx, bid.ID, "Canceled")
		}
	}

	responsibleCount, err := h.repo.GetResponsibleCount(ctx, bid.OrganizationID)
	if err != nil {
		return err
	}

	approvedCount := 0
	for _, d := range decisions {
		if d.Decision == "Approved" {
			approvedCount++
		}
	}

	// Кворум: min(3, количество ответственных за организацию)
	quorum := min(3, responsibleCount)

	// Если достигнут кворум, предложение одобряется и тендер закрывается
	if approvedCount >= quorum {
		err = h.repo.UpdateBidStatus(ctx, bid.ID, "Published")
		if err != nil {
			return err
		}
		return h.repo.CloseTender(ctx, bid.TenderID)
	}
	return nil
}

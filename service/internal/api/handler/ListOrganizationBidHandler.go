package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *BidHandler) ListOrganizationBids(c echo.Context) error {
	var tenderId uuid.UUID
	var err error

	if c.Param("id") != "" {
		tenderIdParam := c.Param("id")
		tenderId, err = uuid.Parse(tenderIdParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат tenderID"})
		}
	}

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

	organizationID, err := h.orgRepo.GetOrganizationFromTender(ctx, tenderId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"reason": "тендер или предложение не найдено"})
	}

	// Проверяем, является ли пользователь ответственным за организацию
	isResponsible, err := h.orgRepo.CheckOrganizationResponsible(ctx, organizationID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "ошибка проверки ответственности"})
	}

	if !isResponsible {
		return c.JSON(http.StatusForbidden, map[string]string{"response": "у вас нет доступа к этим предложениям"})
	}

	// Фильтр для получения предложений организации
	filter := repository.BidFilter{
		TenderId:  tenderId,
		CreatorID: id,
		Limit:     limit,
		Offset:    offset,
	}

	// Получаем список предложений организации
	bids, err := h.repo.ListOrganizationBids(ctx, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"response": "не удалось получить список предложений организации"})
	}

	var response []models.BidResponse
	for _, bid := range bids {
		response = append(response, models.ToBidResponse(*bid))
	}

	if len(response) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"response": "предложения не найдены"})
	}

	return c.JSON(http.StatusOK, response)
}

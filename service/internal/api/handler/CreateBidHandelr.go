package handler

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type BidHandler struct {
	repo    repository.BidRepository
	orgRepo repository.OrganizationRepository
}

func NewBidHandler(repo repository.BidRepository, orgRepo repository.OrganizationRepository) *BidHandler {
	return &BidHandler{repo: repo, orgRepo: orgRepo}
}

func (h *BidHandler) CreateBidHandler(c echo.Context) error {
	bid := new(models.Bid)
	if err := c.Bind(bid); err != nil {
		log.Printf("ошибка привязки данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат данных"})
	}

	if err := c.Validate(bid); err != nil {
		log.Printf("ошибка валидации данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "данные не соответствуют требованиям"})
	}

	creatorUsername := bid.AuthorID
	if creatorUsername == uuid.Nil {
		log.Printf("ошибка: AuthorID отсутствует или пуст")
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "некорректное имя пользователя"})
	}

	ctx := c.Request().Context()
	organizationID, err := h.orgRepo.GetOrganizationFromTender(ctx, bid.TenderID)
	if organizationID == uuid.Nil || err != nil {
		log.Printf("ошибка: organizationID отсутствует или пуст")
		return c.JSON(http.StatusNotFound, map[string]string{"reason": "тендер не найден"})
	}

	bid.OrganizationID = organizationID

	flag, err := h.repo.CheckAuthorID(ctx, bid)
	if err != nil || flag == false {
		return c.JSON(http.StatusForbidden, map[string]string{"reason": "недостаточно прав"})
	}

	bid.Status = "Created"

	createdBid, err := h.repo.CreateBid(c.Request().Context(), bid)
	if err != nil {
		log.Printf("ошибка создания предложения %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "предложение уже создано по этому тендеру"})
	}

	response := models.ToBidResponse(*createdBid)

	return c.JSON(http.StatusOK, response)
}

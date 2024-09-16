package handler

import (
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/models"
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// TenderHandler структура для хранения зависимостей обработчика тендеров
type TenderHandler struct {
	repo    repository.TenderRepository
	orgRepo repository.OrganizationRepository
}

// NewTenderHandler создает новый экземпляр TenderHandler
func NewTenderHandler(repo repository.TenderRepository, orgRepo repository.OrganizationRepository) *TenderHandler {
	return &TenderHandler{repo: repo, orgRepo: orgRepo}
}

func (h *TenderHandler) CreateTenderHandler(c echo.Context) error {
	tender := new(models.Tender)
	if err := c.Bind(tender); err != nil {
		log.Printf("ошибка привязки данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "неверный формат данных"})
	}

	if err := c.Validate(tender); err != nil {
		log.Printf("ошибка валидации данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"reason": "данные не соответствуют требованиям"})
	}

	creatorUsername := tender.CreatorUsername
	if creatorUsername == "" {
		log.Printf("ошибка: creatorUsername отсутствует или пуст")
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "некорректное имя пользователя"})
	}

	organizationID := tender.OrganizationID
	if organizationID == uuid.Nil {
		log.Printf("ошибка: organizationID отсутствует или пуст")
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "некорректный идентификатор организации"})
	}

	// Получаем идентификатор пользователя по username
	id, err := h.orgRepo.GetUserUUID(c.Request().Context(), creatorUsername)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "пользователь не существует или некорректен"})
	}

	ctx := c.Request().Context()

	if flag, _ := h.orgRepo.CheckOrganizationResponsible(ctx, organizationID, id); flag == false {
		log.Printf("ошибка: пользователь не является ответственным за организацию")
		return c.JSON(http.StatusForbidden, map[string]string{"reason": "пользователь не имеет прав на выполнение действия"})
	}

	// Устанавливаем нужные значения в тендер
	tender.CreatorID = id
	tender.Status = "Created"

	// Сохранение тендера в базе данных
	createdTender, err := h.repo.CreateTender(ctx, tender)
	if err != nil {
		log.Printf("ошибка создания тендера %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "не удалось создать тендер"})
	}

	response := models.ToTenderResponse(*createdTender)

	return c.JSON(http.StatusOK, response)
}

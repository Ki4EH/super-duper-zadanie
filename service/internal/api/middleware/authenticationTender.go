package middleware

import (
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func AccessControlTenderMiddleware(orgRepo repository.OrganizationRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ctx := c.Request().Context()

			userName := c.QueryParam("username")

			id, err := orgRepo.GetUserUUID(c.Request().Context(), userName)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"reason": "пользователь не существует или некорректен"})
			}

			tenderIDParam := c.Param("id")
			tenderID, err := uuid.Parse(tenderIDParam)
			if err != nil {
				log.Printf("ошибка: парсинга tenderId: %v", err)
				return c.JSON(http.StatusBadRequest, map[string]string{"reason": "некорректный идентификатор организации"})
			}
			organizationID, err := orgRepo.GetOrganizationFromTender(ctx, tenderID)
			if err != nil || organizationID == uuid.Nil {
				log.Printf("ошибка: некорректный идентификатор организации: %v", err)
				return c.JSON(http.StatusBadRequest, map[string]string{"reason": "некорректный идентификатор организации"})
			}

			isResponsible, err := orgRepo.CheckOrganizationResponsible(ctx, organizationID, id)
			if err != nil {
				log.Printf("ошибка проверки ответственности: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"reason": "ошибка проверки прав доступа"})
			}

			if !isResponsible {
				log.Printf("ошибка: пользователь не является ответственным за организацию")
				return c.JSON(http.StatusForbidden, map[string]string{"reason": "пользователь не имеет прав на выполнение действия"})
			}

			return next(c)
		}
	}
}

package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

var validate = validator.New()

// ValidationMiddleware проверяет входящие данные и возвращает ошибку, если данные некорректны
func ValidationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := validate.Struct(c.Request().Body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверные данные запроса"})
		}
		return next(c)
	}
}

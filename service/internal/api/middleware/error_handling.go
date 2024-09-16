package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func ErrorHandlingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			c.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "внутренняя ошибка сервера"})
		}
		return nil
	}
}

package middleware

import (
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

func RequestLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		log.Printf("Метод: %s, URI: %s, Статус: %d, Время: %v", c.Request().Method, c.Request().RequestURI, c.Response().Status, time.Since(start))
		return err
	}
}

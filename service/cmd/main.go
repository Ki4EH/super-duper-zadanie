package main

import (
	"github.com/Ki4EH/super-duper-zadanie/service/config"
	"github.com/Ki4EH/super-duper-zadanie/service/internal/api"
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Загружаем конфигурацию из .env файла
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Устанавливаем соединение с базой данных
	dbConn, err := db.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Создаем новый экземпляр Echo
	e := echo.New()
	// Регистрация кастомного валидатора
	e.Validator = &CustomValidator{validator: validator.New()}
	// Инициализируем маршруты
	api.InitRoutes(e, dbConn)
	// Запускаем сервер
	log.Printf("Сервер запущен на %s", cfg.ServerAddress)
	if err := e.Start(cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

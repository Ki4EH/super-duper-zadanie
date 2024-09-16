package db

import (
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strings"
)

// NewPostgresConnection создает подключение к базе данных PostgreSQL с учетом различных форматов строки подключения
func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
	var dsn string

	if cfg.PostgresConn != "" {
		dsn = cfg.PostgresConn
		log.Println("Подключение к базе данных через POSTGRES_CONN")
	} else if cfg.PostgresJDBC != "" {
		dsn = convertJDBCToDSN(cfg.PostgresJDBC)
		log.Println("Подключение к базе данных через POSTGRES_JDBC_URL")
	} else {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDBName)
		log.Println("Подключение к базе данных через отдельные параметры")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	log.Println("Успешное подключение к базе данных")
	return db, nil
}

func convertJDBCToDSN(jdbcURL string) string {
	parsedURL := strings.Replace(jdbcURL, "jdbc:postgresql://", "postgres://", 1)
	parsedURL = strings.Replace(parsedURL, "?user=", "/", 1)
	parsedURL = strings.Replace(parsedURL, "&password=", ":", 1)
	return parsedURL
}

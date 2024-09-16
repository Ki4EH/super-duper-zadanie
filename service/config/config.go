package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ServerAddress  string
	PostgresConn   string
	PostgresJDBC   string
	PostgresUser   string
	PostgresPass   string
	PostgresHost   string
	PostgresPort   string
	PostgresDBName string
}

// LoadConfig загружает конфигурацию из .env файла или переменных окружения
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Нет .env файла, используется только переменные окружения: %v", err)
	}

	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")

	cfg := &Config{
		ServerAddress:  viper.GetString("SERVER_ADDRESS"),
		PostgresConn:   viper.GetString("POSTGRES_CONN"),
		PostgresJDBC:   viper.GetString("POSTGRES_JDBC_URL"),
		PostgresUser:   viper.GetString("POSTGRES_USERNAME"),
		PostgresPass:   viper.GetString("POSTGRES_PASSWORD"),
		PostgresHost:   viper.GetString("POSTGRES_HOST"),
		PostgresPort:   viper.GetString("POSTGRES_PORT"),
		PostgresDBName: viper.GetString("POSTGRES_DATABASE"),
	}

	return cfg, nil
}

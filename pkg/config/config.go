package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBot TelegramBot
	Database    DatabaseConfig
}

type TelegramBot struct {
	Token string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d", d.Username, d.Password, d.Host, d.Port)
}

func New() *Config {
	godotenv.Load()

	return &Config{
		TelegramBot: TelegramBot{
			Token: getEnv("TELEGRAM_BOT_TOKEN", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "admin"),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	res, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Wrong env variable value at key %s. Expected INTEGER value, got %s", key, value))
	}
	return res
}

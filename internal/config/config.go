package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	TelegramBot TelegramBot
	Database    DatabaseConfig
	JWT         JWTConfig
	Redis       RedisConfig
}

type TelegramBot struct {
	Token string
}

type JWTConfig struct {
	Secret string
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username, d.Password, d.Host, d.Port, d.Database)
}

func New() *Config {
	return &Config{
		TelegramBot: TelegramBot{
			Token: getEnv("TELEGRAM_BOT_TOKEN", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "admin"),
			Database: getEnv("DB_NAME", "duckwars"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "FILL_SECRET_KEY"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Username: getEnv("REDIS_USERNAME", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DATABASE", 0),
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

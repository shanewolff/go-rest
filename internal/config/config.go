package config

import (
	"os"
)

type Config struct {
	DSN      string
	APIToken string
	Addr     string
	LogLevel string
}

func LoadConfig() Config {
	return Config{
		DSN:      getEnv("DB_DSN", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC"),
		APIToken: getEnv("API_TOKEN", "secret123"),
		Addr:     getEnv("SERVER_ADDR", ":8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

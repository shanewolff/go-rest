package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN      string
	APIToken string
	Addr     string
	LogLevel string
	AppEnv   string
}

func LoadConfig() Config {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return Config{
		DSN:      getEnv("DB_DSN", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC"),
		APIToken: getEnv("API_TOKEN", "secret123"),
		Addr:     getEnv("SERVER_ADDR", ":8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		AppEnv:   getEnv("APP_ENV", "production"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN           string
	APIToken      string
	Addr          string
	LogLevel      string
	AppEnv        string
	JWTSecret     string
	JWTExpiration time.Duration
}

func LoadConfig() Config {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtExpiration, _ := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))

	return Config{
		DSN:           getEnv("DB_DSN", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC"),
		APIToken:      getEnv("API_TOKEN", "secret123"),
		Addr:          getEnv("SERVER_ADDR", ":8080"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		AppEnv:        getEnv("APP_ENV", "production"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key"),
		JWTExpiration: jwtExpiration,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

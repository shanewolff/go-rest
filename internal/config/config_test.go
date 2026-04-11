package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		// Ensure environment is clean
		os.Clearenv()

		cfg := LoadConfig()

		assert.Equal(t, ":8080", cfg.Addr)
		assert.Equal(t, "info", cfg.LogLevel)
		assert.Equal(t, "production", cfg.AppEnv)
		assert.Equal(t, "super-secret-key", cfg.JWTSecret)
		assert.Equal(t, 24*time.Hour, cfg.JWTExpiration)
	})

	t.Run("environment variables override", func(t *testing.T) {
		require.NoError(t, os.Setenv("SERVER_ADDR", ":9090"))
		require.NoError(t, os.Setenv("LOG_LEVEL", "debug"))
		require.NoError(t, os.Setenv("APP_ENV", "development"))
		require.NoError(t, os.Setenv("JWT_SECRET", "custom-secret"))
		require.NoError(t, os.Setenv("JWT_EXPIRATION", "1h"))

		cfg := LoadConfig()

		assert.Equal(t, ":9090", cfg.Addr)
		assert.Equal(t, "debug", cfg.LogLevel)
		assert.Equal(t, "development", cfg.AppEnv)
		assert.Equal(t, "custom-secret", cfg.JWTSecret)
		assert.Equal(t, 1*time.Hour, cfg.JWTExpiration)
	})
}

package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("production", func(t *testing.T) {
		require.NoError(t, os.Setenv("APP_ENV", "production"))
		l, err := NewLogger("info")
		assert.NoError(t, err)
		assert.NotNil(t, l)
	})

	t.Run("development", func(t *testing.T) {
		require.NoError(t, os.Setenv("APP_ENV", "development"))
		l, err := NewLogger("debug")
		assert.NoError(t, err)
		assert.NotNil(t, l)
	})

	t.Run("invalid level", func(t *testing.T) {
		l, err := NewLogger("invalid")
		assert.NoError(t, err)
		assert.NotNil(t, l)
	})
}

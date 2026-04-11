package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/shanewolff/go-rest/internal/domain"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Item{}, &domain.User{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestPostgresRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresRepository(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		item := &domain.Item{
			Title: "Test Item",
			Price: 10.5,
		}

		err := repo.Create(item)
		assert.NoError(t, err)
		assert.NotZero(t, item.ID)

		found, err := repo.GetByID(item.ID)
		assert.NoError(t, err)
		assert.Equal(t, item.Title, found.Title)
		assert.Equal(t, item.Price, found.Price)
	})

	t.Run("GetAll", func(t *testing.T) {
		items, err := repo.GetAll()
		assert.NoError(t, err)
		assert.NotEmpty(t, items)
	})

	t.Run("GetByID - Not Found", func(t *testing.T) {
		_, err := repo.GetByID(999)
		assert.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		item := &domain.Item{Title: "To Delete", Price: 1.0}
		err := repo.Create(item)
		assert.NoError(t, err)

		err = repo.Delete(item.ID)
		assert.NoError(t, err)

		_, err = repo.GetByID(item.ID)
		assert.Error(t, err)
	})
}

func TestPostgresUserRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	t.Run("Create and GetByUsername", func(t *testing.T) {
		user := &domain.User{
			Username:     "testuser",
			PasswordHash: "hashed",
		}

		err := repo.Create(user)
		assert.NoError(t, err)

		found, err := repo.GetByUsername("testuser")
		assert.NoError(t, err)
		assert.Equal(t, user.Username, found.Username)
	})

	t.Run("GetByUsername - Not Found", func(t *testing.T) {
		_, err := repo.GetByUsername("nonexistent")
		assert.Error(t, err)
	})
}

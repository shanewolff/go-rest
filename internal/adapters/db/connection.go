package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shanewolff/go-rest/internal/domain"
)

// InitDB initializes the database connection and performs migrations.
func InitDB(dsn string, logger *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate the schema for our internal DB models
	// This ensures the database schema matches our DB models.
	err = db.AutoMigrate(&Item{}, &domain.User{})
	if err != nil {
		return nil, err
	}

	logger.Info("Database connection established and migrations completed")
	return db, nil
}

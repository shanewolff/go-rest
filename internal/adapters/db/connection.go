package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection and performs migrations.
func InitDB(dsn string, logger *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate is disabled in favor of explicit golang-migrate CLI usage.
	// Ensure you run database migrations via `task migrate:up` before starting the application.

	logger.Info("Database connection established")
	return db, nil
}

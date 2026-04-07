package db

import (
	"gorm.io/gorm"

	"github.com/shanewolff/go-rest/internal/domain"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *postgresUserRepository) Create(user *domain.User) error {
	result := r.db.Create(user)
	return result.Error
}

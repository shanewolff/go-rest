package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/shanewolff/go-rest/internal/domain"
	"github.com/shanewolff/go-rest/internal/mocks"
)

func TestAuthService_Register(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	service := NewAuthService(repo, "secret", 1*time.Hour)

	req := domain.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}

	repo.EXPECT().GetByUsername("testuser").Return(nil, nil)
	repo.EXPECT().Create(mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.Register(req)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.NotEmpty(t, user.PasswordHash)
}

func TestAuthService_Login(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	service := NewAuthService(repo, "secret", 1*time.Hour)

	// To test login, we need a user with a hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := &domain.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	req := domain.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	repo.EXPECT().GetByUsername("testuser").Return(user, nil)

	token, returnedUser, err := service.Login(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.Username, returnedUser.Username)
}

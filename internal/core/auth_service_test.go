package core

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/shanewolff/go-rest/internal/domain"
	"github.com/shanewolff/go-rest/internal/mocks"
)

func TestAuthService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("user already exists", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, "secret", 1*time.Hour)

		repo.EXPECT().GetByUsername("existing").Return(&domain.User{Username: "existing"}, nil)
		_, err := service.Register(domain.RegisterRequest{Username: "existing", Password: "password"})
		assert.ErrorIs(t, err, ErrUserAlreadyExists)
	})

	t.Run("repo error on create", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, "secret", 1*time.Hour)

		repo.EXPECT().GetByUsername("newuser").Return(nil, nil)
		repo.EXPECT().Create(mock.Anything).Return(assert.AnError)

		_, err := service.Register(domain.RegisterRequest{Username: "newuser", Password: "password"})
		assert.Error(t, err)
	})
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, "secret", 1*time.Hour)

		req := domain.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		repo.EXPECT().GetByUsername("testuser").Return(user, nil)

		token, returnedUser, err := service.Login(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, user.Username, returnedUser.Username)
	})

	t.Run("invalid credentials - wrong password", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, "secret", 1*time.Hour)

		repo.EXPECT().GetByUsername("testuser").Return(user, nil)
		_, _, err := service.Login(domain.LoginRequest{Username: "testuser", Password: "wrongpassword"})
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("invalid credentials - user not found", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, "secret", 1*time.Hour)

		repo.EXPECT().GetByUsername("nonexistent").Return(nil, ErrInvalidCredentials)
		_, _, err := service.Login(domain.LoginRequest{Username: "nonexistent", Password: "any"})
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	secret := "secret"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, secret, 1*time.Hour)

		repo.EXPECT().GetByUsername("testuser").Return(user, nil)
		token, _, _ := service.Login(domain.LoginRequest{Username: "testuser", Password: "password123"})

		userID, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), userID)
	})

	t.Run("invalid token", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, secret, 1*time.Hour)
		_, err := service.ValidateToken("invalid-token")
		assert.ErrorIs(t, err, ErrTokenInvalid)
	})

	t.Run("expired token", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		expiredService := NewAuthService(repo, secret, -1*time.Hour)
		repo.EXPECT().GetByUsername("testuser").Return(user, nil)
		expiredToken, _, _ := expiredService.Login(domain.LoginRequest{Username: "testuser", Password: "password123"})

		service := NewAuthService(repo, secret, 1*time.Hour)
		_, err := service.ValidateToken(expiredToken)
		assert.ErrorIs(t, err, ErrTokenInvalid)
	})

	t.Run("invalid signing method", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, secret, 1*time.Hour)

		token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
			"user_id": 1,
		})
		tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

		_, err := service.ValidateToken(tokenString)
		assert.ErrorIs(t, err, ErrTokenInvalid)
	})

	t.Run("missing user_id claim", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, secret, 1*time.Hour)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		_, err := service.ValidateToken(tokenString)
		assert.ErrorIs(t, err, ErrTokenInvalid)
	})

	t.Run("invalid user_id claim type", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, secret, 1*time.Hour)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "not-a-number",
			"exp":     time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		_, err := service.ValidateToken(tokenString)
		assert.ErrorIs(t, err, ErrTokenInvalid)
	})

	t.Run("wrong secret", func(t *testing.T) {
		repo := mocks.NewUserRepository(t)
		service := NewAuthService(repo, secret, 1*time.Hour)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": 1,
			"exp":     time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte("wrong-secret"))

		_, err := service.ValidateToken(tokenString)
		assert.ErrorIs(t, err, ErrTokenInvalid)
	})
}

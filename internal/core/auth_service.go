package core

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/shanewolff/go-rest/internal/domain"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrTokenInvalid       = errors.New("invalid token")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type authService struct {
	repo          domain.UserRepository
	jwtSecret     []byte
	jwtExpiration time.Duration
}

func NewAuthService(repo domain.UserRepository, jwtSecret string, jwtExpiration time.Duration) domain.AuthService {
	return &authService{
		repo:          repo,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
	}
}

func (s *authService) Register(req domain.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	existing, _ := s.repo.GetByUsername(req.Username)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req domain.LoginRequest) (string, *domain.User, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.jwtExpiration).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *authService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrTokenInvalid
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrTokenInvalid
	}

	return uint(userIDFloat), nil
}

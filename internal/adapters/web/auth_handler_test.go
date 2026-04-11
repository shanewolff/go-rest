package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/shanewolff/go-rest/internal/domain"
	"github.com/shanewolff/go-rest/internal/mocks"
)

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := mocks.NewAuthService(t)
		logger := zap.NewNop()
		handler := NewAuthHandler(mockService, logger)

		reqBody := domain.RegisterRequest{
			Username: "testuser",
			Password: "password123",
		}
		expectedUser := &domain.User{
			ID:       1,
			Username: "testuser",
		}

		mockService.EXPECT().Register(reqBody).Return(expectedUser, nil)

		r := gin.Default()
		r.POST("/register", handler.Register)

		jsonValue, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var actualUser domain.User
		err := json.Unmarshal(w.Body.Bytes(), &actualUser)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Username, actualUser.Username)
		assert.Equal(t, expectedUser.ID, actualUser.ID)
	})

	t.Run("bind error", func(t *testing.T) {
		mockService := mocks.NewAuthService(t)
		logger := zap.NewNop()
		handler := NewAuthHandler(mockService, logger)

		r := gin.Default()
		r.POST("/register", handler.Register)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := mocks.NewAuthService(t)
		logger := zap.NewNop()
		handler := NewAuthHandler(mockService, logger)

		reqBody := domain.RegisterRequest{
			Username: "testuser",
			Password: "password123",
		}

		mockService.EXPECT().Register(reqBody).Return(nil, errors.New("registration failed"))

		r := gin.Default()
		r.POST("/register", handler.Register)

		jsonValue, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := mocks.NewAuthService(t)
		logger := zap.NewNop()
		handler := NewAuthHandler(mockService, logger)

		reqBody := domain.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		expectedUser := &domain.User{
			ID:       1,
			Username: "testuser",
		}
		expectedToken := "mock-token"

		mockService.EXPECT().Login(reqBody).Return(expectedToken, expectedUser, nil)

		r := gin.Default()
		r.POST("/login", handler.Login)

		jsonValue, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp domain.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, resp.Token)
		assert.Equal(t, expectedUser.Username, resp.User.Username)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockService := mocks.NewAuthService(t)
		logger := zap.NewNop()
		handler := NewAuthHandler(mockService, logger)

		reqBody := domain.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}

		mockService.EXPECT().Login(reqBody).Return("", nil, errors.New("unauthorized"))

		r := gin.Default()
		r.POST("/login", handler.Login)

		jsonValue, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("bind error", func(t *testing.T) {
		mockService := mocks.NewAuthService(t)
		logger := zap.NewNop()
		handler := NewAuthHandler(mockService, logger)

		r := gin.Default()
		r.POST("/login", handler.Login)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

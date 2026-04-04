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

func TestGetItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewItemService(t)
	logger := zap.NewNop()
	handler := NewItemHandler(mockService, "test-token", logger)

	expectedItems := []domain.Item{
		{ID: 1, Title: "Item 1", Price: 10.0},
	}
	mockService.EXPECT().GetAllItems().Return(expectedItems, nil)

	r := gin.Default()
	r.GET("/items", handler.GetItems)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualItems []domain.Item
	err := json.Unmarshal(w.Body.Bytes(), &actualItems)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actualItems))
	assert.Equal(t, "Item 1", actualItems[0].Title)
}

func TestGetItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewItemService(t)
	logger := zap.NewNop()
	handler := NewItemHandler(mockService, "test-token", logger)

	t.Run("success", func(t *testing.T) {
		expectedItem := &domain.Item{ID: 1, Title: "Item 1", Price: 10.0}
		mockService.EXPECT().GetItem(uint(1)).Return(expectedItem, nil)

		r := gin.Default()
		r.GET("/items/:id", handler.GetItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualItem domain.Item
		err := json.Unmarshal(w.Body.Bytes(), &actualItem)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), actualItem.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.EXPECT().GetItem(uint(2)).Return(nil, errors.New("not found"))

		r := gin.Default()
		r.GET("/items/:id", handler.GetItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/2", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCreateItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewItemService(t)
	logger := zap.NewNop()
	handler := NewItemHandler(mockService, "test-token", logger)

	reqBody := domain.CreateItemRequest{Title: "New Item", Price: 15.0}
	expectedItem := &domain.Item{ID: 1, Title: "New Item", Price: 15.0}

	mockService.EXPECT().CreateItem(reqBody).Return(expectedItem, nil)

	r := gin.Default()
	r.POST("/items", handler.CreateItem)

	jsonValue, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	handler := NewItemHandler(nil, "secret-key", logger)

	r := gin.New()
	r.Use(handler.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("authorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("X-API-Token", "secret-key")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("X-API-Token", "wrong-key")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

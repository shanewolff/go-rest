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
	handler := NewItemHandler(mockService, logger)

	t.Run("success", func(t *testing.T) {
		expectedItems := []domain.Item{
			{ID: 1, Title: "Item 1", Price: 10.0},
		}
		mockService.EXPECT().GetAllItems().Return(expectedItems, nil).Once()

		r := gin.New()
		r.GET("/items", handler.GetItems)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		mockService.EXPECT().GetAllItems().Return(nil, errors.New("service error")).Once()

		r := gin.New()
		r.GET("/items", handler.GetItems)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewItemService(t)
	logger := zap.NewNop()
	handler := NewItemHandler(mockService, logger)

	t.Run("success", func(t *testing.T) {
		expectedItem := &domain.Item{ID: 1, Title: "Item 1", Price: 10.0}
		mockService.EXPECT().GetItem(uint(1)).Return(expectedItem, nil).Once()

		r := gin.New()
		r.GET("/items/:id", handler.GetItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id format", func(t *testing.T) {
		r := gin.New()
		r.GET("/items/:id", handler.GetItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/abc", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.EXPECT().GetItem(uint(2)).Return(nil, errors.New("not found")).Once()

		r := gin.New()
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
	handler := NewItemHandler(mockService, logger)

	t.Run("success", func(t *testing.T) {
		reqBody := domain.CreateItemRequest{Title: "New Item", Price: 15.0}
		expectedItem := &domain.Item{ID: 1, Title: "New Item", Price: 15.0}

		mockService.EXPECT().CreateItem(reqBody).Return(expectedItem, nil).Once()

		r := gin.New()
		r.POST("/items", handler.CreateItem)

		jsonValue, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		r := gin.New()
		r.POST("/items", handler.CreateItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/items", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		reqBody := domain.CreateItemRequest{Title: "Error Item", Price: 10.0}
		mockService.EXPECT().CreateItem(reqBody).Return(nil, errors.New("service error")).Once()

		r := gin.New()
		r.POST("/items", handler.CreateItem)

		jsonValue, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDeleteItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewItemService(t)
	logger := zap.NewNop()
	handler := NewItemHandler(mockService, logger)

	t.Run("success", func(t *testing.T) {
		mockService.EXPECT().DeleteItem(uint(1)).Return(nil)

		r := gin.Default()
		r.DELETE("/items/:id", handler.DeleteItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/items/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		r := gin.Default()
		r.DELETE("/items/:id", handler.DeleteItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/items/abc", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService.EXPECT().DeleteItem(uint(2)).Return(errors.New("not found"))

		r := gin.Default()
		r.DELETE("/items/:id", handler.DeleteItem)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/items/2", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shanewolff/go-rest/internal/domain"
)

// This package acts as the Inbound Adapter for Gin Web Framework.

type ItemHandler struct {
	service domain.ItemService
	logger  *zap.Logger
}

func NewItemHandler(s domain.ItemService, l *zap.Logger) *ItemHandler {
	return &ItemHandler{
		service: s,
		logger:  l,
	}
}

// --- HTTP Handlers mapping web requests to core logic ---

func (h *ItemHandler) GetItems(c *gin.Context) {
	items, err := h.service.GetAllItems()
	if err != nil {
		h.logger.Error("Failed to fetch items", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) GetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	item, err := h.service.GetItem(uint(id))
	if err != nil {
		h.logger.Warn("Item not found", zap.Uint64("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req domain.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdItem, err := h.service.CreateItem(req)
	if err != nil {
		h.logger.Error("Failed to create item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}
	c.JSON(http.StatusCreated, createdItem)
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.service.DeleteItem(uint(id))
	if err != nil {
		h.logger.Error("Failed to delete item", zap.Uint64("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

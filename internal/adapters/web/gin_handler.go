package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-rest/internal/domain"
)

// This package acts as the Inbound Adapter for Gin Web Framework.

type ItemHandler struct {
	service  domain.ItemService
	apiToken string
	logger   *zap.Logger
}

func NewItemHandler(s domain.ItemService, apiToken string, l *zap.Logger) *ItemHandler {
	return &ItemHandler{
		service:  s,
		apiToken: apiToken,
		logger:   l,
	}
}

// SetupRouter initializes a new Gin router, registers middleware, and configures routes for the handler.
func (h *ItemHandler) SetupRouter() *gin.Engine {
	router := h.NewRouter()
	h.RegisterRoutes(router)
	return router
}

// NewRouter initializes a new Gin router with middleware
func (h *ItemHandler) NewRouter() *gin.Engine {
	router := gin.New() // Use gin.New() to have full control over middleware
	router.Use(gin.Recovery())
	router.Use(h.CustomLogger())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return router
}

// RegisterRoutes configures the routes for the item handler
func (h *ItemHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	api.Use(h.AuthMiddleware())
	{
		api.GET("/items", h.GetItems)
		api.GET("/items/:id", h.GetItem)
		api.POST("/items", h.CreateItem)
		api.DELETE("/items/:id", h.DeleteItem)
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

// --- MIDDLEWARE ---

func (h *ItemHandler) CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(startTime)
		h.logger.Info("Incoming Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		)
	}
}

func (h *ItemHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-API-Token")
		if token != h.apiToken {
			h.logger.Warn("Unauthorized access attempt",
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized API token"})
			return
		}
		c.Next()
	}
}

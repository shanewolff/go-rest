package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shanewolff/go-rest/internal/domain"
)

type Middleware struct {
	authService domain.AuthService
	apiToken    string
	logger      *zap.Logger
}

func NewMiddleware(authService domain.AuthService, apiToken string, logger *zap.Logger) *Middleware {
	return &Middleware{
		authService: authService,
		apiToken:    apiToken,
		logger:      logger,
	}
}

func (m *Middleware) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in the format Bearer <token>"})
			return
		}

		userID, err := m.authService.ValidateToken(parts[1])
		if err != nil {
			m.logger.Warn("Invalid token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func (m *Middleware) APITokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-API-Token")
		if token != m.apiToken {
			m.logger.Warn("Unauthorized access attempt",
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized API token"})
			return
		}
		c.Next()
	}
}

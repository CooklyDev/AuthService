package presentation

import (
	"fmt"
	"strings"

	"github.com/CooklyDev/AuthService/internal"
	"github.com/CooklyDev/AuthService/internal/application/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetContainer(container *internal.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("container", container)
		c.Next()
	}
}

func extractAuthService(c *gin.Context) (*usecases.AuthService, error) {
	containerValue, exists := c.Get("container")
	if !exists {
		return nil, fmt.Errorf("container not found")
	}

	container, ok := containerValue.(*internal.Container)
	if !ok {
		return nil, fmt.Errorf("invalid container type")
	}

	authService, err := container.GetAuthService(extractSessionID(c))
	if err != nil {
		return nil, fmt.Errorf("failed to get auth service: %w", err)
	}

	return authService, nil
}

func extractSessionID(c *gin.Context) uuid.UUID {
	sessionIDHeader := strings.TrimSpace(c.GetHeader("X-Session-ID"))
	if sessionIDHeader == "" {
		return uuid.Nil
	}

	sessionID, err := uuid.Parse(sessionIDHeader)
	if err != nil {
		return uuid.Nil
	}

	return sessionID
}

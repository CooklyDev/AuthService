package presentation

import (
	"github.com/CooklyDev/AuthService/internal"

	"github.com/gin-gonic/gin"
)

func GetContainer(container *internal.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("container", container)
		c.Next()
	}
}

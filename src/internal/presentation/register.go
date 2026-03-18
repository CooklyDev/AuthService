package presentation

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/CooklyDev/AuthService/internal"
// 	"github.com/CooklyDev/AuthService/internal/application"
// 	"github.com/CooklyDev/AuthService/internal/domain"

// 	"github.com/gin-gonic/gin"
// )

// func Register(c *gin.Context) {
// 	username := c.PostForm("username")
// 	email := c.PostForm("email")
// 	password := c.PostForm("password")

// 	containerValue, exists := c.Get("container")
// 	if !exists {
// 		FailWithStatus(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
// 		return
// 	}

// 	container, ok := containerValue.(*internal.Container)
// 	if !ok {
// 		FailWithStatus(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
// 		return
// 	}

// 	authService, err := container.GetAuthService()
// 	if err != nil {
// 		FailWithStatus(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
// 		return
// 	}

// 	if err := authService.Register(username, email, password); err != nil {
// 		status, code, message := mapRegisterError(err)
// 		FailWithStatus(c, status, code, message)
// 		return
// 	}

// 	Ok(c, nil)
// }

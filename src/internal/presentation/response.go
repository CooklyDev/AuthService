package presentation

import (
	"errors"
	"net/http"

	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Fail(c *gin.Context, status int, code string, message string) {
	c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func FailServerError(c *gin.Context) {
	Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
}

func MapAppError(err error) (int, string, string) {
	// Map application errors to HTTP status codes and error messages
	if errors.Is(err, domain.ErrBusinessRule) {
		return http.StatusBadRequest, "BUSINESS_RULE_VIOLATION", err.Error()
	}

	if errors.Is(err, adapters.ErrAdapter) {
		return http.StatusBadRequest, "ADAPTER_ERROR", err.Error()
	}

	// Default to internal server error for unrecognized errors
	return http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error"
}

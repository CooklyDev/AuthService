package presentation

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logout godoc
// @Summary Logout user
// @Description Terminates the current user session.
// @Tags Authentication
// @Produce json
// @Param X-Session-ID header string true "Session ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /logout [post]
func Logout(c *gin.Context) {
	sessionIDHeader := strings.TrimSpace(c.GetHeader("X-Session-ID"))
	if sessionIDHeader == "" {
		Fail(c, http.StatusBadRequest, "MISSING_SESSION_ID", "session ID is required")
		return
	}

	sessionID, err := uuid.Parse(sessionIDHeader)
	if err != nil {
		Fail(c, http.StatusBadRequest, "INVALID_SESSION_ID", "session ID must be a valid UUID")
		return
	}

	authService, err := extractAuthService(c)
	if err != nil {
		FailServerError(c)
		return
	}

	if err := authService.Logout(sessionID); err != nil {
		status, code, message := MapAppError(err)
		Fail(c, status, code, message)
		return
	}

	Ok(c, nil)
}

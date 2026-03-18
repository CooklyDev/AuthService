package presentation

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	sessionID := strings.TrimSpace(c.GetHeader("X-Session-ID"))
	if sessionID == "" {
		Fail(c, http.StatusBadRequest, "MISSING_SESSION_ID", "session ID is required")
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

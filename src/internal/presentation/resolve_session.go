package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ResolveSession godoc
// @Summary Resolve current session
// @Description Returns session and user identifiers for the provided session ID.
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param session_id formData string true "Session ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /resolve [post]
func ResolveSession(c *gin.Context) {
	sessionIDRaw := c.PostForm("session_id")
	if sessionIDRaw == "" {
		Fail(c, http.StatusBadRequest, "MISSING_SESSION_ID", "session ID is required")
		return
	}

	sessionID, err := uuid.Parse(sessionIDRaw)
	if err != nil {
		Fail(c, http.StatusBadRequest, "INVALID_SESSION_ID", "session ID must be a valid UUID")
		return
	}

	authService, err := extractAuthService(c)
	if err != nil {
		FailServerError(c)
		return
	}

	resolvedSession, err := authService.ResolveSession(sessionID)
	if err != nil {
		status, code, message := MapAppError(err)
		Fail(c, status, code, message)
		return
	}

	Ok(c, resolvedSession)
}

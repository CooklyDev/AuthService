package presentation

import "github.com/gin-gonic/gin"

// Register godoc
// @Summary Register user
// @Description Creates a new user account and session.
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {object} SessionIDResponse "Returns the created session ID in data"
// @Failure 400 {object} ErrorResponse "Business rule violation or adapter invariant violation"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /register [post]
func Register(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	authService, err := extractAuthService(c)
	if err != nil {
		FailServerError(c)
		return
	}

	session, err := authService.LocalRegister(username, email, password)
	if err != nil {
		status, code, message := MapAppError(err)
		Fail(c, status, code, message)
		return
	}

	Ok(c, session)
}

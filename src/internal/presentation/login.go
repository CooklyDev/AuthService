package presentation

import "github.com/gin-gonic/gin"

// Login godoc
// @Summary Login user
// @Description Authenticates user credentials and creates a session.
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /login [post]
func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	authService, err := extractAuthService(c)
	if err != nil {
		FailServerError(c)
		return
	}

	session, err := authService.Login(email, password)
	if err != nil {
		status, code, message := MapAppError(err)
		Fail(c, status, code, message)
		return
	}

	Ok(c, session)
}

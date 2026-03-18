package presentation

import "github.com/gin-gonic/gin"

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

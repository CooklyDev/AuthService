package presentation

import "github.com/gin-gonic/gin"

func Health(c *gin.Context) {
	Ok(c, gin.H{"status": "ok"})
}

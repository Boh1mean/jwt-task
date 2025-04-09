package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Welcome(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email not found in context"})
		return
	}

	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"Email": email,
	})
}

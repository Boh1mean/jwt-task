package handler

import (
	"jwtservertask/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateTokenHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	tokenStr := strings.Trim(strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer")), "\"")
	claims, err := utils.ValidateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token is valid", "claims": claims})
}

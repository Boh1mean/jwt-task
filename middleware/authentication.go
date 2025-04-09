package middleware

import (
	"fmt"
	"jwtservertask/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("Authorization")
		if clientToken == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(clientToken, "Bearer ", "", 1)
		if tokenString == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		fmt.Println("[AuthenticationMiddleware] User is authenticated")
		ctx.Set("userID", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Next()
	}
}

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
		clientToken := ctx.GetHeader("Authorization")

		if clientToken != "" && strings.HasPrefix(clientToken, "Bearer ") {
			clientToken = strings.TrimPrefix(clientToken, "Bearer ")
		} else {
			// Если нет — пробуем из cookie
			var err error
			clientToken, err = ctx.Cookie("access_token")
			if err != nil || clientToken == "" {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		claims, err := utils.ValidateToken(clientToken)
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

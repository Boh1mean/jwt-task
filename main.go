package main

import (
	"jwtservertask/handler"
	"jwtservertask/initializers"
	"jwtservertask/service"

	"github.com/gin-gonic/gin"
)

func main() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()

	tokenService := service.NewTokenService(initializers.DB)
	authService := service.NewAuthService(tokenService)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	r.POST("/signup", authHandler.SignUp)
	r.POST("/login", authHandler.Login)
	r.GET("/validate", handler.ValidateTokenHandler)
	r.GET("/refresh", authHandler.RefreshTokenHandler)

	r.Run()
}

package main

import (
	"jwtservertask/handler"
	"jwtservertask/initializers"
	"jwtservertask/middleware"
	"jwtservertask/service"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	tokenService := service.NewTokenService(initializers.DB)
	authService := service.NewAuthService(tokenService)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	r.POST("/signup", authHandler.SignUp)
	r.POST("/login", authHandler.Login)
	r.GET("/validate", handler.ValidateTokenHandler)
	r.GET("/refresh", authHandler.RefreshTokenHandler)

	protected := r.Group("/user")
	protected.Use(middleware.Authentication())
	{
		protected.GET("/home", handler.Welcome)
	}

	r.Run()
}

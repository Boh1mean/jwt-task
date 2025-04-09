package main

import (
	"jwtservertask/handler"
	"jwtservertask/initializers"
	"jwtservertask/middleware"
	"jwtservertask/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	tokenService := service.NewTokenService(initializers.DB)
	authService := service.NewAuthService(tokenService)
	//authHandler := handler.NewAuthHandler(authService)
	loginFormHandler := handler.NewLoginFormHandler(authService)

	r := gin.Default()

	// r.POST("/signup", authHandler.SignUp)
	// r.POST("/login", authHandler.Login)
	// r.GET("/login", authHandler.Login)
	// r.GET("/validate", handler.ValidateTokenHandler)
	// r.GET("/refresh", authHandler.RefreshTokenHandler)

	r.LoadHTMLGlob("templates/*")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", loginFormHandler.Login)

	r.GET("/homepage", middleware.Authentication(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "homepage.html", nil)
	})

	r.Run()
}

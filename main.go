package main

import (
	"jwtservertask/handler"
	"jwtservertask/initializers"
	"jwtservertask/middleware"
	"jwtservertask/repository"
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
	tokenRepo := repository.NewPostgresTokenRepository()
	authService := service.NewAuthService(tokenService, tokenRepo)

	// AuthHandler := handler.NewAuthHandler(authService)
	loginFormHandler := handler.NewLoginFormHandler(authService)
	SignUpFormHandler := handler.NewSignUpFormHandler(authService)
	LogoutFormHandler := handler.NewLogoutFormHandler(authService)

	r := gin.Default()

	// r.POST("/signup", authHandler.SignUp)
	// r.POST("/login", authHandler.Login)
	// r.GET("/login", authHandler.Login)
	r.GET("/validate", handler.ValidateTokenHandler)
	// r.GET("/refresh", authHandler.RefreshTokenHandler)

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})
	r.POST("/signup", SignUpFormHandler.SignUp)

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", loginFormHandler.Login)

	r.POST("/logout", LogoutFormHandler.Logout)

	r.GET("/user/homepage", middleware.Authentication(), func(c *gin.Context) {
		email, _ := c.Get("email")
		c.HTML(http.StatusOK, "homepage.html", gin.H{"Email": email})
	})

	// protected := r.Group("/user")
	// protected.Use(middleware.Authentication())
	// {
	// 	protected.GET("/homepage", handler.Welcome)
	// }

	r.Run()
}

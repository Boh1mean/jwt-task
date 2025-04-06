package main

import (
	"jwtservertask/handler"
	"jwtservertask/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()

	r.POST("/signup", handler.SignUp)
	r.POST("/login", handler.Login)

	r.Run()
}

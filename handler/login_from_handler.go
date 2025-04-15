package handler

import (
	"jwtservertask/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginFormHandler struct {
	authService *service.AuthService
}

func NewLoginFormHandler(authService *service.AuthService) *LoginFormHandler {
	return &LoginFormHandler{authService: authService}
}

func (h *LoginFormHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Email and password are required.",
		})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(email, password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", accessToken, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", false, true)

	c.Redirect(http.StatusFound, "/user/homepage")
}

package handler

import (
	"jwtservertask/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignUpFormHandler struct {
	authService *service.AuthService
}

func NewSignUpFormHandler(authService *service.AuthService) *SignUpFormHandler {
	return &SignUpFormHandler{authService: authService}
}

func (h *SignUpFormHandler) SignUp(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	err := h.authService.SignUp(email, password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(email, password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	c.Redirect(http.StatusFound, "/homepage")
}

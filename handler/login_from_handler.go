package handler

import (
	"jwtservertask/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginFormHandler обработка логина через HTML форму
type LoginFormHandler struct {
	authService *service.AuthService
}

func NewLoginFormHandler(authService *service.AuthService) *LoginFormHandler {
	return &LoginFormHandler{authService: authService}
}

func (h *LoginFormHandler) Login(c *gin.Context) {
	// Получаем данные из формы
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
		// Если ошибка, показываем страницу с ошибкой
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Устанавливаем cookies для авторизации
	c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	// Редирект на домашнюю страницу после успешного логина
	c.Redirect(http.StatusFound, "/homepage")
}

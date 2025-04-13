package handler

import (
	"jwtservertask/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogoutFormHandler struct {
	authService *service.AuthService
}

func NewLogoutFormHandler(authService *service.AuthService) *LogoutFormHandler {
	return &LogoutFormHandler{authService: authService}
}
func (h *LogoutFormHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	hash := service.HashToken(refreshToken)
	err = h.authService.Logout(hash)
	if err != nil {
		c.HTML(http.StatusBadRequest, "homepage.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/login")
}

package handler

import (
	"jwtservertask/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type signupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /signup
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req signupReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	err := h.authService.SignUp(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// POST /login
func (h *AuthHandler) Login(c *gin.Context) {
	var req signupReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		3600,
		"/",
		"localhost",
		false,
		true,
	)
	// можно установить refreshToken как httpOnly cookie
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GET /refresh
func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token provided"})
		return
	}

	newAccessToken, newRefreshToken, err := h.authService.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

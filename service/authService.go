package service

import (
	"errors"
	"jwtservertask/initializers"
	"jwtservertask/models"
	"jwtservertask/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	tokenService *TokenService
}

func NewAuthService(tokenService *TokenService) *AuthService {
	return &AuthService{
		tokenService: tokenService,
	}
}

func (s *AuthService) SignUp(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create User
	user := models.User{Email: email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (s *AuthService) Login(email, password string) (tokenString string, refreshToken string, err error) {
	var user models.User
	result := initializers.DB.First(&user, "email = ?", email)

	if result.Error != nil || user.ID == 0 {
		return "", "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", "", errors.New("token not create")
	}

	refreshToken, err = utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = s.tokenService.SaveRefreshToken(refreshToken, user.ID, expiresAt)
	if err != nil {
		return "", "", errors.New("failed to save refresh token ")
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, string, error) {
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	storedToken, err := s.tokenService.FindByToken(refreshToken)
	if err != nil || storedToken == nil || storedToken.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("refresh token expired or not found")
	}

	accessToken, err := utils.GenerateToken(claims.UserID, claims.Email)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(claims.UserID, claims.Email)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 дней
	err = s.tokenService.SaveRefreshToken(newRefreshToken, claims.UserID, expiresAt)
	if err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, newRefreshToken, nil
}

package service

import (
	"errors"
	"jwtservertask/initializers"
	"jwtservertask/models"
	"jwtservertask/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

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

func (s *AuthService) Login(email, password string) (string, error) {
	var user models.User
	result := initializers.DB.First(&user, "email = ?", email)

	if result.Error != nil || user.ID == 0 {
		return "", errors.New("invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", errors.New("token not create")
	}

	return token, nil
}

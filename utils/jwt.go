package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {

	var secretKey = []byte(os.Getenv("SECRET"))

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("[FAIL]: could not sign token: %w", err)
	}

	return tokenString, nil
}

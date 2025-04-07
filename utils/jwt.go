package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("SECRET"))

type SignedDetails struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, userEmail string) (string, error) {

	claims := SignedDetails{
		UserID: userID,
		Email:  userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("[FAIL]: could not sign token: %w", err)
	}

	return tokenString, nil
}

// func ValidateToken(tokenString string) (models.User, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return []byte(secretKey), nil
// 	})

// 	// user := models.User{}
// 	// if err != nil {
// 	// 	return user, err
// 	// }

// 	// payload, ok := token.Claims.(jwt.MapClaims)
// 	// if ok && token.Valid {
// 	// 	user.Email = payload["userID"].(string)

// 	// 	return user, nil
// 	// }

// 	// return user, errors.New("invalid token")
// }

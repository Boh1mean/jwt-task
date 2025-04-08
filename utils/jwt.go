package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "a-string-secret-at-least-256-bits-long" // строка, не []byte

type SignedDetails struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, userEmail string) (string, error) {

	claims := &SignedDetails{
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

	fmt.Println("[OK]: Generated token:", tokenString)
	return tokenString, nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	fmt.Println("[INFO]: Validating token:", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println("[ERROR]: Unexpected signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Println("[ERROR]: Token parse error:", err)
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		fmt.Println("[ERROR]: Claims could not be cast")
		return nil, fmt.Errorf("invalid token claims")
	}

	if !token.Valid {
		fmt.Println("[ERROR]: Token is not valid")
		fmt.Printf("[DEBUG]: Token expires at: %v, now: %v\n", claims.ExpiresAt.Time, time.Now())
		return nil, fmt.Errorf("invalid token")
	}

	fmt.Printf("[OK]: Token is valid. Claims: %+v\n", claims)
	return claims, nil
}

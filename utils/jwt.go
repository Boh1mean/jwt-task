package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "a-string-secret-at-least-256-bits-long" // строка, не []byte

type SignedDetails struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshSignedDetails struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
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

func GenerateRefreshToken(userID uint, userEmail string) (string, error) {

	claims := RefreshSignedDetails{
		UserID:    userID,
		Email:     userEmail,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return refreshToken, nil
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

func ValidateRefreshToken(refreshToken string) (*RefreshSignedDetails, error) {
	fmt.Println("[INFO]: Validating refresh token:", refreshToken)

	token, err := jwt.ParseWithClaims(refreshToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Println("[ERROR]: Refresh token parse error:", err)
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshSignedDetails)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

func SetAccessAndRefreshTokenCookies(c *gin.Context, tokenString, refreshToken string) {
	c.SetCookie(
		"access_token",
		tokenString,
		int(3600),
		"/",
		"localhost",
		false,
		true,
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(time.Hour.Seconds()*24*7),
		"/",
		"localhost",
		false,
		true,
	)
}

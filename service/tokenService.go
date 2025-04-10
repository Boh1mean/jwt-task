package service

import (
	"crypto/sha256"
	"encoding/hex"
	"jwtservertask/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type TokenService struct {
	db *gorm.DB
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{db: db}
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *TokenService) SaveRefreshToken(token string, userID uint, expiresAt time.Time) error {
	if s.db == nil {
		log.Println("Error: DB connection is nil")
	}
	hashed := HashToken(token)
	refresh := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hashed,
		ExpiresAt: expiresAt,
	}

	err := s.db.Create(&refresh).Error
	if err != nil {
		log.Println("Error saving refresh token:", err)
	}
	return err
}

func (s *TokenService) FindByToken(token string) (*models.RefreshToken, error) {
	hashedToken := HashToken(token)
	var refreshToken models.RefreshToken
	result := s.db.Where("token_hash = ?", hashedToken).First(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}

	return &refreshToken, nil
}

func (s *TokenService) ValidateRefreshToken(token string) (*models.RefreshToken, error) {
	hashed := HashToken(token)
	var refresh models.RefreshToken
	err := s.db.Where("token_hash = ?", hashed).First(&refresh).Error
	if err != nil {
		return nil, err
	}

	return &refresh, nil
}

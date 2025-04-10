package repository

import (
	"jwtservertask/models"
)

type TokenRepository interface {
	Save(token *models.RefreshToken) error
	FindByHash(hash string) (*models.RefreshToken, error)
	DeleteByHash(hash string) error
	DeleteAllByUser(userID uint) error
}

package repository

import (
	"jwtservertask/initializers"
	"jwtservertask/models"
)

type PostgresTokenRepository struct{}

func NewPostgresTokenRepository() TokenRepository {
	return &PostgresTokenRepository{}
}

func (r *PostgresTokenRepository) Save(token *models.RefreshToken) error {
	return initializers.DB.Create(token).Error
}

func (r *PostgresTokenRepository) FindByHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := initializers.DB.First(&token, "token_hash = ?", hash).Error
	return &token, err
}

func (r *PostgresTokenRepository) DeleteByHash(hash string) error {
	return initializers.DB.Where("token_hash = ?", hash).Delete(&models.RefreshToken{}).Error
}

func (r *PostgresTokenRepository) DeleteAllByUser(userID uint) error {
	return initializers.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

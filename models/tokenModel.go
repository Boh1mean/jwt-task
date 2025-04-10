package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint       `gorm:"index"`
	TokenHash string     `gorm:"uniqueIndex"`
	DeletedAt *time.Time `json:"deleted_at"`
	ExpiresAt time.Time
}

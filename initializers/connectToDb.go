package initializers

import (
	"jwtservertask/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // ← глобальная переменная

func ConnectToDb() *gorm.DB {
	var err error
	dsn := "host=localhost user=postgres password=pas123 dbname=jwt_test port=5432 sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to db ")
	}

	DB.AutoMigrate(&models.User{})

	return DB
}

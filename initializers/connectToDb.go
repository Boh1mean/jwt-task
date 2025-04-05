package initializers

import (
	"jwtservertask/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDb() *gorm.DB {
	var err error
	dsn := "host=localhost user=postgres password=pas123 dbname=jwt_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to db ")
	}

	db.AutoMigrate(&models.User{})

	return db
}

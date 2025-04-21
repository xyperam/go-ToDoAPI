package utils

import (
	"fmt"
	"log"

	"go-web-server/config"
	"go-web-server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	//auto-migrate model
	if err := database.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal("Failed to auto-migrate database:", err)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = database
	fmt.Println("Connected to database")
}

func CheckConnection() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	return nil
}

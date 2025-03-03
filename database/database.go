package database

import (
	"fmt"
	"log"
	"os"

	"be-porto-v3/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// LoadEnv membaca file .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables instead")
	}
}

// ConnectDB membuat koneksi ke database dan mengembalikan *gorm.DB dan error
func ConnectDB() (*gorm.DB, error) {
	LoadEnv() // Panggil LoadEnv untuk membaca konfigurasi

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate models
	err = db.AutoMigrate(&models.Project{}, &models.User{}, &models.Article{})
	if err != nil {
		return nil, err
	}

	DB = db
	log.Println("Database connected successfully!")
	return DB, nil
}

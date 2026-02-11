package database

import (
	"fmt"
	"hobby-loop/m/internal/models"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func Connect() {
	once.Do(func() {
		// Build DSN from Environment Variables
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSL"),
		)

		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database!", err)
		}

		// Auto migration
		DB.AutoMigrate(&models.User{}, &models.Basket{}, &models.Subscription{}, &models.Order{}, &models.Address{})

		log.Println("Database connected and migrated.")
	})
}

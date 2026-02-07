package database

import (
	"log"
	"subscription-market/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// in prod use os.Getenv("DNS")
	dns := "host=localhost user=postgres password=postgres dbname=market port=5432 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	// Auto migration
	DB.AutoMigrate(&models.User{}, &models.Basket{}, &models.Subscription{}, &models.Order{})

	log.Println("Database connected and migrated.")

}
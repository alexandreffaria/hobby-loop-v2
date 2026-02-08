package main

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/handlers"
	"hobby-loop/m/internal/models"
	"hobby-loop/m/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.Connect()
	database.DB.AutoMigrate(&models.Address{})

	// Start the background worker
	worker.Start()

	// Initialize Gin router
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		_, err := database.DB.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "Database connection error"})
			return
		}
		c.JSON(200, gin.H{"status": "alive", "database": "connected"})
	})

	router.POST("/register", handlers.RegisterUser)

	router.POST("/baskets", handlers.CreateBasket)

	router.POST("/subscriptions", handlers.SubscribeToBasket)

	router.GET("/orders", handlers.GetOrders)

	// Start the server
	router.Run(":8080")
}

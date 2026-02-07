package main

import (
	"hobby-loop/m/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.Connect()

	// Initialize Gin router
	router := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := database.DB.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "Database connection error"})
			return
		}
		c.jSON(200, gin.H{"status": "alive", "database": "connected"})
	})

	// Start the server
	router.Run(":8080")
}
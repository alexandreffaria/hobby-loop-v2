package handlers

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"

	"net/http"
	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {
	var orders []models.Order

	// Fetch all orders from the database
	if result := database.DB.Find(&orders); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)

}
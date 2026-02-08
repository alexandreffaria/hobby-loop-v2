package handlers

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SubscribeInput struct {
	UserID   uint `json:"user_id" binding:"required"`
	BasketID uint `json:"basket_id" binding:"required"`
}

func SubscribeToBasket(c *gin.Context) {
	var input SubscribeInput

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepare subscription data
	subscription := models.Subscription{
		UserID:           input.UserID,
		BasketID:         input.BasketID,
		Status:           "active",
		NextDeliveryDate: time.Now(),
	}

	// Save subscription to database
	if err := database.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

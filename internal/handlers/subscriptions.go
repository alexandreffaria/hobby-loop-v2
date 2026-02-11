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

func CancelSubscription(c *gin.Context) {
	id := c.Param("id")

	// FIX: Explicitly cast the float64 claim from JWT to uint
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(val.(float64))

	var sub models.Subscription

	// Security: Only find the subscription if it belongs to this user
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&sub).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	// Logic: Update status
	sub.Status = "cancelled"
	if err := database.DB.Save(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully", "status": sub.Status})
}

package handlers

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"

	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOrders(c *gin.Context) {

	userID := c.MustGet("user_id").(float64)
	isSeller := c.MustGet("is_seller").(bool)

	var orders []models.Order

	query := database.DB.Model(&models.Order{}).
		Preload("Subscription.User").
		Preload("Subscription.Basket").
		Joins("JOIN subscriptions ON subscriptions.id = orders.subscription_id").
		Joins("JOIN baskets ON baskets.id = subscriptions.basket_id")

	if isSeller {
		query = query.Where("baskets.seller_id = ?", uint(userID))
	} else {
		query = query.Where("subscriptions.user_id = ?", uint(userID))
	}

	// Fetch all orders from the database
	if result := query.Find(&orders); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)

}

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(float64)

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	result := database.DB.Joins("JOIN subscriptions ON subscriptions.id = orders.subscription_id").
		Joins("JOIN baskets ON baskets.id = subscriptions.basket_id").
		Where("orders.id = ? AND baskets.seller_id = ?", id, uint(userID)).
		First(&order)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or you don't have permission to update it"})
		return
	}

	order.Status = input.Status
	database.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

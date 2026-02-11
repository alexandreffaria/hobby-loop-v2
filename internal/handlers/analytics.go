package handlers

import (
	"hobby-loop/m/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SellerDashboardStats struct {
	ActiveSubscribers int64 `json:"active_subscribers"`
	GrossRevenue      float64 `json:"gross_revenue"`
	PlatformFees      float64 `json:"platform_fees"`
	NetEarnings       float64 `json:"net_revenue"`
}

func GetSellerDashboard(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := val.(float64)

	var stats SellerDashboardStats

	database.DB.Table("subscriptions").
		Joins("JOIN baskets ON baskets.id = subscriptions.basket_id").
		Where("baskets.seller_id = ? AND subscriptions.status = ?", userID, "active").
		Count(&stats.ActiveSubscribers)

	type FinancialResult struct {
		Gross float64
		Fee   float64
		Net   float64
	}

	var fin FinancialResult

	database.DB.Table("orders").
		Select("COALESCE(SUM(amount_paid), 0) AS gross, COALESCE(SUM(platform_fee), 0) AS fee, COALESCE(SUM(seller_net), 0) AS net").
		Joins("JOIN subscriptions ON subscriptions.id = orders.subscription_id").
		Joins("JOIN baskets ON baskets.id = subscriptions.basket_id").
		Where("baskets.seller_id = ?", userID).
		Scan(&fin)

	stats.GrossRevenue = fin.Gross
	stats.PlatformFees = fin.Fee
	stats.NetEarnings = fin.Net

	c.JSON(http.StatusOK, stats)
}

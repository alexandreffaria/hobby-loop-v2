package worker

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
	"log"
	"os"
	"strconv"
	"time"
)

func Start() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			ProcessSubscriptions()
		}
	}()
}

func ProcessSubscriptions() {
	var subs []models.Subscription
	now := time.Now()

	// Find subscriptions with next delivery date due
	database.DB.Where("next_delivery_date <= ? AND status = ?", now, "active").Find(&subs)

	for _, sub := range subs {
		// Fetch basket details
		var basket models.Basket
		database.DB.First(&basket, sub.BasketID)

		feePercentage := 0.1 // Default to 10%
		if envFee := os.Getenv("PLATFORM_FEE"); envFee != "" {
			if value, err := strconv.ParseFloat(envFee, 64); err == nil {
				feePercentage = value
			} else {
				log.Printf("Invalid PLATFORM_FEE value: %s, using default 10%%", envFee)
			}
		}

		fee := basket.Price * feePercentage
		net := basket.Price - fee

		// Create order
		order := models.Order{
			SubscriptionID: sub.ID,
			AmountPaid:     basket.Price,
			PlatformFee:    fee,
			SellerNet:      net,
			Status:         "pending_payment",
		}
		database.DB.Create(&order)

		log.Printf("STUB: Order %d created. Status 'pending_payment'. Waiting for external payment event.", order.ID)
		// Update next delivery date based on basket interval
		sub.NextDeliveryDate = calculateNextDeliveryDate(sub.NextDeliveryDate, basket.Interval)
		database.DB.Save(&sub)
	}
}

func calculateNextDeliveryDate(current time.Time, interval string) time.Time {
	switch interval {
	case "weekly":
		return current.AddDate(0, 0, 7)
	case "biweekly":
		return current.AddDate(0, 0, 14)
	case "monthly":
		return current.AddDate(0, 1, 0)
	default:
		return current.AddDate(0, 1, 0) // Default to monthly
	}
}

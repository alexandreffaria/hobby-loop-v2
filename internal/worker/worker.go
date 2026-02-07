package worker

import (
	"log"
	"time"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
)

func Start() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			processSubscriptions()
		}
	}()
}

func processSubscriptions() {
	var subs []models.Subscription
	now := time.Now()

	// Find subscriptions with next delivery date due
	result := database.DB.Where("next_delivery_date <= ? AND status = ?", now, "active").Find(&subs)
	if result.Error != nil {
		log.Println("Worker error: ", result.Error)
		return
	}

	if len (subs) > 0 {
		log.Printf("Worker: Found %d subscriptions to process\n", len(subs))
	}

	for _, sub := range subs {
		// Create order for each subscription
		order := models.Order{
			SubscriptionID: sub.ID,
			AmountPaid:     0.0, // This should be set based on basket price
			Status:         "pending_delivery",
		}
		database.DB.Create(&order)

		// Update next delivery date (for simplicity, adding 30 days)
		sub.NextDeliveryDate = sub.NextDeliveryDate.AddDate(0, 0, 30)
		database.DB.Save(&sub)

		log.Printf("Worker: Created order %d for subscription %d\n", order.ID, sub.ID)
	}
}
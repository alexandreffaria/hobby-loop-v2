package worker

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
	"log"
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

		// Create order
		order := models.Order{
			SubscriptionID: sub.ID,
			AmountPaid:     basket.Price,
			Status:         "processing_payment",
		}
		database.DB.Create(&order)

		// Simulate payment processing (replace with real payment logic)
		go emitInvoice(order.ID)

		// Update next delivery date based on basket interval
		sub.NextDeliveryDate = calculateNextDeliveryDate(sub.NextDeliveryDate, basket.Interval)
		database.DB.Save(&sub)
	}
}

func emitInvoice(orderID uint) {
	// Simulate delay for invoice emission
	time.Sleep(5 * time.Second)

	var order models.Order
	database.DB.First(&order, orderID)

	// Simulate invoice details
	order.InvoiceKey = "35231212345678000199550010000000011000000000"
	order.InvoiceURL = "https://invoices.example.com/" + order.InvoiceKey
	order.Status = "paid_and_invoiced"
	database.DB.Save(&order)

	log.Printf("Invoice emitted for order %d: Key=%s, URL=%s", order.ID, order.InvoiceKey, order.InvoiceURL)
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

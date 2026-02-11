package worker_test

import (
	"github.com/stretchr/testify/assert"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
	"hobby-loop/m/internal/worker" // Import the worker package
	"os"
	"testing"
	"time"
)

func TestProcessSubscriptions_GeneratesOrder(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "market")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL", "disable")

	// 1. Setup
	database.Connect()

	// Clean slate
	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM subscriptions")
	database.DB.Exec("DELETE FROM baskets")
	database.DB.Exec("DELETE FROM users")

	// 2. Seed Data
	user := models.User{Email: "worker_test@test.com", Password: "123"}
	database.DB.Create(&user)

	// We create a subscription that was due YESTERDAY
	basket := models.Basket{Price: 50.00, Interval: "weekly", Active: true}
	database.DB.Create(&basket)

	sub := models.Subscription{
		BasketID:         basket.ID,
		UserID:           1, // Dummy user ID
		Status:           "active",
		NextDeliveryDate: time.Now().AddDate(0, 0, -1), // Due yesterday
	}
	database.DB.Create(&sub)

	// 3. Action: Run the Worker Logic Manually
	// (We export ProcessSubscriptions in worker.go for this exact reason)
	worker.ProcessSubscriptions()

	// 4. Assertions
	var order models.Order
	result := database.DB.Where("subscription_id = ?", sub.ID).First(&order)

	// Check if Order exists
	assert.NoError(t, result.Error, "Worker should have created an order")
	assert.Equal(t, 50.00, order.AmountPaid, "Order price should match basket")
	assert.Equal(t, "processing_payment", order.Status)

	// Check if Subscription Date updated (Should be in the future now)
	var updatedSub models.Subscription
	database.DB.First(&updatedSub, sub.ID)
	assert.True(t, updatedSub.NextDeliveryDate.After(time.Now()), "Next delivery should be in the future")
}

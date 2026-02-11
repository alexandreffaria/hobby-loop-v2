package worker_test

import (
	"github.com/stretchr/testify/assert"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
	"hobby-loop/m/internal/worker"
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
	os.Setenv("PLATFORM_FEE", "0.1") // Ensure fee is set for test

	database.Connect()

	// 1. Clean slate - STRICT ORDER
	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM subscriptions")
	database.DB.Exec("DELETE FROM baskets")
	database.DB.Exec("DELETE FROM addresses") // Must be before users
	database.DB.Exec("DELETE FROM users")

	// 2. Seed Data
	user := models.User{Email: "worker_test@test.com", Password: "123"}
	database.DB.Create(&user)

	basket := models.Basket{Price: 50.00, Interval: "weekly", Active: true}
	database.DB.Create(&basket)

	sub := models.Subscription{
		BasketID:         basket.ID,
		UserID:           user.ID, // Use the ID from the created user
		Status:           "active",
		NextDeliveryDate: time.Now().AddDate(0, 0, -1),
	}
	database.DB.Create(&sub)

	// 3. Action
	worker.ProcessSubscriptions()

	// 4. Assertions
	var order models.Order
	result := database.DB.Where("subscription_id = ?", sub.ID).First(&order)

	assert.NoError(t, result.Error, "Worker should have created an order")
	assert.Equal(t, 50.00, order.AmountPaid)
	assert.Equal(t, "processing_payment", order.Status)

	var updatedSub models.Subscription
	database.DB.First(&updatedSub, sub.ID)
	assert.True(t, updatedSub.NextDeliveryDate.After(time.Now()))
}

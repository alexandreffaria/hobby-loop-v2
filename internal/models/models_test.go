package models_test

import (
	"github.com/stretchr/testify/assert"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
	"os"
	"testing"
)

func TestCreateUser(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "market")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL", "disable")

	database.Connect()

	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM subscriptions")
	database.DB.Exec("DELETE FROM users")

	user := models.User{
		Email:    "test@example.com",
		Password: "securepassword",
		IsSeller: true,
	}

	result := database.DB.Create(&user)

	assert.NoError(t, result.Error, "Should not return an error")
	assert.NotZero(t, user.ID, "User ID should be set after creation")
	assert.Equal(t, int64(1), result.RowsAffected, "One row should be affected")
}

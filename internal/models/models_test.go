package models_test

import (
	"testing"
	"hobby-loop/m/internal/models"
	"hobby-loop/m/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Connect to the database
	database.Connect()

	database.DB.Exec("DELETE FROM users") // Clean up before test
	
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
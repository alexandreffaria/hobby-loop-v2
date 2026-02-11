package auth_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"hobby-loop/m/internal/auth"
)

func TestTokenLifecycle(t *testing.T) {
	// Setup Key
	os.Setenv("JWT_SECRET", "unit-test-secret")

	// 1. Generate
	userID := uint(50)
	isSeller := true
	tokenString, err := auth.GenerateToken(userID, isSeller)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// 2. Validate (Success)
	claims, err := auth.ValidateToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, float64(userID), claims["user_id"]) // JWT numbers are float64
	assert.Equal(t, isSeller, claims["is_seller"])

	// 3. Validate (Fail - Tampered Token)
	// We add some garbage to the end of a valid token
	badToken := tokenString + "fake"
	_, err = auth.ValidateToken(badToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}

package handlers_test

import (
	"hobby-loop/m/internal/database"
	"os"
	"testing"
    "github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// 1. Global Setup for the "handlers" package
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "market")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL", "disable")
	os.Setenv("JWT_SECRET", "test-secret")

    // Set Gin to Test Mode globally
    gin.SetMode(gin.TestMode)

	// 2. Connect Once
	database.Connect()

	// 3. Run All Tests
	code := m.Run()

	// 4. Exit
	os.Exit(code)
}
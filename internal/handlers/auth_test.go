package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/handlers"
	"hobby-loop/m/internal/models"
)

func TestLogin(t *testing.T) {
	// 1. Setup
	// Env vars for DB connection
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "market")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL", "disable")
	os.Setenv("JWT_SECRET", "test-secret")

	database.Connect()

	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM subscriptions")
	database.DB.Exec("DELETE FROM users")

	// Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/login", handlers.Login)

	// 2. Seed a User with a HASHED password
	password := "supersecret"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		Email:    "login_test@example.com",
		Password: string(hashed), // Simulate what RegisterUser does
		IsSeller: false,
	}
	database.DB.Create(&user)

	// 3. Test Cases
	t.Run("Success - Correct Credentials", func(t *testing.T) {
		payload := map[string]string{
			"email":    "login_test@example.com",
			"password": "supersecret",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	t.Run("Fail - Wrong Password", func(t *testing.T) {
		payload := map[string]string{
			"email":    "login_test@example.com",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")
	})

	t.Run("Fail - User Not Found", func(t *testing.T) {
		payload := map[string]string{
			"email":    "ghost@example.com",
			"password": "any",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})
}

package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/logger"

	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/handlers"
	"hobby-loop/m/internal/models"
)

func TestLogin(t *testing.T) {
	// SILENCE THE NOISE: Don't print SQL errors for this test
	// (We check for nil just in case, though TestMain fixes it)
	if database.DB != nil {
		database.DB.Logger = logger.Default.LogMode(logger.Silent)
	}

	// Clean DB
	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM subscriptions")
	database.DB.Exec("DELETE FROM addresses")
	database.DB.Exec("DELETE FROM users")

	// Setup Router
	r := gin.Default()
	r.POST("/login", handlers.Login)

	// 2. Seed a User
	password := "supersecret"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		Email:    "login_test@example.com",
		Password: string(hashed),
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

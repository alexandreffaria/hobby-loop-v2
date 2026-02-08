package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", handlers.RegisterUser)
	return r
}

func TestRegisterUser(t *testing.T) {
	// 1. Setup
	database.Connect()
	// Clean DB to avoid unique constraint errors
	database.DB.Exec("DELETE FROM users")
	database.DB.Exec("DELETE FROM addresses")

	router := SetupRouter()

	// 2. Define Test Cases
	tests := []struct {
		name         string
		payload      map[string]interface{}
		expectedCode int
	}{
		{
			name: "Success - Valid CPF",
			payload: map[string]interface{}{
				"email":         "joao@example.com",
				"password":      "123456",
				"full_name":     "João da Silva",
				"document_type": "CPF",
				"document":      "12345678900", // A valid generated testing CPF
				"is_seller":     false,
			},
			expectedCode: 201,
		},
		{
			name: "Fail - Invalid CPF",
			payload: map[string]interface{}{
				"email":         "bad@example.com",
				"password":      "123456",
				"full_name":     "Bad CPF Guy",
				"document_type": "CPF",
				"document":      "11111111111", // Known invalid CPF
				"is_seller":     false,
			},
			// Note: This will only pass if you implemented the check in handlers/user.go!
			// If it returns 201, it means your validation logic is missing.
			expectedCode: 400,
		},
	}

	// 3. Run Tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

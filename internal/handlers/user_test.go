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
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/handlers"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", handlers.RegisterUser)
	return r
}

func TestRegisterUser(t *testing.T) {
	// 1. Setup
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "market")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL", "disable")
	database.Connect()
	// Clean DB to avoid unique constraint errors
	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM subscriptions")
	database.DB.Exec("DELETE FROM addresses")
	database.DB.Exec("DELETE FROM users")

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
				"document":      "09205193690",
				"is_seller":     false,
				"address": map[string]interface{}{
					"street":   "Rua Teste",
					"number":   "123",
					"city":     "São Paulo",
					"state":    "SP",
					"zip_code": "01001-000",
				},
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
				"document":      "12345678912",
				"is_seller":     false,
			},
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

			assert.Equal(t, tt.expectedCode, w.Code, "Request failed. Response Body: %s", w.Body.String())
		})
	}
}

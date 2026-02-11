package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"hobby-loop/m/internal/auth"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/handlers"
	"hobby-loop/m/internal/models"
	"hobby-loop/m/internal/worker"
)

// Helper to setup the router with all our logic
func SetupFullRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Public
	r.POST("/register", handlers.RegisterUser)

	// Protected
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/baskets", handlers.CreateBasket)
		protected.GET("/baskets", handlers.ListBaskets)
		protected.POST("/subscriptions", handlers.SubscribeToBasket)
		protected.GET("/orders", handlers.GetOrders)
		protected.PATCH("/orders/:id", handlers.UpdateOrderStatus)
	}
	return r
}

func TestTheMarketplaceFlow(t *testing.T) {

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "market")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL", "disable")
	os.Setenv("JWT_SECRET", "test-secret") // Needed for token generation!
	// 1. INFRASTRUCTURE SETUP
	database.Connect()
	// Clean slate
	database.DB.Exec("DELETE FROM addresses")
	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM subscriptions")
	database.DB.Exec("DELETE FROM baskets")
	database.DB.Exec("DELETE FROM users")

	router := SetupFullRouter()

	// 2. ACTORS SETUP (Direct DB insertion for speed)
	seller := models.User{Email: "farmer@test.com", Password: "hashed", IsSeller: true, Document: "111", FullName: "Farmer Joe"}
	buyer := models.User{Email: "hungry@test.com", Password: "hashed", IsSeller: false, Document: "222", FullName: "Hungry Bob"}
	database.DB.Create(&seller)
	database.DB.Create(&buyer)

	// Generate Tokens
	sellerToken, _ := auth.GenerateToken(seller.ID, true)
	buyerToken, _ := auth.GenerateToken(buyer.ID, false)

	// ---------------------------------------------------------
	// SCENARIO START
	// ---------------------------------------------------------

	// STEP 1: Seller creates a "Veggie Box"
	t.Log("Step 1: Seller creates Basket")
	basketPayload := `{"name": "Veggie Box", "description": "Fresh", "price": 100.0, "interval": "weekly", "seller_id": ` + fmt.Sprint(seller.ID) + `}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/baskets", bytes.NewBufferString(basketPayload))
	req.Header.Set("Authorization", sellerToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var createdBasket models.Basket
	json.Unmarshal(w.Body.Bytes(), &createdBasket)
	assert.NotZero(t, createdBasket.ID)

	// STEP 2: Buyer searches for "Veggie"
	t.Log("Step 2: Buyer searches for Baskets")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/baskets?search=Veggie", nil)
	req.Header.Set("Authorization", buyerToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Veggie Box")

	// STEP 3: Buyer Subscribes
	t.Log("Step 3: Buyer Subscribes")
	subPayload := fmt.Sprintf(`{"user_id": %d, "basket_id": %d}`, buyer.ID, createdBasket.ID)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/subscriptions", bytes.NewBufferString(subPayload))
	req.Header.Set("Authorization", buyerToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	// STEP 4: TIME PASSES (Trigger Worker)
	t.Log("Step 4: Worker runs (Generating Orders)")
	// Manually trigger the worker logic to process the new subscription immediately
	worker.ProcessSubscriptions()

	// STEP 5: Buyer checks their Orders
	t.Log("Step 5: Buyer checks Orders")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/orders", nil)
	req.Header.Set("Authorization", buyerToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var buyerOrders []models.Order
	json.Unmarshal(w.Body.Bytes(), &buyerOrders)
	assert.Len(t, buyerOrders, 1)
	assert.Equal(t, 100.0, buyerOrders[0].AmountPaid)
	orderID := buyerOrders[0].ID

	// STEP 6: Seller fulfills the Order
	t.Log("Step 6: Seller marks Order as Shipped")
	updatePayload := `{"status": "shipped"}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/orders/%d", orderID), bytes.NewBufferString(updatePayload))
	req.Header.Set("Authorization", sellerToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Verify status in DB
	var finalOrder models.Order
	database.DB.First(&finalOrder, orderID)
	assert.Equal(t, "shipped", finalOrder.Status)
}

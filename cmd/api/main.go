package main

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/handlers"
	"hobby-loop/m/internal/models"
	"hobby-loop/m/internal/worker"
	"hobby-loop/m/internal/auth"

	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Connect to the database
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system environment variables")
	}
	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Basket{}, &models.Subscription{}, &models.Order{}, &models.Address{})
	// Start the background worker
	worker.Start()

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(RequestLogger(logger))

	router.GET("/health", func(c *gin.Context) {
		sqlDB, err := database.DB.DB()
		if err != nil {
			slog.Error("Database connection failed", "error", err)
			c.JSON(500, gin.H{"status": "Database connection error"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			slog.Error("Database ping failed", "error", err)
			c.JSON(500, gin.H{"status": "Database unreachable"})
			return
		}
		c.JSON(200, gin.H{"status": "alive", "database": "connected"})
	})

	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", handlers.Login)

	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/baskets", handlers.CreateBasket)
		protected.POST("/subscriptions", handlers.SubscribeToBasket)
		protected.GET("/orders", handlers.GetOrders)
	}

	// Start the server
	slog.Info("Starting server on :8080")
	router.Run(":8080")
}

func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		logger.Info("HTTP Request",
			"status", status,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"duration", duration.String(),
			"ip", c.ClientIP(),
			"request_id", requestID,
		)
	}
}


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("is_seller", claims["is_seller"])
		c.Next()
	}
}
package handlers

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBasket(c *gin.Context) {
	var input models.Basket

	// Validate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create in db
	if result := database.DB.Create(&input); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func ListBaskets(c *gin.Context) {
	var baskets []models.Basket

	query := database.DB.Where("active = ?", true)

	search := c.Query("search")
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if result := query.Find(&baskets); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, baskets)
}
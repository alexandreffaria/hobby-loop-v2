package handlers

import (
	"net/http"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"

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

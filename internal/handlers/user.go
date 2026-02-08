package handlers

import (
	"net/http"
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/klassmann/cpfcnpj"
)

func RegisterUser(c *gin.Context) {
	var input struct {
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required"`
		FullName     string `json:"full_name" binding:"required"`
		DocumentType string `json:"document_type" binding:"required"`
		Document     string `json:"document" binding:"required"`
		IsSeller     bool   `json:"is_seller"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cleanDocument := cpfcnpj.Clean(input.Document)
	isValid := false

	if input.DocumentType == "CPF" {
		isValid = cpfcnpj.ValidateCPF(cleanDocument)
	} else if input.DocumentType == "CNPJ" {
		isValid = cpfcnpj.ValidateCNPJ(cleanDocument)
	} 

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document"})
		return
	}
	
	// Validate document
	user := models.User{
		Email: input.Email,
		Password: input.Password,
		FullName: input.FullName,
		DocumentType: input.DocumentType,
		Document: cleanDocument,
		IsSeller: input.IsSeller,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}



	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
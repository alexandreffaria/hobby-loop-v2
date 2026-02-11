package handlers

import (
	"hobby-loop/m/internal/database"
	"hobby-loop/m/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klassmann/cpfcnpj"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	FullName     string `json:"full_name" binding:"required"`
	DocumentType string `json:"document_type" binding:"required"`
	Document     string `json:"document" binding:"required"`
	IsSeller     bool   `json:"is_seller"`

	Address struct {
		Street       string `json:"street"`
		Number       string `json:"number"`
		Complement   string `json:"complement"`
		Neighborhood string `json:"neighborhood"`
		City         string `json:"city"`
		State        string `json:"state"`
		ZipCode      string `json:"zip_code"`
	} `json:"address"`
}

func RegisterUser(c *gin.Context) {
	var input RegisterInput

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Validate document
	user := models.User{
		Email:        input.Email,
		Password:     string(hashedPassword),
		FullName:     input.FullName,
		DocumentType: input.DocumentType,
		Document:     cleanDocument,
		IsSeller:     input.IsSeller,
		Addresses: []models.Address{
			{
				Street:       input.Address.Street,
				Number:       input.Address.Number,
				Complement:   input.Address.Complement,
				Neighborhood: input.Address.Neighborhood,
				City:         input.Address.City,
				State:        input.Address.State,
				ZipCode:      input.Address.ZipCode,
			},
		},
	}
	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

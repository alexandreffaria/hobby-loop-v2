package models

import (
	"gorm.io/gorm"
	"time"
)

type Address struct {
	gorm.Model
	UserID       uint   `json:"user_id"`
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zip_code"`
}

// Usuário
type User struct {
	gorm.Model
	Email        string    `json:"email" gorm:"unique"`
	Password     string    `json:"-"`
	FullName     string    `json:"full_name"`
	DocumentType string    `json:"document_type"`
	Document     string    `json:"document" goorm:"unique;not null"`
	IsSeller     bool      `json:"is_seller"`
	Addresses    []Address `json:"addresses"`
}

// Basket
type Basket struct {
	gorm.Model
	SellerID    uint    `json:"seller_id"`
	Name        string  `json:name`
	Description string  `json:description`
	Price       float64 `json:price`
	Interval    string  `json:"interval"`
	Active      bool    `json:"active" gorm:"default:true"`
}

// Subscription
type Subscription struct {
	gorm.Model
	UserID           uint      `json:"user_id"`
	BasketID         uint      `json:"basket_id"`
	Status           string    `json:"status"`
	NextDeliveryDate time.Time `json:"next_delivery_date"`
}

// Order
type Order struct {
	gorm.Model
	SubscriptionID uint    `json:"subscription_id"`
	AmountPaid     float64 `json:"amount_paid"`
	Status         string  `json:"status"`

	InvoiceKey string `json:"invoice_key"` // Chave de acesso da nota fiscal
	InvoiceURL string `json:"invoice_url"` // URL para consulta da nota fiscal
}

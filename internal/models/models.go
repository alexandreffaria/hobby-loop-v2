package models

import (
	"time"
	"gorm.io/gorm"
)

// Usuário
type User struct {
	gorm.Model
	Email string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	IsSeller bool `json:"is_seller"`
}

// Basket
type Basket struct {
	gorm.Model
	SellerID uint `json:"seller_id"`
	Name string `json:name`
	Description string `json:description`
	Price float64 `json:price`
	Interval string `json:"interval"`
}

// Subscription
type Subscription struct{
	gorm.Model
	UserID uint `json:"user_id"`
	BasketID uint `json:"basket_id"`
	Status string `json:"status"`
	NextDeliveryDate time.Time `json:"next_delivery_date"`
}

type Order struct {
	gorm.Model
	SubscriptionID uint `json:"subscription_id"` 
	AmountPaid float64 `json:"amount_paid"`
	Status string `json:"status"`
}
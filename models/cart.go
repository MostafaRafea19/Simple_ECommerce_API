package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	CustomerID uint       `json:"customer_id"`
	Customer   Customer   `gorm:"foreignKey:customer_id"`
	Products   []Product  `gorm:"many2many:cart_items;"`
	CartItems  []CartItem `gorm:"foreignKey:cart_id"`
	TotalPrice float64    `json:"total_price"`
	IsActive   bool       `json:"is_active"`
	Order      *Order     `gorm:"constraint:OnDelete:CASCADE;"`
}

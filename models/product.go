package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SellerId    uint    `json:"seller_id"`
	Seller      *Seller `gorm:"foreignKey:seller_id"`
}

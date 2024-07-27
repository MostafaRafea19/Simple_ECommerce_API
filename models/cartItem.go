package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID    uint     `json:"cart_id"`
	Cart      *Cart    `gorm:"foreignKey:cart_id;constraint:OnDelete:CASCADE;"`
	ProductID uint     `json:"product_id"`
	Product   *Product `gorm:"foreignKey:product_id;constraint:OnDelete:CASCADE;"`
	Quantity  int      `json:"quantity"`
	Status    string   `json:"status" gorm:"default:'pending'"`
}

package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	OrderID     uint    `json:"order_id"`
	Order       *Order  `gorm:"foreignKey:order_id;constraint:OnDelete:CASCADE;"`
	TotalAmount float64 `json:"total_amount"`
	Paid        bool    `json:"paid"`
}

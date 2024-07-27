package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	CartID       uint          `json:"cart_id"`
	Cart         *Cart         `gorm:"foreignKey:cart_id"`
	TotalAmount  float64       `json:"total_amount"`
	OrderedDate  time.Time     `json:"ordered_date"`
	Status       string        `json:"status"`
	Payment      *Payment      `gorm:"constraint:OnDelete:CASCADE;"`
	ShippingInfo *ShippingInfo `gorm:"constraint:OnDelete:CASCADE;"`
}

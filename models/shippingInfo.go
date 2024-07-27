package models

import (
	"gorm.io/gorm"
)

type ShippingInfo struct {
	gorm.Model
	OrderID uint   `json:"order_id"`
	Order   *Order `gorm:"foreignKey:order_id;constraint:OnDelete:CASCADE;"`
	Address string `json:"address"`
}

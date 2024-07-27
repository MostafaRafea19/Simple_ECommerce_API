package models

type Customer struct {
	User
	Address string `json:"address"`
	Carts   []Cart `gorm:"foreignKey:customer_id"`
}

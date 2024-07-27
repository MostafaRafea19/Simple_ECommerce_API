package models

type Seller struct {
	User
	StoreName string    `json:"store_name"`
	Products  []Product `gorm:"foreignKey:seller_id"`
}

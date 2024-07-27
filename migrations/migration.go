package migrations

import (
	"api/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Customer{},
		&models.Seller{},
		&models.Admin{},
		&models.Product{},
		&models.Cart{},
		&models.Order{},
		&models.Payment{},
		&models.ShippingInfo{},
		&models.CartItem{},
	)
}

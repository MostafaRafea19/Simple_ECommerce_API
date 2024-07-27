package controllers

import (
	"api/database"
	"api/models"
	"api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
)

func AddItemToCart(c *gin.Context) {
	customerID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	var product models.Product
	if err := database.GetDB().First(&product, input.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Product not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var cart models.Cart
	if err := database.GetDB().Where("customer_id = ? AND is_active = ?", customerID, true).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = models.Cart{
				CustomerID: customerID.(uint),
				IsActive:   true,
			}
			database.GetDB().Create(&cart)
		} else {
			utils.InternalServerErrorJSON(c, err.Error())
			return
		}
	}

	var cartItem models.CartItem
	if err := database.GetDB().Where("product_id = ? AND cart_id = ?", product.ID, cart.ID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: input.ProductID,
				Quantity:  input.Quantity,
			}
			database.GetDB().Create(&cartItem)
		} else {
			utils.InternalServerErrorJSON(c, err.Error())
			return
		}
	} else {
		cartItem.Quantity += input.Quantity
		database.GetDB().Save(&cartItem)
	}

	cart.TotalPrice += float64(input.Quantity) * product.Price
	database.GetDB().Save(&cart)

	utils.JSONResponse(c, http.StatusCreated, gin.H{"message": "Product successfully added to cart"})
}

func UpdateCartItem(c *gin.Context) {
	customerID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	var cart models.Cart
	if err := database.GetDB().Where("customer_id = ? AND is_active = ?", customerID, true).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Cart not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var cartItem models.CartItem
	if err := database.GetDB().First(&cartItem, c.Param("cartItemId")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Cart item not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	cartItem.Quantity = input.Quantity
	if err := database.GetDB().Save(&cartItem).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, cartItem)
}

func DeleteCartItem(c *gin.Context) {
	var cartItem models.Customer
	if err := database.GetDB().First(&cartItem, c.Param("cartItemId")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Cart item not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if err := database.GetDB().Unscoped().Delete(&cartItem).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Cart item deleted successfully"})
}

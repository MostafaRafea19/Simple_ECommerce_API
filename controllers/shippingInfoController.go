package controllers

import (
	"api/database"
	"api/models"
	"api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetOrderShippingInfo(c *gin.Context) {
	var shippingInfo models.ShippingInfo
	if err := database.GetDB().Where("order_id = ?", c.Param("id")).First(&shippingInfo).Error; err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Order shipping info not found",
		})
		return
	}

	c.JSON(http.StatusOK, shippingInfo)
}

func GetSellerOrderShippingInfo(c *gin.Context) {
	sellerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	var order models.Order
	if err := database.GetDB().Preload("ShippingInfo").Preload("Cart.CartItems", "product_id IN (?)", database.GetDB().Model(models.Product{})).Where("seller_id = ?", sellerId).First(&order, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Order not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var filteredCartItems []models.CartItem
	for _, item := range order.Cart.CartItems {
		if item.Product.SellerId == sellerId {
			filteredCartItems = append(filteredCartItems, item)
		}
	}

	if len(filteredCartItems) == 0 {
		utils.BadRequestErrorJson(c, "Order doesn't belong to this seller")
		return
	}

	c.JSON(http.StatusOK, order.ShippingInfo)
}

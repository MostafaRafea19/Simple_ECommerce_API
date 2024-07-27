package controllers

import (
	"api/database"
	"api/models"
	"api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOrderPayment(c *gin.Context) {
	orderID := c.Param("id")
	var payment models.Payment
	if err := database.GetDB().Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Payment not found",
		})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func UpdatePaymentStatus(c *gin.Context) {
	orderID := c.Param("id")
	var payment models.Payment
	if err := database.GetDB().Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		utils.ErrorJSON(c, http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "OrderPayment not found",
		})
		return
	}

	var input struct {
		Paid bool `json:"paid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorJSON(c, http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "Invalid input",
		})
		return
	}

	payment.Paid = input.Paid

	if err := database.GetDB().Save(&payment).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, payment)
}

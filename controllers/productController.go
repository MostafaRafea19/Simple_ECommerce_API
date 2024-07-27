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

func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := database.GetDB().Find(&products).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if len(products) == 0 {
		utils.NotFoundRequestErrorJson(c, "No products found")
		return
	}

	utils.JSONResponse(c, http.StatusOK, products)
}

func GetSellerProducts(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Seller is not authenticated")
		return
	}

	var products []models.Product
	if err := database.GetDB().Where("seller_id = ?", sellerID).Find(&products).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if len(products) == 0 {
		utils.NotFoundRequestErrorJson(c, "No products found")
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {
	var product models.Product
	if err := database.GetDB().First(&product, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Product not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, product)
}

func GetSellerProductDetails(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Seller is not authenticated")
		return
	}

	var product models.Product
	if err := database.GetDB().First(&product, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Product not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if product.SellerId != sellerID {
		utils.UnauthorizedRequestJson(c, "Product does not belong to seller")
		return
	}

	utils.JSONResponse(c, http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var productInput struct {
		Name        string  `json:"name" binding:"required"`
		SKU         string  `json:"sku" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&productInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	var product models.Product
	err := database.GetDB().Where("sku = ?", productInput.SKU).First(&product).Error
	if err == nil {
		utils.ConflictRequestErrorJson(c, "Product already exists with the same sku")
		return
	}

	sellerID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Seller is not authenticated")
		return
	}

	newProduct := models.Product{
		Name:        productInput.Name,
		SKU:         productInput.SKU,
		Description: productInput.Description,
		Price:       productInput.Price,
		SellerId:    sellerID.(uint),
	}

	if err := database.GetDB().Create(&newProduct).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, newProduct)
}

func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	var existingProduct models.Product
	if err := database.GetDB().First(&existingProduct, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Product not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var productInput struct {
		SKU         string  `json:"sku" binding:"required"`
		Name        string  `json:"name" binding:"omitempty"`
		Description string  `json:"description" binding:"omitempty"`
		Price       float64 `json:"price" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&productInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	if productInput.SKU != "" {
		var product models.Product
		err := database.GetDB().Where("sku = ?", productInput.SKU).First(&product).Error
		if err == nil {
			utils.ConflictRequestErrorJson(c, "Product already exists with the same sku")
			return
		}
		existingProduct.SKU = productInput.SKU
	}
	if productInput.Name != "" {
		existingProduct.Name = productInput.Name
	}
	if productInput.Description != "" {
		existingProduct.Description = productInput.Description
	}
	if productInput.Price != 0 {
		existingProduct.Price = productInput.Price
	}

	if err := database.GetDB().Save(&existingProduct).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, existingProduct)
}

func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	var existingProduct models.Product
	if err := database.GetDB().First(&existingProduct, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Product not found")
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if err := database.GetDB().Unscoped().Delete(&existingProduct).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

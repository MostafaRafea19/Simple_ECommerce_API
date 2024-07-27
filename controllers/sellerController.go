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

func GetSellers(c *gin.Context) {
	var sellers []models.Seller
	if err := database.GetDB().Find(&sellers).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if len(sellers) == 0 {
		utils.NotFoundRequestErrorJson(c, "No sellers are found")
	}

	utils.JSONResponse(c, http.StatusOK, sellers)
}

func GetSeller(c *gin.Context) {
	var seller models.Seller
	if err := database.GetDB().First(&seller, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "seller not found")
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, seller)
}

func GetSellerProfile(c *gin.Context) {
	sellerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Seller is not authenticated")
		return
	}

	c.AddParam("id", sellerId.(string))
	GetCustomer(c)
}

func CreateSeller(c *gin.Context) {
	var sellerInput struct {
		Email     string `json:"email" binding:"required,email"`
		Name      string `json:"name" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		Password  string `json:"password" binding:"required,min=6"`
		StoreName string `json:"store_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&sellerInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	var seller models.Seller
	err := database.GetDB().Where("email = ?", sellerInput.Email).First(&seller).Error
	if err == nil {
		utils.ConflictRequestErrorJson(c, "Seller already exists with the same email")
		return
	}

	err = database.GetDB().Where("phone = ?", sellerInput.Phone).First(&seller).Error
	if err == nil {
		utils.ConflictRequestErrorJson(c, "Seller already exists with the same Phone")
		return
	}

	hashedPassword, err := utils.HashPassword(sellerInput.Password)
	if err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	newSeller := models.Seller{
		User: models.User{
			Email:    sellerInput.Email,
			Phone:    sellerInput.Phone,
			Password: hashedPassword,
			Name:     sellerInput.Name,
		},
		StoreName: sellerInput.StoreName,
	}

	if err := database.GetDB().Create(&newSeller).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, newSeller)
}

func UpdateSeller(c *gin.Context) {
	sellerID := c.Param("id")

	var existingSeller models.Seller
	if err := database.GetDB().First(&existingSeller, sellerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "seller not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var sellerInput struct {
		Email     string `json:"email" binding:"omitempty,email"`
		Name      string `json:"name" binding:"omitempty"`
		Phone     string `json:"phone" binding:"omitempty"`
		Password  string `json:"password" binding:"omitempty,min=6"`
		StoreName string `json:"store_name" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&sellerInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	if sellerInput.Email != "" {
		var seller models.Seller
		err := database.GetDB().Where("email = ?", sellerInput.Email).First(&seller).Error
		if err == nil {
			utils.ConflictRequestErrorJson(c, "Seller already exists with the same email")
			return
		}
		existingSeller.Email = sellerInput.Email
	}
	if sellerInput.Phone != "" {
		var seller models.Seller
		err := database.GetDB().Where("phone = ?", sellerInput.Phone).First(&seller).Error
		if err == nil {
			utils.ConflictRequestErrorJson(c, "Seller already exists with the same phone")
			return
		}
		existingSeller.Phone = sellerInput.Phone
	}
	if sellerInput.Password != "" {
		hashedPassword, err := utils.HashPassword(sellerInput.Password)
		if err != nil {
			utils.InternalServerErrorJSON(c, err.Error())
			return
		}
		existingSeller.Password = hashedPassword
	}
	if sellerInput.Name != "" {
		existingSeller.Name = sellerInput.Name
	}
	if sellerInput.StoreName != "" {
		existingSeller.StoreName = sellerInput.StoreName
	}

	if err := database.GetDB().Save(&existingSeller).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, existingSeller)
}

func UpdateSellerProfile(c *gin.Context) {
	sellerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Seller is not authenticated")
		return
	}

	c.AddParam("id", sellerId.(string))
	UpdateSeller(c)
}

func DeleteSeller(c *gin.Context) {
	customerID := c.Param("id")

	var existingSeller models.Seller
	if err := database.GetDB().First(&existingSeller, customerID).Error; err != nil {
		utils.NotFoundRequestErrorJson(c, "seller not found")
		return
	}

	if err := database.GetDB().Unscoped().Delete(&existingSeller).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Seller deleted successfully"})
}

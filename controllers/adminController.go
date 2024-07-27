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

func CreateAdmin(c *gin.Context) {
	var adminInput struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&adminInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	var admin models.Admin
	err := database.GetDB().Where("email = ?", adminInput.Email).First(&admin).Error
	if err == nil {
		utils.ConflictRequestErrorJson(c, "Seller already exists with the same email")
		return
	}

	err = database.GetDB().Where("phone = ?", adminInput.Phone).First(&admin).Error
	if err == nil {
		utils.ConflictRequestErrorJson(c, "Seller already exists with the same Phone")
		return
	}

	err = database.GetDB().Where("username = ?", adminInput.Username).First(&admin).Error
	if err == nil {
		utils.ConflictRequestErrorJson(c, "Admin already exists with the same username")
		return
	}

	hashedPassword, err := utils.HashPassword(adminInput.Password)
	if err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	newAdmin := models.Admin{
		User: models.User{
			Email:    adminInput.Email,
			Phone:    adminInput.Phone,
			Password: hashedPassword,
			Name:     adminInput.Name,
		},
		Username: adminInput.Username,
	}

	if err := database.GetDB().Create(&newAdmin).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, newAdmin)
}

func GetAdmin(c *gin.Context) {
	var admin models.Admin
	if err := database.GetDB().First(&admin, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Admin not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, admin)
}

func GetAdminProfile(c *gin.Context) {
	adminId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Admin is not authenticated")
		return
	}

	c.AddParam("id", adminId.(string))
	GetAdmin(c)
}

func UpdateAdminProfile(c *gin.Context) {
	adminId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Admin is not authenticated")
		return
	}

	c.AddParam("id", adminId.(string))
	UpdateAdmin(c)
}

func UpdateAdmin(c *gin.Context) {
	adminID := c.Param("id")

	var existingAdmin models.Seller
	if err := database.GetDB().First(&existingAdmin, adminID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "admin not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var adminInput struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Name     string `json:"name" binding:"omitempty"`
		Phone    string `json:"phone" binding:"omitempty"`
		Password string `json:"password" binding:"omitempty,min=6"`
		Username string `json:"username" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&adminInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	if adminInput.Email != "" {
		var admin models.Admin
		err := database.GetDB().Where("email = ?", adminInput.Email).First(&admin).Error
		if err == nil {
			utils.ConflictRequestErrorJson(c, "Admin already exists with the same email")
			return
		}
		existingAdmin.Email = adminInput.Email
	}
	if adminInput.Phone != "" {
		var admin models.Admin
		err := database.GetDB().Where("email = ?", adminInput.Phone).First(&admin).Error
		if err == nil {
			utils.ConflictRequestErrorJson(c, "Admin already exists with the same email")
			return
		}
		existingAdmin.Phone = adminInput.Phone
	}
	if adminInput.Password != "" {
		hashedPassword, err := utils.HashPassword(adminInput.Password)
		if err != nil {
			utils.InternalServerErrorJSON(c, err.Error())
			return
		}
		existingAdmin.Password = hashedPassword
	}
	if adminInput.Name != "" {
		existingAdmin.Name = adminInput.Name
	}
	if adminInput.Username != "" {
		var admin models.Admin
		err := database.GetDB().Where("username = ?", adminInput.Username).First(&admin).Error
		if err == nil {
			utils.ConflictRequestErrorJson(c, "Admin already exists with the same username")
			return
		}
		existingAdmin.StoreName = adminInput.Username
	}

	if err := database.GetDB().Save(&existingAdmin).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, existingAdmin)
}

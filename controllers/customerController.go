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

func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	if err := database.GetDB().Find(&customers).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if len(customers) == 0 {
		utils.NotFoundRequestErrorJson(c, "No customers found")
		return
	}

	utils.JSONResponse(c, http.StatusOK, customers)
}

func GetCustomerProfile(c *gin.Context) {
	customerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	c.AddParam("id", customerId.(string))
	GetCustomer(c)
}

func GetCustomer(c *gin.Context) {
	var customer models.Customer
	if err := database.GetDB().First(&customer, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Customer not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, customer)
}

func CreateCustomer(c *gin.Context) {
	var customerInput struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Address  string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&customerInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	var customer models.Customer
	err := database.GetDB().Where("email = ?", customerInput.Email).First(&customer).Error
	if err == nil {
		utils.ConflictRequestErrorJson(c, "Customer already exists with the same email")
		return
	}

	err = database.GetDB().Where("phone = ?", customerInput.Phone).First(&customer).Error
	if err == nil {
		utils.ConflictRequestErrorJson(c, "Customer already exists with the same Phone")
		return
	}

	hashedPassword, err := utils.HashPassword(customerInput.Password)
	if err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	newCustomer := models.Customer{
		User: models.User{
			Email:    customerInput.Email,
			Phone:    customerInput.Phone,
			Password: hashedPassword,
			Name:     customerInput.Name,
		},
		Address: customerInput.Address,
	}

	if err := database.GetDB().Create(&newCustomer).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, newCustomer)
}

func UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")

	var existingCustomer models.Customer
	if err := database.GetDB().First(&existingCustomer, customerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Customer not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var customerInput struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Phone    string `json:"phone" binding:"omitempty"`
		Password string `json:"password" binding:"omitempty,min=6"`
		Name     string `json:"name" binding:"omitempty"`
		Address  string `json:"address" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&customerInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	if customerInput.Email != "" {
		var customer models.Customer
		err := database.GetDB().Where("email = ?", customerInput.Email).First(&customer).Error
		if err == nil {
			utils.ConflictRequestErrorJson(c, "Customer already exists with the same email")
			return
		}
		existingCustomer.Email = customerInput.Email
	}
	if customerInput.Phone != "" {
		var customer models.Customer
		err := database.GetDB().Where("phone = ?", customerInput.Phone).First(&customer).Error
		if err == nil {
			utils.ConflictRequestErrorJson(c, "Customer already exists with the same Phone")
			return
		}
		existingCustomer.Phone = customerInput.Phone
	}
	if customerInput.Password != "" {
		hashedPassword, err := utils.HashPassword(customerInput.Password)
		if err != nil {
			utils.InternalServerErrorJSON(c, err.Error())
			return
		}
		existingCustomer.Password = hashedPassword
	}
	if customerInput.Name != "" {
		existingCustomer.Name = customerInput.Name
	}
	if customerInput.Address != "" {
		existingCustomer.Address = customerInput.Address
	}

	if err := database.GetDB().Save(&existingCustomer).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, existingCustomer)
}

func DeleteCustomer(c *gin.Context) {
	customerID := c.Param("id")

	var existingCustomer models.Customer
	if err := database.GetDB().First(&existingCustomer, customerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Customer not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if err := database.GetDB().Unscoped().Delete(&existingCustomer).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

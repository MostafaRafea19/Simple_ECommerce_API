package controllers

import (
	"api/models"
	"api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
)

func Login(c *gin.Context) {
	var loginInput struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	service := os.Getenv("SERVICE")
	var userType string
	var hashedPassword string
	var userID uint

	switch service {
	case "customers":
		var customer models.Customer
		if utils.CheckEmailNotFoundError(c, loginInput.Email, "customer", &customer) {
			return
		}
		userType = "customer"
		hashedPassword = customer.Password
		userID = customer.ID

	case "sellers":
		var seller models.Seller
		if utils.CheckEmailNotFoundError(c, loginInput.Email, "seller", &seller) {
			return
		}
		userType = "seller"
		hashedPassword = seller.Password
		userID = seller.ID

	case "admins":
		var admin models.Admin
		if utils.CheckEmailNotFoundError(c, loginInput.Email, "admin", &admin) {
			return
		}
		userType = "admin"
		hashedPassword = admin.Password
		userID = admin.ID
	}

	if !utils.CheckPasswordHash(loginInput.Password, hashedPassword) {
		utils.UnauthorizedRequestJson(c, "Wrong password.")
		return
	}

	token, err := utils.GenerateJWT(userID, userType)
	if err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

package utils

import (
	"api/database"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func HandleValidationErrors(verr validator.ValidationErrors) []ValidationError {
	var errs []ValidationError

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, ValidationError{Field: f.Field(), Reason: err})
	}

	return errs
}

func CheckEmailNotFoundError(c *gin.Context, email string, resource string, entity interface{}) bool {
	err := database.GetDB().Where("email = ?", email).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			UnauthorizedRequestJson(c, "No"+resource+"found with this email")
		}
		InternalServerErrorJSON(c, err.Error())
	}

	return err != nil
}

func CheckUniqueValidationError(c *gin.Context, attribute string, value string, resource string, entity interface{}) {
	err := database.GetDB().Where(attribute+" = ?", value).First(&entity).Error
	if err == nil {
		ConflictRequestErrorJson(c, resource+"already exists with the same"+attribute)
		return
	}
}

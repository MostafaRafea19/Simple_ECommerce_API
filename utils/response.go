package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func JSONResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

func ErrorJSON(c *gin.Context, statusCode int, response interface{}) {
	c.JSON(statusCode, response)
}

func InternalServerErrorJSON(c *gin.Context, errorMessage string) {
	ErrorJSON(c, http.StatusInternalServerError, gin.H{"error": errorMessage, "message": "Internal Server Error"})
}

func BadRequestErrorJson(c *gin.Context, error string) {
	ErrorJSON(c, http.StatusBadRequest, gin.H{"message": error})
}

func UnauthorizedRequestJson(c *gin.Context, error string) {
	ErrorJSON(c, http.StatusUnauthorized, gin.H{"message": error})
}

func NotFoundRequestErrorJson(c *gin.Context, error string) {
	ErrorJSON(c, http.StatusNotFound, gin.H{"message": error})
}

func ConflictRequestErrorJson(c *gin.Context, error string) {
	ErrorJSON(c, http.StatusConflict, gin.H{"message": error})
}

func ValidationErrorJson(c *gin.Context, verr validator.ValidationErrors) {
	ErrorJSON(c, http.StatusUnprocessableEntity, gin.H{"errors": HandleValidationErrors(verr)})
}

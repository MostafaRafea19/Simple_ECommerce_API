package middlewares

import (
	"api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorJSON(c, http.StatusUnauthorized, gin.H{
				"error": "User is not authorized",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.ErrorJSON(c, http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "admin" {
			utils.ErrorJSON(c, http.StatusForbidden, gin.H{
				"error": "Access forbidden: admins only",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func SellerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "seller" {
			utils.ErrorJSON(c, http.StatusForbidden, gin.H{
				"error": "Access forbidden: sellers only",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CustomerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "customer" {
			utils.ErrorJSON(c, http.StatusForbidden, gin.H{
				"error": "Access forbidden: customers only",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

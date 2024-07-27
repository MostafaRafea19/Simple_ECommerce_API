package routes

import (
	"api/controllers"
	"api/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter(group string) *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/login", controllers.Login)
	}

	switch group {
	case "customers":
		CustomerRoutes(apiGroup)

	case "sellers":
		SellerRoutes(apiGroup)

	case "admins":
		AdminRoutes(apiGroup)
	}

	return router
}

func CustomerRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	apiGroup.POST("/customers", controllers.CreateCustomer)

	customerGroup := apiGroup.Group("/customers")
	customerGroup.Use(middlewares.AuthMiddleware(), middlewares.CustomerMiddleware())
	{
		customerGroup.GET("/profile", controllers.GetCustomerProfile)

		productGroup := customerGroup.Group("/products")
		{
			productGroup.GET("/", controllers.GetProducts)
			productGroup.GET("/:id", controllers.GetProduct)
		}

		cartGroup := customerGroup.Group("/cart")
		{
			cartGroup.POST("/", controllers.AddItemToCart)
			cartGroup.PATCH("cart-items/:cartItemId", controllers.UpdateCartItem)
			cartGroup.DELETE("cart-items/:cartItemId", controllers.DeleteCartItem)
		}

		orderGroup := customerGroup.Group("/orders")
		{
			orderGroup.POST("/", controllers.PlaceOrder)
			orderGroup.GET("/", controllers.GetCustomerOrders)
			orderGroup.GET("/:id", controllers.GetCustomerOrder)
		}
	}

	return customerGroup
}

func SellerRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	apiGroup.POST("/sellers", controllers.CreateSeller)

	sellerGroup := apiGroup.Group("/sellers")
	sellerGroup.Use(middlewares.AuthMiddleware(), middlewares.SellerMiddleware())
	{

		sellerGroup.GET("profile/", controllers.GetSellerProfile)
		sellerGroup.PATCH("profile/", controllers.UpdateSellerProfile)

		productGroup := sellerGroup.Group("/products")
		{
			productGroup.POST("/", controllers.CreateProduct)
			productGroup.GET("/", controllers.GetSellerProducts)
			productGroup.GET("/:id", controllers.GetSellerProductDetails)
			productGroup.PATCH("/:id", controllers.UpdateProduct)
			productGroup.DELETE("/:id", controllers.DeleteProduct)
		}

		orderGroup := sellerGroup.Group("/orders")
		{
			orderGroup.GET("/", controllers.GetSellerOrders)
			orderGroup.GET("/:id", controllers.GetSellerOrderDetails)
			orderGroup.PATCH("/:id/:itemId", controllers.UpdateOrderItemStatus)
			orderGroup.DELETE("/:id", controllers.DeleteOrder)
			orderGroup.GET("/:id/shipping_info", controllers.GetSellerOrderShippingInfo)
		}
	}

	return sellerGroup
}

func AdminRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	adminGroup := apiGroup.Group("admins")
	adminGroup.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		adminGroup.POST("/", controllers.CreateAdmin)
		adminGroup.GET("/", controllers.GetAdminProfile)
		adminGroup.PATCH("/", controllers.UpdateAdminProfile)

		customerGroup := adminGroup.Group("/customers")
		{
			customerGroup.GET("/", controllers.GetCustomers)
			customerGroup.GET("/:id", controllers.GetCustomer)
			customerGroup.PATCH("/:id", controllers.UpdateCustomer)
			customerGroup.DELETE("/:id", controllers.DeleteCustomer)
		}

		sellerGroup := adminGroup.Group("/sellers")
		{
			sellerGroup.GET("/", controllers.GetSellers)
			sellerGroup.GET("/:id", controllers.GetSeller)
			sellerGroup.PATCH("/:id", controllers.UpdateSeller)
			sellerGroup.DELETE("/:id", controllers.DeleteSeller)
		}

		productGroup := adminGroup.Group("/products")
		{
			productGroup.GET("/", controllers.GetProducts)
			productGroup.GET("/:id", controllers.GetProduct)
		}

		orderGroup := adminGroup.Group("/orders")
		{
			orderGroup.GET("/", controllers.GetOrders)
			orderGroup.GET("/:id", controllers.GetOrder)
			orderGroup.GET("/:id/shipping_info", controllers.GetOrderShippingInfo)
			orderGroup.GET("/:id/payment", controllers.GetOrderPayment)
			orderGroup.PATCH("/:id/payment", controllers.UpdatePaymentStatus)
		}
	}

	return adminGroup
}

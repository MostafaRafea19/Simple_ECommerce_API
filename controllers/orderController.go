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
	"time"
)

func GetOrders(c *gin.Context) {
	var orders []models.Order
	if err := database.GetDB().Find(&orders).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if len(orders) == 0 {
		utils.NotFoundRequestErrorJson(c, "No orders found")
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetOrder(c *gin.Context) {
	var order models.Order
	if err := database.GetDB().Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Order not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetCustomerOrders(c *gin.Context) {
	customerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
	}

	var orders []models.Order
	if err := database.GetDB().Where("cart_id IN (SELECT id FROM carts WHERE customer_id = ?)", customerId).Find(&orders).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if len(orders) == 0 {
		utils.NotFoundRequestErrorJson(c, "No orders found for this customer")
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetSellerOrders(c *gin.Context) {
	sellerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	var orders []models.Order
	subQuery := database.GetDB().Model(&models.CartItem{}).Select("cart_id").Joins("JOIN products ON products.id = cart_items.product_id").Where("products.seller_id = ?", sellerId)

	if err := database.GetDB().Preload("Cart.CartItems", "product_id IN (?)", database.GetDB().Select("id").Model(&models.Product{}).Where("seller_id = ?", sellerId)).Where("cart_id IN (?)", subQuery).Find(&orders).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if len(orders) == 0 {
		utils.NotFoundRequestErrorJson(c, "No orders found for this seller")
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetCustomerOrder(c *gin.Context) {
	customerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	var order models.Order
	if err := database.GetDB().First(&order, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Order not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var cart models.Cart
	if err := database.GetDB().First(&cart, order.CartID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Cart not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if cart.CustomerID != customerId {
		utils.NotFoundRequestErrorJson(c, "Order doesn't belong to this customer")
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetSellerOrderDetails(c *gin.Context) {
	sellerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	var order models.Order
	if err := database.GetDB().Preload("Cart.CartItems", "product_id IN (?)", database.GetDB().Model(models.Product{})).Where("seller_id = ?", sellerId).First(&order, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Order not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var filteredCartItems []models.CartItem
	for _, item := range order.Cart.CartItems {
		if item.Product.SellerId == sellerId {
			filteredCartItems = append(filteredCartItems, item)
		}
	}

	if len(filteredCartItems) == 0 {
		utils.BadRequestErrorJson(c, "Order doesn't belong to this seller")
		return
	}

	order.Cart.CartItems = filteredCartItems

	c.JSON(http.StatusOK, order)
}

func PlaceOrder(c *gin.Context) {
	customerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	var input struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	var cart models.Cart
	if err := database.GetDB().Where("customer_id = ? AND is_active = ?", customerId, true).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "There is no active cart for this customer.")
			return
		}
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if err := database.GetDB().Where("cart_id = ?", cart.ID).First(&models.Order{}).Error; err == nil {
		utils.ConflictRequestErrorJson(c, "There is an already placed order for this customer's cart")
		return
	}

	var cartItems []models.CartItem
	if err := database.GetDB().Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if len(cartItems) == 0 {
		utils.NotFoundRequestErrorJson(c, "No items found for this customer's cart.")
		return
	}

	totalAmount := 0.0
	for _, item := range cartItems {
		var product models.Product
		if err := database.GetDB().First(&product, item.ProductID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.NotFoundRequestErrorJson(c, "One of products in cart is not found")
				return
			}

			utils.InternalServerErrorJSON(c, err.Error())
			return
		}
		totalAmount += float64(item.Quantity) * product.Price
	}

	order := models.Order{
		CartID:      cart.ID,
		TotalAmount: totalAmount,
		OrderedDate: time.Now(),
		Status:      utils.StatusPending,
	}

	if err := database.GetDB().Create(&order).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	ShippingInfo := models.ShippingInfo{
		OrderID: order.ID,
		Address: input.Address,
	}

	if err := database.GetDB().Create(&ShippingInfo).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	payment := models.Payment{
		OrderID:     order.ID,
		TotalAmount: totalAmount,
		Paid:        false,
	}

	if err := database.GetDB().Create(&payment).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	cart.IsActive = false
	database.GetDB().Save(&cart)

	if err := database.GetDB().Preload("Payment").Preload("ShippingInfo").First(&order, order.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Order not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

func UpdateOrderItemStatus(c *gin.Context) {
	sellerId, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedRequestJson(c, "Customer is not authenticated")
		return
	}

	var existingOrder models.Order
	if err := database.GetDB().First(&existingOrder, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Order not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	var cartItem models.CartItem
	if err := database.GetDB().First(&cartItem, c.Param("itemId")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Order item not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if existingOrder.CartID != cartItem.CartID {
		utils.BadRequestErrorJson(c, "Order item doesn't belong to this order.")
		return
	}

	var product models.Product
	if err := database.GetDB().First(&product, cartItem.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Product from order item is not found.")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if sellerId != product.SellerId {
		utils.BadRequestErrorJson(c, "Order item doesn't belong to seller.")
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.ValidationErrorJson(c, verr)
			return
		}

		utils.BadRequestErrorJson(c, err.Error())
		return
	}

	switch input.Status {
	case utils.StatusPending, utils.StatusShipped, utils.StatusDelivered, utils.StatusCancelled:
		cartItem.Status = input.Status
	default:
		utils.BadRequestErrorJson(c, "Invalid status, it should be either "+utils.StatusPending+" or "+utils.StatusShipped+" or "+utils.StatusCancelled+" or "+utils.StatusDelivered)
		return
	}

	if err := database.GetDB().Save(&cartItem).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Order item status successfully updated."})
}

func DeleteOrder(c *gin.Context) {
	var existingOrder models.Order
	if err := database.GetDB().First(&existingOrder, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundRequestErrorJson(c, "Order not found")
			return
		}

		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	if err := database.GetDB().Unscoped().Delete(&existingOrder).Error; err != nil {
		utils.InternalServerErrorJSON(c, err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

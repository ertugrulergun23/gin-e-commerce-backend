package handlers

import (
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InputOrder represents the request body for updating an order status.
type InputOrder struct {
	Status *string
}

// InputOrderItem represents the request body for creating an order with a single item.
type InputOrderItem struct {
	Product_id int `json:"product_id"`
	Quantity   int `json:"quantity"`
}

// GetUserOrders retrieves the authenticated user's order and its items.
func (h *Handler) GetUserOrders(c *gin.Context) {
	var order models.Order
	var order_item []models.Order_Item

	owner_id := c.GetInt("user_id")

	if err := h.Db.Where("owner_id = ?", owner_id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := h.Db.Where("order_id = ?", order.Id).Find(&order_item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order, "order_items": order_item})
}

// CreateOrderFromCart converts the user's cart into an order within a transaction.
func (h *Handler) CreateOrderFromCart(c *gin.Context) {
	id := c.GetInt("user_id")
	var cart []models.Cart

	if err := h.Db.Where("user_id = ?", id).Find(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(cart) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Sepetiniz boş"})
		return
	}

	err := h.Db.Transaction(func(tx *gorm.DB) error {
		order := models.Order{
			Owner_id: id,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		var orderItems []models.Order_Item
		for _, item := range cart {
			orderItems = append(orderItems, models.Order_Item{
				Order_id:   order.Id,
				Product_id: item.Product_id,
				Quantity:   item.Quantity,
			})
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		if err := tx.Delete(&cart).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders cannot initalized: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Orders created successfully !"})
}

// CreateOrder creates a new order with a single order item.
func (h *Handler) CreateOrder(c *gin.Context) {
	id := c.GetInt("user_id")
	var input_order_item InputOrderItem
	var order models.Order
	var order_item models.Order_Item

	if err := c.ShouldBindJSON(&input_order_item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.Owner_id = id

	if err := h.Db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order_item.Order_id = order.Id
	order_item.Product_id = input_order_item.Product_id
	order_item.Quantity = input_order_item.Quantity

	if err := h.Db.Create(&order_item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})

}

// UpdateOrderStatus updates the status of an order. Admin only.
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	var input_order InputOrder
	id := c.Param("id")
	role := c.GetString("role")
	var order models.Order

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot update the status of order"})
		return
	}

	if err := c.ShouldBindJSON(&input_order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = *input_order.Status

	if err := h.Db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

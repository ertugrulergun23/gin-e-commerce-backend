package handlers

import (
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InputCart represents the request body for cart operations.
type InputCart struct {
	Product_id int
	Quantity   *int
}

// GetUserCart retrieves all cart items for the authenticated user.
func (h *Handler) GetUserCart(c *gin.Context) {
	var cart []models.Cart
	owner_id := c.GetInt("user_id")

	if err := h.Db.Where("user_id = ?", owner_id).Find(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddProductToCart adds a new product to the authenticated user's cart.
func (h *Handler) AddProductToCart(c *gin.Context) {
	var input InputCart
	user_id := c.GetInt("user_id")
	var cart models.Cart

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart.User_id = user_id
	cart.Product_id = input.Product_id
	cart.Quantity = *input.Quantity

	if err := h.Db.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product successfully added to cart", "cartId": cart.Id})
}

// UpdateCart updates the quantity of a cart item. Removes the item if quantity reaches zero.
func (h *Handler) UpdateCart(c *gin.Context) {
	var cart models.Cart
	id := c.Param("id")
	if err := h.Db.First(&cart, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	user_id := c.GetInt("user_id")

	if cart.User_id != user_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "This user cannot update this cart"})
		return
	}

	var input InputCart
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Quantity != nil {
		updates["quantity"] = cart.Quantity + *input.Quantity
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please insert a quantity value"})
		return
	}

	if updates["quantity"] == 0 {
		if err := h.Db.Delete(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := h.Db.Model(&cart).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart updated", "cartId": cart.Id})

}

// DeleteCart removes a cart item by ID. Only the cart owner can delete.
func (h *Handler) DeleteCart(c *gin.Context) {
	var cart models.Cart
	id := c.Param("id")
	if err := h.Db.First(&cart, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	user_id := c.GetInt("user_id")

	if user_id != cart.User_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "This user cannot delete this cart"})
		return
	}

	if err := h.Db.Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted successfully", "cart_id": cart.Id})
}

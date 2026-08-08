package handlers

import (
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InputCart struct {
	Quantity *int
}

// URL = /cart/:id
func (h *Handler) GetUserCart(c *gin.Context) {
	var carts []models.Cart
	id := c.Param("id")

	if err := h.Db.Where("user_id = ?", id).Find(&carts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, carts)
}

// URL = /cart/add
func (h *Handler) AddProductToCart(c *gin.Context) {
	var cart models.Cart

	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Db.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product successfully added to cart", "cartId": cart.Id})
}

// URL = /cart/:id/update
func (h *Handler) UpdateCart(c *gin.Context) {
	var cart models.Cart
	id := c.Param("id")
	if err := h.Db.First(&cart, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// URL = /cart/:id/delete
func (h *Handler) DeleteCart(c *gin.Context) {
	var cart models.Comment
	id := c.Param("id")
	if err := h.Db.First(&cart, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Db.Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted succesfully", "comment_id": cart.Id})
}

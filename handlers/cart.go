package handlers

import (
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InputCart struct {
	Product_id int
	Quantity   *int
}

// URL = /cart
func (h *Handler) GetUserCart(c *gin.Context) {
	var carts []models.Cart
	owner_id := c.GetInt("user_id")

	if err := h.Db.Where("user_id = ?", owner_id).Find(&carts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, carts)
}

// URL = /cart/add
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

	user_id := c.GetInt("user_id")

	if cart.User_id != user_id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This user cannot update this cart"})
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
	var cart models.Cart
	id := c.Param("id")
	if err := h.Db.First(&cart, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id := c.GetInt("user_id")

	if user_id != cart.User_id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This person cannot delete this comment"})
		return
	}

	if err := h.Db.Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted succesfully", "comment_id": cart.Id})
}

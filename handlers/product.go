package handlers

import (
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InputProduct struct {
	Name  *string  `json:"name"`
	Price *float64 `json:"price"`
	Stock *int     `json:"stock"`
	Point *float64 `json:"point"`
}

func (h *Handler) GetProducts(c *gin.Context) {
	name := "%" + c.Query("name") + "%"
	down_price := c.Query("down_price")
	up_price := c.Query("up_price")
	point := c.Query("point")

	query := h.Db.Model(&models.Product{})

	if name != "" {
		query = query.Where("name LIKE ? ", name)
	}
	if down_price != "" {
		query = query.Where("price >= ?", down_price)
	}
	if up_price != "" {
		query = query.Where("price <= ?", up_price)
	}
	if point != "" {
		query = query.Where("point = ?", point)
	}

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := h.Db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	role := c.GetString("role")

	if role != "seller" && role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can't create product"})
		return
	}
	var input InputProduct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller_id := c.GetInt("user_id")

	var product models.Product

	product.Seller_id = seller_id
	product.Name = *input.Name
	product.Stock = *input.Stock
	product.Price = *input.Price
	product.Point = *input.Point

	if err := h.Db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product_id": product.Id})
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	if err := h.Db.Find(&product, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller_id := c.GetInt("user_id")

	if product.Seller_id != seller_id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can not update this product"})
		return
	}

	var input InputProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}

	if input.Price != nil {
		updates["price"] = *input.Price
	}

	if input.Stock != nil {
		updates["stock"] = *input.Stock
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert at least one field to updated"})
		return
	}

	if err := h.Db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := h.Db.Find(&product, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	seller_id := c.GetInt("user_id")

	if product.Seller_id != seller_id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can not delete this product"})
		return
	}

	if err := h.Db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product successfully deleted", "product_id": product.Id})
}

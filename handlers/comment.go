package handlers

import (
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InputComment struct {
	Content *string `json:"content"`
	Point   *int    `json:"point"`
}

// URL = /comment/create
func (h *Handler) CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment-id": comment.Id})
}

// URL = /comment/:id
func (h *Handler) GetComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")

	if err := h.Db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// URL = /comments?owner=...&product=...
func (h *Handler) GetComments(c *gin.Context) {
	ownerId := c.Query("owner")
	productId := c.Query("product")

	var comments []models.Comment
	query := h.Db.Model(&models.Comment{})

	if ownerId != "" {
		query = query.Where("owner_id = ?", ownerId)
	}

	if productId != "" {
		query = query.Where("product_id = ?", productId)
	}

	if err := query.Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// URL = /comment/:id/update
func (h *Handler) UpdateComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")
	if err := h.Db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input InputComment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Content != nil {
		updates["content"] = *input.Content
	}

	if input.Point != nil {
		updates["point"] = *input.Point
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert at least one field to update"})
		return
	}

	if err := h.Db.Model(&comment).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// URL = /comment/:id/delete
func (h *Handler) DeleteComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")
	if err := h.Db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Db.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted succesfully", "comment_id": comment.Id})
}

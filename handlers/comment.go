package handlers

import (
	"ecommerce/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// InputComment represents the payload for creating or updating a comment.
type InputComment struct {
	Content *string `json:"content"`
	Point   *int    `json:"point"`
}

// CreateComment creates a new comment on a product for the authenticated user.
func (h *Handler) CreateComment(c *gin.Context) {
	owner_id := c.GetInt("user_id")
	product_id, err := strconv.Atoi(c.Param("product_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id parameter"})
		return
	}

	var input InputComment
	var comment models.Comment

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.Content = *input.Content
	comment.Point = *input.Point
	comment.Owner_id = owner_id
	comment.Product_id = product_id

	if err := h.Db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"comment-id": comment.Id})
}

// GetComment retrieves a single comment by ID.
func (h *Handler) GetComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")

	if err := h.Db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetComments retrieves comments with optional owner and product filters.
func (h *Handler) GetComments(c *gin.Context) {
	ownerId := c.Query("owner")
	productId := c.Query("product")

	var comments []models.Comment
	query := h.Db.Model(&models.Comment{})

	// Apply filter by owner ID if provided.
	if ownerId != "" {
		query = query.Where("owner_id = ?", ownerId)
	}

	// Apply filter by product ID if provided.
	if productId != "" {
		query = query.Where("product_id = ?", productId)
	}

	if err := query.Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// UpdateComment updates a comment's content or point. Only the owner can update.
func (h *Handler) UpdateComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")
	if err := h.Db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	owner_id := c.GetInt("user_id")

	// Ensure the current user is the owner of the comment.
	if owner_id != comment.Owner_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "This person cannot update this comment"})
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

// DeleteComment removes a comment by ID. Only the owner can delete.
func (h *Handler) DeleteComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")
	if err := h.Db.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	owner_id := c.GetInt("user_id")

	// Ensure the current user is the owner of the comment.
	if owner_id != comment.Owner_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "This person cannot delete this comment"})
		return
	}

	if err := h.Db.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully", "comment_id": comment.Id})
}

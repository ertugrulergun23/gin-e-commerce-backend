package handlers

import (
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// InputUser represents the payload for user profile updates.
type InputUser struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

// CreateUser registers a new user with a hashed password.
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password cannot hashed"})
		return
	}

	user.Password = string(hashedPassword)

	if err := h.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user-id": user.Id})
}

// GetUser retrieves the authenticated user's profile. Password is masked in the response.
func (h *Handler) GetUser(c *gin.Context) {
	var user models.User
	id := c.GetInt("user_id")

	if err := h.Db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Mask the password before returning the response.
	user.Password = " "

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates the authenticated user's profile fields (username, email, password).
func (h *Handler) UpdateUser(c *gin.Context) {
	var user models.User
	id := c.GetInt("user_id")
	if err := h.Db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input InputUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Username != nil {
		updates["name"] = *input.Username
	}

	if input.Password != nil {
		updates["password"] = *input.Password
	}

	if input.Email != nil {
		updates["email"] = *input.Email
	}

	// Ensure at least one field is provided for update.
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert at least one field to update"})
		return
	}

	if err := h.Db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser permanently deletes the authenticated user's account.
func (h *Handler) DeleteUser(c *gin.Context) {
	var user models.User
	id := c.GetInt("user_id")
	if err := h.Db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := h.Db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "user_id": user.Id})
}

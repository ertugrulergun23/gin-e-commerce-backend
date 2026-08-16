package handlers

import (
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type InputUser struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

// URL = /user/create
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

	c.JSON(http.StatusOK, gin.H{"user-id": user.Id})
}

// URL = /user
func (h *Handler) GetUser(c *gin.Context) {
	var user models.User
	id := c.GetInt("user_id")

	if err := h.Db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Kullanıcı bulunamadı !"})
		return
	}

	user.Password = " "

	c.JSON(http.StatusOK, user)
}

// URL = /user/update
func (h *Handler) UpdateUser(c *gin.Context) {
	var user models.User
	id := c.GetInt("user_id")
	if err := h.Db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert at least one field to update"})
		return
	}

	if err := h.Db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

}

// URL = /user/delete
func (h *Handler) DeleteUser(c *gin.Context) {
	var user models.User
	id := c.GetInt("user_id")
	if err := h.Db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted succesfully", "user_id": user.Id})
}

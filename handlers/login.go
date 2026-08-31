package handlers

import (
	"ecommerce/auth"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// InputUserLogin represents the payload required for user login.
type InputUserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login authenticates a user by validating input, finding the user, verifying password, and generating a JWT.
func (h *Handler) Login(c *gin.Context) {
	var input InputUserLogin
	var user models.User

	// Validate request payload
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing username or password"})
		return
	}

	// Find user by username
	if err := h.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare provided password with stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT authentication token
	token, err := auth.GenerateToken(user.Id, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Token": token})
}

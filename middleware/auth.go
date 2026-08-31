package middleware

import (
	"ecommerce/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired extracts and validates JWT from Authorization header, sets user_id and role in context.
func AuthRequired(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Token"})
		return
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")
	claims, err := auth.ParseToken(tokenStr)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	c.Set("user_id", claims.UserID)
	c.Set("role", claims.Role)
	c.Next()
}

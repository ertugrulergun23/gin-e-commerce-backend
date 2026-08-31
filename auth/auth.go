// Package auth provides authentication and JWT token management.
package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret is the secret key used for signing and parsing JWTs.
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// Claims represents custom JWT claims containing user ID and role.
type Claims struct {
	UserID int    `json:"id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token valid for 24 hours.
func GenerateToken(userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// ParseToken parses and validates a JWT token string, returning its claims.
func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

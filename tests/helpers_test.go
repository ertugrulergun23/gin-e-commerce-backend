package tests

import (
	"database/sql"
	"ecommerce/auth"
	"ecommerce/handlers"
	"ecommerce/routers"
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupMockDB creates a sqlmock database and wraps it with GORM for testing.
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return db, mock, gormDB
}

// setupRouter creates a Gin engine in test mode with all routes registered.
func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return routers.SetupRouters(db)
}

// setupHandler creates a raw Handler without full router (useful for targeted tests).
func setupHandler(db *gorm.DB) *handlers.Handler {
	return handlers.SetHandler(db)
}

// generateTestToken generates a valid JWT for a user in tests.
func generateTestToken(t *testing.T, userID int, role string) string {
	t.Helper()
	if os.Getenv("JWT_SECRET_KEY") == "" {
		os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	}
	token, err := auth.GenerateToken(userID, role)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return token
}

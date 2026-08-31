package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Ensure JWT secret is set for all tests in this package.
	if os.Getenv("JWT_SECRET_KEY") == "" {
		os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	}
}

// TestLogin_Success verifies that valid credentials return a token.
func TestLogin_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role"}).
		AddRow(1, "testuser", "test@example.com", string(hashed), "buyer")

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	body, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "password123",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["Token"] == "" {
		t.Error("expected Token in response but got empty")
	}
}

// TestLogin_InvalidBody verifies that missing body returns 400.
func TestLogin_InvalidBody(t *testing.T) {
	_, _, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// TestLogin_UserNotFound verifies that a non-existent user returns 400.
func TestLogin_UserNotFound(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnError(sqlmock.ErrCancelled)

	body, _ := json.Marshal(map[string]string{
		"username": "unknown",
		"password": "pass",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// TestLogin_WrongPassword verifies that a wrong password returns 401.
func TestLogin_WrongPassword(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role"}).
		AddRow(1, "testuser", "test@example.com", string(hashed), "buyer")

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	body, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// TestLogin_EmptyUsernameOrPassword verifies that empty fields return 401.
func TestLogin_EmptyUsernameOrPassword(t *testing.T) {
	_, _, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	// Only send one field — username missing → ShouldBindJSON won't error here
	// because both fields are plain strings, but the handler should fail on DB lookup.
	// We simulate the DB returning error.
	_, mock2, gormDB2 := setupMockDB(t)
	router2 := setupRouter(gormDB2)
	mock2.ExpectQuery(`SELECT \* FROM "users"`).WillReturnError(sqlmock.ErrCancelled)

	body, _ := json.Marshal(map[string]string{"username": "", "password": ""})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router2.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	_ = router
}

// Ensure time import is used to suppress unused import errors.
var _ = time.Now

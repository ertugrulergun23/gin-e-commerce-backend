package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserCart_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity"}).
		AddRow(1, 1, 10, 2).
		AddRow(2, 1, 11, 1)

	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs(1).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/cart", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetUserCart_DBError(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/cart", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestAddProductToCart_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "carts"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	qty := 3
	body, _ := json.Marshal(map[string]interface{}{
		"Product_id": 10,
		"Quantity":   qty,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/cart/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestAddProductToCart_InvalidBody(t *testing.T) {
	_, _, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/cart/add", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateCart_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity"}).
		AddRow(1, 1, 10, 2)
	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "carts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	qty := 2
	body, _ := json.Marshal(map[string]interface{}{
		"Quantity": qty,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/cart/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCart_QuantityBecomesZeroDeletes(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity"}).
		AddRow(1, 1, 10, 2)
	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "carts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	qty := -2
	body, _ := json.Marshal(map[string]interface{}{
		"Quantity": qty,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/cart/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateCart_NotOwner(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 2, "buyer")

	rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity"}).
		AddRow(1, 1, 10, 2)
	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	qty := 1
	body, _ := json.Marshal(map[string]interface{}{
		"Quantity": qty,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/cart/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDeleteCart_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity"}).
		AddRow(1, 1, 10, 2)
	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "carts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/cart/1/delete", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteCart_NotOwner(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 2, "buyer")

	rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity"}).
		AddRow(1, 1, 10, 2)
	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/cart/1/delete", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

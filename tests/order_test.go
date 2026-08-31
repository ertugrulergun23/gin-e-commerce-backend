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

func TestGetUserOrders_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	orderRows := sqlmock.NewRows([]string{"id", "owner_id", "status"}).
		AddRow(1, 1, "pending")
	mock.ExpectQuery(`SELECT \* FROM "orders"`).
		WithArgs(1, 1).
		WillReturnRows(orderRows)

	itemRows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity"}).
		AddRow(1, 1, 10, 2).
		AddRow(2, 1, 11, 1)
	mock.ExpectQuery(`SELECT \* FROM "order_items"`).
		WithArgs(1).
		WillReturnRows(itemRows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/order", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetUserOrders_NotFound(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	mock.ExpectQuery(`SELECT \* FROM "orders"`).
		WithArgs(1, 1).
		WillReturnError(fmt.Errorf("no orders found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/order", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCreateOrderFromCart_EmptyCart(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity"}))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/order/cart", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestCreateOrderFromCart_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	cartRows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity"}).
		AddRow(1, 1, 10, 2)
	mock.ExpectQuery(`SELECT \* FROM "carts"`).
		WithArgs(1).
		WillReturnRows(cartRows)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "orders"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`INSERT INTO "order_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`DELETE FROM "carts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/order/cart", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreateOrder_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "orders"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "order_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	body, _ := json.Marshal(map[string]interface{}{
		"product_id": 10,
		"quantity":   2,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/order/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateOrderStatus_AdminSuccess(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "admin")

	orderRows := sqlmock.NewRows([]string{"id", "owner_id", "status"}).
		AddRow(1, 2, "pending")
	mock.ExpectQuery(`SELECT \* FROM "orders"`).
		WithArgs("1", 1).
		WillReturnRows(orderRows)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "orders"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	status := "completed"
	body, _ := json.Marshal(map[string]interface{}{
		"Status": &status,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/order/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateOrderStatus_NonAdminUnauthorized(t *testing.T) {
	_, _, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	status := "completed"
	body, _ := json.Marshal(map[string]interface{}{
		"Status": &status,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/order/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

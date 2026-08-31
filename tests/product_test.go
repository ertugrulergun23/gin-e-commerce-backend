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

func TestGetProducts_NoFilter(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "point", "seller_id"}).
		AddRow(1, "ProductA", 19.99, 10, 4.5, 1).
		AddRow(2, "ProductB", 29.99, 5, 3.8, 2)

	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/products", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var products []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &products)
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestGetProducts_WithNameFilter(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "point", "seller_id"}).
		AddRow(1, "ProductA", 19.99, 10, 4.5, 1)

	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/products?name=ProductA&down_price=10&up_price=50&point=4.5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetProducts_DBError(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnError(fmt.Errorf("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/products", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetProduct_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "point", "seller_id"}).
		AddRow(1, "ProductA", 19.99, 10, 4.5, 1)

	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/product/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnError(fmt.Errorf("record not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/product/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCreateProduct_SellerSuccess(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "seller")

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "products"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
	mock.ExpectCommit()

	name := "NewProduct"
	price := 49.99
	stock := 100
	point := 0.0
	body, _ := json.Marshal(map[string]interface{}{
		"name":  name,
		"price": price,
		"stock": stock,
		"point": point,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/product/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreateProduct_BuyerUnauthorized(t *testing.T) {
	_, _, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 2, "buyer")

	body, _ := json.Marshal(map[string]interface{}{
		"name":  "Product",
		"price": 10.0,
		"stock": 5,
		"point": 0.0,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/product/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreateProduct_InvalidBody(t *testing.T) {
	_, _, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "seller")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/product/create", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateProduct_DBError(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "seller")

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "products"`).
		WillReturnError(fmt.Errorf("db insert error"))
	mock.ExpectRollback()

	name := "NewProduct"
	price := 49.99
	stock := 100
	point := 0.0
	body, _ := json.Marshal(map[string]interface{}{
		"name":  name,
		"price": price,
		"stock": stock,
		"point": point,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/product/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "seller")

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "point", "seller_id"}).
		AddRow(1, "ProductA", 19.99, 10, 4.5, 1)
	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	newName := "Updated Name"
	body, _ := json.Marshal(map[string]interface{}{
		"name": newName,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/product/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateProduct_NotOwner(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 2, "seller")

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "point", "seller_id"}).
		AddRow(1, "ProductA", 19.99, 10, 4.5, 1)
	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnRows(rows)

	newName := "Updated Name"
	body, _ := json.Marshal(map[string]interface{}{
		"name": newName,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/product/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestUpdateProduct_NoFields(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "seller")

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "point", "seller_id"}).
		AddRow(1, "ProductA", 19.99, 10, 4.5, 1)
	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnRows(rows)

	body, _ := json.Marshal(map[string]interface{}{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/product/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "seller")

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "point", "seller_id"}).
		AddRow(5, "ProductA", 19.99, 10, 4.5, 1)
	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "products"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/product/5/delete", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteProduct_NotOwner(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 2, "seller")

	rows := sqlmock.NewRows([]string{"id", "name", "price", "stock", "point", "seller_id"}).
		AddRow(1, "ProductA", 19.99, 10, 4.5, 1)
	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/product/1/delete", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

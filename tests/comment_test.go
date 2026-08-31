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

func TestGetComment_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	rows := sqlmock.NewRows([]string{"id", "owner_id", "product_id", "content", "point"}).
		AddRow(1, 1, 2, "Great product!", 5)

	mock.ExpectQuery(`SELECT \* FROM "comments"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/comment/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetComment_NotFound(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "comments"`).
		WillReturnError(fmt.Errorf("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/comment/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetComments_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	rows := sqlmock.NewRows([]string{"id", "owner_id", "product_id", "content", "point"}).
		AddRow(1, 1, 2, "Great!", 5).
		AddRow(2, 1, 3, "Okay", 3)

	mock.ExpectQuery(`SELECT \* FROM "comments"`).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/comments?owner=1&product=2", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetComments_DBError(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)

	mock.ExpectQuery(`SELECT \* FROM "comments"`).
		WillReturnError(fmt.Errorf("db query error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/comments", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestCreateComment_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "comments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
	mock.ExpectCommit()

	content := "Awesome!"
	point := 5
	body, _ := json.Marshal(map[string]interface{}{
		"content": content,
		"point":   point,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/comment/create/2", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreateComment_InvalidProductID(t *testing.T) {
	_, _, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/comment/create/abc", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateComment_InvalidBody(t *testing.T) {
	_, _, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/comment/create/1", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateComment_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	rows := sqlmock.NewRows([]string{"id", "owner_id", "product_id", "content", "point"}).
		AddRow(1, 1, 2, "Original comment", 4)
	mock.ExpectQuery(`SELECT \* FROM "comments"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "comments"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	content := "Updated comment"
	body, _ := json.Marshal(map[string]interface{}{
		"content": content,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/comment/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateComment_NotOwner(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 2, "buyer")

	rows := sqlmock.NewRows([]string{"id", "owner_id", "product_id", "content", "point"}).
		AddRow(1, 1, 2, "Original comment", 4)
	mock.ExpectQuery(`SELECT \* FROM "comments"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	content := "Updated comment"
	body, _ := json.Marshal(map[string]interface{}{
		"content": content,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/comment/1/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDeleteComment_Success(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 1, "buyer")

	rows := sqlmock.NewRows([]string{"id", "owner_id", "product_id", "content", "point"}).
		AddRow(1, 1, 2, "Original comment", 4)
	mock.ExpectQuery(`SELECT \* FROM "comments"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "comments"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/comment/1/delete", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteComment_NotOwner(t *testing.T) {
	_, mock, gormDB := setupMockDB(t)
	router := setupRouter(gormDB)
	token := generateTestToken(t, 2, "buyer")

	rows := sqlmock.NewRows([]string{"id", "owner_id", "product_id", "content", "point"}).
		AddRow(1, 1, 2, "Original comment", 4)
	mock.ExpectQuery(`SELECT \* FROM "comments"`).
		WithArgs("1", 1).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/comment/1/delete", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

package tests

import (
	"ecommerce/auth"
	"os"
	"testing"
)

func TestAuth_GenerateAndParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "my-test-secret")

	token, err := auth.GenerateToken(10, "admin")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := auth.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserID != 10 {
		t.Errorf("expected UserID 10, got %d", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("expected Role admin, got %s", claims.Role)
	}
}

func TestAuth_ParseInvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "my-test-secret")

	_, err := auth.ParseToken("not-a-valid-token")
	if err == nil {
		t.Error("expected error parsing invalid token, got nil")
	}
}

func TestAuth_ParseEmptyToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "my-test-secret")

	_, err := auth.ParseToken("")
	if err == nil {
		t.Error("expected error parsing empty token, got nil")
	}
}

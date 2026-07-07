package auth

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"

	token, err := GenerateToken(1, secret, time.Hour)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}

	if claims.UserID != 1 {
		t.Fatalf("expected user id 1, got %d", claims.UserID)
	}
}
func TestParseTokenWithWrongSecret(t *testing.T) {
	token, err := GenerateToken(1, "right-secret", time.Hour)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	_, err = ParseToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseExpiredToken(t *testing.T) {
	token, err := GenerateToken(1, "test-secret", -time.Hour)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	_, err = ParseToken(token, "test-secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

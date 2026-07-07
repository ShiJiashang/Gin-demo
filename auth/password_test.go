package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("123456")
	if err != nil {
		t.Fatalf("hash password failed: %v", err)
	}

	if hash == "123456" {
		t.Fatal("expected hashed password, got plain text")
	}

	if !CheckPassword("123456", hash) {
		t.Fatal("expected correct password to pass")
	}

	if CheckPassword("wrong", hash) {
		t.Fatal("expected wrong password to fail")
	}
}

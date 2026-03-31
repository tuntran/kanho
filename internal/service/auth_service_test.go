package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/auth"
)

func TestRegister_HashesPassword(t *testing.T) {
	password := "securePassword123!"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	if hash == password {
		t.Error("hash should not equal plaintext password")
	}
	if hash == "" {
		t.Error("hash should not be empty")
	}

	// Verify the hash is valid bcrypt
	if err := auth.VerifyPassword(hash, password); err != nil {
		t.Errorf("VerifyPassword should accept correct password: %v", err)
	}
}

func TestLogin_RejectsWrongPassword(t *testing.T) {
	password := "correctPassword"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	err = auth.VerifyPassword(hash, "wrongPassword")
	if err == nil {
		t.Error("VerifyPassword should reject wrong password")
	}
}

func TestRefreshToken_RejectsExpired(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	// Generate a token pair with very short TTL
	pair, err := auth.GenerateTokenPair(secret, userID, "test@test.com", time.Millisecond, time.Millisecond)
	if err != nil {
		t.Fatalf("GenerateTokenPair: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	_, err = auth.ValidateToken(secret, pair.RefreshToken)
	if err == nil {
		t.Error("ValidateToken should reject expired refresh token")
	}
}

package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateTokenPair_Success(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"

	pair, err := GenerateTokenPair(secret, userID, email, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pair.AccessToken == "" {
		t.Fatal("access token should not be empty")
	}
	if pair.RefreshToken == "" {
		t.Fatal("refresh token should not be empty")
	}
	if pair.RefreshJTI == "" {
		t.Fatal("refresh JTI should not be empty")
	}
}

func TestValidateToken_ValidAccess(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()
	email := "user@test.com"

	pair, err := GenerateTokenPair(secret, userID, email, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	claims, err := ValidateToken(secret, pair.AccessToken)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
}

func TestValidateToken_ValidRefresh(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()
	email := "user@test.com"

	pair, err := GenerateTokenPair(secret, userID, email, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	claims, err := ValidateToken(secret, pair.RefreshToken)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.ID != pair.RefreshJTI {
		t.Errorf("expected JTI %s, got %s", pair.RefreshJTI, claims.ID)
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New()

	// Generate a token that expires immediately
	pair, err := GenerateTokenPair(secret, userID, "user@test.com", -1*time.Second, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	_, err = ValidateToken(secret, pair.AccessToken)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	userID := uuid.New()

	pair, err := GenerateTokenPair("secret-a", userID, "user@test.com", 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	_, err = ValidateToken("secret-b", pair.AccessToken)
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

func TestValidateToken_MalformedToken(t *testing.T) {
	_, err := ValidateToken("secret", "not-a-jwt-token")
	if err == nil {
		t.Fatal("expected error for malformed token, got nil")
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	_, err := ValidateToken("secret", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestHashJTI(t *testing.T) {
	jti := "test-jti-value"
	hash1 := HashJTI(jti)
	hash2 := HashJTI(jti)
	if hash1 != hash2 {
		t.Fatalf("hash should be deterministic: %s != %s", hash1, hash2)
	}
	if len(hash1) != 64 {
		t.Fatalf("expected sha256 hex length 64, got %d", len(hash1))
	}

	// Different input should give different hash
	hash3 := HashJTI("different-jti")
	if hash1 == hash3 {
		t.Fatal("different inputs should produce different hashes")
	}
}

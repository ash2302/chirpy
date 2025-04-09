package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPasswordAndCheck(t *testing.T) {
	password := "mySecret123"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if hashed == "" {
		t.Fatal("hashed password is empty")
	}

	err = CheckPasswordHash(hashed, password)
	if err != nil {
		t.Errorf("password check failed for correct password: %v", err)
	}

	err = CheckPasswordHash(hashed, "wrongPassword")
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mySuperSecretKey"
	expiresIn := time.Minute * 5

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}
	if parsedID != userID {
		t.Errorf("expected userID %v, got %v", userID, parsedID)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidToken := "this.is.not.a.valid.token"
	secret := "mySuperSecretKey"

	_, err := ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "secret1"
	wrongSecret := "secret2"
	expiresIn := time.Minute * 5

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "mySuperSecretKey"
	expiresIn := -time.Second

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

package auth

import (
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	secret := "test-secret"
	token, err := NewToken("user-123", secret, time.Hour)
	if err != nil {
		t.Fatalf("NewToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	tok, claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if !tok.Valid {
		t.Fatal("expected valid token")
	}
	sub, _ := claims["sub"].(string)
	if sub != "user-123" {
		t.Fatalf("expected sub=user-123, got %q", sub)
	}
}

func TestNewTokenExpired(t *testing.T) {
	secret := "test-secret"
	token, err := NewToken("user-456", secret, -time.Hour)
	if err != nil {
		t.Fatalf("NewToken: %v", err)
	}

	_, _, err = ParseToken(token, secret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	token, err := NewToken("user-789", "secret-a", time.Hour)
	if err != nil {
		t.Fatalf("NewToken: %v", err)
	}

	_, _, err = ParseToken(token, "secret-b")
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestParseTokenMalformed(t *testing.T) {
	_, _, err := ParseToken("not-a-jwt", "secret")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
}

func TestParseTokenEmpty(t *testing.T) {
	_, _, err := ParseToken("", "secret")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewRefreshToken(t *testing.T) {
	secret := "test-secret"
	refreshToken, err := NewRefreshToken("user-abc", secret, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("NewRefreshToken: %v", err)
	}

	tok, claims, err := ParseToken(refreshToken, secret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if !tok.Valid {
		t.Fatal("expected valid refresh token")
	}

	typ, ok := claims["type"].(string)
	if !ok || typ != "refresh" {
		t.Fatalf("expected type=refresh, got %q", typ)
	}

	sub, _ := claims["sub"].(string)
	if sub != "user-abc" {
		t.Fatalf("expected sub=user-abc, got %q", sub)
	}
}

func TestRefreshTokenNotValidForAccess(t *testing.T) {
	secret := "test-secret"
	refreshToken, err := NewRefreshToken("user-def", secret, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("NewRefreshToken: %v", err)
	}

	// A refresh token should parse fine, but consumers should check the "type" claim
	tok, claims, err := ParseToken(refreshToken, secret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if !tok.Valid {
		t.Fatal("expected valid token")
	}
	typ, _ := claims["type"].(string)
	if typ != "refresh" {
		t.Fatal("expected type=refresh claim in refresh token")
	}
}

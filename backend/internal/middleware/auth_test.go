package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ailivili/internal/auth"
)

func TestParseBearer(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		wantOK   bool
		wantTok  string
	}{
		{"valid", "Bearer abc123", true, "abc123"},
		{"lowercase", "bearer xyz", true, "xyz"},
		{"mixed case", "BEARER tok", true, "tok"},
		{"extra spaces", "Bearer   tok  ", true, "tok"},
		{"empty", "", false, ""},
		{"no bearer", "abc123", false, ""},
		{"only bearer", "Bearer", false, ""},
		{"basic auth", "Basic YWxhZGRpbjpvcGVuIHNlc2FtZQ==", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok, ok := parseBearer(tt.header)
			if ok != tt.wantOK {
				t.Fatalf("ok: got %v, want %v", ok, tt.wantOK)
			}
			if ok && tok != tt.wantTok {
				t.Fatalf("token: got %q, want %q", tok, tt.wantTok)
			}
		})
	}
}

func TestRequireAuthNoHeader(t *testing.T) {
	mw := RequireAuth("secret")
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuthInvalidToken(t *testing.T) {
	mw := RequireAuth("secret")
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer invalid-token")
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuthValidToken(t *testing.T) {
	secret := "test-secret"
	token, err := auth.NewToken("user-123", secret, time.Hour)
	if err != nil {
		t.Fatalf("NewToken: %v", err)
	}

	mw := RequireAuth(secret)
	var gotUserID string
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, ok := UserIDFromContext(r.Context())
		if !ok {
			t.Fatal("expected user ID in context")
		}
		gotUserID = uid
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	handler.ServeHTTP(w, r)

	if gotUserID != "user-123" {
		t.Fatalf("expected user-123, got %q", gotUserID)
	}
}

func TestRequireAuthWrongSecret(t *testing.T) {
	token, _ := auth.NewToken("user-456", "secret-a", time.Hour)

	mw := RequireAuth("secret-b")
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestUserIDFromContextNotSet(t *testing.T) {
	_, ok := UserIDFromContext(context.Background())
	if ok {
		t.Fatal("expected false when user ID not in context")
	}
}

package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("mypassword123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "mypassword123" {
		t.Fatal("hash should not equal plaintext")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	if err := CheckPassword(hash, "correct-horse-battery-staple"); err != nil {
		t.Fatalf("CheckPassword: expected match, got %v", err)
	}
}

func TestCheckPasswordWrong(t *testing.T) {
	hash, err := HashPassword("password1")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	if err := CheckPassword(hash, "password2"); err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestHashPasswordEmpty(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword empty: %v", err)
	}
	if err := CheckPassword(hash, ""); err != nil {
		t.Fatalf("CheckPassword empty: %v", err)
	}
}

func TestHashPasswordLong(t *testing.T) {
	long := make([]byte, 72) // bcrypt max
	for i := range long {
		long[i] = 'x'
	}
	hash, err := HashPassword(string(long))
	if err != nil {
		t.Fatalf("HashPassword long: %v", err)
	}
	if hash == "" {
		t.Fatal("expected hash for long password")
	}
}

func TestCheckPasswordInvalidHash(t *testing.T) {
	if err := CheckPassword("not-a-valid-bcrypt-hash", "anything"); err == nil {
		t.Fatal("expected error for invalid hash")
	}
}

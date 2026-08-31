package auth

import "testing"

func TestHashAndComparePassword(t *testing.T) {
	hash, err := HashPassword("s3cret-pass")
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	if err := ComparePassword(hash, "s3cret-pass"); err != nil {
		t.Errorf("ComparePassword() dengan password benar error: %v", err)
	}

	if err := ComparePassword(hash, "password-salah"); err == nil {
		t.Error("ComparePassword() dengan password salah harusnya error")
	}
}

func TestHashPassword_Salted(t *testing.T) {
	h1, err := HashPassword("same")
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}
	h2, err := HashPassword("same")
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}
	if h1 == h2 {
		t.Error("dua hash untuk input sama harusnya berbeda (salt)")
	}
}

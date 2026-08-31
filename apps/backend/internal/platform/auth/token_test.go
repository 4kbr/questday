package auth

import (
	"testing"
	"time"
)

func TestJWT_IssueVerifyRoundTrip(t *testing.T) {
	j := NewJWT("secret", time.Hour)

	tok, err := j.Issue("u1", "Asia/Jakarta")
	if err != nil {
		t.Fatalf("Issue() error: %v", err)
	}

	claims, err := j.Verify(tok)
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}
	if claims.UserID != "u1" {
		t.Errorf("UserID = %q, mau %q", claims.UserID, "u1")
	}
	if claims.Timezone != "Asia/Jakarta" {
		t.Errorf("Timezone = %q, mau %q", claims.Timezone, "Asia/Jakarta")
	}
}

func TestJWT_Expired(t *testing.T) {
	j := NewJWT("secret", -time.Hour)

	tok, err := j.Issue("u1", "Asia/Jakarta")
	if err != nil {
		t.Fatalf("Issue() error: %v", err)
	}

	if _, err := j.Verify(tok); err == nil {
		t.Error("Verify() token kedaluwarsa harusnya error")
	}
}

func TestJWT_Tampered(t *testing.T) {
	j := NewJWT("secret", time.Hour)

	tok, err := j.Issue("u1", "Asia/Jakarta")
	if err != nil {
		t.Fatalf("Issue() error: %v", err)
	}

	b := []byte(tok)
	if b[len(b)-1] == 'a' {
		b[len(b)-1] = 'b'
	} else {
		b[len(b)-1] = 'a'
	}

	if _, err := j.Verify(string(b)); err == nil {
		t.Error("Verify() token yang diutak-atik harusnya error")
	}
}

func TestJWT_WrongSecret(t *testing.T) {
	issuer := NewJWT("secret", time.Hour)
	verifier := NewJWT("secret-lain", time.Hour)

	tok, err := issuer.Issue("u1", "Asia/Jakarta")
	if err != nil {
		t.Fatalf("Issue() error: %v", err)
	}

	if _, err := verifier.Verify(tok); err == nil {
		t.Error("Verify() dengan secret berbeda harusnya error")
	}
}

// Package auth menyediakan primitif keamanan lintas-module: token JWT dan
// hashing password. Module user memakainya lewat interface (port), bukan
// bergantung ke detail implementasi.
package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword menghasilkan hash bcrypt dari password plaintext.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("auth: hash password: %w", err)
	}
	return string(b), nil
}

// ComparePassword mencocokkan plaintext dengan hash; error kalau tidak cocok.
func ComparePassword(hash, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return fmt.Errorf("auth: password tidak cocok: %w", err)
	}
	return nil
}

// Package user menangani registrasi, login, dan profil user.
//
// Menyimpan Timezone user — penting karena "quest per hari" bergantung pada
// batas hari lokal user (ADR-006). Module lain membaca timezone ini lewat JWT
// claims (ADR-013), bukan dengan meng-query module user.
package user

import (
	"errors"
	"time"
)

// User adalah entitas inti module ini. Murni domain: tak kenal HTTP maupun SQL.
type User struct {
	ID           string
	Email        string
	PasswordHash string // JANGAN pernah masuk response
	DisplayName  string
	Timezone     string // IANA, mis. "Asia/Jakarta"
	CreatedAt    time.Time
}

// Error domain module user. Dipetakan ke status HTTP lewat
// httpx.RegisterErrorMapping di New() (ADR-023).
var (
	ErrEmailTaken        = errors.New("email sudah terpakai")
	ErrInvalidCredential = errors.New("email atau password salah")
	ErrUserNotFound      = errors.New("user tidak ditemukan")
)

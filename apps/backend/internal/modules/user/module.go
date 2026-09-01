// Package user menangani registrasi, login, dan profil user.
//
// Menyimpan Timezone user — penting karena "quest per hari" bergantung pada
// batas hari lokal user (lihat DECISIONS). Module lain butuh timezone ini untuk
// menghitung "hari ini" dengan benar.
package user

import (
	"database/sql"
	"net/http"

	"questday/internal/platform/auth"
	"questday/internal/platform/httpx"
	"questday/internal/platform/validator"
)

// Module adalah satu-satunya API publik module user. service, handler, dan
// postgresRepository tetap unexported.
type Module struct {
	handler *handler
	svc     *service // disimpan untuk AsUserDirectory() di Phase 3 (T3.4)
}

// New merakit repo -> service -> handler dan mendaftarkan pemetaan error domain
// -> HTTP (ADR-023).
func New(db *sql.DB, issuer auth.Issuer) *Module {
	repo := newPostgresRepository(db)
	svc := newService(repo, issuer)
	h := newHandler(svc, validator.New())

	httpx.RegisterErrorMapping(ErrEmailTaken, http.StatusConflict, "email_taken")
	httpx.RegisterErrorMapping(ErrInvalidCredential, http.StatusUnauthorized, "invalid_credential")
	httpx.RegisterErrorMapping(ErrUserNotFound, http.StatusNotFound, "user_not_found")

	return &Module{handler: h, svc: svc}
}

// AsUserDirectory mengekspos service sebagai penyedia nama user. Dikonsumsi
// module scoring lewat port scoring.UserDirectory (ADR-014). Tipe balik
// *service (unexported) memang disengaja.
func (m *Module) AsUserDirectory() *service {
	return m.svc
}

// Package middleware berisi middleware HTTP kustom milik QuestDay.
//
// Untuk yang generik (RequestID, RealIP, Logger, Recoverer, Timeout) manfaatkan
// bawaan chi (github.com/go-chi/chi/v5/middleware) — dipasang di server/router.go.
// Yang di sini hanya yang spesifik ke domain kita, mis. autentikasi.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"questday/internal/platform/auth"
	"questday/internal/platform/httpx"
)

// ctxKey adalah tipe privat untuk key context — supaya tak bentrok dengan key
// dari package lain.
type ctxKey int

const (
	keyUserID ctxKey = iota
	keyTimezone
)

// Authenticator memvalidasi header "Authorization: Bearer <token>", lalu menaruh
// userID DAN timezone user ke request context (ADR-013). Token hilang/invalid
// -> 401.
func Authenticator(verifier auth.Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := r.Header.Get("Authorization")
			token, ok := bearerToken(raw)
			if !ok {
				httpx.Unauthorized(w, "token tidak ada")
				return
			}

			claims, err := verifier.Verify(token)
			if err != nil {
				httpx.Unauthorized(w, "token tidak valid")
				return
			}

			ctx := WithUserID(r.Context(), claims.UserID)
			ctx = WithTimezone(ctx, claims.Timezone)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// bearerToken mengambil token dari nilai header Authorization.
func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}
	token := strings.TrimSpace(header[len(prefix):])
	return token, token != ""
}

// WithUserID menaruh userID ke context.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyUserID, id)
}

// UserIDFrom mengambil userID dari context.
func UserIDFrom(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(keyUserID).(string)
	return id, ok && id != ""
}

// WithTimezone menaruh timezone user ke context.
func WithTimezone(ctx context.Context, tz string) context.Context {
	return context.WithValue(ctx, keyTimezone, tz)
}

// TimezoneFrom mengambil timezone user dari context. Dipakai handler quest
// (Phase 2) untuk menghitung "hari ini".
func TimezoneFrom(ctx context.Context) (string, bool) {
	tz, ok := ctx.Value(keyTimezone).(string)
	return tz, ok && tz != ""
}

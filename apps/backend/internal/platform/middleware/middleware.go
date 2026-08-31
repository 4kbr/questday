// Package middleware berisi middleware HTTP kustom milik QuestDay.
//
// Untuk yang generik (RequestID, RealIP, Logger, Recoverer, Timeout) manfaatkan
// bawaan chi (github.com/go-chi/chi/v5/middleware) — dipasang di server/router.go.
// Yang di sini hanya yang spesifik ke domain kita, mis. autentikasi.
package middleware

// TODO:
//   // Authenticator memvalidasi Bearer JWT, mengambil userID, lalu menaruhnya
//   // ke request context. Handler membaca userID dari context (lihat helper
//   // di bawah). Tolak 401 kalau token hilang/invalid.
//   func Authenticator(verifier auth.TokenVerifier) func(http.Handler) http.Handler
//
//   // Helper simpan/ambil userID dari context (pakai tipe key privat, JANGAN
//   // string biasa, supaya tak bentrok):
//   //   type ctxKey int
//   //   func WithUserID(ctx, id) context.Context
//   //   func UserIDFrom(ctx) (id, ok)

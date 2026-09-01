package server

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"questday/internal/platform/httpx"
)

// healthHandler = liveness: proses hidup. Selalu 200.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// readyHandler = readiness: dependency penting siap (di sini: DB bisa di-ping).
// Gagal -> 503. Versi disempurnakan di Phase 4 (T4.1).
func readyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			httpx.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}

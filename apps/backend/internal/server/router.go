package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"questday/internal/modules/quest"
	"questday/internal/modules/scoring"
	"questday/internal/modules/user"
	"questday/internal/platform/auth"
	"questday/internal/platform/middleware"
)

// routerDeps adalah bahan yang dibutuhkan buildRouter. Dikumpulkan di server.New
// (composition root).
type routerDeps struct {
	db          *sql.DB
	verifier    auth.Verifier
	userMod     *user.Module
	questMod    *quest.Module
	scoringMod  *scoring.Module
	corsOrigins []string
}

// buildRouter memasang middleware global, CORS, health check, lalu mem-mount
// tiap module di bawah /api/v1 dengan pemisahan rute publik vs terproteksi.
func buildRouter(d routerDeps) http.Handler {
	r := chi.NewRouter()

	// Middleware generik: pakai bawaan chi. chimw.RealIP sengaja TIDAK dipakai —
	// rentan IP-spoofing tanpa trusted proxy (GHSA-3fxj-6jh8-hvhx).
	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))

	if len(d.corsOrigins) > 0 {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   d.corsOrigins,
			AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Authorization", "Content-Type"},
			AllowCredentials: false, // token lewat header, bukan cookie (ADR-020)
			MaxAge:           300,
		}))
	}

	// Health check di luar prefix /api/v1 dan di luar grup terproteksi.
	r.Get("/healthz", healthHandler)
	r.Get("/readyz", readyHandler(d.db))

	r.Route("/api/v1", func(r chi.Router) {
		d.userMod.RegisterPublicRoutes(r)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticator(d.verifier))
			d.userMod.RegisterProtectedRoutes(r)
			d.questMod.RegisterRoutes(r)
			d.scoringMod.RegisterRoutes(r)
		})
	})

	return r
}

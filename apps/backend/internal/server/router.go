package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"questday/internal/modules/quest"
	"questday/internal/modules/scoring"
	"questday/internal/modules/user"
	"questday/internal/platform/auth"
	"questday/internal/platform/middleware"
)

// routerDeps adalah bahan yang dibutuhkan buildRouter. Dikumpulkan di server.New
// (composition root).
type routerDeps struct {
	db         *sql.DB
	verifier   auth.Verifier
	userMod    *user.Module
	questMod   *quest.Module
	scoringMod *scoring.Module
}

// buildRouter memasang middleware global, health check, lalu mem-mount tiap
// module di bawah /api/v1 dengan pemisahan rute publik vs terproteksi.
func buildRouter(d routerDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Health check di luar prefix /api/v1. Versi lengkap di Phase 4 (T4.1).
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

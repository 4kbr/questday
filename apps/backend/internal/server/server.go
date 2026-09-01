// Package server merakit seluruh aplikasi: membuat dependency bersama
// (auth), meng-instansiasi tiap module, lalu menyusun router.
//
// Di sinilah "composition root" — satu-satunya tempat yang tahu semua module
// sekaligus dan menyambungkan mereka.
package server

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"questday/internal/config"
	"questday/internal/modules/quest"
	"questday/internal/modules/scoring"
	"questday/internal/modules/user"
	"questday/internal/platform/auth"
)

// Server membungkus *http.Server yang sudah dirakit.
type Server struct {
	httpServer *http.Server
}

// New merakit dependency, module, dan router jadi satu Server siap jalan.
func New(cfg config.Config, db *sql.DB) *Server {
	jwt := auth.NewJWT(cfg.JWTSecret, cfg.JWTTTL)

	// Urutan penting: user dulu (dibutuhkan scoring lewat UserDirectory),
	// lalu scoring, lalu quest (butuh scoring lewat ScoreAwarder).
	userMod := user.New(db, jwt)
	scoringMod := scoring.New(db, userMod.AsUserDirectory())
	questMod := quest.New(db, scoringMod.AsScoreAwarder())

	handler := buildRouter(routerDeps{
		db:         db,
		verifier:   jwt,
		userMod:    userMod,
		questMod:   questMod,
		scoringMod: scoringMod,
	})

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + cfg.HTTPPort,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

// ListenAndServe menjalankan HTTP server (blocking).
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown mematikan HTTP server dengan anggun. Graceful shutdown penuh
// (termasuk db.Close) disempurnakan di Phase 4 (T4.3/T4.4).
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

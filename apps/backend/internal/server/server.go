// Package server merakit seluruh aplikasi: membuat dependency bersama
// (auth), meng-instansiasi tiap module, lalu menyusun router.
//
// Di sinilah "composition root" — satu-satunya tempat yang tahu semua module
// sekaligus dan menyambungkan mereka.
package server

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"questday/internal/config"
	"questday/internal/modules/quest"
	"questday/internal/modules/scoring"
	"questday/internal/modules/user"
	"questday/internal/platform/auth"
)

// Server membungkus HTTP server yang sudah dirakit plus koneksi DB (ditutup
// saat Shutdown).
type Server struct {
	httpServer *http.Server
	db         *sql.DB
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
		db:          db,
		verifier:    jwt,
		userMod:     userMod,
		questMod:    questMod,
		scoringMod:  scoringMod,
		corsOrigins: cfg.CORSAllowedOrigins,
	})

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + cfg.HTTPPort,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		db: db,
	}
}

// ListenAndServe menjalankan HTTP server (blocking). Mengembalikan
// http.ErrServerClosed saat Shutdown dipanggil — caller memperlakukannya normal.
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown menutup HTTP server dengan anggun lalu menutup koneksi DB.
func (s *Server) Shutdown(ctx context.Context) error {
	httpErr := s.httpServer.Shutdown(ctx)
	dbErr := s.db.Close()
	return errors.Join(httpErr, dbErr)
}

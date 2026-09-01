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

	userMod := user.New(db, jwt)

	handler := buildRouter(routerDeps{
		db:       db,
		verifier: jwt,
		userMod:  userMod,
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

// Command api adalah entrypoint HTTP server QuestDay.
//
// Untuk Phase 0 file ini sengaja minimal: muat config lalu jalankan router kecil
// dengan health check. Perakitan module, middleware stack, dan graceful shutdown
// menyusul di Phase 1 & Phase 4.
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"questday/internal/config"
	"questday/internal/platform/httpx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	// TODO(phase1/phase4): ganti dengan server.New(cfg, db) — mount modul /api/v1,
	// middleware stack, graceful shutdown (T1.9, T4.1–T4.4).
	log.Printf("questday listening on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}

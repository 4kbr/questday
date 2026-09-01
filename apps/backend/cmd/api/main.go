// Command api adalah entrypoint HTTP server QuestDay.
//
// Merakit lewat server.New (composition root). Graceful shutdown penuh
// menyusul di Phase 4 (T4.4).
package main

import (
	"log"

	"questday/internal/config"
	"questday/internal/platform/database"
	"questday/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db := database.MustConnect(cfg.DatabaseURL)
	defer db.Close()

	srv := server.New(cfg, db)

	log.Printf("questday listening on :%s", cfg.HTTPPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server: %v", err)
	}
}

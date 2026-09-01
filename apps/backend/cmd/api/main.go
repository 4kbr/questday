// Command api adalah entrypoint HTTP server QuestDay.
//
// Merakit lewat server.New (composition root), jalan sampai menerima
// SIGINT/SIGTERM, lalu shutdown dengan anggun (HTTP drain + db.Close).
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	srv := server.New(cfg, db)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()
	log.Printf("questday listening on :%s", cfg.HTTPPort)

	<-ctx.Done()
	log.Println("shutdown: sinyal diterima, menutup server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Println("shutdown: selesai")
}

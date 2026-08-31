// Command api adalah entrypoint HTTP server QuestDay.
//
// Tanggung jawab file ini (tipis saja — hanya perakitan & lifecycle):
//   - Muat konfigurasi dari environment  (internal/config).
//   - Buka koneksi database               (internal/platform/database).
//   - Rakit server + semua module         (internal/server).
//   - Jalankan HTTP server.
//   - Graceful shutdown saat terima sinyal (SIGINT/SIGTERM):
//       stop terima request -> drain -> tutup DB.
//
// Logika bisnis TIDAK ada di sini. Kalau file ini mulai gemuk, pindahkan ke
// internal/server.
package main

// TODO:
//   func main() {
//       cfg  := config.Load()
//       db   := database.MustConnect(cfg.DatabaseURL)
//       srv  := server.New(cfg, db)      // merakit router + module
//       // jalankan srv.ListenAndServe() di goroutine
//       // tunggu sinyal, lalu srv.Shutdown(ctx)
//   }

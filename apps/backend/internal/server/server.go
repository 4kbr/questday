// Package server merakit seluruh aplikasi: membuat dependency bersama
// (validator, auth), meng-instansiasi tiap module, lalu menyusun router.
//
// Di sinilah "composition root" — satu-satunya tempat yang tahu semua module
// sekaligus dan menyambungkan mereka (mis. menyuntik scoring sebagai
// implementasi port quest.ScoreAwarder).
package server

// TODO:
//   type Server struct { httpServer *http.Server }
//
//   // New merakit segalanya:
//   //   1. Buat platform deps: auth.NewJWT(...), validator.New(), dll.
//   //   2. Instansiasi module:
//   //        userMod    := user.New(db, jwt, ...)
//   //        scoringMod := scoring.New(db)
//   //        questMod   := quest.New(db, scoringMod.AsScoreAwarder())  // port!
//   //        achMod     := achievement.New(db)   // v2, boleh diskip dulu
//   //   3. Bangun router (router.go) dan mount tiap module.
//   //   4. Bungkus ke *http.Server (addr dari cfg.HTTPPort).
//   func New(cfg config.Config, db *sql.DB) *Server
//
//   // ListenAndServe & Shutdown mem-forward ke http.Server.
//   func (s *Server) ListenAndServe() error
//   func (s *Server) Shutdown(ctx context.Context) error

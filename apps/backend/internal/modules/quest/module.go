// Package quest mengelola definisi quest (template task harian) dan log
// penyelesaiannya per hari.
//
// KEPUTUSAN KUNCI: Quest (definisi) dan QuestLog (instance harian) dipisah.
//   - Quest    : "Lari pagi", 20 poin, health, recurring harian. Dibuat/edit user.
//   - QuestLog : "2026-08-31, user X menyelesaikan Lari pagi". Bertambah tiap hari.
// Semua perhitungan (hari ini selesai apa, streak) berbasis QuestLog.
//
// file module.go = perakitan module: menyambungkan repository -> service ->
// handler, dan mengekspos RegisterRoutes. Inilah satu-satunya API publik module
// yang dipanggil server.
package quest

// TODO:
//   type Module struct { handler *handler }
//
//   // New merakit module. scorer disuntik dari luar (server) sebagai
//   // implementasi port ScoreAwarder — supaya quest TIDAK import package scoring.
//   func New(db *sql.DB, scorer ScoreAwarder) *Module {
//       repo := newPostgresRepository(db)
//       svc  := newService(repo, scorer)
//       h    := newHandler(svc)
//       return &Module{handler: h}
//   }
//
//   // RegisterRoutes memasang endpoint quest ke router yang diberikan.
//   func (m *Module) RegisterRoutes(r chi.Router)

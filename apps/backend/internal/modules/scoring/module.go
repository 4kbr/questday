// Package scoring memegang SEMUA logika gamifikasi: poin, XP, level, streak,
// dan leaderboard. Module ini bereaksi terhadap penyelesaian quest.
//
// Service-nya mengimplementasikan port quest.ScoreAwarder. Server menyuntik
// scoring ke quest saat perakitan, sehingga quest tak perlu import scoring.
package scoring

// TODO:
//   type Module struct {
//       service *service
//       handler *handler
//   }
//   func New(db *sql.DB) *Module {
//       repo := newPostgresRepository(db)
//       svc  := newService(repo)
//       return &Module{service: svc, handler: newHandler(svc)}
//   }
//
//   // AsScoreAwarder mengekspos service sebagai implementasi quest.ScoreAwarder.
//   // Dipakai server: quest.New(db, scoringMod.AsScoreAwarder()).
//   func (m *Module) AsScoreAwarder() *service   // service memenuhi interface itu
//
//   func (m *Module) RegisterRoutes(r chi.Router)

// Package scoring memegang SEMUA logika gamifikasi: poin, XP, level, streak,
// dan leaderboard. Bereaksi terhadap penyelesaian quest.
//
// service-nya mengimplementasikan port quest.ScoreAwarder; server menyuntik
// scoring ke quest saat perakitan sehingga quest tak perlu import scoring.
//
// module.go = perakitan: repository -> service -> handler, plus RegisterRoutes
// dan AsScoreAwarder. Satu-satunya API publik module.
package scoring

import "database/sql"

// Module = API publik module scoring. service/handler/postgresRepository unexported.
type Module struct {
	svc     *service
	handler *handler
}

// New merakit module. dir (UserDirectory) disuntik dari luar (server) —
// diimplementasi module user; scoring tak import user (ADR-014).
func New(db *sql.DB, dir UserDirectory) *Module {
	repo := newPostgresRepository(db)
	svc := newService(repo, dir)
	return &Module{svc: svc, handler: newHandler(svc)}
}

// AsScoreAwarder mengekspos service sebagai implementasi quest.ScoreAwarder.
// Tipe balik *service (unexported) memang disengaja — pemanggil hanya memakainya
// lewat interface quest.ScoreAwarder.
func (m *Module) AsScoreAwarder() *service {
	return m.svc
}

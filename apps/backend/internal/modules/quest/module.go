// Package quest mengelola definisi quest (template task harian) dan log
// penyelesaiannya per hari.
//
// KEPUTUSAN KUNCI: Quest (definisi) dan QuestLog (instance harian) dipisah
// (ADR-004). Semua perhitungan (hari ini selesai apa, streak) berbasis QuestLog.
//
// module.go = perakitan module: repository -> service -> handler, plus
// RegisterRoutes. Satu-satunya API publik module yang dipanggil server.
package quest

import (
	"database/sql"
	"net/http"

	"questday/internal/platform/httpx"
	"questday/internal/platform/validator"
)

// Module = API publik module quest. service/handler/postgresRepository unexported.
type Module struct {
	handler *handler
}

// New merakit module. scorer disuntik dari luar (server) sebagai implementasi
// port ScoreAwarder — supaya quest TIDAK import package scoring (ADR-005).
func New(db *sql.DB, scorer ScoreAwarder) *Module {
	repo := newPostgresRepository(db)
	svc := newService(repo, scorer)
	h := newHandler(svc, validator.New())

	httpx.RegisterErrorMapping(ErrQuestNotFound, http.StatusNotFound, "quest_not_found")
	// ErrNotOwner sengaja memakai kode & status yang sama dengan not-found —
	// jangan bocorkan bahwa quest itu ada dan milik orang lain.
	httpx.RegisterErrorMapping(ErrNotOwner, http.StatusNotFound, "quest_not_found")
	httpx.RegisterErrorMapping(ErrAlreadyCompleted, http.StatusConflict, "already_completed")
	httpx.RegisterErrorMapping(ErrNotCompleted, http.StatusConflict, "not_completed")

	return &Module{handler: h}
}

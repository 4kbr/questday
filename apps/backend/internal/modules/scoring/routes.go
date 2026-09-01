package scoring

import "github.com/go-chi/chi/v5"

// RegisterRoutes memasang endpoint scoring. Dipanggil server DI DALAM grup
// terproteksi (ber-Authenticator). /leaderboard tetap terproteksi untuk MVP
// (papan peringkat internal, bukan halaman publik).
func (m *Module) RegisterRoutes(r chi.Router) {
	r.Get("/me/score", m.handler.score)
	r.Get("/me/streak", m.handler.streak)
	r.Get("/leaderboard", m.handler.leaderboard)
}

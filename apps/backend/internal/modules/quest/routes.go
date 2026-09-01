package quest

import "github.com/go-chi/chi/v5"

// RegisterRoutes memasang endpoint quest. Dipanggil server DI DALAM grup
// terproteksi (ber-Authenticator) — semua rute quest butuh login.
//
// Perangkap chi: "/today" harus didaftarkan SEBELUM "/{questId}", kalau tidak
// "today" tertangkap sebagai questId.
func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/quests", func(r chi.Router) {
		r.Get("/", m.handler.list)
		r.Post("/", m.handler.create)
		r.Get("/today", m.handler.today)
		r.Patch("/{questId}", m.handler.update)
		r.Delete("/{questId}", m.handler.archive)
		r.Post("/{questId}/complete", m.handler.complete)
		r.Post("/{questId}/uncomplete", m.handler.uncomplete)
	})
}

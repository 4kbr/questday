// Package achievement = badge/lencana (mis. "7 hari beruntun", "100 quest selesai").
//
// STATUS: POST-MVP (v2). Scaffold ini sengaja ada supaya strukturnya lengkap,
// tapi JANGAN dikerjakan sampai MVP (user + quest + scoring) selesai. Boleh
// dibiarkan tidak di-mount di router untuk sekarang.
//
// Pola sama seperti module lain: definition (aturan badge) + unlock (badge yang
// sudah didapat user). Bereaksi ke event yang sama dengan scoring (quest
// selesai, streak naik) — pertimbangkan pindah ke event bus saat menggarap ini
// supaya scoring & achievement tidak saling numpuk di orkestrasi (lihat DECISIONS).
package achievement

// TODO (v2):
//   type Module struct { handler *handler }
//   func New(db *sql.DB) *Module
//   func (m *Module) RegisterRoutes(r chi.Router)

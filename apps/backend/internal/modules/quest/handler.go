package quest

// handler.go = lapisan HTTP: baca request -> panggil service -> tulis response.
// TIPIS. Tugasnya cuma: decode+validate DTO, ambil userID dari context,
// konversi DTO<->domain, dan map error domain -> status (via httpx).
//
// TODO:
//   type handler struct { svc *service }
//   func newHandler(svc *service) *handler
//
//   func (h *handler) create(w, r)      // POST   /quests
//   func (h *handler) list(w, r)        // GET    /quests
//   func (h *handler) update(w, r)      // PATCH  /quests/{questId}
//   func (h *handler) archive(w, r)     // DELETE /quests/{questId}
//   func (h *handler) today(w, r)       // GET    /quests/today
//   func (h *handler) complete(w, r)    // POST   /quests/{questId}/complete
//   func (h *handler) uncomplete(w, r)  // POST   /quests/{questId}/uncomplete
//
//   // Catatan: "hari ini" harus dihitung dari timezone user, bukan server.
//   // Ambil timezone dari profil user (atau dari context) sebelum panggil service.

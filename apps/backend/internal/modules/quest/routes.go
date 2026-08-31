package quest

// routes.go = pemetaan URL -> handler. Dipanggil dari module.go RegisterRoutes.
//
// TODO:
//   func (m *Module) RegisterRoutes(r chi.Router) {
//       r.Route("/quests", func(r chi.Router) {
//           r.Post("/", m.handler.create)
//           r.Get("/", m.handler.list)
//           r.Get("/today", m.handler.today)
//           r.Patch("/{questId}", m.handler.update)
//           r.Delete("/{questId}", m.handler.archive)
//           r.Post("/{questId}/complete", m.handler.complete)
//           r.Post("/{questId}/uncomplete", m.handler.uncomplete)
//       })
//   }

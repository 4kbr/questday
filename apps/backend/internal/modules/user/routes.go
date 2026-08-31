package user

// routes.go — perhatikan: register/login itu PUBLIK (tanpa auth), sedangkan
// /me butuh auth. Pemisahan group auth dilakukan di server/router.go; di sini
// cukup daftarkan path-nya.
//
// TODO:
//   func (m *Module) RegisterRoutes(r chi.Router) {
//       r.Post("/auth/register", m.handler.register)
//       r.Post("/auth/login", m.handler.login)
//       r.Get("/me", m.handler.me)   // pastikan berada di group ber-Authenticator
//   }

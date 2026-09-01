package user

import "github.com/go-chi/chi/v5"

// RegisterPublicRoutes mendaftarkan rute tanpa auth: register & login.
func (m *Module) RegisterPublicRoutes(r chi.Router) {
	r.Post("/auth/register", m.handler.register)
	r.Post("/auth/login", m.handler.login)
}

// RegisterProtectedRoutes mendaftarkan rute yang butuh Bearer token. Dipasang
// oleh server di dalam group ber-Authenticator — kalau tidak, handler tak akan
// menemukan userID di context.
func (m *Module) RegisterProtectedRoutes(r chi.Router) {
	r.Get("/me", m.handler.me)
	r.Patch("/me", m.handler.updateMe)
}

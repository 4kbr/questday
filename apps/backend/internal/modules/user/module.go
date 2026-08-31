// Package user menangani registrasi, login, dan profil user.
//
// Menyimpan Timezone user — penting karena "quest per hari" bergantung pada
// batas hari lokal user (lihat DECISIONS). Module lain butuh timezone ini untuk
// menghitung "hari ini" dengan benar.
package user

// TODO:
//   type Module struct { handler *handler }
//   func New(db *sql.DB, issuer auth.Issuer) *Module {
//       repo := newPostgresRepository(db)
//       svc  := newService(repo, issuer)
//       return &Module{handler: newHandler(svc)}
//   }
//   func (m *Module) RegisterRoutes(r chi.Router)

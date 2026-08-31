// GANTI module path ini dengan repo kamu, mis:
//   module github.com/username/questday
// Import internal nanti mengikuti path ini, mis:
//   github.com/username/questday/internal/modules/quest
module github.com/yourorg/questday

go 1.23

// TODO: tambah dependency saat mulai implement, mis:
//   require (
//       github.com/go-chi/chi/v5 v5.x.x
//       github.com/jackc/pgx/v5 v5.x.x           // driver postgres (atau lib/pq)
//       github.com/golang-jwt/jwt/v5 v5.x.x      // token auth
//       github.com/go-playground/validator/v10   // validasi DTO
//       golang.org/x/crypto                       // bcrypt password
//   )

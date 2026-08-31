package user

// repository_postgres.go = implementasi Postgres user. Satu-satunya tempat SQL user.
//
// TODO:
//   type postgresRepository struct { db *sql.DB }
//   func newPostgresRepository(db *sql.DB) *postgresRepository
//   // Implement Create/GetByEmail/GetByID. Manfaatkan unique index email untuk
//   // deteksi ErrEmailTaken.

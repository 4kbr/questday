package scoring

// repository_postgres.go = implementasi Postgres untuk scoring. Satu-satunya
// tempat SQL scoring. Pertimbangkan UPSERT (ON CONFLICT) untuk wallet & streak.
//
// TODO:
//   type postgresRepository struct { db *sql.DB }
//   func newPostgresRepository(db *sql.DB) *postgresRepository
//   // Implement semua method Repository.

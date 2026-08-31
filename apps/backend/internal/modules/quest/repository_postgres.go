package quest

// repository_postgres.go = implementasi Repository memakai Postgres.
// SATU-SATUNYA tempat SQL untuk module quest. Jangan tulis SQL di service/handler.
//
// TODO:
//   type postgresRepository struct { db *sql.DB }
//   func newPostgresRepository(db *sql.DB) *postgresRepository
//   // Implement semua method interface Repository dengan query SQL.
//   // Ingat: filter berdasarkan user_id di setiap query (jangan bocor antar user).
//   // Petakan sql.ErrNoRows -> ErrQuestNotFound.

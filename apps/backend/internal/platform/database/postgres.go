// Package database mengelola koneksi ke Postgres.
//
// Hanya soal koneksi & pool — TIDAK ada query domain di sini. Query milik
// masing-masing module (repository_postgres.go di tiap module).
package database

// TODO:
//   // Connect membuka pool koneksi, set batas (MaxOpenConns, MaxIdleConns,
//   // ConnMaxLifetime), lalu Ping untuk memastikan hidup.
//   func Connect(dsn string) (*sql.DB, error)   // atau *pgxpool.Pool kalau pakai pgx
//
//   // MustConnect: versi yang panic kalau gagal — dipakai di main saat startup.
//   func MustConnect(dsn string) *sql.DB

// Package database mengelola koneksi ke Postgres.
//
// Hanya soal koneksi & pool — TIDAK ada query domain di sini. Query milik
// masing-masing module (repository_postgres.go di tiap module).
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // driver "pgx" untuk database/sql (ADR-012)
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
	pingTimeout     = 5 * time.Second
)

// Connect membuka pool koneksi, set batasnya, lalu Ping untuk memastikan hidup.
// Kembalikan error (bukan panic) kalau DB tidak bisa dihubungi.
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("database: buka koneksi: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	return db, nil
}

// MustConnect memanggil Connect dan panic kalau gagal.
// Dipakai hanya di main saat startup — fail fast.
func MustConnect(dsn string) *sql.DB {
	db, err := Connect(dsn)
	if err != nil {
		panic(err)
	}
	return db
}

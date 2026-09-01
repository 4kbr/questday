// Package config memuat konfigurasi aplikasi dari environment variable.
//
// Satu-satunya tempat membaca os.Getenv. Bagian lain menerima Config lewat
// argumen, bukan membaca env langsung — supaya gampang di-test & jelas.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config adalah seluruh konfigurasi runtime aplikasi.
type Config struct {
	Env                string        // development | production
	HTTPPort           string        // mis. "8080"
	DatabaseURL        string        // DSN Postgres
	JWTSecret          string        // secret penandatangan token
	JWTTTL             time.Duration // masa berlaku token
	CORSAllowedOrigins []string      // origin yang diizinkan; kosong = CORS mati
}

// Load membaca env, memberi default yang wajar, lalu memvalidasi nilai wajib.
// Kembalikan zero Config + error kalau ada yang kurang / tidak valid.
func Load() (Config, error) {
	cfg := Config{
		Env:         getenv("APP_ENV", "development"),
		HTTPPort:    getenv("HTTP_PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", ""),
		JWTSecret:   getenv("JWT_SECRET", ""),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("config: DATABASE_URL wajib diisi")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("config: JWT_SECRET wajib diisi")
	}

	ttlRaw := getenv("JWT_TTL", "24h")
	ttl, err := time.ParseDuration(ttlRaw)
	if err != nil {
		return Config{}, fmt.Errorf("config: JWT_TTL tidak valid (%q): %w", ttlRaw, err)
	}
	cfg.JWTTTL = ttl

	cfg.CORSAllowedOrigins = splitList(getenv("CORS_ALLOWED_ORIGINS", ""))

	return cfg, nil
}

// splitList memecah string comma-separated jadi slice, membuang spasi & entri
// kosong. "" -> nil.
func splitList(raw string) []string {
	var out []string
	for part := range strings.SplitSeq(raw, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// getenv mengembalikan nilai env key, atau def kalau kosong / tidak ada.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

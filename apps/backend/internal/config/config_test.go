package config

import (
	"strings"
	"testing"
	"time"
)

func setRequired(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://localhost/questday")
	t.Setenv("JWT_SECRET", "secret")
	// Netralkan var opsional supaya test tak terpengaruh environment pemanggil
	// (Makefile `make test` meng-export isi .env). Test yang perlu menguji
	// nilai kustom meng-override lagi setelah memanggil helper ini.
	t.Setenv("APP_ENV", "")
	t.Setenv("HTTP_PORT", "")
	t.Setenv("JWT_TTL", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
}

func TestLoad_Defaults(t *testing.T) {
	setRequired(t)
	// APP_ENV / HTTP_PORT / JWT_TTL sengaja dikosongkan → pakai default.

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Env != "development" {
		t.Errorf("Env = %q, mau %q", cfg.Env, "development")
	}
	if cfg.HTTPPort != "8080" {
		t.Errorf("HTTPPort = %q, mau %q", cfg.HTTPPort, "8080")
	}
	if cfg.JWTTTL != 24*time.Hour {
		t.Errorf("JWTTTL = %v, mau %v", cfg.JWTTTL, 24*time.Hour)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	setRequired(t)
	t.Setenv("DATABASE_URL", "") // t.Setenv tidak bisa unset — pakai string kosong.

	_, err := Load()
	if err == nil {
		t.Fatal("Load() harusnya error saat DATABASE_URL kosong")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("error %q tidak menyebut DATABASE_URL", err.Error())
	}
}

func TestLoad_MissingJWTSecret(t *testing.T) {
	setRequired(t)
	t.Setenv("JWT_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() harusnya error saat JWT_SECRET kosong")
	}
	if !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Errorf("error %q tidak menyebut JWT_SECRET", err.Error())
	}
}

func TestLoad_CustomJWTTTL(t *testing.T) {
	setRequired(t)
	t.Setenv("JWT_TTL", "1h")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.JWTTTL != time.Hour {
		t.Errorf("JWTTTL = %v, mau %v", cfg.JWTTTL, time.Hour)
	}
}

func TestLoad_InvalidJWTTTL(t *testing.T) {
	setRequired(t)
	t.Setenv("JWT_TTL", "bukan-durasi")

	if _, err := Load(); err == nil {
		t.Fatal("Load() harusnya error saat JWT_TTL tidak valid")
	}
}

func TestLoad_CORSAllowedOrigins(t *testing.T) {
	setRequired(t)
	t.Setenv("CORS_ALLOWED_ORIGINS", " http://a.com ,http://b.com,, ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	want := []string{"http://a.com", "http://b.com"}
	if len(cfg.CORSAllowedOrigins) != len(want) {
		t.Fatalf("CORSAllowedOrigins = %v, mau %v", cfg.CORSAllowedOrigins, want)
	}
	for i := range want {
		if cfg.CORSAllowedOrigins[i] != want[i] {
			t.Errorf("CORSAllowedOrigins[%d] = %q, mau %q", i, cfg.CORSAllowedOrigins[i], want[i])
		}
	}
}

func TestLoad_CORSAllowedOrigins_Empty(t *testing.T) {
	setRequired(t) // CORS_ALLOWED_ORIGINS dinetralkan ke ""

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(cfg.CORSAllowedOrigins) != 0 {
		t.Errorf("CORSAllowedOrigins = %v, mau kosong", cfg.CORSAllowedOrigins)
	}
}

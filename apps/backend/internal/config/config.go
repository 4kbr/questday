// Package config memuat konfigurasi aplikasi dari environment variable.
//
// Satu-satunya tempat membaca os.Getenv. Bagian lain menerima Config lewat
// argumen, bukan membaca env langsung — supaya gampang di-test & jelas.
package config

// TODO:
//   type Config struct {
//       Env         string        // development | production
//       HTTPPort    string        // mis. "8080"
//       DatabaseURL string
//       JWTSecret   string
//       JWTTTL      time.Duration
//   }
//
//   // Load membaca env, memberi default yang wajar, dan memvalidasi nilai wajib
//   // (mis. DATABASE_URL, JWT_SECRET). Kembalikan error kalau ada yang kurang.
//   func Load() (Config, error)
//
//   // Helper kecil kalau perlu: getenv(key, default), mustEnv(key), dst.

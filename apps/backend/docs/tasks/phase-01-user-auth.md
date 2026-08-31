# Phase 1 — User & Auth

**Tujuan:** module `user` lengkap (register / login / profil) plus middleware
autentikasi. Setelah phase ini server sudah bisa dijalankan dan diuji dengan
curl — ini titik pertama aplikasi "hidup".

**Prasyarat:** Phase 0 selesai (butuh `auth`, `httpx`, `validator`, `database`,
tabel `users`).

**Rujukan pola:** GUIDES §2 (menambah endpoint), ARCHITECTURE "lapisan di dalam
satu module".

---

## T1.1 — `user/domain.go`

- **Sentuh:** `internal/modules/user/domain.go`
- **Isi:**
  ```go
  type User struct {
      ID           string
      Email        string
      PasswordHash string    // JANGAN pernah masuk response
      DisplayName  string
      Timezone     string    // IANA, mis. "Asia/Jakarta"
      CreatedAt    time.Time
  }
  var ErrEmailTaken, ErrInvalidCredential, ErrUserNotFound error
  ```
- **Aturan:** file ini **tidak mengimpor apa pun** selain stdlib (`time`,
  `errors`). Tanpa HTTP, tanpa SQL, tanpa package kita yang lain.
- **DoD:** 1 struct + 3 error.
- **Verifikasi:** `make vet`

## T1.2 — `user/dto.go`

- **Sentuh:** `internal/modules/user/dto.go`
- **Isi:**
  ```go
  type RegisterRequest struct {
      Email       string `json:"email"        validate:"required,email"`
      Password    string `json:"password"     validate:"required,min=8"`
      DisplayName string `json:"display_name" validate:"required"`
      Timezone    string `json:"timezone"     validate:"omitempty,timezone"`
  }
  type LoginRequest  struct{ Email, Password string }
  type AuthResponse  struct{ Token string; User UserResponse }
  type UserResponse  struct{ ID, Email, DisplayName, Timezone string }
  func toUserResponse(u User) UserResponse
  ```
- **Keputusan yang diselesaikan di sini:** komentar scaffold menandai `Timezone`
  sekaligus `required` **dan** "default Asia/Jakarta" — dua hal yang tak bisa
  benar bersamaan. **Pilih: opsional**, default `"Asia/Jakarta"` diisi di service
  (T1.4). Validasi `omitempty,timezone` supaya string ngawur tetap ditolak.
- **Aturan:** `toUserResponse` **tidak menyalin `PasswordHash`**. Ini satu-satunya
  jalan User → HTTP.
- **DoD:** tak ada field sensitif di DTO response.
- **Verifikasi:** `grep -n PasswordHash internal/modules/user/dto.go` → kosong.

## T1.3 — `user/repository.go` + `repository_postgres.go`

- **Sentuh:** `internal/modules/user/repository.go`,
  `internal/modules/user/repository_postgres.go`
- **Isi:**
  ```go
  type Repository interface {
      Create(ctx context.Context, u User) error
      GetByEmail(ctx context.Context, email string) (User, error)
      GetByID(ctx context.Context, id string) (User, error)
  }
  type postgresRepository struct{ db *sql.DB }
  func newPostgresRepository(db *sql.DB) *postgresRepository
  ```
- **Aturan:**
  - SQL **hanya** di `repository_postgres.go`.
  - `Create` mendeteksi pelanggaran `UNIQUE(email)` → kembalikan `ErrEmailTaken`
    (cek `pgconn.PgError.Code == "23505"` atau `strings.Contains` sebagai
    fallback — pilih yang eksplisit).
  - `sql.ErrNoRows` → `ErrUserNotFound`.
  - ID digenerate pemanggil (service) pakai UUIDv7, bukan di SQL (ADR-015).
- **DoD:** 3 method terimplementasi, error domain dipetakan.
- **Verifikasi:** `make vet && make build`

## T1.4 — `user/service.go`

- **Sentuh:** `internal/modules/user/service.go`
- **Isi:**
  ```go
  type service struct{ repo Repository; issuer auth.Issuer }
  func newService(repo Repository, issuer auth.Issuer) *service
  func (s *service) Register(ctx, req RegisterRequest) (AuthResponse, error)
  func (s *service) Login(ctx, req LoginRequest) (AuthResponse, error)
  func (s *service) Profile(ctx, userID string) (UserResponse, error)
  ```
  `Register`: default timezone `"Asia/Jakarta"` kalau kosong → `HashPassword` →
  generate UUIDv7 → `repo.Create` → `issuer.Issue(id, timezone)`.
  `Login`: `GetByEmail` → `ComparePassword` → issue token.
- **Aturan keras:** `Login` mengembalikan **`ErrInvalidCredential` yang sama**
  baik email tak ditemukan maupun password salah. Jangan bocorkan email mana yang
  terdaftar.
- **Aturan:** service tak tahu HTTP dan tak menyentuh SQL.
- **DoD:** 3 use case; token berisi userID + timezone.
- **Verifikasi:** lihat T1.10.

## T1.5 — `user/handler.go`

- **Sentuh:** `internal/modules/user/handler.go`
- **Isi:** `type handler struct{ svc *service }`, `newHandler`, dan 3 method:
  `register`, `login`, `me` (ambil userID dari context lewat
  `middleware.UserIDFrom`).
- **Aturan:** handler **tipis** — hanya decode/validate → panggil service → tulis
  lewat `httpx`. Nol logika bisnis. Petakan error domain ke status:
  `ErrEmailTaken` → 409, `ErrInvalidCredential` → 401, `ErrUserNotFound` → 404.
- **DoD:** 3 handler, tak ada `json.Marshal` langsung.
- **Verifikasi:** `make build`

## T1.6 — `user/routes.go` — pecah publik vs terproteksi

- **Sentuh:** `internal/modules/user/routes.go`, `internal/modules/user/module.go`
- **Isi:**
  ```go
  func (m *Module) RegisterPublicRoutes(r chi.Router)    // POST /auth/register, POST /auth/login
  func (m *Module) RegisterProtectedRoutes(r chi.Router) // GET /me
  ```
- **Kenapa dipecah:** scaffold mendaftarkan ketiga rute dalam satu
  `RegisterRoutes`, sementara `server/router.go` memanggilnya **hanya di grup
  publik** → `GET /me` akan mendarat tanpa `Authenticator` dan handler-nya tak
  akan pernah menemukan userID di context. Ini bug yang tertanam di scaffold.
- **DoD:** dua method terpisah; `router.go` (T1.9) memasang masing-masing di grup
  yang benar.
- **Verifikasi:** `curl` ke `/me` tanpa token → 401 (T1.10).

## T1.7 — `user/module.go`

- **Sentuh:** `internal/modules/user/module.go`
- **Isi:**
  ```go
  type Module struct{ handler *handler; svc *service }
  func New(db *sql.DB, issuer auth.Issuer) *Module   // repo -> service -> handler
  ```
- **Aturan:** `module.go` adalah **satu-satunya API publik** module ini. `service`,
  `handler`, `postgresRepository` tetap unexported.
- **Catatan:** field `svc` disimpan karena Phase 3 butuh `AsUserDirectory()`
  (T3.4). Boleh ditambahkan nanti.
- **DoD:** `user.New(db, jwt)` merakit seluruh rantai.
- **Verifikasi:** `make build`

## T1.8 — `platform/middleware`: Authenticator + context helper

- **Sentuh:** `internal/platform/middleware/middleware.go`
- **Isi:**
  ```go
  type ctxKey int
  const (keyUserID ctxKey = iota; keyTimezone)

  func Authenticator(verifier auth.Verifier) func(http.Handler) http.Handler
  func WithUserID(ctx, id string) context.Context
  func UserIDFrom(ctx) (string, bool)
  func WithTimezone(ctx, tz string) context.Context
  func TimezoneFrom(ctx) (string, bool)
  ```
  `Authenticator`: baca header `Authorization: Bearer <token>` → `Verify` →
  taruh `UserID` **dan** `Timezone` ke context → next. Gagal → 401 lewat
  `httpx.Unauthorized`.
- **Aturan:**
  - Key context wajib **tipe privat**, bukan `string` — supaya tak bentrok.
  - Middleware generik (RequestID, RealIP, Logger, Recoverer, Timeout)
    **jangan ditulis sendiri** — pakai bawaan chi, dipasang di `server/router.go`.
  - Pakai `auth.Verifier` (bukan `auth.TokenVerifier` seperti di komentar
    scaffold — lihat T0.12).
- **Kenapa Timezone di sini:** konsekuensi ADR-013; Phase 2 mengandalkan
  `TimezoneFrom` untuk menghitung "hari ini".
- **DoD:** helper untuk userID & timezone lengkap.
- **Verifikasi:** `make build`

## T1.9 — `server`: rakit minimal supaya bisa dijalankan

- **Sentuh:** `internal/server/server.go`, `internal/server/router.go`,
  `internal/server/health.go`, `cmd/api/main.go`
- **Isi (versi minimal, disempurnakan di Phase 4):**
  ```go
  func New(cfg config.Config, db *sql.DB) *Server
  // 1. jwt := auth.NewJWT(cfg.JWTSecret, cfg.JWTTTL); v := validator.New()
  // 2. userMod := user.New(db, jwt)
  // 3. buildRouter: r.Route("/api/v1", func(r chi.Router){
  //        userMod.RegisterPublicRoutes(r)
  //        r.Group(func(r chi.Router){
  //            r.Use(middleware.Authenticator(jwt))
  //            userMod.RegisterProtectedRoutes(r)
  //        })
  //    })
  // 4. http.Server{Addr: ":" + cfg.HTTPPort}
  func (s *Server) ListenAndServe() error
  func (s *Server) Shutdown(ctx context.Context) error
  ```
  `main.go`: `config.Load()` → `database.MustConnect` → `server.New` → jalan.
  Health check dan graceful shutdown penuh menyusul di T4.1/T4.4 — untuk sekarang
  cukup `GET /healthz` mengembalikan 200 supaya bisa dicek hidup.
- **Aturan:** `server` adalah composition root — satu-satunya tempat yang tahu
  semua module. Module tidak saling tahu.
- **Catatan:** komentar scaffold di `main.go` menulis `cfg := config.Load()`
  (tanpa error), padahal signature-nya `Load() (Config, error)` — tangani
  error-nya.
- **DoD:** `make run` menyalakan server, `GET /healthz` → 200.
- **Verifikasi:** `make run` lalu `curl -i localhost:8080/healthz`

## T1.10 — Test phase 1

- **Sentuh (baru):** `internal/modules/user/service_test.go`
- **Isi:** fake `Repository` (map in-memory) + fake `auth.Issuer`. Kasus:
  - Register happy path → token terbit, timezone default terisi saat request kosong.
  - Register email sudah ada → `ErrEmailTaken`.
  - Login password salah **dan** email tak ada → **error yang sama**
    (`ErrInvalidCredential`).
  - `Profile` mengembalikan `UserResponse` tanpa PasswordHash.
- **Smoke test manual:**
  ```bash
  curl -sX POST localhost:8080/api/v1/auth/register -d '{"email":"a@b.c","password":"rahasia123","display_name":"A"}'
  curl -sX POST localhost:8080/api/v1/auth/login    -d '{"email":"a@b.c","password":"rahasia123"}'
  curl -si  localhost:8080/api/v1/me                                  # -> 401
  curl -si  localhost:8080/api/v1/me -H "Authorization: Bearer $TOKEN" # -> 200
  ```
- **DoD:** `make test` hijau; smoke test sesuai.

---

## Exit criteria Phase 1

- [ ] `make fmt && make vet && make build && make test` bersih.
- [ ] register → login → `GET /me` dengan Bearer token jalan end-to-end.
- [ ] `GET /me` tanpa token → 401.
- [ ] `PasswordHash` tidak muncul di response mana pun
      (`curl ... | grep -i password` → kosong).
- [ ] Token berisi `UserID` **dan** `Timezone` (siap dipakai Phase 2).
- [ ] `internal/modules/user` tidak meng-import module lain.

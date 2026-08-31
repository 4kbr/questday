# Phase 0 — Foundation

**Tujuan:** menyiapkan pondasi yang dipakai semua module: module path & dependency,
skema database, konfigurasi, dan seluruh `internal/platform/*`. Setelah phase ini
belum ada endpoint yang jalan, tapi semua bahan untuk membangun module sudah ada.

**Prasyarat:** Go 1.23+, Docker, `air`, `golang-migrate` (lihat GUIDES §1).
Postgres hidup: `docker compose -f docker-compose.dev.yml up -d`.

**Kenapa duluan:** T0.1 memblokir semua import, T0.2 memblokir semua kompilasi,
T0.3 memblokir semua `repository_postgres.go`, dan **T0.15 memblokir seluruh
`apps/frontend`** (frontend menggenerate type-nya dari kontrak — ADR-019).
Sisanya dipakai lintas module.

---

## T0.1 — Ganti module path Go

- **Sentuh:** `go.mod`
- **Isi:** `module github.com/yourorg/questday` → `module questday` (ADR-011).
  Hapus komentar TODO baris 1-4 yang menyuruh mengganti ini.
- **Aturan:** wajib **sebelum** menulis satu pun import internal — semua import
  mengikuti path ini (`questday/internal/modules/quest`, dst.).
- **DoD:** `go.mod` baris pertama `module questday`.
- **Verifikasi:** `head -1 go.mod`

## T0.2 — Tambah dependency & lahirkan `go.sum`

- **Sentuh:** `go.mod` (blok `require`), `go.sum` (baru)
- **Isi:**
  ```
  github.com/go-chi/chi/v5              // router
  github.com/jackc/pgx/v5               // driver postgres (dipakai via stdlib)
  github.com/golang-jwt/jwt/v5          // token auth
  github.com/go-playground/validator/v10 // validasi DTO
  golang.org/x/crypto                   // bcrypt
  github.com/google/uuid                // UUIDv7 (ADR-015)
  ```
  Lalu `make tidy`.
- **Catatan:** driver dipakai sebagai `database/sql` driver
  (`_ "github.com/jackc/pgx/v5/stdlib"`), bukan `pgxpool` — ADR-012.
- **Kenapa penting:** `Dockerfile` melakukan `COPY go.mod go.sum ./` dan **gagal
  build** selama `go.sum` belum ada.
- **DoD:** `go.sum` ada, 6 dependency di `go.mod`.
- **Verifikasi:** `make tidy && make build`

## T0.3 — Migrasi: DDL 6 tabel (`.up.sql`)

- **Sentuh:** `migrations/000001_init.up.sql` (saat ini **hanya komentar, nol DDL**)
- **Isi:**

  | Tabel | Kolom | Constraint / index |
  |---|---|---|
  | `users` | `id uuid pk`, `email text not null`, `password_hash text not null`, `display_name text not null`, `timezone text not null default 'Asia/Jakarta'`, `created_at timestamptz not null default now()` | **`UNIQUE(email)`** |
  | `quests` | `id uuid pk`, `user_id uuid → users`, `title`, `note`, `category`, `difficulty`, `recurrence`, `active bool default true`, `created_at timestamptz` | `INDEX(user_id)` |
  | `quest_logs` | `id uuid pk`, `quest_id uuid → quests`, `user_id uuid → users`, `date DATE not null`, `status text`, `points_awarded int`, `completed_at timestamptz` | **`UNIQUE(quest_id, date)`**, `INDEX(user_id, date)` |
  | `wallets` | `user_id uuid pk → users`, `total_points int default 0`, `xp int default 0`, `level int default 1` | — |
  | `streaks` | `user_id uuid pk → users`, `current int default 0`, `longest int default 0`, `last_active DATE` | — |
  | `point_transactions` | `id uuid pk`, `user_id uuid → users`, `quest_id uuid null`, `points int` (boleh negatif), `date DATE`, `created_at timestamptz default now()` | `INDEX(user_id, created_at)` |

- **Aturan:**
  - `id` bertipe `uuid` tapi **tanpa** `default gen_random_uuid()` — nilainya
    digenerate di aplikasi (ADR-015), jadi tak perlu ekstensi `pgcrypto` dan tak
    perlu `RETURNING id`.
  - `quest_logs.date` bertipe **DATE** dan berisi tanggal **lokal user**, bukan
    UTC (ADR-006). Jangan pakai `timestamptz` di sini.
  - `UNIQUE(quest_id, date)` adalah yang mencegah dobel-selesai (ADR-004) —
    jangan diandalkan hanya lewat cek di service.
  - `UNIQUE(email)` adalah sumber deteksi `user.ErrEmailTaken`.
  - Tabel `badges`/`unlocks` **tidak** dibuat sekarang (v2).
- **DoD:** 6 tabel + semua index/constraint di atas.
- **Verifikasi:** `make migrate-up` lalu `\dt` di psql — 6 tabel muncul.

## T0.4 — Migrasi: `.down.sql`

- **Sentuh:** `migrations/000001_init.down.sql`
- **Isi:** `DROP TABLE IF EXISTS` urutan terbalik agar FK aman:
  `point_transactions → streaks → wallets → quest_logs → quests → users`.
- **DoD:** `make migrate-down` lalu `make migrate-up` bisa berulang tanpa error.
- **Verifikasi:** `make migrate-down && make migrate-up`
- **Catatan:** kalau state jadi `dirty`, perbaiki SQL lalu
  `make migrate-force version=N` (GUIDES §6).

## T0.5 — `config`: struct + `Load()`

- **Sentuh:** `internal/config/config.go`
- **Isi:**
  ```go
  type Config struct {
      Env         string        // development | production
      HTTPPort    string        // "8080"
      DatabaseURL string
      JWTSecret   string
      JWTTTL      time.Duration
  }
  func Load() (Config, error)
  ```
  Plus helper privat `getenv(key, def string) string`. Beri default wajar
  (`APP_ENV=development`, `HTTP_PORT=8080`, `JWT_TTL=24h`), **validasi wajib**
  untuk `DATABASE_URL` dan `JWT_SECRET` → error yang menyebut nama variabelnya.
- **Aturan:** ini **satu-satunya tempat** yang boleh memanggil `os.Getenv`.
  Bagian lain menerima `Config` sebagai argumen.
- **DoD:** `Load()` mengembalikan error deskriptif kalau env wajib kosong.
- **Verifikasi:** `make vet`; lihat T0.13.
- **Rujukan:** `.env.example` (5 variabel, cocok 1:1 dengan struct ini).

## T0.6 — `platform/database`: koneksi Postgres

- **Sentuh:** `internal/platform/database/postgres.go`
- **Isi:**
  ```go
  func Connect(dsn string) (*sql.DB, error)  // pool limit + Ping
  func MustConnect(dsn string) *sql.DB       // panic; dipakai di main
  ```
  Import blank driver: `_ "github.com/jackc/pgx/v5/stdlib"`. Set
  `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`, lalu `PingContext`
  dengan timeout supaya startup tidak menggantung.
- **Aturan:** `platform/*` netral domain — tak boleh import `modules/*`.
- **DoD:** `Connect` mengembalikan error (bukan panic) kalau DB mati.
- **Verifikasi:** `make build`

## T0.7 — `platform/httpx`: response

- **Sentuh:** `internal/platform/httpx/response.go`
- **Isi:**
  ```go
  func JSON(w http.ResponseWriter, status int, v any)
  func Data(w http.ResponseWriter, status int, data any)  // envelope {"data": ...}
  func NoContent(w http.ResponseWriter)                   // 204
  ```
  Set `Content-Type: application/json` sebelum `WriteHeader`.
- **DoD:** semua response sukses di seluruh app lewat sini — handler tak pernah
  memanggil `json.NewEncoder` langsung.
- **Verifikasi:** `make build`

## T0.8 — `platform/httpx`: error + mapper domain → HTTP

- **Sentuh:** `internal/platform/httpx/errors.go`
- **Isi:**
  ```go
  type ErrorResponse struct{ Error struct{ Code, Message string } }
  // bentuk: {"error": {"code": "not_found", "message": "..."}}
  func Error(w, status int, code, message string)
  func BadRequest(w, msg); Unauthorized(w, msg); Forbidden(w, msg)
  func NotFound(w, msg); Conflict(w, msg); Internal(w)
  ```
  Plus **mapper** error domain → HTTP berbasis `errors.Is`, dipakai semua handler.
- **Aturan:**
  - `Internal(w)` **tidak boleh** membocorkan detail error ke client (log saja).
  - Mapper tidak boleh import `modules/*` (ADR: platform netral). Jadi bentuknya
    tabel yang **didaftarkan dari luar**, mis.
    `httpx.RegisterErrorMapping(err error, status int, code string)` yang
    dipanggil tiap module saat `New()`, **atau** tiap handler memetakan sendiri
    error modulenya. Pilih satu dan konsisten — kalau memilih registry, catat
    sebagai ADR.
- **DoD:** envelope error seragam di seluruh app; cocok dengan yang nanti ditulis
  di `contracts/openapi.yaml` sebagai `ErrorResponse` (T4.6).
- **Verifikasi:** `make build`

## T0.9 — `platform/httpx`: decode request

- **Sentuh:** `internal/platform/httpx/decode.go`
- **Isi:**
  ```go
  func Decode(w http.ResponseWriter, r *http.Request, dst any) error
  func DecodeAndValidate(w http.ResponseWriter, r *http.Request, dst any) error
  ```
  `Decode`: `http.MaxBytesReader` (mis. 1MB), `DisallowUnknownFields()`, pesan
  ramah untuk JSON rusak / field tak dikenal / body kosong.
- **DoD:** handler cukup memanggil `DecodeAndValidate` lalu langsung pakai DTO.
- **Verifikasi:** `make build`
- **Catatan:** `DecodeAndValidate` butuh `*validator.Validator` (T0.10) — putuskan
  apakah disuntik lewat argumen atau package-level var. Hindari global mutable.

## T0.10 — `platform/validator`

- **Sentuh:** `internal/platform/validator/validator.go`
- **Isi:**
  ```go
  type Validator struct{ v *validator.Validate }
  func New() *Validator
  func (v *Validator) Struct(s any) error  // error rapi: field -> pesan
  ```
  Register custom rule bila perlu (mis. `timezone` IANA valid via
  `time.LoadLocation`).
- **DoD:** error validasi bisa dibaca manusia, bukan dump `validator.FieldError`.
- **Verifikasi:** `make build`

## T0.11 — `platform/auth`: password

- **Sentuh:** `internal/platform/auth/password.go`
- **Isi:**
  ```go
  func HashPassword(plain string) (string, error)
  func ComparePassword(hash, plain string) error
  ```
  bcrypt (`golang.org/x/crypto/bcrypt`), cost default.
- **Catatan:** komentar scaffold menyebut nama "Hash"/"Compare" tapi signature-nya
  `HashPassword`/`ComparePassword` — pakai yang panjang, lebih jelas dari luar package.
- **DoD:** hash tak pernah sama untuk input sama (salt), Compare cocok.
- **Verifikasi:** lihat T0.13.

## T0.12 — `platform/auth`: JWT

- **Sentuh:** `internal/platform/auth/token.go`
- **Isi:**
  ```go
  type Claims struct {
      UserID   string
      Timezone string   // ADR-013
      jwt.RegisteredClaims
  }
  type Issuer   interface{ Issue(userID, timezone string) (string, error) }
  type Verifier interface{ Verify(token string) (Claims, error) }
  type JWT struct{ secret []byte; ttl time.Duration }
  func NewJWT(secret string, ttl time.Duration) *JWT
  ```
  `*JWT` mengimplementasi `Issuer` & `Verifier`.
- **Keputusan yang harus diselesaikan di sini:** scaffold menyebut dua nama untuk
  hal yang sama — `auth.Verifier` (token.go) vs `auth.TokenVerifier`
  (middleware.go). **Pakai `auth.Verifier`**, dan sesuaikan signature
  `middleware.Authenticator` di T1.8.
- **Aturan:** `Timezone` masuk claims (ADR-013) supaya handler quest tidak perlu
  query user tiap request. Konsekuensinya: user yang mengubah timezone baru
  merasakan efeknya setelah token baru — dicatat di ADR-013.
- **DoD:** `Verify` menolak token expired, signature salah, dan algoritma tak
  terduga (`alg: none`).
- **Verifikasi:** lihat T0.13.

## T0.13 — Test phase 0

- **Sentuh (baru):** `internal/config/config_test.go`,
  `internal/platform/auth/password_test.go`,
  `internal/platform/auth/token_test.go`
- **Isi:**
  - `config.Load`: default terisi saat env kosong; error saat `DATABASE_URL` /
    `JWT_SECRET` hilang (pakai `t.Setenv`).
  - `HashPassword`/`ComparePassword`: happy path, password salah → error.
  - `JWT`: issue→verify round-trip (UserID & Timezone kembali utuh); token
    expired ditolak; token yang diubah 1 karakter ditolak.
- **DoD:** `make test` hijau. Ini test pertama di repo — belum ada satu pun
  `*_test.go` sebelumnya.
- **Verifikasi:** `make test`

## T0.14 — ADR-011 s/d ADR-015 (SUDAH DITULIS — tinggal ditinjau)

- **Sentuh:** `docs/DECISIONS.md`
- **Status:** kelima ADR **sudah ada** di `docs/DECISIONS.md` (ditulis bersamaan
  dengan task list ini, karena keputusannya diambil sebelum implementasi dimulai):
  - **ADR-011** Module path `questday`.
  - **ADR-012** Akses DB `database/sql` + `pgx/v5/stdlib`.
  - **ADR-013** Timezone user ikut di JWT claims.
  - **ADR-014** Nama di leaderboard lewat port `UserDirectory`, bukan JOIN.
  - **ADR-015** PK `uuid` dengan UUID v7 digenerate di aplikasi.
- **Yang dikerjakan di sini:** baca ulang kelimanya sebelum mulai koding, dan
  kalau implementasi ternyata memaksa menyimpang — **tulis ADR baru yang
  men-supersede**, jangan diam-diam ditabrak (AGENTS.md).
- **DoD:** kelima ADR terbaca dan dipahami; ADR-001..010 tetap utuh.
- **Verifikasi:** `grep -c "^## ADR-" docs/DECISIONS.md` → 16 (15 ADR + template).

## T0.15 — Isi `contracts/openapi.yaml` (kontrak lebih dulu)

- **Sentuh:** `contracts/openapi.yaml`
- **Isi:**
  1. **15 schema** di `components.schemas` (sekarang `{}`): `RegisterRequest`,
     `LoginRequest`, `AuthResponse`, `UserResponse`, **`UpdateProfileRequest`**,
     `CreateQuestRequest`, `UpdateQuestRequest`, `QuestResponse`,
     `CompleteQuestRequest`, `QuestLogResponse`, `TodayQuestsResponse`,
     `ScoreResponse`, `StreakResponse`, `LeaderboardEntry`, `ErrorResponse`.
  2. **`requestBody` + `responses`** untuk 11 operasi yang sudah terdaftar
     (termasuk 400/401/404/409 yang merujuk `ErrorResponse`).
  3. **`parameters`**: path `{questId}` (uuid) di 3 operasi; query `limit` di
     `/leaderboard`.
  4. **Pasang `security: [bearerAuth: []]`** di operasi terproteksi.
     `securitySchemes.bearerAuth` sudah didefinisikan tapi belum dipakai di satu
     operasi pun.
  5. **Tambah endpoint yang belum ada di kontrak:** `GET /me`, **`PATCH /me`**
     (lihat T1.11), serta `GET /healthz` & `GET /readyz` (catat bahwa keduanya
     **di luar** prefix `/api/v1`).
- **Aturan:** `ErrorResponse` **wajib** persis sama bentuknya dengan envelope
  `httpx` (T0.8). Ini yang diminta TODO baris 9 di file kontrak.
- **Kenapa di Phase 0, bukan di akhir:** ini memang aturan repo — "kontrak lebih
  dulu" (AGENTS.md #8, ADR-002). Dan `apps/frontend` **memblokir di sini**:
  task F0.4 frontend menggenerate seluruh type TypeScript-nya dari file ini
  (ADR-019). Selama `components.schemas` masih `{}`, frontend tak bisa mulai.
  Mengisinya sekarang membuat dua app bisa jalan paralel.
- **Catatan:** bentuk DTO di sini menjadi acuan saat menulis `dto.go` di Phase 1-3.
  Kalau saat implementasi ternyata ada yang perlu berubah, **ubah kontraknya
  juga** di saat yang sama — jangan biarkan melenceng sampai Phase 4.
- **DoD:** nol `# TODO` tersisa di `contracts/openapi.yaml`; spec valid.
- **Verifikasi:** `grep -c TODO contracts/openapi.yaml` → 0, dan
  `npx @redocly/cli lint contracts/openapi.yaml` bersih.

---

## Exit criteria Phase 0

- [ ] `make fmt && make vet && make build && make test` bersih.
- [ ] `make migrate-up` sukses dari DB kosong; `make migrate-down` bersih; 6 tabel ada.
- [ ] `go.sum` ada (Dockerfile tidak lagi gagal di tahap `COPY`).
- [ ] Tidak ada file di `internal/platform/*` yang meng-import `internal/modules/*`.
- [ ] ADR-011..015 sudah dibaca; tak ada yang ditabrak diam-diam.
- [ ] `contracts/openapi.yaml` terisi lengkap (nol TODO) — **frontend memblokir di sini**.

# Hands-off — State Backend QuestDay

Dokumen orientasi untuk agent/kolaborator berikutnya yang mengerjakan
`apps/backend`. **Baca sesudah `AGENTS.md` (root), sebelum menyentuh kode.**

**Terakhir diperbarui:** 2026-09-01
**Status:** Phase 0 (Foundation) **selesai & ter-commit**. Phase 1 (User & Auth)
**selesai** (kode + unit test hijau; smoke test end-to-end menunggu Postgres).
Phase 2 (Quest) belum dimulai.

> Root `docs/HANDS-OFF.md` tetap jadi orientasi seluruh proyek (monorepo, titik
> sambung frontend). File ini fokus ke state backend saja.

---

## TL;DR

- Phase 0 tuntas di commit `1e3f737 backend: phase 0 foundation`
  (branch `feature/backend/m`).
- **Phase 1 (User & Auth) selesai (belum di-commit).** Endpoint hidup:
  `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /api/v1/me`,
  `PATCH /api/v1/me`. `make fmt && vet && build && test` bersih; unit test
  service user hijau. **Smoke test curl end-to-end belum dijalankan** —
  butuh Postgres (Docker daemon mati saat implementasi).
- `contracts/openapi.yaml`: semua response 2xx berbadan kini dibungkus
  `{"data": ...}` (ADR-024 baru) — cocok dengan `httpx.Data`. Masih valid
  (redocly: 0 error, 3 warning benign).
- Mulai berikutnya: `apps/backend/docs/tasks/phase-02-quest.md`, kerjakan
  **T2.1 → T2.10 berurutan**. Sebelum itu: jalankan Postgres + smoke test
  Phase 1 (lihat bawah), lalu commit Phase 1.

---

## Peta baca (urutan)

1. `AGENTS.md` (root) — aturan keras & batasan.
2. **File ini** — apa yang sudah/belum ada.
3. `docs/ARCHITECTURE.md` — bentuk sistem, aturan dependensi.
4. `docs/DECISIONS.md` — 23 ADR. **ADR-023 baru** (registry error mapper).
5. `apps/backend/docs/tasks/phase-01-user-auth.md` — task Phase 1.
6. `docs/GUIDES.md` §2 — resep menambah endpoint.
7. `contracts/openapi.yaml` — kontrak; acuan bentuk `dto.go`.

---

## Yang SUDAH jadi (Phase 0)

| Area | Kondisi |
|---|---|
| `go.mod` / `go.sum` | `module questday`, directive **`go 1.25.0`**. Dep langsung: `go-chi/chi/v5`, `jackc/pgx/v5`, `golang-jwt/jwt/v5`, `go-playground/validator/v10`, `golang.org/x/crypto`. `google/uuid` ada di module graph sebagai indirect — jadi direct otomatis begitu kode Phase 1 meng-import-nya (UUIDv7, ADR-015). `go.sum` lengkap. |
| `migrations/000001_init.{up,down}.sql` | 6 tabel: `users`, `quests`, `quest_logs`, `wallets`, `streaks`, `point_transactions`. `up`/`down`/`up` repeatable. UUID PK **tanpa** `DEFAULT` (app generate). `quest_logs.date` & `point_transactions.date` = `date` (tanggal lokal user, ADR-006 — **bukan** `timestamptz`). `UNIQUE(quest_id, date)` di `quest_logs`. Index `quests(user_id)`, `quest_logs(user_id,date)`, `point_transactions(user_id,created_at)`. FK `ON DELETE CASCADE` (kecuali `point_transactions.quest_id` → `SET NULL`). Tanpa `badges`/`unlocks` (v2). |
| `internal/config/config.go` | `Config{Env, HTTPPort, DatabaseURL, JWTSecret string; JWTTTL time.Duration}`, `Load() (Config, error)`. Default `APP_ENV=development`, `HTTP_PORT=8080`, `JWT_TTL=24h`. Wajib: `DATABASE_URL`, `JWT_SECRET` (error menyebut nama var). **Satu-satunya pemanggil `os.Getenv`.** |
| `internal/platform/database/postgres.go` | `Connect(dsn) (*sql.DB, error)` & `MustConnect(dsn) *sql.DB`. Driver `database/sql` + `_ "github.com/jackc/pgx/v5/stdlib"` (ADR-012). Pool 25/25/5m, `PingContext` 5 dtk. |
| `internal/platform/httpx` | `response.go`: `JSON`, `Data` (envelope `{"data": …}`), `NoContent` (204). `errors.go`: `ErrorResponse` (`{"error":{"code","message"}}`), `Error(w,status,code,msg)` + shortcut `BadRequest/Unauthorized/Forbidden/NotFound/Conflict/Internal`; **registry** `RegisterErrorMapping(err, status, code)` + `WriteError(w, err)` (`errors.Is`, fallback `Internal` + log). `decode.go`: `Decode(w,r,dst)` (MaxBytes 1 MiB, `DisallowUnknownFields`, pesan ramah) & `DecodeAndValidate(w, r, dst, v *validator.Validator)` — keduanya **tak menulis response**, hanya kembalikan error. |
| `internal/platform/validator/validator.go` | `Validator` wrapper, `New()`, `Struct(any) error` (pesan `field: pesan; …`). Rule kustom `timezone` (IANA via `time.LoadLocation`, `""` → invalid). |
| `internal/platform/auth` | `password.go`: `HashPassword` / `ComparePassword` (bcrypt, default cost). `token.go`: `Claims{UserID, Timezone string; jwt.RegisteredClaims}` (tz di claims, ADR-013), interface **`Issuer`** & **`Verifier`** (nama `Verifier` — bukan `TokenVerifier`), `JWT` + `NewJWT(secret, ttl)`. HS256; `Verify` menolak kedaluwarsa, signature salah, `alg:none`/algoritma lain (`jwt.WithValidMethods`). |
| `cmd/api/main.go` | **Minimal** (sengaja): `config.Load()` → chi router 2 route health → `http.ListenAndServe`. Ada seam `// TODO(phase1/phase4): server.New(cfg, db)`. Belum ada koneksi DB, belum graceful shutdown. |
| `contracts/openapi.yaml` | **Penuh & valid** (redocly: 0 error, 3 warning benign). `servers` = `http://localhost:8080`; **semua path bisnis di-prefiks `/api/v1`**; `/healthz` & `/readyz` di luar prefiks. 16 operasi (11 asli + `GET /me`, `PATCH /me`, 2 health). 15 schema di `components.schemas`. `security: [{bearerAuth: []}]` di operasi terproteksi; `security: []` eksplisit di 4 operasi publik. **Ini acuan bentuk `dto.go` Phase 1–3** — kalau perlu ubah, ubah kontrak barengan (AGENTS.md #8). |
| Test | `internal/config/config_test.go`, `internal/platform/auth/password_test.go`, `internal/platform/auth/token_test.go`. Test pertama di repo. `make test` hijau. |
| `docs/DECISIONS.md` | **ADR-023** ditambahkan (registry error mapper). ADR-001..022 tak berubah. |

---

## Yang BELUM dikerjakan

| Hal | Kapan digarap |
|---|---|
| `internal/server/{server.go, router.go, health.go}` | **Terisi (Phase 1, minimal).** `server.New(cfg, db)` rakit `auth.NewJWT` + `user.New` + `buildRouter`. Router: chi middleware (RequestID/Logger/Recoverer), `/healthz` + `/readyz` (readyz sudah ping DB), `/api/v1` dengan split publik vs group ber-`Authenticator`. Graceful shutdown penuh & CORS di Phase 4 (T4.2–T4.4). |
| `internal/platform/middleware/middleware.go` | **Terisi (T1.8).** `Authenticator(verifier auth.Verifier)` + helper `WithUserID`/`UserIDFrom`/`WithTimezone`/`TimezoneFrom` (key tipe privat `ctxKey`). Isi userID **dan** timezone dari JWT claims ke context. |
| `internal/modules/user/**` | **Terisi (Phase 1).** 8 file + `service_test.go`. Semua rute lewat `RegisterPublicRoutes` / `RegisterProtectedRoutes`. `New(db, issuer)` daftarkan error mapping (ADR-023): `ErrEmailTaken`→409 `email_taken`, `ErrInvalidCredential`→401 `invalid_credential`, `ErrUserNotFound`→404 `user_not_found`. `Module.svc` disimpan untuk `AsUserDirectory()` (T3.4). |
| `internal/modules/quest/**`, `internal/modules/scoring/**` | Stub. Phase 2 & 3. |
| `internal/modules/achievement/**` | **Jangan disentuh** — v2 (ADR-010). Tetap tak di-mount. |
| Test untuk `httpx`, `validator`, `database` | Belum ada. Boleh ditambah saat menyentuhnya. |

---

## Keputusan mengikat yang lahir di Phase 0

- **ADR-023 — registry error mapper.** Tiap module mendaftarkan sentinel
  error-nya lewat `httpx.RegisterErrorMapping(err, status, code)` di dalam
  `New()` (yang dipanggil dari `server`). Handler cukup `httpx.WriteError(w, err)`.
  **Jangan** menulis tangga `errors.Is` per handler.
  Contoh isi yang diharapkan di Phase 1:
  `httpx.RegisterErrorMapping(user.ErrEmailTaken, 409, "email_taken")`,
  `…ErrInvalidCredential, 401, "invalid_credential"`,
  `…ErrUserNotFound, 404, "user_not_found"`.
- **`ErrorResponse` = `{"error":{"code":string,"message":string}}`** — `httpx`
  sudah cocok dengan kontrak; jaga tetap sinkron.
- **`PATCH /me` mengembalikan `AuthResponse`** (token baru), bukan `UserResponse`
  (ADR-022) — sudah tertulis di kontrak, tinggal diimplementasi di T1.10.

---

## Deviasi Phase 0 dari rencana (biar tak bingung saat baca kode)

1. **`go.mod` directive `go 1.25.0`**, bukan `1.23` — salah satu dependency
   menuntutnya. Tak ada baris `toolchain`. Build butuh Go ≥ 1.25.
2. **`internal/modules/quest/module.go` kena `gofmt -w`** (sisip 1 baris `//`
   pemisah di doc comment) supaya `make fmt` bersih se-repo. Nol perubahan
   logika/perilaku — isinya tetap stub.
3. **`config_test.go` sengaja `t.Setenv("APP_ENV"/"HTTP_PORT"/"JWT_TTL", "")`**
   di helper — `make test` meng-`export` isi `.env`, tanpa netralisasi ini test
   nilai default kepengaruh environment lokal.
4. **3 warning redocly tersisa** (server pakai `localhost`; `/healthz` & `/readyz`
   tanpa response 4xx). Benign; tak bisa dihilangkan tanpa melanggar spec/tugas.

---

## Cara Phase 1 dikerjakan (sudah selesai — untuk konteks pembaca kode)

- **Jebakan `GET /me` tanpa auth: sudah dibereskan.** `user/routes.go` dipecah
  jadi `RegisterPublicRoutes` (register, login) + `RegisterProtectedRoutes`
  (`GET /me`, `PATCH /me`); yang protected dipasang di group ber-`Authenticator`
  di `server/router.go`.
- **`RegisterRequest.Timezone` = opsional** (`validate:"omitempty,timezone"`),
  default `"Asia/Jakarta"` diisi di `service.Register` (`defaultTimezone`).
- **JSON tag `snake_case`** (`display_name`) ikut kontrak, bukan `displayName`
  dari stub.
- **`Login` = `ErrInvalidCredential` yang sama** untuk email tak ada maupun
  password salah.
- **`PATCH /me` menerbitkan token baru** (`AuthResponse`, ADR-022); email tak
  bisa diubah; `UpdateProfileRequest` pakai pointer field.
- **UUIDv7 digenerate di `service`** (`github.com/google/uuid` — sekarang direct
  dep di `go.mod`), bukan di SQL.
- **Deteksi email dobel** di repo lewat `pgconn.PgError.Code == "23505"`
  (`errors.As`), bukan `strings.Contains`.
- `handler` bikin `validator.New()` sendiri lewat `user.New` (signature `New(db,
  issuer)` dipertahankan sesuai T1.7 — validator tak ikut parameter).
- **ADR-024 (baru): amplop response sukses `{"data": ...}`.** Handler sukses
  wajib lewat `httpx.Data` (bukan `httpx.JSON` mentah). `register` pakai `200`
  (bukan `201`) sesuai kontrak. Kontrak sudah dibungkus untuk semua endpoint
  (termasuk Phase 2/3). ADR count kini **25**.

### BELUM dilakukan untuk Phase 1

1. **Smoke test curl end-to-end** (butuh Postgres; Docker daemon mati saat
   implementasi). Jalankan dulu sebelum lanjut Phase 2 — lihat bawah.
2. **Commit Phase 1.**

---

## Perangkap yang menunggu di Phase 2

- **`ErrNotOwner` → HTTP 404**, bukan 403 — jangan bocorkan bahwa quest itu
  milik orang lain (daftarkan mapping-nya begitu).
- **Perangkap chi:** daftarkan `/quests/today` **sebelum** `/quests/{questId}`.
- **`quest.service` terima `localDate` sebagai argumen** — dilarang panggil
  `time.Now()`. Handler hitung "hari ini" dari `middleware.TimezoneFrom`
  (fallback `Asia/Jakarta`).
- **Port `ScoreAwarder` milik `quest`** (bukan `scoring`); Phase 2 inject
  `noopScorer{}` sementara.

---

## Cara memulai sesi berikutnya (smoke test Phase 1 + Phase 2)

```bash
# dari root repo — nyalakan Docker Desktop dulu
docker compose -f docker-compose.dev.yml up -d          # Postgres

cd apps/backend
make migrate-up                                          # 6 tabel
make fmt && make vet && make build && make test          # harus hijau

make run &                                               # server :8080
# smoke test Phase 1:
TOKEN=$(curl -sX POST localhost:8080/api/v1/auth/register \
  -d '{"email":"a@b.c","password":"rahasia123","display_name":"A"}' | sed 's/.*"token":"\([^"]*\)".*/\1/')
curl -sX POST localhost:8080/api/v1/auth/login -d '{"email":"a@b.c","password":"rahasia123"}'
curl -si  localhost:8080/api/v1/me                                   # -> 401
curl -si  localhost:8080/api/v1/me -H "Authorization: Bearer $TOKEN" # -> 200
curl -sX PATCH localhost:8080/api/v1/me -H "Authorization: Bearer $TOKEN" \
  -d '{"timezone":"Asia/Makassar"}'                                  # -> AuthResponse, token baru
curl -s localhost:8080/api/v1/me -H "Authorization: Bearer $TOKEN" | grep -i password  # -> kosong
```

Kalau hijau: commit Phase 1, lalu buka `docs/tasks/phase-02-quest.md` dan
kerjakan **T2.1 → T2.10** berurutan.

---

## Verifikasi baseline Phase 0 (kalau ragu masih sehat)

| Cek | Harapan |
|---|---|
| `make fmt && make vet && make build && make test` | bersih / PASS |
| `make migrate-up` lalu `make migrate-down` lalu `make migrate-up` | tanpa error, 6 tabel |
| `grep -c TODO contracts/openapi.yaml` | `0` |
| `npx @redocly/cli@latest lint contracts/openapi.yaml` | valid (3 warning benign) |
| `grep -rn "questday/internal/modules" internal/platform/` | kosong (platform netral) |
| `grep -c "^## ADR-" docs/DECISIONS.md` | `25` (ADR-001..024 + template) |

# Hands-off — State Backend QuestDay

Dokumen orientasi untuk agent/kolaborator berikutnya yang mengerjakan
`apps/backend`. **Baca sesudah `AGENTS.md` (root), sebelum menyentuh kode.**

**Terakhir diperbarui:** 2026-09-01
**Status:** Phase 0/1 (`5c7ddb2`), Phase 2 (`0e408ab`), Phase 3 (`7a1b57d`)
ter-commit. **Phase 4 (Hardening) selesai → MVP BACKEND DONE** (kode + unit test
+ smoke test end-to-end penuh lolos; **belum di-commit**). **STOP di sini —
jangan mulai `achievement`/v2 tanpa diminta (AGENTS.md).**

> Root `docs/HANDS-OFF.md` tetap jadi orientasi seluruh proyek (monorepo, titik
> sambung frontend). File ini fokus ke state backend saja.

---

## TL;DR

- Phase 0 + Phase 1 ter-commit di `5c7ddb2` (branch `feature/backend/m`, sudah
  di-merge dgn `main`). Endpoint hidup:
  `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /api/v1/me`,
  `PATCH /api/v1/me`. `make -C apps/backend fmt vet build test` bersih; unit
  test service user hijau; smoke test end-to-end lolos.
- **Phase 2:** module `quest` lengkap. Endpoint (semua terproteksi):
  `GET/POST /api/v1/quests`, `GET /api/v1/quests/today`,
  `PATCH/DELETE /api/v1/quests/{questId}`,
  `POST /api/v1/quests/{questId}/complete|uncomplete`.
- **Phase 3:** module `scoring` lengkap + **gamifikasi hidup**. Endpoint
  (terproteksi): `GET /api/v1/me/score`, `GET /api/v1/me/streak`,
  `GET /api/v1/leaderboard?limit=`. `noopScorer` **dihapus** — `quest.New` kini
  dapat `scoringMod.AsScoreAwarder()`. Complete quest menaikkan poin/XP/level +
  streak; uncomplete rollback poin (streak tetap, ADR-009). ADR-016 diputuskan:
  kompensasi manual (hapus log yatim bila scorer gagal).
- **Phase 4 (Hardening):** graceful shutdown (`signal.NotifyContext` di
  `main.go`, `Server.Shutdown` = HTTP drain + `db.Close`), `http.Server` timeouts
  lengkap, **CORS** via `go-chi/cors` (env `CORS_ALLOWED_ORIGINS`, **ADR-027**),
  `chimw.Timeout(30s)`. Validasi body gagal → code **`validation_failed`**
  (bukan `bad_request`). `chimw.RealIP` sengaja tak dipakai (IP-spoofing).
  `Dockerfile` → `golang:1.25-alpine`.
- `contracts/openapi.yaml`: response 2xx ber-envelope `{"data": ...}`
  (**ADR-025**); `QuestResponse.points` (**ADR-026**). `getMyScore`/`getMyStreak`
  **tak lagi punya `404`** (wallet/streak auto-default). 0 TODO, redocly valid
  (3 warning benign). Cocok 1:1 dengan implementasi.
- **Infra dev lewat root `Makefile`** (ADR-024): `make up` / `make backend`.
  Kredensial Postgres di `.env.docker` (gitignored). **Catatan lokal:**
  container `questday_postgres` map host **5433**, `apps/backend/.env`
  `HTTP_PORT=8001` → server `:8001`, DB `localhost:5433`.
- **MVP backend selesai.** Berikutnya cuma: commit Phase 4. Sesudah itu tak ada
  task backend lagi kecuali diminta (v2: `achievement`, event bus, streak
  freeze, dll — lihat `docs/tasks/README.md`).

---

## Peta baca (urutan)

1. `AGENTS.md` (root) — aturan keras & batasan.
2. **File ini** — apa yang sudah/belum ada.
3. `docs/ARCHITECTURE.md` — bentuk sistem, aturan dependensi.
4. `docs/DECISIONS.md` — 28 ADR. Relevan backend: **ADR-016** (atomicity =
   kompensasi manual), **ADR-023** (registry error mapper), **ADR-024** (root
   Makefile + `.env.docker`), **ADR-025** (amplop `{"data":...}`), **ADR-026**
   (`QuestResponse.points`), **ADR-027** (CORS via `go-chi/cors`).
5. `apps/backend/docs/tasks/README.md` — semua phase 0–4 selesai; sisanya v2.
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

## Isi tiap area (state akhir MVP)

| Hal | Kapan digarap |
|---|---|
| `internal/server/{server.go, router.go, health.go}` + `cmd/api/main.go` | **Final (Phase 4).** `server.New` rakit user→scoring→quest; `http.Server` 4 timeout; `Server.Shutdown` = HTTP drain + `db.Close` (`errors.Join`). Router: `RequestID/Logger/Recoverer/Timeout(30s)` + CORS (kalau origin di-set) + `/healthz`+`/readyz` + `/api/v1` (publik vs `Authenticator` group). `main.go`: `signal.NotifyContext` → goroutine `ListenAndServe` → `Shutdown` 10s. |
| `internal/platform/middleware/middleware.go` | **Terisi (T1.8).** `Authenticator(verifier auth.Verifier)` + helper `WithUserID`/`UserIDFrom`/`WithTimezone`/`TimezoneFrom` (key tipe privat `ctxKey`). Isi userID **dan** timezone dari JWT claims ke context. |
| `internal/modules/user/**` | **Terisi (Phase 1).** 8 file + `service_test.go`. Semua rute lewat `RegisterPublicRoutes` / `RegisterProtectedRoutes`. `New(db, issuer)` daftarkan error mapping (ADR-023): `ErrEmailTaken`→409 `email_taken`, `ErrInvalidCredential`→401 `invalid_credential`, `ErrUserNotFound`→404 `user_not_found`. `Module.svc` disimpan untuk `AsUserDirectory()` (T3.4). |
| `internal/modules/quest/**` | **Terisi (Phase 2).** 8 file + `domain_test.go` + `service_test.go`. Port `ScoreAwarder` (milik quest, ADR-005) — kini disuntik `scoring` sungguhan. Error mapping (ADR-023): `ErrQuestNotFound`→404 `quest_not_found`, `ErrNotOwner`→**404 `quest_not_found`** (kode sama, jangan bocor), `ErrAlreadyCompleted`→409 `already_completed`, `ErrNotCompleted`→409 `not_completed`. Poin cuma di `domain.go` (5/10/20). `CompleteQuest` kompensasi ADR-016: hapus log bila `OnQuestCompleted` gagal. |
| `internal/modules/scoring/**` | **Terisi (Phase 3).** 8 file + `domain_test.go` + `service_test.go`. `LevelForXP`/`XPForLevel`/`NextStreak` = satu-satunya sumber kurva level & aturan streak. Port **`UserDirectory`** (milik scoring, ADR-014) diimplementasi `user.AsUserDirectory()` (`ListNamesByIDs`, `WHERE id = ANY($1)`). `New(db, dir)`, `AsScoreAwarder()`. SQL scoring **tak** JOIN `users`. Leaderboard nama fallback `"(pengguna dihapus)"`; `?limit=` clamp [1,100] default 20. Tak ada error domain (wallet/streak auto-default). |
| `internal/server/scorer.go` | **DIHAPUS di Phase 3** (`noopScorer` sudah tak ada). `server.New` urutan: user → scoring → quest. |
| `internal/modules/scoring/**` | Stub. Phase 3. |
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
- **ADR-025: amplop response sukses `{"data": ...}`.** Handler sukses wajib lewat
  `httpx.Data` (bukan `httpx.JSON` mentah). `register` pakai `200` (bukan `201`)
  sesuai kontrak. Kontrak sudah dibungkus untuk semua endpoint (termasuk
  Phase 2/3). *(Merge `main` membawa ADR-024 = root Makefile; ADR "amplop" jadi
  ADR-025.)*
- **Smoke test end-to-end Phase 1: LOLOS** (register→login→`GET /me` 401 tanpa
  token / 200 dgn token → `PATCH /me` menerbitkan token baru ber-`tz` baru →
  tak ada field password di response → password salah = 401 `invalid_credential`).

---

## Cara Phase 2 dikerjakan (sudah selesai — untuk konteks pembaca kode)

- **`ErrNotOwner` → 404 `quest_not_found`** (kode sama dgn not-found; tak bocor).
  Ekstra: `repo.GetQuest` sudah difilter `user_id` → praktis quest orang lain
  langsung `ErrQuestNotFound`; cek `q.UserID != userID` di service tetap ada
  sebagai pertahanan kedua & supaya testable dgn fake repo.
- **`localDateFrom(r)` di `quest/handler.go`** menghitung tanggal lokal user dari
  `middleware.TimezoneFrom` (fallback `Asia/Jakarta`), lalu `time.Date(...,
  time.UTC)` jam 00:00 → dipakai sebagai kolom `date`. `quest/service.go`
  menerima `localDate` sebagai argumen, **nol `time.Now()`**.
- **Port `ScoreAwarder`** didefinisikan di `quest/service.go` (milik quest).
  Phase 2 sempat pakai `noopScorer` di `server`; **dihapus di Phase 3**.
- **`quest.Points()` = satu-satunya sumber poin** (5/10/20). `QuestResponse.points`
  & `QuestLog.PointsAwarded` diisi dari sana. Request tak pernah kirim poin.
- **Status kode per kontrak:** `POST /quests` → 201; `PATCH` → 200;
  `complete` → 200 `QuestLogResponse`; `uncomplete` & `DELETE` → 204.
- **`quest.New(db, scorer)`** — `validator.New()` dibuat di dalamnya (pola sama
  `user.New`).
- **Smoke test end-to-end Phase 2: LOLOS** — create(points=10)→today(false)→
  list→complete(200)→today(true)→complete lagi(409 `already_completed`)→
  uncomplete(204)→patch difficulty=hard(points→20)→archive(204)→today items
  kosong→quest user lain complete = 404 `quest_not_found`.

---

## Cara Phase 3 dikerjakan (sudah selesai — untuk konteks pembaca kode)

- **Kurva level** (`scoring/domain.go`, satu-satunya sumber, ADR-007):
  `XPForLevel(n) = 25*n*(n+1) - 50` → ambang `0/100/250/450/700/1000/…`.
  `LevelForXP` loop naik dari 1, tanpa cap. `PointsToNextLevel(xp)` diturunkan.
- **XP == TotalPoints** untuk MVP; keduanya naik/turun bersama. `Level` selalu
  dari `LevelForXP`. `OnQuestUncompleted` clamp `TotalPoints`/`XP` ≥ 0 dan
  **tak menyentuh streak** (ADR-009).
- **`NextStreak`** murni (tak `time.Now()`): hari sama → tetap; +1 hari → +1;
  bolong/pertama → 1; `Longest` tak pernah turun.
- **Port `UserDirectory`** (milik `scoring`, ADR-014) diimplementasi
  `user.AsUserDirectory()` → `user.service.NamesByIDs` →
  `repo.ListNamesByIDs` (`SELECT id, display_name FROM users WHERE id = ANY($1)`,
  `[]string` langsung — pgx v5 stdlib meng-encode ke array). `scoring` & `user`
  tak saling import; SQL scoring tak sebut `users`.
- **ADR-016 = kompensasi manual.** `quest.CompleteQuest`: kalau
  `scorer.OnQuestCompleted` gagal → `repo.DeleteLog(...)` lalu return error.
  Dikunci test `TestCompleteQuest_ScorerFails_NoOrphanLog`.
- **Server wiring** (`server.New`): `user.New` → `scoring.New(db,
  userMod.AsUserDirectory())` → `quest.New(db, scoringMod.AsScoreAwarder())`.
  `noopScorer` + `internal/server/scorer.go` dihapus.
- **`/leaderboard`** tetap terproteksi (MVP). `?limit=` default 20, clamp
  [1,100]. Nama ID tak dikenal → `"(pengguna dihapus)"`.
- **Smoke test end-to-end Phase 3: LOLOS** — score/streak default
  `{0,0,1,ptl:100}` / `{0,0,null}` → complete(hard=20) → score
  `{20,20,1,ptl:80}`, streak `{1,1,"2026-09-01"}` → leaderboard
  `[{rank:1,"P3 Alice",20}]` → uncomplete → score balik `{0,0,1,100}`, **streak
  tetap** `{1,1,…}` → `?limit=999` clamp, 200.

---

## Cara Phase 4 dikerjakan (MVP DONE — untuk konteks pembaca kode)

- **Graceful shutdown:** `cmd/api/main.go` — `signal.NotifyContext(SIGINT,
  SIGTERM)`, `ListenAndServe` di goroutine (abaikan `http.ErrServerClosed`),
  `<-ctx.Done()` → `srv.Shutdown(ctx 10s)`. `Server.Shutdown` = HTTP drain +
  `db.Close()` digabung `errors.Join`. `main` tak lagi `defer db.Close()`.
- **`http.Server` timeouts:** ReadHeader 5s, Read 15s, Write 15s, Idle 60s.
- **Middleware stack** (`buildRouter`): `RequestID`, `Logger`, `Recoverer`,
  `Timeout(30s)`. **`chimw.RealIP` sengaja TIDAK dipakai** — IP-spoofing tanpa
  trusted proxy (GHSA-3fxj-6jh8-hvhx). Komentar ada di `router.go`.
- **CORS (ADR-027):** `github.com/go-chi/cors`. Aktif hanya bila
  `config.CORSAllowedOrigins` tak kosong (env `CORS_ALLOWED_ORIGINS`,
  comma-separated). `AllowCredentials:false` (token di header, ADR-020).
  `config.splitList` parser; `config_test.go` menguji parse & kosong.
- **`validation_failed`:** `httpx.ValidationFailed(w, msg)` (400) baru; semua
  5 titik `DecodeAndValidate` gagal di `user`/`quest` handler pakai ini
  (dulu `httpx.BadRequest` = `bad_request`). Sekarang cocok dgn enum
  `ErrorResponse` di kontrak.
- **Kontrak (T4.6):** `getMyScore` & `getMyStreak` **kehilangan blok `404`**
  (impl auto-default, tak pernah 404). 13 rute bisnis impl = 13 path kontrak;
  `bearerAuth` hanya di operasi grup protected. 0 TODO. `servers` tetap
  `http://localhost:8080` (contoh; warning benign).
- **`Dockerfile`:** `FROM golang:1.25-alpine AS build` (cocok `go 1.25.0`).
- **Smoke test MVP penuh: LOLOS** — CORS preflight (`Access-Control-Allow-Origin`)
  → register/login → `validation_failed` untuk body jelek → 2 quest (easy=5,
  hard=20) → today 2×false → complete hard → score `{20,20,1,ptl:80}` +
  streak `{1,1}` → leaderboard `[{rank:1,"MVP User",20}]` → uncomplete → score
  `{0,0,1,100}`, streak tetap → `PATCH /me` tz baru (token baru) → today OK →
  no password leak → **SIGTERM → shutdown bersih** ("sinyal diterima" →
  "selesai" → exit).

---

## Kalau lanjut kerja backend (semua phase MVP sudah selesai)

- Baseline sehat: `make up` → `make -C apps/backend migrate-up` →
  `make -C apps/backend fmt vet build test` (harus PASS) → `make backend`.
- Lokal: server `:8001`, DB `localhost:5433` (`apps/backend/.env`). Kalau
  `make up` bilang "container name in use", Postgres sudah jalan — lanjut saja.
- **v2 backlog** (jangan dikerjakan tanpa diminta — AGENTS.md): module
  `achievement`, in-process event bus, streak freeze/grace, streak rollback saat
  uncomplete, leaderboard paging/periode/cache, atomicity tx sungguhan
  (supersede ADR-016), CORS default aktif. Detail di
  `apps/backend/docs/tasks/README.md`.

---

## Cara memulai sesi (infra + baseline)

```bash
# dari root repo
make up                                  # Postgres (docker compose, --env-file .env.docker)
make -C apps/backend migrate-up          # 6 tabel
make -C apps/backend fmt vet build test  # harus hijau

make backend                             # server via air (port dari apps/backend/.env)
```

Lokal saat ini: DB `localhost:5433`, server `:8001` (lihat `apps/backend/.env`).
Kalau `make up` gagal "container name in use", container Postgres sudah jalan —
lanjut saja. Smoke test Phase 1 & Phase 2 ada di file tugas masing-masing.

---

## Verifikasi baseline (kalau ragu masih sehat)

| Cek | Harapan |
|---|---|
| `make -C apps/backend fmt vet build test` | bersih / PASS |
| `make -C apps/backend migrate-up` lalu `migrate-down` lalu `migrate-up` | tanpa error, 6 tabel |
| `grep -c TODO contracts/openapi.yaml` | `0` |
| `npx @redocly/cli@latest lint contracts/openapi.yaml` | valid (3 warning benign) |
| `grep -rn "questday/internal/modules" internal/platform/` | kosong (platform netral) |
| `grep -c "^## ADR-" docs/DECISIONS.md` | `28` (ADR-001..027 + template) |

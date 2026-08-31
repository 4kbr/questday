# Phase 4 — Hardening & Kontrak

**Tujuan:** merapikan sisa-sisa yang membuat backend layak disebut selesai:
health check, graceful shutdown, error mapping seragam, dan **`contracts/openapi.yaml`
yang benar-benar cocok dengan implementasi**. Setelah phase ini → **MVP DONE**.

**Prasyarat:** Phase 3 selesai (semua endpoint MVP sudah jalan).

**Catatan soal kontrak:** `contracts/openapi.yaml` sudah diisi lengkap di
**T0.15 (Phase 0)** — sesuai aturan "kontrak lebih dulu" (AGENTS.md #8, ADR-002)
dan karena `apps/frontend` menggenerate type-nya dari sana (ADR-019). Yang
tersisa di phase ini hanyalah **memverifikasi** bahwa implementasi akhir benar-
benar cocok dengan kontrak itu (T4.6).

---

## T4.1 — `server/health.go`

- **Sentuh:** `internal/server/health.go`, `internal/server/router.go`
- **Isi:**
  ```go
  func healthHandler() http.HandlerFunc               // 200 "ok" — liveness
  func readyHandler(db *sql.DB) http.HandlerFunc      // db.PingContext — readiness
  ```
  Daftarkan `GET /healthz` dan `GET /readyz` **di luar** prefix `/api/v1` dan di
  luar grup terproteksi.
- **Aturan:** `readyHandler` pakai context bertimeout supaya tidak menggantung
  saat DB lambat; balas 503 kalau gagal.
- **DoD:** kedua endpoint jalan tanpa token.
- **Verifikasi:** `curl -i localhost:8080/healthz; curl -i localhost:8080/readyz`

## T4.2 — `server/router.go` final

- **Sentuh:** `internal/server/router.go`
- **Isi:** `buildRouter` lengkap:
  ```go
  r.Use(chimw.RequestID, chimw.RealIP, chimw.Logger, chimw.Recoverer)
  r.Use(chimw.Timeout(30 * time.Second))
  // CORS bila frontend sudah ada (baca origin dari config)
  r.Get("/healthz", ...); r.Get("/readyz", ...)
  r.Route("/api/v1", func(r chi.Router) {
      userMod.RegisterPublicRoutes(r)
      r.Group(func(r chi.Router) {
          r.Use(middleware.Authenticator(jwt))
          userMod.RegisterProtectedRoutes(r)
          questMod.RegisterRoutes(r)
          scoringMod.RegisterRoutes(r)
      })
  })
  ```
- **Aturan:** middleware generik **pakai bawaan chi**, jangan tulis sendiri.
  Milik kita hanya `Authenticator`.
- **DoD:** rute publik & terproteksi terpisah benar; panic tak menjatuhkan proses.
- **Verifikasi:** `/me` tanpa token → 401; `/auth/login` tanpa token → 200/401
  (bukan 401 karena middleware).

## T4.3 — `server/server.go` final

- **Sentuh:** `internal/server/server.go`
- **Isi:**
  ```go
  type Server struct{ httpServer *http.Server; db *sql.DB }
  func New(cfg config.Config, db *sql.DB) *Server
  func (s *Server) ListenAndServe() error
  func (s *Server) Shutdown(ctx context.Context) error   // http shutdown + db.Close
  ```
  Set `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout` pada
  `http.Server` (jangan biarkan default tak terbatas).
- **Aturan:** `server` adalah composition root — satu-satunya yang tahu semua
  module. Module tetap tidak saling kenal.
- **DoD:** `Shutdown` menutup HTTP lalu DB, tanpa kebocoran.
- **Verifikasi:** `make build`

## T4.4 — `cmd/api/main.go`: graceful shutdown

- **Sentuh:** `cmd/api/main.go`
- **Isi:**
  ```go
  cfg, err := config.Load()          // tangani error
  db  := database.MustConnect(cfg.DatabaseURL)
  srv := server.New(cfg, db)
  go srv.ListenAndServe()            // log error kecuali http.ErrServerClosed
  // tunggu SIGINT/SIGTERM (signal.NotifyContext)
  // ctx timeout ~10s -> srv.Shutdown(ctx)
  ```
- **Catatan:** komentar scaffold menulis `cfg := config.Load()` (mengabaikan
  error) padahal `Load` mengembalikan `(Config, error)` — selaraskan, jangan ikut
  komentarnya.
- **DoD:** Ctrl+C menghentikan server tanpa memutus request yang sedang jalan.
- **Verifikasi:** `make run`, kirim request panjang, tekan Ctrl+C — request
  selesai, proses keluar bersih.

## T4.5 — Error mapping seragam lintas module

- **Sentuh:** `internal/platform/httpx/errors.go`, semua `handler.go`
- **Isi:** pastikan setiap error domain punya pasangan status + `code` yang
  konsisten:

  | Error | Status | `code` |
  |---|---|---|
  | `user.ErrEmailTaken` | 409 | `email_taken` |
  | `user.ErrInvalidCredential` | 401 | `invalid_credential` |
  | `user.ErrUserNotFound` | 404 | `user_not_found` |
  | `quest.ErrQuestNotFound` | 404 | `quest_not_found` |
  | `quest.ErrNotOwner` | 404 | `quest_not_found` (sengaja disamarkan) |
  | `quest.ErrAlreadyCompleted` | 409 | `already_completed` |
  | `quest.ErrNotCompleted` | 409 | `not_completed` |
  | validasi gagal | 400 | `validation_failed` |
  | lain-lain | 500 | `internal_error` |

- **Aturan:** `ErrNotOwner` sengaja dipetakan ke 404 dengan code yang sama seperti
  not-found — jangan bocorkan bahwa resource itu ada milik orang lain. 500 tak
  pernah memuat detail internal di body (log saja).
- **DoD:** tak ada handler yang menulis status/code sendiri di luar tabel ini.
- **Verifikasi:** `grep -rn "http.Error\|WriteHeader(5" internal/modules/` → kosong.

## T4.6 — Verifikasi kontrak vs implementasi

- **Sentuh:** `contracts/openapi.yaml` (hanya bila ada yang melenceng)
- **Isi:** kontrak sudah lengkap sejak T0.15. Di sini dicocokkan dengan hasil
  akhir implementasi:
  1. **Setiap endpoint yang jalan ada di kontrak, dan sebaliknya** — tak ada
     rute yang terlupa, tak ada operasi yang didokumentasikan tapi tak
     diimplementasi.
  2. **Bentuk field cocok**: nama, tipe, dan required setiap schema sama dengan
     DTO di `dto.go` masing-masing module.
  3. **Status & `code` error** yang benar-benar dikembalikan handler cocok dengan
     yang tercantum di kontrak (bandingkan dengan tabel di T4.5).
  4. **`security: bearerAuth`** terpasang persis di operasi yang memang berada di
     grup terproteksi `router.go` — tak lebih, tak kurang.
  5. Envelope sukses `{"data": ...}` dan error `{"error":{code,message}}` sesuai
     `ErrorResponse` di kontrak.
- **Aturan:** kalau ada yang melenceng, **perbaiki kontrak dan/atau
  implementasi** — dan beri tahu frontend, karena mereka harus menjalankan
  `npm run gen:api` lagi (ADR-019). Jangan biarkan salah satunya "benar sendiri".
- **DoD:** nol selisih antara kontrak dan implementasi; `grep -c TODO
  contracts/openapi.yaml` → 0.
- **Verifikasi:** `npx @redocly/cli lint contracts/openapi.yaml` bersih, lalu
  bandingkan daftar rute:
  ```bash
  grep -rn "r.Get\|r.Post\|r.Patch\|r.Delete" internal/modules/ internal/server/
  grep -n "^  /" contracts/openapi.yaml
  ```

## T4.7 — `Dockerfile` + `.env.example`

- **Sentuh:** `apps/backend/Dockerfile`, `apps/backend/.env.example`
- **Isi:**
  - `Dockerfile`: sesuaikan versi Go (cocokkan `go.mod`) & module path
    (`questday`, ADR-011). `COPY go.mod go.sum ./` sekarang aman karena `go.sum`
    lahir di T0.2.
  - `.env.example`: tambah variabel yang muncul selama implementasi (mis.
    `LOG_LEVEL`, `CORS_ALLOWED_ORIGINS`) — dan pastikan `Config` (T0.5) ikut
    diperbarui.
- **Aturan:** jangan pernah menaruh nilai asli di `.env.example` — hanya
  placeholder.
- **DoD:** `docker build -t questday-api apps/backend` sukses.
- **Verifikasi:** `docker build` + jalankan container melawan Postgres dev.

## T4.8 — Smoke test end-to-end MVP

- **Sentuh:** — (verifikasi; boleh dituangkan jadi skrip
  `apps/backend/scripts/smoke.sh` kalau berguna)
- **Isi:** satu alur penuh dari database kosong:
  ```
  make migrate-down && make migrate-up
  register -> login -> buat 2 quest (easy, hard)
  GET /quests/today          -> 2 item, completed: false
  complete quest hard        -> 200
  GET /quests/today          -> 1 completed
  GET /me/score              -> poin = Points(hard), level sesuai LevelForXP
  GET /me/streak             -> current: 1, longest: 1
  GET /leaderboard           -> user muncul dengan display_name
  uncomplete quest hard      -> 200
  GET /me/score              -> poin kembali 0
  GET /me/streak             -> current tetap 1 (ADR-009)
  PATCH /me {"timezone":"Asia/Makassar"} -> token baru, claims.Timezone berubah
  GET /quests/today          -> tanggal mengikuti timezone baru
  ```
- **DoD:** seluruh alur di atas sesuai harapan.

---

## Exit criteria Phase 4 — **MVP BACKEND DONE**

- [ ] `make fmt && make vet && make build && make test` bersih.
- [ ] `make migrate-up` dari DB kosong sukses; smoke test T4.8 lulus penuh.
- [ ] `/healthz` & `/readyz` jalan; Ctrl+C shutdown bersih.
- [ ] Envelope response & error seragam di seluruh endpoint.
- [ ] `contracts/openapi.yaml` nol TODO dan cocok 1:1 dengan implementasi
      (termasuk `PATCH /me`).
- [ ] `docker build` sukses.
- [ ] Tak ada import lintas-module, tak ada SQL di luar `repository_postgres.go`,
      tak ada `time.Now()` mentah untuk logika harian.
- [ ] `PasswordHash` tak pernah muncul di response.
- [ ] Semua keputusan baru selama implementasi tercatat di `docs/DECISIONS.md`.

Setelah ini: **jangan lanjut ke `achievement` atau fitur v2 tanpa diminta**
(AGENTS.md). Backlog v2 ada di [README.md](README.md).

# Hands-off — State Proyek QuestDay

Dokumen orientasi cepat untuk agent/kolaborator berikutnya. **Baca ini dulu**,
lalu masuk ke dokumen yang relevan dengan tugasmu.

**Terakhir diperbarui:** 2026-08-31

---

## TL;DR

QuestDay = pencatat quest harian dengan gamifikasi (poin, XP, level, streak).
Monorepo: backend Go (modular monolith) + frontend React (belum ada).

**State sekarang: perencanaan selesai, implementasi belum dimulai.**

- Backend = **scaffold murni** — 42 file `.go` yang isinya hanya `package x` +
  komentar TODO. Nol tipe, nol fungsi, nol dependency.
- Frontend = **kosong** — hanya `README.md` dan `docs/tasks/`.
- Yang sudah jadi: **arsitektur, 22 ADR, dan 109 task terperinci** yang dibagi
  per phase untuk kedua app.

**Kalau kamu diminta mulai koding:** buka `apps/backend/docs/tasks/README.md`,
kerjakan **T0.1** dan seterusnya berurutan. Jangan lompat.

---

## Peta dokumen — baca sesuai kebutuhan

| Dokumen | Isi | Kapan dibaca |
|---|---|---|
| `AGENTS.md` (root) | Aturan main & batasan untuk agent | **Selalu, pertama** |
| `docs/ARCHITECTURE.md` | Bentuk sistem, lapisan, aturan dependensi | Sebelum menyentuh struktur |
| `docs/DECISIONS.md` | 22 ADR — kenapa segala sesuatunya begini | Sebelum mengubah keputusan |
| `docs/GUIDES.md` | Resep: nambah endpoint/module/migrasi | Saat mengerjakan tugas umum |
| `contracts/openapi.yaml` | Kontrak API — sumber kebenaran bersama | Saat menyentuh endpoint |
| `apps/backend/docs/tasks/` | 55 task backend, 5 phase | Saat ngoding backend |
| `apps/frontend/docs/tasks/` | 54 task frontend, 5 phase | Saat ngoding frontend |

---

## Kondisi tiap bagian

### `apps/backend` — scaffold, belum diimplementasi

| Hal | Kondisi |
|---|---|
| File `.go` | 42 file, **semuanya hanya `package x` + komentar TODO** |
| `go.mod` | module path masih placeholder `github.com/yourorg/questday`, **nol dependency** |
| `go.sum` | **belum ada** — karena itu `Dockerfile` akan gagal di tahap `COPY` |
| `migrations/000001_init.up.sql` | **hanya komentar, nol DDL** |
| Test | **tidak ada satu pun `*_test.go`** (padahal `make test` sudah pakai `-race -cover`) |
| `Makefile`, `.air.toml`, `.env.example` | sudah lengkap & siap pakai |

### `apps/frontend` — belum ada apa-apa

Hanya `README.md` dan `docs/tasks/`. Belum ada `package.json`, belum ada `src/`.

### `contracts/openapi.yaml` — kerangka

11 operasi terdaftar (hanya `tags` + `summary`), `components.schemas` masih `{}`,
`securitySchemes.bearerAuth` sudah didefinisikan tapi **belum dipasang** di satu
operasi pun. 14 marker TODO. Mengisinya adalah **backend T0.15**.

---

## Urutan kerja yang direncanakan

```
BACKEND                                FRONTEND
Phase 0  Foundation (15 task)  ──T0.15──►  Phase 0  Setup (12 task)
  go.mod, migrasi, config,     kontrak     vite, tailwind, shadcn,
  platform/*, ISI KONTRAK                  gen:api, client.ts, MSW
     │                                        │
Phase 1  User & Auth (11)      ──T1.10──►  Phase 1  Auth & shell (11)
     │                         PATCH /me      │
Phase 2  Quest (10)                        Phase 2  Quest (12)
     │                                        │
Phase 3  Scoring (11)                      Phase 3  Scoring (9)
     │                                        │
Phase 4  Hardening (8)                     Phase 4  Polish (10)
     ▼                                        ▼
   MVP BACKEND DONE                       MVP FRONTEND DONE
```

**Dua app bisa dikerjakan paralel**, dengan dua titik sambung:

| Frontend butuh | Dari backend | Kalau belum ada |
|---|---|---|
| `contracts/openapi.yaml` terisi | **T0.15** | frontend F0.4 memblokir total |
| `PATCH /me` | **T1.10** | frontend F4.1–F4.3 memblokir |

Selain dua itu, frontend pakai mock (MSW, `VITE_USE_MOCK=true`) sampai endpoint
terkait jadi.

---

## Keputusan yang mengikat (jangan ditabrak diam-diam)

Semua ada di `docs/DECISIONS.md` dengan konteks & konsekuensinya. Yang paling
sering menggigit:

| ADR | Keputusan | Konsekuensi praktis |
|---|---|---|
| 004 | `Quest` (template) ≠ `QuestLog` (penyelesaian harian) | `UNIQUE(quest_id, date)` mencegah dobel |
| 005 | Antar-module lewat **port**, bukan import | `quest.ScoreAwarder` diimplementasi `scoring` |
| 006 | "Hari ini" = **timezone user** | **Dilarang `time.Now()` mentah** untuk logika harian |
| 007 | Poin di `Quest.Points()`, level di `LevelForXP` | Jangan sebar angka ajaib |
| 008 | Streak reset saat bolong, belum ada freeze | Logika terpusat di `NextStreak` |
| 009 | Uncomplete rollback poin, **streak dibiarkan** | Ini disengaja, bukan bug |
| 010 | `achievement` = v2 | **Jangan dikerjakan** kecuali diminta |
| 011 | Module path Go = `questday` | Import: `questday/internal/...` |
| 012 | `database/sql` + `pgx/v5/stdlib` | Semua repo terima `*sql.DB` |
| 013 | Timezone ikut di **JWT claims** | Ubah timezone → **wajib token baru** |
| 014 | Leaderboard pakai port `UserDirectory` | SQL scoring **tak boleh** JOIN ke `users` |
| 015 | PK `uuid`, UUIDv7 digenerate di app | Tak butuh `pgcrypto`, tak butuh `RETURNING id` |
| 016 | **(dipesan, belum diputuskan)** | Atomicity log+poin — diputuskan di backend T3.10 |
| 018 | TanStack Query (server state) vs Zustand (client state) | **Dilarang simpan hasil API ke Zustand** |
| 019 | Type frontend **digenerate** dari kontrak | `schema.gen.ts` tak pernah diedit tangan |
| 020 | Token di localStorage (MVP) | Risiko XSS diterima sadar; v2 → refresh token |
| 022 | `PATCH /me` masuk MVP | Mengembalikan `AuthResponse`, bukan `UserResponse` |

---

## Aturan keras yang paling sering dilanggar

**Backend:**
1. `domain.go` tak mengimpor HTTP/SQL/package kita yang lain.
2. SQL **hanya** di `repository_postgres.go`.
3. Handler tipis — decode/validate → service → `httpx`. Nol logika bisnis.
4. Satu module **tak pernah** meng-import package module lain. Butuh module lain
   → definisikan interface di module **peminta**, `server` yang menyuntik.
5. `platform/*` netral domain — tak boleh import `modules/*`.
6. Query selalu difilter `user_id`; `PasswordHash` tak pernah masuk response.

**Frontend:**
1. **HTTP hanya di `src/apis/`.** Halaman & komponen tak pernah panggil
   `axios`/`fetch` — selalu lewat hook di `features/*/queries/`.
2. `src/apis/schema.gen.ts` hasil generate — jangan diedit tangan.
3. Server state → TanStack Query. Token & UI → Zustand. Jangan dicampur.
4. Tanggal "hari ini" pakai yang dikirim backend, bukan `new Date()`.

**Keduanya:**
- Ubah/tambah endpoint → **ubah `contracts/openapi.yaml` dulu**.
- Keputusan desain baru → **tambah ADR**, jangan hapus/ubah yang lama.
- Ragu soal desain? **Berhenti dan tanya**, jangan tebak lalu memaksakan.

---

## Cara memulai (setelah Phase 0 dikerjakan)

```bash
# Infra
docker compose -f docker-compose.dev.yml up -d      # Postgres

# Backend
cd apps/backend
cp .env.example .env
make migrate-up
make dev                                            # air hot-reload, :8080

# Frontend
cd apps/frontend
npm install && cp .env.example .env
npm run gen:api                                     # type dari contracts/
npm run dev                                         # :5173
```

Belum bisa dijalankan sekarang — `make build` akan gagal karena belum ada kode.

---

## Jebakan yang sudah ditemukan (jangan diulang)

Ini temuan saat membaca scaffold. Semuanya sudah jadi task, tapi gampang
terlewat kalau task-nya tidak dibaca:

1. **`user/routes.go`** mendaftarkan `/auth/register`, `/auth/login`, dan `/me`
   dalam satu `RegisterRoutes`, sementara `server/router.go` memanggilnya hanya
   di grup publik → **`GET /me` akan mendarat tanpa auth**. Pecah jadi
   `RegisterPublicRoutes` / `RegisterProtectedRoutes` (backend T1.6).
2. **Nama tabrakan:** `token.go` mendefinisikan `auth.Verifier`, tapi
   `middleware.go` menulis `auth.TokenVerifier`. Dipilih `auth.Verifier`
   (backend T0.12).
3. **`Dockerfile` akan gagal build** karena `COPY go.mod go.sum ./` sedangkan
   `go.sum` belum ada. Lahir di backend T0.2.
4. **`cmd/api/main.go`** mencontohkan `cfg := config.Load()` (tanpa error),
   padahal signature-nya `Load() (Config, error)`. Ikuti signature, bukan komentar.
5. **`RegisterRequest.Timezone`** ditandai `required` **dan** "default
   Asia/Jakarta" sekaligus — tak bisa dua-duanya. Dipilih opsional + default di
   service (backend T1.2).
6. **Perangkap chi:** `/quests/today` harus didaftarkan **sebelum**
   `/quests/{questId}`, kalau tidak `today` tertangkap sebagai id.
7. **`ErrNotOwner` dipetakan ke 404**, bukan 403 — jangan bocorkan bahwa quest
   itu ada milik orang lain.
8. **Logout frontend wajib `queryClient.clear()`** — tanpa itu user berikutnya di
   browser yang sama melihat data user sebelumnya dari cache.

---

## Yang BELUM diputuskan

Ini terbuka dan perlu keputusan saat implementasi menyentuhnya:

- **ADR-016 — atomicity "buat log + tambah poin".** Dua tulisan DB berurutan;
  kalau yang kedua gagal, log terlanjur ada dan poin tak masuk. Pilihannya
  dikerjakan (transaksi lintas-module) atau ditunda (kompensasi manual). **Wajib
  ditulis sebagai ADR** apa pun yang dipilih. Lihat backend T3.10.
- **Mapper error domain → HTTP**: registry di `httpx` (yang dipanggil tiap module
  saat `New()`) atau tiap handler memetakan sendiri. Backend T0.8 — kalau memilih
  registry, catat sebagai ADR.
- **`/leaderboard` butuh login atau tidak.** Untuk MVP diputuskan tetap
  terproteksi; kalau dibuka publik → ADR baru.
- **Stack deployment** — belum dibahas sama sekali.

---

## Batasan scope

**Jangan dikerjakan kecuali diminta eksplisit:**

- Module `achievement` dan seluruh fitur v2 (ADR-010).
- Event bus in-process, streak freeze, leaderboard paging, Redis, background job.
- Halaman statistik/grafik di frontend (butuh endpoint riwayat yang belum ada).

Backlog v2 lengkap ada di `apps/backend/docs/tasks/README.md` dan
`apps/frontend/docs/tasks/README.md`.

**Catatan sikap (dari AGENTS.md):** pemilik kode ingin menulis implementasinya
sendiri. Scaffold sengaja berisi TODO. Jangan mengisi implementasi besar-besaran
tanpa diminta — kalau diminta, kerjakan sepotong dan jelaskan.

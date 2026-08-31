# Task Backend QuestDay

Peta pekerjaan `apps/backend` dari nol sampai **MVP selesai**, dibagi per phase.
Dokumen ini adalah *index* + papan progres. Detail tiap task ada di file phase.

Sumber kebenaran lain (baca dulu kalau belum): `docs/ARCHITECTURE.md` (bentuk
sistem), `docs/DECISIONS.md` (kenapa begini), `docs/GUIDES.md` (cara ngerjain),
`AGENTS.md` (aturan main).

---

## Kondisi awal (per 2026-08-31)

Backend masih **scaffold murni**: 42 file `.go` yang isinya hanya `package x` +
komentar TODO (62 marker TODO total), `go.mod` tanpa satu pun dependency,
`migrations/000001_init.up.sql` tanpa satu baris DDL, dan
`contracts/openapi.yaml` dengan 11 operasi tanpa schema (`components.schemas: {}`).

Artinya: **belum ada kode yang perlu dibongkar.** Semua task di bawah adalah
mengisi yang kosong, bukan mengubah yang sudah jalan.

---

## Peta phase & dependensi

```
Phase 0  Foundation ─────────────────────────────► blocker semua phase
   │     go.mod, migrasi, config, platform/*
   ▼
Phase 1  User & Auth ────────────────────────────► butuh: auth, httpx, validator, db
   │     modules/user + middleware Authenticator
   ▼
Phase 2  Quest ──────────────────────────────────► butuh: middleware (userID+tz)
   │     modules/quest, ScoreAwarder masih no-op
   ▼
Phase 3  Scoring ────────────────────────────────► butuh: port quest.ScoreAwarder
   │     modules/scoring + sambungkan port + UserDirectory
   ▼
Phase 4  Hardening & Kontrak ────────────────────► MVP DONE
         health, shutdown, error mapping, openapi.yaml lengkap
```

Urutan ini mengikuti `docs/GUIDES.md` §8. Jangan lompat: phase 2 sengaja bisa
diuji tanpa scoring (pakai ScoreAwarder no-op) supaya tidak ada dua module
setengah jadi sekaligus.

### Yang ditunggu `apps/frontend`

Frontend dikerjakan paralel dan memblokir di dua titik:

| Frontend butuh | Dari | Kenapa |
|---|---|---|
| `contracts/openapi.yaml` terisi | **T0.15** | frontend menggenerate seluruh type-nya dari sana (ADR-019) |
| `PATCH /me` | **T1.10** | halaman Settings — user harus bisa mengubah timezone (ADR-006, ADR-022) |

Sisanya frontend pakai mock (MSW) sampai endpoint terkait jadi. Lihat
`apps/frontend/docs/tasks/README.md`.

---

## Keputusan yang mengikat semua task

Lima keputusan diambil sebelum task ini disusun, dicatat sebagai
**ADR-011 s/d ADR-015** di `docs/DECISIONS.md`:

| ADR | Topik | Keputusan |
|-----|-------|-----------|
| 011 | Module path Go | `questday` (import: `questday/internal/...`) |
| 012 | Akses DB | `database/sql` + driver `pgx/v5/stdlib`; semua repo terima `*sql.DB` |
| 013 | Timezone user | Ikut di **JWT claims**; middleware isi `userID` + `timezone` ke context |
| 014 | Nama di leaderboard | Port **`UserDirectory`** milik `scoring`, diimplement `user` |
| 015 | Primary key | `uuid`, **UUIDv7 digenerate di aplikasi** (`github.com/google/uuid`) |
| 022 | `PATCH /me` | Masuk MVP; mengubah timezone **menerbitkan token baru** |

Kalau salah satu terasa salah saat implementasi: **jangan ditabrak diam-diam** —
tulis ADR baru yang men-supersede.

---

## Papan progres

Status: `[ ]` belum · `[~]` jalan · `[x]` selesai

### Phase 0 — Foundation → [phase-00-foundation.md](phase-00-foundation.md)

| | ID | Task |
|---|---|---|
| [ ] | T0.1 | `go.mod`: ganti module path → `questday` |
| [ ] | T0.2 | `go.mod`: tambah dependency + `make tidy` (lahirkan `go.sum`) |
| [ ] | T0.3 | Migrasi: DDL 6 tabel (`.up.sql`) |
| [ ] | T0.4 | Migrasi: `.down.sql` urutan terbalik |
| [ ] | T0.5 | `config`: `Config` + `Load()` |
| [ ] | T0.6 | `platform/database`: `Connect` / `MustConnect` |
| [ ] | T0.7 | `platform/httpx/response.go` |
| [ ] | T0.8 | `platform/httpx/errors.go` + mapper error domain → HTTP |
| [ ] | T0.9 | `platform/httpx/decode.go` |
| [ ] | T0.10 | `platform/validator` |
| [ ] | T0.11 | `platform/auth/password.go` (bcrypt) |
| [ ] | T0.12 | `platform/auth/token.go` (Claims + JWT) |
| [ ] | T0.13 | Test phase 0 |
| [x] | T0.14 | ADR-011..015 di `docs/DECISIONS.md` (sudah ditulis, tinggal ditinjau) |
| [ ] | T0.15 | **Isi `contracts/openapi.yaml`** — memblokir `apps/frontend` |

### Phase 1 — User & Auth → [phase-01-user-auth.md](phase-01-user-auth.md)

| | ID | Task |
|---|---|---|
| [ ] | T1.1 | `user/domain.go` |
| [ ] | T1.2 | `user/dto.go` |
| [ ] | T1.3 | `user/repository.go` + `repository_postgres.go` |
| [ ] | T1.4 | `user/service.go` |
| [ ] | T1.5 | `user/handler.go` |
| [ ] | T1.6 | `user/routes.go` — pecah publik vs terproteksi |
| [ ] | T1.7 | `user/module.go` |
| [ ] | T1.8 | `platform/middleware`: `Authenticator` + context helper |
| [ ] | T1.9 | `server`: rakit minimal (user saja) supaya bisa dijalankan |
| [ ] | T1.10 | **`PATCH /me`** — ubah profil & timezone, terbitkan token baru |
| [ ] | T1.11 | Test phase 1 |

### Phase 2 — Quest → [phase-02-quest.md](phase-02-quest.md)

| | ID | Task |
|---|---|---|
| [ ] | T2.1 | `quest/domain.go` — entitas, enum, `Points()`, error |
| [ ] | T2.2 | `quest/dto.go` |
| [ ] | T2.3 | Migrasi: verifikasi kolom quest & quest_logs cukup |
| [ ] | T2.4 | `quest/repository.go` + `repository_postgres.go` (10 method) |
| [ ] | T2.5 | `quest/service.go` — port `ScoreAwarder` + 7 use case |
| [ ] | T2.6 | `quest/handler.go` — hitung "hari ini" dari timezone user |
| [ ] | T2.7 | `quest/routes.go` |
| [ ] | T2.8 | `quest/module.go` |
| [ ] | T2.9 | `server`: mount quest + suntik ScoreAwarder no-op |
| [ ] | T2.10 | Test phase 2 |

### Phase 3 — Scoring → [phase-03-scoring.md](phase-03-scoring.md)

| | ID | Task |
|---|---|---|
| [ ] | T3.1 | `scoring/domain.go` — Wallet, Streak, Transaction, `LevelForXP`, `NextStreak` |
| [ ] | T3.2 | `scoring/dto.go` |
| [ ] | T3.3 | `scoring/repository.go` + `repository_postgres.go` (6 method, UPSERT) |
| [ ] | T3.4 | Port `UserDirectory` di scoring + `user.AsUserDirectory()` |
| [ ] | T3.5 | `scoring/service.go` — `OnQuestCompleted` / `OnQuestUncompleted` |
| [ ] | T3.6 | `scoring/service.go` — `GetScore`, `GetStreak`, `Leaderboard` |
| [ ] | T3.7 | `scoring/handler.go` + `routes.go` |
| [ ] | T3.8 | `scoring/module.go` — `AsScoreAwarder()` |
| [ ] | T3.9 | `server`: ganti no-op jadi scoring sungguhan |
| [ ] | T3.10 | Atomicity: log + poin dalam satu transaksi (atau ADR utang) |
| [ ] | T3.11 | Test phase 3 |

### Phase 4 — Hardening & Kontrak → [phase-04-hardening.md](phase-04-hardening.md)

| | ID | Task |
|---|---|---|
| [ ] | T4.1 | `server/health.go` — `/healthz` + `/readyz` |
| [ ] | T4.2 | `server/router.go` final |
| [ ] | T4.3 | `server/server.go` — `New/ListenAndServe/Shutdown` |
| [ ] | T4.4 | `cmd/api/main.go` — graceful shutdown |
| [ ] | T4.5 | Error mapping seragam lintas module |
| [ ] | T4.6 | Verifikasi kontrak vs implementasi (pengisiannya di T0.15) |
| [ ] | T4.7 | `Dockerfile` + `.env.example` |
| [ ] | T4.8 | Smoke test end-to-end MVP |

---

## Definition of Done — MVP backend

Semua terpenuhi:

1. `make fmt && make vet && make build && make test` bersih.
2. `make migrate-up` sukses dari database kosong.
3. Alur lengkap jalan lewat HTTP: register → login → buat quest → complete →
   `GET /me/score` naik → `GET /me/streak` benar → `GET /leaderboard` tampil →
   `PATCH /me` mengubah timezone & menerbitkan token baru.
4. `contracts/openapi.yaml` cocok dengan implementasi (tak ada endpoint yang ada
   di kode tapi hilang dari kontrak, dan sebaliknya).
5. Tak ada import lintas-module (`modules/a` → `modules/b`), tak ada SQL di luar
   `repository_postgres.go`, tak ada `time.Now()` mentah untuk logika harian.
6. `PasswordHash` tidak pernah muncul di response mana pun.
7. Keputusan baru selama implementasi sudah dicatat di `docs/DECISIONS.md`.

Cek cepat pelanggaran batas module:

```bash
cd apps/backend
grep -rn "questday/internal/modules/" internal/modules/ | grep -v "^internal/modules/\([a-z]*\)/.*questday/internal/modules/\1"
# harus kosong
```

---

## Backlog v2 (sengaja TIDAK dipecah jadi task)

Jangan dikerjakan sampai MVP selesai (AGENTS.md, ADR-010):

- **Module `achievement`** — `Badge`/`Unlock`, `GET /me/achievements`, tabel
  `badges` & `unlocks`. Scaffold sudah ada di `internal/modules/achievement/`
  tapi belum di-mount.
- **Event bus in-process** menggantikan orkestrasi langsung, supaya `scoring` &
  `achievement` tidak menumpuk di `quest.service` (ADR-005 & ADR-010).
- **Streak freeze / grace day** (ADR-008 menunda ini).
- **Rollback streak saat uncomplete** (ADR-009 menunda ini — utang yang disadari).
- **Leaderboard lanjutan**: paging, periode (mingguan/bulanan), cache.
- **Redis** (cache / rate-limit / job queue), background job.
- **Tooling**: target `lint` (golangci-lint) di Makefile, adminer/pgweb di
  `docker-compose.dev.yml`.
- **`apps/frontend`** — masih kosong, butuh keputusan stack dulu.

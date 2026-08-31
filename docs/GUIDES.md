# Panduan Kerja (How-To)

Resep praktis: kalau mau X, sentuh file apa saja. Baca `ARCHITECTURE.md` dulu
untuk gambaran besar. Setiap kali mengambil keputusan tak-sepele, catat di
`DECISIONS.md`.

Konvensi penulisan di panduan ini: **Buat** = file baru, **Ubah** = file lama
yang disesuaikan.

---

## 1. Memulai project dari nol

Prasyarat: Go 1.23+, Docker, dan CLI berikut:

```bash
go install github.com/air-verse/air@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Langkah:

1. **Ganti module path** di `apps/backend/go.mod` (`github.com/yourorg/questday`
   → repo kamu). Semua import internal mengikuti path ini.
2. Dari root: `docker compose -f docker-compose.dev.yml up -d` (nyalakan Postgres).
3. `cd apps/backend && cp .env.example .env`, sesuaikan bila perlu.
4. Tambahkan dependency awal lalu `go mod tidy`:
   chi, driver postgres (pgx/lib/pq), jwt, validator, bcrypt (lihat komentar di `go.mod`).
5. Isi implementasi sesuai TODO — urutan yang disarankan ada di bagian 8.
6. `make migrate-up` setelah menulis migrasi.
7. `make dev` untuk jalan dengan hot reload.

Cek cepat sehat: `GET /healthz` (setelah kamu implement health handler).

---

## 2. Menambah endpoint pada module yang sudah ada

Contoh: menambah `GET /quests/{id}` (detail satu quest).

1. **Ubah `contracts/openapi.yaml`** — tambah path-nya dulu (kontrak lebih dulu).
2. **Ubah `modules/quest/dto.go`** — bila butuh bentuk response baru.
3. **Ubah `modules/quest/service.go`** — tambah use case, mis. `GetQuest(ctx,userID,id)`.
4. **Ubah `modules/quest/repository.go`** — tambah method interface bila perlu query baru.
5. **Ubah `modules/quest/repository_postgres.go`** — implement query-nya.
6. **Ubah `modules/quest/handler.go`** — tambah handler yang memanggil service.
7. **Ubah `modules/quest/routes.go`** — daftarkan route ke handler.

Tak ada file yang perlu disentuh di luar module + kontrak. Itu tanda batas
module-nya sehat.

---

## 3. Menambah use case (tanpa endpoint baru)

Use case = method baru di `service.go`. Contoh: aturan bisnis internal yang
dipanggil use case lain.

1. **Ubah `service.go`** — tambah method; taruh logika di sini.
2. Kalau butuh data baru: **Ubah `repository.go`** (+ `repository_postgres.go`).
3. Kalau aturan murni domain (mis. perhitungan): taruh di **`domain.go`** sebagai
   fungsi/metode entitas, lalu panggil dari service. Jaga domain tetap tanpa
   dependensi luar.

Pedoman: HTTP → handler, orkestrasi → service, aturan murni → domain, data → repo.

---

## 4. Menambah module baru

Contoh: module `reminder`.

1. **Buat folder** `internal/modules/reminder/` dengan berkas pola standar:
   `module.go, domain.go, dto.go, repository.go, repository_postgres.go,
   service.go, handler.go, routes.go`. Salin kerangka dari module lain.
2. **Isi `module.go`** — `New(...)` merakit repo→service→handler, plus
   `RegisterRoutes`.
3. **Ubah `internal/server/server.go`** — instansiasi module & simpan dependency.
4. **Ubah `internal/server/router.go`** — `reminderMod.RegisterRoutes(r)` di
   group yang tepat (publik vs butuh-auth).
5. Butuh tabel? **Buat migrasi** (`make migrate-create name=create_reminders`).
6. **Ubah `contracts/openapi.yaml`** untuk endpoint barunya.
7. **Catat di `DECISIONS.md`** bila module ini membawa keputusan desain baru.

**Butuh module lain?** Jangan meng-import package-nya. Definisikan interface
(port) di module peminta, implementasikan di module penyedia, suntik lewat
`server`. (Pola: `quest.ScoreAwarder` ← `scoring`.)

---

## 5. Membuat dua module saling bicara (pola port)

1. Di module **peminta**, `service.go`: definisikan interface kecil berisi hanya
   yang dibutuhkan (mis. `type Notifier interface { Notify(...) error }`).
2. Simpan interface itu sebagai field service; panggil di use case terkait.
3. Di module **penyedia**: pastikan service-nya punya method yang cocok, dan
   ekspos lewat `module.go` (mis. `func (m *Module) AsNotifier() *service`).
4. Di `server/server.go`: `requesterMod := requester.New(db, providerMod.AsNotifier())`.

Aturan: interface milik peminta, bukan penyedia. Ini yang menjaga arah dependensi.

---

## 6. Menambah migrasi database

1. `make migrate-create name=deskripsi_singkat` → menghasilkan pasangan
   `NNNNNN_deskripsi.up.sql` & `.down.sql` di `migrations/`.
2. Tulis DDL di `.up.sql`; tulis kebalikannya di `.down.sql` (urutan DROP dibalik,
   hormati foreign key).
3. `make migrate-up` untuk menerapkan. `make migrate-down` untuk mundur satu.
4. Kalau state jadi `dirty` karena migrasi gagal: perbaiki SQL, lalu
   `make migrate-force version=N` ke versi yang benar, ulangi.

Jangan pernah mengedit migrasi yang sudah dipakai bersama — buat migrasi baru.

---

## 7. Menambah field pada DTO / entitas

- **DTO saja** (bentuk API): ubah `dto.go` + mapper-nya + `contracts/openapi.yaml`.
- **Entitas domain** (butuh kolom DB): ubah `domain.go`, buat migrasi,
  sesuaikan `repository_postgres.go` (query select/insert/update), lalu mapper
  di `dto.go`. Ingat: jangan bocorkan field sensitif (mis. `PasswordHash`) ke
  response.

---

## 8. Urutan implementasi MVP yang disarankan

1. `config` + `platform/database` + `platform/httpx` — pondasi.
2. `platform/auth` (jwt + bcrypt) + `platform/validator`.
3. Module **user** (register/login/me) + migrasi `users`. Uji auth dulu.
4. `platform/middleware` Authenticator, pasang di `server/router.go`.
5. Module **quest** (CRUD + today) + migrasi `quests`, `quest_logs` — tanpa
   scoring dulu (suntik ScoreAwarder no-op sementara).
6. Module **scoring** + migrasi `wallets`, `streaks`, `point_transactions`.
   Sambungkan sebagai `ScoreAwarder` sungguhan; uji complete/uncomplete.
7. Rapikan: health check, graceful shutdown, error mapping seragam.
8. (v2) **achievement**, leaderboard lanjutan, freeze streak, dll.

Setiap langkah: perbarui `contracts/openapi.yaml` seiring jalan, catat keputusan
baru di `DECISIONS.md`.

---

## 9. Perintah harian (Makefile)

```
make dev            # server + hot reload (air)
make run            # jalan sekali tanpa reload
make build          # compile ke ./bin
make test           # test + race + cover
make fmt / vet / tidy
make migrate-up / migrate-down / migrate-create name=... / migrate-force version=...
make help           # daftar lengkap
```

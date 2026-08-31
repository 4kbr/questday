# Arsitektur QuestDay

Dokumen ini menjelaskan *bentuk* sistem: lapisan, anatomi module, alur request,
dan aturan dependensi. Untuk *alasan di balik* keputusan, lihat `DECISIONS.md`.
Untuk *cara mengerjakan* sesuatu, lihat `GUIDES.md`.

## Gaya arsitektur

**Modular monolith.** Satu binary, satu proses, satu database — tapi kode
dipecah jadi module yang batasnya tegas. Tiap module memiliki domain, logika,
dan penyimpanannya sendiri. Module berkomunikasi lewat *interface (port)*, bukan
saling mengintip isi.

Kenapa monolith dulu, bukan microservice: lebih sederhana untuk MVP solo, satu
deploy, transaksi DB gampang. Batas antar-module yang rapi membuat pemecahan ke
service terpisah tetap mungkin nanti tanpa merombak total.

## Peta folder

```
apps/backend/
├── cmd/api/main.go            # entrypoint: rakit & jalankan, lalu shutdown
├── internal/
│   ├── config/                # muat konfigurasi dari env
│   ├── platform/              # kode teknis lintas-module (bukan domain)
│   │   ├── database/          # koneksi Postgres
│   │   ├── httpx/             # helper response/decode/error seragam
│   │   ├── middleware/        # middleware kustom (auth)
│   │   ├── auth/              # JWT + hashing password
│   │   └── validator/         # validasi DTO
│   ├── server/                # composition root: rakit router + module
│   └── modules/               # DOMAIN — inti aplikasi
│       ├── user/              # auth & profil (menyimpan timezone)
│       ├── quest/             # definisi quest + log harian
│       ├── scoring/           # poin, XP, level, streak, leaderboard
│       └── achievement/       # badge (v2, post-MVP)
└── migrations/                # skema DB (golang-migrate)
```

## Lapisan di dalam satu module

Tiap module (contoh `quest`) memakai lapisan yang sama:

```
handler.go   HTTP: decode/validate request, panggil service, tulis response
   │            (tipis; tak ada logika bisnis)
   ▼
service.go   Use case / logika aplikasi. Orkestrasi domain + repository.
   │            (tak tahu HTTP maupun SQL)
   ▼
repository.go        Kontrak penyimpanan (interface)
repository_postgres.go   Implementasi SQL (satu-satunya tempat SQL module ini)

domain.go    Entitas, enum, aturan invarian, error domain (Go murni)
dto.go       Bentuk data HTTP (request/response), terpisah dari entitas
module.go    Perakitan module + RegisterRoutes (satu-satunya API publik module)
routes.go    Pemetaan URL -> handler
```

## Aturan dependensi (penting)

Panah = "boleh bergantung pada". Yang tidak digambar = tidak boleh.

```
handler  ->  service  ->  repository(interface)  <-  repository_postgres
                │
                └────────────►  domain
```

Aturan keras:

1. **Domain tidak bergantung ke apa pun.** `domain.go` tak mengimpor HTTP/SQL/
   package lain milik kita. Ini inti yang stabil.
2. **Arah dalam→luar.** Domain ← service ← handler. Tidak boleh terbalik.
3. **SQL hanya di `repository_postgres.go`.** Service & handler tak menyentuh SQL.
4. **Antar-module lewat port.** Satu module TIDAK meng-import package module lain.
   Kalau butuh module lain, definisikan interface (port) di module peminta, dan
   `server` yang menyuntik implementasinya. Lihat pola `quest.ScoreAwarder`.
5. **`platform/*` tidak mengimpor `modules/*`.** Platform itu netral domain.
6. **`server` boleh tahu semua module** — dialah composition root.

## Contoh keterhubungan antar-module (gamifikasi)

Saat quest diselesaikan, poin harus bertambah. Tapi `quest` tidak boleh
bergantung ke `scoring`. Solusinya port:

```
quest.service  ──(panggil)──►  quest.ScoreAwarder (interface, milik quest)
                                        ▲
                                        │ diimplementasi oleh
                                 scoring.service

server.New():  quest.New(db, scoringMod.AsScoreAwarder())
```

Jadi `quest` hanya kenal interface-nya sendiri; `scoring` yang menyesuaikan diri;
`server` yang menjahit keduanya.

## Alur request (contoh: selesaikan quest)

```
POST /api/v1/quests/{id}/complete
  │
  ▼ chi router (server/router.go) + middleware auth (isi userID ke context)
  ▼ quest.handler.complete   -> ambil userID & timezone, hitung "hari ini"
  ▼ quest.service.CompleteQuest -> cek kepemilikan, buat QuestLog
  │      └─► quest.ScoreAwarder.OnQuestCompleted  (= scoring.service)
  │              -> tambah poin, hitung level, update streak, catat transaksi
  ▼ httpx tulis response seragam
```

Idealnya pembuatan log + penambahan poin berada dalam satu transaksi DB agar
konsisten (lihat DECISIONS soal transaksi lintas-module).

## Batas hari & timezone

"Quest per hari" berarti butuh definisi *hari* yang jelas. Hari dihitung dari
**timezone user** (disimpan di `users.timezone`), bukan timezone server. Handler
mengubah waktu sekarang ke tanggal lokal user sebelum memanggil service. Detail
& konsekuensinya di `DECISIONS.md`.

## Yang belum ada (sadar, sengaja)

- Event bus in-process (sekarang pakai orkestrasi langsung; lihat DECISIONS).
- Cache/Redis, background job.
- Module achievement (v2).

Semua itu titik ekstensi yang sudah diantisipasi struktur, bukan utang yang
memaksa perombakan.

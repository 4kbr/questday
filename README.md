# QuestDay

Pencatat quest/task harian dengan konsep gamifikasi — selesaikan quest (olahraga,
lari, tidur cukup, belajar, ngoding, dll), kumpulkan poin, naikkan XP/level, dan
jaga streak.

Monorepo. Backend berupa **modular monolith** (Go + chi). Frontend menyusul.

## Struktur

```
questday/
├── apps/
│   ├── backend/          # Go modular monolith (fokus MVP saat ini)
│   └── frontend/         # TODO — masih kosong
├── contracts/            # Kontrak API (OpenAPI) — sumber kebenaran bersama app
├── docs/                 # HANDS-OFF, ARCHITECTURE, DECISIONS, GUIDES
├── docker-compose.dev.yml# Infra dev (Postgres, dll)
└── AGENTS.md             # Aturan main untuk AI coding agent
```

## Mulai cepat

Lihat `docs/GUIDES.md` bagian "Memulai project dari nol".

Ringkas:

```bash
docker compose -f docker-compose.dev.yml up -d   # nyalakan Postgres
cd apps/backend
cp .env.example .env                              # isi kredensial
make migrate-up                                   # jalankan migrasi
make dev                                          # jalankan server + hot reload (air)
```

## Dokumentasi

- `docs/HANDS-OFF.md` — **state proyek saat ini** & titik mulai (baca duluan).
- `docs/ARCHITECTURE.md` — bentuk sistem, lapisan, aturan dependensi.
- `docs/DECISIONS.md` — catatan keputusan (ADR) selama pengembangan.
- `docs/GUIDES.md` — cara nambah module / endpoint / usecase / migrasi / dll.
- `AGENTS.md` — panduan & batasan untuk agent yang bantu ngoding.

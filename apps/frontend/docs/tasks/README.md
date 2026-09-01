# Task Frontend QuestDay

Peta pekerjaan `apps/frontend` dari nol sampai **MVP selesai**, dibagi per phase.
Dokumen ini adalah *index* + papan progres. Detail tiap task ada di file phase.

Baca dulu kalau belum: `contracts/openapi.yaml` (kontrak API — sumber kebenaran
bersama), `docs/DECISIONS.md` (kenapa begini), `AGENTS.md` (aturan main),
dan `apps/backend/docs/tasks/README.md` (apa yang backend sediakan & kapan).

---

## Kondisi awal

`apps/frontend` **kosong total** — hanya `README.md`. Belum ada `package.json`,
belum ada satu baris kode. Semua task di bawah membangun dari nol.

Backend juga belum ada kodenya (baru task list). Karena itu **MSW** dipakai
supaya frontend tidak menganggur menunggu backend — lihat ADR-021.

---

## Stack

| Bagian | Pilihan | Kenapa |
|---|---|---|
| Core | React + Vite + TypeScript | ADR-017 |
| Server state | TanStack Query | cache, refetch, invalidate lintas-fitur — ADR-018 |
| Client state | Zustand | hanya token auth & UI state — ADR-018 |
| HTTP | axios, terkumpul di `src/apis/*.api.ts` | halaman tak pernah panggil HTTP langsung |
| Types | `openapi-typescript` dari kontrak | ADR-019, menegakkan ADR-002 |
| Routing | React Router v7 | ekosistem paling matang |
| UI | shadcn/ui + Tailwind CSS | tampilan ala SaaS |
| Form | react-hook-form + zod | pasangan bawaan shadcn `Form` |
| Token | localStorage (zustand persist) | ADR-020 — backend belum punya refresh token |
| Mock dev | MSW di balik `VITE_USE_MOCK` | ADR-021 |

---

## Peta phase & dependensi

```
Phase 0  Setup & fondasi ────────────────────────► blocker semua phase
   │     vite, tailwind, shadcn, apis/client.ts, query client, MSW
   ▼
Phase 1  Auth & shell ───────────────────────────► butuh: client.ts, router
   │     login, register, auth.store, ProtectedRoute, AppShell
   ▼
Phase 2  Quest ──────────────────────────────────► butuh: shell + token
   │     dashboard (today), CRUD quest, complete/uncomplete
   ▼
Phase 3  Scoring ────────────────────────────────► butuh: questKeys (Phase 2)
   │     kartu score/streak, leaderboard, invalidasi lintas-fitur
   ▼
Phase 4  Settings & polish ──────────────────────► MVP DONE
         PATCH /me, dark mode, error/empty state, build
```

### Ketergantungan ke backend

Frontend **tidak menunggu backend selesai**, tapi ada dua titik sambung:

| Butuh | Dari backend | Kapan |
|---|---|---|
| `contracts/openapi.yaml` terisi schema | **T0.15** (Phase 0 backend) | sebelum FE T0.4 (`gen:api`) |
| Endpoint auth jalan | Phase 1 backend | sebelum FE Phase 1 lepas dari mock |
| Endpoint quest jalan | Phase 2 backend | sebelum FE Phase 2 lepas dari mock |
| Endpoint scoring jalan | Phase 3 backend | sebelum FE Phase 3 lepas dari mock |
| `PATCH /me` | **T1.11** (Phase 1 backend) | sebelum FE Phase 4 |

Selama endpoint terkait belum ada, kerjakan dengan `VITE_USE_MOCK=true`.

---

## Struktur `src/`

```
src/
├── apis/              # SATU-SATUNYA tempat panggil HTTP
│   ├── client.ts      # axios instance + interceptor (token, 401, unwrap envelope)
│   ├── schema.gen.ts  # HASIL GENERATE — jangan diedit tangan
│   ├── types.ts       # alias ramah dari schema.gen.ts
│   ├── auth.api.ts
│   ├── quest.api.ts
│   └── scoring.api.ts
├── features/
│   ├── auth/{components,queries,schemas}
│   ├── quest/{components,queries,schemas}
│   └── scoring/{components,queries}
├── components/
│   ├── ui/            # shadcn (generated)
│   └── layout/        # AppShell, Sidebar, Topbar
├── pages/             # satu file per route, tipis
├── routes/            # router, ProtectedRoute, konstanta path
├── stores/            # auth.store.ts, ui.store.ts
├── lib/               # utils, date, query-client.ts
└── mocks/             # MSW handlers + setup
```

---

## Aturan keras (berlaku di semua phase)

1. **HTTP hanya di `src/apis/`.** Halaman & komponen tak pernah memanggil
   `axios`/`fetch` langsung — selalu lewat hook di `features/*/queries/`, yang
   memanggil fungsi di `src/apis/*.api.ts`. Ini permintaan eksplisit pemilik kode
   dan dikunci ADR-018.
2. **`schema.gen.ts` tak pernah diedit tangan.** Kontrak berubah → jalankan
   `npm run gen:api`.
3. **Tak ada type request/response yang diketik ulang.** Semua turun dari
   `schema.gen.ts` lewat alias di `apis/types.ts`.
4. **Server state ≠ client state.** Data dari API → TanStack Query. Token & UI →
   Zustand. Jangan menyimpan hasil API ke zustand.
5. **"Hari ini" ikut timezone user** (ADR-006), bukan `new Date()` browser.
   Backend yang menentukan tanggal; frontend menampilkan apa yang dikirim
   `/quests/today`.
6. **Token hanya disentuh** lewat `auth.store` + interceptor di `client.ts`.
   Tak ada komponen yang membaca `localStorage` langsung.

---

## Papan progres

Status: `[ ]` belum · `[~]` jalan · `[x]` selesai

### Phase 0 — Setup & fondasi → [phase-00-setup.md](phase-00-setup.md)

| | ID | Task |
|---|---|---|
| [x] | F0.1 | Scaffold Vite + React + TypeScript |
| [x] | F0.2 | Tailwind CSS + path alias `@/` |
| [x] | F0.3 | shadcn/ui init + token tema SaaS |
| [x] | F0.4 | `gen:api` — generate `schema.gen.ts` dari kontrak |
| [x] | F0.5 | `apis/types.ts` — alias type ramah |
| [x] | F0.6 | `apis/client.ts` — axios + interceptor |
| [x] | F0.7 | `lib/query-client.ts` + provider |
| [x] | F0.8 | Struktur folder + `routes/` kerangka |
| [x] | F0.9 | MSW setup + handler dasar |
| [x] | F0.10 | `.env.example`, `.gitignore`, oxlint + Prettier |
| [x] | F0.11 | Script `package.json` + README frontend |
| [x] | F0.12 | Test phase 0 (Vitest + Testing Library) |

### Phase 1 — Auth & shell → [phase-01-auth.md](phase-01-auth.md)

| | ID | Task |
|---|---|---|
| [x] | F1.1 | `apis/auth.api.ts` |
| [x] | F1.2 | `stores/auth.store.ts` (zustand + persist) |
| [x] | F1.3 | `features/auth/schemas` — zod, cocok dgn validasi backend |
| [x] | F1.4 | `features/auth/queries` — `useLogin`, `useRegister`, `useMe` |
| [x] | F1.5 | Halaman Login |
| [x] | F1.6 | Halaman Register |
| [x] | F1.7 | `routes/` — `ProtectedRoute`, `GuestRoute`, konstanta path |
| [x] | F1.8 | `components/layout/AppShell` — sidebar + topbar |
| [x] | F1.9 | Logout: bersihkan token **dan** cache Query |
| [x] | F1.10 | MSW handler auth |
| [x] | F1.11 | Test phase 1 |

### Phase 2 — Quest → [phase-02-quest.md](phase-02-quest.md)

| | ID | Task |
|---|---|---|
| [ ] | F2.1 | `apis/quest.api.ts` (7 fungsi) |
| [ ] | F2.2 | `questKeys` — query key terpusat |
| [ ] | F2.3 | `features/quest/queries` — 7 hook |
| [ ] | F2.4 | `features/quest/schemas` — zod create/update |
| [ ] | F2.5 | `QuestItem` — baris quest + checkbox complete |
| [ ] | F2.6 | Optimistic update complete/uncomplete + rollback |
| [ ] | F2.7 | Halaman Dashboard — quest hari ini |
| [ ] | F2.8 | `QuestFormDialog` — create & edit |
| [ ] | F2.9 | Halaman Quests — daftar + arsip |
| [ ] | F2.10 | Empty / loading / error state |
| [ ] | F2.11 | MSW handler quest |
| [ ] | F2.12 | Test phase 2 |

### Phase 3 — Scoring → [phase-03-scoring.md](phase-03-scoring.md)

| | ID | Task |
|---|---|---|
| [ ] | F3.1 | `apis/scoring.api.ts` |
| [ ] | F3.2 | `scoringKeys` + `features/scoring/queries` |
| [ ] | F3.3 | `ScoreCard` — poin, XP, level, progress bar |
| [ ] | F3.4 | `StreakCard` — current & longest |
| [ ] | F3.5 | Pasang kedua kartu di Dashboard |
| [ ] | F3.6 | **Invalidasi lintas-fitur** setelah complete/uncomplete |
| [ ] | F3.7 | Halaman Leaderboard |
| [ ] | F3.8 | MSW handler scoring |
| [ ] | F3.9 | Test phase 3 |

### Phase 4 — Settings & polish → [phase-04-polish.md](phase-04-polish.md)

| | ID | Task |
|---|---|---|
| [ ] | F4.1 | `updateMe` + `useUpdateProfile` (butuh backend T1.11) |
| [ ] | F4.2 | Halaman Settings + pemilih timezone IANA |
| [ ] | F4.3 | Simpan token baru setelah timezone berubah (ADR-022) |
| [ ] | F4.4 | Dark mode toggle (`ui.store`) |
| [ ] | F4.5 | Toast seragam untuk mutation |
| [ ] | F4.6 | `ErrorBoundary` + halaman 404 |
| [ ] | F4.7 | Responsive — sidebar jadi drawer di mobile |
| [ ] | F4.8 | Audit a11y ringan |
| [ ] | F4.9 | Lepas mock: uji melawan backend asli |
| [ ] | F4.10 | Build produksi + Dockerfile/nginx (opsional) |

---

## Definition of Done — MVP frontend

1. `npm run lint && npm run build && npm run test` bersih.
2. Dengan `VITE_USE_MOCK=false` dan backend asli jalan, alur penuh sukses:
   register → login → buat quest → complete → poin & streak naik → leaderboard →
   ubah timezone di settings → tetap login dengan token baru.
3. Refresh browser tidak melempar user ke login; token expired (401) melempar ke
   login dengan bersih.
4. Tak ada satu pun `axios`/`fetch` di luar `src/apis/`:
   ```bash
   grep -rn "axios\|fetch(" src/ --include=*.tsx --include=*.ts | grep -v "^src/apis/" | grep -v "^src/mocks/"
   # harus kosong
   ```
5. Tak ada type API yang diketik tangan di luar `src/apis/types.ts`.
6. Setiap halaman punya loading, empty, dan error state — tak ada layar putih.
7. Tampil benar di light & dark mode, dan di lebar 375px.
8. Keputusan baru selama implementasi tercatat di `docs/DECISIONS.md`.

---

## Backlog v2 (sengaja TIDAK dipecah jadi task)

- Halaman **achievement** (menunggu module `achievement` backend — ADR-010).
- **Grafik statistik** poin/streak harian — butuh endpoint riwayat yang belum ada.
- **Streak freeze** UI (menunggu ADR-008 ditinjau ulang).
- **Leaderboard lanjutan**: paging, filter periode.
- **Refresh token / cookie httpOnly** — menggantikan ADR-020.
- PWA & offline, i18n, E2E test (Playwright), Storybook.

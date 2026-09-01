# Hands-off — State `apps/frontend`

Dokumen orientasi cepat khusus **frontend**. Untuk state seluruh proyek lihat
`docs/HANDS-OFF.md` (root). Baca `AGENTS.md` dulu, lalu ini.

**Terakhir diperbarui:** 2026-09-01 · **Branch:** `feature/frontend/master`

---

## TL;DR

**Phase 0 (Setup & fondasi) SELESAI & terverifikasi.** Belum ada fitur berisi —
tapi seluruh lapisan pondasi (Vite + tema + lapisan API + router + MSW + test)
sudah jalan.

**Belum di-commit.** Semua file `apps/frontend/*` masih untracked; `README.md`
termodifikasi. `docs/tasks/`, `contracts/`, `AGENTS.md`, root `docs/` tak tersentuh.

**Mulai dari sini:** `apps/frontend/docs/tasks/phase-01-auth.md`, kerjakan
**F1.1** dst berurutan.

---

## Progres per phase

| Phase | File task | Status |
|---|---|---|
| 0 — Setup & fondasi | `docs/tasks/phase-00-setup.md` | ✅ **selesai** (F0.1–F0.12) |
| 1 — Auth & shell | `docs/tasks/phase-01-auth.md` | ⬜ belum mulai |
| 2 — Quest | `docs/tasks/phase-02-quest.md` | ⬜ belum mulai |
| 3 — Scoring | `docs/tasks/phase-03-scoring.md` | ⬜ belum mulai |
| 4 — Settings & polish | `docs/tasks/phase-04-polish.md` | ⬜ belum mulai |

Papan progres per-task ada di `docs/tasks/README.md`.

---

## Stack terpasang (⚠️ lebih baru dari asumsi task docs)

`create-vite` v9 men-scaffold stack yang lebih baru dari yang ditulis task docs.
Semua bekerja; batch berikutnya harus sadar ini:

| Bagian | Versi / pilihan | Catatan |
|---|---|---|
| Core | **Vite 8.2**, **React 19.2**, **TypeScript 6.0** | `@vitejs/plugin-react` 6 |
| tsconfig | split (`app`/`node`), `verbatimModuleSyntax`, **`erasableSyntaxOnly`** | TS "parameter properties" di constructor **dilarang** — pakai field eksplisit |
| Styling | **Tailwind v4** (`@tailwindcss/vite`) + **shadcn/ui** v4 | style `radix-nova`, base color `neutral`, font Inter |
| Server state | **TanStack Query 5** | `src/lib/query-client.ts` |
| Client state | **Zustand** | belum dipasang — lahir di F1.2 |
| HTTP | **axios 1.20** | hanya di `src/apis/client.ts` |
| Types | **openapi-typescript 7** | butuh `--legacy-peer-deps` (peer minta TS 5.x) |
| Routing | **react-router-dom 7** | `createBrowserRouter` |
| Lint | **oxlint** (BUKAN ESLint) | keputusan pemilik kode; `.oxlintrc.json` |
| Format | **Prettier** | format saja, bukan linting |
| Mock | **MSW 2** | di balik flag `VITE_USE_MOCK` |
| Test | **Vitest 4** + Testing Library + jsdom | |

---

## Keputusan yang mengikat (dibuat selama Phase 0 — jangan ditabrak diam-diam)

1. **Kontrak = sumber kebenaran. TIDAK ADA unwrap `{data}`.** `api` mengembalikan
   `res.data` **mentah** sesuai schema. Hanya **error** yang ber-amplop:
   `{"error":{code,message}}` → diubah jadi `ApiError(status, code, message)` di
   `src/apis/client.ts`. Menyimpang dari teks draf F0.6/F0.9 yang berasumsi ada
   `{data}`. Fitur membedakan perlakuan lewat **`ApiError.code`**
   (`already_completed`, `email_taken`, dst — lihat deskripsi `ErrorResponse` di
   `contracts/openapi.yaml`), **bukan** mencocokkan string pesan. Kalau backend
   nanti benar-benar pakai `{data}`, revisi **kontrak + `client.ts` + test**
   bersama, dan catat ADR.
2. **oxlint, bukan ESLint** (menyimpang dari F0.10, disetujui pemilik kode).
   Prettier tetap dipakai untuk format.
3. **`no-restricted-imports` mem-blokir SELURUH `@/features/*`.** oxlint tak bisa
   membatasi "fitur lain" relatif terhadap file saat ini. Konsekuensi:
   **import di dalam satu fitur WAJIB relatif (`./`)**, bukan `@/features/...`.
   Kode bersama tetap lewat `@/components/*` atau `@/lib/*`. Ada `TODO(oxlint)`
   di `.oxlintrc.json`.
4. **`.npmrc` → `legacy-peer-deps=true`** supaya `npm install` polos jalan untuk
   dev baru (konflik `openapi-typescript@7` ↔ `typescript@6`).
5. **`typecheck` = `tsc -b --noEmit`** (root tsconfig solution-style; `tsc
   --noEmit` polos lolos tanpa mengecek apa-apa).
6. **shadcn CLI di `devDependencies`.** `next-themes` + `tw-animate-css` sengaja
   tetap di `dependencies` (runtime).

Semua keputusan ini juga terangkum di `apps/frontend/README.md`.

---

## Yang sudah ada di `src/`

```
src/
├── apis/
│   ├── client.ts        # axios instance TUNGGAL + ApiError + setTokenGetter/setUnauthorizedHandler
│   ├── client.test.ts   # map error→ApiError, success raw, header Bearer
│   ├── schema.gen.ts     # GENERATE dari contracts/openapi.yaml — jangan diedit tangan
│   └── types.ts          # alias ramah (User, AuthResponse, Quest, Score, Streak, ...)
├── lib/
│   ├── query-client.ts  # QueryClient: staleTime 30s, no-retry 4xx, refetchOnWindowFocus
│   └── utils.ts          # cn() dari shadcn
├── routes/
│   ├── paths.ts         # PATHS — konstanta, jangan tulis path string literal di komponen
│   ├── index.tsx         # createBrowserRouter, 6 route → placeholder page
│   └── router.test.tsx
├── pages/                # 6 placeholder tipis (Login/Register/Dashboard/Quests/Leaderboard/Settings)
├── mocks/
│   ├── handlers.ts      # GET /healthz + helper errorBody()/errorResponse() untuk reuse
│   └── browser.ts        # setupWorker(...handlers)
├── components/
│   ├── ui/               # 16 komponen shadcn (generated — tak ada logika bisnis di sini)
│   └── layout/           # kosong (.gitkeep) — AppShell lahir di F1.8
├── features/{auth,quest,scoring}/  # kerangka folder kosong (.gitkeep)
├── stores/               # kosong (.gitkeep) — auth.store & ui.store menyusul
├── test/setup.ts         # @testing-library/jest-dom/vitest
├── main.tsx              # enableMocking() → QueryClientProvider → RouterProvider → Devtools(DEV)
├── index.css             # Tailwind v4 + token tema shadcn (:root + .dark) + Inter
└── vite-env.d.ts         # tipe VITE_API_BASE_URL, VITE_USE_MOCK
```

Config root `apps/frontend/`: `.env.example`, `.gitignore`, `.npmrc`,
`.oxlintrc.json`, `.prettierrc`, `.prettierignore`, `components.json`,
`vite.config.ts`, `vitest.config.ts`, `tsconfig*.json`.

Script: `dev`, `build`, `preview`, `lint`, `format`, `format:check`, `test`,
`test:watch`, `gen:api`, `typecheck`.

---

## Cara jalan

```bash
cd apps/frontend
npm install                 # .npmrc sudah set legacy-peer-deps
cp .env.example .env
npm run gen:api             # regenerate src/apis/schema.gen.ts dari kontrak
npm run dev                 # http://localhost:5173  (VITE_USE_MOCK=true → MSW aktif)
```

Gate verifikasi (semua hijau per 2026-09-01):
`npm run typecheck && npm run lint && npm run test && npm run format:check && npm run build`

---

## Titik sambung ke backend (belum aktif)

Frontend jalan penuh dgn `VITE_USE_MOCK=true` sampai endpoint backend siap.

| Butuh | Dari backend | Dipakai FE saat |
|---|---|---|
| `contracts/openapi.yaml` terisi schema | ✅ sudah (16 schema) | F0.4 — **tidak lagi memblokir** |
| Endpoint auth jalan | Phase 1 backend | FE Phase 1 lepas mock |
| Endpoint quest jalan | Phase 2 backend | FE Phase 2 lepas mock |
| Endpoint scoring jalan | Phase 3 backend | FE Phase 3 lepas mock |
| `PATCH /me` | backend T1.11 | FE Phase 4 |

> Catatan: root `docs/HANDS-OFF.md` masih bilang `components.schemas: {}` — itu
> **stale**. Kontrak sudah terisi; `npm run gen:api` menghasilkan 16 schema.

---

## Pesan untuk agent berikutnya (Phase 1)

- **F1.2 `auth.store`**: zustand `persist` **WAJIB** `name: 'questday.auth'` dan
  simpan token di `state.token` — default `tokenGetter` di `client.ts` membaca
  key itu. Setelah store jadi, panggil `setTokenGetter(() => useAuth.getState().token)`
  dan `setUnauthorizedHandler(() => { logout(); redirect ke PATHS.login })` sekali
  saat bootstrap (mis. di `main.tsx` atau modul store).
- **F1.9 Logout**: wajib `queryClient.clear()` selain hapus token (jebakan #8 di
  root HANDS-OFF).
- HTTP baru → tambah fungsi di `src/apis/*.api.ts` (buat file baru), bungkus
  dengan hook di `src/features/*/queries/`. Komponen **tak pernah** sentuh `api`
  langsung.
- Route baru → tambah di `src/routes/paths.ts` + `src/routes/index.tsx`. Jangan
  string literal.
- Type API kurang → perbaiki `contracts/openapi.yaml`, jalankan `npm run gen:api`,
  tambah alias di `src/apis/types.ts`. Jangan bikin `interface` tangan.
- Keputusan desain baru → tambah ADR di root `docs/DECISIONS.md`, dan perbarui
  dokumen ini.

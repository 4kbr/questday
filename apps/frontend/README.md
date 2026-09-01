# QuestDay — Frontend

Web app QuestDay. Mengonsumsi API dari `apps/backend` sesuai kontrak di
`contracts/openapi.yaml` (sumber kebenaran bersama).

## Stack

- **Vite 8** + **React 19** + **TypeScript 6**
- **Tailwind CSS v4** (`@tailwindcss/vite`) + **shadcn/ui**
- **TanStack Query** — server state (cache, refetch, invalidate lintas-fitur)
- **Zustand** — client state (token auth & UI saja) — _menyusul di Phase 1_
- **React Router v7** — routing
- **react-hook-form + zod** — form
- **oxlint** (lint) + **Prettier** (format saja — bukan ESLint)
- **MSW** — mock backend di dev, di balik flag `VITE_USE_MOCK`
- **Vitest** + Testing Library + jsdom — test

## Mulai cepat

```bash
cd apps/frontend
npm install            # .npmrc memaksa legacy-peer-deps (openapi-typescript masih minta TS 5.x)
cp .env.example .env
npm run gen:api        # generate src/apis/schema.gen.ts dari contracts/openapi.yaml
npm run dev            # http://localhost:5173
```

## Flag `VITE_USE_MOCK`

Di `.env`:

| Nilai   | Efek                                                                                                                                         |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `true`  | **MSW aktif**. Request di-intercept service worker, dijawab handler di `src/mocks/handlers.ts`. Dipakai selama endpoint backend belum jalan. |
| `false` | MSW mati. Request menembus ke backend asli di `VITE_API_BASE_URL`.                                                                           |

MSW **tak pernah ikut ke production build** — `enableMocking()` di `src/main.tsx`
langsung `return` kalau flag ≠ `'true'`, dan bundler membuang dynamic-import
`@/mocks/browser` (tidak ada jejak `msw` di `dist/assets`).

## Catatan envelope

Per **ADR-025**, setiap response **sukses** berbadan JSON dibungkus
`{"data": <payload>}`, dan **`src/apis/client.ts` meng-unwrap-nya secara
terpusat** di interceptor response. Jadi `api.get(...)` / `*.api.ts` selalu
melihat `res.data` = payload langsung sesuai schema, tanpa perlu menyentuh
amplop. `register` mengembalikan **200** (bukan 201).

**Error** tetap ber-amplop `{"error":{code,message}}` → dipetakan jadi `ApiError`
(`status`, `code`, `message`) di `src/apis/client.ts`. Fitur membedakan perlakuan
lewat `ApiError.code` (mis. `already_completed`), jangan mencocokkan string
pesan.

`/healthz` **di luar amplop** — response polos (`{status:'ok'}`) dibiarkan apa
adanya oleh interceptor (tak ada key `data` → passthrough).

`src/apis/client.ts` adalah **satu-satunya** tempat membuat instance HTTP.
Halaman & komponen tak pernah memanggil `axios`/`fetch` langsung — selalu lewat
hook di `features/*/queries/` yang memanggil `src/apis/*.api.ts`.

## Build produksi

```bash
npm run build     # tsc -b && vite build → dist/ (file statis siap deploy)
npm run preview   # sajikan dist/ lokal untuk verifikasi (default port dari vite.config)
```

- **`VITE_*` di-_inline_ saat build**, bukan dibaca saat runtime. `VITE_API_BASE_URL`
  (dan `VITE_USE_MOCK`, `VITE_PORT`) ditentukan **sebelum `npm run build`** —
  nilainya ikut ter-_bake_ ke dalam bundle/image. Mengganti env di
  container/server **tidak** berpengaruh; kalau target API berbeda, build ulang.
- **Route code-splitting**: tiap halaman jadi chunk terpisah (`React.lazy` +
  `Suspense`), jadi bundle awal kecil dan halaman dimuat sesuai kebutuhan.
- **Deep link** (mis. reload di `/quests`) bekerja karena router punya route
  `*` (catch-all). Untuk static hosting sungguhan, server tetap butuh **SPA
  fallback rewrite** (mis. nginx `try_files $uri /index.html;`) supaya semua
  path mengembalikan `index.html`. `npm run preview` sudah melakukan fallback ini
  sendiri.

## Script

| Script                 | Fungsi                                         |
| ---------------------- | ---------------------------------------------- |
| `npm run dev`          | Vite dev server                                |
| `npm run build`        | `tsc -b && vite build`                         |
| `npm run preview`      | Preview hasil build                            |
| `npm run lint`         | oxlint                                         |
| `npm run format`       | `prettier --write .`                           |
| `npm run format:check` | `prettier --check .`                           |
| `npm run test`         | `vitest run`                                   |
| `npm run test:watch`   | `vitest` (watch)                               |
| `npm run typecheck`    | `tsc -b --noEmit`                              |
| `npm run gen:api`      | generate `src/apis/schema.gen.ts` dari kontrak |

## Aturan yang tak boleh dilanggar

1. **HTTP hanya di `src/apis/`.** (mock di `src/mocks/` hanya menyebut path, tak
   membuat instance axios.)
2. **`src/apis/schema.gen.ts` hasil generate** — jangan diedit tangan. Kontrak
   berubah → `npm run gen:api`.
3. **Tak ada type API diketik ulang di luar `src/apis/types.ts`.**
4. **Server state ≠ client state.** Data API → TanStack Query. Token & UI →
   Zustand.
5. **"Hari ini" ikut timezone user** (ADR-006) — pakai tanggal dari response
   backend, bukan `new Date()`.
6. **Import silang antar-fitur dilarang** (`.oxlintrc.json` →
   `no-restricted-imports`). Kode bersama lewat `@/components/*` atau `@/lib/*`;
   di dalam satu fitur pakai import relatif.

## Struktur & rencana kerja

Struktur `src/` dan peta pekerjaan per phase (Phase 0–4, sampai MVP) ada di
[`docs/tasks/`](docs/tasks/README.md). Alasan tiap pilihan stack dicatat sebagai
ADR-017..022 di `docs/DECISIONS.md`.
